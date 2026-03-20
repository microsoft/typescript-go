# Investigation: Stack Overflow with Self-Referential ReturnType (TypeScript#63273)

## Issue
TypeScript crashes with "Maximum call stack size exceeded" (or Go stack overflow) when compiling:

```typescript
function clone(): <T>(obj: T) => T extends any ? ReturnType<typeof clone>[0] : 0;
```

This affects both the TypeScript compiler (microsoft/TypeScript#63273) and the Go port (typescript-go).

## Reproduction
- **Reproduced in Go**: Yes — running `tsgo test.ts` produces `fatal error: stack overflow` (goroutine stack exceeds 1000000000-byte limit)
- **Reproduced in TypeScript 5.9.3**: Yes — `tsc test.ts` produces `RangeError: Maximum call stack size exceeded`
- **Test case added**: `testdata/tests/cases/compiler/stackOverflowSelfReferentialReturnType.ts`

## Root Cause Analysis

### The Infinite Recursion Cycle

The recursion involves a cycle between the **NodeBuilder** (type-to-node serializer) and the **Checker** (type resolver):

1. **`conditionalTypeToTypeNode`** (`nodebuilderimpl.go:2601`) — serializing the conditional type `T extends any ? ReturnType<typeof clone>[0] : 0`
2. → calls **`getTrueTypeFromConditionalType`** (`checker.go:23864`) to resolve the true branch type
3. → calls **`getTypeFromTypeNode`** → **`getTypeFromIndexedAccessTypeNode`** (`checker.go:22348`) to resolve `ReturnType<typeof clone>[0]`
4. → calls **`getIndexedAccessTypeEx`** → **`getPropertyTypeForIndexType`** (`checker.go:26434`)
5. → Property `0` doesn't exist on the resolved type, so an **error** is generated
6. → Error message calls **`c.TypeToString(objectType)`** to produce `"Property '0' does not exist on type 'X'"`
7. → **`typeToStringEx`** (`printer.go:192`) creates a **brand new** `NodeBuilder` with fresh context
8. → New NodeBuilder calls **`typeToTypeNode`** for the same anonymous function type
9. → **`createAnonymousTypeNode`** → **`visitAndTransformType`** → **`createTypeNodeFromObjectType`**
10. → Serializes the function signature, including its return type (which is the same conditional type)
11. → **GOTO step 1** — infinite loop!

### Why Existing Guards Don't Work

The Go codebase has recursion guards in the NodeBuilder:
- **`visitedTypes` set** — tracks type IDs being serialized to detect cycles
- **`symbolDepth` map** — limits depth to 10 per composite symbol identity
- **Truncation length** — soft limit based on approximate output length

**None of these guards work** because `TypeToString` at step 7 creates a **brand new** `NodeBuilder` with an empty context:
- `visitedTypes` is empty
- `symbolDepth` is empty
- All counters reset to 0

Each re-entrant call starts from scratch, so the cycle is never detected.

### Key Files and Functions Involved

| File | Function | Role |
|------|----------|------|
| `internal/checker/nodebuilderimpl.go:2601` | `conditionalTypeToTypeNode` | Serializes conditional type, calls checker to resolve true/false types |
| `internal/checker/checker.go:23864` | `getTrueTypeFromConditionalType` | Resolves the true type of a conditional (triggers type resolution) |
| `internal/checker/checker.go:22348` | `getTypeFromIndexedAccessTypeNode` | Resolves indexed access type `ReturnType<typeof clone>[0]` |
| `internal/checker/checker.go:26434` | `getPropertyTypeForIndexType` | Finds property on type, generates error if not found |
| `internal/checker/printer.go:192` | `typeToStringEx` | Creates new NodeBuilder for error messages |
| `internal/checker/nodebuilderimpl.go:2524` | `createAnonymousTypeNode` | Serializes anonymous type (the function type) |
| `internal/checker/nodebuilderimpl.go:2404` | `createTypeNodeFromObjectType` | Serializes object type, including signatures |
| `internal/checker/nodebuilderimpl.go:1740` | `signatureToSignatureDeclarationHelper` | Serializes return type → re-enters conditional serialization |

## Fix Applied

### Approach: Depth Guard on the Checker

Added a `nodeBuilderDepth` counter on the `Checker` struct that tracks re-entrant calls to `typeToStringEx`. This counter persists across NodeBuilder instances because it lives on the Checker:

**`internal/checker/checker.go`**: Added `nodeBuilderDepth int` field to `Checker` struct.

**`internal/checker/printer.go`**: Added depth guard at the top of `typeToStringEx`:
```go
const maxNodeBuilderDepth = 10
if c.nodeBuilderDepth >= maxNodeBuilderDepth {
    return "..."
}
c.nodeBuilderDepth++
defer func() { c.nodeBuilderDepth-- }()
```

When the depth limit is hit, `typeToStringEx` returns `"..."` instead of recursing, which:
1. Prevents the stack overflow
2. Still produces meaningful error messages (e.g., `Property '0' does not exist on type '...'`)
3. Has no impact on normal (non-recursive) type serialization

### Why This Fix Is Correct

- The depth counter is on the `Checker` struct, not the `NodeBuilder`, so it survives across NodeBuilder context resets
- Depth limit of 10 is generous enough to handle any legitimate nesting
- Returning `"..."` is consistent with the existing truncation behavior used elsewhere in the NodeBuilder
- The fix is minimal and doesn't affect the hot path (non-recursive cases just increment/decrement a counter)

### Test Results

- New test case `stackOverflowSelfReferentialReturnType.ts` passes with correct error output
- All existing compiler tests (local and submodule) pass
- All fourslash tests pass
