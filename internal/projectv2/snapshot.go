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
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
)

var snapshotID atomic.Uint64

// !!! create some type safety for this to ensure caching
func newSnapshotFS(fs vfs.FS, overlays map[tspath.Path]*overlay, positionEncoding lsproto.PositionEncodingKind, toPath func(string) tspath.Path) *overlayFS {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	return newOverlayFS(cachedFS, overlays, positionEncoding, toPath)
}

type Snapshot struct {
	id uint64

	// Session options are immutable for the server lifetime,
	// so can be a pointer.
	sessionOptions *SessionOptions
	toPath         func(fileName string) tspath.Path

	// Immutable state, cloned between snapshots
	overlayFS                          *overlayFS
	ProjectCollection                  *ProjectCollection
	ConfigFileRegistry                 *ConfigFileRegistry
	compilerOptionsForInferredProjects *core.CompilerOptions
	builderLogs                        *logCollector
}

// NewSnapshot
func NewSnapshot(
	fs *overlayFS,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
	extendedConfigCache *extendedConfigCache,
	configFileRegistry *ConfigFileRegistry,
	compilerOptionsForInferredProjects *core.CompilerOptions,
	toPath func(fileName string) tspath.Path,
) *Snapshot {

	id := snapshotID.Add(1)
	s := &Snapshot{
		id: id,

		sessionOptions: sessionOptions,
		toPath:         toPath,

		overlayFS:                          fs,
		ConfigFileRegistry:                 configFileRegistry,
		ProjectCollection:                  &ProjectCollection{},
		compilerOptionsForInferredProjects: compilerOptionsForInferredProjects,
	}

	return s
}

func (s *Snapshot) GetDefaultProject(uri lsproto.DocumentUri) *Project {
	fileName := ls.DocumentURIToFileName(uri)
	path := s.toPath(fileName)
	return s.ProjectCollection.GetDefaultProject(fileName, path)
}

func (s *Snapshot) ID() uint64 {
	return s.id
}

func (s *Snapshot) GetFile(uri lsproto.DocumentUri) FileHandle {
	fileName := ls.DocumentURIToFileName(uri)
	return s.overlayFS.getFile(fileName)
}

type SnapshotChange struct {
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

func (s *Snapshot) Clone(ctx context.Context, change SnapshotChange, session *Session) *Snapshot {
	var logger *logCollector
	if session.options.LoggingEnabled {
		var close func()
		logger, close = NewLogCollector(fmt.Sprintf("Cloning snapshot %d", s.id))
		defer close()
	}

	fs := newSnapshotFS(session.fs.fs, session.fs.overlays, session.fs.positionEncoding, s.toPath)
	compilerOptionsForInferredProjects := s.compilerOptionsForInferredProjects
	if change.compilerOptionsForInferredProjects != nil {
		// !!! mark inferred projects as dirty?
		compilerOptionsForInferredProjects = change.compilerOptionsForInferredProjects
	}

	projectCollectionBuilder := newProjectCollectionBuilder(
		ctx,
		fs,
		s.ProjectCollection,
		s.ConfigFileRegistry,
		compilerOptionsForInferredProjects,
		s.sessionOptions,
		session.parseCache,
		session.extendedConfigCache,
		logger,
	)

	for file, hash := range change.fileChanges.Closed {
		projectCollectionBuilder.DidCloseFile(file, hash)
	}

	for uri := range change.fileChanges.Opened.Keys() {
		projectCollectionBuilder.DidOpenFile(uri)
	}

	projectCollectionBuilder.DidChangeFiles(slices.Collect(maps.Keys(change.fileChanges.Changed.M)))

	for _, uri := range change.requestedURIs {
		projectCollectionBuilder.DidRequestFile(uri)
	}

	newSnapshot := NewSnapshot(
		fs,
		s.sessionOptions,
		session.parseCache,
		session.extendedConfigCache,
		nil,
		s.compilerOptionsForInferredProjects,
		s.toPath,
	)

	newSnapshot.ProjectCollection, newSnapshot.ConfigFileRegistry = projectCollectionBuilder.Finalize()

	return newSnapshot
}
