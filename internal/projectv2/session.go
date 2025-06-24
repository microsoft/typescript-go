package projectv2

import (
	"context"
	"sync"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

type SessionOptions struct {
	tspath.ComparePathsOptions
	DefaultLibraryPath string
	TypingsLocation    string
	PositionEncoding   lsproto.PositionEncodingKind
	WatchEnabled       bool
	NewLine            string
}

type Session struct {
	options    SessionOptions
	fs         *overlayFS
	parseCache *parseCache
	converters *ls.Converters

	snapshotMu sync.Mutex
	snapshot   *Snapshot

	pendingFileChangesMu sync.Mutex
	pendingFileChanges   []FileChange
}

func NewSession(options SessionOptions) *Session {
	overlayFS := newOverlayFS(bundled.WrapFS(osvfs.FS()), make(map[lsproto.DocumentUri]*overlay))
	parseCache := &parseCache{options: options.ComparePathsOptions}
	converters := ls.NewConverters(options.PositionEncoding, func(fileName string) *ls.LineMap {
		// !!! cache
		return ls.ComputeLineStarts(overlayFS.getFile(ls.FileNameToDocumentURI(fileName)).Content())
	})

	return &Session{
		options:    options,
		fs:         overlayFS,
		parseCache: parseCache,
		converters: converters,
	}
}

func (s *Session) DidOpenFile(ctx context.Context, uri lsproto.DocumentUri, version int32, content string, languageKind lsproto.LanguageKind) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind:         FileChangeKindOpen,
		URI:          uri,
		Version:      version,
		Content:      content,
		LanguageKind: languageKind,
	})
	s.fs.updateOverlays(s.pendingFileChanges, s.converters)
	s.pendingFileChanges = nil
	// !!! update snapshot
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

func (s *Session) GetLanguageService(ctx context.Context, uri lsproto.DocumentUri) *ls.LanguageService {
	changed := s.flushChanges(ctx)
}

func (s *Session) flushChanges(ctx context.Context) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()

	if len(s.pendingFileChanges) == 0 {
		return
	}

	changed := s.fs.updateOverlays(s.pendingFileChanges, s.converters)
	s.pendingFileChanges = nil
	return changed
}
