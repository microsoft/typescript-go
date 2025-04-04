package cachedvfs

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/vfs"
)

type FS struct {
	fs vfs.FS

	directoryExistsCache      sync.Map // map[string]bool
	fileExistsCache           sync.Map // map[string]bool
	getAccessibleEntriesCache sync.Map // map[string]vfs.Entries
	realpathCache             sync.Map // map[string]string
	statCache                 sync.Map // map[string]vfs.FileInfo
}

var _ vfs.FS = (*FS)(nil)

func From(fs vfs.FS) *FS {
	return &FS{fs: fs}
}

func (vfs *FS) ClearCache() {
	vfs.directoryExistsCache.Clear()
	vfs.fileExistsCache.Clear()
	vfs.getAccessibleEntriesCache.Clear()
	vfs.realpathCache.Clear()
	vfs.statCache.Clear()
}

func (vfs *FS) DirectoryExists(path string) bool {
	return cached(&vfs.directoryExistsCache, path, vfs.fs.DirectoryExists)
}

func (vfs *FS) FileExists(path string) bool {
	return cached(&vfs.fileExistsCache, path, vfs.fs.FileExists)
}

func (vfs *FS) GetAccessibleEntries(path string) vfs.Entries {
	return cached(&vfs.getAccessibleEntriesCache, path, vfs.fs.GetAccessibleEntries)
}

func (vfs *FS) ReadFile(path string) (contents string, ok bool) {
	return vfs.fs.ReadFile(path)
}

func (vfs *FS) Realpath(path string) string {
	return cached(&vfs.realpathCache, path, vfs.fs.Realpath)
}

func (vfs *FS) Remove(path string) error {
	// TODO: should this call ClearCache?
	return vfs.fs.Remove(path)
}

func (vfs *FS) Stat(path string) vfs.FileInfo {
	return cached(&vfs.statCache, path, vfs.fs.Stat)
}

func (vfs *FS) UseCaseSensitiveFileNames() bool {
	return vfs.fs.UseCaseSensitiveFileNames()
}

func (vfs *FS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return vfs.fs.WalkDir(root, walkFn)
}

func (vfs *FS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	// TODO: should this call ClearCache?
	return vfs.fs.WriteFile(path, data, writeByteOrderMark)
}

func cached[Arg any, Ret any](cache *sync.Map, key Arg, fn func(Arg) Ret) Ret {
	if ret, ok := cache.Load(key); ok {
		return ret.(Ret)
	}
	ret := fn(key)
	cache.Store(key, ret)
	return ret
}
