package projectv2

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync/atomic"

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
	sessionOptions *SessionOptions
	parseCache     *parseCache
	logger         *project.Logger

	// Immutable state, cloned between snapshots
	overlayFS          *overlayFS
	compilerFS         *compilerFS
	projectCollection  *projectCollection
	configFileRegistry *configFileRegistry
}

// NewSnapshot
func NewSnapshot(
	fs vfs.FS,
	overlays map[tspath.Path]*overlay,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
	logger *project.Logger,
	configFileRegistry *configFileRegistry,
) *Snapshot {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	id := snapshotID.Add(1)
	s := &Snapshot{
		id: id,

		sessionOptions:     sessionOptions,
		parseCache:         parseCache,
		logger:             logger,
		configFileRegistry: configFileRegistry,
		projectCollection:  &projectCollection{},

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

func (s *Snapshot) ConfiguredProjects() []*Project {
	projects := make([]*Project, 0, len(s.projectCollection.configuredProjects))
	s.fillConfiguredProjects(&projects)
	return projects
}

func (s *Snapshot) fillConfiguredProjects(projects *[]*Project) {
	for _, p := range s.projectCollection.configuredProjects {
		*projects = append(*projects, p)
	}
	slices.SortFunc(*projects, func(a, b *Project) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

func (s *Snapshot) Projects() []*Project {
	if s.projectCollection.inferredProject == nil {
		return s.ConfiguredProjects()
	}
	projects := make([]*Project, 0, len(s.projectCollection.configuredProjects)+1)
	s.fillConfiguredProjects(&projects)
	projects = append(projects, s.projectCollection.inferredProject)
	return projects
}

func (s *Snapshot) GetDefaultProject(uri lsproto.DocumentUri) *Project {
	fileName := ls.DocumentURIToFileName(uri)
	path := s.toPath(fileName)
	var (
		containingProjects                       []*Project
		firstConfiguredProject                   *Project
		firstNonSourceOfProjectReferenceRedirect *Project
		multipleDirectInclusions                 bool
	)
	for _, p := range s.ConfiguredProjects() {
		if p.containsFile(path) {
			containingProjects = append(containingProjects, p)
			if !multipleDirectInclusions && !p.IsSourceFromProjectReference(path) {
				if firstNonSourceOfProjectReferenceRedirect == nil {
					firstNonSourceOfProjectReferenceRedirect = p
				} else {
					multipleDirectInclusions = true
				}
			}
			if firstConfiguredProject == nil {
				firstConfiguredProject = p
			}
		}
	}
	if len(containingProjects) == 1 {
		return containingProjects[0]
	}
	if len(containingProjects) == 0 {
		if s.projectCollection.inferredProject != nil && s.projectCollection.inferredProject.containsFile(path) {
			return s.projectCollection.inferredProject
		}
		return nil
	}
	if !multipleDirectInclusions {
		if firstNonSourceOfProjectReferenceRedirect != nil {
			// Multiple projects include the file, but only one is a direct inclusion.
			return firstNonSourceOfProjectReferenceRedirect
		}
		// Multiple projects include the file, and none are direct inclusions.
		return firstConfiguredProject
	}
	// Multiple projects include the file directly.
	// !!! temporary!
	builder := newProjectCollectionBuilder(context.Background(), s, s.projectCollection, s.configFileRegistry, snapshotChange{})
	defer func() {
		p, c := builder.finalize()
		if p != s.projectCollection || c != s.configFileRegistry {
			panic("temporary builder should have collected no changes for a find lookup")
		}
	}()

	if defaultProject := builder.findDefaultConfiguredProject(fileName, path); defaultProject != nil {
		return defaultProject
	}
	return firstConfiguredProject
}

type snapshotChange struct {
	// fileChanges are the changes that have occurred since the last snapshot.
	fileChanges FileChangeSummary
	// requestedURIs are URIs that were requested by the client.
	// The new snapshot should ensure projects for these URIs have loaded programs.
	requestedURIs []lsproto.DocumentUri
}

func (s *Snapshot) Clone(ctx context.Context, change snapshotChange, session *Session) *Snapshot {
	newSnapshot := NewSnapshot(
		session.fs.fs,
		session.fs.overlays,
		s.sessionOptions,
		s.parseCache,
		s.logger,
		nil,
	)

	projectCollectionBuilder := newProjectCollectionBuilder(
		ctx,
		newSnapshot,
		s.projectCollection,
		s.configFileRegistry,
		change,
	)

	for uri := range change.fileChanges.Opened.Keys() {
		fileName := uri.FileName()
		path := s.toPath(fileName)
		projectCollectionBuilder.tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(fileName, path, projectLoadKindCreate)
	}

	newSnapshot.projectCollection, newSnapshot.configFileRegistry = projectCollectionBuilder.finalize()

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
