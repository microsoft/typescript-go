//go:build !noembed

package bundled

import (
	"io/fs"
	"slices"
	"strings"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
)

const embedded = true

const scheme = "bundled:///"

func splitPath(path string) (rest string, ok bool) {
	rest, ok = strings.CutPrefix(path, scheme)
	if !ok {
		return "", false
	}
	return rest, true
}

func libPath() string {
	return scheme + "libs"
}

// wrappedFS is implemented directly rather than going through [io/fs.FS].
// Our vfs.FS works with file contents in terms of strings, and that's
// what go:embed does under the hood, but going through fs.FS will cause
// copying to []byte and back.

type wrappedFS struct {
	fs vfs.FS
}

var _ vfs.FS = (*wrappedFS)(nil)

func wrapFS(fs vfs.FS) vfs.FS {
	return &wrappedFS{fs: fs}
}

func (vfs *wrappedFS) UseCaseSensitiveFileNames() bool {
	return vfs.fs.UseCaseSensitiveFileNames()
}

func (vfs *wrappedFS) FileExists(path string) bool {
	if rest, ok := splitPath(path); ok {
		_, ok := readEmbeddedFile(rest)
		return ok
	}
	return vfs.fs.FileExists(path)
}

func (vfs *wrappedFS) ReadFile(path string) (contents string, ok bool) {
	if rest, ok := splitPath(path); ok {
		return readEmbeddedFile(rest)
	}
	return vfs.fs.ReadFile(path)
}

func (vfs *wrappedFS) DirectoryExists(path string) bool {
	if rest, ok := splitPath(path); ok {
		return rest == "libs"
	}
	return vfs.fs.DirectoryExists(path)
}

func (vfs *wrappedFS) GetDirectories(path string) []string {
	if rest, ok := splitPath(path); ok {
		if rest == "" {
			return []string{"libs"}
		}
		return []string{}
	}
	return vfs.fs.GetDirectories(path)
}

var rootEntries = []fs.DirEntry{
	fs.FileInfoToDirEntry(&embeddedFileInfo{name: "libs", mode: fs.ModeDir}),
}

func (vfs *wrappedFS) GetEntries(path string) []fs.DirEntry {
	if rest, ok := splitPath(path); ok {
		if rest == "" {
			return slices.Clone(rootEntries)
		}
		if rest == "libs" {
			return slices.Clone(libsEntries)
		}
		return []fs.DirEntry{}
	}
	return vfs.fs.GetEntries(path)
}

//nolint:errorlint
func (vfs *wrappedFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	if rest, ok := splitPath(root); ok {
		if err := vfs.walkDir(rest, walkFn); err != nil {
			if err == fs.SkipAll {
				return nil
			}
			return err
		}
		return nil
	}
	return vfs.fs.WalkDir(root, walkFn)
}

//nolint:errorlint
func (vfs *wrappedFS) walkDir(rest string, walkFn vfs.WalkDirFunc) error {
	var entries []fs.DirEntry
	switch rest {
	case "":
		entries = rootEntries
	case "libs":
		entries = libsEntries
	default:
		return nil
	}

	for _, entry := range entries {
		name := rest + "/" + entry.Name()

		if err := walkFn(scheme+name, entry, nil); err != nil {
			if err == fs.SkipAll {
				return fs.SkipAll
			}
			if err == fs.SkipDir {
				continue
			}
			return err
		}
		if entry.IsDir() {
			if err := vfs.walkDir(name, walkFn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vfs *wrappedFS) Realpath(path string) string {
	if _, ok := splitPath(path); ok {
		return path
	}
	return vfs.fs.Realpath(path)
}

func (vfs *wrappedFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	if _, ok := splitPath(path); ok {
		panic("cannot write to embedded file system")
	}
	return vfs.fs.WriteFile(path, data, writeByteOrderMark)
}

type embeddedFileInfo struct {
	mode fs.FileMode
	name string
	size int64
}

var _ fs.FileInfo = (*embeddedFileInfo)(nil)

func (e *embeddedFileInfo) IsDir() bool {
	return e.mode.IsDir()
}

func (e *embeddedFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (e *embeddedFileInfo) Mode() fs.FileMode {
	return e.mode
}

func (e *embeddedFileInfo) Name() string {
	return e.name
}

func (e *embeddedFileInfo) Size() int64 {
	return e.size
}

func (e *embeddedFileInfo) Sys() any {
	return nil
}
