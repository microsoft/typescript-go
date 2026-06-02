package execute

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/trackingvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

type cachedSourceFile struct {
	file    *ast.SourceFile
	modTime time.Time
}

type watchCompilerHost struct {
	compiler.CompilerHost
	cache *collections.SyncMap[tspath.Path, *cachedSourceFile]
}

func (h *watchCompilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	info := h.CompilerHost.FS().Stat(opts.FileName)

	if cached, ok := h.cache.Load(opts.Path); ok {
		if info != nil && info.ModTime().Equal(cached.modTime) {
			return cached.file
		}
	}

	file := h.CompilerHost.GetSourceFile(opts)
	if file != nil {
		if info != nil {
			h.cache.Store(opts.Path, &cachedSourceFile{
				file:    file,
				modTime: info.ModTime(),
			})
		}
	} else {
		h.cache.Delete(opts.Path)
	}
	return file
}

type Watcher struct {
	mu                             sync.Mutex
	sys                            tsc.System
	configFileName                 string
	config                         *tsoptions.ParsedCommandLine
	compilerOptionsFromCommandLine *core.CompilerOptions
	reportDiagnostic               tsc.DiagnosticReporter
	reportErrorSummary             tsc.DiagnosticsReporter
	reportWatchStatus              tsc.DiagnosticReporter
	testing                        tsc.CommandLineTesting

	program             *incremental.Program
	extendedConfigCache *tsc.ExtendedConfigCache
	configModified      bool
	configHasErrors     bool
	configFilePaths     []string

	sourceFileCache *collections.SyncMap[tspath.Path, *cachedSourceFile]
	fileWatcher     *vfswatch.FileWatcher
	watches         []fswatch.Watch // targeted directory/file watches
	watchTerminated chan struct{}   // closed when a watch terminates
	doCycleCh       chan struct{}   // buffered signal to run DoCycle off the callback goroutine
	debugLog        io.Writer       // nil = silent; set via TS_WATCH_DEBUG

	changedMu       sync.Mutex
	changedPaths    map[string]fswatch.EventKind // event path → last event kind
	changedOverflow bool                         // true on ErrOverflow; forces full scan fallback
}

var _ tsc.Watcher = (*Watcher)(nil)

func createWatcher(
	sys tsc.System,
	configParseResult *tsoptions.ParsedCommandLine,
	compilerOptionsFromCommandLine *core.CompilerOptions,
	reportDiagnostic tsc.DiagnosticReporter,
	reportErrorSummary tsc.DiagnosticsReporter,
	testing tsc.CommandLineTesting,
) *Watcher {
	w := &Watcher{
		sys:                            sys,
		config:                         configParseResult,
		compilerOptionsFromCommandLine: compilerOptionsFromCommandLine,
		reportDiagnostic:               reportDiagnostic,
		reportErrorSummary:             reportErrorSummary,
		reportWatchStatus:              tsc.CreateWatchStatusReporter(sys, configParseResult.Locale(), configParseResult.CompilerOptions(), testing),
		testing:                        testing,
		sourceFileCache:                &collections.SyncMap[tspath.Path, *cachedSourceFile]{},
		watchTerminated:                make(chan struct{}),
		doCycleCh:                      make(chan struct{}, 1),
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	w.fileWatcher = vfswatch.NewFileWatcher(sys.FS())
	return w
}

func (w *Watcher) start() {
	w.mu.Lock()
	w.extendedConfigCache = &tsc.ExtendedConfigCache{}
	host := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	w.program = incremental.ReadBuildInfoProgram(w.config, incremental.NewBuildInfoReader(host), host)

	if w.configFileName != "" {
		w.configFilePaths = append([]string{w.configFileName}, w.config.ExtendedSourceFiles()...)
	}

	if w.sys.GetEnvironmentVariable("TS_WATCH_DEBUG") != "" {
		w.debugLog = w.sys.Writer()
	}

	w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Starting_compilation_in_watch_mode))
	w.doBuild()
	w.mu.Unlock()

	if w.testing == nil {
		if err := w.subscribe(); err != nil {
			fmt.Fprintf(w.sys.Writer(), "Error: Failed to start file watcher: %v\n", err)
			return
		}
		// Process DoCycle signals on a dedicated goroutine so fswatch
		// callbacks return immediately and don't block event delivery.
		go func() {
			for range w.doCycleCh {
				w.DoCycle()
			}
		}()
		// Block until a watch terminates (e.g. watched directory deleted),
		// then clean up.
		<-w.watchTerminated
		w.closeWatches()
		close(w.doCycleCh)
	}
}

func (w *Watcher) subscribe() error {
	watcher := fswatch.Default()
	if w.debugLog != nil {
		fmt.Fprintf(w.debugLog, "[watch] using %s backend\n", watcher.Name())
	}

	// Watch wildcard directories from tsconfig (recursive or non-recursive
	// based on the include patterns). This detects new/deleted source files.
	if w.config.ConfigFile != nil {
		for dir, recursive := range w.config.WildcardDirectories() {
			realDir := w.sys.FS().Realpath(dir)
			opts := []fswatch.WatchOption{fswatch.WithIgnore(shouldIgnoreWatchPath)}
			if recursive {
				opts = append(opts, fswatch.WithRecursive())
			}
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] watching directory %s (recursive=%v)\n", realDir, recursive)
			}
			watch, err := watcher.WatchDirectory(realDir, w.onWatchEvents, opts...)
			if err != nil {
				w.closeWatches()
				return fmt.Errorf("watching %s: %w", realDir, err)
			}
			w.watches = append(w.watches, watch)
		}
	}

	// Watch config files (tsconfig.json and extended configs) for changes.
	for _, path := range w.configFilePaths {
		realPath := w.sys.FS().Realpath(path)
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] watching file %s\n", realPath)
		}
		watch, err := watcher.WatchFile(realPath, w.onWatchEvents)
		if err != nil {
			w.closeWatches()
			return fmt.Errorf("watching %s: %w", realPath, err)
		}
		w.watches = append(w.watches, watch)
	}

	if len(w.watches) == 0 {
		// No config file — watch the current directory non-recursively.
		dir := w.sys.FS().Realpath(w.sys.GetCurrentDirectory())
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] no tsconfig, watching %s\n", dir)
		}
		watch, err := watcher.WatchDirectory(
			dir, w.onWatchEvents,
			fswatch.WithIgnore(shouldIgnoreWatchPath),
		)
		if err != nil {
			return err
		}
		w.watches = append(w.watches, watch)
	}

	return nil
}

func (w *Watcher) closeWatches() {
	for _, watch := range w.watches {
		watch.Close()
	}
	w.watches = nil
}

func shouldIgnoreWatchPath(path string) bool {
	p := tspath.NormalizeSlashes(path)
	return strings.HasSuffix(p, "/.git") ||
		strings.Contains(p, "/.git/") ||
		strings.Contains(p, "/node_modules/.") ||
		strings.Contains(p, "/.#")
}

func (w *Watcher) hasRelevantChanges(changedPaths map[string]fswatch.EventKind) bool {
	for eventPath := range changedPaths {
		p := tspath.NormalizeSlashes(eventPath)
		if _, ok := w.fileWatcher.WatchStateEntry(p); ok {
			return true
		}
	}
	return false
}

func (w *Watcher) evictChangedSourceFiles(changedPaths map[string]fswatch.EventKind) {
	caseSensitive := w.sys.FS().UseCaseSensitiveFileNames()
	cwd := w.sys.GetCurrentDirectory()
	for eventPath := range changedPaths {
		p := tspath.ToPath(eventPath, cwd, caseSensitive)
		if _, ok := w.sourceFileCache.Load(p); ok {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] evicting cached source file: %s\n", p)
			}
			w.sourceFileCache.Delete(p)
		}
	}
}

func (w *Watcher) onWatchEvents(events []fswatch.Event, err error) {
	if err != nil {
		if errors.Is(err, fswatch.ErrOverflow) {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] event overflow, triggering rebuild\n")
			}
			w.changedMu.Lock()
			w.changedOverflow = true
			w.changedMu.Unlock()
			w.signalDoCycle()
			return
		}
		if errors.Is(err, fswatch.ErrWatchTerminated) {
			fmt.Fprintf(w.sys.Writer(), "Warning: File watcher terminated: %v\n", err)
			close(w.watchTerminated)
			return
		}
		fmt.Fprintf(w.sys.Writer(), "Warning: File watch error: %v\n", err)
		return
	}

	if len(events) > 0 {
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] %d event(s): ", len(events))
			for i, e := range events {
				if i > 0 {
					fmt.Fprint(w.debugLog, ", ")
				}
				if i >= 5 {
					fmt.Fprintf(w.debugLog, "... and %d more", len(events)-i)
					break
				}
				fmt.Fprintf(w.debugLog, "%s %s", e.Kind, e.Path)
			}
			fmt.Fprintln(w.debugLog)
		}
		w.changedMu.Lock()
		if w.changedPaths == nil {
			w.changedPaths = make(map[string]fswatch.EventKind, len(events))
		}
		for _, e := range events {
			w.changedPaths[e.Path] = e.Kind
		}
		w.changedMu.Unlock()
		w.signalDoCycle()
	}
}

// signalDoCycle sends a non-blocking signal to the DoCycle goroutine.
// If a signal is already pending, additional signals are coalesced.
func (w *Watcher) signalDoCycle() {
	select {
	case w.doCycleCh <- struct{}{}:
		// Signal sent; the DoCycle goroutine will pick it up.
	default:
		// A signal is already pending; this event will be covered
		// by the next DoCycle's filesystem scan.
	}
}

func (w *Watcher) DoCycle() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.changedMu.Lock()
	changedPaths := w.changedPaths
	overflow := w.changedOverflow
	w.changedPaths = nil
	w.changedOverflow = false
	w.changedMu.Unlock()

	hasEvents := len(changedPaths) > 0 || overflow

	if w.recheckTsConfig() {
		return
	}

	if hasEvents && !overflow && !w.configModified {
		if w.hasRelevantChanges(changedPaths) {
			w.evictChangedSourceFiles(changedPaths)
		} else if !w.fileWatcher.WatchStateUninitialized() && !w.fileWatcher.HasChangesFromWatchState() {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] DoCycle: %d event(s) not relevant to compilation, skipping rebuild\n", len(changedPaths))
			}
			if w.testing != nil {
				w.testing.OnProgram(w.program)
			}
			return
		}
	} else if !hasEvents {
		if !w.fileWatcher.WatchStateUninitialized() && !w.configModified && !w.fileWatcher.HasChangesFromWatchState() {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] DoCycle: no tracked files changed, skipping rebuild\n")
			}
			if w.testing != nil {
				w.testing.OnProgram(w.program)
			}
			return
		}
	}

	w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.File_change_detected_Starting_incremental_compilation))
	w.doBuild()
}

func (w *Watcher) doBuild() {
	if w.configModified {
		w.sourceFileCache = &collections.SyncMap[tspath.Path, *cachedSourceFile]{}
	}

	cached := cachedvfs.From(w.sys.FS())
	tfs := &trackingvfs.FS{Inner: cached}
	innerHost := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), tfs, w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	host := &watchCompilerHost{CompilerHost: innerHost, cache: w.sourceFileCache}

	var wildcardDirs map[string]bool
	if w.config.ConfigFile != nil {
		wildcardDirs = w.config.WildcardDirectories()
		for dir := range wildcardDirs {
			tfs.SeenFiles.Add(dir)
		}
		if len(wildcardDirs) > 0 {
			w.config = w.config.ReloadFileNamesOfParsedCommandLine(w.sys.FS())
		}
	}
	for _, path := range w.configFilePaths {
		tfs.SeenFiles.Add(path)
	}

	w.program = incremental.NewProgram(compiler.NewProgram(compiler.ProgramOptions{
		Config: w.config,
		Host:   host,
	}), w.program, nil, w.testing != nil)

	result := w.compileAndEmit()
	cached.DisableAndClearCache()
	w.fileWatcher.UpdateWatchState(tfs.SeenFiles.ToSlice(), wildcardDirs)
	w.configModified = false

	programFiles := w.program.GetProgram().FilesByPath()
	w.sourceFileCache.Range(func(path tspath.Path, _ *cachedSourceFile) bool {
		if _, ok := programFiles[path]; !ok {
			w.sourceFileCache.Delete(path)
		}
		return true
	})

	errorCount := len(result.Diagnostics)
	if errorCount == 1 {
		w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Found_1_error_Watching_for_file_changes))
	} else {
		w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Found_0_errors_Watching_for_file_changes, errorCount))
	}

	if w.testing != nil {
		w.testing.OnProgram(w.program)
	}
}

func (w *Watcher) compileAndEmit() tsc.CompileAndEmitResult {
	return tsc.EmitFilesAndReportErrors(tsc.EmitInput{
		Sys:                w.sys,
		ProgramLike:        w.program,
		Program:            w.program.GetProgram(),
		Config:             w.config,
		ReportDiagnostic:   w.reportDiagnostic,
		ReportErrorSummary: w.reportErrorSummary,
		Writer:             w.sys.Writer(),
		CompileTimes:       &tsc.CompileTimes{},
		Testing:            w.testing,
	})
}

func (w *Watcher) recheckTsConfig() bool {
	if w.configFileName == "" {
		return false
	}

	if !w.configHasErrors && len(w.configFilePaths) > 0 {
		changed := false
		for _, path := range w.configFilePaths {
			old, ok := w.fileWatcher.WatchStateEntry(path)
			if !ok {
				changed = true
				break
			}
			s := w.sys.FS().Stat(path)
			if !old.Exists {
				if s != nil {
					changed = true
					break
				}
			} else {
				if s == nil || !s.ModTime().Equal(old.ModTime) {
					changed = true
					break
				}
			}
		}
		if !changed {
			return false
		}
	}

	extendedConfigCache := &tsc.ExtendedConfigCache{}
	configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(w.configFileName, w.compilerOptionsFromCommandLine, nil, w.sys, extendedConfigCache)
	if len(errors) > 0 {
		for _, e := range errors {
			w.reportDiagnostic(e)
		}
		w.configHasErrors = true
		errorCount := len(errors)
		if errorCount == 1 {
			w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Found_1_error_Watching_for_file_changes))
		} else {
			w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Found_0_errors_Watching_for_file_changes, errorCount))
		}
		return true
	}
	if w.configHasErrors {
		w.configModified = true
	}
	w.configHasErrors = false
	w.configFilePaths = append([]string{w.configFileName}, configParseResult.ExtendedSourceFiles()...)
	if !reflect.DeepEqual(w.config.ParsedConfig, configParseResult.ParsedConfig) {
		w.configModified = true
	}
	w.config = configParseResult
	w.extendedConfigCache = extendedConfigCache
	return false
}
