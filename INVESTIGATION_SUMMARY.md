# Investigation Summary: Import Elision Bug in tsgo

## Issue
`tsgo` incorrectly emits `require()` statements when an ES6 import contains both:
1. A `type` import specifier (e.g., `type ValueData`)
2. A regular import specifier that is ONLY used in type positions (e.g., `Value`)

## Reproduction
```typescript
// provider.ts
export class Value { data: string = ""; }
export type ValueData = { data: string };

// broken.ts - BUG
import { Value, type ValueData } from "./provider";
export function test(value: Value): Value { return null as any; }
// Output: require("./provider"); ❌ INCORRECT

// working.ts - OK
import { Value } from "./provider";
export function test(value: Value): Value { return null as any; }
// Output: No require() ✓ CORRECT
```

## Key Findings

### Trigger Condition
The bug ONLY occurs when:
- Import statement has BOTH `type` and non-`type` specifiers
- The non-`type` specifier is used ONLY in type positions (not at runtime)

### Root Cause
The import processing logic incorrectly marks the non-type import's `aliasLinks.referenced` 
as `true` when it should remain `false`. This happens during import declaration checking,
where the code fails to distinguish between:
- **Type-level references**: Should NOT mark alias as referenced
- **Value-level references**: Should mark alias as referenced

### Code Flow
1. Import specifier processed by `ImportElisionTransformer.visit` (importelision.go:79-84)
2. Calls `shouldEmitAliasDeclaration` which checks `isReferencedAliasDeclaration`
3. `isReferencedAliasDeclaration` returns `true` because `aliasLinks.referenced` is `true`
4. Import specifier is NOT elided, causing `require()` to be emitted

### Missing Logic
Comparison with TypeScript shows that tsgo's `markSymbolOfAliasDeclarationIfTypeOnly`
(checker.go:14701) is missing logic to:
1. Check if the alias target resolves to a type-only declaration
2. Properly handle mixed type/non-type imports in the same statement

TypeScript has `markSymbolOfAliasDeclarationIfTypeOnlyWorker` that checks the target
symbol's declarations for type-only status, which tsgo is missing.

## Investigation Files
- `/INVESTIGATION.md` - Detailed investigation notes
- `/testdata/repro-import-elision-bug/` - Minimal reproduction case
- `/testdata/tests/cases/compiler/importValueOnlyUsedAsType.ts` - Original test case

## Next Steps for Fix
1. Enhance `markSymbolOfAliasDeclarationIfTypeOnly` to check alias targets
2. Ensure type-level references don't set `aliasLinks.referenced = true`
3. Add logic to handle mixed type/non-type imports correctly
4. Ensure the fix doesn't break existing functionality
5. Run all tests to verify no regressions

## Files to Modify
- `/internal/checker/checker.go` - `markSymbolOfAliasDeclarationIfTypeOnly` and related functions
- Possibly `/internal/checker/emitresolver.go` - `IsReferencedAliasDeclaration`
- Possibly `/internal/transformers/tstransforms/importelision.go` - Import elision logic

## References
- TypeScript source: `src/compiler/checker.ts` lines 4385-4450
- Test case: `testdata/tests/cases/compiler/importValueOnlyUsedAsType.ts`
- Repro case: `testdata/repro-import-elision-bug/`
