package execute

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	ts "github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	dw "github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func CommandLine(sys System, cb cbType, commandLineArgs []string) ExitStatus {
	_, exitStatus := TestCommandLine(sys, cb, commandLineArgs)
	return exitStatus
}

func TestCommandLine(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, ExitStatus) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys.FS(), sys.Host().GetCurrentDirectory())
	return parsedCommandLine, executeCommandLineWorker(sys, cb, parsedCommandLine)
}

func executeCommandLineWorker(sys System, cb cbType, commandLine *tsoptions.ParsedCommandLine) ExitStatus {
	reportDiagnostic := CreateDiagnosticReporter(sys, false)
	configFileName := ""
	// if commandLine.Options().Locale != nil

	if len(commandLine.Errors) > 0 {
		for _, e := range commandLine.Errors {
			reportDiagnostic(e)
		}
		return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
	}

	// if commandLine.Options().Init != nil
	// if commandLine.Options().Version != nil
	// if commandLine.Options().Help != nil || commandLine.Options().All != nil
	// if commandLine.Options().Watch != nil  && commandLine.Options().ListFilesOnly

	if commandLine.CompilerOptions().Project != "" {
		if len(commandLine.FileNames()) != 0 {
			reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Option_project_cannot_be_mixed_with_source_files_on_a_command_line))
			return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
		}

		fileOrDirectory := tspath.NormalizePath(commandLine.CompilerOptions().Project)
		if fileOrDirectory != "" || sys.FS().DirectoryExists(fileOrDirectory) {
			configFileName = tspath.CombinePaths(fileOrDirectory, "tsconfig.json")
			if !sys.FS().FileExists(configFileName) {
				reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Cannot_find_a_tsconfig_json_file_at_the_current_directory_Colon_0, configFileName))
				return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
			}
		} else {
			configFileName = fileOrDirectory
			if !sys.FS().FileExists(configFileName) {
				reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.The_specified_path_does_not_exist_Colon_0, fileOrDirectory))
				return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
			}
		}
	} else if len(commandLine.FileNames()) == 0 {
		searchPath := tspath.NormalizePath(sys.Host().GetCurrentDirectory())
		configFileName = findConfigFile(searchPath, sys.FS().FileExists, "tsconfig.json")
	}

	if configFileName != "" && len(commandLine.FileNames()) > 0 {
		if commandLine.CompilerOptions().ShowConfig {
			reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Cannot_find_a_tsconfig_json_file_at_the_current_directory_Colon_0, tspath.NormalizePath(sys.Host().GetCurrentDirectory())))
		} else {
			// print version
			// print help
		}
		return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
	}

	// currentDirectory := sys.Host().GetCurrentDirectory()
	compilerOptionsFromCommandLine := commandLine.CompilerOptions() // todo: convert to options with absolute paths

	if configFileName != "" {
		// extendedConfigCache := map[string]any{}
		configParseResult := parseConfigFileWithSystem(
			configFileName,
			compilerOptionsFromCommandLine,
			sys,
		)
		// if commandLineOptions.ShowConfig
		if isWatchSet(configParseResult.CompilerOptions()) {
			// todo watch
			return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
		} else if isIncrementalCompilation(configParseResult.CompilerOptions()) {
			// todo performIncrementalCompilation
		} else {
			performCompilation(
				sys,
				cb,
				reportDiagnostic,
				configParseResult,
			)
		}
	} else {
		if compilerOptionsFromCommandLine.ShowConfig {
			// write show config
			return sys.Exit(ExitStatusSuccess)
		}
		// todo update reportDiagnostic
		if isWatchSet(compilerOptionsFromCommandLine) {
			// todo watch
			return sys.Exit(ExitStatusDiagnosticsPresent_OutputsSkipped)
		} else if isIncrementalCompilation(compilerOptionsFromCommandLine) {
			// todo incremental
		} else {
			commandLine.SetCompilerOptions(compilerOptionsFromCommandLine)
			performCompilation(
				sys,
				cb,
				reportDiagnostic,
				commandLine,
			)
		}
	}

	return sys.Exit(ExitStatusSuccess)
}

func parseConfigFileWithSystem(
	configFileName string,
	optionsToExtend *core.CompilerOptions,
	sys System,
	// reportDiagnostic DiagnosticReporter,
	// extendedConfigCache
) *tsoptions.ParsedCommandLine {
	// todo: on unrecoverable diagnostic-- needs dignosticreporter

	return getParsedCommandLineOfConfigFile(configFileName, optionsToExtend, sys.Host())
}

func getParsedCommandLineOfConfigFile(configFileName string, options *core.CompilerOptions, host ts.CompilerHost) *tsoptions.ParsedCommandLine {
	// todo implement when tsconfigParsing
	errors := []*ast.Diagnostic{}
	configFileText, errors := tsoptions.TryReadFile(configFileName, host.FS().ReadFile, errors)
	if len(errors) > 0 {
		// on unrecoverable diagnostic (errors)
		return nil
	}

	tsConfigSourceFile := parser.ParseJSONText(configFileName, configFileText)
	cwd := host.GetCurrentDirectory()

	tsConfigSourceFile.SetPath(tspath.ToPath(configFileName, cwd, host.FS().UseCaseSensitiveFileNames()))
	// todo tsconfigParsing
	// result.resolvedPath = tspath.NormalizePath(result.resolvedPath)
	// result.originalFileName = result.FileName()
	// return tsoptions.ParseJsonSourceFileConfigFileContent(
	// 	tsConfigSourceFile,
	// 	host,
	// 	tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(configFileName), cwd),
	// 	options,
	// 	tspath.GetNormalizedAbsolutePath(configFileName, cwd),
	// )
	return nil
}

func performCompilation(sys System, cb cbType, reportDiagnostic DiagnosticReporter, config *tsoptions.ParsedCommandLine) ExitStatus {
	program := ts.NewProgramFromParsedCommandLine(config, sys.Host())
	options := program.Options()
	allDiagnostics := program.ConfigParsingDiagnostics()

	// todo: early exit logic and append diagnostics
	diagnostics := program.GetSyntacticDiagnostics(nil)
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, program.GetOptionsDiagnostics()...)
		if !options.ListFilesOnly {
			// program.GetBindDiagnostics(nil)
			diagnostics = append(diagnostics, program.GetGlobalDiagnostics()...)
		}
	}
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, program.GetSemanticDiagnostics(nil)...)
	}
	// todo declaration diagnostics
	// if len(diagnostics) == 0 && options.NoEmit == core.TSTrue && getEmitDeclarations(options) {
	// 	addRange(allDiagnostics, program.getDeclarationDiagnostics(/*sourceFile*/ undefined, cancellationToken));
	// }

	emitResult := &ts.EmitResult{EmitSkipped: true, Diagnostics: []*ast.Diagnostic{}}
	if !options.ListFilesOnly {
		// todo emit not fully implemented
		// emitResult = program.Emit(&ts.EmitOptions{})
	}
	diagnostics = append(diagnostics, emitResult.Diagnostics...)

	allDiagnostics = append(allDiagnostics, diagnostics...)
	if allDiagnostics != nil {
		allDiagnostics = ts.SortAndDeduplicateDiagnostics(allDiagnostics)
		for _, diagnostic := range allDiagnostics {
			reportDiagnostic(diagnostic)
		}
		createReportErrorSummary(sys, config.CompilerOptions())(allDiagnostics)
	}
	reportStatistics(sys, program)
	if cb != nil {
		cb(program)
	}
	return sys.Exit(getExitStatus(emitResult, allDiagnostics))
}

func getExitStatus(emitResult *ts.EmitResult, diagnostics []*ast.Diagnostic) ExitStatus {
	if emitResult.EmitSkipped && diagnostics != nil && len(diagnostics) > 0 {
		return ExitStatusDiagnosticsPresent_OutputsSkipped
	} else if len(diagnostics) > 0 {
		return ExitStatusDiagnosticsPresent_OutputsGenerated
	}
	return ExitStatusSuccess
}

func isBuildCommand(args []string) bool {
	return len(args) > 0 && args[0] == "build"
}

func isWatchSet(options *core.CompilerOptions) bool {
	return options.Watch
}

func isIncrementalCompilation(options *core.CompilerOptions) bool {
	return options.Incremental
}

type (
	DiagnosticReporter = func(diagnostic *ast.Diagnostic)
	cbType             = func(p any) any
)

func CreateDiagnosticReporter(sys System, pretty bool) DiagnosticReporter {
	if !pretty {
		return func(diagnostic *ast.Diagnostic) {
			dw.WriteFormatDiagnostic(sys, diagnostic, sys.GetFormatOpts())
			sys.EndWrite()
		}
	}

	diagArr := [1]*ast.Diagnostic{}
	return func(diagnostic *ast.Diagnostic) {
		diagArr[0] = diagnostic
		dw.FormatDiagnosticsWithColorAndContext(sys, diagArr[:], sys.GetFormatOpts())
		sys.EndWrite()
		diagArr[0] = nil
	}
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

func createReportErrorSummary(sys System, options *core.CompilerOptions) func(diagnostics []*ast.Diagnostic) {
	if shouldBePretty(sys, options) {
		formatOpts := sys.GetFormatOpts()
		return func(diagnostics []*ast.Diagnostic) {
			dw.WriteErrorSummaryText(sys, diagnostics, formatOpts)
			sys.EndWrite()
		}
	}
	return nil
}

func shouldBePretty(sys System, options *core.CompilerOptions) bool {
	if options == nil || !options.Pretty {
		// todo: return defaultIsPretty(sys);
		return true
	}
	return options.Pretty
}

func reportStatistics(sys System, program *ts.Program) {
	// todo
	stats := []statistic{
		newStatistic("Files", len(program.SourceFiles())),
		// newStatistic("Identifiers", program.IdentifierCount()),
		// newStatistic("Symbols", program.getSymbolCount()),
		newStatistic("Types", program.TypeCount()),
		// newStatistic("Instantiations", program.getInstantiationCount()),
	}

	for _, stat := range stats {
		fmt.Fprintf(sys, "%s:"+strings.Repeat(" ", 20-len(stat.name))+"%v\n", stat.name, stat.value)
	}
}

type statistic struct {
	name  string
	value int
}

func newStatistic(name string, count int) statistic {
	return statistic{name, count}
}
