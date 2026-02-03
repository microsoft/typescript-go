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
3. BUT: Something is incorrectly marking `Value` as referenced

The presence of the `type` import specifier seems to be triggering some code path that marks the non-type import as referenced, even though it's only used in type positions.

## Investigation Needed

### Hypothesis
When the import clause contains a mix of `type` and non-`type` imports, there may be:

1. **Shared alias resolution**: The import clause might be sharing some state or resolution logic that causes the regular import to be marked as referenced when the type import is processed.

2. **Import clause level check**: There might be a check at the import clause level that marks all non-type imports as referenced if any import is used.

3. **Symbol merging issue**: The symbols for `Value` and `ValueData` might be getting merged or confused, causing the reference tracking to conflate them.

### Next Steps

1. **Add debugging**: Insert logging in `markAliasSymbolAsReferenced` to trace when and why `Value` is being marked as referenced.

2. **Check import clause processing**: Investigate how `KindImportClause` handles mixed type/non-type imports.

3. **Symbol resolution**: Check if there's symbol confusion between `Value` (the class) and its usage in type positions.

4. **Compare with TypeScript**: Look at TypeScript's `checker.ts` to see how it handles this case differently.

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
