package incremental

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/zeebo/xxh3"
)

type FileInfo struct {
	version            string
	signature          string
	affectsGlobalScope bool
	impliedNodeFormat  core.ResolutionMode
}

func (f *FileInfo) Version() string                        { return f.version }
func (f *FileInfo) Signature() string                      { return f.signature }
func (f *FileInfo) AffectsGlobalScope() bool               { return f.affectsGlobalScope }
func (f *FileInfo) ImpliedNodeFormat() core.ResolutionMode { return f.impliedNodeFormat }

func ComputeHash(text string, hashWithText bool) string {
	hashBytes := xxh3.HashString128(text).Bytes()
	hash := hex.EncodeToString(hashBytes[:])
	if hashWithText {
		hash += "-" + text
	}
	return hash
}

type FileEmitKind uint32

const (
	FileEmitKindNone        FileEmitKind = 0
	FileEmitKindJs          FileEmitKind = 1 << 0 // emit js file
	FileEmitKindJsMap       FileEmitKind = 1 << 1 // emit js.map file
	FileEmitKindJsInlineMap FileEmitKind = 1 << 2 // emit inline source map in js file
	FileEmitKindDtsErrors   FileEmitKind = 1 << 3 // emit dts errors
	FileEmitKindDtsEmit     FileEmitKind = 1 << 4 // emit d.ts file
	FileEmitKindDtsMap      FileEmitKind = 1 << 5 // emit d.ts.map file

	FileEmitKindDts        = FileEmitKindDtsErrors | FileEmitKindDtsEmit
	FileEmitKindAllJs      = FileEmitKindJs | FileEmitKindJsMap | FileEmitKindJsInlineMap
	FileEmitKindAllDtsEmit = FileEmitKindDtsEmit | FileEmitKindDtsMap
	FileEmitKindAllDts     = FileEmitKindDts | FileEmitKindDtsMap
	FileEmitKindAll        = FileEmitKindAllJs | FileEmitKindAllDts
)

func GetFileEmitKind(options *core.CompilerOptions) FileEmitKind {
	result := FileEmitKindJs
	if options.SourceMap.IsTrue() {
		result |= FileEmitKindJsMap
	}
	if options.InlineSourceMap.IsTrue() {
		result |= FileEmitKindJsInlineMap
	}
	if options.GetEmitDeclarations() {
		result |= FileEmitKindDts
	}
	if options.DeclarationMap.IsTrue() {
		result |= FileEmitKindDtsMap
	}
	if options.EmitDeclarationOnly.IsTrue() {
		result &= FileEmitKindAllDts
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
		return FileEmitKindNone
	}
	if oldEmitKind == 0 || emitKind == 0 {
		return emitKind
	}
	diff := oldEmitKind ^ emitKind
	result := FileEmitKindNone
	// If there is diff in Js emit, pending emit is js emit flags
	if (diff & FileEmitKindAllJs) != 0 {
		result |= emitKind & FileEmitKindAllJs
	}
	// If dts errors pending, add dts errors flag
	if (diff & FileEmitKindDtsErrors) != 0 {
		result |= emitKind & FileEmitKindAllDts
	}
	// If there is diff in Dts emit, pending emit is dts emit flags
	if (diff & FileEmitKindAllDtsEmit) != 0 {
		result |= emitKind & FileEmitKindAllDtsEmit
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
	messageKey         diagnostics.Key
	messageArgs        []string
	messageChain       []*buildInfoDiagnosticWithFileName
	relatedInformation []*buildInfoDiagnosticWithFileName
	reportsUnnecessary bool
	reportsDeprecated  bool
	skippedOnNoEmit    bool
	repopulateInfo     *ast.RepopulateDiagnosticInfo
}

type DiagnosticsOrBuildInfoDiagnosticsWithFileName struct {
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

	if b.repopulateInfo != nil {
		return repopulateDiagnosticChain(b, p, fileForDiagnostic)
	}

	var messageChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		messageChain = append(messageChain, msg.toDiagnostic(p, fileForDiagnostic))
	}
	var relatedInformation []*ast.Diagnostic
	for _, info := range b.relatedInformation {
		relatedInformation = append(relatedInformation, info.toDiagnostic(p, fileForDiagnostic))
	}
	return ast.NewDiagnosticFromSerialized(
		fileForDiagnostic,
		core.NewTextRange(b.pos, b.end),
		b.code,
		b.category,
		b.messageKey,
		b.messageArgs,
		messageChain,
		relatedInformation,
		b.reportsUnnecessary,
		b.reportsDeprecated,
		b.skippedOnNoEmit,
	)
}

// repopulateDiagnosticChain recomputes a diagnostic chain entry that depends on
// program state which may have changed between incremental builds.
func repopulateDiagnosticChain(b *buildInfoDiagnosticWithFileName, p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	info := b.repopulateInfo
	switch info.Kind {
	case ast.RepopulateModeMismatch:
		return repopulateModeMismatchChain(b, p, file)
	case ast.RepopulateModuleNotFound:
		return repopulateModuleNotFoundChain(b, p, file, info)
	default:
		// Fall back to using the stored (possibly stale) data
		return b.toDiagnosticWithoutRepopulate(p, file)
	}
}

func (b *buildInfoDiagnosticWithFileName) toDiagnosticWithoutRepopulate(p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	var messageChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		messageChain = append(messageChain, msg.toDiagnostic(p, file))
	}
	var relatedInformation []*ast.Diagnostic
	for _, info := range b.relatedInformation {
		relatedInformation = append(relatedInformation, info.toDiagnostic(p, file))
	}
	return ast.NewDiagnosticFromSerialized(
		file,
		core.NewTextRange(b.pos, b.end),
		b.code,
		b.category,
		b.messageKey,
		b.messageArgs,
		messageChain,
		relatedInformation,
		b.reportsUnnecessary,
		b.reportsDeprecated,
		b.skippedOnNoEmit,
	)
}

func repopulateModeMismatchChain(b *buildInfoDiagnosticWithFileName, p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	if file == nil {
		return b.toDiagnosticWithoutRepopulate(p, file)
	}
	ext := tspath.TryGetExtensionFromPath(file.FileName())
	targetExt := core.IfElse(ext == tspath.ExtensionTs, tspath.ExtensionMts, core.IfElse(ext == tspath.ExtensionJs, tspath.ExtensionMjs, ""))
	meta := p.GetSourceFileMetaData(file.Path())
	packageJsonType := meta.PackageJsonType
	packageJsonDirectory := meta.PackageJsonDirectory

	var messageKey diagnostics.Key
	var code int32
	var category diagnostics.Category
	var messageArgs []string

	if packageJsonDirectory != "" && packageJsonType == "" {
		if targetExt != "" {
			messageKey = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_add_the_field_type_Colon_module_to_1.Key()
			code = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_add_the_field_type_Colon_module_to_1.Code()
			category = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_add_the_field_type_Colon_module_to_1.Category()
			messageArgs = []string{targetExt, tspath.CombinePaths(packageJsonDirectory, "package.json")}
		} else {
			messageKey = diagnostics.To_convert_this_file_to_an_ECMAScript_module_add_the_field_type_Colon_module_to_0.Key()
			code = diagnostics.To_convert_this_file_to_an_ECMAScript_module_add_the_field_type_Colon_module_to_0.Code()
			category = diagnostics.To_convert_this_file_to_an_ECMAScript_module_add_the_field_type_Colon_module_to_0.Category()
			messageArgs = []string{tspath.CombinePaths(packageJsonDirectory, "package.json")}
		}
	} else if targetExt != "" {
		messageKey = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_create_a_local_package_json_file_with_type_Colon_module.Key()
		code = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_create_a_local_package_json_file_with_type_Colon_module.Code()
		category = diagnostics.To_convert_this_file_to_an_ECMAScript_module_change_its_file_extension_to_0_or_create_a_local_package_json_file_with_type_Colon_module.Category()
		messageArgs = []string{targetExt}
	} else {
		messageKey = diagnostics.To_convert_this_file_to_an_ECMAScript_module_create_a_local_package_json_file_with_type_Colon_module.Key()
		code = diagnostics.To_convert_this_file_to_an_ECMAScript_module_create_a_local_package_json_file_with_type_Colon_module.Code()
		category = diagnostics.To_convert_this_file_to_an_ECMAScript_module_create_a_local_package_json_file_with_type_Colon_module.Category()
		messageArgs = nil
	}

	var nextChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		nextChain = append(nextChain, msg.toDiagnostic(p, file))
	}

	return ast.NewDiagnosticFromSerialized(
		file,
		core.NewTextRange(b.pos, b.end),
		code,
		category,
		messageKey,
		messageArgs,
		nextChain,
		nil,
		false,
		false,
		false,
	)
}

func repopulateModuleNotFoundChain(b *buildInfoDiagnosticWithFileName, p *compiler.Program, file *ast.SourceFile, info *ast.RepopulateDiagnosticInfo) *ast.Diagnostic {
	if file == nil {
		return b.toDiagnosticWithoutRepopulate(p, file)
	}

	moduleReference := info.ModuleReference
	mode := info.Mode
	packageName := info.PackageName
	if packageName == "" {
		packageName = moduleReference
	}

	resolvedModule := p.GetResolvedModule(file, moduleReference, mode)

	var messageKey diagnostics.Key
	var code int32
	var category diagnostics.Category
	var messageArgs []string

	if resolvedModule != nil && resolvedModule.AlternateResult != "" {
		alternatePackageName := packageName
		if strings.Contains(resolvedModule.AlternateResult, "/node_modules/@types/") {
			alternatePackageName = "@types/" + module.MangleScopedPackageName(packageName)
		}
		messageKey = diagnostics.There_are_types_at_0_but_this_result_could_not_be_resolved_when_respecting_package_json_exports_The_1_library_may_need_to_update_its_package_json_or_typings.Key()
		code = diagnostics.There_are_types_at_0_but_this_result_could_not_be_resolved_when_respecting_package_json_exports_The_1_library_may_need_to_update_its_package_json_or_typings.Code()
		category = diagnostics.There_are_types_at_0_but_this_result_could_not_be_resolved_when_respecting_package_json_exports_The_1_library_may_need_to_update_its_package_json_or_typings.Category()
		messageArgs = []string{resolvedModule.AlternateResult, alternatePackageName}
	} else {
		packagesMap := getPackagesMap(p)
		if _, ok := packagesMap[module.GetTypesPackageName(packageName)]; ok {
			messageKey = diagnostics.If_the_0_package_actually_exposes_this_module_consider_sending_a_pull_request_to_amend_https_Colon_Slash_Slashgithub_com_SlashDefinitelyTyped_SlashDefinitelyTyped_Slashtree_Slashmaster_Slashtypes_Slash_1.Key()
			code = diagnostics.If_the_0_package_actually_exposes_this_module_consider_sending_a_pull_request_to_amend_https_Colon_Slash_Slashgithub_com_SlashDefinitelyTyped_SlashDefinitelyTyped_Slashtree_Slashmaster_Slashtypes_Slash_1.Code()
			category = diagnostics.If_the_0_package_actually_exposes_this_module_consider_sending_a_pull_request_to_amend_https_Colon_Slash_Slashgithub_com_SlashDefinitelyTyped_SlashDefinitelyTyped_Slashtree_Slashmaster_Slashtypes_Slash_1.Category()
			messageArgs = []string{packageName, module.MangleScopedPackageName(packageName)}
		} else if hasTypes, _ := packagesMap[packageName]; hasTypes {
			messageKey = diagnostics.If_the_0_package_actually_exposes_this_module_try_adding_a_new_declaration_d_ts_file_containing_declare_module_1.Key()
			code = diagnostics.If_the_0_package_actually_exposes_this_module_try_adding_a_new_declaration_d_ts_file_containing_declare_module_1.Code()
			category = diagnostics.If_the_0_package_actually_exposes_this_module_try_adding_a_new_declaration_d_ts_file_containing_declare_module_1.Category()
			messageArgs = []string{packageName, moduleReference}
		} else {
			messageKey = diagnostics.Try_npm_i_save_dev_types_Slash_1_if_it_exists_or_add_a_new_declaration_d_ts_file_containing_declare_module_0.Key()
			code = diagnostics.Try_npm_i_save_dev_types_Slash_1_if_it_exists_or_add_a_new_declaration_d_ts_file_containing_declare_module_0.Code()
			category = diagnostics.Try_npm_i_save_dev_types_Slash_1_if_it_exists_or_add_a_new_declaration_d_ts_file_containing_declare_module_0.Category()
			messageArgs = []string{moduleReference, module.MangleScopedPackageName(packageName)}
		}
	}

	var nextChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		nextChain = append(nextChain, msg.toDiagnostic(p, file))
	}

	return ast.NewDiagnosticFromSerialized(
		file,
		core.NewTextRange(b.pos, b.end),
		code,
		category,
		messageKey,
		messageArgs,
		nextChain,
		nil,
		false,
		false,
		false,
	)
}

// getPackagesMap builds a map of package names to whether they bundle types,
// similar to the checker's getPackagesMap.
func getPackagesMap(p *compiler.Program) map[string]bool {
	packagesMap := make(map[string]bool)
	resolvedModules := p.GetResolvedModules()
	for _, resolvedModulesInFile := range resolvedModules {
		for _, mod := range resolvedModulesInFile {
			if mod.PackageId.Name != "" {
				packagesMap[mod.PackageId.Name] = packagesMap[mod.PackageId.Name] || mod.Extension == tspath.ExtensionDts
			}
		}
	}
	return packagesMap
}

func (d *DiagnosticsOrBuildInfoDiagnosticsWithFileName) getDiagnostics(p *compiler.Program, file *ast.SourceFile) []*ast.Diagnostic {
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
	fileInfos collections.SyncMap[tspath.Path, *FileInfo]
	options   *core.CompilerOptions
	//  Contains the map of ReferencedSet=Referenced files of the file if module emit is enabled
	referencedMap referenceMap
	// Cache of semantic diagnostics for files with their Path being the key
	semanticDiagnosticsPerFile collections.SyncMap[tspath.Path, *DiagnosticsOrBuildInfoDiagnosticsWithFileName]
	// Cache of dts emit diagnostics for files with their Path being the key
	emitDiagnosticsPerFile collections.SyncMap[tspath.Path, *DiagnosticsOrBuildInfoDiagnosticsWithFileName]
	// The map has key by source file's path that has been changed
	changedFilesSet collections.SyncSet[tspath.Path]
	// Files pending to be emitted
	affectedFilesPendingEmit collections.SyncMap[tspath.Path, FileEmitKind]
	// Name of the file whose dts was the latest to change
	latestChangedDtsFile string
	// Hash of d.ts emitted for the file, use to track when emit of d.ts changes
	emitSignatures collections.SyncMap[tspath.Path, *emitSignature]
	// Recorded if program had errors that need to be reported even with --noCheck
	hasErrors core.Tristate
	// Recorded if program had semantic errors only for non incremental build
	hasSemanticErrors bool
	// If semantic diagnostic check is pending
	checkPending bool

	// Additional fields that are not serialized but needed to track state

	// true if build info emit is pending
	buildInfoEmitPending                    atomic.Bool
	hasErrorsFromOldState                   core.Tristate
	hasSemanticErrorsFromOldState           bool
	allFilesExcludingDefaultLibraryFileOnce sync.Once
	//  Cache of all files excluding default library file for the current program
	allFilesExcludingDefaultLibraryFile []*ast.SourceFile
	hasChangedDtsFile                   bool
	hasEmitDiagnostics                  bool

	// Used with testing to add text of hash for better comparison
	hashWithText bool
}

func (s *snapshot) addFileToChangeSet(filePath tspath.Path) {
	s.changedFilesSet.Add(filePath)
	s.buildInfoEmitPending.Store(true)
}

func (s *snapshot) addFileToAffectedFilesPendingEmit(filePath tspath.Path, emitKind FileEmitKind) {
	existingKind, _ := s.affectedFilesPendingEmit.Load(filePath)
	s.affectedFilesPendingEmit.Store(filePath, existingKind|emitKind)
	if emitKind&FileEmitKindDtsErrors != 0 {
		s.emitDiagnosticsPerFile.Delete(filePath)
	}
	s.buildInfoEmitPending.Store(true)
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

func getTextHandlingSourceMapForSignature(text string, data *compiler.WriteFileData) string {
	if data.SourceMapUrlPos != -1 {
		return text[:data.SourceMapUrlPos]
	}
	return text
}

func (s *snapshot) computeSignatureWithDiagnostics(file *ast.SourceFile, text string, data *compiler.WriteFileData) string {
	var builder strings.Builder
	builder.WriteString(getTextHandlingSourceMapForSignature(text, data))
	for _, diag := range data.Diagnostics {
		diagnosticToStringBuilder(diag, file, &builder)
	}
	return s.computeHash(builder.String())
}

func diagnosticToStringBuilder(diagnostic *ast.Diagnostic, file *ast.SourceFile, builder *strings.Builder) {
	if diagnostic == nil {
		return
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
	builder.WriteString(string(diagnostic.MessageKey()))
	builder.WriteString("\n")
	for _, arg := range diagnostic.MessageArgs() {
		builder.WriteString(arg)
		builder.WriteString("\n")
	}
	for _, chain := range diagnostic.MessageChain() {
		diagnosticToStringBuilder(chain, file, builder)
	}
	for _, info := range diagnostic.RelatedInformation() {
		diagnosticToStringBuilder(info, file, builder)
	}
}

func (s *snapshot) computeHash(text string) string {
	return ComputeHash(text, s.hashWithText)
}

func (s *snapshot) canUseIncrementalState() bool {
	if !s.options.IsIncremental() && s.options.Build.IsTrue() {
		// If not incremental build (with tsc -b), we don't need to track state except diagnostics per file so we can use it
		return false
	}
	return true
}
