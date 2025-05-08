package project

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"slices"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/semver"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type CachedTyping struct {
	Location string
	Version  semver.Version
}

func IsTypingUpToDate(cachedTyping *CachedTyping, availableTypingVersions map[string]string) bool {
	useVersion, ok := availableTypingVersions["ts"+core.VersionMajorMinor]
	if !ok {
		useVersion = availableTypingVersions["latest"]
	}
	availableVersion := semver.MustParse(useVersion)
	return availableVersion.Compare(&cachedTyping.Version) <= 0
}

func DiscoverTypings(
	p *Project,
	typingsInfo *TypingsCacheInfo,
	fileNames []string,
	projectRootPath string,
	safeList map[string]string,
	packageNameToTypingLocation *collections.SyncMap[string, *CachedTyping],
	typesRegistry map[string]map[string]string,
) (cachedTypingPaths []string, newTypingNames []string, filesToWatch []string) {
	// A typing name to typing file path mapping
	inferredTypings := map[string]string{}

	// Only infer typings for .js and .jsx files
	fileNames = core.Filter(core.Map(fileNames, func(fileName string) string {
		if tspath.HasJSFileExtension(fileName) {
			return fileName
		}
		return ""
	}), func(fileName string) bool {
		return fileName != ""
	})

	if typingsInfo.typeAcquisition.Include != nil {
		addInferredTypings(p, inferredTypings, typingsInfo.typeAcquisition.Include, "Explicitly included types")
	}
	exclude := typingsInfo.typeAcquisition.Exclude

	// Directories to search for package.json, bower.json and other typing information
	if typingsInfo.compilerOptions.Types != nil {
		possibleSearchDirs := map[string]bool{}
		for _, fileName := range fileNames {
			possibleSearchDirs[tspath.GetDirectoryPath(fileName)] = true
		}
		possibleSearchDirs[projectRootPath] = true
		for searchDir := range possibleSearchDirs {
			filesToWatch = addTypingNamesAndGetFilesToWatch(p, inferredTypings, filesToWatch, searchDir, "bower.json", "bower_components")
			filesToWatch = addTypingNamesAndGetFilesToWatch(p, inferredTypings, filesToWatch, searchDir, "package.json", "node_modules")
		}
	}

	if !typingsInfo.typeAcquisition.DisableFilenameBasedTypeAcquisition.IsTrue() {
		getTypingNamesFromSourceFileNames(p, inferredTypings, safeList, fileNames)
	}

	// add typings for unresolved imports
	modules := slices.Compact(core.Map(typingsInfo.unresolvedImports, core.NonRelativeModuleNameForTypingCache))
	addInferredTypings(p, inferredTypings, modules, "Inferred typings from unresolved imports")

	// Remove typings that the user has added to the exclude list
	for _, excludeTypingName := range exclude {
		delete(inferredTypings, excludeTypingName)
		p.Logf("Typing for %s is in exclude list, will be ignored.", excludeTypingName)
	}

	// Add the cached typing locations for inferred typings that are already installed
	packageNameToTypingLocation.Range(func(name string, typing *CachedTyping) bool {
		registryEntry := typesRegistry[name]
		if inferredTypings[name] == "" && registryEntry != nil && IsTypingUpToDate(typing, registryEntry) {
			inferredTypings[name] = typing.Location
		}
		return true
	})

	for typing, inferred := range inferredTypings {
		if inferred != "" {
			cachedTypingPaths = append(cachedTypingPaths, inferred)
		} else {
			newTypingNames = append(newTypingNames, typing)
		}
	}
	p.Logf("Finished typings discovery: cachedTypingsPaths: %v newTypingNames: %v, filesToWatch %v", cachedTypingPaths, newTypingNames, filesToWatch)
	return cachedTypingPaths, newTypingNames, filesToWatch
}

func addInferredTyping(inferredTypings map[string]string, typingName string) {
	if _, ok := inferredTypings[typingName]; !ok {
		inferredTypings[typingName] = ""
	}
}

func addInferredTypings(p *Project, inferredTypings map[string]string, typingNames []string, message string) {
	p.Logf("%s: %v", message, typingNames)
	for _, typingName := range typingNames {
		addInferredTyping(inferredTypings, typingName)
	}
}

/**
 * Infer typing names from given file names. For example, the file name "jquery-min.2.3.4.js"
 * should be inferred to the 'jquery' typing name; and "angular-route.1.2.3.js" should be inferred
 * to the 'angular-route' typing name.
 * @param fileNames are the names for source files in the project
 */
func getTypingNamesFromSourceFileNames(
	p *Project,
	inferredTypings map[string]string,
	safeList map[string]string,
	fileNames []string,
) {
	hasJsxFile := false
	var fromFileNames []string
	for _, fileName := range fileNames {
		hasJsxFile = hasJsxFile || tspath.FileExtensionIs(fileName, tspath.ExtensionJs)
		inferredTypingName := tspath.RemoveFileExtension(tspath.ToFileNameLowerCase(tspath.GetBaseFileName(fileName)))
		cleanedTypingName := removeMinAndVersionNumbers(inferredTypingName)
		if safeName, ok := safeList[cleanedTypingName]; ok {
			fromFileNames = append(fromFileNames, safeName)
		}
	}
	if len(fromFileNames) > 0 {
		addInferredTypings(p, inferredTypings, fromFileNames, "Inferred typings from file names")
	}
	if hasJsxFile {
		p.Logf("Inferred 'react' typings due to presence of '.jsx' extension")
		addInferredTyping(inferredTypings, "react")
	}
}

func getManifestNamesFromDependencies(manifestNames []string, dependencies map[string]string) []string {
	if dependencies != nil {
		for dependency := range dependencies {
			manifestNames = append(manifestNames, dependency)
		}
	}
	return manifestNames
}

/**
 * Adds inferred typings from manifest/module pairs (think package.json + node_modules)
 *
 * @param projectRootPath is the path to the directory where to look for package.json, bower.json and other typing information
 * @param manifestName is the name of the manifest (package.json or bower.json)
 * @param modulesDirName is the directory name for modules (node_modules or bower_components). Should be lowercase!
 * @param filesToWatch are the files to watch for changes. We will push things into this array.
 */
func addTypingNamesAndGetFilesToWatch(
	p *Project,
	inferredTypings map[string]string,
	filesToWatch []string,
	projectRootPath string,
	manifestName string,
	modulesDirName string,
) []string {
	// First, we check the manifests themselves. They're not
	// _required_, but they allow us to do some filtering when dealing
	// with big flat dep directories.
	manifestPath := tspath.CombinePaths(projectRootPath, manifestName)
	var manifestTypingNames []string
	manifestContents, ok := p.FS().ReadFile(manifestPath)
	if ok {
		var manifest packagejson.DependencyFields
		filesToWatch = append(filesToWatch, manifestPath)
		// var manifest map[string]any
		err := json.Unmarshal([]byte(manifestContents), &manifest)
		if err != nil {
			manifestTypingNames = getManifestNamesFromDependencies(manifestTypingNames, manifest.Dependencies.Value)
			manifestTypingNames = getManifestNamesFromDependencies(manifestTypingNames, manifest.DevDependencies.Value)
			manifestTypingNames = getManifestNamesFromDependencies(manifestTypingNames, manifest.OptionalDependencies.Value)
			manifestTypingNames = getManifestNamesFromDependencies(manifestTypingNames, manifest.PeerDependencies.Value)
			addInferredTypings(p, inferredTypings, manifestTypingNames, "Typing names in '"+manifestPath+"' dependencies")
		}
	}

	// Now we scan the directories for typing information in
	// already-installed dependencies (if present). Note that this
	// step happens regardless of whether a manifest was present,
	// which is certainly a valid configuration, if an unusual one.
	packagesFolderPath := tspath.CombinePaths(projectRootPath, modulesDirName)
	filesToWatch = append(filesToWatch, packagesFolderPath)
	if !p.FS().DirectoryExists(packagesFolderPath) {
		return filesToWatch
	}

	// There's two cases we have to take into account here:
	// 1. If manifest is undefined, then we're not using a manifest.
	//    That means that we should scan _all_ dependencies at the top
	//    level of the modulesDir.
	// 2. If manifest is defined, then we can do some special
	//    filtering to reduce the amount of scanning we need to do.
	//
	// Previous versions of this algorithm checked for a `_requiredBy`
	// field in the package.json, but that field is only present in
	// `npm@>=3 <7`.

	// Package names that do **not** provide their own typings, so
	// we'll look them up.
	var packageNames []string

	var dependencyManifestNames []string
	if len(manifestTypingNames) > 0 {
		// This is #1 described above.
		for _, typingName := range manifestTypingNames {
			dependencyManifestNames = append(dependencyManifestNames, tspath.CombinePaths(packagesFolderPath, typingName, manifestName))
		}
	} else {
		// And #2. Depth = 3 because scoped packages look like `node_modules/@foo/bar/package.json`
		depth := 3
		for _, manifestPath := range vfs.ReadDirectory(p.FS(), p.GetCurrentDirectory(), packagesFolderPath, []string{tspath.ExtensionJson}, nil, nil, &depth) {
			if tspath.GetBaseFileName(manifestPath) != manifestName {
				continue
			}

			// It's ok to treat
			// `node_modules/@foo/bar/package.json` as a manifest,
			// but not `node_modules/jquery/nested/package.json`.
			// We only assume depth 3 is ok for formally scoped
			// packages. So that needs this dance here.

			pathComponents := tspath.GetPathComponents(manifestPath, "")
			lenPathComponents := len(pathComponents)
			isScoped := rune(pathComponents[lenPathComponents-3][0]) == '@'

			if isScoped && tspath.ToFileNameLowerCase(pathComponents[lenPathComponents-4]) == modulesDirName || // `node_modules/@foo/bar`
				!isScoped && tspath.ToFileNameLowerCase(pathComponents[lenPathComponents-3]) == modulesDirName { // `node_modules/foo`
				dependencyManifestNames = append(dependencyManifestNames, manifestPath)
			}
		}

	}

	p.Logf("Searching for typing names in %s; all files: %v", packagesFolderPath, dependencyManifestNames)

	// Once we have the names of things to look up, we iterate over
	// and either collect their included typings, or add them to the
	// list of typings we need to look up separately.
	for _, manifestPath := range dependencyManifestNames {
		manifestContents, ok := p.FS().ReadFile(manifestPath)
		if !ok {
			continue
		}
		manifest, err := packagejson.Parse([]byte(manifestContents))
		// If the package has its own d.ts typings, those will take precedence. Otherwise the package name will be used
		// to download d.ts files from DefinitelyTyped
		if err != nil || len(manifest.Name.Value) == 0 {
			continue
		}
		ownTypes := manifest.Types.Value
		if len(ownTypes) == 0 {
			ownTypes = manifest.Typings.Value
		}
		if len(ownTypes) != 0 {
			absolutePath := tspath.GetNormalizedAbsolutePath(ownTypes, tspath.GetDirectoryPath(manifestPath))
			if p.FS().FileExists(absolutePath) {
				p.Logf("    Package '%s' provides its own types.", manifest.Name.Value)
				inferredTypings[manifest.Name.Value] = absolutePath
			} else {
				p.Logf("    Package '%s' provides its own types but they are missing.", manifest.Name.Value)
			}
		} else {
			packageNames = append(packageNames, manifest.Name.Value)
		}
	}
	addInferredTypings(p, inferredTypings, packageNames, "    Found package names")
	return filesToWatch
}

/**
 * Takes a string like "jquery-min.4.2.3" and returns "jquery"
 *
 * @internal
 */
func removeMinAndVersionNumbers(fileName string) string {
	// We used to use the regex /[.-]((min)|(\d+(\.\d+)*))$/ and would just .replace it twice.
	// Unfortunately, that regex has O(n^2) performance because v8 doesn't match from the end of the string.
	// Instead, we now essentially scan the filename (backwards) ourselves.

	end := len(fileName)
	for pos := end - 1; pos > 0; pos-- {
		ch := rune(fileName[pos])
		if ch >= '0' && ch <= '9' {
			// Match a \d+ segment
			for {
				pos--
				ch = rune(fileName[pos])
				if pos <= 0 || ch < '0' || ch > '9' {
					break
				}
			}
		} else if pos > 4 && (ch == 'n' || ch == 'N') {
			// Looking for "min" or "min"
			// Already matched the 'n'
			pos--
			ch = rune(fileName[pos])
			if ch != 'i' && ch != 'I' {
				break
			}
			pos--
			ch = rune(fileName[pos])
			if ch != 'm' && ch != 'M' {
				break
			}
			pos--
			ch = rune(fileName[pos])
		} else {
			// This character is not part of either suffix pattern
			break
		}

		if ch != '-' && ch != '.' {
			break
		}
		end = pos
	}

	// end might be fileName.length, in which case this should internally no-op
	if end == len(fileName) {
		return fileName
	}
	return fileName[0:end]
}

type NameValidationResult int

const (
	NameOk NameValidationResult = iota
	EmptyName
	NameTooLong
	NameStartsWithDot
	NameStartsWithUnderscore
	NameContainsNonURISafeCharacters
)

const maxPackageNameLength = 214

/**
 * Validates package name using rules defined at https://docs.npmjs.com/files/package.json
 *
 * @internal
 */
func ValidatePackageName(packageName string) (result NameValidationResult, name string, isScopeName bool) {
	return validatePackageNameWorker(packageName /*supportScopedPackage*/, true)
}

var packageNameRegexp = regexp.MustCompile("^@([^/]+)/([^/]+)$")

func validatePackageNameWorker(packageName string, supportScopedPackage bool) (result NameValidationResult, name string, isScopeName bool) {
	packageNameLen := len(packageName)
	if packageNameLen == 0 {
		return EmptyName, "", false
	}
	if packageNameLen > maxPackageNameLength {
		return NameTooLong, "", false
	}
	firstChar := rune(packageName[0])
	if firstChar == '.' {
		return NameStartsWithDot, "", false
	}
	if firstChar == '_' {
		return NameStartsWithUnderscore, "", false
	}
	// check if name is scope package like: starts with @ and has one '/' in the middle
	// scoped packages are not currently supported
	if supportScopedPackage {
		matches := packageNameRegexp.FindStringSubmatch(packageName)
		if matches != nil {
			scopeResult, _, _ := validatePackageNameWorker(matches[1] /*supportScopedPackage*/, false)
			if scopeResult != NameOk {
				return scopeResult, matches[1], true
			}
			packageResult, _, _ := validatePackageNameWorker(matches[2] /*supportScopedPackage*/, false)
			if packageResult != NameOk {
				return packageResult, matches[2], false
			}
			return NameOk, "", false
		}
	}
	if url.QueryEscape(packageName) != packageName {
		return NameContainsNonURISafeCharacters, "", false
	}
	return NameOk, "", false
}

/** @internal */
func RenderPackageNameValidationFailure(typing string, result NameValidationResult, name string, isScopeName bool) string {
	var kind string
	if isScopeName {
		kind = "Scope"
	} else {
		kind = "Package"
	}
	if name == "" {
		name = typing
	}
	switch result {
	case EmptyName:
		return fmt.Sprintf("'%s':: %s name '%s' cannot be empty", typing, kind, name)
	case NameTooLong:
		return fmt.Sprintf("'%s':: %s name '%s' should be less than %d characters", typing, kind, name, maxPackageNameLength)
	case NameStartsWithDot:
		return fmt.Sprintf("'%s':: %s name '%s' cannot start with '.'", typing, kind, name)
	case NameStartsWithUnderscore:
		return fmt.Sprintf("'%s':: %s name '%s' cannot start with '_'", typing, kind, name)
	case NameContainsNonURISafeCharacters:
		return fmt.Sprintf("'%s':: %s name '%s' contains non URI safe characters", typing, kind, name)
	case NameOk:
		panic("Unexpected Ok result")
	default:
		panic("Unknown package name validation result")
	}
}
