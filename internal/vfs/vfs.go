package vfs

import (
	"io/fs"
	"os"
	"runtime"
	"strings"
	"sync"
	"unicode"

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

func FromOS(cwd string) FS {
	if cwd == "" {
		panic("cwd must be provided")
	}

	useCaseSensitiveFileNames := isFileSystemCaseSensitive()
	return &adapter{
		cwd:                       tspath.ToPath(cwd, "", useCaseSensitiveFileNames),
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		rootFor:                   os.DirFS,
	}
}

// FromIOFS creates a new FS from an [fs.FS].
// For paths like `c:/foo/bar`, fsys will be used as though the path is `/c:/foo/bar`.
func FromIOFS(cwd string, useCaseSensitiveFileNames bool, fsys fs.FS) FS {
	if cwd == "" {
		panic("cwd must be provided")
	}

	return &adapter{
		cwd:                       tspath.ToPath(cwd, "", useCaseSensitiveFileNames),
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		rootFor: func(root string) fs.FS {
			if root == "/" {
				return fsys
			}

			sub, err := fs.Sub(fsys, tspath.RemoveTrailingDirectorySeparator(root))
			if err != nil {
				panic(err)
			}
			return sub
		},
	}
}

var isFileSystemCaseSensitive = sync.OnceValue(func() bool {
	// win32/win64 are case insensitive platforms
	if runtime.GOOS == "windows" {
		return false
	}

	// If the current executable exists under a different case, we must be case-insensitve.
	if _, err := os.Stat(swapCase(os.Args[0])); os.IsNotExist(err) {
		return false
	}
	return true
})

// Convert all lowercase chars to uppercase, and vice-versa
func swapCase(str string) string {
	return strings.Map(func(r rune) rune {
		upper := unicode.ToUpper(r)
		if upper == r {
			return unicode.ToLower(r)
		} else {
			return upper
		}
	}, str)
}

type adapter struct {
	cwd                       tspath.Path
	useCaseSensitiveFileNames bool

	rootFor func(root string) fs.FS
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

func (a *adapter) rootAndPath(path tspath.Path) (fs.FS, string) {
	rootName, rest := splitPath(path)
	return a.rootFor(rootName), rest
}

func (a *adapter) stat(path tspath.Path) fs.FileInfo {
	fsys, rest := a.rootAndPath(path)
	if fsys == nil {
		return nil
	}
	stat, err := fs.Stat(fsys, rest)
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
	fsys, rest := a.rootAndPath(path)
	if fsys == nil {
		return nil
	}

	entries, err := fs.ReadDir(fsys, rest)
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
	fsys, rest := a.rootAndPath(path)
	if fsys == nil {
		return ""
	}
	contents, err := fs.ReadFile(fsys, rest)
	if err != nil {
		return ""
	}
	return string(contents)
}
