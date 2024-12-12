package compiler

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type EmitOptions struct {
	TargetSourceFile *ast.SourceFile // Single file to emit. If `nil`, emits all files
	forceDtsEmit     bool
}

type EmitResult struct {
	EmitSkipped  bool
	Diagnostics  []*ast.Diagnostic      // Contains declaration emit diagnostics
	EmittedFiles []string               // Array of files the compiler wrote to disk
	sourceMaps   []*sourceMapEmitResult // Array of sourceMapData if compiler emitted sourcemaps
}

type sourceMapEmitResult struct {
	inputSourceFileNames []string // Input source file (which one can use on program to get the file), 1:1 mapping with the sourceMap.sources list
	sourceMap            *sourcemap.RawSourceMap
}

func (p *Program) Emit(options *EmitOptions) *EmitResult {
	// !!! performance measurement

	host := &emitHost{program: p}

	writerPool := &sync.Pool{
		New: func() any {
			return printer.NewTextWriter(host.Options().NewLine.GetNewLineCharacter())
		},
	}
	wg := core.NewWorkGroup(p.programOptions.SingleThreaded)

	var emitters []*emitter
	sourceFiles := getSourceFilesToEmit(host, options.TargetSourceFile, options.forceDtsEmit)
	for _, sourceFile := range sourceFiles {
		emitter := &emitter{
			host:              host,
			emittedFilesList:  nil,
			sourceMapDataList: nil,
			writer:            nil,
			sourceFile:        sourceFile,
		}
		emitters = append(emitters, emitter)
		wg.Run(func() {
			// take an unused writer
			writer := writerPool.Get().(printer.EmitTextWriter)
			writer.Clear()

			// attach writer and perform emit
			emitter.writer = writer
			emitter.paths = getOutputPathsFor(sourceFile, host, options.forceDtsEmit)
			emitter.emit()
			emitter.writer = nil

			// put the writer back in the pool
			writerPool.Put(writer)
		})
	}

	// wait for emit to complete
	wg.Wait()

	// collect results from emit, preserving input order
	result := &EmitResult{}
	for _, emitter := range emitters {
		if emitter.emitSkipped {
			result.EmitSkipped = true
		}
		result.Diagnostics = append(result.Diagnostics, emitter.emitterDiagnostics.GetDiagnostics()...)
		if emitter.emittedFilesList != nil {
			result.EmittedFiles = append(result.EmittedFiles, emitter.emittedFilesList...)
		}
		if emitter.sourceMapDataList != nil {
			result.sourceMaps = append(result.sourceMaps, emitter.sourceMapDataList...)
		}
	}
	return result
}

type emitOnly byte

const (
	emitAll emitOnly = iota
	emitOnlyJs
	emitOnlyDts
	emitOnlyBuildInfo
)

type emitter struct {
	host               EmitHost
	emitOnly           emitOnly
	emittedFilesList   []string
	emitterDiagnostics DiagnosticsCollection
	emitSkipped        bool
	sourceMapDataList  []*sourceMapEmitResult
	writer             printer.EmitTextWriter
	paths              *outputPaths
	sourceFile         *ast.SourceFile
}

func (e *emitter) emit() {
	// !!! tracing
	e.emitJsFile(e.sourceFile, e.paths.jsFilePath, e.paths.sourceMapFilePath)
	e.emitDeclarationFile(e.sourceFile, e.paths.declarationFilePath, e.paths.declarationMapPath)
	e.emitBuildInfo(e.paths.buildInfoPath)
}

func (e *emitter) emitJsFile(sourceFile *ast.SourceFile, jsFilePath string, sourceMapFilePath string) {
	options := e.host.Options()

	if sourceFile == nil || e.emitOnly != emitAll && e.emitOnly != emitOnlyJs || len(jsFilePath) == 0 {
		return
	}
	if options.NoEmit == core.TSTrue || e.host.IsEmitBlocked(jsFilePath) {
		return
	}

	// !!! mark linked references
	// !!! transform the source files?

	printerOptions := printer.PrinterOptions{
		NewLine: options.NewLine,
		// !!!
	}

	// create a printer to print the nodes
	printer := printer.NewPrinter(printerOptions, printer.PrintHandlers{
		// !!!
	})

	e.printSourceFile(jsFilePath, sourceMapFilePath, sourceFile, printer)

	if e.emittedFilesList != nil {
		e.emittedFilesList = append(e.emittedFilesList, jsFilePath)
		if sourceMapFilePath != "" {
			e.emittedFilesList = append(e.emittedFilesList, sourceMapFilePath)
		}
	}
}

func (e *emitter) emitDeclarationFile(sourceFile *ast.SourceFile, declarationFilePath string, declarationMapPath string) {
	// !!!
}

func (e *emitter) emitBuildInfo(buildInfoPath string) {
	// !!!
}

func (e *emitter) printSourceFile(jsFilePath string, sourceMapFilePath string, sourceFile *ast.SourceFile, printer *printer.Printer) bool {
	// !!! sourceMapGenerator
	// !!! bundles not implemented, may be deprecated
	sourceFiles := []*ast.SourceFile{sourceFile}

	printer.Write(sourceFile.AsNode(), sourceFile, e.writer /*, sourceMapGenerator*/)

	// !!! add sourceMapGenerator to sourceMapDataList
	// !!! append sourceMappingURL to output
	// !!! write the source map
	e.writer.WriteLine()

	// Write the output file
	text := e.writer.String()
	data := &WriteFileData{} // !!!
	err := e.host.WriteFile(jsFilePath, text, e.host.Options().EmitBOM == core.TSTrue, sourceFiles, data)
	if err != nil {
		e.emitterDiagnostics.add(ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, jsFilePath, err.Error()))
	}

	// Reset state
	e.writer.Clear()
	return !data.SkippedDtsWrite
}

func getOutputExtension(fileName string, jsx core.JsxEmit) string {
	switch {
	case tspath.FileExtensionIs(fileName, tspath.ExtensionJson):
		return tspath.ExtensionJson
	case jsx == core.JsxEmitPreserve && tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionJsx, tspath.ExtensionTsx}):
		return tspath.ExtensionJsx
	case tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionMts, tspath.ExtensionMjs}):
		return tspath.ExtensionMjs
	case tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionCts, tspath.ExtensionCjs}):
		return tspath.ExtensionCjs
	default:
		return tspath.ExtensionJs
	}
}

func getSourceFilePathInNewDir(fileName string, newDirPath string, currentDirectory string, commonSourceDirectory string, useCaseSensitiveFileNames bool) string {
	sourceFilePath := tspath.GetNormalizedAbsolutePath(fileName, currentDirectory)
	commonSourceDirectory = tspath.EnsureTrailingDirectorySeparator(commonSourceDirectory)
	isSourceFileInCommonSourceDirectory := tspath.ContainsPath(commonSourceDirectory, sourceFilePath, tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: useCaseSensitiveFileNames,
		CurrentDirectory:          currentDirectory,
	})
	if isSourceFileInCommonSourceDirectory {
		sourceFilePath = sourceFilePath[len(commonSourceDirectory):]
	}
	return tspath.CombinePaths(newDirPath, sourceFilePath)
}

func getOwnEmitOutputFilePath(fileName string, host EmitHost, extension string) string {
	compilerOptions := host.Options()
	var emitOutputFilePathWithoutExtension string
	if len(compilerOptions.OutDir) > 0 {
		currentDirectory := host.GetCurrentDirectory()
		emitOutputFilePathWithoutExtension = tspath.RemoveFileExtension(getSourceFilePathInNewDir(
			fileName,
			compilerOptions.OutDir,
			currentDirectory,
			host.CommonSourceDirectory(),
			host.UseCaseSensitiveFileNames(),
		))
	} else {
		emitOutputFilePathWithoutExtension = tspath.RemoveFileExtension(fileName)
	}
	return emitOutputFilePathWithoutExtension + extension
}

func getSourceMapFilePath(jsFilePath string, options *core.CompilerOptions) string {
	// !!!
	return ""
}

func getDeclarationEmitOutputFilePath(file string, host EmitHost) string {
	// !!!
	return ""
}

type outputPaths struct {
	jsFilePath          string
	sourceMapFilePath   string
	declarationFilePath string
	declarationMapPath  string
	buildInfoPath       string
}

func getOutputPathsFor(sourceFile *ast.SourceFile, host EmitHost, forceDtsEmit bool) *outputPaths {
	options := host.Options()
	// !!! bundle not implemented, may be deprecated
	ownOutputFilePath := getOwnEmitOutputFilePath(sourceFile.FileName(), host, getOutputExtension(sourceFile.FileName(), options.Jsx))
	isJsonFile := isJsonSourceFile(sourceFile)
	// If json file emits to the same location skip writing it, if emitDeclarationOnly skip writing it
	isJsonEmittedToSameLocation := isJsonFile &&
		tspath.ComparePaths(sourceFile.FileName(), ownOutputFilePath, tspath.ComparePathsOptions{
			CurrentDirectory:          host.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: host.UseCaseSensitiveFileNames(),
		}) == 0
	paths := &outputPaths{}
	if options.EmitDeclarationOnly != core.TSTrue && !isJsonEmittedToSameLocation {
		paths.jsFilePath = ownOutputFilePath
		if !isJsonSourceFile(sourceFile) {
			paths.sourceMapFilePath = getSourceMapFilePath(paths.jsFilePath, options)
		}
	}
	if forceDtsEmit || options.GetEmitDeclarations() && !isJsonFile {
		paths.declarationFilePath = getDeclarationEmitOutputFilePath(sourceFile.FileName(), host)
		if options.GetAreDeclarationMapsEnabled() {
			paths.declarationMapPath = paths.declarationFilePath + ".map"
		}
	}
	return paths
}

func forEachEmittedFile(host EmitHost, action func(emitFileNames *outputPaths, sourceFile *ast.SourceFile) bool, sourceFiles []*ast.SourceFile, options *EmitOptions) bool {
	// !!! outFile not yet implemented, may be deprecated
	for _, sourceFile := range sourceFiles {
		if action(getOutputPathsFor(sourceFile, host, options.forceDtsEmit), sourceFile) {
			return true
		}
	}
	return false
}

func sourceFileMayBeEmitted(sourceFile *ast.SourceFile, host EmitHost, forceDtsEmit bool) bool {
	// !!! Js files are emitted only if option is enabled

	// Declaration files are not emitted
	if sourceFile.IsDeclarationFile {
		return false
	}

	// !!! Source file from node_modules are not emitted

	// forcing dts emit => file needs to be emitted
	if forceDtsEmit {
		return true
	}

	// !!! Source files from referenced projects are not emitted

	// Any non json file should be emitted
	if !isJsonSourceFile(sourceFile) {
		return true
	}

	// !!! Should JSON input files be emitted
	return false
}

func getSourceFilesToEmit(host EmitHost, targetSourceFile *ast.SourceFile, forceDtsEmit bool) []*ast.SourceFile {
	// !!! outFile not yet implemented, may be deprecated
	var sourceFiles []*ast.SourceFile
	if targetSourceFile != nil {
		sourceFiles = []*ast.SourceFile{targetSourceFile}
	} else {
		sourceFiles = host.SourceFiles()
	}
	return core.Filter(sourceFiles, func(sourceFile *ast.SourceFile) bool {
		return sourceFileMayBeEmitted(sourceFile, host, forceDtsEmit)
	})
}
