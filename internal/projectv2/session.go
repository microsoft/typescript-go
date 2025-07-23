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
	LoggingEnabled     bool
}

type Session struct {
	options                            SessionOptions
	toPath                             func(string) tspath.Path
	client                             Client
	logger                             Logger
	fs                                 *overlayFS
	parseCache                         *parseCache
	extendedConfigCache                *extendedConfigCache
	compilerOptionsForInferredProjects *core.CompilerOptions
	programCounter                     *programCounter

	snapshotMu sync.RWMutex
	snapshot   *Snapshot

	pendingFileChangesMu sync.Mutex
	pendingFileChanges   []FileChange
}

func NewSession(options SessionOptions, fs vfs.FS, client Client, logger Logger) *Session {
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
		toPath:              toPath,
		client:              client,
		logger:              logger,
		fs:                  overlayFS,
		parseCache:          parseCache,
		extendedConfigCache: extendedConfigCache,
		programCounter:      &programCounter{},
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
		Hash: s.fs.getFile(uri.FileName()).Hash(),
	})
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

func (s *Session) DidChangeWatchedFiles(ctx context.Context, changes []*lsproto.FileEvent) {
	fileChanges := make([]FileChange, 0, len(changes))
	for _, change := range changes {
		var kind FileChangeKind
		switch change.Type {
		case lsproto.FileChangeTypeCreated:
			kind = FileChangeKindWatchCreate
		case lsproto.FileChangeTypeChanged:
			kind = FileChangeKindWatchChange
		case lsproto.FileChangeTypeDeleted:
			kind = FileChangeKindWatchDelete
		default:
			continue // Ignore unknown change types.
		}
		fileChanges = append(fileChanges, FileChange{
			Kind: kind,
			URI:  change.Uri,
		})
	}
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, fileChanges...)
}

func (s *Session) Snapshot() (*Snapshot, func()) {
	s.snapshotMu.RLock()
	defer s.snapshotMu.RUnlock()
	snapshot := s.snapshot
	snapshot.Ref()
	return snapshot, func() {
		if snapshot.Deref() {
			// The session itself accounts for one reference to the snapshot, and it derefs
			// in UpdateSnapshot while holding the snapshotMu lock, so the only way to end
			// up here is for an external caller to release the snapshot after the session
			// has already dereferenced it and moved to a new snapshot. In other words, we
			// can assume that `snapshot != s.snapshot`, and therefor there's no way for
			// anyone else to acquire a reference to this snapshot again.
			snapshot.dispose(s)
		}
	}
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
		// !!! Allow multiple projects to update in parallel
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
	oldSnapshot := s.snapshot
	newSnapshot := oldSnapshot.Clone(ctx, change, s)
	s.snapshot = newSnapshot
	shouldDispose := newSnapshot != oldSnapshot && oldSnapshot.Deref()
	s.snapshotMu.Unlock()
	if shouldDispose {
		oldSnapshot.dispose(s)
	}
	go func() {
		if s.options.LoggingEnabled {
			newSnapshot.builderLogs.WriteLogs(s.logger)
			s.logger.Log("")
		}
		if s.options.WatchEnabled {
			if err := s.updateWatches(oldSnapshot, newSnapshot); err != nil && s.options.LoggingEnabled {
				s.logger.Log(err)
			}
		}
	}()
	return newSnapshot
}

func updateWatch[T any](ctx context.Context, client Client, logger Logger, oldWatcher, newWatcher *WatchedFiles[T]) []error {
	var errors []error
	if newWatcher != nil {
		if id, watchers := newWatcher.Watchers(); len(watchers) > 0 {
			if err := client.WatchFiles(ctx, id, watchers); err != nil {
				errors = append(errors, err)
			}
			if logger != nil {
				if oldWatcher == nil {
					logger.Log(fmt.Sprintf("Added new watch: %s", id))
				} else {
					logger.Log(fmt.Sprintf("Updated watch: %s", id))
				}
				for _, watcher := range watchers {
					logger.Log(fmt.Sprintf("\t%s", *watcher.GlobPattern.Pattern))
				}
				logger.Log("")
			}
		}
	}
	if oldWatcher != nil {
		if id, watchers := oldWatcher.Watchers(); len(watchers) > 0 {
			if err := client.UnwatchFiles(ctx, id); err != nil {
				errors = append(errors, err)
			}
			if logger != nil && newWatcher == nil {
				logger.Log(fmt.Sprintf("Removed watch: %s", id))
			}
		}
	}
	return errors
}

func (s *Session) updateWatches(oldSnapshot *Snapshot, newSnapshot *Snapshot) error {
	var errors []error
	ctx := context.Background()
	core.DiffMapsFunc(
		oldSnapshot.ConfigFileRegistry.configs,
		newSnapshot.ConfigFileRegistry.configs,
		func(a, b *configFileEntry) bool {
			return a.rootFilesWatch.ID() == b.rootFilesWatch.ID()
		},
		func(_ tspath.Path, addedEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s.client, s.logger, nil, addedEntry.rootFilesWatch)...)
		},
		func(_ tspath.Path, removedEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s.client, s.logger, removedEntry.rootFilesWatch, nil)...)
		},
		func(_ tspath.Path, oldEntry, newEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s.client, s.logger, oldEntry.rootFilesWatch, newEntry.rootFilesWatch)...)
		},
	)

	core.DiffMaps(
		oldSnapshot.ProjectCollection.configuredProjects,
		newSnapshot.ProjectCollection.configuredProjects,
		func(_ tspath.Path, addedProject *Project) {
			errors = append(errors, updateWatch(ctx, s.client, s.logger, nil, addedProject.affectingLocationsWatch)...)
			errors = append(errors, updateWatch(ctx, s.client, s.logger, nil, addedProject.failedLookupsWatch)...)
		},
		func(_ tspath.Path, removedProject *Project) {
			errors = append(errors, updateWatch(ctx, s.client, s.logger, removedProject.affectingLocationsWatch, nil)...)
			errors = append(errors, updateWatch(ctx, s.client, s.logger, removedProject.failedLookupsWatch, nil)...)
		},
		func(_ tspath.Path, oldProject, newProject *Project) {
			if oldProject.affectingLocationsWatch.ID() != newProject.affectingLocationsWatch.ID() {
				errors = append(errors, updateWatch(ctx, s.client, s.logger, oldProject.affectingLocationsWatch, newProject.affectingLocationsWatch)...)
			}
			if oldProject.failedLookupsWatch.ID() != newProject.failedLookupsWatch.ID() {
				errors = append(errors, updateWatch(ctx, s.client, s.logger, oldProject.failedLookupsWatch, newProject.failedLookupsWatch)...)
			}
		},
	)

	if len(errors) > 0 {
		return fmt.Errorf("errors updating watches: %v", errors)
	}
	return nil
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
