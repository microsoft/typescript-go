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

// FS is a file system abstraction.
type FS interface {
	// UseCaseSensitiveFileNames returns true if the file system is case-sensitive.
	UseCaseSensitiveFileNames() bool

	// GetCurrentDirectory returns the current directory.
	GetCurrentDirectory() string

	// FileExists returns true if the file exists.
	FileExists(path string) bool

	// ReadFile reads the file specified by path and returns the content.
	// If the file fails to be read, ok will be false.
	ReadFile(path string) (contents string, ok bool)

	// DirectoryExists returns true if the path is a directory.
	DirectoryExists(path string) bool

	// GetDirectories returns the names of the directories in the specified directory.
	GetDirectories(path string) []string

	// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
	// It is has the same behavior as [fs.WalkDir], but with paths as [string].
	WalkDir(root string, walkFn WalkDirFunc) error
}

// DirEntry is [fs.DirEntry].
type DirEntry = fs.DirEntry

// WalkDirFunc is [fs.WalkDirFunc].
type WalkDirFunc = fs.WalkDirFunc

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
		cwd:                       cwd,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		rootFor:                   os.DirFS,
	}
}

// FromIOFS creates a new FS from an [fs.FS].
// For paths like `c:/foo/bar`, fsys will be used as though it's rooted at `/` and the path is `/c:/foo/bar`.
func FromIOFS(cwd string, useCaseSensitiveFileNames bool, fsys fs.FS) FS {
	if cwd == "" {
		panic("cwd must be provided")
	}

	return &adapter{
		readSema:                  osReadSema,
		cwd:                       cwd,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		rootFor: func(root string) fs.FS {
			if root == "/" {
				return fsys
			}

			p := tspath.RemoveTrailingDirectorySeparator(root)
			sub, err := fs.Sub(fsys, p)
			if err != nil {
				panic(err)
			}
			return sub
		},
	}
}

var osReadSema = make(chan struct{}, 128)

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
	readSema chan struct{}

	cwd                       string
	useCaseSensitiveFileNames bool

	rootFor func(root string) fs.FS
}

func (a *adapter) UseCaseSensitiveFileNames() bool {
	return a.useCaseSensitiveFileNames
}

func (a *adapter) GetCurrentDirectory() string {
	return a.cwd
}

func splitPath(p string) (rootName, rest string) {
	l := tspath.GetEncodedRootLength(string(p))
	if l < 0 {
		panic("FS does not support URLs")
	}
	return string(p[:l]), string(p[l:])
}

func (a *adapter) rootAndPath(path string) (fs.FS, string) {
	rootName, rest := splitPath(path)
	return a.rootFor(rootName), rest
}

func (a *adapter) stat(path string) fs.FileInfo {
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

func (a *adapter) FileExists(path string) bool {
	stat := a.stat(path)
	return stat != nil && !stat.IsDir()
}

func (a *adapter) ReadFile(path string) (contents string, ok bool) {
	if a.readSema != nil {
		a.readSema <- struct{}{}
		defer func() { <-a.readSema }()
	}

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

func (a *adapter) DirectoryExists(path string) bool {
	stat := a.stat(path)
	return stat != nil && stat.IsDir()
}

func (a *adapter) GetDirectories(path string) []string {
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

func (a *adapter) WalkDir(root string, walkFn WalkDirFunc) error {
	fsys, rest := a.rootAndPath(root)
	if fsys == nil {
		return nil
	}
	return fs.WalkDir(fsys, rest, walkFn)
}
