package incremental

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type fileInfo struct {
	version            string
	signature          string
	affectsGlobalScope bool
	impliedNodeFormat  core.ResolutionMode
}

func (f *fileInfo) Version() string                        { return f.version }
func (f *fileInfo) Signature() string                      { return f.signature }
func (f *fileInfo) AffectsGlobalScope() bool               { return f.affectsGlobalScope }
func (f *fileInfo) ImpliedNodeFormat() core.ResolutionMode { return f.impliedNodeFormat }

type fileEmitKind uint32

const (
	fileEmitKindNone        fileEmitKind = 0
	fileEmitKindJs          fileEmitKind = 1 << 0 // emit js file
	fileEmitKindJsMap       fileEmitKind = 1 << 1 // emit js.map file
	fileEmitKindJsInlineMap fileEmitKind = 1 << 2 // emit inline source map in js file
	fileEmitKindDtsErrors   fileEmitKind = 1 << 3 // emit dts errors
	fileEmitKindDtsEmit     fileEmitKind = 1 << 4 // emit d.ts file
	fileEmitKindDtsMap      fileEmitKind = 1 << 5 // emit d.ts.map file

	fileEmitKindDts        = fileEmitKindDtsErrors | fileEmitKindDtsEmit
	fileEmitKindAllJs      = fileEmitKindJs | fileEmitKindJsMap | fileEmitKindJsInlineMap
	fileEmitKindAllDtsEmit = fileEmitKindDtsEmit | fileEmitKindDtsMap
	fileEmitKindAllDts     = fileEmitKindDts | fileEmitKindDtsMap
	fileEmitKindAll        = fileEmitKindAllJs | fileEmitKindAllDts
)

func getFileEmitKind(options *core.CompilerOptions) fileEmitKind {
	result := fileEmitKindJs
	if options.SourceMap.IsTrue() {
		result |= fileEmitKindJsMap
	}
	if options.InlineSourceMap.IsTrue() {
		result |= fileEmitKindJsInlineMap
	}
	if options.GetEmitDeclarations() {
		result |= fileEmitKindDts
	}
	if options.DeclarationMap.IsTrue() {
		result |= fileEmitKindDtsMap
	}
	if options.EmitDeclarationOnly.IsTrue() {
		result &= fileEmitKindAllDts
	}
	return result
}

func getPendingEmitKindWithOptions(options *core.CompilerOptions, oldOptions *core.CompilerOptions) fileEmitKind {
	oldEmitKind := getFileEmitKind(oldOptions)
	newEmitKind := getFileEmitKind(options)
	return getPendingEmitKind(newEmitKind, oldEmitKind)
}

func getPendingEmitKind(emitKind fileEmitKind, oldEmitKind fileEmitKind) fileEmitKind {
	if oldEmitKind == emitKind {
		return fileEmitKindNone
	}
	if oldEmitKind == 0 || emitKind == 0 {
		return emitKind
	}
	diff := oldEmitKind ^ emitKind
	result := fileEmitKindNone
	// If there is diff in Js emit, pending emit is js emit flags
	if (diff & fileEmitKindAllJs) != 0 {
		result |= emitKind & fileEmitKindAllJs
	}
	// If dts errors pending, add dts errors flag
	if (diff & fileEmitKindDtsErrors) != 0 {
		result |= emitKind & fileEmitKindDtsErrors
	}
	// If there is diff in Dts emit, pending emit is dts emit flags
	if (diff & fileEmitKindAllDtsEmit) != 0 {
		result |= emitKind & fileEmitKindAllDtsEmit
	}
	return result
}

/**
 * Determining what all is pending to be emitted based on previous options or previous file emit flags
 *  @internal
 */
func getPendingEmitKindWithSeen(emitKind fileEmitKind, seenEmitKind fileEmitKind, options compiler.EmitOptions, isForDtsErrors bool) fileEmitKind {
	pendingKind := getPendingEmitKind(emitKind, seenEmitKind)
	if options.EmitOnly == compiler.EmitOnlyDts {
		pendingKind &= fileEmitKindAllDts
	}
	if isForDtsErrors {
		pendingKind &= fileEmitKindDtsErrors
	}
	return pendingKind
}

func getFileEmitKindAllDts(isForDtsErrors bool) fileEmitKind {
	return core.IfElse(isForDtsErrors, fileEmitKindDtsErrors, fileEmitKindAllDts)
}

/**
 * Signature (Hash of d.ts emitted), is string if it was emitted using same d.ts.map option as what compilerOptions indicate,
 * otherwise tuple of string
 */
type emitSignature struct {
	signature                     string
	signatureWithDifferentOptions []string
}

/**
 * Covert to Emit signature based on oldOptions and EmitSignature format
 * If d.ts map options differ then swap the format, otherwise use as is
 */
func (e *emitSignature) getNewEmitSignature(oldOptions *core.CompilerOptions, newOptions *core.CompilerOptions) *emitSignature {
	if oldOptions.DeclarationMap.IsTrue() == newOptions.DeclarationMap.IsTrue() {
		return e
	}
	if e.signatureWithDifferentOptions == nil {
		return &emitSignature{
			signatureWithDifferentOptions: []string{e.signature},
		}
	} else {
		return &emitSignature{
			signature: e.signatureWithDifferentOptions[0],
		}
	}
}

type diagnosticsOrBuildInfoDiagnostics struct {
	diagnostics          []*ast.Diagnostic
	buildInfoDiagnostics []*BuildInfoDiagnostic
}

func (d *diagnosticsOrBuildInfoDiagnostics) getDiagnostics(p *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
	if d.diagnostics != nil {
		return d.diagnostics
	}
	// Convert and cache the diagnostics
	d.diagnostics = core.Map(d.buildInfoDiagnostics, func(diag *BuildInfoDiagnostic) *ast.Diagnostic {
		return diag.toDiagnostic(p, file)
	})
	return d.diagnostics
}

type programState struct {
	// State that is serialized as buildinfo
	/**
	 * Information of the file eg. its version, signature etc
	 */
	fileInfos map[tspath.Path]*fileInfo
	options   *core.CompilerOptions
	/**
	 * Contains the map of ReferencedSet=Referenced files of the file if module emit is enabled
	 */
	referencedMap *collections.ManyToManyMap[tspath.Path, tspath.Path]
	/**
	 * Cache of bind and check diagnostics for files with their Path being the key
	 */
	semanticDiagnosticsPerFile map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics
	/** Cache of dts emit diagnostics for files with their Path being the key */
	emitDiagnosticsPerFile map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics
	/**
	 * The map has key by source file's path that has been changed
	 */
	changedFilesSet *collections.Set[tspath.Path]
	/**
	 * Files pending to be emitted
	 */
	affectedFilesPendingEmit map[tspath.Path]fileEmitKind
	/**
	 * Name of the file whose dts was the latest to change
	 */
	latestChangedDtsFile string
	/**
	 * Hash of d.ts emitted for the file, use to track when emit of d.ts changes
	 */
	emitSignatures map[tspath.Path]*emitSignature
	/** Recorded if program had errors */
	hasErrors core.Tristate
	/** If semantic diagnsotic check is pending */
	checkPending bool

	// Used during incremental updates
	/**
	 * true if file version is used as signature
	 * This helps in delaying the calculation of the d.ts hash as version for the file till reasonable time
	 */
	useFileVersionAsSignature bool
	/**
	 * true if build info is emitted
	 */
	buildInfoEmitPending  bool
	hasErrorsFromOldState core.Tristate
	/**
	 * True if the semantic diagnostics were copied from the old state
	 */
	semanticDiagnosticsFromOldState collections.Set[tspath.Path]
	/**
	 * Cache of all files excluding default library file for the current program
	 */
	allFilesExcludingDefaultLibraryFile []*ast.SourceFile

	// !!! sheetal handle parallel updates and state sanity
	/**
	* Map of files that have already called update signature.
	* That means hence forth these files are assumed to have
	* no change in their signature for this version of the program
	 */
	hasCalledUpdateShapeSignature collections.Set[tspath.Path]
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
	seenEmittedFiles map[tspath.Path]fileEmitKind
	/**
	 * Records if change in dts emit was detected
	 */
	hasChangedEmitSignature bool
}

func (p *programState) tracksReferences() bool {
	return p.options.Module != core.ModuleKindNone
}

func (p *programState) createReferenceMap() {
	if p.tracksReferences() {
		p.referencedMap = &collections.ManyToManyMap[tspath.Path, tspath.Path]{}
	}
}

func (p *programState) createEmitSignaturesMap() {
	if p.emitSignatures == nil && p.options.Composite.IsTrue() {
		p.emitSignatures = make(map[tspath.Path]*emitSignature)
	}
}

func (p *programState) addFileToChangeSet(filePath tspath.Path) {
	p.changedFilesSet.Add(filePath)
	p.buildInfoEmitPending = true
}

func (p *programState) addFileToAffectedFilesPendingEmit(filePath tspath.Path, emitKind fileEmitKind) {
	existingKind := p.affectedFilesPendingEmit[filePath]

	if p.affectedFilesPendingEmit == nil {
		p.affectedFilesPendingEmit = make(map[tspath.Path]fileEmitKind)
	}
	p.affectedFilesPendingEmit[filePath] = existingKind | emitKind
	delete(p.emitDiagnosticsPerFile, filePath)
}

func (p *programState) getAllFilesExcludingDefaultLibraryFile(program *compiler.Program, firstSourceFile *ast.SourceFile) []*ast.SourceFile {
	// Use cached result
	if p.allFilesExcludingDefaultLibraryFile != nil {
		return p.allFilesExcludingDefaultLibraryFile
	}

	files := program.GetSourceFiles()
	p.allFilesExcludingDefaultLibraryFile = make([]*ast.SourceFile, 0, len(files))
	addSourceFile := func(file *ast.SourceFile) {
		if !program.IsSourceFileDefaultLibrary(file.Path()) {
			p.allFilesExcludingDefaultLibraryFile = append(p.allFilesExcludingDefaultLibraryFile, file)
		}
	}
	if firstSourceFile != nil {
		addSourceFile(firstSourceFile)
	}
	for _, file := range files {
		if file != firstSourceFile {
			addSourceFile(file)
		}
	}
	return p.allFilesExcludingDefaultLibraryFile
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
	programEmitKind := getFileEmitKind(p.options)
	var emitKind fileEmitKind
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
			p.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics)
		}
		p.emitDiagnosticsPerFile[affected.Path()] = &diagnosticsOrBuildInfoDiagnostics{
			diagnostics: result.Diagnostics,
		}
	}
	return result, false
}

/**
 * Returns next file to be emitted from files that retrieved semantic diagnostics but did not emit yet
 */
func (p *programState) getNextAffectedFilePendingEmit(program *compiler.Program, options compiler.EmitOptions, isForDtsErrors bool) (*ast.SourceFile, fileEmitKind) {
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

func (p *programState) getNextPendingEmitDiagnosticsFile(program *compiler.Program, isForDtsErrors bool) (*ast.SourceFile, *diagnosticsOrBuildInfoDiagnostics, fileEmitKind) {
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
		p.hasChangedEmitSignature = true
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

func (p *programState) emitBuildInfo(ctx context.Context, program *compiler.Program, options compiler.EmitOptions) *compiler.EmitResult {
	buildInfoFileName := outputpaths.GetBuildInfoFileName(p.options, tspath.ComparePathsOptions{
		CurrentDirectory:          program.GetCurrentDirectory(),
		UseCaseSensitiveFileNames: program.UseCaseSensitiveFileNames(),
	})
	if buildInfoFileName == "" {
		return nil
	}

	p.ensureHasErrorsForState(ctx, program)
	if !p.buildInfoEmitPending && p.hasErrorsFromOldState == p.hasErrors {
		return nil
	}
	p.buildInfoEmitPending = false
	p.hasErrorsFromOldState = p.hasErrors
	buildInfo := programStateToBuildInfo(p, program, buildInfoFileName)
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
	var emittedFiles []string
	if p.options.ListEmittedFiles.IsTrue() {
		emittedFiles = []string{buildInfoFileName}
	}
	return &compiler.EmitResult{
		EmitSkipped:  false,
		EmittedFiles: emittedFiles,
	}
}

func (p *programState) ensureHasErrorsForState(ctx context.Context, program *compiler.Program) {
	if p.hasErrors != core.TSUnknown {
		return
	}

	// Check semantic and emit diagnostics first as we dont need to ask program about it
	if slices.ContainsFunc(program.GetSourceFiles(), func(file *ast.SourceFile) bool {
		semanticDiagnostics := p.semanticDiagnosticsPerFile[file.Path()]
		if semanticDiagnostics == nil {
			// Missing semantic diagnostics in cache will be encoded in incremental buildInfo
			return p.options.IsIncremental()
		}
		if len(semanticDiagnostics.diagnostics) > 0 || len(semanticDiagnostics.buildInfoDiagnostics) > 0 {
			// cached semantic diagnostics will be encoded in buildInfo
			return true
		}
		if _, ok := p.emitDiagnosticsPerFile[file.Path()]; ok {
			// emit diagnostics will be encoded in buildInfo;
			return true
		}
		return false
	}) {
		// Because semantic diagnostics are recorded in buildInfo, we dont need to encode hasErrors in incremental buildInfo
		// But encode as errors in non incremental buildInfo
		p.hasErrors = core.IfElse(p.options.IsIncremental(), core.TSFalse, core.TSTrue)
		return
	}
	if len(program.GetConfigFileParsingDiagnostics()) > 0 ||
		len(program.GetSyntacticDiagnostics(ctx, nil)) > 0 ||
		len(program.GetBindDiagnostics(ctx, nil)) > 0 ||
		len(program.GetOptionsDiagnostics(ctx)) > 0 {
		p.hasErrors = core.TSTrue
	} else {
		p.hasErrors = core.TSFalse
	}
}

func (p *programState) getSemanticDiagnostics(ctx context.Context, program *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
	if file != nil {
		return p.getSemanticDiagnosticsOfFile(ctx, program, file)
	}

	for {
		_, done := p.getSemanticDiagnosticsOfNextAffectedFile(ctx, program)
		if done {
			break
		}
	}

	var diagnostics []*ast.Diagnostic
	for _, file := range program.GetSourceFiles() {
		diagnostics = append(diagnostics, p.getSemanticDiagnosticsOfFile(ctx, program, file)...)
	}
	if p.checkPending && !p.options.NoCheck.IsTrue() {
		p.checkPending = false
		p.buildInfoEmitPending = true
	}
	return diagnostics
}

/*
* Gets the semantic diagnostics either from cache if present, or otherwise from program and caches it
* Note that it is assumed that when asked about checker diagnostics, the file has been taken out of affected files/changed file set
 */
func (p *programState) getSemanticDiagnosticsOfFile(ctx context.Context, program *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
	if p.options.NoCheck.IsTrue() {
		return nil
	}

	// !!! this is different from strada where we were adding program diagnostics but
	// but with blank slate it would be good to call that directly instead of unnecessarily concatenating

	// Report the check diagnostics from the cache if we already have those diagnostics present
	if cachedDiagnostics, ok := p.semanticDiagnosticsPerFile[file.Path()]; ok {
		return compiler.FilterNoEmitSemanticDiagnostics(cachedDiagnostics.getDiagnostics(program, file), p.options)
	}

	// Diagnostics werent cached, get them from program, and cache the result
	diagnostics := program.GetSemanticDiagnostics(ctx, file)
	p.semanticDiagnosticsPerFile[file.Path()] = &diagnosticsOrBuildInfoDiagnostics{diagnostics: diagnostics}
	p.buildInfoEmitPending = true
	return compiler.FilterNoEmitSemanticDiagnostics(diagnostics, p.options)
}

/**
 * Return the semantic diagnostics for the next affected file, done if iteration is complete
 */
func (p *programState) getSemanticDiagnosticsOfNextAffectedFile(ctx context.Context, program *compiler.Program) ([]*ast.Diagnostic, bool) {
	for {
		affected := p.getNextAffectedFile(ctx, program)
		if affected == nil {
			if p.checkPending && !p.options.NoCheck.IsTrue() {
				p.checkPending = false
				p.buildInfoEmitPending = true
			}
			return nil, true
		}
		// Get diagnostics for the affected file if its not ignored
		result := p.getSemanticDiagnosticsOfFile(ctx, program, affected)
		p.seenAffectedFiles.Add(affected.Path())
		p.affectedFilesIndex++
		p.buildInfoEmitPending = true
		if result == nil {
			continue
		}
		return result, false
	}
}

/**
 * This function returns the next affected file to be processed.
 * Note that until doneAffected is called it would keep reporting same result
 * This is to allow the callers to be able to actually remove affected file only when the operation is complete
 * eg. if during diagnostics check cancellation token ends up cancelling the request, the affected file should be retained
 */
func (p *programState) getNextAffectedFile(ctx context.Context, program *compiler.Program) *ast.SourceFile {
	for {
		if p.affectedFiles != nil {
			for p.affectedFilesIndex < len(p.affectedFiles) {
				affectedFile := p.affectedFiles[p.affectedFilesIndex]
				if !p.seenAffectedFiles.Has(affectedFile.Path()) {
					// Set the next affected file as seen and remove the cached semantic diagnostics
					p.addFileToAffectedFilesPendingEmit(affectedFile.Path(), getFileEmitKind(p.options))
					p.handleDtsMayChangeOfAffectedFile(ctx, program, affectedFile)
					return affectedFile
				}
				p.affectedFilesIndex++
			}

			// Remove the changed file from the change set
			p.changedFilesSet.Delete(p.currentChangedFilePath)
			p.currentChangedFilePath = ""
			// Commit the changes in file signature
			p.oldSignatures = nil
			p.affectedFiles = nil
		}

		// Get next changed file
		var file tspath.Path
		for file = range p.changedFilesSet.Keys() {
			// Get next batch of affected files
			p.affectedFiles = p.getFilesAffectedBy(ctx, program, file)
			p.currentChangedFilePath = file
			p.affectedFilesIndex = 0
			break
		}

		// Done if there are no more changed files
		if file == "" {
			return nil
		}
	}
}

func (p *programState) getFilesAffectedBy(ctx context.Context, program *compiler.Program, path tspath.Path) []*ast.SourceFile {
	file := program.GetSourceFileByPath(path)
	if file == nil {
		return nil
	}

	if !p.updateShapeSignature(ctx, program, file, p.useFileVersionAsSignature) {
		return []*ast.SourceFile{file}
	}

	if !p.tracksReferences() {
		return p.getAllFilesExcludingDefaultLibraryFile(program, file)
	}

	if info := p.fileInfos[file.Path()]; info.affectsGlobalScope {
		p.getAllFilesExcludingDefaultLibraryFile(program, file)
	}

	if p.options.IsolatedModules.IsTrue() {
		return []*ast.SourceFile{file}
	}

	// Now we need to if each file in the referencedBy list has a shape change as well.
	// Because if so, its own referencedBy files need to be saved as well to make the
	// emitting result consistent with files on disk.
	seenFileNamesMap := p.forEachFileReferencedBy(
		program,
		file,
		func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool) {
			// If the current file is not nil and has a shape change, we need to queue it for processing
			if currentFile != nil && p.updateShapeSignature(ctx, program, currentFile, p.useFileVersionAsSignature) {
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

func (p *programState) forEachFileReferencedBy(
	program *compiler.Program,
	file *ast.SourceFile,
	fn func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool),
) map[tspath.Path]*ast.SourceFile {
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
			currentFile := program.GetSourceFileByPath(currentPath)
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

/**
 * Gets the files referenced by the the file path
 */
func (p *programState) getReferencedByPaths(file tspath.Path) map[tspath.Path]struct{} {
	keys, ok := p.referencedMap.GetKeys(file)
	if !ok {
		return nil
	}
	return keys.Keys()
}

func (p *programState) updateShapeSignature(ctx context.Context, program *compiler.Program, file *ast.SourceFile, useFileVersionAsSignature bool) bool {
	// If we have cached the result for this file, that means hence forth we should assume file shape is uptodate
	if p.hasCalledUpdateShapeSignature.Has(file.Path()) {
		return false
	}

	info := p.fileInfos[file.Path()]
	prevSignature := info.signature
	var latestSignature string
	if !file.IsDeclarationFile && !useFileVersionAsSignature {
		latestSignature = p.computeDtsSignature(ctx, program, file)
	}
	// Default is to use file version as signature
	if latestSignature == "" {
		latestSignature = info.version
	}
	if p.oldSignatures == nil {
		p.oldSignatures = make(map[tspath.Path]string)
	}
	p.oldSignatures[file.Path()] = prevSignature
	p.hasCalledUpdateShapeSignature.Add(file.Path())
	info.signature = latestSignature
	return latestSignature != prevSignature
}

func (p *programState) computeDtsSignature(ctx context.Context, program *compiler.Program, file *ast.SourceFile) string {
	var signature string
	program.Emit(ctx, compiler.EmitOptions{
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

func (p *programState) isChangedSignature(path tspath.Path) bool {
	oldSignature := p.oldSignatures[path]
	newSignature := p.fileInfos[path].signature
	return newSignature != oldSignature
}

/**
 * Removes semantic diagnostics for path and
 * returns true if there are no more semantic diagnostics from the old state
 */
func (p *programState) removeSemanticDiagnosticsOf(path tspath.Path) {
	if p.semanticDiagnosticsFromOldState.Has(path) {
		p.semanticDiagnosticsFromOldState.Delete(path)
		delete(p.semanticDiagnosticsPerFile, path)
	}
}

func (p *programState) removeDiagnosticsOfLibraryFiles(program *compiler.Program) {
	if !p.cleanedDiagnosticsOfLibFiles {
		p.cleanedDiagnosticsOfLibFiles = true
		for _, file := range program.GetSourceFiles() {
			if program.IsSourceFileDefaultLibrary(file.Path()) && !checker.SkipTypeChecking(file, p.options, program, true) {
				p.removeSemanticDiagnosticsOf(file.Path())
			}
		}
	}
}

/**
 *  Handles semantic diagnostics and dts emit for affectedFile and files, that are referencing modules that export entities from affected file
 *  This is because even though js emit doesnt change, dts emit / type used can change resulting in need for dts emit and js change
 */
func (p *programState) handleDtsMayChangeOfAffectedFile(ctx context.Context, program *compiler.Program, affectedFile *ast.SourceFile) {
	p.removeSemanticDiagnosticsOf(affectedFile.Path())

	// If affected files is everything except default library, then nothing more to do
	if slices.Equal(p.allFilesExcludingDefaultLibraryFile, p.affectedFiles) {
		p.removeDiagnosticsOfLibraryFiles(program)
		// When a change affects the global scope, all files are considered to be affected without updating their signature
		// That means when affected file is handled, its signature can be out of date
		// To avoid this, ensure that we update the signature for any affected file in this scenario.
		p.updateShapeSignature(ctx, program, affectedFile, p.useFileVersionAsSignature)
		return
	}

	if p.options.AssumeChangesOnlyAffectDirectDependencies.IsTrue() {
		return
	}

	// Iterate on referencing modules that export entities from affected file and delete diagnostics and add pending emit
	// If there was change in signature (dts output) for the changed file,
	// then only we need to handle pending file emit
	if !p.tracksReferences() ||
		!p.changedFilesSet.Has(affectedFile.Path()) ||
		!p.isChangedSignature(affectedFile.Path()) {
		return
	}

	// Since isolated modules dont change js files, files affected by change in signature is itself
	// But we need to cleanup semantic diagnostics and queue dts emit for affected files
	if p.options.IsolatedModules.IsTrue() {
		p.forEachFileReferencedBy(
			program,
			affectedFile,
			func(currentFile *ast.SourceFile, currentPath tspath.Path) (queueForFile bool, fastReturn bool) {
				if p.handleDtsMayChangeOfGlobalScope(ctx, program, currentPath /*invalidateJsFiles*/, false) {
					return false, true
				}
				p.handleDtsMayChangeOf(ctx, program, currentPath /*invalidateJsFiles*/, false)
				if p.isChangedSignature(currentPath) {
					return true, false
				}
				return false, false
			},
		)
	}

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
				typeChecker, done = program.GetTypeCheckerForFile(ctx, affectedFile)
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
	if keys, ok := p.referencedMap.GetKeys(affectedFile.Path()); ok {
		for exportedFromPath := range keys.Keys() {
			if p.handleDtsMayChangeOfGlobalScope(ctx, program, exportedFromPath, invalidateJsFiles) {
				return
			}
			if references, ok := p.referencedMap.GetKeys(exportedFromPath); ok {
				for filePath := range references.Keys() {
					if p.handleDtsMayChangeOfFileAndExportsOfFile(ctx, program, filePath, invalidateJsFiles, &seenFileAndExportsOfFile) {
						return
					}
				}
			}
		}
	}
}

/**
 * handle dts and semantic diagnostics on file and iterate on anything that exports this file
 * return true when all work is done and we can exit handling dts emit and semantic diagnostics
 */
func (p *programState) handleDtsMayChangeOfFileAndExportsOfFile(ctx context.Context, program *compiler.Program, filePath tspath.Path, invalidateJsFiles bool, seenFileAndExportsOfFile *collections.Set[tspath.Path]) bool {
	if seenFileAndExportsOfFile.AddIfAbsent(filePath) == false {
		return false
	}
	if p.handleDtsMayChangeOfGlobalScope(ctx, program, filePath, invalidateJsFiles) {
		return true
	}
	p.handleDtsMayChangeOf(ctx, program, filePath, invalidateJsFiles)

	// Remove the diagnostics of files that import this file and handle all its exports too
	if keys, ok := p.referencedMap.GetKeys(filePath); ok {
		for referencingFilePath := range keys.Keys() {
			if p.handleDtsMayChangeOfFileAndExportsOfFile(ctx, program, referencingFilePath, invalidateJsFiles, seenFileAndExportsOfFile) {
				return true
			}
		}
	}
	return false
}

func (p *programState) handleDtsMayChangeOfGlobalScope(ctx context.Context, program *compiler.Program, filePath tspath.Path, invalidateJsFiles bool) bool {
	if info, ok := p.fileInfos[filePath]; !ok || !info.affectsGlobalScope {
		return false
	}
	// Every file needs to be handled
	for _, file := range p.getAllFilesExcludingDefaultLibraryFile(program, nil) {
		p.handleDtsMayChangeOf(ctx, program, file.Path(), invalidateJsFiles)
	}
	p.removeDiagnosticsOfLibraryFiles(program)
	return true
}

/**
 * Handle the dts may change, so they need to be added to pending emit if dts emit is enabled,
 * Also we need to make sure signature is updated for these files
 */
func (p *programState) handleDtsMayChangeOf(ctx context.Context, program *compiler.Program, path tspath.Path, invalidateJsFiles bool) {
	p.removeSemanticDiagnosticsOf(path)
	if p.changedFilesSet.Has(path) {
		return
	}
	file := program.GetSourceFileByPath(path)
	if file == nil {
		return
	}
	// Even though the js emit doesnt change and we are already handling dts emit and semantic diagnostics
	// we need to update the signature to reflect correctness of the signature(which is output d.ts emit) of this file
	// This ensures that we dont later during incremental builds considering wrong signature.
	// Eg where this also is needed to ensure that .tsbuildinfo generated by incremental build should be same as if it was first fresh build
	// But we avoid expensive full shape computation, as using file version as shape is enough for correctness.
	p.updateShapeSignature(ctx, program, file, true)
	// If not dts emit, nothing more to do
	if invalidateJsFiles {
		p.addFileToAffectedFilesPendingEmit(path, getFileEmitKind(p.options))
	} else if p.options.GetEmitDeclarations() {
		p.addFileToAffectedFilesPendingEmit(path, core.IfElse(p.options.DeclarationMap.IsTrue(), fileEmitKindAllDts, fileEmitKindDts))
	}
}

func newProgramState(program *compiler.Program, oldProgram *Program) *programState {
	if oldProgram != nil && oldProgram.program == program {
		return oldProgram.state
	}
	files := program.GetSourceFiles()
	state := &programState{
		options:                    program.Options(),
		semanticDiagnosticsPerFile: make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics, len(files)),
		seenEmittedFiles:           make(map[tspath.Path]fileEmitKind, len(files)),
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
						state.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics, len(files))
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
			var pendingEmitKind fileEmitKind
			if tsoptions.CompilerOptionsAffectEmit(oldProgram.state.options, state.options) {
				pendingEmitKind = getFileEmitKind(state.options)
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
