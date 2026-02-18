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

## Root Cause: Symlink Path Mismatch in `autoImportBuilderFS` Cache

**Confirmed via test**: `TestAutoImportBuilderFS/symlink_cache_mismatch` in
`internal/project/snapshotfs_test.go`.

### Mechanism

The root cause is a cache key mismatch between the path used during module
resolution's `FileExists` check and the path used in the subsequent
`GetSourceFile` call.

`autoImportBuilderFS.GetFileByPath` stores file handles in its `untrackedFiles`
cache keyed by `path` (derived from `toPath(fileName)`). When module resolution
and `GetSourceFile` use **different paths** for the same physical file (symlink
vs realpath), the cache is bypassed.

### Step-by-Step Scenario

1. **Module resolution calls `FileExists(symlinkPath)`**:
   - `sourceFS.FileExists` → `sourceFS.GetFile(symlinkPath)` → 
     `autoImportBuilderFS.GetFile(symlinkPath)` → 
     `GetFileByPath(symlinkPath, toPath(symlinkPath))`
   - File is read from disk and stored in `untrackedFiles[toPath(symlinkPath)]`
   - `FileExists` returns `true` → module resolution succeeds

2. **Module resolution resolves the symlink**:
   - `createResolvedModuleHandlingSymlink` calls `getOriginalAndResolvedFileName`
   - `realPath(symlinkPath)` → returns the **realpath** (different from symlink path)
   - `ResolvedModule.ResolvedFileName = realpathPath`
   - `ResolvedModule.OriginalPath = symlinkPath`

3. **File gets deleted from disk** (e.g., concurrent `npm install`):
   - The real file and the symlink both become unavailable

4. **Checker calls `GetSourceFileForResolvedModule(realpathPath)`**:
   - `GetSourceFile(realpathPath)` → `host.GetSourceFile(realpathPath, toPath(realpathPath))`
   - `autoImportBuilderFS.GetFile(realpathPath)` → 
     `GetFileByPath(realpathPath, toPath(realpathPath))`
   - `toPath(realpathPath) ≠ toPath(symlinkPath)` → **cache miss** in `untrackedFiles`
   - Falls through to `snapshotFSBuilder.fs.ReadFile(realpathPath)` → **file deleted** → nil
   - `GetSourceFile` returns nil
   - `binder.BindSourceFile(nil)` → **CRASH**

### Why This Doesn't Happen with `compiler.Program`

`compiler.Program.GetSourceFile` looks up files from a pre-built `filesByPath` map.
All files are loaded during program construction, so symlink resolution happens once
and the map contains entries for the correct paths. There is no on-demand loading
that can miss a cache entry due to path differences.

### VFS Architecture

```
autoImportRegistryCloneHost
  └── sourceFS (implements vfs.FS)
        ├── GetFile(path) → autoImportBuilderFS.GetFile(path)
        ├── FileExists(path) → GetFile(path) || source.FS().FileExists(path)  
        └── autoImportBuilderFS (implements FileSource)
              ├── overlays (open files, keyed by path)
              ├── diskFiles (dirty.SyncMap, from previous snapshots, keyed by path)
              ├── untrackedFiles (SyncMap, for on-demand loads, keyed by path)
              └── snapshotFSBuilder.fs (cachedvfs.FS → OS filesystem)
```

The critical detail: `untrackedFiles` is keyed by `toPath(fileName)`. When
`fileName` is a symlink path, the key is different from the key for the realpath.
The file is only stored once, at whichever path was used first.

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

The nil filter at lines 1080-1082 handles nil entries but only *after* the crash.
This is fixed in this PR.

## Fix

1. Add nil guard in `aliasResolver.GetSourceFile` before `binder.BindSourceFile`
2. Add nil guard in `registry.go:updateIndexes` before `extractFromFile`
3. Add nil guard in `registry.go:createAliasResolver` before `binder.BindSourceFile`
