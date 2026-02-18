# Investigation: nil source file in aliasResolve.GetSourceFile

## Issue Summary

A nil pointer dereference (SIGSEGV) occurs in `aliasResolver.GetSourceFile` when
`binder.BindSourceFile` is called on a nil `*ast.SourceFile`. This crashes the
language server entirely, requiring a manual restart.

## Immediate Cause

In `internal/ls/autoimport/aliasresolver.go:80-84`:

```go
func (r *aliasResolver) GetSourceFile(fileName string) *ast.SourceFile {
    file := r.host.GetSourceFile(fileName, r.toPath(fileName))
    binder.BindSourceFile(file)  // panics if file == nil
    return file
}
```

`r.host.GetSourceFile()` can return `nil`, but `GetSourceFile` does not check for
this before passing it to `binder.BindSourceFile()`, which immediately calls
`file.IsBound()` — a nil pointer dereference.

The checker at `checker.go:14809` already handles nil returns from
`GetSourceFileForResolvedModule`, so returning nil is correct behavior. The bug
is that `aliasResolver.GetSourceFile` crashes *internally* before it gets a
chance to return nil.

## Call Chain

```
checker.resolveAlias
  → checker.getTargetOfAliasDeclaration
    → checker.getTargetOfImportEqualsDeclaration
      → checker.resolveExternalModuleName
        → checker.resolveExternalModule (checker.go:14806)
          → program.GetSourceFileForResolvedModule
            → aliasResolver.GetSourceFileForResolvedModule (aliasresolver.go:127)
              → aliasResolver.GetSourceFile (aliasresolver.go:80)
                → binder.BindSourceFile(nil) → CRASH
```

## Root Cause Analysis: Why Module Resolution Succeeds but GetSourceFile Returns nil

### Background

The `aliasResolver` acts as a lightweight `checker.Program` for auto-import. Unlike
`compiler.Program` (which pre-loads all files), the alias resolver loads files
on-demand. The key question is: how can module resolution succeed (finding a file
via `FileExists`) but `GetSourceFile` subsequently fail for the same file?

### VFS Architecture

The auto-import host wraps `autoImportBuilderFS` with `sourceFS`:

```
autoImportRegistryCloneHost
  └── sourceFS (implements vfs.FS)
        └── autoImportBuilderFS (implements FileSource)
              └── snapshotFSBuilder
                    ├── overlays (open files)
                    ├── diskFiles (dirty.SyncMap, from previous snapshots)
                    └── fs (cachedvfs.FS → OS filesystem)
```

The module resolver uses `host.FS()` = `sourceFS` for file existence checks.
`host.GetSourceFile()` also goes through `sourceFS.GetFile()` → `autoImportBuilderFS`.

### The `sourceFS.FileExists` Two-Stage Check

`sourceFS.FileExists` has a critical fallback path:

```go
func (fs *sourceFS) FileExists(path string) bool {
    if fh := fs.GetFile(path); fh != nil {
        return true
    }
    return fs.source.FS().FileExists(path)  // fallback to raw VFS
}
```

The fallback `fs.source.FS().FileExists(path)` goes to `snapshotFSBuilder.fs`,
which is a `cachedvfs.FS`. Its `FileExists` is **cached** but its `ReadFile` is
**NOT cached** (always passes through to OS).

### Leading Theory: Stale `cachedvfs.FileExists` Cache

The most plausible scenario involves the `cachedvfs` having a stale `FileExists`
cache entry:

1. A file (e.g., `@types/node/fs.d.ts`) exists in `diskFiles` from a previous 
   snapshot and is also cached as `true` in `cachedvfs.fileExistsCache`.
2. The file gets deleted from disk (e.g., `npm install` running).
3. `invalidateNodeModulesCache()` marks the `diskFiles` entry as `needsReload=true`.
4. During auto-import registry building, the checker calls module resolution via
   `aliasResolver.GetResolvedModule`.
5. Module resolution calls `sourceFS.FileExists(path)`:
   - `GetFile` → `autoImportBuilderFS` → `diskFiles.Load` finds the dirty entry → 
     `reloadEntryIfNeeded` tries to re-read from disk → fails → deletes entry → 
     returns nil → `GetFile` returns nil.
   - **Fallback**: `snapshotFSBuilder.fs.FileExists(path)` → `cachedvfs` returns 
     stale `true` from its cache.
   - `FileExists` returns `true` → module resolution succeeds.
6. The checker calls `GetSourceFileForResolvedModule(resolvedFileName)` → 
   `GetSourceFile(resolvedFileName)`.
7. `sourceFS.GetFile` → `autoImportBuilderFS.GetFile`:
   - `diskFiles.Load` now finds the base map entry (dirty entry was deleted, base 
     entry from previous snapshot still exists with `needsReload=false`) → returns 
     stale file handle → `GetSourceFile` succeeds with stale data.

This scenario actually results in stale data being served rather than nil. But
there are edge cases within the `dirty.SyncMap` that might not fall back to the
base map cleanly, especially under concurrent access patterns.

### Alternative Theory: Path Mismatch from Symlink Resolution

Module resolution resolves symlinks via `createResolvedModuleHandlingSymlink`:

```go
func (r *resolutionState) getOriginalAndResolvedFileName(fileName string) (string, string) {
    resolvedFileName := r.realPath(fileName)
    ...
    return fileName, resolvedFileName  // originalPath=symlink, resolvedFileName=realpath
}
```

The `ResolvedFileName` stored in the resolved module is the **realpath**, but
`FileExists` during resolution was called on the **original (symlink) path**.

If the file handle was stored in `autoImportBuilderFS.untrackedFiles` under the
symlink path's key, then `GetSourceFile` called with the realpath would miss the
cache. It would then attempt to read from disk at the realpath. If the realpath
is valid, this should succeed. But if there's any path normalization difference
between what `ReadFile` receives and the actual filesystem path, this could fail.

### Incomplete Understanding

The exact triggering condition remains uncertain. The scenarios investigated 
don't perfectly explain a nil return from `GetSourceFile` when module resolution
succeeded via the same VFS. Possible remaining gaps:

1. **Concurrent access to `dirty.SyncMap`**: During the second pass in
   `updateIndexes`, `diskFiles` entries might be mutated by concurrent operations
   from other snapshot building phases.
2. **`wrapvfs` Realpath transformation**: When the module resolver uses a custom
   `Realpath` (via `getModuleResolver`), the resolved path might not correspond
   to a path that `autoImportBuilderFS.GetFile` can read.
3. **Edge case in `dirty.SyncMap` after deletion + base map interaction**: After
   an entry is deleted in the dirty map, the next `Load` might create a new
   entry from the base map. The exact behavior under rapid concurrent access
   needs further investigation.

## Additional Bug Found

In `registry.go:createAliasResolver` (line ~1072-1076), there's another nil-unsafe
call to `binder.BindSourceFile`:

```go
wg.Go(func() {
    file := b.host.GetSourceFile(realpathFileName, realpathPath)
    binder.BindSourceFile(file)  // <-- also crashes if file is nil
    rootFiles[i] = file
})
```

However, this is partly mitigated by the nil filter at lines 1080-1082:
```go
rootFiles = slices.DeleteFunc(rootFiles, func(f *ast.SourceFile) bool {
    return f == nil
})
```

This filter handles the case where `GetSourceFile` returns nil, but only
*after* the crash in `binder.BindSourceFile`. This is a separate bug that
should also be fixed.

## Fix

Add nil guards in `GetSourceFile` (the immediate crash site):

```go
func (r *aliasResolver) GetSourceFile(fileName string) *ast.SourceFile {
    file := r.host.GetSourceFile(fileName, r.toPath(fileName))
    if file == nil {
        return nil
    }
    binder.BindSourceFile(file)
    return file
}
```

And in `registry.go:updateIndexes` (line ~783) where `GetSourceFile` result 
is used without nil check:

```go
sourceFile := aliasResolver.GetSourceFile(source.fileName)
if sourceFile == nil {
    return true
}
```
