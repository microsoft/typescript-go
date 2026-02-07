package execute

import (
	"fmt"
	"reflect"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type Watcher struct {
	sys                            tsc.System
	configFileName                 string
	config                         *tsoptions.ParsedCommandLine
	compilerOptionsFromCommandLine *core.CompilerOptions
	reportDiagnostic               tsc.DiagnosticReporter
	reportErrorSummary             tsc.DiagnosticsReporter
	testing                        tsc.CommandLineTesting

	host           compiler.CompilerHost
	program        *incremental.Program
	prevModified   map[string]time.Time
	configModified bool
	deletedFiles   []string // source files deleted since last cycle
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
		testing:                        testing,
		// reportWatchStatus: createWatchStatusReporter(sys, configParseResult.CompilerOptions().Pretty),
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	return w
}

func (w *Watcher) start() {
	w.host = compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), nil, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	w.program = incremental.ReadBuildInfoProgram(w.config, incremental.NewBuildInfoReader(w.host), w.host)

	if w.testing == nil {
		watchInterval := w.config.ParsedConfig.WatchOptions.WatchInterval()
		for {
			w.DoCycle()
			time.Sleep(watchInterval)
		}
	} else {
		// Initial compilation in test mode
		w.DoCycle()
	}
}

func (w *Watcher) DoCycle() {
	// if this function is updated, make sure to update `RunWatchCycle` in export_test.go as needed

	if w.hasErrorsInTsConfig() {
		// these are unrecoverable errors--report them and do not build
		return
	}
	// updateProgram()
	w.program = incremental.NewProgram(compiler.NewProgram(compiler.ProgramOptions{
		Config:           w.config,
		Host:             w.host,
		JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
	}), w.program, nil, w.testing != nil)

	if w.hasBeenModified(w.program.GetProgram()) {
		fmt.Fprintln(w.sys.Writer(), "build starting at", w.sys.Now().Format("03:04:05 PM"))
		timeStart := w.sys.Now()
		w.compileAndEmit()
		w.cleanupDeletedOutputs()
		fmt.Fprintf(w.sys.Writer(), "build finished in %.3fs\n", w.sys.Now().Sub(timeStart).Seconds())
	} else {
		// print something???
		// fmt.Fprintln(w.sys.Writer(), "no changes detected at ", w.sys.Now())
	}
	if w.testing != nil {
		w.testing.OnProgram(w.program)
	}
}

func (w *Watcher) compileAndEmit() {
	// !!! output/error reporting is currently the same as non-watch mode
	// diagnostics, emitResult, exitStatus :=
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
	// only need to check and reparse tsconfig options/update host if we are watching a config file
	extendedConfigCache := &tsc.ExtendedConfigCache{}
	if w.configFileName != "" {
		// !!! need to check that this merges compileroptions correctly. This differs from non-watch, since we allow overriding of previous options
		configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(w.configFileName, w.compilerOptionsFromCommandLine, nil, w.sys, extendedConfigCache)
		if len(errors) > 0 {
			for _, e := range errors {
				w.reportDiagnostic(e)
			}
			return true
		}
		// CompilerOptions contain fields which should not be compared; clone to get a copy without those set.
		if !reflect.DeepEqual(w.config.CompilerOptions().Clone(), configParseResult.CompilerOptions().Clone()) {
			// fmt.Fprintln(w.sys.Writer(), "build triggered due to config change")
			w.configModified = true
		}
		w.config = configParseResult
	}
	w.host = compiler.NewCompilerHost(w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), extendedConfigCache, getTraceFromSys(w.sys, w.config.Locale(), w.testing))
	return false
}

func (w *Watcher) hasBeenModified(program *compiler.Program) bool {
	// checks watcher's snapshot against program file modified times
	currState := map[string]time.Time{}
	filesModified := w.configModified
	for _, sourceFile := range program.SourceFiles() {
		fileName := sourceFile.FileName()
		s := w.sys.FS().Stat(fileName)
		if s == nil {
			// do nothing; if file is in program.SourceFiles() but is not found when calling Stat, file has been very recently deleted.
			// deleted files are handled outside of this loop
			continue
		}
		currState[fileName] = s.ModTime()
		if !filesModified {
			if currState[fileName] != w.prevModified[fileName] {
				// fmt.Fprint(w.sys.Writer(), "build triggered from ", fileName, ": ", w.prevModified[fileName], " -> ", currState[fileName], "\n")
				filesModified = true
			}
			// catch cases where no files are modified, but some were deleted
			delete(w.prevModified, fileName)
		}
	}
	if !filesModified && len(w.prevModified) > 0 {
		// fmt.Fprintln(w.sys.Writer(), "build triggered due to deleted file")
		filesModified = true
	}

	// Capture names of deleted source files before overwriting prevModified.
	// When filesModified was set early (e.g. due to a modified file), the loop
	// above may not have removed all current files from prevModified. We must
	// compare against currState to find files that are truly gone from disk.
	w.deletedFiles = w.deletedFiles[:0]
	for fileName := range w.prevModified {
		if _, exists := currState[fileName]; !exists {
			w.deletedFiles = append(w.deletedFiles, fileName)
		}
	}

	w.prevModified = currState

	// reset state for next cycle
	w.configModified = false
	return filesModified
}

// cleanupDeletedOutputs removes output files (.js, .js.map, .d.ts, .d.ts.map)
// corresponding to source files that were deleted since the last watch cycle.
// This prevents stale outputs from remaining in the output directory.
// Fixes: https://github.com/microsoft/TypeScript/issues/16057
func (w *Watcher) cleanupDeletedOutputs() {
	if len(w.deletedFiles) == 0 {
		return
	}

	options := w.config.CompilerOptions()
	program := w.program.GetProgram()

	// Skip cleanup when no emit is expected.
	if options.NoEmit.IsTrue() {
		return
	}

	for _, deletedFile := range w.deletedFiles {
		// Compute and remove JS output file.
		jsPath := outputpaths.GetOutputJSFileName(deletedFile, options, program)
		if jsPath != "" {
			w.sys.FS().Remove(jsPath)

			// Remove source map if enabled.
			if mapPath := outputpaths.GetSourceMapFilePath(jsPath, options); mapPath != "" {
				w.sys.FS().Remove(mapPath)
			}
		}

		// Compute and remove declaration output files.
		if options.GetEmitDeclarations() {
			dtsPath := outputpaths.GetOutputDeclarationFileNameWorker(deletedFile, options, program)
			if dtsPath != "" {
				w.sys.FS().Remove(dtsPath)

				// Remove declaration map if enabled.
				if options.GetAreDeclarationMapsEnabled() {
					w.sys.FS().Remove(dtsPath + ".map")
				}
			}
		}
	}
}
