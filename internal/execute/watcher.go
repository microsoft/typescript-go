package execute

import (
	"fmt"
	"reflect"
	"time"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs"
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
func (fs *trackingFS) Stat(path string) vfs.FileInfo { return fs.inner.Stat(path) }
func (fs *trackingFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.inner.WalkDir(root, walkFn)
}
func (fs *trackingFS) Realpath(path string) string { return fs.inner.Realpath(path) }

type watchEntry struct {
	modTime time.Time
	exists  bool
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

	watchState map[string]watchEntry
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
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	return w
}

func (w *Watcher) start() {
	w.extendedConfigCache = &tsc.ExtendedConfigCache{}
	host := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	w.program = incremental.ReadBuildInfoProgram(w.config, incremental.NewBuildInfoReader(host), host)

	w.doBuild()

	if w.testing == nil {
		for {
			time.Sleep(w.pollInterval())
			w.DoCycle()
		}
	}
}

func (w *Watcher) DoCycle() {
	if w.hasErrorsInTsConfig() {
		// these are unrecoverable errors--report them and do not build
		return
	}
	if w.watchState != nil && !w.configModified && !w.hasWatchedFilesChanged() {
		if w.testing != nil {
			w.testing.OnProgram(w.program)
		}
		return
	}

	if w.testing == nil {
		w.refreshWatchState()
		settledAt := w.sys.Now()
		for w.sys.Now().Sub(settledAt) < watchDebounceWait {
			time.Sleep(w.pollInterval())
			if w.hasWatchedFilesChanged() {
				w.refreshWatchState()
				settledAt = w.sys.Now()
			}
		}
	}

	w.doBuild()
}

func (w *Watcher) doBuild() {
	tfs := &trackingFS{inner: w.sys.FS()}
	host := compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), tfs, w.sys.DefaultLibraryPath(), w.extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))

	fmt.Fprintln(w.sys.Writer(), "build starting at", w.sys.Now().Format("03:04:05 PM"))
	timeStart := w.sys.Now()

	w.program = incremental.NewProgram(compiler.NewProgram(compiler.ProgramOptions{
		Config: w.config,
		Host:   host,
	}), w.program, nil, w.testing != nil)

	w.compileAndEmit()
	w.buildWatchState(tfs)
	w.configModified = false
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

	extendedConfigCache := &tsc.ExtendedConfigCache{}
	configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(w.configFileName, w.compilerOptionsFromCommandLine, nil, w.sys, extendedConfigCache)
	if len(errors) > 0 {
		if !w.configHasErrors {
			for _, e := range errors {
				w.reportDiagnostic(e)
			}
			w.configHasErrors = true
		}
		return true
	}
	if w.configHasErrors {
		w.configModified = true
	}
	w.configHasErrors = false
	if !reflect.DeepEqual(w.config.ParsedConfig, configParseResult.ParsedConfig) {
		w.configModified = true
	}
	w.config = configParseResult
	w.extendedConfigCache = extendedConfigCache
	return false
}

func (w *Watcher) hasWatchedFilesChanged() bool {
	for path, old := range w.watchState {
		s := w.sys.FS().Stat(path)
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
	return false
}

func (w *Watcher) buildWatchState(tfs *trackingFS) {
	w.watchState = make(map[string]watchEntry)
	tfs.seenFiles.Range(func(fn string) bool {
		if s := w.sys.FS().Stat(fn); s != nil {
			w.watchState[fn] = watchEntry{modTime: s.ModTime(), exists: true}
		} else {
			w.watchState[fn] = watchEntry{exists: false}
		}
		return true
	})
}

func (w *Watcher) refreshWatchState() {
	for path := range w.watchState {
		if s := w.sys.FS().Stat(path); s != nil {
			w.watchState[path] = watchEntry{modTime: s.ModTime(), exists: true}
		} else {
			w.watchState[path] = watchEntry{exists: false}
		}
	}
}

func (w *Watcher) pollInterval() time.Duration {
	return w.config.ParsedConfig.WatchOptions.WatchInterval()
}

// Testing helpers — exported for use by test packages.

func (w *Watcher) HasWatchedFilesChanged() bool {
	return w.hasWatchedFilesChanged()
}

func (w *Watcher) RefreshWatchState() {
	w.refreshWatchState()
}

func (w *Watcher) WatchStateLen() int {
	return len(w.watchState)
}

func (w *Watcher) WatchStateHas(path string) bool {
	_, ok := w.watchState[path]
	return ok
}

func (w *Watcher) DebugWatchState(fn func(path string, modTime time.Time, exists bool)) {
	for path, entry := range w.watchState {
		fn(path, entry.modTime, entry.exists)
	}
}
