# Investigation: nil source file in aliasResolve.GetSourceFile

## Issue Summary

A nil pointer dereference (SIGSEGV) occurs in `aliasResolver.GetSourceFile` when
`binder.BindSourceFile` is called on a nil `*ast.SourceFile`. This crashes the
language server entirely, requiring a manual restart.

## Root Cause

### Immediate Cause

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

### Why `host.GetSourceFile()` Returns nil

The host implementation is `autoImportRegistryCloneHost.GetSourceFile` in
`internal/project/autoimport.go:148-167`:

```go
func (a *autoImportRegistryCloneHost) GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile {
    fh := a.fs.GetFile(fileName)
    if fh == nil {
        return nil  // <-- returns nil when file doesn't exist
    }
    // ... parse and return
}
```

The underlying `SnapshotFS.GetFileByPath` can return nil when:
1. The file doesn't exist on disk
2. The file was deleted between snapshots
3. A module resolves to a path that doesn't exist in the current snapshot

### The Trigger Scenario

The crash happens during auto-import alias resolution. The call chain is:

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

The `aliasResolver` is a lightweight implementation of `checker.Program` used
specifically for auto-import. Unlike the real `compiler.Program`:

1. **`compiler.Program.GetSourceFile`** looks up files from `filesByPath` map — 
   it only contains files that were successfully loaded during program construction,
   so it never returns a file that needs binding but doesn't exist.

2. **`aliasResolver.GetSourceFile`** resolves files on-demand via the host's 
   `SnapshotFS`. Module resolution can resolve to a file path, but the file may
   not actually exist in the current snapshot (e.g., the file was just deleted,
   or module resolution points to a file that hasn't been loaded yet).

### Why This Doesn't Happen with `compiler.Program`

In `compiler.Program.GetSourceFileForResolvedModule` (`program.go:1564-1578`):

```go
func (p *Program) GetSourceFileForResolvedModule(fileName string) *ast.SourceFile {
    file := p.GetSourceFile(fileName)     // map lookup, may return nil
    if file == nil {
        filename := p.GetParseFileRedirect(fileName)
        if filename != "" {
            return p.GetSourceFile(filename)
        }
    }
    return file  // may return nil, but caller (checker) handles nil
}
```

The checker at line 14806-14809 already handles a nil return:

```go
sourceFile = c.program.GetSourceFileForResolvedModule(resolvedModule.ResolvedFileName)
if sourceFile != nil {  // properly guards against nil
    // ...
}
```

So the checker is fine with `GetSourceFileForResolvedModule` returning nil.
The problem is that `aliasResolver.GetSourceFile` crashes *internally* before
it even gets a chance to return nil.

### Why Module Resolution Succeeds But File Doesn't Exist

Module resolution (`aliasResolver.GetResolvedModule`) uses a `module.Resolver`
which resolves based on filesystem state (checking `FileExists`, reading 
`package.json`, etc.). However, the `aliasResolver` is used during snapshot 
transitions — when the user is actively editing code (e.g., pasting into a git
diff viewer). Between:

1. Module resolution checking `FileExists` (which reads from the VFS)
2. `GetSourceFile` trying to read/parse the file

...the snapshot state may have changed, or the file may exist at the VFS layer
but not be loadable via the `SnapshotFS` file handle mechanism (e.g., the file
exists on disk but the overlay/cache hasn't populated it).

There's also a second possible scenario: module resolution resolved to a `.d.ts`
or other file that exists on disk but whose path normalization differs from what
`SnapshotFS.GetFile` expects.

## Fix

Add nil guards in both `GetSourceFile` and the `updateIndexes` call site:

### `aliasresolver.go:GetSourceFile`
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

### `registry.go:updateIndexes` (line ~783)
```go
sourceFile := aliasResolver.GetSourceFile(source.fileName)
if sourceFile == nil {
    return true
}
```

## Impact

- The fix is minimal and defensive
- The checker already handles nil returns from `GetSourceFileForResolvedModule`
  (see checker.go:14809), so returning nil is the correct behavior
- The `extractFromFile` function also doesn't handle nil (would crash on 
  `file.Symbol` access), so the guard in `updateIndexes` is also needed
