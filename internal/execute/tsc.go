package execute

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/pprof"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type cbType = func(p any) any

func applyBulkEdits(text string, edits []core.TextChange) string {
	b := strings.Builder{}
	b.Grow(len(text))
	lastEnd := 0
	for _, e := range edits {
		start := e.TextRange.Pos()
		if start != lastEnd {
			b.WriteString(text[lastEnd:e.TextRange.Pos()])
		}
		b.WriteString(e.NewText)

		lastEnd = e.TextRange.End()
	}
	b.WriteString(text[lastEnd:])

	return b.String()
}

func CommandLine(sys System, commandLineArgs []string) ExitStatus {
	status, _, watcher := commandLineWorker(sys, commandLineArgs)
	if watcher == nil {
		return status
	}
	return start(watcher)
}

func commandLineWorker(sys System, commandLineArgs []string) (ExitStatus, *tsoptions.ParsedCommandLine, *watcher) {
	if len(commandLineArgs) > 0 {
		// !!! build mode
		switch strings.ToLower(commandLineArgs[0]) {
		case "-b", "--b", "-build", "--build":
			fmt.Fprint(sys.Writer(), "Build mode is currently unsupported."+sys.NewLine())
			sys.EndWrite()
			return ExitStatusNotImplemented, nil, nil
			// case "-f":
			// 	return fmtMain(sys, commandLineArgs[1], commandLineArgs[1])
		}
	}

	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	status, watcher := tscCompilation(sys, parsedCommandLine)
	return status, parsedCommandLine, watcher
}

func fmtMain(sys System, input, output string) ExitStatus {
	ctx := format.WithFormatCodeSettings(context.Background(), format.GetDefaultFormatCodeSettings(sys.NewLine()), sys.NewLine())
	input = string(tspath.ToPath(input, sys.GetCurrentDirectory(), sys.FS().UseCaseSensitiveFileNames()))
	output = string(tspath.ToPath(output, sys.GetCurrentDirectory(), sys.FS().UseCaseSensitiveFileNames()))
	fileContent, ok := sys.FS().ReadFile(input)
	if !ok {
		fmt.Fprint(sys.Writer(), "File not found: "+input+sys.NewLine())
		return ExitStatusNotImplemented
	}
	text := fileContent
	pathified := tspath.ToPath(input, sys.GetCurrentDirectory(), true)
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName:         string(pathified),
		Path:             pathified,
		JSDocParsingMode: ast.JSDocParsingModeParseAll,
	}, text, core.GetScriptKindFromFileName(string(pathified)))
	edits := format.FormatDocument(ctx, sourceFile)
	newText := applyBulkEdits(text, edits)

	if err := sys.FS().WriteFile(output, newText, false); err != nil {
		fmt.Fprint(sys.Writer(), err.Error()+sys.NewLine())
		return ExitStatusNotImplemented
	}
	return ExitStatusSuccess
}

func tscCompilation(sys System, commandLine *tsoptions.ParsedCommandLine) (ExitStatus, *watcher) {
	configFileName := ""
	reportDiagnostic := createDiagnosticReporter(sys, commandLine.CompilerOptions())
	// if commandLine.Options().Locale != nil

	if len(commandLine.Errors) > 0 {
		for _, e := range commandLine.Errors {
			reportDiagnostic(e)
		}
		return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
	}

	if pprofDir := commandLine.CompilerOptions().PprofDir; pprofDir != "" {
		// !!! stderr?
		profileSession := pprof.BeginProfiling(pprofDir, sys.Writer())
		defer profileSession.Stop()
	}

	if commandLine.CompilerOptions().Init.IsTrue() {
		return ExitStatusNotImplemented, nil
	}

	if commandLine.CompilerOptions().Version.IsTrue() {
		printVersion(sys)
		return ExitStatusSuccess, nil
	}

	if commandLine.CompilerOptions().Help.IsTrue() || commandLine.CompilerOptions().All.IsTrue() {
		printHelp(sys, commandLine)
		return ExitStatusSuccess, nil
	}

	if commandLine.CompilerOptions().Watch.IsTrue() && commandLine.CompilerOptions().ListFilesOnly.IsTrue() {
		reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Options_0_and_1_cannot_be_combined, "watch", "listFilesOnly"))
		return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
	}

	if commandLine.CompilerOptions().Project != "" {
		if len(commandLine.FileNames()) != 0 {
			reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Option_project_cannot_be_mixed_with_source_files_on_a_command_line))
			return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
		}

		fileOrDirectory := tspath.NormalizePath(commandLine.CompilerOptions().Project)
		if sys.FS().DirectoryExists(fileOrDirectory) {
			configFileName = tspath.CombinePaths(fileOrDirectory, "tsconfig.json")
			if !sys.FS().FileExists(configFileName) {
				reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Cannot_find_a_tsconfig_json_file_at_the_current_directory_Colon_0, configFileName))
				return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
			}
		} else {
			configFileName = fileOrDirectory
			if !sys.FS().FileExists(configFileName) {
				reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.The_specified_path_does_not_exist_Colon_0, fileOrDirectory))
				return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
			}
		}
	} else if len(commandLine.FileNames()) == 0 {
		searchPath := tspath.NormalizePath(sys.GetCurrentDirectory())
		configFileName = findConfigFile(searchPath, sys.FS().FileExists, "tsconfig.json")
	}

	if configFileName == "" && len(commandLine.FileNames()) == 0 {
		if commandLine.CompilerOptions().ShowConfig.IsTrue() {
			reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Cannot_find_a_tsconfig_json_file_at_the_current_directory_Colon_0, tspath.NormalizePath(sys.GetCurrentDirectory())))
		} else {
			printVersion(sys)
			printHelp(sys, commandLine)
		}
		return ExitStatusDiagnosticsPresent_OutputsSkipped, nil
	}

	// !!! convert to options with absolute paths is usually done here, but for ease of implementation, it's done in `tsoptions.ParseCommandLine()`
	compilerOptionsFromCommandLine := commandLine.CompilerOptions()
	configForCompilation := commandLine
	var extendedConfigCache collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
	var configTime time.Duration
	if configFileName != "" {
		configStart := sys.Now()
		configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(configFileName, compilerOptionsFromCommandLine, sys, &extendedConfigCache)
		configTime = sys.Now().Sub(configStart)
		if len(errors) != 0 {
			// these are unrecoverable errors--exit to report them as diagnostics
			for _, e := range errors {
				reportDiagnostic(e)
			}
			return ExitStatusDiagnosticsPresent_OutputsGenerated, nil
		}
		configForCompilation = configParseResult
		// Updater to reflect pretty
		reportDiagnostic = createDiagnosticReporter(sys, commandLine.CompilerOptions())
	}

	if compilerOptionsFromCommandLine.ShowConfig.IsTrue() {
		showConfig(sys, configForCompilation.CompilerOptions())
		return ExitStatusSuccess, nil
	}
	if configForCompilation.CompilerOptions().Watch.IsTrue() {
		return ExitStatusSuccess, createWatcher(sys, configForCompilation, reportDiagnostic)
	} else if configForCompilation.CompilerOptions().IsIncremental() {
		return performIncrementalCompilation(
			sys,
			configForCompilation,
			reportDiagnostic,
			&extendedConfigCache,
			configTime,
		), nil
	}
	return performCompilation(
		sys,
		configForCompilation,
		reportDiagnostic,
		&extendedConfigCache,
		configTime,
	), nil
}

func findConfigFile(searchPath string, fileExists func(string) bool, configName string) string {
	result, ok := tspath.ForEachAncestorDirectory(searchPath, func(ancestor string) (string, bool) {
		fullConfigName := tspath.CombinePaths(ancestor, configName)
		if fileExists(fullConfigName) {
			return fullConfigName, true
		}
		return fullConfigName, false
	})
	if !ok {
		return ""
	}
	return result
}

func performIncrementalCompilation(
	sys System,
	config *tsoptions.ParsedCommandLine,
	reportDiagnostic diagnosticReporter,
	extendedConfigCache *collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry],
	configTime time.Duration,
) ExitStatus {
	host := compiler.NewCachedFSCompilerHost(config.CompilerOptions(), sys.GetCurrentDirectory(), sys.FS(), sys.DefaultLibraryPath(), extendedConfigCache)
	oldProgram := incremental.ReadBuildInfoProgram(config, incremental.NewBuildInfoReader(host))
	// todo: cache, statistics, tracing
	parseStart := sys.Now()
	program := compiler.NewProgram(compiler.ProgramOptions{
		Config:           config,
		Host:             host,
		JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
	})
	parseTime := sys.Now().Sub(parseStart)
	incrementalProgram := incremental.NewProgram(program, oldProgram)
	return emitAndReportStatistics(
		sys,
		incrementalProgram,
		incrementalProgram.GetProgram,
		config,
		reportDiagnostic,
		configTime,
		parseTime,
	)
}

func performCompilation(
	sys System,
	config *tsoptions.ParsedCommandLine,
	reportDiagnostic diagnosticReporter,
	extendedConfigCache *collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry],
	configTime time.Duration,
) ExitStatus {
	host := compiler.NewCachedFSCompilerHost(config.CompilerOptions(), sys.GetCurrentDirectory(), sys.FS(), sys.DefaultLibraryPath(), extendedConfigCache)
	// todo: cache, statistics, tracing
	parseStart := sys.Now()
	program := compiler.NewProgram(compiler.ProgramOptions{
		Config:           config,
		Host:             host,
		JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
	})
	parseTime := sys.Now().Sub(parseStart)
	return emitAndReportStatistics(
		sys,
		program,
		func() *compiler.Program {
			return program
		},
		config,
		reportDiagnostic,
		configTime,
		parseTime,
	)
}

func emitAndReportStatistics(
	sys System,
	program compiler.AnyProgram,
	getCoreProgram func() *compiler.Program,
	config *tsoptions.ParsedCommandLine,
	reportDiagnostic diagnosticReporter,
	configTime time.Duration,
	parseTime time.Duration,
) ExitStatus {
	result := emitFilesAndReportErrors(sys, program, reportDiagnostic)
	if result.status != ExitStatusSuccess {
		// compile exited early
		return result.status
	}

	result.configTime = configTime
	result.parseTime = parseTime
	result.totalTime = sys.SinceStart()

	if config.CompilerOptions().Diagnostics.IsTrue() || config.CompilerOptions().ExtendedDiagnostics.IsTrue() {
		var memStats runtime.MemStats
		// GC must be called twice to allow things to settle.
		runtime.GC()
		runtime.GC()
		runtime.ReadMemStats(&memStats)

		reportStatistics(sys, getCoreProgram(), result, &memStats)
	}

	if result.emitResult.EmitSkipped && len(result.diagnostics) > 0 {
		return ExitStatusDiagnosticsPresent_OutputsSkipped
	} else if len(result.diagnostics) > 0 {
		return ExitStatusDiagnosticsPresent_OutputsGenerated
	}
	return ExitStatusSuccess
}

type compileAndEmitResult struct {
	diagnostics []*ast.Diagnostic
	emitResult  *compiler.EmitResult
	status      ExitStatus
	configTime  time.Duration
	parseTime   time.Duration
	bindTime    time.Duration
	checkTime   time.Duration
	totalTime   time.Duration
	emitTime    time.Duration
}

func emitFilesAndReportErrors(
	sys System,
	program compiler.AnyProgram,
	reportDiagnostic diagnosticReporter,
) (result compileAndEmitResult) {
	ctx := context.Background()

	allDiagnostics := compiler.GetDiagnosticsOfAnyProgram(
		ctx,
		program,
		nil,
		false,
		func(name string, start bool, startTime time.Time) time.Time {
			if !start {
				switch name {
				case "bind":
					result.bindTime = sys.Now().Sub(startTime)
				case "check":
					result.checkTime = sys.Now().Sub(startTime)
				}
			}
			return sys.Now()
		},
	)

	emitResult := &compiler.EmitResult{EmitSkipped: true, Diagnostics: []*ast.Diagnostic{}}
	if !program.Options().ListFilesOnly.IsTrue() {
		emitStart := sys.Now()
		emitResult = program.Emit(ctx, compiler.EmitOptions{})
		result.emitTime = sys.Now().Sub(emitStart)
	}
	allDiagnostics = append(allDiagnostics, emitResult.Diagnostics...)

	allDiagnostics = compiler.SortAndDeduplicateDiagnostics(allDiagnostics)
	for _, diagnostic := range allDiagnostics {
		reportDiagnostic(diagnostic)
	}

	if sys.Writer() != nil {
		for _, file := range emitResult.EmittedFiles {
			fmt.Fprint(sys.Writer(), "TSFILE: ", tspath.GetNormalizedAbsolutePath(file, sys.GetCurrentDirectory()))
		}
		listFiles(sys, program)
	}

	createReportErrorSummary(sys, program.Options())(allDiagnostics)
	result.diagnostics = allDiagnostics
	result.emitResult = emitResult
	result.status = ExitStatusSuccess
	return result
}

// func isBuildCommand(args []string) bool {
// 	return len(args) > 0 && args[0] == "build"
// }

func showConfig(sys System, config *core.CompilerOptions) {
	// !!!
	enc := json.NewEncoder(sys.Writer())
	enc.SetIndent("", "    ")
	enc.Encode(config) //nolint:errcheck,errchkjson
}

func listFiles(sys System, program compiler.AnyProgram) {
	options := program.Options()
	// !!! explainFiles
	if options.ListFiles.IsTrue() || options.ListFilesOnly.IsTrue() {
		for _, file := range program.GetSourceFiles() {
			fmt.Fprintf(sys.Writer(), "%s%s", file.FileName(), sys.NewLine())
		}
	}
}
