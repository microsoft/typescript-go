package execute

import (
	"reflect"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/trackingvfs"
)

const watchDebounceWait = 250 * time.Millisecond

type FileWatcher struct {
	fs                  vfs.FS
	pollInterval        time.Duration
	testing             bool
	callback            func()
	watchState          map[string]trackingvfs.WatchEntry
	wildcardDirectories map[string]bool
}

func newFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func()) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		pollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
	}
}

func (fw *FileWatcher) updateWatchedFiles(tfs *trackingvfs.FS) {
	fw.watchState = make(map[string]trackingvfs.WatchEntry)
	tfs.SeenFiles.Range(func(fn string) bool {
		if s := fw.fs.Stat(fn); s != nil {
			fw.watchState[fn] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			fw.watchState[fn] = trackingvfs.WatchEntry{Exists: false, ChildCount: -1}
		}
		return true
	})
	for dir, recursive := range fw.wildcardDirectories {
		if !recursive {
			continue
		}
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entries := fw.fs.GetAccessibleEntries(path)
			count := len(entries.Files) + len(entries.Directories)
			if existing, ok := fw.watchState[path]; ok {
				existing.ChildCount = count
				fw.watchState[path] = existing
			} else {
				if s := fw.fs.Stat(path); s != nil {
					fw.watchState[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
				}
			}
			return nil
		})
	}
}

func (fw *FileWatcher) WaitForSettled(now func() time.Time) {
	if fw.testing {
		return
	}
	current := fw.currentState()
	settledAt := now()
	for now().Sub(settledAt) < watchDebounceWait {
		time.Sleep(fw.pollInterval)
		if fw.HasChanges(current) {
			current = fw.currentState()
			settledAt = now()
		}
	}
}

func (fw *FileWatcher) currentState() map[string]trackingvfs.WatchEntry {
	state := make(map[string]trackingvfs.WatchEntry, len(fw.watchState))
	for path := range fw.watchState {
		if s := fw.fs.Stat(path); s != nil {
			state[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: -1}
		} else {
			state[path] = trackingvfs.WatchEntry{Exists: false, ChildCount: -1}
		}
	}
	for dir, recursive := range fw.wildcardDirectories {
		if !recursive {
			continue
		}
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entries := fw.fs.GetAccessibleEntries(path)
			count := len(entries.Files) + len(entries.Directories)
			if existing, ok := state[path]; ok {
				existing.ChildCount = count
				state[path] = existing
			} else {
				if s := fw.fs.Stat(path); s != nil {
					state[path] = trackingvfs.WatchEntry{ModTime: s.ModTime(), Exists: true, ChildCount: count}
				}
			}
			return nil
		})
	}
	return state
}

func (fw *FileWatcher) HasChanges(baseline map[string]trackingvfs.WatchEntry) bool {
	for path, old := range baseline {
		s := fw.fs.Stat(path)
		if !old.Exists {
			if s != nil {
				return true
			}
		} else {
			if s == nil || !s.ModTime().Equal(old.ModTime) {
				return true
			}
		}
	}
	for dir, recursive := range fw.wildcardDirectories {
		if !recursive {
			continue
		}
		found := false
		_ = fw.fs.WalkDir(dir, func(path string, d vfs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			entry, ok := baseline[path]
			if !ok {
				found = true
				return vfs.SkipAll
			}
			if entry.ChildCount >= 0 {
				entries := fw.fs.GetAccessibleEntries(path)
				if len(entries.Files)+len(entries.Directories) != entry.ChildCount {
					found = true
					return vfs.SkipAll
				}
			}
			return nil
		})
		if found {
			return true
		}
	}
	return false
}

func (fw *FileWatcher) Run(now func() time.Time) {
	for {
		time.Sleep(fw.pollInterval)
		if fw.watchState == nil || fw.HasChanges(fw.watchState) {
			fw.WaitForSettled(now)
			fw.callback()
		}
	}
}

type cachedSourceFile struct {
	file    *ast.SourceFile
	modTime time.Time
}

type watchCompilerHost struct {
	inner compiler.CompilerHost
	cache *collections.SyncMap[tspath.Path, *cachedSourceFile]
}

var _ compiler.CompilerHost = (*watchCompilerHost)(nil)

func (h *watchCompilerHost) FS() vfs.FS                  { return h.inner.FS() }
func (h *watchCompilerHost) DefaultLibraryPath() string  { return h.inner.DefaultLibraryPath() }
func (h *watchCompilerHost) GetCurrentDirectory() string { return h.inner.GetCurrentDirectory() }
func (h *watchCompilerHost) Trace(msg *diagnostics.Message, args ...any) {
	h.inner.Trace(msg, args...)
}

func (h *watchCompilerHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	return h.inner.GetResolvedProjectReference(fileName, path)
}

func (h *watchCompilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	info := h.inner.FS().Stat(opts.FileName)

	if cached, ok := h.cache.Load(opts.Path); ok {
		if info != nil && info.ModTime().Equal(cached.modTime) {
			return cached.file
		}
	}

	file := h.inner.GetSourceFile(opts)
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
	fileWatcher     *FileWatcher
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
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	w.fileWatcher = newFileWatcher(
		sys.FS(),
		w.config.ParsedConfig.WatchOptions.WatchInterval(),
		testing != nil,
		w.DoCycle,
	)
	return w
}

func (w *Watcher) start() {
	w.extendedConfigCache = &tsc.ExtendedConfigCache{}
	host := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	w.program = incremental.ReadBuildInfoProgram(w.config, incremental.NewBuildInfoReader(host), host)

	if w.configFileName != "" {
		w.configFilePaths = append([]string{w.configFileName}, w.config.ExtendedSourceFiles()...)
	}

	w.reportWatchStatus(ast.NewCompilerDiagnostic(diagnostics.Starting_compilation_in_watch_mode))
	w.doBuild()

	if w.testing == nil {
		w.fileWatcher.Run(w.sys.Now)
	}
}

func (w *Watcher) DoCycle() {
	if w.hasErrorsInTsConfig() {
		return
	}
	if w.fileWatcher.watchState != nil && !w.configModified && !w.fileWatcher.HasChanges(w.fileWatcher.watchState) {
		if w.testing != nil {
			w.testing.OnProgram(w.program)
		}
		return
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
	host := &watchCompilerHost{inner: innerHost, cache: w.sourceFileCache}

	if w.config.ConfigFile != nil {
		wildcardDirs := w.config.WildcardDirectories()
		for dir := range wildcardDirs {
			tfs.SeenFiles.Add(dir)
		}
		w.fileWatcher.wildcardDirectories = wildcardDirs
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
	w.fileWatcher.updateWatchedFiles(tfs)
	w.fileWatcher.pollInterval = w.config.ParsedConfig.WatchOptions.WatchInterval()
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

func (w *Watcher) hasErrorsInTsConfig() bool {
	if w.configFileName == "" {
		return false
	}

	if !w.configHasErrors && len(w.configFilePaths) > 0 {
		changed := false
		for _, path := range w.configFilePaths {
			if old, ok := w.fileWatcher.watchState[path]; ok {
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

// Testing helpers — exported for use by test packages

func (w *Watcher) HasWatchedFilesChanged() bool {
	return w.fileWatcher.HasChanges(w.fileWatcher.watchState)
}

func (w *Watcher) WatchStateLen() int {
	return len(w.fileWatcher.watchState)
}

func (w *Watcher) WatchStateHas(path string) bool {
	_, ok := w.fileWatcher.watchState[path]
	return ok
}

func (w *Watcher) DebugWatchState(fn func(path string, modTime time.Time, exists bool)) {
	for path, entry := range w.fileWatcher.watchState {
		fn(path, entry.ModTime, entry.Exists)
	}
}
