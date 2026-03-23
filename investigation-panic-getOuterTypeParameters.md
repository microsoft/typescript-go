# Investigation: Panic in `getOuterTypeParameters` (20260321.1 release)

## Issue Summary
A nil pointer dereference panic occurs in `getOuterTypeParameters` when a symbol with `SymbolFlagsClass` has a nil `ValueDeclaration`. The panic was introduced in release `20260321.1` by PR #3157 ("Fix named export aliases merging with `export =`").

## Stack Trace Analysis
The crash occurs at `checker.go:23121` in `getOuterTypeParameters` when accessing `node.Parent` where `node` is nil. The call chain is:

```
checkClassDeclaration → getTypeOfSymbol → getTypeOfFuncClassEnumModule →
getBaseTypeVariableOfClass → getBaseConstructorTypeOfClass → checkExpression →
checkIdentifier → getTypeOfAlias → getTypeOfSymbol → getTypeOfFuncClassEnumModule →
getBaseTypeVariableOfClass → getDeclaredTypeOfClassOrInterface →
getOuterTypeParametersOfClassOrInterface → getOuterTypeParameters → CRASH
```

Key observations:
1. There are TWO calls to `getBaseTypeVariableOfClass` - the outer one for the class being checked, and the inner one for the base class resolved through an alias
2. The inner class symbol has `SymbolFlagsClass` but nil `ValueDeclaration`
3. The `debug.AssertIsDefined` at line 23113 is a no-op in production builds (built with `noassert` tag)

## Root Cause Analysis

### PR #3157 Changes
The PR made three changes to `internal/checker/checker.go`:

1. **`hasExportedMembersOfKind`**: Changed from `symbol.Flags` to `c.getSymbolFlags(symbol)` to resolve aliases when checking flags
2. **`getExternalModuleMember`**: Changed `getExportOfModule` to use `moduleSymbol` instead of `targetSymbol` for `export =` modules
3. **`getExportsOfModuleWorker`**: Expanded filter to include namespace exports and resolve aliases via `getSymbolFlags`

### Identified Issue: Missing Equality Check
The Go code in `getExternalModuleMember` is missing a check that exists in the TypeScript source. 

**TypeScript** (checker.ts):
```ts
const symbol = symbolFromModule && symbolFromVariable && symbolFromModule !== symbolFromVariable ?
    combineValueAndTypeSymbols(symbolFromVariable, symbolFromModule) :
    symbolFromModule || symbolFromVariable;
```

**Go** (before fix):
```go
symbol := symbolFromVariable
if symbolFromModule != nil {
    symbol = symbolFromModule
    if symbolFromVariable != nil {
        symbol = c.combineValueAndTypeSymbols(symbolFromVariable, symbolFromModule)
    }
}
```

The Go version was missing `symbolFromModule !== symbolFromVariable` - it always called `combineValueAndTypeSymbols` when both were non-nil, even when they were the same symbol. Before PR #3157, this didn't matter because `symbolFromModule` was always nil for `export =` modules (since the resolved target class doesn't have `SymbolFlagsModule`). After the PR, `symbolFromModule` can be non-nil, making this missing check potentially impactful.

`combineValueAndTypeSymbols` creates a new transient symbol with `valueSymbol.Flags | typeSymbol.Flags` and `ValueDeclaration = valueSymbol.ValueDeclaration`. If the value symbol's `ValueDeclaration` is nil (possible for synthesized/transient properties), and the type symbol has `SymbolFlagsClass`, the merged result would have `SymbolFlagsClass` but nil `ValueDeclaration`.

### Defensive Fix
Added a nil guard in `getOuterTypeParametersOfClassOrInterface` to prevent the crash even if a class symbol somehow ends up without a `ValueDeclaration`. This matches the existing `debug.AssertIsDefined` pattern but provides runtime safety in production builds where assertions are disabled.

## Changes Made
1. **Fixed `getExternalModuleMember`** to match TypeScript's logic: only call `combineValueAndTypeSymbols` when both symbols are non-nil AND different (pointer inequality check)
2. **Added defensive nil check** in `getOuterTypeParametersOfClassOrInterface` to return nil instead of crashing when `declaration` is nil

## Reproduction Status
Could not reproduce with small test cases. The issue likely requires a specific combination of:
- A module with `export =` pointing to a class
- Named type/namespace exports from the same module  
- A class in another module extending something imported through this module
- The export = target having static members or inherited properties with the same name as the named exports

The reporter mentions a "large private monorepo" which likely has complex module resolution chains that create the specific conditions for the crash.
