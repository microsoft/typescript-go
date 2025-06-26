package incrementaltestutil

import (
	"encoding/json"
	"fmt"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var fakeTsVersion = "FakeTSVersion"

type FsHandlingBuildInfo struct {
	fs vfs.FS
}

var _ vfs.FS = (*FsHandlingBuildInfo)(nil)

func NewFsHandlingBuildInfo(fs vfs.FS) *FsHandlingBuildInfo {
	return &FsHandlingBuildInfo{
		fs: fs,
	}
}

func (f *FsHandlingBuildInfo) FS() vfs.FS {
	return f.fs
}

func (f *FsHandlingBuildInfo) UseCaseSensitiveFileNames() bool {
	return f.fs.UseCaseSensitiveFileNames()
}

// FileExists returns true if the file exists.
func (f *FsHandlingBuildInfo) FileExists(path string) bool {
	return f.fs.FileExists(path)
}

// ReadFile reads the file specified by path and returns the content.
// If the file fails to be read, ok will be false.
func (f *FsHandlingBuildInfo) ReadFile(path string) (contents string, ok bool) {
	contents, ok = f.fs.ReadFile(path)
	if ok && tspath.FileExtensionIs(path, tspath.ExtensionTsBuildInfo) {
		// read buildinfo and modify version
		var buildInfo incremental.BuildInfo
		err := json.Unmarshal([]byte(contents), &buildInfo)
		if err == nil && buildInfo.Version == fakeTsVersion {
			buildInfo.Version = core.Version()
			newContents, err := json.Marshal(&buildInfo)
			if err != nil {
				panic("testFs.ReadFile: failed to marshal build info after fixing version: " + err.Error())
			}
			contents = string(newContents)
		}
	}
	return contents, ok
}

func (f *FsHandlingBuildInfo) WriteFile(path string, data string, writeByteOrderMark bool) error {
	if tspath.FileExtensionIs(path, tspath.ExtensionTsBuildInfo) {
		var buildInfo incremental.BuildInfo
		err := json.Unmarshal([]byte(data), &buildInfo)
		if err == nil && buildInfo.Version == core.Version() {
			// Change it to fakeTsVersion
			buildInfo.Version = fakeTsVersion
			newData, err := json.Marshal(&buildInfo)
			if err != nil {
				return fmt.Errorf("testFs.WriteFile: failed to marshal build info after fixing version: %w", err)
			}
			data = string(newData)
		}
		if err == nil {
			// Write readable build info version
			f.fs.WriteFile(path+".readable.baseline.txt", toReadableBuildInfo(&buildInfo, data), false)
		}
	}
	return f.fs.WriteFile(path, data, writeByteOrderMark)
}

// Removes `path` and all its contents. Will return the first error it encounters.
func (f *FsHandlingBuildInfo) Remove(path string) error {
	return f.fs.Remove(path)
}

// DirectoryExists returns true if the path is a directory.
func (f *FsHandlingBuildInfo) DirectoryExists(path string) bool {
	return f.fs.DirectoryExists(path)
}

// GetAccessibleEntries returns the files/directories in the specified directory.
// If any entry is a symlink, it will be followed.
func (f *FsHandlingBuildInfo) GetAccessibleEntries(path string) vfs.Entries {
	return f.fs.GetAccessibleEntries(path)
}

func (f *FsHandlingBuildInfo) Stat(path string) vfs.FileInfo {
	return f.fs.Stat(path)
}

// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
// It is has the same behavior as [fs.WalkDir], but with paths as [string].
func (f *FsHandlingBuildInfo) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return f.fs.WalkDir(root, walkFn)
}

// Realpath returns the "real path" of the specified path,
// following symlinks and correcting filename casing.
func (f *FsHandlingBuildInfo) Realpath(path string) string {
	return f.fs.Realpath(path)
}
