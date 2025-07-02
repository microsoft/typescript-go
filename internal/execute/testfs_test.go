package execute_test

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type testFsTrackingLibs struct {
	vfs.FS
	defaultLibs *collections.SyncSet[string]
}

func NewFSTrackingLibs(fs vfs.FS) *testFsTrackingLibs {
	return &testFsTrackingLibs{FS: fs}
}

func (f *testFsTrackingLibs) removeIgnoreLibPath(path string) {
	if f.defaultLibs != nil && f.defaultLibs.Has(path) {
		f.defaultLibs.Delete(path)
	}
}

// ReadFile reads the file specified by path and returns the content.
// If the file fails to be read, ok will be false.
func (f *testFsTrackingLibs) ReadFile(path string) (contents string, ok bool) {
	f.removeIgnoreLibPath(path)
	return f.FS.ReadFile(path)
}

func (f *testFsTrackingLibs) WriteFile(path string, data string, writeByteOrderMark bool) error {
	f.removeIgnoreLibPath(path)
	return f.FS.WriteFile(path, data, writeByteOrderMark)
}

// Removes `path` and all its contents. Will return the first error it encounters.
func (f *testFsTrackingLibs) Remove(path string) error {
	f.removeIgnoreLibPath(path)
	return f.FS.Remove(path)
}
