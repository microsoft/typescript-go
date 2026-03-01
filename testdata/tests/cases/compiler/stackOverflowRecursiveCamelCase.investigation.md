# Investigation: Stack Overflow in Type Instantiation (#1093)

## Summary

A stack overflow occurs in tsgo when checking code that uses recursive generic functions
with complex conditional/mapped types (specifically from type-fest's `CamelCase` type).
The crash is a Go goroutine stack overflow (exceeds 1GB limit), not a graceful error.

## Reproduction

The test case in `testdata/tests/cases/compiler/stackOverflowRecursiveCamelCase.ts`
reproduces the crash. It requires **multiple files** (using `// @filename:` directives) —
the same types inlined into a single file do NOT crash.

The crash can also be reproduced with:
```bash
npm install type-fest@4.38.0
echo 'import { CamelCase } from "type-fest";
const transform = <TResult>(iteratee: any): TResult => undefined as any;
const camelize = <T extends object>(): CamelCase<T> => transform(camelize);' > test.ts
tsgo --noEmit --strict --moduleResolution bundler test.ts
```

## Root Cause Analysis

### The Recursive Cycle

The stack overflow is caused by an unbounded recursive cycle between **type relation checking**
and **variance computation**:

1. `checkReturnExpression` triggers type comparison (`checkTypeRelatedToAndOptionallyElaborate`)
2. `structuredTypeRelatedToWorker` encounters alias types with type arguments and calls
   `getAliasVariances()` to compute variance (line 3339 of relater.go)
3. `getVariancesWorker` calls `isTypeAssignableTo()` to compare marker type instantiations
   (line 1364 of relater.go)
4. `isTypeAssignableTo` → `isTypeRelatedTo` → `checkTypeRelatedToEx` creates a **fresh Relater**
   with empty source/target stacks
5. The new Relater's `recursiveTypeRelatedTo` → `structuredTypeRelatedTo` →
   `structuredTypeRelatedToWorker` encounters more alias types...
6. → Back to step 2 with a different alias symbol

### Why the Existing Guards Don't Help

**`instantiationDepth` (limit: 100):** The instantiation depth counter goes up during type
instantiation but comes back down when each instantiation chain completes. The individual
instantiation chains triggered within each variance computation are short (depth ~5-10).
The counter resets between variance computation cycles.

**`sourceStack`/`targetStack` (limit: 100):** These stacks are per-Relater instance. Each call
to `checkTypeRelatedToEx` creates a **new Relater** (via `getRelater()`), which has fresh empty
stacks. So the depth within any single Relater never reaches 100.

**`relationCount` (budget: ~16M):** This tracks total comparisons but doesn't prevent stack
overflow — it's a complexity guard, not a stack depth guard.

### Why It Works in TypeScript (JavaScript)

In JavaScript, the call stack has a hard limit (typically ~10,000-30,000 frames). This acts
as an implicit recursion guard that catches the unbounded cycle before it becomes a problem.
Go goroutines can grow their stacks up to 1GB (~5-20 million frames), so this implicit
protection doesn't exist.

### Why Multi-File is Required

When all types are in a single file, the `getOuterTypeParameters` function (line 23219 of
checker.go) may compute different outer type parameter sets, which can cause some
instantiations to short-circuit (via the `isTypeParameterPossiblyReferenced` filter in
`getObjectTypeInstantiation`). In multi-file scenarios, imported type aliases may track
more type parameters through the import chain, leading to more instantiation work and
triggering the variance computation cycle.

### Why type-fest v4.38.0 Specifically

The crash was introduced by the addition of `ApplyDefaultOptions` in type-fest v4.38.0
(PR #1081). This type alias uses complex generic machinery (involving `Simplify`, `Merge`,
`RequiredKeysOf`, `OptionalKeysOf`, `IfAny`, `IfNever`) that creates enough distinct alias
types with type arguments to trigger the variance computation cycle. The `CamelCase` type
uses `ApplyDefaultOptions` in two places (for both `Words` and `CamelCaseFromArray` options).

## Suggested Fixes

1. **Add a global type-checking depth counter** that tracks the total nesting depth across
   all Relater instances and `instantiateType` calls. This would be analogous to JavaScript's
   implicit call stack limit.

2. **Track Relater nesting depth** in the Checker: increment a counter when entering
   `checkTypeRelatedToEx` and decrement when leaving. If the counter exceeds a threshold
   (e.g., 100 or 200), set overflow and return early.

3. **Track goroutine stack usage** using Go runtime introspection (e.g., comparing current
   stack pointer against a baseline). This is more complex but directly addresses the
   underlying issue.

Option 2 is likely the most straightforward fix, as it directly addresses the unbounded
nesting of Relater instances through variance computation.

## Stack Trace Pattern

The repeating pattern in the stack trace is:
```
getObjectTypeInstantiation → instantiateTypeWorker → instantiateTypeWithAlias → 
  instantiateType → instantiateList → instantiateTypes → 
  instantiateTypeWorker → instantiateTypeWithAlias → instantiateType →
  CompositeTypeMapper.Map → ... (repeat object/conditional instantiation)
  
→ structuredTypeRelatedTo → recursiveTypeRelatedTo → isRelatedToEx →
  structuredTypeRelatedToWorker → getAliasVariances → getVariancesWorker →
  isTypeAssignableTo → isTypeRelatedTo → checkTypeRelatedToEx →
  (creates new Relater) → isRelatedToEx → recursiveTypeRelatedTo →
  structuredTypeRelatedTo → ... (repeat)
```
