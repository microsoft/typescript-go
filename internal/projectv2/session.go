package projectv2

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2/ata"
	"github.com/microsoft/typescript-go/internal/projectv2/logging"
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
	DebounceDelay      time.Duration
}

type SessionInit struct {
	Options     *SessionOptions
	FS          vfs.FS
	Client      Client
	Logger      logging.Logger
	NpmExecutor ata.NpmExecutor
}

type Session struct {
	options                            *SessionOptions
	toPath                             func(string) tspath.Path
	client                             Client
	logger                             logging.Logger
	npmExecutor                        ata.NpmExecutor
	fs                                 *overlayFS
	parseCache                         *parseCache
	extendedConfigCache                *extendedConfigCache
	compilerOptionsForInferredProjects *core.CompilerOptions
	programCounter                     *programCounter
	typingsInstaller                   *ata.TypingsInstaller
	backgroundTasks                    *BackgroundQueue

	snapshotMu sync.RWMutex
	snapshot   *Snapshot

	pendingFileChangesMu sync.Mutex
	pendingFileChanges   []FileChange

	pendingATAChangesMu sync.Mutex
	pendingATAChanges   map[tspath.Path]*ATAStateChange

	// Debouncing fields for snapshot updates
	snapshotUpdateMu     sync.Mutex
	snapshotUpdateCancel context.CancelFunc
}

func NewSession(init *SessionInit) *Session {
	currentDirectory := init.Options.CurrentDirectory
	useCaseSensitiveFileNames := init.FS.UseCaseSensitiveFileNames()
	toPath := func(fileName string) tspath.Path {
		return tspath.ToPath(fileName, currentDirectory, useCaseSensitiveFileNames)
	}
	overlayFS := newOverlayFS(init.FS, make(map[tspath.Path]*overlay), init.Options.PositionEncoding, toPath)
	parseCache := &parseCache{options: tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: init.FS.UseCaseSensitiveFileNames(),
		CurrentDirectory:          init.Options.CurrentDirectory,
	}}
	extendedConfigCache := &extendedConfigCache{}

	session := &Session{
		options:             init.Options,
		toPath:              toPath,
		client:              init.Client,
		logger:              init.Logger,
		npmExecutor:         init.NpmExecutor,
		fs:                  overlayFS,
		parseCache:          parseCache,
		extendedConfigCache: extendedConfigCache,
		programCounter:      &programCounter{},
		backgroundTasks:     newBackgroundQueue(),
		snapshot: NewSnapshot(
			&snapshotFS{
				toPath: toPath,
				fs:     init.FS,
			},
			init.Options,
			parseCache,
			extendedConfigCache,
			&ConfigFileRegistry{},
			nil,
			toPath,
		),
		pendingATAChanges: make(map[tspath.Path]*ATAStateChange),
	}

	session.typingsInstaller = ata.NewTypingsInstaller(&ata.TypingsInstallerOptions{
		TypingsLocation: init.Options.TypingsLocation,
		ThrottleLimit:   5,
	}, session)

	return session
}

// FS implements module.ResolutionHost
func (s *Session) FS() vfs.FS {
	return s.fs.fs
}

// GetCurrentDirectory implements module.ResolutionHost
func (s *Session) GetCurrentDirectory() string {
	return s.options.CurrentDirectory
}

// Trace implements module.ResolutionHost
func (s *Session) Trace(msg string) {
	panic("ATA module resolution should not use tracing")
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
	changes, overlays := s.flushChangesLocked(ctx)
	s.pendingFileChangesMu.Unlock()
	s.UpdateSnapshot(ctx, overlays, SnapshotChange{
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

func (s *Session) DidChangeFile(ctx context.Context, uri lsproto.DocumentUri, version int32, changes []lsproto.TextDocumentContentChangePartialOrWholeDocument) {
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
	s.pendingFileChanges = append(s.pendingFileChanges, fileChanges...)
	s.pendingFileChangesMu.Unlock()

	// Schedule a debounced snapshot update
	s.ScheduleSnapshotUpdate()
}

func (s *Session) DidChangeCompilerOptionsForInferredProjects(ctx context.Context, options *core.CompilerOptions) {
	s.compilerOptionsForInferredProjects = options
	s.UpdateSnapshot(ctx, s.fs.Overlays(), SnapshotChange{
		compilerOptionsForInferredProjects: options,
	})
}

// ScheduleSnapshotUpdate schedules a debounced snapshot update.
// If there's already a pending update, it will be cancelled and a new one scheduled.
// This is useful for batching rapid changes like file watch events.
func (s *Session) ScheduleSnapshotUpdate() {
	s.snapshotUpdateMu.Lock()
	defer s.snapshotUpdateMu.Unlock()

	// Cancel any existing scheduled update
	if s.snapshotUpdateCancel != nil {
		s.snapshotUpdateCancel()
		s.logger.Log("Delaying scheduled snapshot update...")
	} else {
		s.logger.Log("Scheduling new snapshot update...")
	}

	// Create a new cancellable context for the debounce task
	debounceCtx, cancel := context.WithCancel(context.Background())
	s.snapshotUpdateCancel = cancel

	// Enqueue the debounced snapshot update
	s.backgroundTasks.Enqueue(debounceCtx, func(ctx context.Context) {
		// Sleep for the debounce delay
		select {
		case <-time.After(s.options.DebounceDelay):
			// Delay completed, proceed with update
		case <-ctx.Done():
			// Context was cancelled, newer events arrived
			return
		}

		// Clear the cancel function since we're about to execute the update
		s.snapshotUpdateMu.Lock()
		s.snapshotUpdateCancel = nil
		s.snapshotUpdateMu.Unlock()

		// Process the accumulated changes
		changeSummary, overlays, ataChanges := s.flushChanges(context.Background())
		if !changeSummary.IsEmpty() || len(ataChanges) > 0 {
			if s.options.LoggingEnabled {
				s.logger.Log("Running scheduled snapshot update")
			}
			s.UpdateSnapshot(context.Background(), overlays, SnapshotChange{
				fileChanges: changeSummary,
				ataChanges:  ataChanges,
			})
		} else if s.options.LoggingEnabled {
			s.logger.Log("Scheduled snapshot update skipped (no changes)")
		}
	})
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
	fileChanges, overlays, ataChanges := s.flushChanges(ctx)
	updateSnapshot := !fileChanges.IsEmpty() || len(ataChanges) > 0
	if updateSnapshot {
		// If there are pending file changes, we need to update the snapshot.
		// Sending the requested URI ensures that the project for this URI is loaded.
		snapshot = s.UpdateSnapshot(ctx, overlays, SnapshotChange{
			fileChanges:   fileChanges,
			ataChanges:    ataChanges,
			requestedURIs: []lsproto.DocumentUri{uri},
		})
	} else {
		// If there are no pending file changes, we can try to use the current snapshot.
		s.snapshotMu.RLock()
		snapshot = s.snapshot
		s.snapshotMu.RUnlock()
	}

	project := snapshot.GetDefaultProject(uri)
	if project == nil && !updateSnapshot || project != nil && project.dirty {
		// The current snapshot does not have an up to date project for the URI,
		// so we need to update the snapshot to ensure the project is loaded.
		// !!! Allow multiple projects to update in parallel
		snapshot = s.UpdateSnapshot(ctx, overlays, SnapshotChange{requestedURIs: []lsproto.DocumentUri{uri}})
		project = snapshot.GetDefaultProject(uri)
	}
	if project == nil {
		return nil, fmt.Errorf("no project found for URI %s", uri)
	}
	return ls.NewLanguageService(project, snapshot.Converters()), nil
}

func (s *Session) UpdateSnapshot(ctx context.Context, overlays map[tspath.Path]*overlay, change SnapshotChange) *Snapshot {
	// Cancel any pending scheduled update since we're doing an immediate update
	s.snapshotUpdateMu.Lock()
	if s.snapshotUpdateCancel != nil {
		s.logger.Log("Canceling scheduled snapshot update and performing one now")
		s.snapshotUpdateCancel()
		s.snapshotUpdateCancel = nil
	}
	s.snapshotUpdateMu.Unlock()

	s.snapshotMu.Lock()
	oldSnapshot := s.snapshot
	newSnapshot := oldSnapshot.Clone(ctx, change, overlays, s)
	s.snapshot = newSnapshot
	shouldDispose := newSnapshot != oldSnapshot && oldSnapshot.Deref()
	s.snapshotMu.Unlock()
	if shouldDispose {
		oldSnapshot.dispose(s)
	}

	// Enqueue ATA updates if needed
	if s.npmExecutor != nil {
		s.triggerATAForUpdatedProjects(newSnapshot)
	}

	// Enqueue logging, watch updates, and diagnostic refresh tasks
	s.backgroundTasks.Enqueue(context.Background(), func(ctx context.Context) {
		if s.options.LoggingEnabled {
			s.logger.Write(newSnapshot.builderLogs.String())
			s.logProjectChanges(oldSnapshot, newSnapshot)
			s.logger.Write("")
		}
		if s.options.WatchEnabled {
			if err := s.updateWatches(oldSnapshot, newSnapshot); err != nil && s.options.LoggingEnabled {
				s.logger.Log(err)
			}
		}
		if change.fileChanges.IncludesWatchChangesOnly {
			if err := s.client.RefreshDiagnostics(context.Background()); err != nil && s.options.LoggingEnabled {
				s.logger.Log(fmt.Sprintf("Error refreshing diagnostics: %v", err))
			}
		}
	})

	return newSnapshot
}

// WaitForBackgroundTasks waits for all background tasks to complete.
// This is intended to be used only for testing purposes.
func (s *Session) WaitForBackgroundTasks() {
	s.backgroundTasks.WaitForEmpty()
}

func updateWatch[T any](ctx context.Context, client Client, logger logging.Logger, oldWatcher, newWatcher *WatchedFiles[T]) []error {
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
	// Cancel any pending snapshot update
	s.snapshotUpdateMu.Lock()
	if s.snapshotUpdateCancel != nil {
		s.snapshotUpdateCancel()
		s.snapshotUpdateCancel = nil
	}
	s.snapshotUpdateMu.Unlock()
	s.backgroundTasks.Close()
}

func (s *Session) flushChanges(ctx context.Context) (FileChangeSummary, map[tspath.Path]*overlay, map[tspath.Path]*ATAStateChange) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingATAChangesMu.Lock()
	defer s.pendingATAChangesMu.Unlock()
	pendingATAChanges := s.pendingATAChanges
	s.pendingATAChanges = make(map[tspath.Path]*ATAStateChange)
	fileChanges, overlays := s.flushChangesLocked(ctx)
	return fileChanges, overlays, pendingATAChanges
}

// flushChangesLocked should only be called with s.pendingFileChangesMu held.
func (s *Session) flushChangesLocked(ctx context.Context) (FileChangeSummary, map[tspath.Path]*overlay) {
	if len(s.pendingFileChanges) == 0 {
		return FileChangeSummary{}, s.fs.Overlays()
	}

	start := time.Now()
	changes, overlays := s.fs.processChanges(s.pendingFileChanges)
	if s.options.LoggingEnabled {
		s.logger.Log(fmt.Sprintf("Processed %d file changes in %v", len(s.pendingFileChanges), time.Since(start)))
	}
	s.pendingFileChanges = nil
	return changes, overlays
}

// logProjectChanges logs information about projects that have changed between snapshots
func (s *Session) logProjectChanges(oldSnapshot *Snapshot, newSnapshot *Snapshot) {
	logProject := func(project *Project) {
		var builder strings.Builder
		project.print(s.logger.IsVerbose() /*writeFileNames*/, s.logger.IsVerbose() /*writeFileExplanation*/, &builder)
		s.logger.Log(builder.String())
	}
	core.DiffMaps(
		oldSnapshot.ProjectCollection.configuredProjects,
		newSnapshot.ProjectCollection.configuredProjects,
		func(path tspath.Path, addedProject *Project) {
			// New project added
			logProject(addedProject)
		},
		func(path tspath.Path, removedProject *Project) {
			// Project removed
			s.logger.Logf("\nProject '%s' removed\n%s", removedProject.Name(), hr)
		},
		func(path tspath.Path, oldProject, newProject *Project) {
			// Project updated
			if newProject.ProgramUpdateKind == ProgramUpdateKindNewFiles {
				logProject(newProject)
			}
		},
	)

	oldInferred := oldSnapshot.ProjectCollection.inferredProject
	newInferred := newSnapshot.ProjectCollection.inferredProject

	if oldInferred != nil && newInferred == nil {
		// Inferred project removed
		s.logger.Logf("\nProject '%s' removed\n%s", oldInferred.Name(), hr)
	} else if newInferred != nil && newInferred.ProgramUpdateKind == ProgramUpdateKindNewFiles {
		// Inferred project updated
		logProject(newInferred)
	}
}

func (s *Session) NpmInstall(cwd string, npmInstallArgs []string) ([]byte, error) {
	return s.npmExecutor.NpmInstall(cwd, npmInstallArgs)
}

func (s *Session) triggerATAForUpdatedProjects(newSnapshot *Snapshot) {
	for _, project := range newSnapshot.ProjectCollection.Projects() {
		if project.ShouldTriggerATA() {
			s.backgroundTasks.Enqueue(context.Background(), func(ctx context.Context) {
				var logTree *logging.LogTree
				if s.options.LoggingEnabled {
					logTree = logging.NewLogTree(fmt.Sprintf("Triggering ATA for project %s", project.Name()))
				}

				typingsInfo := project.ComputeTypingsInfo()
				request := &ata.TypingsInstallRequest{
					ProjectID:        project.configFilePath,
					TypingsInfo:      &typingsInfo,
					FileNames:        core.Map(project.Program.GetSourceFiles(), func(file *ast.SourceFile) string { return file.FileName() }),
					ProjectRootPath:  project.currentDirectory,
					CompilerOptions:  project.CommandLine.CompilerOptions(),
					CurrentDirectory: s.options.CurrentDirectory,
					GetScriptKind:    core.GetScriptKindFromFileName,
					FS:               s.fs.fs,
					Logger:           logTree,
				}

				if typingsFiles, err := s.typingsInstaller.InstallTypings(request); err != nil && logTree != nil {
					s.logger.Log(fmt.Sprintf("ATA installation failed for project %s: %v", project.Name(), err))
					s.logger.Log(logTree.String())
				} else {
					s.pendingATAChangesMu.Lock()
					defer s.pendingATAChangesMu.Unlock()
					s.pendingATAChanges[project.configFilePath] = &ATAStateChange{
						TypingsInfo:  &typingsInfo,
						TypingsFiles: typingsFiles,
						Logs:         logTree,
					}
				}
			})
		}
	}
}
