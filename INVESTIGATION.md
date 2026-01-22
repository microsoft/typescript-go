# Investigation: TestCodeFixPromoteTypeOnlyOrderingCrash

## Summary

Fixed a crash in `TestCodeFixPromoteTypeOnlyOrderingCrash` which occurred when promoting a type-only import to a value import with `verbatimModuleSyntax: true` and import reordering.

## Root Causes

The crash was caused by TWO separate bugs:

### Bug 1: `getVisualListRange` finding wrong token (Critical)

**Location**: `internal/format/indent.go`, function `getVisualListRange`

**Problem**: The function was using `astnav.FindNextToken(prior, node, sourceFile)` to find the token after a list. However, this was finding the first child token WITHIN the list instead of the token AFTER the list (e.g., finding "AAA" instead of the closing brace `}`).

**Impact**: When `GetContainingList` checked if a node was within the visual list range, it returned `nil` because the visual range was too small (e.g., 13-18 instead of 13-32), causing the crash with the assertion "containingList should not be nil".

**Fix**: Changed the logic to use `astnav.FindPrecedingToken(sourceFile, searchPos)` with incremental searching to find the actual closing token after the list:

```go
// Find the token that comes after the list ends
searchPos := list.End() + 1
next := astnav.FindPrecedingToken(sourceFile, searchPos)
// Advance until we find a token that starts after the list end
maxSearchDistance := 10 // Limit for performance
for next != nil && next.End() <= list.End() && searchDistance < maxSearchDistance {
    searchPos++
    next = astnav.FindPrecedingToken(sourceFile, searchPos)
}
```

### Bug 2: Incorrect sorting assumption during import promotion

**Location**: `internal/ls/autoimport/fix.go`, function `promoteImportClause`

**Problem**: When promoting an import specifier from type-only to value (e.g., promoting BBB in `import type { AAA, BBB }`), the code assumed value imports should always go first (index 0). This was incorrect because:
1. It didn't respect the actual sort order preference
2. When using inline type modifiers (required by `verbatimModuleSyntax`), alphabetical sorting is typically expected
3. The comparison was using incompatible sorting preferences (default "Last" vs needed "Inline")

**Impact**: Specifiers were being moved to the wrong position, causing incorrect output with mangled formatting.

**Fix**: 
1. Changed to use `OrganizeImportsTypeOrderInline` when converting to inline type modifiers
2. Created synthetic specifiers with type modifiers to properly compute the correct insertion index
3. Only reorder if the insertion index actually changes
4. Extracted `createSyntheticImportSpecifier` helper to reduce code duplication

```go
// Use inline type ordering when converting to inline type modifiers
prefsClone := *prefsForInlineType
prefsClone.OrganizeImportsTypeOrder = lsutil.OrganizeImportsTypeOrderInline

// Create synthetic specifiers to find correct position
specsWithTypeModifiers := core.Map(namedImportsData.Elements.Nodes, func(e *ast.Node) *ast.Node {
    // ... create specifiers with type modifiers ...
    return createSyntheticImportSpecifier(changes.NodeFactory, s, true)
})

insertionIndex := organizeimports.GetImportSpecifierInsertionIndex(
    specsWithTypeModifiers, newSpecifier, specifierComparer)
    
// Only reorder if position changes
if insertionIndex != aliasIndex {
    changes.Delete(sourceFile, aliasDeclaration)
    changes.InsertImportSpecifierAtIndex(sourceFile, newSpecifier, namedImports, insertionIndex)
}
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
2. `internal/ls/autoimport/fix.go` - Fixed import reordering logic to use correct sorting preference and added helper function

## Code Review Feedback Addressed

1. ✅ Added search distance limit (10 chars) to prevent performance issues in large files
2. ✅ Extracted `createSyntheticImportSpecifier` helper to reduce code duplication
3. ✅ Added clarifying comments

## Verification

- ✅ `TestCodeFixPromoteTypeOnlyOrderingCrash` now passes
- ✅ All autoimport tests pass (0.449s)
- ✅ All format tests pass (0.457s)
- ✅ CodeQL security scan: No issues detected

## Security Summary

No security vulnerabilities were introduced or discovered during this investigation. The changes are purely correctness fixes for existing functionality.
