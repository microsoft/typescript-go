package vfs

import (
	"io/fs"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type FS interface {
	UseCaseSensitiveFileNames() bool

	GetCurrentDirectory() tspath.Path
	ToPath(p string) tspath.Path

	FileExists(path tspath.Path) bool
	ReadFile(path tspath.Path) string

	DirectoryExists(path tspath.Path) bool
	GetDirectories(path tspath.Path) []string
}

var _ FS = (*adapter)(nil)

type adapter struct {
	cwd                       tspath.Path
	useCaseSensitiveFileNames bool

	// Map from the "encoded root" to the FS that should be used for that root.
	roots map[string]fs.FS
}

func (a *adapter) UseCaseSensitiveFileNames() bool {
	return a.useCaseSensitiveFileNames
}

func (a *adapter) GetCurrentDirectory() tspath.Path {
	return a.cwd
}

func (a *adapter) ToPath(p string) tspath.Path {
	return tspath.ToPath(p, string(a.cwd), a.useCaseSensitiveFileNames)
}

func splitPath(p tspath.Path) (rootName, rest string) {
	l := tspath.GetEncodedRootLength(string(p))
	if l < 0 {
		panic("FS does not support URLs")
	}
	return string(p[:l]), string(p[l:])
}

func (a *adapter) rootFor(path tspath.Path) (fs.FS, string) {
	rootName, rest := splitPath(path)
	return a.roots[rootName], rest
}

func (a *adapter) stat(path tspath.Path) fs.FileInfo {
	root, rest := a.rootFor(path)
	if root == nil {
		return nil
	}
	stat, err := fs.Stat(root, rest)
	if err != nil {
		return nil
	}
	return stat
}

func (a *adapter) FileExists(path tspath.Path) bool {
	stat := a.stat(path)
	return stat != nil && !stat.IsDir()
}

func (a *adapter) DirectoryExists(path tspath.Path) bool {
	stat := a.stat(path)
	return stat != nil && stat.IsDir()
}

func (a *adapter) GetDirectories(path tspath.Path) []string {
	root, rest := a.rootFor(path)
	if root == nil {
		return nil
	}

	entries, err := fs.ReadDir(root, rest)
	if err != nil {
		return nil
	}

	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs
}

func (a *adapter) ReadFile(path tspath.Path) string {
	root, rest := a.rootFor(path)
	if root == nil {
		return ""
	}
	contents, err := fs.ReadFile(root, rest)
	if err != nil {
		return ""
	}
	return string(contents)
}
