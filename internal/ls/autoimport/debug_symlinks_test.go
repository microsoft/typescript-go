package autoimport

import (
    "strings"
    "testing"
    
    "github.com/microsoft/typescript-go/internal/collections"
    "github.com/microsoft/typescript-go/internal/symlinks"
    "github.com/microsoft/typescript-go/internal/tspath"
)

func TestDebugSymlinkCacheFromResolution(t *testing.T) {
    cache := symlinks.NewKnownSymlink("/", true)
    
    // Simulate the resolution of @packages/foo from bar
    // These are what module resolution would produce:
    // originalPath: the path through the symlink (in node_modules)
    // resolvedFileName: the realpath (the actual file location)
    originalPath := "/home/src/workspaces/project/node_modules/@packages/foo/src/index.ts"
    resolvedFileName := "/home/src/workspaces/project/packages/foo/src/index.ts"
    
    t.Logf("Processing resolution:")
    t.Logf("  originalPath: %s", originalPath)
    t.Logf("  resolvedFileName: %s", resolvedFileName)
    
    cache.ProcessResolution(originalPath, resolvedFileName)
    
    // Check directoriesByRealpath
    directoriesByRealpath := cache.DirectoriesByRealpath()
    
    t.Logf("Directories by realpath:")
    directoriesByRealpath.Range(func(realPath tspath.Path, symlinkSet *collections.SyncSet[string]) bool {
        t.Logf("  RealPath: %s", realPath)
        symlinkSet.Range(func(symlinkPath string) bool {
            t.Logf("    -> Symlink: %s", symlinkPath)
            return true
        })
        return true
    })
    
    // Now check if we can find the symlink for an internal file
    filePath := tspath.Path("/home/src/workspaces/project/packages/foo/src/internal/index.ts")
    
    // Simulate hasSymlinkToNodeModules
    found := false
    tspath.ForEachAncestorDirectoryPath(filePath, func(dirPath tspath.Path) (any, bool) {
        key := dirPath.EnsureTrailingDirectorySeparator()
        symlinkPaths, ok := directoriesByRealpath.Load(key)
        if !ok {
            return nil, false
        }
        symlinkPaths.Range(func(symlinkPath string) bool {
            if strings.Contains(symlinkPath, "/node_modules/") {
                found = true
                return false
            }
            return true
        })
        return nil, found
    })
    
    if !found {
        t.Error("Expected to find symlink to node_modules but didn't")
    } else {
        t.Log("SUCCESS: Found symlink to node_modules for internal file")
    }
}
