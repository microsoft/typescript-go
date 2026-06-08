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
)

// WatchBackend abstracts fswatch.Watcher for testing
type WatchBackend interface {
	WatchDirectory(dir string, fn fswatch.WatchCallback, recursive bool, ignore func(string) bool) (io.Closer, error)
	WatchFile(path string, fn fswatch.WatchCallback) (io.Closer, error)
}

type commandLineTestingWithWatchBackend interface {
	WatchBackend() WatchBackend
}

type fswatchBackend struct{ inner fswatch.Watcher }

func (b *fswatchBackend) WatchDirectory(dir string, fn fswatch.WatchCallback, recursive bool, ignore func(string) bool) (io.Closer, error) {
	var opts []fswatch.WatchOption
	if recursive {
		opts = append(opts, fswatch.WithRecursive())
	}
	if ignore != nil {
		opts = append(opts, fswatch.WithIgnore(ignore))
	}
	return b.inner.WatchDirectory(dir, fn, opts...)
}

func (b *fswatchBackend) WatchFile(path string, fn fswatch.WatchCallback) (io.Closer, error) {
	return b.inner.WatchFile(path, fn)
}

type watchedDir struct {
	closer    io.Closer
	recursive bool
}

type watchedFile struct {
	closer io.Closer
}

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
	commandLineRaw                 *collections.OrderedMap[string, any]
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

	backend         WatchBackend
	watchedDirs     map[string]*watchedDir   // dir path → watch state
	watchedFiles    map[string]*watchedFile  // file path → watch state
	seenFiles       map[tspath.Path]struct{} // all build dependencies (for event filtering)
	seenFileDirs    map[string]struct{}      // parent dirs of seen files (for watch registration)
	configMtimes    map[string]time.Time
	watchTerminated chan struct{}
	doCycleCh       chan struct{}
	debugLog        io.Writer // nil = silent; set via TS_WATCH_DEBUG

	changedMu       sync.Mutex
	changedPaths    map[string]fswatch.EventKind // event path → last event kind
	changedOverflow bool                         // true on ErrOverflow; forces full scan fallback
}

var _ tsc.Watcher = (*Watcher)(nil)

func createWatcher(
	sys tsc.System,
	configParseResult *tsoptions.ParsedCommandLine,
	compilerOptionsFromCommandLine *core.CompilerOptions,
	commandLineRaw *collections.OrderedMap[string, any],
	reportDiagnostic tsc.DiagnosticReporter,
	reportErrorSummary tsc.DiagnosticsReporter,
	testing tsc.CommandLineTesting,
) *Watcher {
	w := &Watcher{
		sys:                            sys,
		config:                         configParseResult,
		compilerOptionsFromCommandLine: compilerOptionsFromCommandLine,
		commandLineRaw:                 commandLineRaw,
		reportDiagnostic:               reportDiagnostic,
		reportErrorSummary:             reportErrorSummary,
		reportWatchStatus:              tsc.CreateWatchStatusReporter(sys, configParseResult.Locale(), configParseResult.CompilerOptions(), testing),
		testing:                        testing,
		sourceFileCache:                &collections.SyncMap[tspath.Path, *cachedSourceFile]{},
		watchTerminated:                make(chan struct{}),
		doCycleCh:                      make(chan struct{}, 1),
		watchedDirs:                    make(map[string]*watchedDir),
		watchedFiles:                   make(map[string]*watchedFile),
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	if t, ok := testing.(commandLineTestingWithWatchBackend); ok {
		w.backend = t.WatchBackend()
	}
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

	if w.testing == nil && w.backend == nil {
		fsw := fswatch.Default()
		w.backend = &fswatchBackend{inner: fsw}
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] using %s backend\n", fsw.Name())
		}
	}

	w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Starting_compilation_in_watch_mode))
	w.doBuild()
	w.mu.Unlock()

	if w.testing == nil {
		go func() {
			for {
				select {
				case <-w.watchTerminated:
					return
				case <-w.doCycleCh:
					w.DoCycle()
				}
			}
		}()
		// Block until a watch terminates (e.g. watched directory deleted),
		// then clean up.
		<-w.watchTerminated
		w.closeAllWatches()
	}
}

func (w *Watcher) reconcileWatches() {
	if w.backend == nil {
		return
	}

	cwd := w.sys.GetCurrentDirectory()

	desiredDirs := make(map[string]bool) // path → recursive

	// 1. Wildcard directories from tsconfig (recursive or non-recursive)
	if w.config.ConfigFile != nil {
		for dir, recursive := range w.config.WildcardDirectories() {
			realDir := w.sys.FS().Realpath(dir)
			desiredDirs[realDir] = recursive
		}
	}

	// 2. Parent directories of all seen files
	opts := w.comparePathsOptions()
	for dir := range w.seenFileDirs {
		if _, already := desiredDirs[dir]; already {
			continue
		}
		if w.coveredByRecursiveWildcard(dir, desiredDirs) {
			continue
		}
		if isStrictAncestorDir(cwd, dir, opts) {
			continue
		}
		desiredDirs[dir] = false
	}

	// 3. For no-config CLI mode, ensure CWD is watched
	if w.config.ConfigFile == nil && len(desiredDirs) == 0 {
		dir := w.sys.FS().Realpath(cwd)
		desiredDirs[dir] = false
	}

	// Compute desired file watches for config files
	desiredFiles := make(map[string]struct{})
	for _, path := range w.configFilePaths {
		realPath := w.sys.FS().Realpath(path)
		desiredFiles[realPath] = struct{}{}
	}

	// For no-config CLI mode, also watch the CLI-specified files directly
	if w.config.ConfigFile == nil {
		for _, path := range w.config.FileNames() {
			absPath := tspath.GetNormalizedAbsolutePath(path, cwd)
			realPath := w.sys.FS().Realpath(absPath)
			desiredFiles[realPath] = struct{}{}
		}
	}

	// Reconcile directory watches: collect stale ones, close outside critical path
	var staleDirClosers []io.Closer
	for dir, wd := range w.watchedDirs {
		wantRecursive, want := desiredDirs[dir]
		if !want || wd.recursive != wantRecursive {
			if w.debugLog != nil {
				if !want {
					fmt.Fprintf(w.debugLog, "[watch] closing stale dir watch: %s\n", dir)
				} else {
					fmt.Fprintf(w.debugLog, "[watch] recreating dir watch %s (recursive %v→%v)\n", dir, wd.recursive, wantRecursive)
				}
			}
			staleDirClosers = append(staleDirClosers, wd.closer)
			delete(w.watchedDirs, dir)
		}
	}
	for _, c := range staleDirClosers {
		c.Close()
	}
	for dir, recursive := range desiredDirs {
		if _, have := w.watchedDirs[dir]; have {
			continue
		}
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] watching directory %s (recursive=%v)\n", dir, recursive)
		}
		entry := &watchedDir{recursive: recursive}
		dirPath := dir
		cb := func(events []fswatch.Event, err error) {
			if err != nil && errors.Is(err, fswatch.ErrWatchTerminated) {
				w.handleWatchTerminated(dirPath, entry, true)
				return
			}
			w.onWatchEvents(events, err)
		}
		watchDir := dir
		watch, err := w.backend.WatchDirectory(watchDir, cb, recursive, shouldIgnoreWatchPath)
		for err != nil && watchDir != cwd && !isStrictAncestorDir(cwd, watchDir, opts) {
			parent := tspath.GetDirectoryPath(watchDir)
			if parent == watchDir {
				break
			}
			watchDir = parent
			if _, have := w.watchedDirs[watchDir]; have {
				break
			}
			watch, err = w.backend.WatchDirectory(watchDir, cb, recursive, shouldIgnoreWatchPath)
		}
		if err != nil {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] failed to watch directory %s: %v\n", dir, err)
			}
			continue
		}
		if watchDir != dir {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] watching ancestor %s for missing %s\n", watchDir, dir)
			}
			dirPath = watchDir
		}
		entry.closer = watch
		w.watchedDirs[dirPath] = entry
	}

	// Reconcile file watches
	var staleFileClosers []io.Closer
	for path, wf := range w.watchedFiles {
		if _, want := desiredFiles[path]; !want {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] closing stale file watch: %s\n", path)
			}
			staleFileClosers = append(staleFileClosers, wf.closer)
			delete(w.watchedFiles, path)
		}
	}
	for _, c := range staleFileClosers {
		c.Close()
	}
	for path := range desiredFiles {
		if _, have := w.watchedFiles[path]; have {
			continue
		}
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] watching file %s\n", path)
		}
		entry := &watchedFile{}
		filePath := path // capture for closure
		cb := func(events []fswatch.Event, err error) {
			if err != nil && errors.Is(err, fswatch.ErrWatchTerminated) {
				w.handleWatchTerminated(filePath, entry, false)
				return
			}
			w.onWatchEvents(events, err)
		}
		watch, err := w.backend.WatchFile(path, cb)
		if err != nil {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] failed to watch file %s: %v\n", path, err)
			}
			continue
		}
		entry.closer = watch
		w.watchedFiles[path] = entry
	}
}

func (w *Watcher) coveredByRecursiveWildcard(dir string, desiredDirs map[string]bool) bool {
	opts := w.comparePathsOptions()
	for wdir, recursive := range desiredDirs {
		if !recursive {
			continue
		}
		if tspath.ContainsPath(wdir, dir, opts) {
			return true
		}
	}
	return false
}

func (w *Watcher) comparePathsOptions() tspath.ComparePathsOptions {
	return tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: w.sys.FS().UseCaseSensitiveFileNames(),
		CurrentDirectory:          w.sys.GetCurrentDirectory(),
	}
}

// isStrictAncestorDir reports whether candidate is a strict ancestor of
// path (i.e. path is contained within candidate, but they are not equal).
func isStrictAncestorDir(path, candidate string, opts tspath.ComparePathsOptions) bool {
	return path != candidate && tspath.ContainsPath(candidate, path, opts)
}

func (w *Watcher) closeAllWatches() {
	w.mu.Lock()
	dirs := make([]io.Closer, 0, len(w.watchedDirs))
	for dir, wd := range w.watchedDirs {
		dirs = append(dirs, wd.closer)
		delete(w.watchedDirs, dir)
	}
	files := make([]io.Closer, 0, len(w.watchedFiles))
	for path, wf := range w.watchedFiles {
		files = append(files, wf.closer)
		delete(w.watchedFiles, path)
	}
	w.mu.Unlock()
	for _, c := range dirs {
		c.Close()
	}
	for _, c := range files {
		c.Close()
	}
}

func (w *Watcher) handleWatchTerminated(path string, identity any, isDir bool) {
	if w.debugLog != nil {
		fmt.Fprintf(w.debugLog, "[watch] watch terminated: %s\n", path)
	}
	var staleCloser io.Closer
	w.mu.Lock()
	if isDir {
		if wd, ok := w.watchedDirs[path]; ok && wd == identity {
			staleCloser = wd.closer
			delete(w.watchedDirs, path)
		}
	} else {
		if wf, ok := w.watchedFiles[path]; ok && wf == identity {
			staleCloser = wf.closer
			delete(w.watchedFiles, path)
		}
	}
	w.mu.Unlock()
	if staleCloser != nil {
		staleCloser.Close()
	}
	w.changedMu.Lock()
	w.changedOverflow = true
	w.changedMu.Unlock()
	w.signalDoCycle()
}

func shouldIgnoreWatchPath(path string) bool {
	p := tspath.NormalizeSlashes(path)
	return strings.HasSuffix(p, "/.git") ||
		strings.Contains(p, "/.git/") ||
		strings.Contains(p, "/node_modules/.") ||
		strings.Contains(p, "/.#")
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

func (w *Watcher) signalDoCycle() {
	select {
	case w.doCycleCh <- struct{}{}:
		// Signal sent; the DoCycle goroutine will pick it up.
	case <-w.watchTerminated:
		// Watch has terminated; drop the signal.
	default:
		// A signal is already pending; coalesced.
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
		// Filter fswatch events against known dependencies
		if w.isRelevantChange(changedPaths) {
			w.evictChangedSourceFiles(changedPaths)
		} else {
			if w.debugLog != nil {
				fmt.Fprintf(w.debugLog, "[watch] DoCycle: %d event(s) not relevant to compilation, skipping rebuild\n", len(changedPaths))
			}
			if w.testing != nil {
				w.testing.OnProgram(w.program)
			}
			return
		}
	} else if overflow {
		// Overflow: evict the entire source file cache to force re-build
		w.sourceFileCache = &collections.SyncMap[tspath.Path, *cachedSourceFile]{}
	} else if !hasEvents && !w.configModified {
		// No events and no config change
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] DoCycle: no events, skipping\n")
		}
		if w.testing != nil {
			w.testing.OnProgram(w.program)
		}
		return
	}

	w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.File_change_detected_Starting_incremental_compilation))
	w.doBuild()
}

func (w *Watcher) isRelevantChange(changedPaths map[string]fswatch.EventKind) bool {
	caseSensitive := w.sys.FS().UseCaseSensitiveFileNames()
	cwd := w.sys.GetCurrentDirectory()
	for eventPath := range changedPaths {
		p := tspath.ToPath(eventPath, cwd, caseSensitive)
		if _, ok := w.seenFiles[p]; ok {
			return true
		}
		if w.config.ConfigFile != nil && w.config.PossiblyMatchesFileName(eventPath) {
			return true
		}
	}
	return false
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

	caseSensitive := w.sys.FS().UseCaseSensitiveFileNames()
	cwd := w.sys.GetCurrentDirectory()
	seenSlice := tfs.SeenFiles.ToSlice()
	w.seenFiles = make(map[tspath.Path]struct{}, len(seenSlice)*2)
	w.seenFileDirs = make(map[string]struct{}, len(seenSlice))
	for _, p := range seenSlice {
		w.seenFiles[tspath.ToPath(p, cwd, caseSensitive)] = struct{}{}
		if rp := w.sys.FS().Realpath(p); rp != p {
			w.seenFiles[tspath.ToPath(rp, cwd, caseSensitive)] = struct{}{}
		}
		dir := tspath.GetDirectoryPath(p)
		if tspath.IsRootedDiskPath(dir) {
			w.seenFileDirs[dir] = struct{}{}
		}
	}
	for _, sf := range w.program.GetProgram().GetSourceFiles() {
		realFile := w.sys.FS().Realpath(sf.FileName())
		if realFile != sf.FileName() {
			realFileDir := tspath.GetDirectoryPath(realFile)
			if tspath.IsRootedDiskPath(realFileDir) {
				w.seenFileDirs[realFileDir] = struct{}{}
			}
		}
	}

	w.configMtimes = make(map[string]time.Time, len(w.configFilePaths))
	for _, cfgPath := range w.configFilePaths {
		if s := w.sys.FS().Stat(cfgPath); s != nil {
			w.configMtimes[cfgPath] = s.ModTime()
		}
	}

	w.reconcileWatches()
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
			oldMtime, ok := w.configMtimes[path]
			s := w.sys.FS().Stat(path)
			if !ok {
				if s != nil {
					changed = true
					break
				}
			} else if s == nil || !s.ModTime().Equal(oldMtime) {
				changed = true
				break
			}
		}
		if !changed {
			return false
		}
	}

	configParseResult, parseErr := w.tryParseConfigFile()
	if parseErr != nil {
		if w.debugLog != nil {
			fmt.Fprintf(w.debugLog, "[watch] config parse recovered from panic: %v\n", parseErr)
		}
		return false
	}
	if configParseResult == nil {
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
	return false
}

func (w *Watcher) tryParseConfigFile() (result *tsoptions.ParsedCommandLine, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			recovered = r
		}
	}()

	extendedConfigCache := &tsc.ExtendedConfigCache{}
	configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(w.configFileName, w.compilerOptionsFromCommandLine, w.commandLineRaw, w.sys, extendedConfigCache)
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
		return nil, nil
	}
	w.extendedConfigCache = extendedConfigCache
	return configParseResult, nil
}
