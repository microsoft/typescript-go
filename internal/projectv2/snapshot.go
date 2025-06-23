package projectv2

import (
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
)

var snapshotID atomic.Uint64

type StateChangeKind int

const (
	StateChangeKindFile StateChangeKind = iota
	StateChangeKindProgramLoad
)

var _ vfs.FS = (*compilerFS)(nil)

type compilerFS struct {
	snapshot *Snapshot
}

// DirectoryExists implements vfs.FS.
func (fs *compilerFS) DirectoryExists(path string) bool {
	return fs.snapshot.overlayFS.fs.DirectoryExists(path)
}

// FileExists implements vfs.FS.
func (fs *compilerFS) FileExists(path string) bool {
	if fh := fs.snapshot.GetFile(ls.FileNameToDocumentURI(path)); fh != nil {
		return true
	}
	return fs.snapshot.overlayFS.fs.FileExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (fs *compilerFS) GetAccessibleEntries(path string) vfs.Entries {
	return fs.snapshot.overlayFS.fs.GetAccessibleEntries(path)
}

// ReadFile implements vfs.FS.
func (fs *compilerFS) ReadFile(path string) (contents string, ok bool) {
	if fh := fs.snapshot.GetFile(ls.FileNameToDocumentURI(path)); fh != nil {
		return fh.Content(), true
	}
	return "", false
}

// Realpath implements vfs.FS.
func (fs *compilerFS) Realpath(path string) string {
	return fs.snapshot.overlayFS.fs.Realpath(path)
}

// Stat implements vfs.FS.
func (fs *compilerFS) Stat(path string) vfs.FileInfo {
	return fs.snapshot.overlayFS.fs.Stat(path)
}

// UseCaseSensitiveFileNames implements vfs.FS.
func (fs *compilerFS) UseCaseSensitiveFileNames() bool {
	return fs.snapshot.overlayFS.fs.UseCaseSensitiveFileNames()
}

// WalkDir implements vfs.FS.
func (fs *compilerFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	panic("unimplemented")
}

// WriteFile implements vfs.FS.
func (fs *compilerFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	panic("unimplemented")
}

// Remove implements vfs.FS.
func (fs *compilerFS) Remove(path string) error {
	panic("unimplemented")
}

type Snapshot struct {
	id uint64

	// Session options are immutable for the server lifetime,
	// so can be a pointer.
	sessionOptions *SessionOptions
	parseCache     *parseCache

	overlayFS          *overlayFS
	compilerFS         *compilerFS
	files              *fileMap
	configuredProjects map[tspath.Path]*Project
}

func NewSnapshot(
	fs vfs.FS,
	overlays *overlays,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
) *Snapshot {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	id := snapshotID.Add(1)
	s := &Snapshot{
		id:                 id,
		sessionOptions:     sessionOptions,
		overlayFS:          newOverlayFSFromOverlays(cachedFS, overlays),
		parseCache:         parseCache,
		files:              newFileMap(),
		configuredProjects: make(map[tspath.Path]*Project),
	}

	s.compilerFS = &compilerFS{snapshot: s}

	return s
}

// GetFile is stable over the lifetime of the snapshot. It first consults its
// own cache (which includes keys for missing files), and only delegates to the
// file system if the key is not known to the cache. GetFile respects the state
// of overlays.
func (s *Snapshot) GetFile(uri lsproto.DocumentUri) fileHandle {
	if f, ok := s.files.Get(uri); ok {
		return f // may be nil, a file known to be missing
	}
	fh := s.overlayFS.getFile(uri)
	s.files.Set(uri, fh)
	return fh
}
