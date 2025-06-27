package incremental

import (
	"context"
	"crypto/sha256"
	"fmt"
	"maps"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type programState struct {
	// Used during incremental updates

	// !!! sheetal handle parallel updates and state sanity
	/**
	* Map of files that have already called update signature.
	* That means hence forth these files are assumed to have
	* no change in their signature for this version of the program
	 */
	hasCalledUpdateShapeSignature collections.SyncSet[tspath.Path]
	/**
	 * whether this program has cleaned semantic diagnostics cache for lib files
	 */
	cleanedDiagnosticsOfLibFiles bool
	/**
	 * Stores signatures before before the update till affected file is committed
	 */
	oldSignatures map[tspath.Path]string
	/**
	 * Current changed file for iterating over affected files
	 */
	currentChangedFilePath tspath.Path
	/**
	 * Set of affected files being iterated
	 */
	affectedFiles []*ast.SourceFile
	/**
	 * Current index to retrieve affected file from
	 */
	affectedFilesIndex int
	/**
	 * Already seen affected files
	 */
	seenAffectedFiles collections.Set[tspath.Path]
	/**
	 * Already seen emitted files
	 */
	seenEmittedFiles map[tspath.Path]FileEmitKind
}

func (p *programState) emit(ctx context.Context, program *compiler.Program, options compiler.EmitOptions) *compiler.EmitResult {
	if result := compiler.HandleNoEmitOptions(ctx, program, options.TargetSourceFile); result != nil {
		if options.TargetSourceFile != nil {
			return result
		}

		// Emit buildInfo and combine result
		buildInfoResult := p.emitBuildInfo(ctx, program, options)
		if buildInfoResult != nil && buildInfoResult.EmittedFiles != nil {
			result.Diagnostics = append(result.Diagnostics, buildInfoResult.Diagnostics...)
			result.EmittedFiles = append(result.EmittedFiles, buildInfoResult.EmittedFiles...)
		}
		return result
	}

	// Emit only affected files if using builder for emit
	if options.TargetSourceFile != nil {
		return program.Emit(ctx, p.getEmitOptions(program, options))
	}

	var results []*compiler.EmitResult
	for {
		affectedEmitResult, done := p.emitNextAffectedFile(ctx, program, options, false)
		if done {
			break
		}
		results = append(results, affectedEmitResult)
	}
	return compiler.CombineEmitResults(results)
}

func (p *programState) getDeclarationDiagnostics(ctx context.Context, program *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
	var diagnostics []*ast.Diagnostic
	for {
		affectedEmitResult, done := p.emitNextAffectedFile(ctx, program, compiler.EmitOptions{}, true)
		if done {
			break
		}
		if file == nil {
			diagnostics = append(diagnostics, affectedEmitResult.Diagnostics...)
		}
	}
	if file == nil {
		return diagnostics
	}
	if emitDiagnostics, ok := p.emitDiagnosticsPerFile[file.Path()]; ok {
		// If diagnostics are present for the file, return them
		return emitDiagnostics.getDiagnostics(program, file)
	}
	return nil
}

/**
 * Emits the next affected file's emit result (EmitResult and sourceFiles emitted) or returns undefined if iteration is complete
 * The first of writeFile if provided, writeFile of BuilderProgramHost if provided, writeFile of compiler host
 * in that order would be used to write the files
 */
func (p *programState) emitNextAffectedFile(ctx context.Context, program *compiler.Program, options compiler.EmitOptions, isForDtsErrors bool) (*compiler.EmitResult, bool) {
	affected := p.getNextAffectedFile(ctx, program)
	programEmitKind := GetFileEmitKind(p.options)
	var emitKind FileEmitKind
	if affected == nil {
		// file pending emit
		pendingAffectedFile, pendingEmitKind := p.getNextAffectedFilePendingEmit(program, options, isForDtsErrors)
		if pendingAffectedFile != nil {
			affected = pendingAffectedFile
			emitKind = pendingEmitKind
		} else {
			// File whose diagnostics need to be reported
			affectedFile, pendingDiagnostics, seenKind := p.getNextPendingEmitDiagnosticsFile(program, isForDtsErrors)
			if affectedFile != nil {
				p.seenEmittedFiles[affectedFile.Path()] = seenKind | getFileEmitKindAllDts(isForDtsErrors)
				return &compiler.EmitResult{
					EmitSkipped: true,
					Diagnostics: pendingDiagnostics.getDiagnostics(program, affectedFile),
				}, false
			}
		}
		if affected == nil {
			// Emit buildinfo if pending
			if isForDtsErrors {
				return nil, true
			}
			result := p.emitBuildInfo(ctx, program, options)
			if result != nil {
				return result, false
			}
			return nil, true
		}
	} else {
		if isForDtsErrors {
			emitKind = fileEmitKindDtsErrors
		} else if options.EmitOnly == compiler.EmitOnlyDts {
			emitKind = programEmitKind & fileEmitKindAllDts
		} else {
			emitKind = programEmitKind
		}
	}
	// Determine if we can do partial emit
	var emitOnly compiler.EmitOnly
	if (emitKind & fileEmitKindAllJs) != 0 {
		emitOnly = compiler.EmitOnlyJs
	}
	if (emitKind & fileEmitKindAllDts) != 0 {
		if emitOnly == compiler.EmitOnlyJs {
			emitOnly = compiler.EmitAll
		} else {
			emitOnly = compiler.EmitOnlyDts
		}
	}
	// // Actual emit without buildInfo as we want to emit it later so the state is updated
	var result *compiler.EmitResult
	if !isForDtsErrors {
		result = program.Emit(ctx, p.getEmitOptions(program, compiler.EmitOptions{
			TargetSourceFile: affected,
			EmitOnly:         emitOnly,
			WriteFile:        options.WriteFile,
		}))
	} else {
		result = &compiler.EmitResult{
			EmitSkipped: true,
			Diagnostics: program.GetDeclarationDiagnostics(ctx, affected),
		}
	}

	// update affected files
	p.seenAffectedFiles.Add(affected.Path())
	p.affectedFilesIndex++
	// Change in changeSet/affectedFilesPendingEmit, buildInfo needs to be emitted
	p.buildInfoEmitPending = true
	// Update the pendingEmit for the file
	existing := p.seenEmittedFiles[affected.Path()]
	p.seenEmittedFiles[affected.Path()] = emitKind | existing
	existingPending, ok := p.affectedFilesPendingEmit[affected.Path()]
	if !ok {
		existingPending = programEmitKind
	}
	pendingKind := getPendingEmitKind(existingPending, emitKind|existing)
	if pendingKind != 0 {
		p.affectedFilesPendingEmit[affected.Path()] = pendingKind
	} else {
		delete(p.affectedFilesPendingEmit, affected.Path())
	}
	if len(result.Diagnostics) != 0 {
		if p.emitDiagnosticsPerFile == nil {
			p.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName)
		}
		p.emitDiagnosticsPerFile[affected.Path()] = &diagnosticsOrBuildInfoDiagnosticsWithFileName{
			diagnostics: result.Diagnostics,
		}
	}
	return result, false
}

/**
 * Returns next file to be emitted from files that retrieved semantic diagnostics but did not emit yet
 */
func (p *programState) getNextAffectedFilePendingEmit(program *compiler.Program, options compiler.EmitOptions, isForDtsErrors bool) (*ast.SourceFile, FileEmitKind) {
	if len(p.affectedFilesPendingEmit) == 0 {
		return nil, 0
	}
	for path, emitKind := range p.affectedFilesPendingEmit {
		affectedFile := program.GetSourceFileByPath(path)
		if affectedFile == nil || !program.SourceFileMayBeEmitted(affectedFile, false) {
			delete(p.affectedFilesPendingEmit, path)
			continue
		}
		seenKind := p.seenEmittedFiles[affectedFile.Path()]
		pendingKind := getPendingEmitKindWithSeen(emitKind, seenKind, options, isForDtsErrors)
		if pendingKind != 0 {
			return affectedFile, pendingKind
		}
	}
	return nil, 0
}

func (p *programState) getNextPendingEmitDiagnosticsFile(program *compiler.Program, isForDtsErrors bool) (*ast.SourceFile, *diagnosticsOrBuildInfoDiagnosticsWithFileName, FileEmitKind) {
	if len(p.emitDiagnosticsPerFile) == 0 {
		return nil, nil, 0
	}
	for path, diagnostics := range p.emitDiagnosticsPerFile {
		affectedFile := program.GetSourceFileByPath(path)
		if affectedFile == nil || !program.SourceFileMayBeEmitted(affectedFile, false) {
			delete(p.emitDiagnosticsPerFile, path)
			continue
		}
		seenKind := p.seenEmittedFiles[affectedFile.Path()]
		if (seenKind & getFileEmitKindAllDts(isForDtsErrors)) != 0 {
			return affectedFile, diagnostics, seenKind
		}
	}
	return nil, nil, 0
}

func (p *programState) getEmitOptions(program *compiler.Program, options compiler.EmitOptions) compiler.EmitOptions {
	if !p.options.GetEmitDeclarations() {
		return options
	}
	return compiler.EmitOptions{
		TargetSourceFile: options.TargetSourceFile,
		EmitOnly:         options.EmitOnly,
		WriteFile: func(fileName string, text string, writeByteOrderMark bool, data *compiler.WriteFileData) error {
			if tspath.IsDeclarationFileName(fileName) {
				var emitSignature string
				info := p.fileInfos[options.TargetSourceFile.Path()]
				if info.signature == info.version {
					signature := computeSignatureWithDiagnostics(options.TargetSourceFile, text, data)
					// With d.ts diagnostics they are also part of the signature so emitSignature will be different from it since its just hash of d.ts
					if len(data.Diagnostics) == 0 {
						emitSignature = signature
					}
					if signature != info.version { // Update it
						if p.affectedFiles != nil {
							// Keep old signature so we know what to undo if cancellation happens
							if _, ok := p.oldSignatures[options.TargetSourceFile.Path()]; !ok {
								if p.oldSignatures == nil {
									p.oldSignatures = make(map[tspath.Path]string)
								}
								p.oldSignatures[options.TargetSourceFile.Path()] = info.signature
							}
						}
						info.signature = signature
					}
				}

				// Store d.ts emit hash so later can be compared to check if d.ts has changed.
				// Currently we do this only for composite projects since these are the only projects that can be referenced by other projects
				// and would need their d.ts change time in --build mode
				if p.skipDtsOutputOfComposite(program, options.TargetSourceFile, fileName, text, data, emitSignature) {
					return nil
				}
			}

			if options.WriteFile != nil {
				return options.WriteFile(fileName, text, writeByteOrderMark, data)
			}
			return program.Host().FS().WriteFile(fileName, text, writeByteOrderMark)
		},
	}
}

/**
 * Compare to existing computed signature and store it or handle the changes in d.ts map option from before
 * returning undefined means that, we dont need to emit this d.ts file since its contents didnt change
 */
func (p *programState) skipDtsOutputOfComposite(program *compiler.Program, file *ast.SourceFile, outputFileName string, text string, data *compiler.WriteFileData, newSignature string) bool {
	if !p.options.Composite.IsTrue() {
		return false
	}
	var oldSignature string
	oldSignatureFormat, ok := p.emitSignatures[file.Path()]
	if ok {
		if oldSignatureFormat.signature != "" {
			oldSignature = oldSignatureFormat.signature
		} else {
			oldSignature = oldSignatureFormat.signatureWithDifferentOptions[0]
		}
	}
	if newSignature == "" {
		newSignature = computeHash(getTextHandlingSourceMapForSignature(text, data))
	}
	// Dont write dts files if they didn't change
	if newSignature == oldSignature {
		// If the signature was encoded as string the dts map options match so nothing to do
		if oldSignatureFormat != nil && oldSignatureFormat.signature == oldSignature {
			data.SkippedDtsWrite = true
			return true
		} else {
			// Mark as differsOnlyInMap so that --build can reverse the timestamp so that
			// the downstream projects dont detect this as change in d.ts file
			data.DiffersOnlyInMap = true
		}
	} else {
		p.latestChangedDtsFile = outputFileName
	}
	if p.emitSignatures == nil {
		p.emitSignatures = make(map[tspath.Path]*emitSignature)
	}
	p.emitSignatures[file.Path()] = &emitSignature{
		signature: newSignature,
	}
	return false
}





func newProgramState(program *compiler.Program, oldProgram *Program) *programState {
	if oldProgram != nil && oldProgram.program == program {
		return oldProgram.state
	}
	files := program.GetSourceFiles()
	state := &programState{
		options:                    program.Options(),
		semanticDiagnosticsPerFile: make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(files)),
		seenEmittedFiles:           make(map[tspath.Path]FileEmitKind, len(files)),
	}
	state.createReferenceMap()
	if oldProgram != nil && state.options.Composite.IsTrue() {
		state.latestChangedDtsFile = oldProgram.state.latestChangedDtsFile
	}
	if state.options.NoCheck.IsTrue() {
		state.checkPending = true
	}

	canUseStateFromOldProgram := oldProgram != nil && state.tracksReferences() == oldProgram.state.tracksReferences()
	if canUseStateFromOldProgram {
		// Copy old state's changed files set
		state.changedFilesSet = oldProgram.state.changedFilesSet.Clone()
		if len(oldProgram.state.affectedFilesPendingEmit) != 0 {
			state.affectedFilesPendingEmit = maps.Clone(oldProgram.state.affectedFilesPendingEmit)
		}
		state.hasErrorsFromOldState = oldProgram.state.hasErrors
	} else {
		state.changedFilesSet = &collections.Set[tspath.Path]{}
		state.useFileVersionAsSignature = true
		state.buildInfoEmitPending = state.options.IsIncremental()
	}

	canCopySemanticDiagnostics := canUseStateFromOldProgram &&
		!tsoptions.CompilerOptionsAffectSemanticDiagnostics(oldProgram.state.options, program.Options())
	// // We can only reuse emit signatures (i.e. .d.ts signatures) if the .d.ts file is unchanged,
	// // which will eg be depedent on change in options like declarationDir and outDir options are unchanged.
	// // We need to look in oldState.compilerOptions, rather than oldCompilerOptions (i.e.we need to disregard useOldState) because
	// // oldCompilerOptions can be undefined if there was change in say module from None to some other option
	// // which would make useOldState as false since we can now use reference maps that are needed to track what to emit, what to check etc
	// // but that option change does not affect d.ts file name so emitSignatures should still be reused.
	canCopyEmitSignatures := state.options.Composite.IsTrue() &&
		oldProgram != nil &&
		oldProgram.state.emitSignatures != nil &&
		!tsoptions.CompilerOptionsAffectDeclarationPath(oldProgram.state.options, program.Options())
	copyDeclarationFileDiagnostics := canCopySemanticDiagnostics &&
		state.options.SkipLibCheck.IsTrue() == oldProgram.state.options.SkipLibCheck.IsTrue()
	copyLibFileDiagnostics := copyDeclarationFileDiagnostics &&
		state.options.SkipDefaultLibCheck.IsTrue() == oldProgram.state.options.SkipDefaultLibCheck.IsTrue()
	state.fileInfos = make(map[tspath.Path]*fileInfo, len(files))
	for _, file := range files {
		version := computeHash(file.Text())
		impliedNodeFormat := program.GetSourceFileMetaData(file.Path()).ImpliedNodeFormat
		affectsGlobalScope := fileAffectsGlobalScope(file)
		var signature string
		if canUseStateFromOldProgram {
			var hasOldUncommitedSignature bool
			signature, hasOldUncommitedSignature = oldProgram.state.oldSignatures[file.Path()]
			if oldFileInfo, ok := oldProgram.state.fileInfos[file.Path()]; ok {
				if !hasOldUncommitedSignature {
					signature = oldFileInfo.signature
				}
				if oldFileInfo.version == version || oldFileInfo.affectsGlobalScope != affectsGlobalScope || oldFileInfo.impliedNodeFormat != impliedNodeFormat {
					state.addFileToChangeSet(file.Path())
				}
			} else {
				state.addFileToChangeSet(file.Path())
			}
			if state.referencedMap != nil {
				newReferences := getReferencedFiles(program, file)
				if newReferences != nil {
					state.referencedMap.Add(file.Path(), newReferences)
				}
				oldReferences, _ := oldProgram.state.referencedMap.GetValues(file.Path())
				// Referenced files changed
				if !newReferences.Equals(oldReferences) {
					state.addFileToChangeSet(file.Path())
				} else {
					for refPath := range newReferences.Keys() {
						if program.GetSourceFileByPath(refPath) == nil {
							// Referenced file was deleted in the new program
							state.addFileToChangeSet(file.Path())
							break
						}
					}
				}
			}
			if !state.changedFilesSet.Has(file.Path()) {
				if emitDiagnostics, ok := oldProgram.state.emitDiagnosticsPerFile[file.Path()]; ok {
					if state.emitDiagnosticsPerFile == nil {
						state.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(files))
					}
					state.emitDiagnosticsPerFile[file.Path()] = emitDiagnostics
				}
				if canCopySemanticDiagnostics {
					if (!file.IsDeclarationFile || copyDeclarationFileDiagnostics) &&
						(!program.IsSourceFileDefaultLibrary(file.Path()) || copyLibFileDiagnostics) {
						// Unchanged file copy diagnostics
						if diagnostics, ok := oldProgram.state.semanticDiagnosticsPerFile[file.Path()]; ok {
							state.semanticDiagnosticsPerFile[file.Path()] = diagnostics
							state.semanticDiagnosticsFromOldState.Add(file.Path())
						}
					}
				}
			}
			if canCopyEmitSignatures {
				if oldEmitSignature, ok := oldProgram.state.emitSignatures[file.Path()]; ok {
					state.createEmitSignaturesMap()
					state.emitSignatures[file.Path()] = oldEmitSignature.getNewEmitSignature(oldProgram.state.options, state.options)
				}
			}
		} else {
			state.addFileToChangeSet(file.Path())
		}
		state.fileInfos[file.Path()] = &fileInfo{
			version:            version,
			signature:          signature,
			affectsGlobalScope: affectsGlobalScope,
			impliedNodeFormat:  impliedNodeFormat,
		}
	}
	if canUseStateFromOldProgram {
		// If the global file is removed, add all files as changed
		allFilesExcludingDefaultLibraryFileAddedToChangeSet := false
		for filePath, oldInfo := range oldProgram.state.fileInfos {
			if _, ok := state.fileInfos[filePath]; !ok {
				if oldInfo.affectsGlobalScope {
					for _, file := range state.getAllFilesExcludingDefaultLibraryFile(program, nil) {
						state.addFileToChangeSet(file.Path())
					}
					allFilesExcludingDefaultLibraryFileAddedToChangeSet = true
				} else {
					state.buildInfoEmitPending = true
				}
				break
			}
		}
		if !allFilesExcludingDefaultLibraryFileAddedToChangeSet {
			// If options affect emit, then we need to do complete emit per compiler options
			// otherwise only the js or dts that needs to emitted because its different from previously emitted options
			var pendingEmitKind FileEmitKind
			if tsoptions.CompilerOptionsAffectEmit(oldProgram.state.options, state.options) {
				pendingEmitKind = GetFileEmitKind(state.options)
			} else {
				pendingEmitKind = getPendingEmitKindWithOptions(state.options, oldProgram.state.options)
			}
			if pendingEmitKind != fileEmitKindNone {
				// Add all files to affectedFilesPendingEmit since emit changed
				for _, file := range files {
					// Add to affectedFilesPending emit only if not changed since any changed file will do full emit
					if !state.changedFilesSet.Has(file.Path()) {
						state.addFileToAffectedFilesPendingEmit(file.Path(), pendingEmitKind)
					}
				}
				state.buildInfoEmitPending = true
			}
		}
	}
	if canUseStateFromOldProgram &&
		len(state.semanticDiagnosticsPerFile) != len(state.fileInfos) &&
		oldProgram.state.checkPending != state.checkPending {
		state.buildInfoEmitPending = true
	}
	return state
}

func fileAffectsGlobalScope(file *ast.SourceFile) bool {
	// if file contains anything that augments to global scope we need to build them as if
	// they are global files as well as module
	if core.Some(file.ModuleAugmentations, func(augmentation *ast.ModuleName) bool {
		return ast.IsGlobalScopeAugmentation(augmentation.Parent)
	}) {
		return true
	}

	if ast.IsExternalOrCommonJSModule(file) || ast.IsJsonSourceFile(file) {
		return false
	}

	/**
	 * For script files that contains only ambient external modules, although they are not actually external module files,
	 * they can only be consumed via importing elements from them. Regular script files cannot consume them. Therefore,
	 * there are no point to rebuild all script files if these special files have changed. However, if any statement
	 * in the file is not ambient external module, we treat it as a regular script file.
	 */
	return file.Statements != nil &&
		file.Statements.Nodes != nil &&
		core.Some(file.Statements.Nodes, func(stmt *ast.Node) bool {
			return !ast.IsModuleWithStringLiteralName(stmt)
		})
}

func getTextHandlingSourceMapForSignature(text string, data *compiler.WriteFileData) string {
	if data.SourceMapUrlPos != -1 {
		return text[:data.SourceMapUrlPos]
	}
	return text
}

func computeSignatureWithDiagnostics(file *ast.SourceFile, text string, data *compiler.WriteFileData) string {
	var builder strings.Builder
	builder.WriteString(getTextHandlingSourceMapForSignature(text, data))
	for _, diag := range data.Diagnostics {
		diagnosticToStringBuilder(diag, file, &builder)
	}
	return computeHash(builder.String())
}

func diagnosticToStringBuilder(diagnostic *ast.Diagnostic, file *ast.SourceFile, builder *strings.Builder) string {
	if diagnostic == nil {
		return ""
	}
	builder.WriteString("\n")
	if diagnostic.File() != file {
		builder.WriteString(tspath.EnsurePathIsNonModuleName(tspath.GetRelativePathFromDirectory(
			tspath.GetDirectoryPath(string(file.Path())),
			string(diagnostic.File().Path()),
			tspath.ComparePathsOptions{},
		)))
	}
	if diagnostic.File() != nil {
		builder.WriteString(fmt.Sprintf("(%d,%d): ", diagnostic.Pos(), diagnostic.Len()))
	}
	builder.WriteString(diagnostic.Category().Name())
	builder.WriteString(fmt.Sprintf("%d: ", diagnostic.Code()))
	builder.WriteString(diagnostic.Message())
	for _, chain := range diagnostic.MessageChain() {
		diagnosticToStringBuilder(chain, file, builder)
	}
	for _, info := range diagnostic.RelatedInformation() {
		diagnosticToStringBuilder(info, file, builder)
	}
	return builder.String()
}

func computeHash(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

/**
* Get the module source file and all augmenting files from the import name node from file
 */
func addReferencedFilesFromImportLiteral(file *ast.SourceFile, referencedFiles *collections.Set[tspath.Path], checker *checker.Checker, importName *ast.LiteralLikeNode) {
	symbol := checker.GetSymbolAtLocation(importName)
	if symbol == nil {
		return
	}
	for _, declaration := range symbol.Declarations {
		fileOfDecl := ast.GetSourceFileOfNode(declaration)
		if fileOfDecl == nil {
			continue
		}
		if file != fileOfDecl {
			referencedFiles.Add(fileOfDecl.Path())
		}
	}
}

/**
* Gets the path to reference file from file name, it could be resolvedPath if present otherwise path
 */
func addReferencedFileFromFileName(program *compiler.Program, fileName string, referencedFiles *collections.Set[tspath.Path], sourceFileDirectory string) {
	if redirect := program.GetParseFileRedirect(fileName); redirect != "" {
		referencedFiles.Add(tspath.ToPath(redirect, program.GetCurrentDirectory(), program.UseCaseSensitiveFileNames()))
	} else {
		referencedFiles.Add(tspath.ToPath(fileName, sourceFileDirectory, program.UseCaseSensitiveFileNames()))
	}
}

/**
 * Gets the referenced files for a file from the program with values for the keys as referenced file's path to be true
 */
func getReferencedFiles(program *compiler.Program, file *ast.SourceFile) *collections.Set[tspath.Path] {
	referencedFiles := collections.Set[tspath.Path]{}

	// We need to use a set here since the code can contain the same import twice,
	// but that will only be one dependency.
	// To avoid invernal conversion, the key of the referencedFiles map must be of type Path
	if len(file.Imports()) > 0 || len(file.ModuleAugmentations) > 0 {
		checker, done := program.GetTypeCheckerForFile(context.TODO(), file)
		for _, importName := range file.Imports() {
			addReferencedFilesFromImportLiteral(file, &referencedFiles, checker, importName)
		}
		// Add module augmentation as references
		for _, moduleName := range file.ModuleAugmentations {
			if !ast.IsStringLiteral(moduleName) {
				continue
			}
			addReferencedFilesFromImportLiteral(file, &referencedFiles, checker, moduleName)
		}
		done()
	}

	sourceFileDirectory := tspath.GetDirectoryPath(file.FileName())
	// Handle triple slash references
	for _, referencedFile := range file.ReferencedFiles {
		addReferencedFileFromFileName(program, referencedFile.FileName, &referencedFiles, sourceFileDirectory)
	}

	// Handle type reference directives
	if typeRefsInFile, ok := program.GetResolvedTypeReferenceDirectives()[file.Path()]; ok {
		for _, typeRef := range typeRefsInFile {
			if typeRef.ResolvedFileName != "" {
				addReferencedFileFromFileName(program, typeRef.ResolvedFileName, &referencedFiles, sourceFileDirectory)
			}
		}
	}

	// !!! sheetal
	// // From ambient modules
	// for (const ambientModule of program.getTypeChecker().getAmbientModules()) {
	//     if (ambientModule.declarations && ambientModule.declarations.length > 1) {
	//         addReferenceFromAmbientModule(ambientModule);
	//     }
	// }
	return core.IfElse(referencedFiles.Len() > 0, &referencedFiles, nil)
}
