package projectv2

import (
	"cmp"
	"context"
	"slices"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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
	sessionOptions *SessionOptions
	parseCache     *parseCache

	overlayFS          *overlayFS
	compilerFS         *compilerFS
	configuredProjects map[tspath.Path]*Project
	inferredProject    *Project
}

// NewSnapshot
func NewSnapshot(
	fs vfs.FS,
	overlays map[lsproto.DocumentUri]*overlay,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
) *Snapshot {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	id := snapshotID.Add(1)
	s := &Snapshot{
		id:                 id,
		sessionOptions:     sessionOptions,
		overlayFS:          newOverlayFS(cachedFS, overlays),
		parseCache:         parseCache,
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
	return s.overlayFS.getFile(uri)
}

func (s *Snapshot) Projects() []*Project {
	projects := make([]*Project, 0, len(s.configuredProjects)+1)
	for _, p := range s.configuredProjects {
		projects = append(projects, p)
	}
	slices.SortFunc(projects, func(a, b *Project) int {
		return cmp.Compare(a.Name, b.Name)
	})
	if s.inferredProject != nil {
		projects = append(projects, s.inferredProject)
	}
	return projects
}

func (s *Snapshot) GetDefaultProject(uri lsproto.DocumentUri) *Project {
	// !!!
	fileName := ls.DocumentURIToFileName(uri)
	path := s.toPath(fileName)
	for _, p := range s.Projects() {
		if p.containsFile(path) {
			return p
		}
	}
	return nil
}

type snapshotChange struct {
	// changedURIs are URIs that have changed since the last snapshot.
	changedURIs collections.Set[lsproto.DocumentUri]
	// requestedURIs are URIs that were requested by the client.
	// The new snapshot should ensure projects for these URIs have loaded programs.
	requestedURIs []lsproto.DocumentUri
}

func (c snapshotChange) toProjectChange(snapshot *Snapshot) projectChange {
	changedURIs := make([]tspath.Path, c.changedURIs.Len())
	requestedURIs := make([]struct {
		path           tspath.Path
		defaultProject *Project
	}, len(c.requestedURIs))
	for i, uri := range c.requestedURIs {
		requestedURIs[i] = struct {
			path           tspath.Path
			defaultProject *Project
		}{
			path:           snapshot.toPath(ls.DocumentURIToFileName(uri)),
			defaultProject: snapshot.GetDefaultProject(uri),
		}
	}
	return projectChange{
		changedURIs:   changedURIs,
		requestedURIs: requestedURIs,
	}
}

func (s *Snapshot) Clone(ctx context.Context, change snapshotChange, session *Session) *Snapshot {
	newSnapshot := NewSnapshot(
		session.fs.fs,
		session.fs.overlays,
		s.sessionOptions,
		s.parseCache,
	)

	projectChange := change.toProjectChange(s)

	for configFilePath, project := range s.configuredProjects {
		newProject, _ := project.Clone(ctx, projectChange, newSnapshot)
		if newProject != nil {
			newSnapshot.configuredProjects[configFilePath] = newProject
		}
	}
	// !!! update inferred project if needed
	return newSnapshot
}

func (s *Snapshot) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, s.sessionOptions.CurrentDirectory, s.sessionOptions.UseCaseSensitiveFileNames)
}
