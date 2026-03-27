package project

import (
	"context"
	"fmt"
	"math"
	"runtime"
	gometrics "runtime/metrics"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	osmemory "github.com/mackerelio/go-osstat/memory"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/ata"
	"github.com/microsoft/typescript-go/internal/project/background"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type UpdateReason int

const (
	UpdateReasonUnknown UpdateReason = iota
	UpdateReasonDidOpenFile
	UpdateReasonDidChangeCompilerOptionsForInferredProjects
	UpdateReasonRequestedLanguageServicePendingChanges
	UpdateReasonRequestedLanguageServiceProjectNotLoaded
	UpdateReasonRequestedLanguageServiceForFileNotOpen
	UpdateReasonRequestedLanguageServiceProjectDirty
	UpdateReasonRequestedLoadProjectTree
	UpdateReasonRequestedLanguageServiceWithAutoImports
	UpdateReasonIdleCleanDiskCache
)

// SessionOptions are the immutable initialization options for a session.
// Snapshots may reference them as a pointer since they never change.
type SessionOptions struct {
	CurrentDirectory       string
	DefaultLibraryPath     string
	TypingsLocation        string
	PositionEncoding       lsproto.PositionEncodingKind
	WatchEnabled           bool
	LoggingEnabled         bool
	TelemetryEnabled       bool
	PushDiagnosticsEnabled bool
	DebounceDelay          time.Duration
	Locale                 locale.Locale
}

type SessionInit struct {
	BackgroundCtx context.Context
	Options       *SessionOptions
	FS            vfs.FS
	Client        Client
	Logger        logging.Logger
	NpmExecutor   ata.NpmExecutor
	ParseCache    *ParseCache
}

// Session manages the state of an LSP session. It receives textDocument
// events and requests for LanguageService objects from the LPS server
// and processes them into immutable snapshots as the data source for
// LanguageServices. When Session transitions from one snapshot to the
// next, it diffs them and updates file watchers and Automatic Type
// Acquisition (ATA) state accordingly.
type Session struct {
	backgroundCtx context.Context
	options       *SessionOptions
	toPath        func(string) tspath.Path
	client        Client
	logger        logging.Logger
	npmExecutor   ata.NpmExecutor
	fs            *overlayFS

	// parseCache is the ref-counted cache of source files used when
	// creating programs during snapshot cloning.
	parseCache *ParseCache
	// extendedConfigCache is the ref-counted cache of tsconfig ASTs
	// that are used in the "extends" of another tsconfig.
	extendedConfigCache *ExtendedConfigCache
	// programCounter counts how many snapshots reference a program.
	// When a program is no longer referenced, its source files are
	// released from the parseCache.
	programCounter *programCounter

	// read-only after initialization
	initialUserConfig *lsutil.UserConfig
	// current preferences
	workspaceUserConfig                *lsutil.UserConfig
	compilerOptionsForInferredProjects *core.CompilerOptions
	typingsInstaller                   *ata.TypingsInstaller
	backgroundQueue                    *background.Queue

	// snapshotID is the counter for snapshot IDs. It does not necessarily
	// equal the `snapshot.ID`. It is stored on Session instead of globally
	// so IDs are predictable in tests.
	snapshotID atomic.Uint64

	// snapshot is the current immutable state of all projects.
	snapshot         *Snapshot
	snapshotMu       sync.RWMutex
	snapshotUpdateMu sync.Mutex

	pendingUserConfigChanges bool
	userConfigRWMu           sync.Mutex

	// pendingFileChanges are accumulated from textDocument/* events delivered
	// by the LSP server through DidOpenFile(), DidChangeFile(), etc. They are
	// applied to the next snapshot update.
	pendingFileChanges   []FileChange
	pendingFileChangesMu sync.Mutex

	// pendingATAChanges are produced by Automatic Type Acquisition (ATA)
	// installations and applied to the next snapshot update.
	pendingATAChanges   map[tspath.Path]*ATAStateChange
	pendingATAChangesMu sync.Mutex

	// diagnosticsRefreshCancel is the cancelation function for a scheduled
	// diagnostics refresh. Diagnostics refreshes are scheduled and debounced
	// after file watch changes and ATA updates.
	diagnosticsRefreshCancel context.CancelFunc
	diagnosticsRefreshMu     sync.Mutex

	// idleCacheCleanTimer is a resettable timer for scheduling idle disk
	// cache cleans. The timer resets on any file event (open, close,
	// change, save, watch) and fires after 30 seconds of inactivity.
	idleCacheCleanTimer *time.Timer
	idleCacheCleanMu    sync.Mutex

	// performanceTelemetryCancel cancels the periodic performance telemetry ticker.
	performanceTelemetryCancel context.CancelFunc

	// seenProjects tracks projects that have already had telemetry sent.
	seenProjects collections.SyncSet[tspath.Path]

	// watches tracks the current watch globs and how many individual WatchedFiles
	// are using each glob.
	watches   map[fileSystemWatcherKey]*fileSystemWatcherValue
	watchesMu sync.Mutex
}

func NewSession(init *SessionInit) *Session {
	currentDirectory := init.Options.CurrentDirectory
	useCaseSensitiveFileNames := init.FS.UseCaseSensitiveFileNames()
	toPath := func(fileName string) tspath.Path {
		return tspath.ToPath(fileName, currentDirectory, useCaseSensitiveFileNames)
	}
	overlayFS := newOverlayFS(init.FS, make(map[tspath.Path]*Overlay), init.Options.PositionEncoding, toPath)
	parseCache := init.ParseCache
	if parseCache == nil {
		parseCache = NewParseCache(RefCountCacheOptions{})
	}
	extendedConfigCache := NewExtendedConfigCache()

	session := &Session{
		backgroundCtx:       init.BackgroundCtx,
		options:             init.Options,
		toPath:              toPath,
		client:              init.Client,
		logger:              init.Logger,
		npmExecutor:         init.NpmExecutor,
		fs:                  overlayFS,
		parseCache:          parseCache,
		extendedConfigCache: extendedConfigCache,
		programCounter:      &programCounter{},
		backgroundQueue:     background.NewQueue(),
		snapshot: NewSnapshot(
			uint64(0),
			&SnapshotFS{
				toPath: toPath,
				fs:     init.FS,
			},
			init.Options,
			&ConfigFileRegistry{},
			nil,
			nil,
			nil,
			NewWatchedFiles(
				"auto-import",
				lsproto.WatchKindCreate|lsproto.WatchKindChange|lsproto.WatchKindDelete,
				func(nodeModulesDirs map[tspath.Path]string) PatternsAndIgnored {
					patterns := make([]string, 0, len(nodeModulesDirs))
					for _, dir := range nodeModulesDirs {
						patterns = append(patterns, getRecursiveGlobPattern(dir))
					}
					slices.Sort(patterns)
					return PatternsAndIgnored{
						patterns: patterns,
					}
				},
			),
			toPath,
		),
		initialUserConfig:   lsutil.NewUserConfig(nil),
		workspaceUserConfig: lsutil.NewUserConfig(nil), // initialize so all `config`s are non-nil
		pendingATAChanges:   make(map[tspath.Path]*ATAStateChange),
		watches:             make(map[fileSystemWatcherKey]*fileSystemWatcherValue),
	}

	if init.Options.TypingsLocation != "" && init.NpmExecutor != nil {
		session.typingsInstaller = ata.NewTypingsInstaller(&ata.TypingsInstallerOptions{
			TypingsLocation: init.Options.TypingsLocation,
			ThrottleLimit:   5,
		}, session)
	}

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

// Gets copy of current configuration
func (s *Session) Config() *lsutil.UserConfig {
	s.userConfigRWMu.Lock()
	defer s.userConfigRWMu.Unlock()
	return s.workspaceUserConfig.Copy()
}

// Trace implements module.ResolutionHost
func (s *Session) Trace(msg string) {
	panic("ATA module resolution should not use tracing")
}

func (s *Session) Configure(config *lsutil.UserConfig) {
	// `config` should never be nil
	s.userConfigRWMu.Lock()
	defer s.userConfigRWMu.Unlock()
	s.pendingUserConfigChanges = true
	oldConfig := s.workspaceUserConfig.Copy()
	s.workspaceUserConfig = s.workspaceUserConfig.Merge(config)

	// Tell the client to re-request certain commands depending on user preference changes.
	if oldConfig != config {
		s.refreshInlayHintsIfNeeded(oldConfig, config)
		s.refreshCodeLensIfNeeded(oldConfig, config)
		s.refreshDiagnosticsIfNeeded(oldConfig, config)
	}
}

func (s *Session) InitializeWithUserConfig(config *lsutil.UserConfig) {
	s.initialUserConfig = config.Copy() // initializes with non-nil config
	s.Configure(s.initialUserConfig)
}

func (s *Session) DidOpenFile(ctx context.Context, uri lsproto.DocumentUri, version int32, content string, languageKind lsproto.LanguageKind) {
	s.cancelDiagnosticsRefresh()
	s.scheduleIdleCacheClean()
	s.snapshotUpdateMu.Lock()
	defer s.snapshotUpdateMu.Unlock()
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
		reason:      UpdateReasonDidOpenFile,
		fileChanges: changes,
		ResourceRequest: ResourceRequest{
			Documents: []lsproto.DocumentUri{uri},
		},
	})
}

func (s *Session) DidCloseFile(ctx context.Context, uri lsproto.DocumentUri) {
	s.cancelDiagnosticsRefresh()
	s.scheduleIdleCacheClean()
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingFileChanges = append(s.pendingFileChanges, FileChange{
		Kind: FileChangeKindClose,
		URI:  uri,
	})
}

func (s *Session) DidChangeFile(ctx context.Context, uri lsproto.DocumentUri, version int32, changes []lsproto.TextDocumentContentChangePartialOrWholeDocument) {
	s.cancelDiagnosticsRefresh()
	s.scheduleIdleCacheClean()
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
	s.cancelDiagnosticsRefresh()
	s.scheduleIdleCacheClean()
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

	// Schedule a debounced diagnostics refresh
	s.ScheduleDiagnosticsRefresh()
	s.scheduleIdleCacheClean()
}

func (s *Session) DidChangeCompilerOptionsForInferredProjects(ctx context.Context, options *core.CompilerOptions) {
	s.compilerOptionsForInferredProjects = options
	s.UpdateSnapshot(ctx, s.fs.Overlays(), SnapshotChange{
		reason:                             UpdateReasonDidChangeCompilerOptionsForInferredProjects,
		compilerOptionsForInferredProjects: options,
	})
}

func (s *Session) ScheduleDiagnosticsRefresh() {
	s.diagnosticsRefreshMu.Lock()
	defer s.diagnosticsRefreshMu.Unlock()

	// Cancel any existing scheduled diagnostics refresh
	if s.diagnosticsRefreshCancel != nil {
		s.diagnosticsRefreshCancel()
		s.logger.Log("Delaying scheduled diagnostics refresh...")
	} else {
		s.logger.Log("Scheduling new diagnostics refresh...")
	}

	// Create a new cancellable context for the debounce task
	debounceCtx, cancel := context.WithCancel(s.backgroundCtx)
	s.diagnosticsRefreshCancel = cancel

	// Enqueue the debounced diagnostics refresh
	s.backgroundQueue.Enqueue(debounceCtx, func(ctx context.Context) {
		// Sleep for the debounce delay
		select {
		case <-time.After(s.options.DebounceDelay):
			// Delay completed, proceed with refresh
		case <-ctx.Done():
			// Context was cancelled, newer events arrived
			return
		}

		// Clear the cancel function since we're about to execute the refresh
		s.diagnosticsRefreshMu.Lock()
		s.diagnosticsRefreshCancel = nil
		s.diagnosticsRefreshMu.Unlock()

		if s.options.LoggingEnabled {
			s.logger.Log("Running scheduled diagnostics refresh")
		}
		if err := s.client.RefreshDiagnostics(s.backgroundCtx); err != nil && s.options.LoggingEnabled {
			s.logger.Logf("Error refreshing diagnostics: %v", err)
		}
	})
}

func (s *Session) cancelDiagnosticsRefresh() {
	s.diagnosticsRefreshMu.Lock()
	defer s.diagnosticsRefreshMu.Unlock()
	if s.diagnosticsRefreshCancel != nil {
		s.diagnosticsRefreshCancel()
		s.logger.Log("Canceled scheduled diagnostics refresh")
		s.diagnosticsRefreshCancel = nil
	}
}

const idleCacheCleanDelay = 30 * time.Second

func (s *Session) scheduleIdleCacheClean() {
	s.idleCacheCleanMu.Lock()
	defer s.idleCacheCleanMu.Unlock()

	if s.idleCacheCleanTimer != nil {
		s.idleCacheCleanTimer.Stop()
	}

	s.idleCacheCleanTimer = time.AfterFunc(idleCacheCleanDelay, func() {
		s.idleCacheCleanMu.Lock()
		s.idleCacheCleanTimer = nil
		s.idleCacheCleanMu.Unlock()

		s.snapshotUpdateMu.Lock()
		defer s.snapshotUpdateMu.Unlock()

		ctx := s.backgroundCtx
		fileChanges, overlays, ataChanges, newConfig := s.flushChanges(ctx)
		s.UpdateSnapshot(ctx, overlays, SnapshotChange{
			reason:         UpdateReasonIdleCleanDiskCache,
			fileChanges:    fileChanges,
			ataChanges:     ataChanges,
			newConfig:      newConfig,
			cleanDiskCache: true,
		})

		runtime.GC()
	})
}

func (s *Session) cancelIdleCacheClean() {
	s.idleCacheCleanMu.Lock()
	defer s.idleCacheCleanMu.Unlock()
	if s.idleCacheCleanTimer != nil {
		s.idleCacheCleanTimer.Stop()
		s.idleCacheCleanTimer = nil
	}
}

const performanceTelemetryInterval = 5 * time.Minute

// StartPerformanceTelemetry begins periodic collection and sending of performance
// telemetry. It should be called once after the session is initialized.
func (s *Session) StartPerformanceTelemetry() {
	if !s.options.TelemetryEnabled {
		return
	}
	ctx, cancel := context.WithCancel(s.backgroundCtx)
	s.performanceTelemetryCancel = cancel
	s.backgroundQueue.Enqueue(ctx, func(ctx context.Context) {
		ticker := time.NewTicker(performanceTelemetryInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.sendPerformanceTelemetry(ctx)
			}
		}
	})
}

func (s *Session) stopPerformanceTelemetry() {
	if s.performanceTelemetryCancel != nil {
		s.performanceTelemetryCancel()
		s.performanceTelemetryCancel = nil
	}
}

func (s *Session) sendPerformanceTelemetry(ctx context.Context) {
	if s.client == nil || !s.options.TelemetryEnabled {
		return
	}
	s.snapshotMu.RLock()
	snapshot := s.snapshot
	s.snapshotMu.RUnlock()

	// Read Go runtime metrics in a single call
	const (
		sMemoryUsedBytes = iota
		sGoMemLimit
		sGoGCPercent
		sHeapGoalBytes
		sHeapLiveBytes
		sGoroutineCount
		sGcCyclesTotal
		sGcCPUSeconds
		sUserCPUSeconds
		sMetricCount
	)
	samples := make([]gometrics.Sample, sMetricCount)
	samples[sMemoryUsedBytes].Name = "/memory/classes/total:bytes"
	samples[sGoMemLimit].Name = "/gc/gomemlimit:bytes"
	samples[sGoGCPercent].Name = "/gc/gogc:percent"
	samples[sHeapGoalBytes].Name = "/gc/heap/goal:bytes"
	samples[sHeapLiveBytes].Name = "/gc/heap/live:bytes"
	samples[sGoroutineCount].Name = "/sched/goroutines:goroutines"
	samples[sGcCyclesTotal].Name = "/gc/cycles/total:gc-cycles"
	samples[sGcCPUSeconds].Name = "/cpu/classes/gc/total:cpu-seconds"
	samples[sUserCPUSeconds].Name = "/cpu/classes/user:cpu-seconds"
	gometrics.Read(samples)

	measurements := &lsproto.PerformanceStatsTelemetryMeasurements{
		OpenFileCount:       float64(len(snapshot.fs.overlays)),
		ProjectCount:        float64(len(snapshot.ProjectCollection.Projects())),
		ConfigCount:         float64(len(snapshot.ConfigFileRegistry.configs)),
		CachedDiskFileCount: float64(len(snapshot.fs.diskFiles)),
	}

	readUint64 := func(s gometrics.Sample) float64 {
		if s.Value.Kind() == gometrics.KindUint64 {
			return float64(s.Value.Uint64())
		}
		return 0
	}
	readFloat64 := func(s gometrics.Sample) float64 {
		if s.Value.Kind() == gometrics.KindFloat64 {
			return s.Value.Float64()
		}
		return 0
	}

	measurements.MemoryUsedBytes = readUint64(samples[sMemoryUsedBytes])
	if samples[sGoMemLimit].Value.Kind() == gometrics.KindUint64 {
		v := samples[sGoMemLimit].Value.Uint64()
		if v < uint64(math.MaxInt64) {
			measurements.GoMemLimit = float64(v)
		}
	}
	measurements.GoGCPercent = readUint64(samples[sGoGCPercent])
	measurements.HeapGoalBytes = readUint64(samples[sHeapGoalBytes])
	measurements.HeapLiveBytes = readUint64(samples[sHeapLiveBytes])
	measurements.GoroutineCount = readUint64(samples[sGoroutineCount])
	measurements.GcCyclesTotal = readUint64(samples[sGcCyclesTotal])
	measurements.GcCPUSeconds = readFloat64(samples[sGcCPUSeconds])
	measurements.UserCPUSeconds = readFloat64(samples[sUserCPUSeconds])

	// Read system memory stats
	if sysMem, err := osmemory.Get(); err == nil {
		measurements.SystemMemTotal = float64(sysMem.Total)
		measurements.SystemMemUsed = float64(sysMem.Used)
	}

	if err := s.client.SendTelemetry(ctx, lsproto.TelemetryEvent{
		PerformanceStatsTelemetryEvent: &lsproto.PerformanceStatsTelemetryEvent{
			Measurements: measurements,
		},
	}); err != nil && s.options.LoggingEnabled {
		s.logger.Logf("Error sending performance telemetry: %v", err)
	}
}

func (s *Session) sendProjectInfoTelemetryForNewProjects(oldSnapshot *Snapshot, newSnapshot *Snapshot) {
	if !s.options.TelemetryEnabled {
		return
	}
	ctx := s.backgroundCtx
	collections.DiffOrderedMaps(
		oldSnapshot.ProjectCollection.ProjectsByPath(),
		newSnapshot.ProjectCollection.ProjectsByPath(),
		func(_ tspath.Path, addedProject *Project) {
			s.sendProjectInfoTelemetry(ctx, addedProject)
		},
		func(_ tspath.Path, _ *Project) {},
		func(_ tspath.Path, _, _ *Project) {},
	)
}

func (s *Session) sendProjectInfoTelemetry(ctx context.Context, project *Project) {
	if s.client == nil || !s.options.TelemetryEnabled {
		return
	}
	if s.seenProjects.Has(project.configFilePath) {
		return
	}
	s.seenProjects.Add(project.configFilePath)

	if project.Program == nil || project.CommandLine == nil {
		return
	}

	info := s.collectProjectInfoTelemetry(project)
	if err := s.client.SendTelemetry(ctx, info); err != nil && s.options.LoggingEnabled {
		s.logger.Logf("Error sending project info telemetry: %v", err)
	}
}

func (s *Session) collectProjectInfoTelemetry(project *Project) lsproto.TelemetryEvent {
	opts := project.CommandLine.CompilerOptions()
	if opts == nil {
		opts = &core.CompilerOptions{}
	}

	configFileName := "other"
	if project.Kind == KindConfigured {
		baseName := tspath.GetBaseFileName(project.configFileName)
		if baseName == "tsconfig.json" || baseName == "jsconfig.json" {
			configFileName = baseName
		}
	}

	projectType := "inferred"
	if project.Kind == KindConfigured {
		projectType = "configured"
	}

	props := map[string]string{
		"configFileName": configFileName,
		"projectType":    projectType,
		"version":        core.Version(),
	}

	// Compiler options — same approach as Strada's convertCompilerOptionsForTelemetry:
	// booleans and enum string names, no paths.
	compilerOptions := map[string]any{}
	setTristate(compilerOptions, "strict", opts.Strict)
	setTristate(compilerOptions, "allowJs", opts.AllowJs)
	setTristate(compilerOptions, "checkJs", opts.CheckJs)
	setTristate(compilerOptions, "noEmit", opts.NoEmit)
	setTristate(compilerOptions, "declaration", opts.Declaration)
	setTristate(compilerOptions, "composite", opts.Composite)
	setTristate(compilerOptions, "isolatedModules", opts.IsolatedModules)
	setTristate(compilerOptions, "esModuleInterop", opts.ESModuleInterop)
	setTristate(compilerOptions, "skipLibCheck", opts.SkipLibCheck)
	setTristate(compilerOptions, "incremental", opts.Incremental)
	if opts.Target != core.ScriptTargetNone {
		compilerOptions["target"] = opts.Target.String()
	}
	if opts.Module != core.ModuleKindNone {
		compilerOptions["module"] = opts.Module.String()
	}
	if name := moduleResolutionKindName(opts.ModuleResolution); name != "" {
		compilerOptions["moduleResolution"] = name
	}
	if opts.Jsx != core.JsxEmitNone {
		compilerOptions["jsx"] = fmt.Sprintf("%d", opts.Jsx)
	}
	if b, err := json.Marshal(compilerOptions); err == nil {
		props["compilerOptions"] = string(b)
	}

	// Config file shape
	if raw, ok := project.CommandLine.Raw.(*collections.OrderedMap[string, any]); ok {
		props["extends"] = boolTelemetry(raw.Has("extends"))
		props["files"] = boolTelemetry(raw.Has("files"))
		props["include"] = boolTelemetry(raw.Has("include"))
		props["exclude"] = boolTelemetry(raw.Has("exclude"))
	}

	return lsproto.TelemetryEvent{
		ProjectInfoTelemetryEvent: &lsproto.ProjectInfoTelemetryEvent{
			Properties:   props,
			Measurements: countFileStats(project.Program.GetSourceFiles()),
		},
	}
}

func setTristate(m map[string]any, key string, v core.Tristate) {
	if v == core.TSTrue {
		m[key] = true
	} else if v == core.TSFalse {
		m[key] = false
	}
}

func boolTelemetry(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func countFileStats(sourceFiles []*ast.SourceFile) *lsproto.ProjectInfoTelemetryMeasurements {
	var stats lsproto.ProjectInfoTelemetryMeasurements
	for _, sf := range sourceFiles {
		fileName := sf.FileName()
		size := float64(sf.End())
		switch core.GetScriptKindFromFileName(fileName) {
		case core.ScriptKindJS:
			stats.JsFileCount++
			stats.JsFileSize += size
		case core.ScriptKindJSX:
			stats.JsxFileCount++
			stats.JsxFileSize += size
		case core.ScriptKindTS:
			if tspath.IsDeclarationFileName(fileName) {
				stats.DtsFileCount++
				stats.DtsFileSize += size
			} else {
				stats.TsFileCount++
				stats.TsFileSize += size
			}
		case core.ScriptKindTSX:
			stats.TsxFileCount++
			stats.TsxFileSize += size
		}
	}
	return &stats
}

func moduleResolutionKindName(kind core.ModuleResolutionKind) string {
	switch kind {
	case core.ModuleResolutionKindUnknown:
		return ""
	case core.ModuleResolutionKindClassic:
		return "Classic"
	case core.ModuleResolutionKindNode10:
		return "Node10"
	case core.ModuleResolutionKindNode16:
		return "Node16"
	case core.ModuleResolutionKindNodeNext:
		return "NodeNext"
	case core.ModuleResolutionKindBundler:
		return "Bundler"
	default:
		return ""
	}
}

func (s *Session) Snapshot() *Snapshot {
	s.snapshotMu.RLock()
	defer s.snapshotMu.RUnlock()
	return s.snapshot
}

func (s *Session) getSnapshot(
	ctx context.Context,
	request ResourceRequest,
) *Snapshot {
	var snapshot *Snapshot
	s.snapshotUpdateMu.Lock()
	defer s.snapshotUpdateMu.Unlock()

	fileChanges, overlays, ataChanges, newConfig := s.flushChanges(ctx)
	updateSnapshot := !fileChanges.IsEmpty() || len(ataChanges) > 0 || newConfig != nil
	if updateSnapshot {
		// If there are pending file changes, we need to update the snapshot.
		// Sending the requested URI ensures that the project for this URI is loaded.
		return s.UpdateSnapshot(ctx, overlays, SnapshotChange{
			reason:          UpdateReasonRequestedLanguageServicePendingChanges,
			fileChanges:     fileChanges,
			ataChanges:      ataChanges,
			newConfig:       newConfig,
			ResourceRequest: request,
		})
	}
	// If there are no pending file changes, we can try to use the current snapshot.
	s.snapshotMu.RLock()
	snapshot = s.snapshot
	s.snapshotMu.RUnlock()

	var updateReason UpdateReason
	if len(request.Projects) > 0 {
		updateReason = UpdateReasonRequestedLanguageServiceProjectDirty
	} else if request.ProjectTree != nil {
		updateReason = UpdateReasonRequestedLoadProjectTree
	} else if request.AutoImports != "" {
		updateReason = UpdateReasonRequestedLanguageServiceWithAutoImports
	} else {
		for _, document := range request.Documents {
			project := snapshot.GetDefaultProject(document)
			if project == nil {
				updateReason = UpdateReasonRequestedLanguageServiceProjectNotLoaded
			} else if project.dirty {
				updateReason = UpdateReasonRequestedLanguageServiceProjectDirty
			}
		}
		if updateReason == UpdateReasonUnknown {
			for _, document := range request.ConfiguredProjectDocuments {
				if snapshot.fs.isOpenFile(document.FileName()) {
					project := snapshot.GetDefaultProject(document)
					if project == nil {
						updateReason = UpdateReasonRequestedLanguageServiceProjectNotLoaded
					} else if project.dirty {
						updateReason = UpdateReasonRequestedLanguageServiceProjectDirty
					}
				} else {
					updateReason = UpdateReasonRequestedLanguageServiceForFileNotOpen
				}
			}
		}
	}

	if updateReason != UpdateReasonUnknown {
		snapshot = s.UpdateSnapshot(ctx, overlays, SnapshotChange{
			reason:          updateReason,
			ResourceRequest: request,
		})
	}
	return snapshot
}

func (s *Session) getSnapshotAndDefaultProject(ctx context.Context, uri lsproto.DocumentUri) (*Snapshot, *Project, *ls.LanguageService, error) {
	snapshot := s.getSnapshot(
		ctx,
		ResourceRequest{Documents: []lsproto.DocumentUri{uri}},
	)
	project := snapshot.GetDefaultProject(uri)
	if project == nil {
		return nil, nil, nil, fmt.Errorf("no project found for URI %s", uri)
	}
	return snapshot, project, ls.NewLanguageService(project.configFilePath, project.GetProgram(), snapshot, uri.FileName()), nil
}

func (s *Session) GetLanguageService(ctx context.Context, uri lsproto.DocumentUri) (*ls.LanguageService, error) {
	_, _, languageService, err := s.getSnapshotAndDefaultProject(ctx, uri)
	if err != nil {
		return nil, err
	}
	return languageService, nil
}

// GetLanguageServiceAndSnapshot returns a LanguageService and the current snapshot.
func (s *Session) GetLanguageServiceAndSnapshot(ctx context.Context, uri lsproto.DocumentUri) (*ls.LanguageService, *Snapshot, error) {
	snapshot, _, languageService, err := s.getSnapshotAndDefaultProject(ctx, uri)
	if err != nil {
		return nil, nil, err
	}
	return languageService, snapshot, nil
}

func (s *Session) GetLanguageServiceAndProjectsForFile(ctx context.Context, uri lsproto.DocumentUri) (*Project, *ls.LanguageService, []ls.Project, error) {
	snapshot, project, defaultLs, err := s.getSnapshotAndDefaultProject(ctx, uri)
	if err != nil {
		return nil, nil, nil, err
	}
	// !!! TODO: sheetal:  Get other projects that contain the file with symlink
	allProjects := snapshot.GetProjectsContainingFile(uri)
	return project, defaultLs, allProjects, nil
}

func (s *Session) GetProjectsForFile(ctx context.Context, uri lsproto.DocumentUri) ([]ls.Project, error) {
	snapshot := s.getSnapshot(
		ctx,
		ResourceRequest{ConfiguredProjectDocuments: []lsproto.DocumentUri{uri}},
	)

	// !!! TODO: sheetal:  Get other projects that contain the file with symlink
	allProjects := snapshot.GetProjectsContainingFile(uri)
	return allProjects, nil
}

func (s *Session) GetLanguageServiceForProjectWithFile(ctx context.Context, project *Project, uri lsproto.DocumentUri) *ls.LanguageService {
	snapshot := s.getSnapshot(
		ctx,
		ResourceRequest{Projects: []tspath.Path{project.Id()}},
	)
	// Ensure we have updated project
	project = snapshot.ProjectCollection.GetProjectByPath(project.Id())
	if project == nil {
		return nil
	}
	// if program doesnt contain this file any more ignore it
	if !project.HasFile(uri.FileName()) {
		return nil
	}
	return ls.NewLanguageService(project.configFilePath, project.GetProgram(), snapshot, uri.FileName())
}

func (s *Session) GetSnapshotLoadingProjectTree(
	ctx context.Context,
	// If null, all project trees need to be loaded, otherwise only those that are referenced
	requestedProjectTrees *collections.Set[tspath.Path],
) *Snapshot {
	snapshot := s.getSnapshot(
		ctx,
		ResourceRequest{ProjectTree: &ProjectTreeRequest{requestedProjectTrees}},
	)
	return snapshot
}

// GetCurrentLanguageServiceWithAutoImports flushes pending file changes, clones the
// current snapshot with auto-import preparation for the given URI, then returns a
// LanguageService for the default project. Use this only outside of request handling
// (e.g. cache warming). For request handlers, use GetLanguageServiceWithAutoImports
// with the request-level snapshot instead.
func (s *Session) GetCurrentLanguageServiceWithAutoImports(ctx context.Context, uri lsproto.DocumentUri) (*ls.LanguageService, error) {
	snapshot := s.getSnapshot(ctx, ResourceRequest{
		Documents:   []lsproto.DocumentUri{uri},
		AutoImports: uri,
	})
	project := snapshot.GetDefaultProject(uri)
	if project == nil {
		return nil, fmt.Errorf("no project found for URI %s", uri)
	}
	return ls.NewLanguageService(project.configFilePath, project.GetProgram(), snapshot, uri.FileName()), nil
}

// GetLanguageServiceWithAutoImports clones the given snapshot with auto-import
// preparation for the given URI, without flushing pending file changes. This
// keeps the request bound to its original snapshot.
func (s *Session) GetLanguageServiceWithAutoImports(ctx context.Context, baseSnapshot *Snapshot, uri lsproto.DocumentUri) (*ls.LanguageService, error) {
	change := SnapshotChange{
		reason: UpdateReasonRequestedLanguageServiceWithAutoImports,
		ResourceRequest: ResourceRequest{
			Documents:   []lsproto.DocumentUri{uri},
			AutoImports: uri,
		},
	}
	newSnapshot := baseSnapshot.Clone(ctx, change, baseSnapshot.fs.overlays, s)

	project := newSnapshot.GetDefaultProject(uri)
	if project == nil {
		newSnapshot.Dispose(s)
		return nil, fmt.Errorf("no project found for URI %s", uri)
	}

	s.backgroundQueue.Enqueue(s.backgroundCtx, func(ctx context.Context) {
		s.adoptSnapshotChange(baseSnapshot, newSnapshot)
	})

	return ls.NewLanguageService(project.configFilePath, project.GetProgram(), newSnapshot, uri.FileName()), nil
}

// adoptSnapshotChange promotes a cloned snapshot as the session's current
// snapshot so future requests benefit from the work already done. If the
// session has moved on, the snapshot is discarded; the next request needing
// auto-imports will redo the work on the latest snapshot.
func (s *Session) adoptSnapshotChange(baseSnapshot, newSnapshot *Snapshot) {
	s.snapshotMu.Lock()
	oldSnapshot := s.snapshot
	if oldSnapshot == baseSnapshot {
		s.snapshot = newSnapshot
		s.snapshotMu.Unlock()
		oldSnapshot.Dispose(s)
	} else {
		s.snapshotMu.Unlock()
		newSnapshot.Dispose(s)
	}
}

func (s *Session) UpdateSnapshot(ctx context.Context, overlays map[tspath.Path]*Overlay, change SnapshotChange) *Snapshot {
	s.snapshotMu.Lock()
	oldSnapshot := s.snapshot
	newSnapshot := oldSnapshot.Clone(ctx, change, overlays, s)
	s.snapshot = newSnapshot
	s.snapshotMu.Unlock()

	if newSnapshot != oldSnapshot {
		oldSnapshot.Dispose(s)
	}

	// Enqueue ATA updates if needed
	if s.typingsInstaller != nil {
		s.triggerATAForUpdatedProjects(newSnapshot)
	}

	// Enqueue logging, watch updates, and diagnostic refresh tasks
	// !!! userPreferences/configuration updates
	s.backgroundQueue.Enqueue(s.backgroundCtx, func(ctx context.Context) {
		if s.options.LoggingEnabled {
			s.logger.Log(newSnapshot.builderLogs.String())
			s.logProjectChanges(oldSnapshot, newSnapshot)
			s.logRuntimeMetrics()
			s.logger.Log("")
		}
		if s.options.WatchEnabled {
			if err := s.updateWatches(oldSnapshot, newSnapshot); err != nil && s.options.LoggingEnabled {
				s.logger.Log(err)
			}
		}
		s.publishProgramDiagnostics(oldSnapshot, newSnapshot)
		s.sendProjectInfoTelemetryForNewProjects(oldSnapshot, newSnapshot)
		s.warmAutoImportCache(ctx, change, oldSnapshot, newSnapshot)
	})

	return newSnapshot
}

// WaitForBackgroundTasks waits for all background tasks to complete.
// This is intended to be used only for testing purposes.
func (s *Session) WaitForBackgroundTasks() {
	s.cancelIdleCacheClean()
	s.backgroundQueue.Wait()
}

func updateWatch[T any](ctx context.Context, session *Session, logger logging.Logger, oldWatcher, newWatcher *WatchedFiles[T]) []error {
	var errors []error
	session.watchesMu.Lock()
	defer session.watchesMu.Unlock()
	if newWatcher != nil {
		if id, watchers, ignored := newWatcher.Watchers(); len(watchers) > 0 {
			var newWatchers collections.OrderedMap[WatcherID, *lsproto.FileSystemWatcher]
			for i, watcher := range watchers {
				key := toFileSystemWatcherKey(watcher)
				value := session.watches[key]
				globId := WatcherID(fmt.Sprintf("%s.%d", id, i))
				if value == nil {
					value = &fileSystemWatcherValue{id: globId}
					session.watches[key] = value
				}
				value.count++
				if value.count == 1 {
					newWatchers.Set(globId, watcher)
				}
			}
			for id, watcher := range newWatchers.Entries() {
				if err := session.client.WatchFiles(ctx, id, []*lsproto.FileSystemWatcher{watcher}); err != nil {
					errors = append(errors, err)
				} else if logger != nil {
					if oldWatcher == nil {
						logger.Log(fmt.Sprintf("Added new watch: %s", id))
					} else {
						logger.Log(fmt.Sprintf("Updated watch: %s", id))
					}
					logger.Log("\t" + *watcher.GlobPattern.Pattern)
					logger.Log("")
				}
			}
			if len(ignored) > 0 {
				logger.Logf("%d paths ineligible for watching", len(ignored))
				if logger.IsVerbose() {
					for path := range ignored {
						logger.Log("\t" + path)
					}
				}
			}
		}
	}
	if oldWatcher != nil {
		if _, watchers, _ := oldWatcher.Watchers(); len(watchers) > 0 {
			var removedWatchers []WatcherID
			for _, watcher := range watchers {
				key := toFileSystemWatcherKey(watcher)
				value := session.watches[key]
				if value == nil {
					continue
				}
				if value.count <= 1 {
					delete(session.watches, key)
					removedWatchers = append(removedWatchers, value.id)
				} else {
					value.count--
				}
			}
			for _, id := range removedWatchers {
				if err := session.client.UnwatchFiles(ctx, id); err != nil {
					errors = append(errors, err)
				} else if logger != nil && newWatcher == nil {
					logger.Log(fmt.Sprintf("Removed watch: %s", id))
				}
			}
		}
	}
	return errors
}

func (s *Session) updateWatches(oldSnapshot *Snapshot, newSnapshot *Snapshot) error {
	var errors []error
	start := time.Now()
	ctx := s.backgroundCtx
	core.DiffMapsFunc(
		oldSnapshot.ConfigFileRegistry.configs,
		newSnapshot.ConfigFileRegistry.configs,
		func(a, b *configFileEntry) bool {
			return a.rootFilesWatch.ID() == b.rootFilesWatch.ID()
		},
		func(_ tspath.Path, addedEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s, s.logger, nil, addedEntry.rootFilesWatch)...)
		},
		func(_ tspath.Path, removedEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s, s.logger, removedEntry.rootFilesWatch, nil)...)
		},
		func(_ tspath.Path, oldEntry, newEntry *configFileEntry) {
			errors = append(errors, updateWatch(ctx, s, s.logger, oldEntry.rootFilesWatch, newEntry.rootFilesWatch)...)
		},
	)

	collections.DiffOrderedMaps(
		oldSnapshot.ProjectCollection.ProjectsByPath(),
		newSnapshot.ProjectCollection.ProjectsByPath(),
		func(_ tspath.Path, addedProject *Project) {
			errors = append(errors, updateWatch(ctx, s, s.logger, nil, addedProject.programFilesWatch)...)
			errors = append(errors, updateWatch(ctx, s, s.logger, nil, addedProject.typingsWatch)...)
		},
		func(_ tspath.Path, removedProject *Project) {
			errors = append(errors, updateWatch(ctx, s, s.logger, removedProject.programFilesWatch, nil)...)
			errors = append(errors, updateWatch(ctx, s, s.logger, removedProject.typingsWatch, nil)...)
		},
		func(_ tspath.Path, oldProject, newProject *Project) {
			if oldProject.programFilesWatch.ID() != newProject.programFilesWatch.ID() {
				errors = append(errors, updateWatch(ctx, s, s.logger, oldProject.programFilesWatch, newProject.programFilesWatch)...)
			}
			if oldProject.typingsWatch.ID() != newProject.typingsWatch.ID() {
				errors = append(errors, updateWatch(ctx, s, s.logger, oldProject.typingsWatch, newProject.typingsWatch)...)
			}
		},
	)

	if oldSnapshot.autoImportsWatch.ID() != newSnapshot.autoImportsWatch.ID() {
		errors = append(errors, updateWatch(ctx, s, s.logger, oldSnapshot.autoImportsWatch, newSnapshot.autoImportsWatch)...)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors updating watches: %v", errors)
	} else if s.options.LoggingEnabled {
		s.logger.Log(fmt.Sprintf("Updated watches in %v", time.Since(start)))
	}
	return nil
}

func (s *Session) Close() {
	// Cancel any pending diagnostics refresh
	s.cancelDiagnosticsRefresh()
	// Cancel any pending idle cache clean
	s.cancelIdleCacheClean()
	// Cancel periodic performance telemetry
	s.stopPerformanceTelemetry()
	s.backgroundQueue.Close()
}

func (s *Session) flushChanges(ctx context.Context) (FileChangeSummary, map[tspath.Path]*Overlay, map[tspath.Path]*ATAStateChange, *lsutil.UserConfig) {
	s.pendingFileChangesMu.Lock()
	defer s.pendingFileChangesMu.Unlock()
	s.pendingATAChangesMu.Lock()
	defer s.pendingATAChangesMu.Unlock()
	pendingATAChanges := s.pendingATAChanges
	s.pendingATAChanges = make(map[tspath.Path]*ATAStateChange)
	fileChanges, overlays := s.flushChangesLocked(ctx)
	s.userConfigRWMu.Lock()
	defer s.userConfigRWMu.Unlock()
	var newUserConfig *lsutil.UserConfig
	if s.pendingUserConfigChanges {
		newUserConfig = s.workspaceUserConfig.Copy()
	}
	s.pendingUserConfigChanges = false
	return fileChanges, overlays, pendingATAChanges, newUserConfig
}

// flushChangesLocked should only be called with s.pendingFileChangesMu held.
func (s *Session) flushChangesLocked(ctx context.Context) (FileChangeSummary, map[tspath.Path]*Overlay) {
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
	var loggedProjectChanges bool
	logProject := func(project *Project) {
		var builder strings.Builder
		project.print(s.logger.IsVerbose() /*writeFileNames*/, s.logger.IsVerbose() /*writeFileExplanation*/, &builder)
		s.logger.Log(builder.String())
		loggedProjectChanges = true
	}
	collections.DiffOrderedMaps(
		oldSnapshot.ProjectCollection.ProjectsByPath(),
		newSnapshot.ProjectCollection.ProjectsByPath(),
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

	if loggedProjectChanges || s.logger.IsVerbose() {
		s.logCacheStats(newSnapshot)
	}
}

var runtimeMetricsSamples = sync.OnceValue(func() []gometrics.Sample {
	descs := gometrics.All()
	var samples []gometrics.Sample
	for _, desc := range descs {
		name := desc.Name
		if strings.HasPrefix(name, "/memory/") || strings.HasPrefix(name, "/gc/") {
			samples = append(samples, gometrics.Sample{Name: name})
		}
	}
	return samples
})

func (s *Session) logRuntimeMetrics() {
	samples := slices.Clone(runtimeMetricsSamples())
	gometrics.Read(samples)

	var builder strings.Builder
	builder.WriteString("\n======== Runtime Metrics ========")
	for _, sample := range samples {
		switch sample.Value.Kind() {
		case gometrics.KindUint64:
			fmt.Fprintf(&builder, "\n%s = %d", sample.Name, sample.Value.Uint64())
		case gometrics.KindFloat64:
			fmt.Fprintf(&builder, "\n%s = %f", sample.Name, sample.Value.Float64())
		case gometrics.KindFloat64Histogram:
			// Skip histograms for log readability
		}
	}
	s.logger.Log(builder.String())
}

func (s *Session) logCacheStats(snapshot *Snapshot) {
	var parseCacheSize int
	var extendedConfigCount int
	if s.logger.IsVerbose() {
		s.parseCache.entries.Range(func(_ ParseCacheKey, _ *refCountCacheEntry[*ast.SourceFile]) bool {
			parseCacheSize++
			return true
		})
		s.extendedConfigCache.entries.Range(func(_ tspath.Path, _ *refCountCacheEntry[*ExtendedConfigCacheEntry]) bool {
			extendedConfigCount++
			return true
		})
	}
	s.logger.Log("\n======== Cache Statistics ========")
	s.logger.Logf("Open file count:   %6d", len(snapshot.fs.overlays))
	s.logger.Logf("Cached disk files: %6d", len(snapshot.fs.diskFiles))
	s.logger.Logf("Project count:     %6d", len(snapshot.ProjectCollection.Projects()))
	s.logger.Logf("Config count:      %6d", len(snapshot.ConfigFileRegistry.configs))
	if s.logger.IsVerbose() {
		s.logger.Logf("Parse cache size:           %6d", parseCacheSize)
		s.logger.Logf("Program count:              %6d", s.programCounter.Len())
		s.logger.Logf("Extended config cache size: %6d", extendedConfigCount)

		s.logger.Log("Auto Imports:")
		autoImportStats := snapshot.AutoImportRegistry().GetCacheStats()
		s.logger.Logf("\tUnique packages (by realpath): %d", autoImportStats.UniquePackageCount)
		if len(autoImportStats.ProjectBuckets) > 0 {
			s.logger.Log("\tProject buckets:")
			for _, bucket := range autoImportStats.ProjectBuckets {
				s.logger.Logf("\t\t%s%s:", bucket.Path, core.IfElse(bucket.State.Dirty(), " (dirty)", ""))
				s.logger.Logf("\t\t\tFiles: %d", bucket.FileCount)
				s.logger.Logf("\t\t\tExports: %d", bucket.ExportCount)
			}
		}
		if len(autoImportStats.NodeModulesBuckets) > 0 {
			s.logger.Log("\tnode_modules buckets:")
			for _, bucket := range autoImportStats.NodeModulesBuckets {
				s.logger.Logf("\t\t%s%s:", bucket.Path, core.IfElse(bucket.State.Dirty(), " (dirty)", ""))
				for packageName := range bucket.State.DirtyPackages().Keys() {
					s.logger.Logf("\t\t\tNeeds granular update: %s", packageName)
				}
				if bucket.DependencyNames != nil {
					s.logger.Logf("\t\t\tCollected packages: %d", bucket.DependencyNames.Len())
				} else {
					s.logger.Logf("\t\t\tCollected packages: all, due to no package.json!")
				}
				s.logger.Logf("\t\t\tTotal packages: %d", bucket.PackageNames.Len())
				s.logger.Logf("\t\t\tFiles: %d", bucket.FileCount)
				s.logger.Logf("\t\t\tExports: %d", bucket.ExportCount)
			}
		}
	}
}

func (s *Session) NpmInstall(cwd string, npmInstallArgs []string) ([]byte, error) {
	return s.npmExecutor.NpmInstall(cwd, npmInstallArgs)
}

func (s *Session) refreshInlayHintsIfNeeded(oldPrefs *lsutil.UserConfig, newPrefs *lsutil.UserConfig) {
	if oldPrefs.JS().InlayHints != newPrefs.JS().InlayHints || oldPrefs.TS().InlayHints != newPrefs.TS().InlayHints {
		if err := s.client.RefreshInlayHints(s.backgroundCtx); err != nil && s.options.LoggingEnabled {
			s.logger.Logf("Error refreshing inlay hints: %v", err)
		}
	}
}

func (s *Session) refreshCodeLensIfNeeded(oldPrefs *lsutil.UserConfig, newPrefs *lsutil.UserConfig) {
	if oldPrefs.JS().CodeLens != newPrefs.JS().CodeLens || oldPrefs.TS().CodeLens != newPrefs.TS().CodeLens {
		if err := s.client.RefreshCodeLens(s.backgroundCtx); err != nil && s.options.LoggingEnabled {
			s.logger.Logf("Error refreshing code lens: %v", err)
		}
	}
}

func (s *Session) refreshDiagnosticsIfNeeded(oldPrefs *lsutil.UserConfig, newPrefs *lsutil.UserConfig) {
	if oldPrefs.TS().CustomConfigFileName != newPrefs.TS().CustomConfigFileName {
		s.ScheduleDiagnosticsRefresh()
	}
}

func (s *Session) publishProgramDiagnostics(oldSnapshot *Snapshot, newSnapshot *Snapshot) {
	if !s.options.PushDiagnosticsEnabled {
		return
	}

	ctx := s.backgroundCtx
	collections.DiffOrderedMaps(
		oldSnapshot.ProjectCollection.ProjectsByPath(),
		newSnapshot.ProjectCollection.ProjectsByPath(),
		func(configFilePath tspath.Path, addedProject *Project) {
			if !shouldPublishProgramDiagnostics(addedProject, newSnapshot.ID()) {
				return
			}
			s.publishProjectDiagnostics(ctx, string(configFilePath), addedProject.Program.GetProgramDiagnostics(), newSnapshot.converters)
		},
		func(configFilePath tspath.Path, removedProject *Project) {
			if removedProject.Kind != KindConfigured {
				return
			}
			s.publishProjectDiagnostics(ctx, string(configFilePath), nil, oldSnapshot.converters)
		},
		func(configFilePath tspath.Path, oldProject, newProject *Project) {
			if !shouldPublishProgramDiagnostics(newProject, newSnapshot.ID()) {
				return
			}
			s.publishProjectDiagnostics(ctx, string(configFilePath), newProject.Program.GetProgramDiagnostics(), newSnapshot.converters)
		},
	)
}

func shouldPublishProgramDiagnostics(p *Project, snapshotID uint64) bool {
	if p.Kind != KindConfigured || p.Program == nil || p.ProgramLastUpdate != snapshotID {
		return false
	}
	return p.ProgramUpdateKind > ProgramUpdateKindCloned
}

func (s *Session) publishProjectDiagnostics(ctx context.Context, configFilePath string, diagnostics []*ast.Diagnostic, converters *lsconv.Converters) {
	lspDiagnostics := make([]*lsproto.Diagnostic, 0, len(diagnostics))
	for _, diag := range diagnostics {
		lspDiagnostics = append(lspDiagnostics, lsconv.DiagnosticToLSPPush(ctx, converters, diag))
	}

	if err := s.client.PublishDiagnostics(ctx, &lsproto.PublishDiagnosticsParams{
		Uri:         lsconv.FileNameToDocumentURI(configFilePath),
		Diagnostics: lspDiagnostics,
	}); err != nil && s.options.LoggingEnabled {
		s.logger.Logf("Error publishing diagnostics: %v", err)
	}
}

func (s *Session) triggerATAForUpdatedProjects(newSnapshot *Snapshot) {
	for _, project := range newSnapshot.ProjectCollection.Projects() {
		if project.ShouldTriggerATA(newSnapshot.ID()) {
			s.backgroundQueue.Enqueue(s.backgroundCtx, func(ctx context.Context) {
				var logTree *logging.LogTree
				if s.options.LoggingEnabled {
					logTree = logging.NewLogTree("Triggering ATA for project " + project.Name())
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

				projectDisplayName := project.DisplayName(s.options.CurrentDirectory)
				if s.client != nil {
					s.client.ProgressStart(diagnostics.Installing_types_for_0, projectDisplayName)
				}
				result, err := s.typingsInstaller.InstallTypings(request)
				if s.client != nil {
					s.client.ProgressFinish(diagnostics.Installing_types_for_0, projectDisplayName)
				}
				if err != nil {
					if logTree != nil {
						s.logger.Log(fmt.Sprintf("ATA installation failed for project %s: %v", project.Name(), err))
						s.logger.Log(logTree.String())
					}
				} else {
					if !slices.Equal(result.TypingsFiles, project.typingsFiles) {
						s.pendingATAChangesMu.Lock()
						defer s.pendingATAChangesMu.Unlock()
						s.pendingATAChanges[project.configFilePath] = &ATAStateChange{
							TypingsInfo:         &typingsInfo,
							TypingsFiles:        result.TypingsFiles,
							TypingsFilesToWatch: result.FilesToWatch,
							Logs:                logTree,
						}
						s.ScheduleDiagnosticsRefresh()
					}
				}
			})
		}
	}
}

func (s *Session) warmAutoImportCache(ctx context.Context, change SnapshotChange, oldSnapshot, newSnapshot *Snapshot) {
	if change.fileChanges.Changed.Len() == 1 {
		var changedFile lsproto.DocumentUri
		for uri := range change.fileChanges.Changed.Keys() {
			changedFile = uri
		}
		if !newSnapshot.fs.isOpenFile(changedFile.FileName()) {
			return
		}
		if newSnapshot.GetPreferences(changedFile.FileName()).IncludeCompletionsForModuleExports.IsFalse() {
			return
		}
		project := newSnapshot.GetDefaultProject(changedFile)
		if project == nil {
			return
		}
		if newSnapshot.AutoImports.IsPreparedForImportingFile(
			changedFile.FileName(),
			project.configFilePath,
			newSnapshot.allUserPreferences.GetPreferences(changedFile.FileName()),
		) {
			return
		}
		_, _ = s.GetCurrentLanguageServiceWithAutoImports(ctx, changedFile)
	}
}
