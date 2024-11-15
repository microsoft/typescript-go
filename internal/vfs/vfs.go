package vfs

import (
	"bytes"
	"encoding/binary"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf16"

	"github.com/microsoft/typescript-go/internal/tspath"
)

// FS is a file system abstraction. All paths are handled as [tspath.Path].
type FS interface {
	// UseCaseSensitiveFileNames returns true if the file system is case-sensitive.
	UseCaseSensitiveFileNames() bool

	// GetCurrentDirectory returns the current directory.
	GetCurrentDirectory() tspath.Path

	// ToPath converts a string to a path. If the path is relative, it is resolved against the current directory.
	ToPath(p string) tspath.Path

	// FileExists returns true if the file exists.
	FileExists(path tspath.Path) bool

	// ReadFile reads the file specified by path and returns the content.
	// If the file fails to be read, ok will be false.
	ReadFile(path tspath.Path) (contents string, ok bool)

	// DirectoryExists returns true if the path is a directory.
	DirectoryExists(path tspath.Path) bool

	// GetDirectories returns the names of the directories in the specified directory.
	GetDirectories(path tspath.Path) []string

	// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
	// It is has the same behavior as [fs.WalkDir], but with paths as [tspath.Path].
	WalkDir(root tspath.Path, walkFn WalkDirFunc) error
}

// DirEntry is [fs.DirEntry].
type DirEntry = fs.DirEntry

// WalkDirFunc is [fs.WalkDirFunc] but with paths as [tspath.Path].
type WalkDirFunc func(path tspath.Path, d DirEntry, err error) error

var (
	// SkipAll is [fs.SkipAll].
	SkipAll = fs.SkipAll

	// SkipDir is [fs.SkipDir].
	SkipDir = fs.SkipDir
)

var _ FS = (*adapter)(nil)

// FromOS creates a new FS from the OS file system.
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

func (a *adapter) ReadFile(path tspath.Path) (contents string, ok bool) {
	fsys, rest := a.rootAndPath(path)
	if fsys == nil {
		return "", false
	}

	b, err := fs.ReadFile(fsys, rest)
	if err != nil {
		return "", false
	}

	var bom [2]byte
	if len(b) >= 2 {
		bom = [2]byte{b[0], b[1]}
		switch bom {
		case [2]byte{0xFF, 0xFE}:
			return decodeUtf16(b[2:], binary.LittleEndian), true
		case [2]byte{0xFE, 0xFF}:
			return decodeUtf16(b[2:], binary.BigEndian), true
		}
	}
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		b = b[3:]
	}

	return string(b), true
}

func decodeUtf16(b []byte, order binary.ByteOrder) string {
	ints := make([]uint16, len(b)/2)
	if err := binary.Read(bytes.NewReader(b), order, &ints); err != nil {
		return ""
	}
	return string(utf16.Decode(ints))
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

	// TODO: should this really exist? ReadDir with manual filtering seems like a better idea.
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs
}

func (a *adapter) WalkDir(root tspath.Path, walkFn WalkDirFunc) error {
	fsys, rest := a.rootAndPath(root)
	if fsys == nil {
		return nil
	}

	toPath := a.ToPath
	return fs.WalkDir(fsys, rest, func(path string, d fs.DirEntry, err error) error {
		return walkFn(toPath(path), d, err)
	})
}
