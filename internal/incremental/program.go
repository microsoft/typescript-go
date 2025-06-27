package incremental

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Program struct {
	snapshot *snapshot
	state    *programState
	program  *compiler.Program
}

var _ compiler.AnyProgram = (*Program)(nil)

func NewProgram(program *compiler.Program, oldProgram *Program) *Program {
	return &Program{
		state:   newProgramState(program, oldProgram),
		program: program,
	}
}

func (p *Program) commitChange(ctx context.Context, change snapshotChange) {
	if change == nil || ctx.Err() != nil {
		return
	}
	change.commit(p.snapshot)
	p.snapshot.buildInfoEmitPending = true
}

func (p *Program) panicIfNoProgram(method string) {
	if p.program == nil {
		panic(fmt.Sprintf("%s should not be called without program", method))
	}
}

func (p *Program) GetProgram() *compiler.Program {
	p.panicIfNoProgram("GetProgram")
	return p.program
}

// Options implements compiler.AnyProgram interface.
func (p *Program) Options() *core.CompilerOptions {
	return p.snapshot.options
}

// GetSourceFiles implements compiler.AnyProgram interface.
func (p *Program) GetSourceFiles() []*ast.SourceFile {
	p.panicIfNoProgram("GetSourceFiles")
	return p.program.GetSourceFiles()
}

// GetConfigFileParsingDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	p.panicIfNoProgram("GetConfigFileParsingDiagnostics")
	return p.program.GetConfigFileParsingDiagnostics()
}

// GetSyntacticDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetSyntacticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetSyntacticDiagnostics")
	return p.program.GetSyntacticDiagnostics(ctx, file)
}

// GetBindDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetBindDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetBindDiagnostics")
	return p.program.GetBindDiagnostics(ctx, file)
}

// GetOptionsDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic {
	p.panicIfNoProgram("GetOptionsDiagnostics")
	return p.program.GetOptionsDiagnostics(ctx)
}

// GetGlobalDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic {
	p.panicIfNoProgram("GetGlobalDiagnostics")
	return p.program.GetGlobalDiagnostics(ctx)
}

// GetSemanticDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetSemanticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	if p.snapshot.options.NoCheck.IsTrue() {
		return nil
	}

	p.panicIfNoProgram("GetSemanticDiagnostics")
	if file != nil {
		// !!! sheetal Ensure all affected files are collected

		// use cached if present otherwise cache and get from program

		diagnostics, change := p.getSemanticDiagnosticsOfFile(ctx, file)

		// This should not result in any changes
		p.commitChange(ctx, change)
		return diagnostics
	}

	// Ensure all the diagnsotics are cached
	p.getSemanticDiagnosticsOfAffectedFiles(ctx)
	if ctx.Err() != nil {
		return nil
	}

	// Return result from cache
	var diagnostics []*ast.Diagnostic
	for _, file := range p.program.GetSourceFiles() {
		diagnosticsOfFile, change := p.getSemanticDiagnosticsOfFile(ctx, file)
		if change != nil {
			panic("After handling all the affected files, there shouldnt be more changes")
		}
		diagnostics = append(diagnostics, diagnosticsOfFile...)
	}
	return diagnostics
}

// GetDeclarationDiagnostics implements compiler.AnyProgram interface.
func (p *Program) GetDeclarationDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetDeclarationDiagnostics")
	return p.state.getDeclarationDiagnostics(ctx, p.program, file)
}

// GetModeForUsageLocation implements compiler.AnyProgram interface.
func (p *Program) Emit(ctx context.Context, options compiler.EmitOptions) *compiler.EmitResult {
	p.panicIfNoProgram("Emit")
	return p.state.emit(ctx, p.program, options)
}

func (p *Program) computeDtsSignature(ctx context.Context, file *ast.SourceFile) string {
	var signature string
	p.program.Emit(ctx, compiler.EmitOptions{
		TargetSourceFile: file,
		EmitOnly:         compiler.EmitOnlyForcedDts,
		WriteFile: func(fileName string, text string, writeByteOrderMark bool, data *compiler.WriteFileData) error {
			if !tspath.IsDeclarationFileName(fileName) {
				panic("File extension for signature expected to be dts, got : " + fileName)
			}
			signature = computeSignatureWithDiagnostics(file, text, data)
			return nil
		},
	})
	return signature
}

func (p *Program) updateShapeSignature(change *affectedFilesChange, file *ast.SourceFile, useFileVersionAsSignature bool) bool {
	// If we have cached the result for this file, that means hence forth we should assume file shape is uptodate
	if _, ok := change.updatedSignatures.Load(file.Path()); ok {
		return false
	}

	info := p.snapshot.fileInfos[file.Path()]
	prevSignature := info.signature
	var latestSignature string
	if !file.IsDeclarationFile && !useFileVersionAsSignature {
		latestSignature = p.computeDtsSignature(change.ctx, file)
	}
	// Default is to use file version as signature
	if latestSignature == "" {
		latestSignature = info.version
	}
	change.updatedSignatures.Store(file.Path(), latestSignature)
	return latestSignature != prevSignature
}

// This function collects all the affected files to be processed.
func (p *Program) collectAllAffectedFiles(ctx context.Context) {
	if p.snapshot.changedFilesSet.Len() == 0 {
		return
	}

	wg := core.NewWorkGroup(p.program.SingleThreaded())
	change := affectedFilesChange{ctx: ctx, program: p}
	var result collections.SyncSet[*ast.SourceFile]
	for file := range p.snapshot.changedFilesSet.Keys() {
		wg.Queue(func() {
			for _, affectedFile := range p.getFilesAffectedBy(&change, file) {
				result.Add(affectedFile)
			}
		})
	}
	wg.RunAndWait()

	if ctx.Err() != nil {
		return
	}

	// For all the affected files, get all the files that would need to change their dts or js files,
	// update their diagnostics
	wg = core.NewWorkGroup(p.program.SingleThreaded())
	emitKind := GetFileEmitKind(p.snapshot.options)
	result.Range(func(file *ast.SourceFile) bool {
		// remove the cached semantic diagnostics and handle dts emit and js emit if needed
		dtsMayChange := change.getDtsMayChange(file.Path(), emitKind)
		wg.Queue(func() {
			p.handleDtsMayChangeOfAffectedFile(dtsMayChange, file)
		})
		return true
	})
	wg.RunAndWait()
	if ctx.Err() != nil {
		return
	}
	change.commit(p.snapshot)
	return
}

func (p *Program) getFilesAffectedBy(change *affectedFilesChange, path tspath.Path) []*ast.SourceFile {
	file := p.program.GetSourceFileByPath(path)
	if file == nil {
		return nil
	}

	if !p.updateShapeSignature(change, file, p.snapshot.useFileVersionAsSignature) {
		return []*ast.SourceFile{file}
	}

	if !p.snapshot.tracksReferences() {
		change.hasAllFilesExcludingDefaultLibraryFile.Store(true)
		return p.snapshot.getAllFilesExcludingDefaultLibraryFile(p.program, file)
	}

	if info := p.snapshot.fileInfos[file.Path()]; info.affectsGlobalScope {
		change.hasAllFilesExcludingDefaultLibraryFile.Store(true)
		p.snapshot.getAllFilesExcludingDefaultLibraryFile(p.program, file)
	}

	if p.snapshot.options.IsolatedModules.IsTrue() {
		return []*ast.SourceFile{file}
	}

	// Now we need to if each file in the referencedBy list has a shape change as well.
	// Because if so, its own referencedBy files need to be saved as well to make the
	// emitting result consistent with files on disk.
	seenFileNamesMap := p.forEachFileReferencedBy(
		file,
		func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool) {
			// If the current file is not nil and has a shape change, we need to queue it for processing
			if currentFile != nil && p.updateShapeSignature(change, currentFile, p.snapshot.useFileVersionAsSignature) {
				return true, false
			}
			return false, false
		},
	)
	// Return array of values that needs emit
	return core.Filter(slices.Collect(maps.Values(seenFileNamesMap)), func(file *ast.SourceFile) bool {
		return file != nil
	})
}

// Gets the files referenced by the the file path
func (p *Program) getReferencedByPaths(file tspath.Path) map[tspath.Path]struct{} {
	keys, ok := p.snapshot.referencedMap.GetKeys(file)
	if !ok {
		return nil
	}
	return keys.Keys()
}

func (p *Program) forEachFileReferencedBy(file *ast.SourceFile, fn func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool)) map[tspath.Path]*ast.SourceFile {
	// Now we need to if each file in the referencedBy list has a shape change as well.
	// Because if so, its own referencedBy files need to be saved as well to make the
	// emitting result consistent with files on disk.
	seenFileNamesMap := map[tspath.Path]*ast.SourceFile{}
	// Start with the paths this file was referenced by
	seenFileNamesMap[file.Path()] = file
	references := p.getReferencedByPaths(file.Path())
	queue := slices.Collect(maps.Keys(references))
	for len(queue) > 0 {
		currentPath := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if _, ok := seenFileNamesMap[currentPath]; !ok {
			currentFile := p.program.GetSourceFileByPath(currentPath)
			seenFileNamesMap[currentPath] = currentFile
			queueForFile, fastReturn := fn(currentFile, currentPath)
			if fastReturn {
				return seenFileNamesMap
			}
			if queueForFile {
				for ref := range p.getReferencedByPaths(currentFile.Path()) {
					queue = append(queue, ref)
				}
			}
		}
	}
	return seenFileNamesMap
}

// Handles semantic diagnostics and dts emit for affectedFile and files, that are referencing modules that export entities from affected file
// This is because even though js emit doesnt change, dts emit / type used can change resulting in need for dts emit and js change
func (p *Program) handleDtsMayChangeOfAffectedFile(dtsMayChange *dtsMayChange, affectedFile *ast.SourceFile) {
	dtsMayChange.change.removeSemanticDiagnosticsOf(affectedFile.Path())

	// If affected files is everything except default library, then nothing more to do
	if dtsMayChange.change.hasAllFilesExcludingDefaultLibraryFile.Load() {
		dtsMayChange.change.removeDiagnosticsOfLibraryFiles()
		// When a change affects the global scope, all files are considered to be affected without updating their signature
		// That means when affected file is handled, its signature can be out of date
		// To avoid this, ensure that we update the signature for any affected file in this scenario.
		p.updateShapeSignature(dtsMayChange.change, affectedFile, p.snapshot.useFileVersionAsSignature)
		return
	}

	if p.snapshot.options.AssumeChangesOnlyAffectDirectDependencies.IsTrue() {
		return
	}

	// Iterate on referencing modules that export entities from affected file and delete diagnostics and add pending emit
	// If there was change in signature (dts output) for the changed file,
	// then only we need to handle pending file emit
	if !p.snapshot.tracksReferences() ||
		!p.snapshot.changedFilesSet.Has(affectedFile.Path()) ||
		!dtsMayChange.change.isChangedSignature(affectedFile.Path()) {
		return
	}

	// Since isolated modules dont change js files, files affected by change in signature is itself
	// But we need to cleanup semantic diagnostics and queue dts emit for affected files
	if p.snapshot.options.IsolatedModules.IsTrue() {
		p.forEachFileReferencedBy(
			affectedFile,
			func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool) {
				if p.handleDtsMayChangeOfGlobalScope(dtsMayChange, currentPath /*invalidateJsFiles*/, false) {
					return false, true
				}
				p.handleDtsMayChangeOf(dtsMayChange, currentPath /*invalidateJsFiles*/, false)
				if dtsMayChange.change.isChangedSignature(currentPath) {
					return true, false
				}
				return false, false
			},
		)
	}

	// !!! sheetal - do i need to pull it off onto change so we arent visiting again and again
	seenFileAndExportsOfFile := collections.Set[tspath.Path]{}
	invalidateJsFiles := false
	var typeChecker *checker.Checker
	var done func()
	// If exported const enum, we need to ensure that js files are emitted as well since the const enum value changed
	if affectedFile.Symbol != nil {
		for _, exported := range affectedFile.Symbol.Exports {
			if exported.Flags&ast.SymbolFlagsConstEnum != 0 {
				invalidateJsFiles = true
				break
			}
			if typeChecker == nil {
				typeChecker, done = p.program.GetTypeCheckerForFile(dtsMayChange.change.ctx, affectedFile)
			}
			aliased := checker.SkipAlias(exported, typeChecker)
			if aliased == exported {
				continue
			}
			if (aliased.Flags & ast.SymbolFlagsConstEnum) != 0 {
				if slices.ContainsFunc(aliased.Declarations, func(d *ast.Node) bool {
					return ast.GetSourceFileOfNode(d) == affectedFile
				}) {
					invalidateJsFiles = true
					break
				}
			}
		}
	}
	if done != nil {
		done()
	}

	// Go through files that reference affected file and handle dts emit and semantic diagnostics for them and their references
	if keys, ok := p.snapshot.referencedMap.GetKeys(affectedFile.Path()); ok {
		for exportedFromPath := range keys.Keys() {
			if p.handleDtsMayChangeOfGlobalScope(dtsMayChange, exportedFromPath, invalidateJsFiles) {
				return
			}
			if references, ok := p.snapshot.referencedMap.GetKeys(exportedFromPath); ok {
				for filePath := range references.Keys() {
					if p.handleDtsMayChangeOfFileAndExportsOfFile(dtsMayChange, filePath, invalidateJsFiles, &seenFileAndExportsOfFile) {
						return
					}
				}
			}
		}
	}
}

func (p *Program) handleDtsMayChangeOfFileAndExportsOfFile(dtsMayChange *dtsMayChange, filePath tspath.Path, invalidateJsFiles bool, seenFileAndExportsOfFile *collections.Set[tspath.Path]) bool {
	if seenFileAndExportsOfFile.AddIfAbsent(filePath) == false {
		return false
	}
	if p.handleDtsMayChangeOfGlobalScope(dtsMayChange, filePath, invalidateJsFiles) {
		return true
	}
	p.handleDtsMayChangeOf(dtsMayChange, filePath, invalidateJsFiles)

	// Remove the diagnostics of files that import this file and handle all its exports too
	if keys, ok := p.snapshot.referencedMap.GetKeys(filePath); ok {
		for referencingFilePath := range keys.Keys() {
			if p.handleDtsMayChangeOfFileAndExportsOfFile(dtsMayChange, referencingFilePath, invalidateJsFiles, seenFileAndExportsOfFile) {
				return true
			}
		}
	}
	return false
}

func (p *Program) handleDtsMayChangeOfGlobalScope(dtsMayChange *dtsMayChange, filePath tspath.Path, invalidateJsFiles bool) bool {
	if info, ok := p.snapshot.fileInfos[filePath]; !ok || !info.affectsGlobalScope {
		return false
	}
	// Every file needs to be handled
	for _, file := range p.snapshot.getAllFilesExcludingDefaultLibraryFile(p.program, nil) {
		p.handleDtsMayChangeOf(dtsMayChange, file.Path(), invalidateJsFiles)
	}
	dtsMayChange.change.removeDiagnosticsOfLibraryFiles()
	return true
}

// Handle the dts may change, so they need to be added to pending emit if dts emit is enabled,
// Also we need to make sure signature is updated for these files
func (p *Program) handleDtsMayChangeOf(dtsMayChange *dtsMayChange, path tspath.Path, invalidateJsFiles bool) {
	dtsMayChange.change.removeSemanticDiagnosticsOf(path)
	if p.snapshot.changedFilesSet.Has(path) {
		return
	}
	file := p.program.GetSourceFileByPath(path)
	if file == nil {
		return
	}
	// Even though the js emit doesnt change and we are already handling dts emit and semantic diagnostics
	// we need to update the signature to reflect correctness of the signature(which is output d.ts emit) of this file
	// This ensures that we dont later during incremental builds considering wrong signature.
	// Eg where this also is needed to ensure that .tsbuildinfo generated by incremental build should be same as if it was first fresh build
	// But we avoid expensive full shape computation, as using file version as shape is enough for correctness.
	p.updateShapeSignature(dtsMayChange.change, file, true)
	// If not dts emit, nothing more to do
	if invalidateJsFiles {
		dtsMayChange.addFileToAffectedFilesPendingEmit(path, GetFileEmitKind(p.snapshot.options))
	} else if p.snapshot.options.GetEmitDeclarations() {
		dtsMayChange.addFileToAffectedFilesPendingEmit(path, core.IfElse(p.snapshot.options.DeclarationMap.IsTrue(), fileEmitKindAllDts, fileEmitKindDts))
	}
}

// Gets the semantic diagnostics either from cache if present, or otherwise from program and caches it
// Note that it is assumed that when asked about checker diagnostics, the file has been taken out of affected files/changed file set
func (p *Program) getSemanticDiagnosticsOfFile(ctx context.Context, file *ast.SourceFile) ([]*ast.Diagnostic, snapshotChange) {
	// Report the check diagnostics from the cache if we already have those diagnostics present
	if cachedDiagnostics, ok := p.snapshot.semanticDiagnosticsPerFile[file.Path()]; ok {
		return compiler.FilterNoEmitSemanticDiagnostics(cachedDiagnostics.getDiagnostics(p.program, file), p.snapshot.options), nil
	}

	// Diagnostics werent cached, get them from program, and cache the result
	change := &semanticDiagnosticChange{semanticDiagnosticsPerFile: make(map[tspath.Path][]*ast.Diagnostic)}
	diagnostics := p.program.GetSemanticDiagnostics(ctx, file)
	change.semanticDiagnosticsPerFile[file.Path()] = diagnostics
	return compiler.FilterNoEmitSemanticDiagnostics(diagnostics, p.snapshot.options), change
}

// Handle affected files and cache the semantic diagnostics
func (p *Program) getSemanticDiagnosticsOfAffectedFiles(ctx context.Context) {
	if len(p.snapshot.semanticDiagnosticsPerFile) == len(p.program.GetSourceFiles()) {
		// If we have all the files,
		return
	}

	// Get all affected files
	p.collectAllAffectedFiles(ctx)
	if ctx.Err() != nil {
		return
	}

	var affectedFiles []*ast.SourceFile
	for _, file := range p.program.GetSourceFiles() {
		if _, ok := p.snapshot.semanticDiagnosticsPerFile[file.Path()]; !ok {
			affectedFiles = append(affectedFiles, file)
		}
	}

	// Get their diagnostics and cache them

	// commit changes if no err
	if ctx.Err() != nil {
		return
	}

	if p.snapshot.checkPending && !p.snapshot.options.NoCheck.IsTrue() {
		p.snapshot.checkPending = false
		p.snapshot.buildInfoEmitPending = true
	}

	// for {
	// 	affected := p.getNextAffectedFile(ctx)
	// 	if affected == nil {
	// 		if p.snapshot.checkPending && !p.snapshot.options.NoCheck.IsTrue() {
	// 			p.checkPending = false
	// 			p.buildInfoEmitPending = true
	// 		}
	// 		return nil, true
	// 	}
	// 	// Get diagnostics for the affected file if its not ignored
	// 	result := p.getSemanticDiagnosticsOfFile(ctx, program, affected)
	// 	p.seenAffectedFiles.Add(affected.Path())
	// 	p.affectedFilesIndex++
	// 	p.buildInfoEmitPending = true
	// 	if result == nil {
	// 		continue
	// 	}
	// 	return result, false
	// }
}

func (p *Program) emitBuildInfo(ctx context.Context, program *compiler.Program, options compiler.EmitOptions) *compiler.EmitResult {
	buildInfoFileName := outputpaths.GetBuildInfoFileName(p.snapshot.options, tspath.ComparePathsOptions{
		CurrentDirectory:          program.GetCurrentDirectory(),
		UseCaseSensitiveFileNames: program.UseCaseSensitiveFileNames(),
	})
	if buildInfoFileName == "" {
		return nil
	}

	hasErrors := p.ensureHasErrorsForState(ctx, program)
	if !p.snapshot.buildInfoEmitPending && p.snapshot.hasErrors == hasErrors {
		return nil
	}
	p.snapshot.hasErrors = hasErrors
	p.snapshot.buildInfoEmitPending = true
	if ctx.Err() != nil {
		return &compiler.EmitResult{
			EmitSkipped: true,
			Diagnostics: []*ast.Diagnostic{
				ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, buildInfoFileName, ctx.Err()),
			},
		}
	}
	buildInfo := snapshotToBuildInfo(p.snapshot, program, buildInfoFileName)
	text, err := json.Marshal(buildInfo)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal build info: %v", err))
	}
	if options.WriteFile != nil {
		err = options.WriteFile(buildInfoFileName, string(text), false, &compiler.WriteFileData{
			BuildInfo: &buildInfo,
		})
	} else {
		err = program.Host().FS().WriteFile(buildInfoFileName, string(text), false)
	}
	if err != nil {
		return &compiler.EmitResult{
			EmitSkipped: true,
			Diagnostics: []*ast.Diagnostic{
				ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, buildInfoFileName, err.Error()),
			},
		}
	}
	p.snapshot.buildInfoEmitPending = false

	var emittedFiles []string
	if p.snapshot.options.ListEmittedFiles.IsTrue() {
		emittedFiles = []string{buildInfoFileName}
	}
	return &compiler.EmitResult{
		EmitSkipped:  false,
		EmittedFiles: emittedFiles,
	}
}

func (p *Program) ensureHasErrorsForState(ctx context.Context, program *compiler.Program) core.Tristate {
	if p.snapshot.hasErrors != core.TSUnknown {
		return p.snapshot.hasErrors
	}

	// Check semantic and emit diagnostics first as we dont need to ask program about it
	if slices.ContainsFunc(program.GetSourceFiles(), func(file *ast.SourceFile) bool {
		semanticDiagnostics := p.snapshot.semanticDiagnosticsPerFile[file.Path()]
		if semanticDiagnostics == nil {
			// Missing semantic diagnostics in cache will be encoded in incremental buildInfo
			return p.snapshot.options.IsIncremental()
		}
		if len(semanticDiagnostics.diagnostics) > 0 || len(semanticDiagnostics.buildInfoDiagnostics) > 0 {
			// cached semantic diagnostics will be encoded in buildInfo
			return true
		}
		if _, ok := p.snapshot.emitDiagnosticsPerFile[file.Path()]; ok {
			// emit diagnostics will be encoded in buildInfo;
			return true
		}
		return false
	}) {
		// Because semantic diagnostics are recorded in buildInfo, we dont need to encode hasErrors in incremental buildInfo
		// But encode as errors in non incremental buildInfo
		return core.IfElse(p.snapshot.options.IsIncremental(), core.TSFalse, core.TSTrue)
	}
	if len(program.GetConfigFileParsingDiagnostics()) > 0 ||
		len(program.GetSyntacticDiagnostics(ctx, nil)) > 0 ||
		len(program.GetBindDiagnostics(ctx, nil)) > 0 ||
		len(program.GetOptionsDiagnostics(ctx)) > 0 {
		return core.TSTrue
	} else {
		return core.TSFalse
	}
}
