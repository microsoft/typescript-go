package projectv2

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
)

var snapshotID atomic.Uint64

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
	sessionOptions      *SessionOptions
	parseCache          *parseCache
	extendedConfigCache *extendedConfigCache
	logger              *project.Logger

	// Immutable state, cloned between snapshots
	overlayFS                          *overlayFS
	compilerFS                         *compilerFS
	projectCollection                  *ProjectCollection
	configFileRegistry                 *ConfigFileRegistry
	compilerOptionsForInferredProjects *core.CompilerOptions
}

// NewSnapshot
func NewSnapshot(
	fs vfs.FS,
	overlays map[tspath.Path]*overlay,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
	extendedConfigCache *extendedConfigCache,
	logger *project.Logger,
	configFileRegistry *ConfigFileRegistry,
	compilerOptionsForInferredProjects *core.CompilerOptions,
) *Snapshot {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	id := snapshotID.Add(1)
	s := &Snapshot{
		id: id,

		sessionOptions:      sessionOptions,
		parseCache:          parseCache,
		extendedConfigCache: extendedConfigCache,
		logger:              logger,
		configFileRegistry:  configFileRegistry,
		projectCollection:   &ProjectCollection{},

		overlayFS: newOverlayFS(cachedFS, sessionOptions.PositionEncoding, overlays),
	}

	s.compilerFS = &compilerFS{snapshot: s}

	return s
}

// GetFile is stable over the lifetime of the snapshot. It first consults its
// own cache (which includes keys for missing files), and only delegates to the
// file system if the key is not known to the cache. GetFile respects the state
// of overlays.
func (s *Snapshot) GetFile(uri lsproto.DocumentUri) fileHandle {
	return s.overlayFS.getFile(uri)
}

func (s *Snapshot) Overlays() map[tspath.Path]*overlay {
	return s.overlayFS.overlays
}

func (s *Snapshot) IsOpenFile(path tspath.Path) bool {
	// An open file is one that has an overlay.
	_, ok := s.overlayFS.overlays[path]
	return ok
}

func (s *Snapshot) GetDefaultProject(uri lsproto.DocumentUri) *Project {
	fileName := ls.DocumentURIToFileName(uri)
	path := s.toPath(fileName)
	return s.projectCollection.GetDefaultProject(fileName, path)
}

type snapshotChange struct {
	// fileChanges are the changes that have occurred since the last snapshot.
	fileChanges FileChangeSummary
	// requestedURIs are URIs that were requested by the client.
	// The new snapshot should ensure projects for these URIs have loaded programs.
	requestedURIs []lsproto.DocumentUri
	// compilerOptionsForInferredProjects is the compiler options to use for inferred projects.
	// It should only be set the value in the next snapshot should be changed. If nil, the
	// value from the previous snapshot will be copied to the new snapshot.
	compilerOptionsForInferredProjects *core.CompilerOptions
}

func (s *Snapshot) Clone(ctx context.Context, change snapshotChange, session *Session) *Snapshot {
	newSnapshot := NewSnapshot(
		session.fs.fs,
		session.fs.overlays,
		s.sessionOptions,
		s.parseCache,
		s.extendedConfigCache,
		s.logger,
		nil,
		s.compilerOptionsForInferredProjects,
	)

	if change.compilerOptionsForInferredProjects != nil {
		// !!! mark inferred projects as dirty?
		newSnapshot.compilerOptionsForInferredProjects = change.compilerOptionsForInferredProjects
	}

	projectCollectionBuilder := newProjectCollectionBuilder(
		ctx,
		newSnapshot,
		s.projectCollection,
		s.configFileRegistry,
	)

	for uri := range change.fileChanges.Opened.Keys() {
		projectCollectionBuilder.DidOpenFile(uri)
	}

	projectCollectionBuilder.DidChangeFiles(slices.Collect(maps.Keys(change.fileChanges.Changed.M)))

	for _, uri := range change.requestedURIs {
		projectCollectionBuilder.DidRequestFile(uri)
	}

	newSnapshot.projectCollection, newSnapshot.configFileRegistry = projectCollectionBuilder.Finalize()

	return newSnapshot
}

func (s *Snapshot) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, s.sessionOptions.CurrentDirectory, s.overlayFS.fs.UseCaseSensitiveFileNames())
}

func (s *Snapshot) Log(msg string) {
	s.logger.Info(msg)
}

func (s *Snapshot) Logf(format string, args ...any) {
	s.logger.Info(fmt.Sprintf(format, args...))
}
