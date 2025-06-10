package modulespecifiers

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type SymlinkedDirectory struct {
	/**
	 * Matches the casing returned by `realpath`.  Used to compute the `realpath` of children.
	 * Always has trailing directory separator
	 */
	Real string
	/**
	 * toPath(real).  Stored to avoid repeated recomputation.
	 * Always has trailing directory separator
	 */
	RealPath tspath.Path
}

type SymlinkCache struct {
	symlinkedDirectories collections.SyncMap[tspath.Path, *SymlinkedDirectory]
	symlinkedFiles       collections.SyncMap[tspath.Path, string]
}

/** Gets a map from symlink to realpath. Keys have trailing directory separators. */
func (cache *SymlinkCache) SymlinkedDirectories() *collections.SyncMap[tspath.Path, *SymlinkedDirectory] {
	return &cache.symlinkedDirectories
}

/** Gets a map from symlink to realpath */
func (cache *SymlinkCache) SymlinkedFiles() *collections.SyncMap[tspath.Path, string] {
	return &cache.symlinkedFiles
}

// all callers should check !containsIgnoredPath(symlinkPath)
func (cache *SymlinkCache) SetSymlinkedDirectory(symlink string, symlinkPath tspath.Path, realDirectory *SymlinkedDirectory) {
	// Large, interconnected dependency graphs in pnpm will have a huge number of symlinks
	// where both the realpath and the symlink path are inside node_modules/.pnpm. Since
	// this path is never a candidate for a module specifier, we can ignore it entirely.

	// !!!
	// if realDirectory != nil {
	// 	if _, ok := cache.symlinkedDirectories.Load(symlinkPath); !ok {
	// 		cache.symlinkedDirectoriesByRealpath.Add(realDirectory.RealPath, symlink)
	// 	}
	// }
	cache.symlinkedDirectories.Store(symlinkPath, realDirectory)
}

func (cache *SymlinkCache) SetSymlinkedFile(symlinkPath tspath.Path, realpath string) {
	cache.symlinkedFiles.Store(symlinkPath, realpath)
}
