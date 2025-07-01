package incremental

import (
	"context"
	"maps"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
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

type FileEmitKind uint32

const (
	fileEmitKindNone        FileEmitKind = 0
	fileEmitKindJs          FileEmitKind = 1 << 0 // emit js file
	fileEmitKindJsMap       FileEmitKind = 1 << 1 // emit js.map file
	fileEmitKindJsInlineMap FileEmitKind = 1 << 2 // emit inline source map in js file
	fileEmitKindDtsErrors   FileEmitKind = 1 << 3 // emit dts errors
	fileEmitKindDtsEmit     FileEmitKind = 1 << 4 // emit d.ts file
	fileEmitKindDtsMap      FileEmitKind = 1 << 5 // emit d.ts.map file

	fileEmitKindDts        = fileEmitKindDtsErrors | fileEmitKindDtsEmit
	fileEmitKindAllJs      = fileEmitKindJs | fileEmitKindJsMap | fileEmitKindJsInlineMap
	fileEmitKindAllDtsEmit = fileEmitKindDtsEmit | fileEmitKindDtsMap
	fileEmitKindAllDts     = fileEmitKindDts | fileEmitKindDtsMap
	fileEmitKindAll        = fileEmitKindAllJs | fileEmitKindAllDts
)

func (fileEmitKind FileEmitKind) String() string {
	var builder strings.Builder
	addFlags := func(flags string) {
		if builder.Len() == 0 {
			builder.WriteString(flags)
		} else {
			builder.WriteString("|")
			builder.WriteString(flags)
		}
	}
	if fileEmitKind != 0 {
		if (fileEmitKind & fileEmitKindJs) != 0 {
			addFlags("Js")
		}
		if (fileEmitKind & fileEmitKindJsMap) != 0 {
			addFlags("JsMap")
		}
		if (fileEmitKind & fileEmitKindJsInlineMap) != 0 {
			addFlags("JsInlineMap")
		}
		if (fileEmitKind & fileEmitKindDts) == fileEmitKindDts {
			addFlags("Dts")
		} else {
			if (fileEmitKind & fileEmitKindDtsEmit) != 0 {
				addFlags("DtsEmit")
			}
			if (fileEmitKind & fileEmitKindDtsErrors) != 0 {
				addFlags("DtsErrors")
			}
		}
		if (fileEmitKind & fileEmitKindDtsMap) != 0 {
			addFlags("DtsMap")
		}
	}
	if builder.Len() != 0 {
		return builder.String()
	}
	return "None"
}

func GetFileEmitKind(options *core.CompilerOptions) FileEmitKind {
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

func getPendingEmitKindWithOptions(options *core.CompilerOptions, oldOptions *core.CompilerOptions) FileEmitKind {
	oldEmitKind := GetFileEmitKind(oldOptions)
	newEmitKind := GetFileEmitKind(options)
	return getPendingEmitKind(newEmitKind, oldEmitKind)
}

func getPendingEmitKind(emitKind FileEmitKind, oldEmitKind FileEmitKind) FileEmitKind {
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

// Signature (Hash of d.ts emitted), is string if it was emitted using same d.ts.map option as what compilerOptions indicate,
// otherwise tuple of string
type emitSignature struct {
	signature                     string
	signatureWithDifferentOptions []string
}

// Covert to Emit signature based on oldOptions and EmitSignature format
// If d.ts map options differ then swap the format, otherwise use as is
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

type buildInfoDiagnosticWithFileName struct {
	// filename if it is for a File thats other than its stored for
	file               tspath.Path
	noFile             bool
	pos                int
	end                int
	code               int32
	category           diagnostics.Category
	message            string
	messageChain       []*buildInfoDiagnosticWithFileName
	relatedInformation []*buildInfoDiagnosticWithFileName
	reportsUnnecessary bool
	reportsDeprecated  bool
	skippedOnNoEmit    bool
}

type diagnosticsOrBuildInfoDiagnosticsWithFileName struct {
	diagnostics          []*ast.Diagnostic
	buildInfoDiagnostics []*buildInfoDiagnosticWithFileName
}

func (b *buildInfoDiagnosticWithFileName) toDiagnostic(p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	var fileForDiagnostic *ast.SourceFile
	if b.file != "" {
		fileForDiagnostic = p.GetSourceFileByPath(b.file)
	} else if !b.noFile {
		fileForDiagnostic = file
	}
	var messageChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		messageChain = append(messageChain, msg.toDiagnostic(p, fileForDiagnostic))
	}
	var relatedInformation []*ast.Diagnostic
	for _, info := range b.relatedInformation {
		relatedInformation = append(relatedInformation, info.toDiagnostic(p, fileForDiagnostic))
	}
	return ast.NewDiagnosticWith(
		fileForDiagnostic,
		core.NewTextRange(b.pos, b.end),
		b.code,
		b.category,
		b.message,
		messageChain,
		relatedInformation,
		b.reportsUnnecessary,
		b.reportsDeprecated,
		b.skippedOnNoEmit,
	)
}

func (d *diagnosticsOrBuildInfoDiagnosticsWithFileName) getDiagnostics(p *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
	if d.diagnostics != nil {
		return d.diagnostics
	}
	// Convert and cache the diagnostics
	d.diagnostics = core.Map(d.buildInfoDiagnostics, func(diag *buildInfoDiagnosticWithFileName) *ast.Diagnostic {
		return diag.toDiagnostic(p, file)
	})
	return d.diagnostics
}

type snapshot struct {
	// These are the fields that get serialized

	// Information of the file eg. its version, signature etc
	fileInfos map[tspath.Path]*fileInfo
	options   *core.CompilerOptions
	//  Contains the map of ReferencedSet=Referenced files of the file if module emit is enabled
	referencedMap *collections.ManyToManyMap[tspath.Path, tspath.Path]
	// Cache of semantic diagnostics for files with their Path being the key
	semanticDiagnosticsPerFile map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName
	// Cache of dts emit diagnostics for files with their Path being the key
	emitDiagnosticsPerFile map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName
	// The map has key by source file's path that has been changed
	changedFilesSet *collections.Set[tspath.Path]
	// Files pending to be emitted
	affectedFilesPendingEmit map[tspath.Path]FileEmitKind
	// Name of the file whose dts was the latest to change
	latestChangedDtsFile string
	// Hash of d.ts emitted for the file, use to track when emit of d.ts changes
	emitSignatures map[tspath.Path]*emitSignature
	// Recorded if program had errors
	hasErrors core.Tristate
	// If semantic diagnsotic check is pending
	checkPending bool

	// Additional fields that are not serialized but needed to track state

	// true if build info emit is pending
	buildInfoEmitPending                    bool
	allFilesExcludingDefaultLibraryFileOnce sync.Once
	//  Cache of all files excluding default library file for the current program
	allFilesExcludingDefaultLibraryFile []*ast.SourceFile
}

func (s *snapshot) tracksReferences() bool {
	return s.options.Module != core.ModuleKindNone
}

func (s *snapshot) createReferenceMap() {
	if s.tracksReferences() {
		s.referencedMap = &collections.ManyToManyMap[tspath.Path, tspath.Path]{}
	}
}

func (s *snapshot) createEmitSignaturesMap() {
	if s.emitSignatures == nil && s.options.Composite.IsTrue() {
		s.emitSignatures = make(map[tspath.Path]*emitSignature)
	}
}

func (s *snapshot) addFileToChangeSet(filePath tspath.Path) {
	s.changedFilesSet.Add(filePath)
	s.buildInfoEmitPending = true
}

func (s *snapshot) addFileToAffectedFilesPendingEmit(filePath tspath.Path, emitKind FileEmitKind) {
	existingKind := s.affectedFilesPendingEmit[filePath]
	if s.affectedFilesPendingEmit == nil {
		s.affectedFilesPendingEmit = make(map[tspath.Path]FileEmitKind)
	}
	s.affectedFilesPendingEmit[filePath] = existingKind | emitKind
	delete(s.emitDiagnosticsPerFile, filePath)
}

func (s *snapshot) getAllFilesExcludingDefaultLibraryFile(program *compiler.Program, firstSourceFile *ast.SourceFile) []*ast.SourceFile {
	s.allFilesExcludingDefaultLibraryFileOnce.Do(func() {
		files := program.GetSourceFiles()
		s.allFilesExcludingDefaultLibraryFile = make([]*ast.SourceFile, 0, len(files))
		addSourceFile := func(file *ast.SourceFile) {
			if !program.IsSourceFileDefaultLibrary(file.Path()) {
				s.allFilesExcludingDefaultLibraryFile = append(s.allFilesExcludingDefaultLibraryFile, file)
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
	})
	return s.allFilesExcludingDefaultLibraryFile
}

func newSnapshotForProgram(program *compiler.Program, oldProgram *Program) *snapshot {
	if oldProgram != nil && oldProgram.program == program {
		return oldProgram.snapshot
	}
	files := program.GetSourceFiles()
	snapshot := &snapshot{
		options:                    program.Options(),
		semanticDiagnosticsPerFile: make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(files)),
	}
	snapshot.createReferenceMap()
	if oldProgram != nil && snapshot.options.Composite.IsTrue() {
		snapshot.latestChangedDtsFile = oldProgram.snapshot.latestChangedDtsFile
	}
	if snapshot.options.NoCheck.IsTrue() {
		snapshot.checkPending = true
	}

	canUseStateFromOldProgram := oldProgram != nil && snapshot.tracksReferences() == oldProgram.snapshot.tracksReferences()
	if canUseStateFromOldProgram {
		// Copy old snapshot's changed files set
		snapshot.changedFilesSet = oldProgram.snapshot.changedFilesSet.Clone()
		if len(oldProgram.snapshot.affectedFilesPendingEmit) != 0 {
			snapshot.affectedFilesPendingEmit = maps.Clone(oldProgram.snapshot.affectedFilesPendingEmit)
		}
		snapshot.buildInfoEmitPending = oldProgram.snapshot.buildInfoEmitPending
	} else {
		snapshot.changedFilesSet = &collections.Set[tspath.Path]{}
		snapshot.buildInfoEmitPending = snapshot.options.IsIncremental()
	}

	canCopySemanticDiagnostics := canUseStateFromOldProgram &&
		!tsoptions.CompilerOptionsAffectSemanticDiagnostics(oldProgram.snapshot.options, program.Options())
	// We can only reuse emit signatures (i.e. .d.ts signatures) if the .d.ts file is unchanged,
	// which will eg be depedent on change in options like declarationDir and outDir options are unchanged.
	// We need to look in oldState.compilerOptions, rather than oldCompilerOptions (i.e.we need to disregard useOldState) because
	// oldCompilerOptions can be undefined if there was change in say module from None to some other option
	// which would make useOldState as false since we can now use reference maps that are needed to track what to emit, what to check etc
	// but that option change does not affect d.ts file name so emitSignatures should still be reused.
	canCopyEmitSignatures := snapshot.options.Composite.IsTrue() &&
		oldProgram != nil &&
		oldProgram.snapshot.emitSignatures != nil &&
		!tsoptions.CompilerOptionsAffectDeclarationPath(oldProgram.snapshot.options, program.Options())
	copyDeclarationFileDiagnostics := canCopySemanticDiagnostics &&
		snapshot.options.SkipLibCheck.IsTrue() == oldProgram.snapshot.options.SkipLibCheck.IsTrue()
	copyLibFileDiagnostics := copyDeclarationFileDiagnostics &&
		snapshot.options.SkipDefaultLibCheck.IsTrue() == oldProgram.snapshot.options.SkipDefaultLibCheck.IsTrue()
	snapshot.fileInfos = make(map[tspath.Path]*fileInfo, len(files))
	for _, file := range files {
		version := computeHash(file.Text())
		impliedNodeFormat := program.GetSourceFileMetaData(file.Path()).ImpliedNodeFormat
		affectsGlobalScope := fileAffectsGlobalScope(file)
		var signature string
		var newReferences *collections.Set[tspath.Path]
		if snapshot.referencedMap != nil {
			newReferences = getReferencedFiles(program, file)
			if newReferences != nil {
				snapshot.referencedMap.Add(file.Path(), newReferences)
			}
		}
		if canUseStateFromOldProgram {
			if oldFileInfo, ok := oldProgram.snapshot.fileInfos[file.Path()]; ok {
				signature = oldFileInfo.signature
				if oldFileInfo.version == version || oldFileInfo.affectsGlobalScope != affectsGlobalScope || oldFileInfo.impliedNodeFormat != impliedNodeFormat {
					snapshot.addFileToChangeSet(file.Path())
				} else if snapshot.referencedMap != nil {
					oldReferences, _ := oldProgram.snapshot.referencedMap.GetValues(file.Path())
					// Referenced files changed
					if !newReferences.Equals(oldReferences) {
						snapshot.addFileToChangeSet(file.Path())
					} else {
						for refPath := range newReferences.Keys() {
							if program.GetSourceFileByPath(refPath) == nil {
								// Referenced file was deleted in the new program
								snapshot.addFileToChangeSet(file.Path())
								break
							}
						}
					}
				}
			} else {
				snapshot.addFileToChangeSet(file.Path())
			}
			if !snapshot.changedFilesSet.Has(file.Path()) {
				if emitDiagnostics, ok := oldProgram.snapshot.emitDiagnosticsPerFile[file.Path()]; ok {
					if snapshot.emitDiagnosticsPerFile == nil {
						snapshot.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(files))
					}
					snapshot.emitDiagnosticsPerFile[file.Path()] = emitDiagnostics
				}
				if canCopySemanticDiagnostics {
					if (!file.IsDeclarationFile || copyDeclarationFileDiagnostics) &&
						(!program.IsSourceFileDefaultLibrary(file.Path()) || copyLibFileDiagnostics) {
						// Unchanged file copy diagnostics
						if diagnostics, ok := oldProgram.snapshot.semanticDiagnosticsPerFile[file.Path()]; ok {
							snapshot.semanticDiagnosticsPerFile[file.Path()] = diagnostics
						}
					}
				}
			}
			if canCopyEmitSignatures {
				if oldEmitSignature, ok := oldProgram.snapshot.emitSignatures[file.Path()]; ok {
					snapshot.createEmitSignaturesMap()
					snapshot.emitSignatures[file.Path()] = oldEmitSignature.getNewEmitSignature(oldProgram.snapshot.options, snapshot.options)
				}
			}
		} else {
			snapshot.addFileToAffectedFilesPendingEmit(file.Path(), GetFileEmitKind(snapshot.options))
			signature = version
		}
		snapshot.fileInfos[file.Path()] = &fileInfo{
			version:            version,
			signature:          signature,
			affectsGlobalScope: affectsGlobalScope,
			impliedNodeFormat:  impliedNodeFormat,
		}
	}
	if canUseStateFromOldProgram {
		// If the global file is removed, add all files as changed
		allFilesExcludingDefaultLibraryFileAddedToChangeSet := false
		for filePath, oldInfo := range oldProgram.snapshot.fileInfos {
			if _, ok := snapshot.fileInfos[filePath]; !ok {
				if oldInfo.affectsGlobalScope {
					for _, file := range snapshot.getAllFilesExcludingDefaultLibraryFile(program, nil) {
						snapshot.addFileToChangeSet(file.Path())
					}
					allFilesExcludingDefaultLibraryFileAddedToChangeSet = true
				} else {
					snapshot.buildInfoEmitPending = true
				}
				break
			}
		}
		if !allFilesExcludingDefaultLibraryFileAddedToChangeSet {
			// If options affect emit, then we need to do complete emit per compiler options
			// otherwise only the js or dts that needs to emitted because its different from previously emitted options
			var pendingEmitKind FileEmitKind
			if tsoptions.CompilerOptionsAffectEmit(oldProgram.snapshot.options, snapshot.options) {
				pendingEmitKind = GetFileEmitKind(snapshot.options)
			} else {
				pendingEmitKind = getPendingEmitKindWithOptions(snapshot.options, oldProgram.snapshot.options)
			}
			if pendingEmitKind != fileEmitKindNone {
				// Add all files to affectedFilesPendingEmit since emit changed
				for _, file := range files {
					// Add to affectedFilesPending emit only if not changed since any changed file will do full emit
					if !snapshot.changedFilesSet.Has(file.Path()) {
						snapshot.addFileToAffectedFilesPendingEmit(file.Path(), pendingEmitKind)
					}
				}
				snapshot.buildInfoEmitPending = true
			}
		}
	}
	if canUseStateFromOldProgram &&
		len(snapshot.semanticDiagnosticsPerFile) != len(snapshot.fileInfos) &&
		oldProgram.snapshot.checkPending != snapshot.checkPending {
		snapshot.buildInfoEmitPending = true
	}
	return snapshot
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

	// For script files that contains only ambient external modules, although they are not actually external module files,
	// they can only be consumed via importing elements from them. Regular script files cannot consume them. Therefore,
	// there are no point to rebuild all script files if these special files have changed. However, if any statement
	// in the file is not ambient external module, we treat it as a regular script file.
	return file.Statements != nil &&
		file.Statements.Nodes != nil &&
		core.Some(file.Statements.Nodes, func(stmt *ast.Node) bool {
			return !ast.IsModuleWithStringLiteralName(stmt)
		})
}

func addReferencedFilesFromSymbol(file *ast.SourceFile, referencedFiles *collections.Set[tspath.Path], symbol *ast.Symbol) {
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

// Get the module source file and all augmenting files from the import name node from file
func addReferencedFilesFromImportLiteral(file *ast.SourceFile, referencedFiles *collections.Set[tspath.Path], checker *checker.Checker, importName *ast.LiteralLikeNode) {
	symbol := checker.GetSymbolAtLocation(importName)
	addReferencedFilesFromSymbol(file, referencedFiles, symbol)
}

// Gets the path to reference file from file name, it could be resolvedPath if present otherwise path
func addReferencedFileFromFileName(program *compiler.Program, fileName string, referencedFiles *collections.Set[tspath.Path], sourceFileDirectory string) {
	if redirect := program.GetParseFileRedirect(fileName); redirect != "" {
		referencedFiles.Add(tspath.ToPath(redirect, program.GetCurrentDirectory(), program.UseCaseSensitiveFileNames()))
	} else {
		referencedFiles.Add(tspath.ToPath(fileName, sourceFileDirectory, program.UseCaseSensitiveFileNames()))
	}
}

// Gets the referenced files for a file from the program with values for the keys as referenced file's path to be true
func getReferencedFiles(program *compiler.Program, file *ast.SourceFile) *collections.Set[tspath.Path] {
	referencedFiles := collections.Set[tspath.Path]{}

	// We need to use a set here since the code can contain the same import twice,
	// but that will only be one dependency.
	// To avoid invernal conversion, the key of the referencedFiles map must be of type Path
	checker, done := program.GetTypeCheckerForFile(context.TODO(), file)
	defer done()
	for _, importName := range file.Imports() {
		addReferencedFilesFromImportLiteral(file, &referencedFiles, checker, importName)
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

	// Add module augmentation as references
	for _, moduleName := range file.ModuleAugmentations {
		if !ast.IsStringLiteral(moduleName) {
			continue
		}
		addReferencedFilesFromImportLiteral(file, &referencedFiles, checker, moduleName)
	}

	// From ambient modules
	for _, ambientModule := range checker.GetAmbientModules() {
		addReferencedFilesFromSymbol(file, &referencedFiles, ambientModule)
	}
	return core.IfElse(referencedFiles.Len() > 0, &referencedFiles, nil)
}
