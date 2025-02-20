package vfs

import (
	"io/fs"
)

// FS is a file system abstraction.
type FS interface {
	// UseCaseSensitiveFileNames returns true if the file system is case-sensitive.
	UseCaseSensitiveFileNames() bool

	// FileExists returns true if the file exists.
	FileExists(path string) bool

	// ReadFile reads the file specified by path and returns the content.
	// If the file fails to be read, ok will be false.
	ReadFile(path string) (contents string, ok bool)

	WriteFile(path string, data string, writeByteOrderMark bool) error

	// DirectoryExists returns true if the path is a directory.
	DirectoryExists(path string) bool

	// GetDirectories returns the names of the directories in the specified directory.
	GetDirectories(path string) []string

	// GetEntries returns the entries in the specified directory.
	GetEntries(path string) []fs.DirEntry

	// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
	// It is has the same behavior as [fs.WalkDir], but with paths as [string].
	WalkDir(root string, walkFn WalkDirFunc) error

	// Realpath returns the "real path" of the specified path,
	// following symlinks and correcting filename casing.
	Realpath(path string) string
}

// DirEntry is [fs.DirEntry].
type DirEntry = fs.DirEntry

var (
	ErrInvalid    = fs.ErrInvalid    // "invalid argument"
	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"
)

// WalkDirFunc is [fs.WalkDirFunc].
type WalkDirFunc = fs.WalkDirFunc

var (
	// SkipAll is [fs.SkipAll].
	SkipAll = fs.SkipAll //nolint:errname

	// SkipDir is [fs.SkipDir].
	SkipDir = fs.SkipDir //nolint:errname
)
