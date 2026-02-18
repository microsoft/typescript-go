# Investigation: Panic with autoImportFileExcludePatterns

## Problem Summary

The LSP server panics with the message:
```
panic handling request textDocument/codeAction: textDocument/codeAction returned ErrNeedsAutoImports even after enabling auto imports
```

## Root Cause Analysis

### Key Files and Functions

1. **`internal/lsp/server.go:728-762`** - `registerLanguageServiceWithAutoImportsRequestHandler`
   - This function handles LSP requests that need auto-imports
   - When it gets `ErrNeedsAutoImports`, it calls `GetLanguageServiceWithAutoImports` to rebuild
   - If it gets `ErrNeedsAutoImports` again after rebuild, it panics (line 752)

2. **`internal/ls/languageservice.go:96-100`** - `getPreparedAutoImportView`
   - Returns `ErrNeedsAutoImports` if `registry.IsPreparedForImportingFile` returns false

3. **`internal/ls/autoimport/registry.go:221-248`** - `IsPreparedForImportingFile`
   - Calls `possiblyNeedsRebuildForFile` on both project and node_modules buckets (lines 230, 237)
   - Returns false if any bucket needs rebuild

4. **`internal/ls/autoimport/registry.go:86-91`** - `possiblyNeedsRebuildForFile`
   ```go
   func (b BucketState) possiblyNeedsRebuildForFile(file tspath.Path, preferences *lsutil.UserPreferences) bool {
       return b.newProgramStructure > 0 ||
           b.hasDirtyFileBesides(file) ||
           !core.UnorderedEqual(b.fileExcludePatterns, preferences.AutoImportFileExcludePatterns) || // Line 89 - CHECKS EXCLUDE PATTERNS
           b.dirtyPackages.Len() > 0
   }
   ```
   - **Line 89 checks if `fileExcludePatterns` changed**

5. **`internal/ls/autoimport/registry.go:637-728`** - `updateIndexes`
   - For node_modules buckets (lines 670-699):
     - Line 676 checks: `bucketState.multipleFilesDirty || !nodeModulesBucket.Value().DependencyNames.Equals(dependencies)`
     - **DOES NOT CHECK fileExcludePatterns**
   - For project buckets (lines 701-728):
     - Line 704: `shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile)`
     - Lines 705-713: Only checks `newProgramStructure`
     - **DOES NOT CHECK fileExcludePatterns**

### The Bug

There is an inconsistency between what `possiblyNeedsRebuildForFile` considers as "needs rebuild" vs what `updateIndexes` actually rebuilds:

1. `possiblyNeedsRebuildForFile` includes fileExcludePatterns check (line 89)
2. `updateIndexes` does NOT check fileExcludePatterns when deciding to rebuild

**Flow of the bug:**
1. User has `autoImportFileExcludePatterns` set in `.vscode/settings.json`
2. First request: `IsPreparedForImportingFile` returns `false` because `possiblyNeedsRebuildForFile` detects changed patterns
3. This triggers `ErrNeedsAutoImports`
4. LSP calls `GetLanguageServiceWithAutoImports` to rebuild
5. In `updateIndexes`, buckets are NOT rebuilt because only dirty files and program structure are checked
6. Second request: `IsPreparedForImportingFile` STILL returns `false` (patterns still differ)
7. `ErrNeedsAutoImports` is returned again → **PANIC**

### Code Locations

- **Panic location:** `internal/lsp/server.go:752`
- **ErrNeedsAutoImports definition:** `internal/ls/completions.go:32`
- **Check that causes false negative:** `internal/ls/autoimport/registry.go:89`
- **Missing check in rebuild logic:**
  - Project buckets: `internal/ls/autoimport/registry.go:704-713`
  - Node modules buckets: `internal/ls/autoimport/registry.go:676`

### User Preferences Structure

From `internal/ls/lsutil/userpreferences.go:73`:
```go
AutoImportFileExcludePatterns []string
```

This is populated from VSCode's `typescript.preferences.autoImportFileExcludePatterns` setting.

## Solution

The fix needs to ensure that `updateIndexes` checks for changed `fileExcludePatterns` and rebuilds buckets when they differ.

### For Project Buckets (line 704):

Currently:
```go
shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile)
```

Should check exclude patterns:
```go
shouldRebuild := project.Value().state.hasDirtyFileBesides(change.RequestedFile) ||
    !core.UnorderedEqual(project.Value().state.fileExcludePatterns, b.userPreferences.AutoImportFileExcludePatterns)
```

### For Node Modules Buckets (line 676):

Currently:
```go
needsFullRebuild := bucketState.multipleFilesDirty || !nodeModulesBucket.Value().DependencyNames.Equals(dependencies)
```

Should check exclude patterns:
```go
needsFullRebuild := bucketState.multipleFilesDirty || 
    !nodeModulesBucket.Value().DependencyNames.Equals(dependencies) ||
    !core.UnorderedEqual(bucketState.fileExcludePatterns, b.userPreferences.AutoImportFileExcludePatterns)
```

## Reproduction Steps

1. Create a pnpm monorepo with @typescript/native-preview
2. Add `.vscode/settings.json` with:
   ```json
   {
     "typescript.preferences.autoImportFileExcludePatterns": ["**/node_modules"]
   }
   ```
3. Open a TypeScript file
4. Write code like `fs.readFileSync()` without importing `fs`
5. Try to open Quick Fix menu
6. The LSP server should panic with the error message above
