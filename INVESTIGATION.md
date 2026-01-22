# Investigation: TestCodeFixPromoteTypeOnlyOrderingCrash

## Summary

Fixed a crash in `TestCodeFixPromoteTypeOnlyOrderingCrash` which occurred when promoting a type-only import to a value import with `verbatimModuleSyntax: true` and import reordering.

## Root Causes

The crash was caused by TWO separate bugs:

### Bug 1: `getVisualListRange` finding wrong token (Critical)

**Location**: `internal/format/indent.go`, function `getVisualListRange`

**Problem**: The function was using `astnav.FindNextToken(prior, node, sourceFile)` to find the token after a list. However, this was finding the first child token WITHIN the list instead of the token AFTER the list (e.g., finding "AAA" instead of the closing brace `}`).

**Impact**: When `GetContainingList` checked if a node was within the visual list range, it returned `nil` because the visual range was too small (e.g., 13-18 instead of 13-32), causing the crash.

**Fix**: Changed the logic to use `astnav.FindPrecedingToken(sourceFile, searchPos)` with incremental searching to find the actual closing token after the list:

```go
// Find the token that comes after the list ends
searchPos := list.End() + 1
next := astnav.FindPrecedingToken(sourceFile, searchPos)
for next != nil && next.End() <= list.End() && searchPos < sourceFile.End() {
    searchPos++
    next = astnav.FindPrecedingToken(sourceFile, searchPos)
}
```

### Bug 2: Incorrect sorting assumption during import promotion

**Location**: `internal/ls/autoimport/fix.go`, function `promoteImportClause`

**Problem**: When promoting an import specifier from type-only to value (e.g., promoting BBB in `import type { AAA, BBB }`), the code assumed value imports should always go first (index 0). This was incorrect because:
1. It didn't respect the actual sort order preference
2. When using inline type modifiers, alphabetical sorting is typically expected

**Impact**: Specifiers were being moved to the wrong position, and the comparison logic was using incompatible sorting preferences.

**Fix**: 
1. Changed to use `OrganizeImportsTypeOrderInline` when converting to inline type modifiers
2. Created synthetic specifiers with type modifiers to properly compute the correct insertion index
3. Only reorder if the insertion index actually changes

```go
// Use inline type ordering when converting to inline type modifiers
prefsClone := *prefsForInlineType
prefsClone.OrganizeImportsTypeOrder = lsutil.OrganizeImportsTypeOrderInline

// Create synthetic specifiers to find correct position
specsWithTypeModifiers := core.Map(namedImportsData.Elements.Nodes, func(e *ast.Node) *ast.Node {
    // ... add type modifiers to non-promoted elements ...
})

insertionIndex := organizeimports.GetImportSpecifierInsertionIndex(
    specsWithTypeModifiers, newSpecifier, specifierComparer)
```

## Test Case

The test promotes BBB from `import type { AAA, BBB }` to:
```typescript
import {
    type AAA,
    BBB,
} from "./bar";
```

With `verbatimModuleSyntax: true`, this requires:
1. Removing the top-level `type` keyword
2. Adding inline `type` modifier to AAA
3. Keeping BBB as a value import
4. Maintaining alphabetical order (AAA before BBB)

## Files Changed

1. `internal/format/indent.go` - Fixed `getVisualListRange` to find correct closing token
2. `internal/ls/autoimport/fix.go` - Fixed import reordering logic to use correct sorting preference
3. (No changes to delete.go needed - the assertion was correct, it was just catching the bug)

## Verification

- ✅ `TestCodeFixPromoteTypeOnlyOrderingCrash` now passes
- ✅ All autoimport tests pass
- ✅ All format tests pass
