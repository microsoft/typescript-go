package projectv2

import (
	"context"
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

type SessionOptions struct {
	CurrentDirectory   string
	DefaultLibraryPath string
	TypingsLocation    string
	PositionEncoding   lsproto.PositionEncodingKind
	WatchEnabled       bool
	NewLine            string
}

type Session struct {
	options    SessionOptions
	fs         *overlayFS
	logger     *project.Logger
	parseCache *parseCache
	converters *ls.Converters

	snapshotMu sync.Mutex
	snapshot   *Snapshot

	pendingFileChangesMu sync.Mutex
	pendingFileChanges   []FileChange
}

func NewSession(options SessionOptions, fs vfs.FS, logger *project.Logger) *Session {
	overlayFS := newOverlayFS(bundled.WrapFS(osvfs.FS()), make(map[tspath.Path]*overlay))
	parseCache := &parseCache{options: tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: fs.UseCaseSensitiveFileNames(),
		CurrentDirectory:          options.CurrentDirectory,
	}}
	converters := ls.NewConverters(options.PositionEncoding, func(fileName string) *ls.LineMap {
		// !!! cache
		return ls.ComputeLineStarts(overlayFS.getFile(ls.FileNameToDocumentURI(fileName)).Content())
	})

	return &Session{
		options:    options,
		fs:         overlayFS,
		logger:     logger,
		parseCache: parseCache,
		converters: converters,
		snapshot: NewSnapshot(
			overlayFS.fs,
			overlayFS.overlays,
			&options,
			parseCache,
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
	s.UpdateSnapshot(ctx, snapshotChange{
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

func (s *Session) GetLanguageService(ctx context.Context, uri lsproto.DocumentUri) (*ls.LanguageService, error) {
	changes := s.flushChanges(ctx)
	if !changes.IsEmpty() {
		s.UpdateSnapshot(ctx, snapshotChange{
			fileChanges:   changes,
			requestedURIs: []lsproto.DocumentUri{uri},
		})
	}

	project := s.snapshot.GetDefaultProject(uri)
	if project == nil {
		return nil, fmt.Errorf("no project found for URI %s", uri)
	}
	if project.LanguageService == nil {
		panic("project language service is nil")
	}
	return project.LanguageService, nil
}

func (s *Session) UpdateSnapshot(ctx context.Context, change snapshotChange) {
	s.snapshotMu.Lock()
	defer s.snapshotMu.Unlock()
	s.snapshot = s.snapshot.Clone(ctx, change, s)
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

	changes := s.fs.processChanges(s.pendingFileChanges, s.converters)
	s.pendingFileChanges = nil
	return changes
}
