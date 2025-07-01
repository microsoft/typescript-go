package execute_test

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type testFsTrackingLibs struct {
	fs          vfs.FS
	defaultLibs *collections.SyncSet[string]
}

var _ vfs.FS = (*testFsTrackingLibs)(nil)

func NewFSTrackingLibs(fs vfs.FS) *testFsTrackingLibs {
	return &testFsTrackingLibs{
		fs: fs,
	}
}

func (f *testFsTrackingLibs) FS() vfs.FS {
	return f.fs
}

func (f *testFsTrackingLibs) removeIgnoreLibPath(path string) {
	if f.defaultLibs != nil && f.defaultLibs.Has(path) {
		f.defaultLibs.Delete(path)
	}
}

func (f *testFsTrackingLibs) UseCaseSensitiveFileNames() bool {
	return f.fs.UseCaseSensitiveFileNames()
}

// FileExists returns true if the file exists.
func (f *testFsTrackingLibs) FileExists(path string) bool {
	return f.fs.FileExists(path)
}

// ReadFile reads the file specified by path and returns the content.
// If the file fails to be read, ok will be false.
func (f *testFsTrackingLibs) ReadFile(path string) (contents string, ok bool) {
	f.removeIgnoreLibPath(path)
	return f.fs.ReadFile(path)
}

func (f *testFsTrackingLibs) WriteFile(path string, data string, writeByteOrderMark bool) error {
	f.removeIgnoreLibPath(path)
	return f.fs.WriteFile(path, data, writeByteOrderMark)
}

// Removes `path` and all its contents. Will return the first error it encounters.
func (f *testFsTrackingLibs) Remove(path string) error {
	f.removeIgnoreLibPath(path)
	return f.fs.Remove(path)
}

// DirectoryExists returns true if the path is a directory.
func (f *testFsTrackingLibs) DirectoryExists(path string) bool {
	return f.fs.DirectoryExists(path)
}

// GetAccessibleEntries returns the files/directories in the specified directory.
// If any entry is a symlink, it will be followed.
func (f *testFsTrackingLibs) GetAccessibleEntries(path string) vfs.Entries {
	return f.fs.GetAccessibleEntries(path)
}

func (f *testFsTrackingLibs) Stat(path string) vfs.FileInfo {
	return f.fs.Stat(path)
}

// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
// It is has the same behavior as [fs.WalkDir], but with paths as [string].
func (f *testFsTrackingLibs) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return f.fs.WalkDir(root, walkFn)
}

// Realpath returns the "real path" of the specified path,
// following symlinks and correcting filename casing.
func (f *testFsTrackingLibs) Realpath(path string) string {
	return f.fs.Realpath(path)
}
