package projectv2

import (
	"context"
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type SessionOptions struct {
	CurrentDirectory   string
	DefaultLibraryPath string
	TypingsLocation    string
	PositionEncoding   lsproto.PositionEncodingKind
	WatchEnabled       bool
	NewLine            string
	LoggingEnabled     bool
}

type Session struct {
	options                            SessionOptions
	fs                                 *overlayFS
	parseCache                         *parseCache
	extendedConfigCache                *extendedConfigCache
	compilerOptionsForInferredProjects *core.CompilerOptions

	snapshotMu sync.RWMutex
	snapshot   *Snapshot

	pendingFileChangesMu sync.Mutex
	pendingFileChanges   []FileChange
}

func NewSession(options SessionOptions, fs vfs.FS) *Session {
	currentDirectory := options.CurrentDirectory
	useCaseSensitiveFileNames := fs.UseCaseSensitiveFileNames()
	toPath := func(fileName string) tspath.Path {
		return tspath.ToPath(fileName, currentDirectory, useCaseSensitiveFileNames)
	}
	overlayFS := newOverlayFS(fs, make(map[tspath.Path]*overlay), options.PositionEncoding, toPath)
	parseCache := &parseCache{options: tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: fs.UseCaseSensitiveFileNames(),
		CurrentDirectory:          options.CurrentDirectory,
	}}
	extendedConfigCache := &extendedConfigCache{}

	return &Session{
		options:             options,
		fs:                  overlayFS,
		parseCache:          parseCache,
		extendedConfigCache: extendedConfigCache,
		snapshot: NewSnapshot(
			newSnapshotFS(overlayFS.fs, overlayFS.overlays, options.PositionEncoding, toPath),
			&options,
			parseCache,
			extendedConfigCache,
			&ConfigFileRegistry{},
			nil,
			toPath,
		),
	}
}

func (s *Session) DidOpenFile(ctx context.Context, uri lsproto.DocumentUri, version int32, content string, languageKind lsproto.LanguageKind) {
	s.pendingFileChangesMu.Lock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind:         FileChangeKindOpen,
		URI:          uri,
		Version:      version,
		Content:      content,
		LanguageKind: languageKind,
	})
	changes := s.flushChangesLocked(ctx)
	s.pendingFileChangesMu.Unlock()
	s.UpdateSnapshot(ctx, SnapshotChange{
		fileChanges:   changes,
		requestedURIs: []lsproto.DocumentUri{uri},
	})
}

func (s *Session) DidCloseFile(ctx context.Context, uri lsproto.DocumentUri) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind: FileChangeKindClose,
		URI:  uri,
	})
	// !!! immediate update if file does not exist
}

func (s *Session) DidChangeFile(ctx context.Context, uri lsproto.DocumentUri, version int32, changes []lsproto.TextDocumentContentChangeEvent) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind:    FileChangeKindChange,
		URI:     uri,
		Version: version,
		Changes: changes,
	})
}

func (s *Session) DidSaveFile(ctx context.Context, uri lsproto.DocumentUri) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind: FileChangeKindSave,
		URI:  uri,
	})
}

// !!! ref count and release
func (s *Session) Snapshot() *Snapshot {
	s.snapshotMu.RLock()
	defer s.snapshotMu.RUnlock()
	return s.snapshot
}

func (s *Session) GetLanguageService(ctx context.Context, uri lsproto.DocumentUri) (*ls.LanguageService, error) {
	var snapshot *Snapshot
	changes := s.flushChanges(ctx)
	updateSnapshot := !changes.IsEmpty()
	if updateSnapshot {
		// If there are pending file changes, we need to update the snapshot.
		// Sending the requested URI ensures that the project for this URI is loaded.
		snapshot = s.UpdateSnapshot(ctx, SnapshotChange{
			fileChanges:   changes,
			requestedURIs: []lsproto.DocumentUri{uri},
		})
	} else {
		// If there are no pending file changes, we can try to use the current snapshot.
		s.snapshotMu.RLock()
		snapshot = s.snapshot
		s.snapshotMu.RUnlock()
	}

	project := snapshot.GetDefaultProject(uri)
	if project == nil && !updateSnapshot {
		// The current snapshot does not have the project for the URI,
		// so we need to update the snapshot to ensure the project is loaded.
		snapshot = s.UpdateSnapshot(ctx, SnapshotChange{requestedURIs: []lsproto.DocumentUri{uri}})
		project = snapshot.GetDefaultProject(uri)
	}
	if project == nil {
		return nil, fmt.Errorf("no project found for URI %s", uri)
	}
	if project.LanguageService == nil {
		panic("project language service is nil")
	}
	return project.LanguageService, nil
}

func (s *Session) UpdateSnapshot(ctx context.Context, change SnapshotChange) *Snapshot {
	s.snapshotMu.Lock()
	defer s.snapshotMu.Unlock()
	s.snapshot = s.snapshot.Clone(ctx, change, s)
	return s.snapshot
}

func (s *Session) Close() {
	// !!!
}

func (s *Session) flushChanges(ctx context.Context) FileChangeSummary {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	return s.flushChangesLocked(ctx)
}

// flushChangesLocked should only be called with s.pendingFileChangesMu held.
func (s *Session) flushChangesLocked(ctx context.Context) FileChangeSummary {
	if len(s.pendingFileChanges) == 0 {
		return FileChangeSummary{}
	}

	changes := s.fs.processChanges(s.pendingFileChanges)
	s.pendingFileChanges = nil
	return changes
}
