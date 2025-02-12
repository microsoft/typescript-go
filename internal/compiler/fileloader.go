package compiler

import (
	"cmp"
	"iter"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/module"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type fileLoader struct {
	host            CompilerHost
	programOptions  ProgramOptions
	compilerOptions *core.CompilerOptions

	resolver             *module.Resolver
	resolvedModulesMutex sync.Mutex
	resolvedModules      map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule]

	mu                      sync.Mutex
	wg                      *core.WorkGroup
	tasksByFileName         map[string]*parseTask
	currentNodeModulesDepth int
	defaultLibraryPath      string
	comparePathsOptions     tspath.ComparePathsOptions
	rootTasks               []*parseTask
}

type sourceFileMetadata struct {
	file              *ast.SourceFile
	impliedNodeFormat core.ResolutionMode
	packageJsonScope  *packagejson.InfoCacheEntry
}

func processAllProgramFiles(
	host CompilerHost,
	programOptions ProgramOptions,
	compilerOptions *core.CompilerOptions,
	resolver *module.Resolver,
	rootFiles []string,
	libs []string,
) (files []*ast.SourceFile, resolvedModules map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule]) {
	loader := fileLoader{
		host:               host,
		programOptions:     programOptions,
		compilerOptions:    compilerOptions,
		resolver:           resolver,
		tasksByFileName:    make(map[string]*parseTask),
		defaultLibraryPath: host.DefaultLibraryPath(),
		comparePathsOptions: tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		},
		wg:        core.NewWorkGroup(programOptions.SingleThreaded),
		rootTasks: make([]*parseTask, 0, len(rootFiles)+len(libs)),
	}

	loader.addRootTasks(rootFiles, false)
	loader.addRootTasks(libs, true)
	loader.addAutomaticTypeDirectiveTasks()

	loader.startTasks(loader.rootTasks)

	loader.wg.Wait()

	files, libFiles := []*ast.SourceFile{}, []*ast.SourceFile{}
	for task := range loader.collectTasks(loader.rootTasks) {
		if task.isLib {
			libFiles = append(libFiles, task.file)
		} else {
			files = append(files, task.file)
		}
	}
	loader.sortLibs(libFiles)

	return append(libFiles, files...), loader.resolvedModules
}

func (p *fileLoader) addRootTasks(files []string, isLib bool) {
	for _, fileName := range files {
		absPath := tspath.GetNormalizedAbsolutePath(fileName, p.host.GetCurrentDirectory())
		p.rootTasks = append(p.rootTasks, &parseTask{normalizedFilePath: absPath, isLib: isLib})
	}
}

func (p *fileLoader) addAutomaticTypeDirectiveTasks() {
	var containingDirectory string
	if p.compilerOptions.ConfigFilePath != "" {
		containingDirectory = tspath.GetDirectoryPath(p.compilerOptions.ConfigFilePath)
	} else {
		containingDirectory = p.host.GetCurrentDirectory()
	}
	containingFileName := tspath.CombinePaths(containingDirectory, module.InferredTypesContainingFile)

	automaticTypeDirectiveNames := module.GetAutomaticTypeDirectiveNames(p.compilerOptions, p.host)
	for _, name := range automaticTypeDirectiveNames {
		resolved := p.resolver.ResolveTypeReferenceDirective(name, containingFileName, core.ResolutionModeNone, nil)
		if resolved.IsResolved() {
			p.rootTasks = append(p.rootTasks, &parseTask{normalizedFilePath: resolved.ResolvedFileName, isLib: false})
		}
	}
}

func (p *fileLoader) startTasks(tasks []*parseTask) {
	if len(tasks) > 0 {
		p.mu.Lock()
		defer p.mu.Unlock()
		for i, task := range tasks {
			// dedup tasks to ensure correct file order, regardless of which task would be started first
			if existingTask, ok := p.tasksByFileName[task.normalizedFilePath]; ok {
				tasks[i] = existingTask
			} else {
				p.tasksByFileName[task.normalizedFilePath] = task
				task.start(p)
			}
		}
	}
}

func (p *fileLoader) collectTasks(tasks []*parseTask) iter.Seq[*parseTask] {
	return func(yield func(*parseTask) bool) {
		p.collectTasksWorker(tasks, yield)
	}
}

func (p *fileLoader) collectTasksWorker(tasks []*parseTask, yield func(*parseTask) bool) bool {
	for _, task := range tasks {
		if _, ok := p.tasksByFileName[task.normalizedFilePath]; ok {
			// ensure we only walk each task once
			delete(p.tasksByFileName, task.normalizedFilePath)

			if len(task.subTasks) > 0 {
				if !p.collectTasksWorker(task.subTasks, yield) {
					return false
				}
			}

			if task.file != nil {
				if !yield(task) {
					return false
				}
			}
		}
	}
	return true
}

func (p *fileLoader) sortLibs(libFiles []*ast.SourceFile) {
	slices.SortFunc(libFiles, func(f1 *ast.SourceFile, f2 *ast.SourceFile) int {
		return cmp.Compare(p.getDefaultLibFilePriority(f1), p.getDefaultLibFilePriority(f2))
	})
}

func (p *fileLoader) getDefaultLibFilePriority(a *ast.SourceFile) int {
	if tspath.ContainsPath(p.defaultLibraryPath, a.FileName(), p.comparePathsOptions) {
		basename := tspath.GetBaseFileName(a.FileName())
		if basename == "lib.d.ts" || basename == "lib.es6.d.ts" {
			return 0
		}
		name := strings.TrimSuffix(strings.TrimPrefix(basename, "lib."), ".d.ts")
		index := slices.Index(tsoptions.Libs, name)
		if index != -1 {
			return index + 1
		}
	}
	return len(tsoptions.Libs) + 2
}

type parseTask struct {
	normalizedFilePath string
	file               *ast.SourceFile
	isLib              bool
	subTasks           []*parseTask
}

func (t *parseTask) start(loader *fileLoader) {
	loader.wg.Run(func() {
		file := loader.parseSourceFile(t.normalizedFilePath)
		metadata := loader.buildFileMetadata(file)

		// !!! if noResolve, skip all of this
		loader.collectExternalModuleReferences(file)

		t.subTasks = make([]*parseTask, 0, len(file.ReferencedFiles)+len(file.Imports)+len(file.ModuleAugmentations))

		for _, ref := range file.ReferencedFiles {
			resolvedPath := resolveTripleslashPathReference(ref.FileName, file.FileName())
			t.addSubTask(resolvedPath, false)
		}

		for _, ref := range file.TypeReferenceDirectives {
			resolutionMode := getModeForTypeReferenceDirectiveInFile(ref, metadata, loader.compilerOptions)
			resolved := loader.resolver.ResolveTypeReferenceDirective(ref.FileName, file.FileName(), resolutionMode, nil)
			if resolved.IsResolved() {
				t.addSubTask(resolved.ResolvedFileName, false)
			}
		}

		if loader.compilerOptions.NoLib != core.TSTrue {
			for _, lib := range file.LibReferenceDirectives {
				name, ok := tsoptions.GetLibFileName(lib.FileName)
				if !ok {
					continue
				}
				t.addSubTask(tspath.CombinePaths(loader.defaultLibraryPath, name), true)
			}
		}

		for _, imp := range loader.resolveImportsAndModuleAugmentations(metadata) {
			t.addSubTask(imp, false)
		}

		t.file = file
		loader.startTasks(t.subTasks)
	})
}

func (p *fileLoader) parseSourceFile(fileName string) *ast.SourceFile {
	path := tspath.ToPath(fileName, p.host.GetCurrentDirectory(), p.host.FS().UseCaseSensitiveFileNames())
	sourceFile := p.host.GetSourceFile(fileName, p.compilerOptions.GetEmitScriptTarget())
	sourceFile.SetPath(path)
	return sourceFile
}

func (t *parseTask) addSubTask(fileName string, isLib bool) {
	normalizedFilePath := tspath.NormalizePath(fileName)
	t.subTasks = append(t.subTasks, &parseTask{normalizedFilePath: normalizedFilePath, isLib: isLib})
}

func (p *fileLoader) buildFileMetadata(sourceFile *ast.SourceFile) *sourceFileMetadata {
	fileMetadata := &sourceFileMetadata{file: sourceFile}
	fileName := sourceFile.FileName()
	moduleResolution := p.compilerOptions.GetModuleResolutionKind()
	shouldLookUpFromPackageJson := core.ModuleResolutionKindNode16 <= moduleResolution && moduleResolution <= core.ModuleResolutionKindNodeNext || tspath.ContainsPath(fileName, "/node_modules/", p.comparePathsOptions)

	if hasEsModuleExtension(fileName) {
		fileMetadata.impliedNodeFormat = core.ResolutionModeESM
	} else if hasCommonJsExtension(fileName) {
		fileMetadata.impliedNodeFormat = core.ResolutionModeCommonJS
	} else if shouldLookUpFromPackageJson && hasStandardExtension(fileName) {
		pkgJsonScope := p.resolver.GetPackageScopeForPath(tspath.GetDirectoryPath(fileName))

		if pkgJsonScope != nil {
			if pkgJsonScope.Contents.Type.Value == "module" {
				fileMetadata.impliedNodeFormat = core.ResolutionModeESM
			} else {
				fileMetadata.impliedNodeFormat = core.ResolutionModeCommonJS
			}
		}

		fileMetadata.packageJsonScope = pkgJsonScope
	}

	return fileMetadata
}

func hasEsModuleExtension(fileName string) bool {
	return tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionDmts, tspath.ExtensionMts, tspath.ExtensionMjs})
}

func hasCommonJsExtension(fileName string) bool {
	return tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionCjs, tspath.ExtensionCts, tspath.ExtensionDcts})
}

func hasStandardExtension(fileName string) bool {
	return tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx})
}

func (p *fileLoader) collectExternalModuleReferences(file *ast.SourceFile) {
	if file.ModuleReferencesProcessed {
		return
	}
	file.ModuleReferencesProcessed = true
	// !!!
	// If we are importing helpers, we need to add a synthetic reference to resolve the
	// helpers library. (A JavaScript file without `externalModuleIndicator` set might be
	// a CommonJS module; `commonJsModuleIndicator` doesn't get set until the binder has
	// run. We synthesize a helpers import for it just in case; it will never be used if
	// the binder doesn't find and set a `commonJsModuleIndicator`.)
	// if (isJavaScriptFile || (!file.isDeclarationFile && (getIsolatedModules(options) || isExternalModule(file)))) {
	// 	if (options.importHelpers) {
	// 		// synthesize 'import "tslib"' declaration
	// 		imports = [createSyntheticImport(externalHelpersModuleNameText, file)];
	// 	}
	// 	const jsxImport = getJSXRuntimeImport(getJSXImplicitImportBase(options, file), options);
	// 	if (jsxImport) {
	// 		// synthesize `import "base/jsx-runtime"` declaration
	// 		(imports ||= []).push(createSyntheticImport(jsxImport, file));
	// 	}
	// }
	for _, node := range file.Statements.Nodes {
		p.collectModuleReferences(file, node, false /*inAmbientModule*/)
	}

	if file.Flags&ast.NodeFlagsPossiblyContainsDynamicImport != 0 {
		file.CollectDynamicImportOrRequireOrJsDocImportCalls()
	}
}

func (p *fileLoader) collectModuleReferences(file *ast.SourceFile, node *ast.Statement, inAmbientModule bool) {
	if ast.IsAnyImportOrReExport(node) {
		moduleNameExpr := ast.GetExternalModuleName(node)
		// TypeScript 1.0 spec (April 2014): 12.1.6
		// An ExternalImportDeclaration in an AmbientExternalModuleDeclaration may reference other external modules
		// only through top - level external module names. Relative external module names are not permitted.
		if moduleNameExpr != nil && ast.IsStringLiteral(moduleNameExpr) {
			moduleName := moduleNameExpr.AsStringLiteral().Text
			if moduleName != "" && (!inAmbientModule || !tspath.IsExternalModuleNameRelative(moduleName)) {
				ast.SetParentInChildren(node) // we need parent data on imports before the program is fully bound, so we ensure it's set here
				file.Imports = append(file.Imports, moduleNameExpr)
				if file.UsesUriStyleNodeCoreModules != core.TSTrue && p.currentNodeModulesDepth == 0 && !file.IsDeclarationFile {
					if strings.HasPrefix(moduleName, "node:") && !exclusivelyPrefixedNodeCoreModules[moduleName] {
						// Presence of `node:` prefix takes precedence over unprefixed node core modules
						file.UsesUriStyleNodeCoreModules = core.TSTrue
					} else if file.UsesUriStyleNodeCoreModules == core.TSUnknown && unprefixedNodeCoreModules[moduleName] {
						// Avoid `unprefixedNodeCoreModules.has` for every import
						file.UsesUriStyleNodeCoreModules = core.TSFalse
					}
				}
			}
		}
		return
	}
	if ast.IsModuleDeclaration(node) && ast.IsAmbientModule(node) && (inAmbientModule || ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) || file.IsDeclarationFile) {
		ast.SetParentInChildren(node)
		nameText := node.AsModuleDeclaration().Name().Text()
		// Ambient module declarations can be interpreted as augmentations for some existing external modules.
		// This will happen in two cases:
		// - if current file is external module then module augmentation is a ambient module declaration defined in the top level scope
		// - if current file is not external module then module augmentation is an ambient module declaration with non-relative module name
		//   immediately nested in top level ambient module declaration .
		if ast.IsExternalModule(file) || (inAmbientModule && !tspath.IsExternalModuleNameRelative(nameText)) {
			file.ModuleAugmentations = append(file.ModuleAugmentations, node.AsModuleDeclaration().Name())
		} else if !inAmbientModule {
			if file.IsDeclarationFile {
				// for global .d.ts files record name of ambient module
				file.AmbientModuleNames = append(file.AmbientModuleNames, nameText)
			}
			// An AmbientExternalModuleDeclaration declares an external module.
			// This type of declaration is permitted only in the global module.
			// The StringLiteral must specify a top - level external module name.
			// Relative external module names are not permitted
			// NOTE: body of ambient module is always a module block, if it exists
			if node.AsModuleDeclaration().Body != nil {
				for _, statement := range node.AsModuleDeclaration().Body.AsModuleBlock().Statements.Nodes {
					p.collectModuleReferences(file, statement, true /*inAmbientModule*/)
				}
			}
		}
	}
}

func resolveTripleslashPathReference(moduleName string, containingFile string) string {
	basePath := tspath.GetDirectoryPath(containingFile)
	referencedFileName := moduleName

	if !tspath.IsRootedDiskPath(moduleName) {
		referencedFileName = tspath.CombinePaths(basePath, moduleName)
	}
	return tspath.NormalizePath(referencedFileName)
}

func (p *fileLoader) resolveImportsAndModuleAugmentations(item *sourceFileMetadata) []string {
	file := item.file
	toParse := make([]string, 0, len(file.Imports))
	if len(file.Imports) > 0 || len(file.ModuleAugmentations) > 0 {
		moduleNames := getModuleNames(file)
		resolutions := p.resolveModuleNames(moduleNames, item)
		optionsForFile := p.getCompilerOptionsForFile(file)
		resolutionsInFile := make(module.ModeAwareCache[*module.ResolvedModule], len(resolutions))

		p.resolvedModulesMutex.Lock()
		defer p.resolvedModulesMutex.Unlock()
		if p.resolvedModules == nil {
			p.resolvedModules = make(map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule])
		}
		p.resolvedModules[file.Path()] = resolutionsInFile

		for _, resolution := range resolutions {
			resolvedFileName := resolution.resolvedModule.ResolvedFileName
			// TODO(ercornel): !!!: check if from node modules

			mode := getModeForUsageLocation(item, resolution.node, optionsForFile)
			resolutionsInFile[module.ModeAwareCacheKey{Name: resolution.node.Text(), Mode: mode}] = resolution.resolvedModule

			// add file to program only if:
			// - resolution was successful
			// - noResolve is falsy
			// - module name comes from the list of imports
			// - it's not a top level JavaScript module that exceeded the search max

			// const elideImport = isJsFileFromNodeModules && currentNodeModulesDepth > maxNodeModuleJsDepth;

			// Don't add the file if it has a bad extension (e.g. 'tsx' if we don't have '--allowJs')
			// This may still end up being an untyped module -- the file won't be included but imports will be allowed.

			shouldAddFile := resolution.resolvedModule.IsResolved() && tspath.FileExtensionIsOneOf(resolvedFileName, []string{".ts", ".tsx", ".mts", ".cts"})
			// TODO(ercornel): !!!: other checks on whether or not to add the file

			if shouldAddFile {
				// p.findSourceFile(resolvedFileName, FileIncludeReason{Import, 0})
				toParse = append(toParse, resolvedFileName)
			}
		}
	}
	return toParse
}

type resolution struct {
	node           *ast.Node
	resolvedModule *module.ResolvedModule
}

func (p *fileLoader) resolveModuleNames(entries []*ast.Node, item *sourceFileMetadata) []*resolution {
	if len(entries) == 0 {
		return nil
	}

	resolvedModules := make([]*resolution, 0, len(entries))

	for _, entry := range entries {
		moduleName := entry.Text()
		if moduleName == "" {
			continue
		}
		resolvedModule := p.resolver.ResolveModuleName(moduleName, item.file.FileName(), item.impliedNodeFormat /* !!! */, nil)
		resolvedModules = append(resolvedModules, &resolution{node: entry, resolvedModule: resolvedModule})
	}

	return resolvedModules
}

func (p *fileLoader) getCompilerOptionsForFile(file *ast.SourceFile) *core.CompilerOptions {
	// !!! return getRedirectReferenceForResolution(file)?.commandLine.options || options;
	return p.compilerOptions
}

func getModeForTypeReferenceDirectiveInFile(ref *ast.FileReference, item *sourceFileMetadata, options *core.CompilerOptions) core.ResolutionMode {
	if ref.ResolutionMode != core.ResolutionModeNone {
		return ref.ResolutionMode
	} else {
		return getDefaultResolutionModeForFile(item, options)
	}
}

func getDefaultResolutionModeForFile(item *sourceFileMetadata, options *core.CompilerOptions) core.ResolutionMode {
	if importSyntaxAffectsModuleResolution(options) {
		return getImpliedNodeFormatForEmitWorker(item, options)
	} else {
		return core.ResolutionModeNone
	}
}

/*
Use `program.getModeForUsageLocation`, which retrieves the correct `compilerOptions`, instead of this function whenever possible.
Calculates the final resolution mode for a given module reference node. This function only returns a result when module resolution
settings allow differing resolution between ESM imports and CJS requires, or when a mode is explicitly provided via import attributes,
which cause an `import` or `require` condition to be used during resolution regardless of module resolution settings. In absence of
overriding attributes, and in modes that support differing resolution, the result indicates the syntax the usage would emit to JavaScript.
Some examples:

	    // tsc foo.mts --module nodenext
	    import {} from "mod";
	    // Result: ESNext - the import emits as ESM due to `impliedNodeFormat` set by .mts file extension

		// tsc foo.cts --module nodenext
		import {} from "mod";
		// Result: CommonJS - the import emits as CJS due to `impliedNodeFormat` set by .cts file extension

		// tsc foo.ts --module preserve --moduleResolution bundler
		import {} from "mod";
		// Result: ESNext - the import emits as ESM due to `--module preserve` and `--moduleResolution bundler`
		// supports conditional imports/exports

		// tsc foo.ts --module preserve --moduleResolution node10
		import {} from "mod";
		// Result: undefined - the import emits as ESM due to `--module preserve`, but `--moduleResolution node10`
		// does not support conditional imports/exports

		// tsc foo.ts --module commonjs --moduleResolution node10
		import type {} from "mod" with { "resolution-mode": "import" };
		// Result: ESNext - conditional imports/exports always supported with "resolution-mode" attribute

Parameters:

- file: The file the import or import-like reference is contained within

- usage: The module reference string

- compilerOptions: The compiler options for the program that owns the file. If the file belongs to a referenced project, the compiler options
should be the options of the referenced project, not the referencing project.

Returns:

- The final resolution mode of the import
*/
func getModeForUsageLocation(item *sourceFileMetadata, usage *ast.Node, options *core.CompilerOptions) core.ResolutionMode {
	if ast.IsImportDeclaration(usage.Parent) || ast.IsExportDeclaration(usage.Parent) || usage.Parent.IsJSDocImportTag() {
		isTypeOnly := ast.IsExclusivelyTypeOnlyImportOrExport(usage.Parent)
		if isTypeOnly {
			var override core.ResolutionMode
			var ok bool
			switch usage.Parent.Kind {
			case ast.KindImportDeclaration:
				override, ok = usage.Parent.AsImportDeclaration().Attributes.GetResolutionModeOverride()
			case ast.KindExportDeclaration:
				override, ok = usage.Parent.AsExportDeclaration().Attributes.GetResolutionModeOverride()
			case ast.KindJSDocImportTag:
				override, ok = usage.Parent.AsJSDocImportTag().Attributes.GetResolutionModeOverride()
			}
			if ok {
				return override
			}
		}
	}
	if usage.Parent.Parent != nil && ast.IsImportTypeNode(usage.Parent.Parent) {
		if override, ok := usage.Parent.Parent.AsImportTypeNode().Attributes.GetResolutionModeOverride(); ok {
			return override
		}
	}

	if options != nil && importSyntaxAffectsModuleResolution(options) {
		return getEmitSyntaxForUsageLocationWorker(item, usage, options)
	}

	return core.ResolutionModeNone
}

func importSyntaxAffectsModuleResolution(options *core.CompilerOptions) bool {
	moduleResolution := options.ModuleResolution
	return core.ModuleResolutionKindNode16 <= moduleResolution && moduleResolution <= core.ModuleResolutionKindNodeNext ||
		options.ResolvePackageJsonExports.IsTrue() || options.ResolvePackageJsonImports.IsTrue()
}

func getEmitSyntaxForUsageLocationWorker(item *sourceFileMetadata, usage *ast.Node, options *core.CompilerOptions) core.ResolutionMode {
	if options == nil {
		// This should always be provided, but we try to fail somewhat
		// gracefully to allow projects like ts-node time to update.
		return core.ResolutionModeNone
	}

	exprParentParent := ast.WalkUpParenthesizedExpressions(usage.Parent).Parent
	if exprParentParent != nil && ast.IsImportEqualsDeclaration(exprParentParent) || ast.IsRequireCall(usage.Parent /*requireStringLiteralLikeArgument*/, false) {
		return core.ModuleKindCommonJS
	}
	if ast.IsImportCall(ast.WalkUpParenthesizedExpressions(usage.Parent)) {
		if shouldTransformImportCallWorker(item, options) {
			return core.ModuleKindCommonJS
		} else {
			return core.ModuleKindESNext
		}
	}
	// If we're in --module preserve on an input file, we know that an import
	// is an import. But if this is a declaration file, we'd prefer to use the
	// impliedNodeFormat. Since we want things to be consistent between the two,
	// we need to issue errors when the user writes ESM syntax in a definitely-CJS
	// file, until/unless declaration emit can indicate a true ESM import. On the
	// other hand, writing CJS syntax in a definitely-ESM file is fine, since declaration
	// emit preserves the CJS syntax.
	fileEmitMode := getEmitModuleFormatOfFile(item, options)
	if fileEmitMode == core.ModuleKindCommonJS {
		return core.ModuleKindCommonJS
	} else {
		if fileEmitMode.IsNonNodeESM() || fileEmitMode == core.ModuleKindPreserve {
			return core.ModuleKindESNext
		}
	}
	return core.ModuleKindNone
}

func getEmitModuleFormatOfFile(item *sourceFileMetadata, options *core.CompilerOptions) core.ResolutionMode {
	implied := getImpliedNodeFormatForEmitWorker(item, options)
	if implied != core.ResolutionModeNone {
		return implied
	} else {
		return options.GetEmitModuleKind()
	}
}

func shouldTransformImportCallWorker(metadata *sourceFileMetadata, options *core.CompilerOptions) bool {
	moduleKind := options.GetEmitModuleKind()
	if core.ModuleKindNode16 <= moduleKind && moduleKind <= core.ModuleKindNodeNext || moduleKind == core.ModuleKindPreserve {
		return false
	}
	return getImpliedNodeFormatForEmitWorker(metadata, options) < core.ModuleKindES2015
}

func getImpliedNodeFormatForEmitWorker(metadata *sourceFileMetadata, options *core.CompilerOptions) core.ResolutionMode {
	moduleKind := options.ModuleKind
	if core.ModuleKindNode16 <= moduleKind && moduleKind <= core.ModuleKindNodeNext {
		return metadata.impliedNodeFormat
	}
	if metadata.impliedNodeFormat == core.ModuleKindCommonJS &&
		(metadata.packageJsonScope.Contents.Type.Value == "commonjs" ||
			tspath.FileExtensionIsOneOf(metadata.file.FileName(), []string{tspath.ExtensionCjs, tspath.ExtensionCts})) {
		return core.ModuleKindCommonJS
	}
	if metadata.impliedNodeFormat == core.ModuleKindESNext &&
		(metadata.packageJsonScope.Contents.Type.Value == "module" ||
			tspath.FileExtensionIsOneOf(metadata.file.FileName(), []string{tspath.ExtensionMjs, tspath.ExtensionMts})) {
		return core.ModuleKindESNext
	}
	return core.ResolutionModeNone
}
