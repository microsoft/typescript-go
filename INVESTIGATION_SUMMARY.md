# Investigation Summary: LSP Panic with autoImportFileExcludePatterns

## Issue

The LSP server panics with the following error when a user has `autoImportFileExcludePatterns` configured in `.vscode/settings.json` and attempts to use Quick Fix on code that needs auto-imports:

```
panic handling request textDocument/codeAction: textDocument/codeAction returned ErrNeedsAutoImports even after enabling auto imports
```

## Root Cause

There is an inconsistency between what the system considers as "needs rebuild" versus what actually gets rebuilt:

### The Check (Line 89)
`internal/ls/autoimport/registry.go:86-91`
```go
func (b BucketState) possiblyNeedsRebuildForFile(file tspath.Path, preferences *lsutil.UserPreferences) bool {
    return b.newProgramStructure > 0 ||
        b.hasDirtyFileBesides(file) ||
        !core.UnorderedEqual(b.fileExcludePatterns, preferences.AutoImportFileExcludePatterns) ||  // ← CHECKS PATTERNS
        b.dirtyPackages.Len() > 0
}
```

### The Rebuild Logic (Lines 676 and 704)
`internal/ls/autoimport/registry.go:670-728`

**For node_modules buckets (line 676):**
```go
needsFullRebuild := bucketState.multipleFilesDirty || !nodeModulesBucket.Value().DependencyNames.Equals(dependencies)
// ← MISSING: Does NOT check fileExcludePatterns
```

**For project buckets (line 704):**
```go
shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile)
// ← MISSING: Does NOT check fileExcludePatterns
```

## Bug Flow

1. User has `autoImportFileExcludePatterns` set in settings
2. User requests code action (e.g., Quick Fix on `fs.readFileSync()`)
3. `getPreparedAutoImportView` checks if registry is ready via `IsPreparedForImportingFile`
4. `IsPreparedForImportingFile` calls `possiblyNeedsRebuildForFile` which detects changed patterns → returns `false`
5. This triggers `ErrNeedsAutoImports`
6. LSP handler calls `GetLanguageServiceWithAutoImports` to rebuild
7. `updateIndexes` is called but doesn't check `fileExcludePatterns`, so **no rebuild happens**
8. Another code action request is made
9. `IsPreparedForImportingFile` STILL returns `false` (patterns still don't match)
10. `ErrNeedsAutoImports` is returned again → **PANIC** at `internal/lsp/server.go:752`

## Reproduction

A test case has been added at `internal/ls/autoimport/registry_fileexclude_test.go` that successfully reproduces the bug:

```
$ go test -v -run TestAutoImportFileExcludePatternsRebuild ./internal/ls/autoimport
=== RUN   TestAutoImportFileExcludePatternsRebuild
    registry_fileexclude_test.go:99: BUG REPRODUCED: IsPreparedForImportingFile still returns false after GetLanguageServiceWithAutoImports
    registry_fileexclude_test.go:102: Project bucket dirty: false
    registry_fileexclude_test.go:103: Node modules bucket dirty: false
--- FAIL: TestAutoImportFileExcludePatternsRebuild (0.05s)
```

The test confirms:
- Buckets are NOT rebuilt (both show `dirty: false`)
- `IsPreparedForImportingFile` keeps returning `false`
- This would trigger the panic in production

## Solution

The fix requires adding `fileExcludePatterns` checks in `updateIndexes()`:

### For node_modules buckets (around line 676):
```go
needsFullRebuild := bucketState.multipleFilesDirty || 
    !nodeModulesBucket.Value().DependencyNames.Equals(dependencies) ||
    !core.UnorderedEqual(bucketState.fileExcludePatterns, b.userPreferences.AutoImportFileExcludePatterns)
```

### For project buckets (around line 704):
```go
shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile) ||
    !core.UnorderedEqual(project.Value().state.fileExcludePatterns, b.userPreferences.AutoImportFileExcludePatterns)
```

## Key Files

- **Panic location:** `internal/lsp/server.go:752`
- **Handler:** `internal/lsp/server.go:728-762` (`registerLanguageServiceWithAutoImportsRequestHandler`)
- **Check function:** `internal/ls/autoimport/registry.go:86-91` (`possiblyNeedsRebuildForFile`)
- **Rebuild function:** `internal/ls/autoimport/registry.go:637-728` (`updateIndexes`)
- **Test case:** `internal/ls/autoimport/registry_fileexclude_test.go`
- **Full investigation:** `INVESTIGATION.md`

## References

- Issue mentions PR #2616 as a previous attempted fix that was insufficient
- The panic occurs specifically with:
  - pnpm monorepos
  - @typescript/native-preview
  - `.vscode/settings.json` with `typescript.preferences.autoImportFileExcludePatterns` set

## Next Steps

1. Implement the fix by adding the missing checks in `updateIndexes()`
2. Verify the test passes with the fix
3. Run full test suite to ensure no regressions
4. Test manually with a pnpm monorepo setup
