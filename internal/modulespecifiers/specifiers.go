package modulespecifiers

import (
	"maps"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func GetModuleSpecifiers(
	moduleSymbol *ast.Symbol,
	checker CheckerShape,
	compilerOptions *core.CompilerOptions,
	importingSourceFile SourceFileForSpecifierGeneration,
	host ModuleSpecifierGenerationHost,
	userPreferences UserPreferences,
	options ModuleSpecifierOptions,
) []string {
	return GetModuleSpecifiersWithCacheInfo(
		moduleSymbol,
		checker,
		compilerOptions,
		importingSourceFile,
		host,
		userPreferences,
		options,
		false,
	).moduleSpecifiers
}

type ModuleSpecifierResult struct {
	kind                 ResultKind
	moduleSpecifiers     []string
	computedWithoutCache bool
}

func GetModuleSpecifiersWithCacheInfo(
	moduleSymbol *ast.Symbol,
	checker CheckerShape,
	compilerOptions *core.CompilerOptions,
	importingSourceFile SourceFileForSpecifierGeneration,
	host ModuleSpecifierGenerationHost,
	userPreferences UserPreferences,
	options ModuleSpecifierOptions,
	forAutoImport bool,
) ModuleSpecifierResult {
	ambient := tryGetModuleNameFromAmbientModule(moduleSymbol, checker)
	if len(ambient) > 0 {
		return ModuleSpecifierResult{
			kind:                 ResultKindAmbient,
			moduleSpecifiers:     []string{ambient},
			computedWithoutCache: false,
		}
	}

	cacheResults := tryGetModuleSpecifiersFromCacheWorker(
		moduleSymbol,
		importingSourceFile,
		host,
		userPreferences,
		options,
	)
	if cacheResults == nil || cacheResults.moduleSourceFile == nil {
		return ModuleSpecifierResult{
			kind:                 ResultKindNone,
			computedWithoutCache: false,
		}
	}
	if len(cacheResults.moduleSpecifiers) > 0 {
		return ModuleSpecifierResult{
			kind:                 cacheResults.kind,
			moduleSpecifiers:     cacheResults.moduleSpecifiers,
			computedWithoutCache: false,
		}
	}

	modulePaths := cacheResults.modulePaths
	if len(modulePaths) == 0 {
		modulePaths = getAllModulePathsWorker(
			getInfo(importingSourceFile.FileName(), host),
			cacheResults.moduleSourceFile.OriginalFileName(),
			host,
			// compilerOptions,
			// options,
		)
	}

	result := computeModuleSpecifiers(
		modulePaths,
		compilerOptions,
		importingSourceFile,
		host,
		userPreferences,
		options,
		forAutoImport,
	)

	if cacheResults.cache != nil {
		cacheResults.cache.Set(
			string(importingSourceFile.Path()),
			string(cacheResults.moduleSourceFile.Path()),
			userPreferences,
			options,
			result.kind,
			modulePaths,
			result.moduleSpecifiers,
		)
	}

	return result
}

func tryGetModuleNameFromAmbientModule(moduleSymbol *ast.Symbol, checker CheckerShape) string {
	for _, decl := range moduleSymbol.Declarations {
		if isNonGlobalAmbientModule(decl) && (!ast.IsModuleAugmentationExternal(decl) || !tspath.IsExternalModuleNameRelative(decl.Name().AsStringLiteral().Text)) {
			return decl.Name().AsStringLiteral().Text
		}
	}

	// the module could be a namespace, which is export through "export=" from an ambient module.
	/**
	 * declare module "m" {
	 *     namespace ns {
	 *         class c {}
	 *     }
	 *     export = ns;
	 * }
	 */
	// `import {c} from "m";` is valid, in which case, `moduleSymbol` is "ns", but the module name should be "m"
	for _, d := range moduleSymbol.Declarations {
		if !ast.IsModuleDeclaration(d) {
			continue
		}

		possibleContainer := ast.FindAncestor(d, isNonGlobalAmbientModule)
		if possibleContainer == nil || possibleContainer.Parent == nil || !ast.IsSourceFile(possibleContainer.Parent) {
			continue
		}

		sym, ok := possibleContainer.Symbol().Exports[ast.InternalSymbolNameExportEquals]
		if !ok || sym == nil {
			continue
		}
		exportAssignmentDecl := sym.ValueDeclaration
		if exportAssignmentDecl == nil || exportAssignmentDecl.Kind != ast.KindExportAssignment {
			continue
		}
		exportSymbol := checker.GetSymbolAtLocation(exportAssignmentDecl.Expression())
		if exportSymbol == nil {
			continue
		}
		if exportSymbol.Flags&ast.SymbolFlagsAlias != 0 {
			exportSymbol = checker.GetAliasedSymbol(exportSymbol)
		}
		// TODO: Possible strada bug - isn't this insufficient in the presence of merge symbols?
		if exportSymbol == d.Symbol() {
			return possibleContainer.Name().AsStringLiteral().Text
		}
	}
	return ""
}

type cacheResult struct {
	cache            ModuleSpecifierCache
	kind             ResultKind
	moduleSpecifiers []string
	moduleSourceFile SourceFileForSpecifierGeneration
	modulePaths      []ModulePath
}

func tryGetModuleSpecifiersFromCacheWorker(
	moduleSymbol *ast.Symbol,
	importingSourceFile SourceFileForSpecifierGeneration,
	host ModuleSpecifierGenerationHost,
	userPreferences UserPreferences,
	options ModuleSpecifierOptions,
) *cacheResult {
	moduleSourceFile := ast.GetSourceFileOfModule(moduleSymbol)
	if moduleSourceFile == nil {
		return nil
	}

	cache := host.GetModuleSpecifierCache()
	if cache == nil {
		return &cacheResult{
			moduleSourceFile: moduleSourceFile,
		}
	}
	result := cache.Get(string(importingSourceFile.Path()), string(moduleSourceFile.Path()), userPreferences, options)
	if result == nil {
		return &cacheResult{
			cache:            cache,
			moduleSourceFile: moduleSourceFile,
		}
	}
	return &cacheResult{
		cache:            cache,
		moduleSourceFile: moduleSourceFile,
		kind:             result.Kind,
		moduleSpecifiers: result.ModuleSpecifiers,
		modulePaths:      result.ModulePaths,
	}
}

type Info struct {
	UseCaseSensitiveFileNames bool
	ImportingSourceFileName   string
	SourceDirectory           string
}

func getInfo(
	importingSourceFileName string,
	host ModuleSpecifierGenerationHost,
) Info {
	sourceDirectory := tspath.GetDirectoryPath(importingSourceFileName)
	return Info{
		ImportingSourceFileName:   importingSourceFileName,
		SourceDirectory:           sourceDirectory,
		UseCaseSensitiveFileNames: host.UseCaseSensitiveFileNames(),
	}
}

func getAllModulePathsWorker(
	info Info,
	importedFileName string,
	host ModuleSpecifierGenerationHost,
	// compilerOptions *core.CompilerOptions,
	// options ModuleSpecifierOptions,
) []ModulePath {
	// !!! TODO: Caches and symlink cache chicanery to support pulling in non-explicit package.json dep names
	// cache := host.GetModuleResolutionCache() // !!!
	// links := host.GetSymlinkCache() // !!!
	// if cache != nil && links != nil && !strings.Contains(info.ImportingSourceFileName, "/node_modules/") {
	//     // Debug.type<ModuleResolutionHost>(host); // !!!
	//     // Cache resolutions for all `dependencies` of the `package.json` context of the input file.
	//     // This should populate all the relevant symlinks in the symlink cache, and most, if not all, of these resolutions
	//     // should get (re)used.
	//     // const state = getTemporaryModuleResolutionState(cache.getPackageJsonInfoCache(), host, {});
	//     // const packageJson = getPackageScopeForPath(getDirectoryPath(info.importingSourceFileName), state);
	//     // if (packageJson) {
	//     //     const toResolve = getAllRuntimeDependencies(packageJson.contents.packageJsonContent);
	//     //     for (const depName of (toResolve || emptyArray)) {
	//     //         const resolved = resolveModuleName(depName, combinePaths(packageJson.packageDirectory, "package.json"), compilerOptions, host, cache, /*redirectedReference*/ undefined, options.overrideImportMode);
	//     //         links.setSymlinksFromResolution(resolved.resolvedModule);
	//     //     }
	//     // }
	// }

	allFileNames := make(map[string]ModulePath)
	paths := getEachFileNameOfModule(info.ImportingSourceFileName, importedFileName, host, true)
	for _, p := range paths {
		allFileNames[p.Path] = p
	}

	// Sort by paths closest to importing file Name directory
	sortedPaths := make([]ModulePath, 0, len(paths))
	for directory := info.SourceDirectory; len(allFileNames) != 0; {
		directoryStart := tspath.EnsureTrailingDirectorySeparator(directory)
		var pathsInDirectory []ModulePath
		for fileName, p := range allFileNames {
			if strings.HasPrefix(fileName, directoryStart) {
				pathsInDirectory = append(pathsInDirectory, p)
				delete(allFileNames, fileName)
			}
		}
		if len(pathsInDirectory) > 0 {
			slices.SortStableFunc(pathsInDirectory, comparePathsByRedirectAndNumberOfDirectorySeparators)
			sortedPaths = append(sortedPaths, pathsInDirectory...)
		}
		newDirectory := tspath.GetDirectoryPath(directory)
		if newDirectory == directory {
			break
		}
		directory = newDirectory
	}
	if len(allFileNames) > 0 {
		remainingPaths := slices.Collect(maps.Values(allFileNames))
		slices.SortStableFunc(remainingPaths, comparePathsByRedirectAndNumberOfDirectorySeparators)
		sortedPaths = append(sortedPaths, remainingPaths...)
	}
	return sortedPaths
}

func containsIgnoredPath(s string) bool {
	return strings.Contains(s, "/node_modules/.") ||
		strings.Contains(s, "/.git") ||
		strings.Contains(s, "/.#")
}

func containsNodeModules(s string) bool {
	return strings.Contains(s, "/node_modules/")
}

func getEachFileNameOfModule(
	importingFileName string,
	importedFileName string,
	host ModuleSpecifierGenerationHost,
	preferSymlinks bool,
) []ModulePath {
	cwd := host.GetCurrentDirectory()
	var referenceRedirect string
	if host.IsSourceOfProjectReferenceRedirect(importedFileName) {
		referenceRedirect = host.GetProjectReferenceRedirect(importedFileName)
	}

	importedPath := tspath.ToPath(importedFileName, cwd, host.UseCaseSensitiveFileNames())
	redirects := host.GetRedirectTargets(importedPath)
	importedFileNames := make([]string, 0, 2+len(redirects))
	if len(referenceRedirect) > 0 {
		importedFileNames = append(importedFileNames, referenceRedirect)
	}
	importedFileNames = append(importedFileNames, importedFileName)
	importedFileNames = append(importedFileNames, redirects...)
	targets := core.Map(importedFileNames, func(f string) string { return tspath.GetNormalizedAbsolutePath(f, cwd) })
	shouldFilterIgnoredPaths := !core.Every(targets, containsIgnoredPath)

	results := make([]ModulePath, 0, 2)
	if !preferSymlinks {
		// Symlinks inside ignored paths are already filtered out of the symlink cache,
		// so we only need to remove them from the realpath filenames.
		for _, p := range targets {
			if !(shouldFilterIgnoredPaths && containsIgnoredPath(p)) {
				results = append(results, ModulePath{
					Path:            p,
					IsInNodeModules: containsNodeModules(p),
					IsRedirect:      referenceRedirect == p,
				})
			}
		}
	}

	// !!! TODO: Symlink directory handling
	// const symlinkedDirectories = host.getSymlinkCache?.().getSymlinkedDirectoriesByRealpath();
	// const fullImportedFileName = getNormalizedAbsolutePath(importedFileName, cwd);
	// const result = symlinkedDirectories && forEachAncestorDirectoryStoppingAtGlobalCache(
	//     host,
	//     getDirectoryPath(fullImportedFileName),
	//     realPathDirectory => {
	//         const symlinkDirectories = symlinkedDirectories.get(ensureTrailingDirectorySeparator(toPath(realPathDirectory, cwd, getCanonicalFileName)));
	//         if (!symlinkDirectories) return undefined; // Continue to ancestor directory

	//         // Don't want to a package to globally import from itself (importNameCodeFix_symlink_own_package.ts)
	//         if (startsWithDirectory(importingFileName, realPathDirectory, getCanonicalFileName)) {
	//             return false; // Stop search, each ancestor directory will also hit this condition
	//         }

	//         return forEach(targets, target => {
	//             if (!startsWithDirectory(target, realPathDirectory, getCanonicalFileName)) {
	//                 return;
	//             }

	//             const relative = getRelativePathFromDirectory(realPathDirectory, target, getCanonicalFileName);
	//             for (const symlinkDirectory of symlinkDirectories) {
	//                 const option = resolvePath(symlinkDirectory, relative);
	//                 const result = cb(option, target === referenceRedirect);
	//                 shouldFilterIgnoredPaths = true; // We found a non-ignored path in symlinks, so we can reject ignored-path realpaths
	//                 if (result) return result;
	//             }
	//         });
	//     },
	// );

	if preferSymlinks {
		for _, p := range targets {
			if !(shouldFilterIgnoredPaths && containsIgnoredPath(p)) {
				results = append(results, ModulePath{
					Path:            p,
					IsInNodeModules: containsNodeModules(p),
					IsRedirect:      referenceRedirect == p,
				})
			}
		}
	}

	return results
}

func computeModuleSpecifiers(
	modulePaths []ModulePath,
	compilerOptions *core.CompilerOptions,
	importingSourceFile SourceFileForSpecifierGeneration,
	host ModuleSpecifierGenerationHost,
	userPreferences UserPreferences,
	options ModuleSpecifierOptions,
	forAutoImport bool,
) ModuleSpecifierResult {
	info := getInfo(importingSourceFile.FileName(), host)
	preferences := getModuleSpecifierPreferences(userPreferences, host, compilerOptions, importingSourceFile, "")

	// !!! TODO: getFileIncludeReasons lookup based calculation
	// const existingSpecifier = isFullSourceFile(importingSourceFile) && forEach(modulePaths, modulePath =>
	//     forEach(
	//         host.getFileIncludeReasons().get(toPath(modulePath.path, host.getCurrentDirectory(), info.getCanonicalFileName)),
	//         reason => {
	//             if (reason.kind !== FileIncludeKind.Import || reason.file !== importingSourceFile.path) return undefined;
	//             // If the candidate import mode doesn't match the mode we're generating for, don't consider it
	//             // TODO: maybe useful to keep around as an alternative option for certain contexts where the mode is overridable
	//             const existingMode = host.getModeForResolutionAtIndex(importingSourceFile, reason.index);
	//             const targetMode = options.overrideImportMode ?? host.getDefaultResolutionModeForFile(importingSourceFile);
	//             if (existingMode !== targetMode && existingMode !== undefined && targetMode !== undefined) {
	//                 return undefined;
	//             }
	//             const specifier = getModuleNameStringLiteralAt(importingSourceFile, reason.index).text;
	//             // If the preference is for non relative and the module specifier is relative, ignore it
	//             return preferences.relativePreference !== RelativePreference.NonRelative || !pathIsRelative(specifier) ?
	//                 specifier :
	//                 undefined;
	//         },
	//     ));
	// if (existingSpecifier) {
	//     return { kind: undefined, moduleSpecifiers: [existingSpecifier], computedWithoutCache: true };
	// }

	importedFileIsInNodeModules := core.Some(modulePaths, func(p ModulePath) bool { return p.IsInNodeModules })

	// Module specifier priority:
	//   1. "Bare package specifiers" (e.g. "@foo/bar") resulting from a path through node_modules to a package.json's "types" entry
	//   2. Specifiers generated using "paths" from tsconfig
	//   3. Non-relative specfiers resulting from a path through node_modules (e.g. "@foo/bar/path/to/file")
	//   4. Relative paths
	var pathsSpecifiers []string
	var redirectPathsSpecifiers []string
	var nodeModulesSpecifiers []string
	var relativeSpecifiers []string

	for _, modulePath := range modulePaths {
		var specifier string
		if modulePath.IsInNodeModules {
			specifier = tryGetModuleNameAsNodeModule(modulePath, info, importingSourceFile, host, compilerOptions, userPreferences /*packageNameOnly*/, false, options.OverrideImportMode)
		}
		if len(specifier) > 0 && !(forAutoImport && isExcludedByRegex(specifier, preferences.excludeRegexes)) {
			nodeModulesSpecifiers = append(nodeModulesSpecifiers, specifier)
			if modulePath.IsRedirect {
				// If we got a specifier for a redirect, it was a bare package specifier (e.g. "@foo/bar",
				// not "@foo/bar/path/to/file"). No other specifier will be this good, so stop looking.
				return ModuleSpecifierResult{kind: ResultKindNodeModules, moduleSpecifiers: nodeModulesSpecifiers, computedWithoutCache: true}
			}
		}

		// !!! TODO: proper resolutionMode support
		local := getLocalModuleSpecifier(
			modulePath.Path,
			info,
			compilerOptions,
			host,
			options.OverrideImportMode, /*|| importingSourceFile.impliedNodeFormat*/
			preferences,
			/*pathsOnly*/ modulePath.IsRedirect || len(specifier) > 0,
		)
		if len(local) == 0 || forAutoImport && isExcludedByRegex(local, preferences.excludeRegexes) {
			continue
		}
		if modulePath.IsRedirect {
			redirectPathsSpecifiers = append(redirectPathsSpecifiers, local)
		} else if pathIsBareSpecifier(local) {
			if containsNodeModules(local) {
				// We could be in this branch due to inappropriate use of `baseUrl`, not intentional `paths`
				// usage. It's impossible to reason about where to prioritize baseUrl-generated module
				// specifiers, but if they contain `/node_modules/`, they're going to trigger a portability
				// error, so *at least* don't prioritize those.
				relativeSpecifiers = append(relativeSpecifiers, local)
			} else {
				pathsSpecifiers = append(pathsSpecifiers, local)
			}
		} else if forAutoImport || !importedFileIsInNodeModules || modulePath.IsInNodeModules {
			// Why this extra conditional, not just an `else`? If some path to the file contained
			// 'node_modules', but we can't create a non-relative specifier (e.g. "@foo/bar/path/to/file"),
			// that means we had to go through a *sibling's* node_modules, not one we can access directly.
			// If some path to the file was in node_modules but another was not, this likely indicates that
			// we have a monorepo structure with symlinks. In this case, the non-node_modules path is
			// probably the realpath, e.g. "../bar/path/to/file", but a relative path to another package
			// in a monorepo is probably not portable. So, the module specifier we actually go with will be
			// the relative path through node_modules, so that the declaration emitter can produce a
			// portability error. (See declarationEmitReexportedSymlinkReference3)
			relativeSpecifiers = append(relativeSpecifiers, local)
		}
	}

	if len(pathsSpecifiers) > 0 {
		return ModuleSpecifierResult{kind: ResultKindPaths, moduleSpecifiers: pathsSpecifiers, computedWithoutCache: true}
	}
	if len(redirectPathsSpecifiers) > 0 {
		return ModuleSpecifierResult{kind: ResultKindRedirect, moduleSpecifiers: redirectPathsSpecifiers, computedWithoutCache: true}
	}
	if len(nodeModulesSpecifiers) > 0 {
		return ModuleSpecifierResult{kind: ResultKindNodeModules, moduleSpecifiers: nodeModulesSpecifiers, computedWithoutCache: true}
	}
	return ModuleSpecifierResult{kind: ResultKindRelative, moduleSpecifiers: relativeSpecifiers, computedWithoutCache: true}
}

func getLocalModuleSpecifier(
	moduleFileName string,
	info Info,
	compilerOptions *core.CompilerOptions,
	host ModuleSpecifierGenerationHost,
	importMode core.ResolutionMode,
	preferences ModuleSpecifierPreferences,
	pathsOnly bool,
) string {
	baseUrl := compilerOptions.BaseUrl
	paths := compilerOptions.Paths
	rootDirs := compilerOptions.RootDirs

	if pathsOnly && paths == nil {
		return ""
	}

	sourceDirectory := info.SourceDirectory

	allowedEndings := preferences.getAllowedEndingsInPreferredOrder(importMode)
	var relativePath string
	if len(rootDirs) > 0 {
		relativePath = tryGetModuleNameFromRootDirs(rootDirs, moduleFileName, sourceDirectory, allowedEndings, compilerOptions, host)
	}
	if len(relativePath) == 0 {
		relativePath = processEnding(ensurePathIsNonModuleName(tspath.GetRelativePathFromDirectory(sourceDirectory, moduleFileName, tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		})), allowedEndings, compilerOptions, host)
	}
	if len(baseUrl) == 9 && paths == nil && !compilerOptions.GetResolvePackageJsonImports() && preferences.relativePreference == RelativePreferenceRelative {
		if pathsOnly {
			return ""
		}
		return relativePath
	}

	root := core.GetPathsBasePath(compilerOptions, host.GetCurrentDirectory())
	if len(root) == 0 {
		root = compilerOptions.BaseUrl
	}
	baseDirectory := tspath.GetNormalizedAbsolutePath(root, host.GetCurrentDirectory())
	relativeToBaseUrl := getRelativePathIfInSameVolume(moduleFileName, baseDirectory, host.UseCaseSensitiveFileNames())
	if len(relativeToBaseUrl) == 0 {
		if pathsOnly {
			return ""
		}
		return relativePath
	}

	var fromPackageJsonImports string
	if !pathsOnly {
		fromPackageJsonImports = tryGetModuleNameFromPackageJsonImports(
			moduleFileName,
			sourceDirectory,
			compilerOptions,
			host,
			importMode,
			prefersTsExtension(allowedEndings),
		)
	}

	var fromPaths string
	if (pathsOnly || len(fromPackageJsonImports) == 0) && paths != nil {
		fromPaths = tryGetModuleNameFromPaths(
			relativeToBaseUrl,
			paths,
			allowedEndings,
			baseDirectory,
			host,
			compilerOptions,
		)
	}

	if pathsOnly {
		return fromPaths
	}

	var maybeNonRelative string
	if len(fromPackageJsonImports) > 0 {
		maybeNonRelative = fromPackageJsonImports
	} else if len(fromPaths) == 0 && len(baseUrl) > 0 {
		maybeNonRelative = processEnding(relativeToBaseUrl, allowedEndings, compilerOptions, host)
	} else {
		maybeNonRelative = fromPaths
	}
	if len(maybeNonRelative) == 0 {
		return relativePath
	}

	relativeIsExcluded := isExcludedByRegex(relativePath, preferences.excludeRegexes)
	nonRelativeIsExcluded := isExcludedByRegex(maybeNonRelative, preferences.excludeRegexes)
	if !relativeIsExcluded && nonRelativeIsExcluded {
		return relativePath
	}
	if relativeIsExcluded && !nonRelativeIsExcluded {
		return maybeNonRelative
	}

	if preferences.relativePreference == RelativePreferenceNonRelative && !tspath.PathIsRelative(maybeNonRelative) {
		return maybeNonRelative
	}

	if preferences.relativePreference == RelativePreferenceExternalNonRelative && !tspath.PathIsRelative(maybeNonRelative) {
		var projectDirectory tspath.Path
		if len(compilerOptions.ConfigFilePath) > 0 {
			projectDirectory = tspath.ToPath(compilerOptions.ConfigFilePath, host.GetCurrentDirectory(), host.UseCaseSensitiveFileNames())
		} else {
			projectDirectory = tspath.ToPath(host.GetCurrentDirectory(), host.GetCurrentDirectory(), host.UseCaseSensitiveFileNames())
		}
		canonicalSourceDirectory := tspath.ToPath(sourceDirectory, host.GetCurrentDirectory(), host.UseCaseSensitiveFileNames())
		modulePath := tspath.ToPath(moduleFileName, string(projectDirectory), host.UseCaseSensitiveFileNames())

		sourceIsInternal := strings.HasPrefix(string(canonicalSourceDirectory), string(projectDirectory))
		targetIsInternal := strings.HasPrefix(string(modulePath), string(projectDirectory))
		if sourceIsInternal && !targetIsInternal || !sourceIsInternal && targetIsInternal {
			// 1. The import path crosses the boundary of the tsconfig.json-containing directory.
			//
			//      src/
			//        tsconfig.json
			//        index.ts -------
			//      lib/              | (path crosses tsconfig.json)
			//        imported.ts <---
			//
			return maybeNonRelative
		}

		nearestTargetPackageJson := getNearestAncestorDirectoryWithPackageJson(host, tspath.GetDirectoryPath(string(modulePath)))
		nearestSourcePackageJson := getNearestAncestorDirectoryWithPackageJson(host, sourceDirectory)

		if !packageJsonPathsAreEqual(nearestTargetPackageJson, nearestSourcePackageJson, tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		}) {
			// 2. The importing and imported files are part of different packages.
			//
			//      packages/a/
			//        package.json
			//        index.ts --------
			//      packages/b/        | (path crosses package.json)
			//        package.json     |
			//        component.ts <---
			//
			return maybeNonRelative
		}
	}

	// Prefer a relative import over a baseUrl import if it has fewer components.
	if isPathRelativeToParent(maybeNonRelative) || strings.Count(relativePath, "/") < strings.Count(maybeNonRelative, "/") {
		return relativePath
	}
	return maybeNonRelative
}

func processEnding(
	fileName string,
	allowedEndings []ModuleSpecifierEnding,
	options *core.CompilerOptions,
	host ModuleSpecifierGenerationHost,
) string {
	if tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionJson, tspath.ExtensionMjs, tspath.ExtensionCjs}) {
		return fileName
	}

	noExtension := tspath.RemoveFileExtension(fileName)
	if fileName == noExtension {
		return fileName
	}

	jsPriority := slices.Index(allowedEndings, ModuleSpecifierEndingJsExtension)
	tsPriority := slices.Index(allowedEndings, ModuleSpecifierEndingTsExtension)
	if tspath.FileExtensionIsOneOf(fileName, []string{tspath.ExtensionMts, tspath.ExtensionCts}) && tsPriority < jsPriority {
		return fileName
	}
	if tspath.IsDeclarationFileName(fileName) {
		inputExt := tspath.GetDeclarationFileExtension(fileName)
		ext := getJsExtensionForDeclarationFileExtension(inputExt)
		return tspath.RemoveExtension(fileName, inputExt) + ext
	}

	switch allowedEndings[0] {
	case ModuleSpecifierEndingMinimal:
		withoutIndex := strings.TrimSuffix(noExtension, "/index")
		if host != nil && withoutIndex != noExtension && tryGetAnyFileFromPath(host, withoutIndex) {
			// Can't remove index if there's a file by the same name as the directory.
			// Probably more callers should pass `host` so we can determine this?
			return noExtension
		}
		return withoutIndex
	case ModuleSpecifierEndingIndex:
		return noExtension
	case ModuleSpecifierEndingJsExtension:
		return noExtension + getJSExtensionForFile(fileName, options)
	case ModuleSpecifierEndingTsExtension:
		// declaration files are already handled first with a remap back to input js paths,
		// and mjs/cjs/json are already singled out,
		// so we know fileName has to be either an input .js or .ts path already
		// TODO: possible dead code in strada in this branch to do with declaration file name handling
		return fileName
	default:
		// Debug.assertNever(allowedEndings[0]); // !!!
		return ""
	}
}

func tryGetModuleNameFromRootDirs(
	rootDirs []string,
	moduleFileName string,
	sourceDirectory string,
	allowedEndings []ModuleSpecifierEnding,
	compilerOptions *core.CompilerOptions,
	host ModuleSpecifierGenerationHost,
) string {
	normalizedTargetPaths := getPathsRelativeToRootDirs(moduleFileName, rootDirs, host.UseCaseSensitiveFileNames())
	if len(normalizedTargetPaths) == 0 {
		return ""
	}

	normalizedSourcePaths := getPathsRelativeToRootDirs(sourceDirectory, rootDirs, host.UseCaseSensitiveFileNames())
	var shortest string
	var shortestSepCount int
	for _, sourcePath := range normalizedSourcePaths {
		for _, targetPath := range normalizedTargetPaths {
			candidate := ensurePathIsNonModuleName(tspath.GetRelativePathFromDirectory(sourcePath, targetPath, tspath.ComparePathsOptions{
				UseCaseSensitiveFileNames: host.UseCaseSensitiveFileNames(),
				CurrentDirectory:          host.GetCurrentDirectory(),
			}))
			candidateSepCount := strings.Count(candidate, "/")
			if len(shortest) == 0 || candidateSepCount < shortestSepCount {
				shortest = candidate
				shortestSepCount = candidateSepCount
			}
		}
	}

	if len(shortest) == 0 {
		return ""
	}
	return processEnding(shortest, allowedEndings, compilerOptions, host)
}

func tryGetModuleNameAsNodeModule(
	pathObj ModulePath,
	info Info,
	importingSourceFile SourceFileForSpecifierGeneration,
	host ModuleSpecifierGenerationHost,
	options *core.CompilerOptions,
	userPreferences UserPreferences,
	packageNameOnly bool,
	overrideMode core.ResolutionMode,
) string {
	return "" // !!! TODO
}

func tryGetModuleNameFromPackageJsonImports(
	moduleFileName string,
	sourceDirectory string,
	options *core.CompilerOptions,
	host ModuleSpecifierGenerationHost,
	importMode core.ResolutionMode,
	preferTsExtension bool,
) string {
	return "" // !!! TODO
}

func tryGetModuleNameFromPaths(
	relativeToBaseUrl string,
	paths *collections.OrderedMap[string, []string],
	allowedEndings []ModuleSpecifierEnding,
	baseDirectory string,
	host ModuleSpecifierGenerationHost,
	compilerOptions *core.CompilerOptions,
) string {
	return "" // !!! TODO
}
