# Investigation: Incorrect `require` Emission for Type-Only Imports

## Problem Summary
`tsgo` incorrectly emits a `require` statement when an ES6 import contains:
1. A `type` import specifier (e.g., `type ValueData`)
2. A regular import specifier (e.g., `Value`) that is ONLY used in type positions

## Reproduction

### Minimal Test Case
```typescript
// provider.ts
export class Value {
    data: string = "";
}
export type ValueData = { data: string };

// consumer.ts - BROKEN
import { Value, type ValueData } from "./provider";
export function test(value: Value): Value {
    return null as any;
}

// consumer-alt.ts - WORKS
import { Value } from "./provider";
export function test(value: Value): Value {
    return null as any;
}
```

### Actual Output (tsgo)
```javascript
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.test = test;
require("./provider");  // ❌ THIS SHOULD NOT BE HERE
function test(value) {
    return value;
}
```

### Expected Output (TypeScript)
```javascript
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.test = test;
// No require() statement
function test(value) {
    return value;
}
```

## Key Finding

The bug ONLY occurs when ALL of the following conditions are met:
1. The import statement contains both a `type` import specifier and a regular import specifier
2. The regular import specifier is ONLY used in type positions (not runtime positions)

### Working Cases (No Bug)
```typescript
// Case 1: Only regular import, type-only usage - WORKS!
import { Value } from "./provider";
export function test(value: Value): Value { return value; }
// Output: No require() ✓

// Case 2: Only type import - WORKS!
import { type Value } from "./provider";
export function test(value: Value): Value { return value; }
// Output: No require() ✓
```

### Broken Case (Bug)
```typescript
// Both type import and regular import (type-only usage) - BROKEN!
import { Value, type ValueData } from "./provider";
export function test(value: Value): Value { return value; }
// Output: require("./provider"); ❌
```

## Root Cause Analysis

## Root Cause Analysis

### Import Elision Logic
The import elision logic is in `/internal/transformers/tstransforms/importelision.go`:

1. **KindImportDeclaration** (line 42-55): Visits the import declaration
   - If `ImportClause` exists, it visits the import clause
   - If the visited import clause is `nil`, the entire import is elided

2. **KindImportClause** (line 56-64): Processes the import clause
   - Checks if default import (`name`) should be emitted via `shouldEmitAliasDeclaration`
   - Visits `namedBindings` (which contains the import specifiers)
   - If both are `nil`, returns `nil` to elide the entire import

3. **KindNamedImports** (line 71-78): Processes named imports (`{ Value, type ValueData }`)
   - Visits each import specifier
   - If all specifiers are elided (empty `elements`), returns `nil`

4. **KindImportSpecifier** (line 79-84): Processes individual import specifiers
   - Calls `shouldEmitAliasDeclaration` to determine if specifier should be kept
   - If false, returns `nil` to elide this specifier

### The `shouldEmitAliasDeclaration` Check
This function (line 130-132) checks:
```go
func (tx *ImportElisionTransformer) shouldEmitAliasDeclaration(node *ast.Node) bool {
    return ast.IsInJSFile(node) || tx.isReferencedAliasDeclaration(node)
}
```

It calls `isReferencedAliasDeclaration` which checks `aliasLinks.referenced` in the checker.

### The Bug

The issue is that when processing an import clause with both type and non-type imports:
1. The `type ValueData` specifier is correctly elided (marked as type-only)
2. The `Value` specifier should also be elided (only used in types)
3. BUT: `Value`'s `aliasLinks.referenced` is incorrectly being set to `true`

The presence of the `type` import specifier in the same import statement seems to be triggering some code path that marks the non-type import as referenced, even though it's only used in type positions.

### Comparison with TypeScript

TypeScript's `markSymbolOfAliasDeclarationIfTypeOnly` function (in `checker.ts` around line 4385) has additional parameters `immediateTarget` and `finalTarget` and includes logic to check if the target resolves to a type-only declaration:

```typescript
function markSymbolOfAliasDeclarationIfTypeOnlyWorker(
    aliasDeclarationLinks: SymbolLinks, 
    target: Symbol | undefined, 
    overwriteEmpty: boolean
): boolean {
    if (target && (aliasDeclarationLinks.typeOnlyDeclaration === undefined || overwriteEmpty && aliasDeclarationLinks.typeOnlyDeclaration === false)) {
        const exportSymbol = target.exports?.get(InternalSymbolName.ExportEquals) ?? target;
        const typeOnly = exportSymbol.declarations && find(exportSymbol.declarations, isTypeOnlyImportOrExportDeclaration);
        aliasDeclarationLinks.typeOnlyDeclaration = typeOnly ?? getSymbolLinks(exportSymbol).typeOnlyDeclaration ?? false;
    }
    return !!aliasDeclarationLinks.typeOnlyDeclaration;
}
```

The tsgo version is missing this additional logic to check if the alias target is type-only.

## Hypothesis

There are two possible root causes:

### Hypothesis 1: Missing Type-Only Target Check
The `markSymbolOfAliasDeclarationIfTypeOnly` function in tsgo is missing logic to check if the alias target resolves to a type-only declaration. When `Value` is imported alongside `type ValueData`, the function should check if all uses of `Value` are in type-only positions, but this check is not implemented.

**Evidence**: TypeScript's version has `markSymbolOfAliasDeclarationIfTypeOnlyWorker` that checks the target symbol's declarations for type-only status.

### Hypothesis 2: Incorrect Reference Marking
When type references are resolved (via `resolveEntityName` called from `getTypeFromTypeReference`), something is incorrectly marking the import alias as referenced. The code should distinguish between:
- Type-level reference: Should NOT mark `aliasLinks.referenced = true`
- Value-level reference: Should mark `aliasLinks.referenced = true`

**Evidence**: The import works correctly when there's NO `type` import in the same statement, suggesting that the presence of the type-only import is affecting how the non-type import is processed.

## Most Likely Cause

Based on the investigation, **Hypothesis 2** appears most likely. The bug manifests ONLY when there's a mix of type-only and non-type-only imports in the same import statement, which suggests that something in the import processing logic is incorrectly treating all non-type-only imports as referenced when any import from that module is used (even if only in type positions).

The reference marking likely happens during the initial checking of the import declaration or during binding, where the code doesn't properly distinguish between imports that will be used at runtime vs imports that will only be used in type positions.

## Files to Investigate

1. `/internal/transformers/tstransforms/importelision.go` - Import elision transformer
2. `/internal/checker/emitresolver.go` - `IsReferencedAliasDeclaration` implementation
3. `/internal/checker/checker.go` - `markAliasSymbolAsReferenced`, `markAliasReferenced`, etc.
4. `/internal/ast/utilities.go` - `IsPartOfTypeNode` and related functions

## Test Case for Regression

Add this test to ensure the fix works:

```typescript
// @module: commonjs
// @target: es2020

// @filename: provider.ts
export class Value {
    data: string = "";
}
export type ValueData = { data: string };

// @filename: consumer.ts
import { Value, type ValueData } from "./provider";

export function test(value: Value): Value {
    return value;
}
```

Expected output for `consumer.js` should NOT contain `require("./provider")`.
