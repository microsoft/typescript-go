package execute

import (
	"fmt"
	"reflect"
	"slices"
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
)

const watchDebounceWait = 250 * time.Millisecond

type trackingFS struct {
	inner     vfs.FS
	seenFiles collections.SyncSet[string]
}

func (fs *trackingFS) ReadFile(path string) (string, bool) {
	fs.seenFiles.Add(path)
	return fs.inner.ReadFile(path)
}

func (fs *trackingFS) FileExists(path string) bool {
	fs.seenFiles.Add(path)
	return fs.inner.FileExists(path)
}
func (fs *trackingFS) UseCaseSensitiveFileNames() bool { return fs.inner.UseCaseSensitiveFileNames() }
func (fs *trackingFS) WriteFile(path string, data string) error {
	return fs.inner.WriteFile(path, data)
}
func (fs *trackingFS) Remove(path string) error { return fs.inner.Remove(path) }
func (fs *trackingFS) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return fs.inner.Chtimes(path, aTime, mTime)
}

func (fs *trackingFS) DirectoryExists(path string) bool {
	fs.seenFiles.Add(path)
	return fs.inner.DirectoryExists(path)
}

func (fs *trackingFS) GetAccessibleEntries(path string) vfs.Entries {
	fs.seenFiles.Add(path)
	return fs.inner.GetAccessibleEntries(path)
}

func (fs *trackingFS) Stat(path string) vfs.FileInfo {
	fs.seenFiles.Add(path)
	return fs.inner.Stat(path)
}

func (fs *trackingFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.inner.WalkDir(root, walkFn)
}
func (fs *trackingFS) Realpath(path string) string { return fs.inner.Realpath(path) }

type WatchEntry struct {
	modTime time.Time
	exists  bool
}

type FileWatcher struct {
	fs                  vfs.FS
	pollInterval        time.Duration
	testing             bool
	callback            func()
	watchState          map[string]WatchEntry
	wildcardDirectories map[string]bool // dir path -> recursive flag
}

func newFileWatcher(fs vfs.FS, pollInterval time.Duration, testing bool, callback func()) *FileWatcher {
	return &FileWatcher{
		fs:           fs,
		pollInterval: pollInterval,
		testing:      testing,
		callback:     callback,
	}
}

func (fw *FileWatcher) updateWatchedFiles(tfs *trackingFS) {
	fw.watchState = make(map[string]WatchEntry)
	tfs.seenFiles.Range(func(fn string) bool {
		if s := fw.fs.Stat(fn); s != nil {
			fw.watchState[fn] = WatchEntry{modTime: s.ModTime(), exists: true}
		} else {
			fw.watchState[fn] = WatchEntry{exists: false}
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
			if _, ok := fw.watchState[path]; !ok {
				if s := fw.fs.Stat(path); s != nil {
					fw.watchState[path] = WatchEntry{modTime: s.ModTime(), exists: true}
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

func (fw *FileWatcher) currentState() map[string]WatchEntry {
	state := make(map[string]WatchEntry, len(fw.watchState))
	for path := range fw.watchState {
		if s := fw.fs.Stat(path); s != nil {
			state[path] = WatchEntry{modTime: s.ModTime(), exists: true}
		} else {
			state[path] = WatchEntry{exists: false}
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
			if _, ok := state[path]; !ok {
				if s := fw.fs.Stat(path); s != nil {
					state[path] = WatchEntry{modTime: s.ModTime(), exists: true}
				}
			}
			return nil
		})
	}
	return state
}

func (fw *FileWatcher) HasChanges(baseline map[string]WatchEntry) bool {
	for path, old := range baseline {
		s := fw.fs.Stat(path)
		if !old.exists {
			if s != nil {
				return true
			}
		} else {
			if s == nil || !s.ModTime().Equal(old.modTime) {
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
			if _, ok := baseline[path]; !ok {
				found = true
				return vfs.SkipAll
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
	if cached, ok := h.cache.Load(opts.Path); ok {
		info := h.inner.FS().Stat(opts.FileName)
		if info != nil && info.ModTime().Equal(cached.modTime) {
			return cached.file
		}
	}

	info := h.inner.FS().Stat(opts.FileName)
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
	testing                        tsc.CommandLineTesting

	program             *incremental.Program
	extendedConfigCache *tsc.ExtendedConfigCache
	configModified      bool
	configHasErrors     bool
	configFilePaths     []string

	sourceFileCache *collections.SyncMap[tspath.Path, *cachedSourceFile]
	fileWatcher     *FileWatcher
}

var (
	_ tsc.Watcher = (*Watcher)(nil)
	_ vfs.FS      = (*trackingFS)(nil)
)

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

	w.doBuild(false)

	if w.testing == nil {
		w.fileWatcher.Run(w.sys.Now)
	}
}

func (w *Watcher) DoCycle() {
	if w.hasErrorsInTsConfig() {
		return
	}
	if w.fileWatcher.watchState != nil && !w.configModified && !w.fileWatcher.HasChanges(w.fileWatcher.watchState) {
		if w.config.ConfigFile != nil && len(w.config.WildcardDirectories()) > 0 {
			updated := w.config.ReloadFileNamesOfParsedCommandLine(w.sys.FS())
			if !slices.Equal(w.config.FileNames(), updated.FileNames()) {
				w.config = updated
				w.doBuild(true)
				return
			}
		}
		if w.testing != nil {
			w.testing.OnProgram(w.program)
		}
		return
	}

	w.doBuild(false)
}

func (w *Watcher) doBuild(fileNamesReloaded bool) {
	if w.configModified {
		w.sourceFileCache = &collections.SyncMap[tspath.Path, *cachedSourceFile]{}
	}

	cached := cachedvfs.From(w.sys.FS())
	tfs := &trackingFS{inner: cached}
	innerHost := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), tfs, w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	host := &watchCompilerHost{inner: innerHost, cache: w.sourceFileCache}

	if w.config.ConfigFile != nil {
		wildcardDirs := w.config.WildcardDirectories()
		for dir := range wildcardDirs {
			tfs.seenFiles.Add(dir)
		}
		w.fileWatcher.wildcardDirectories = wildcardDirs
		if len(wildcardDirs) > 0 && !fileNamesReloaded {
			w.config = w.config.ReloadFileNamesOfParsedCommandLine(w.sys.FS())
		}
	}
	for _, path := range w.configFilePaths {
		tfs.seenFiles.Add(path)
	}

	fmt.Fprintln(w.sys.Writer(), "build starting at", w.sys.Now().Format("03:04:05 PM"))
	timeStart := w.sys.Now()

	w.program = incremental.NewProgram(compiler.NewProgram(compiler.ProgramOptions{
		Config: w.config,
		Host:   host,
	}), w.program, nil, w.testing != nil)

	w.compileAndEmit()
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

	fmt.Fprintf(w.sys.Writer(), "build finished in %.3fs\n", w.sys.Now().Sub(timeStart).Seconds())

	if w.testing != nil {
		w.testing.OnProgram(w.program)
	}
}

func (w *Watcher) compileAndEmit() {
	tsc.EmitFilesAndReportErrors(tsc.EmitInput{
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
				if !old.exists {
					if s != nil {
						changed = true
						break
					}
				} else {
					if s == nil || !s.ModTime().Equal(old.modTime) {
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
		fn(path, entry.modTime, entry.exists)
	}
}
