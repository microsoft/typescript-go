package modulespecifiers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func isNonGlobalAmbientModule(node *ast.Node) bool {
	return ast.IsModuleDeclaration(node) && ast.IsStringLiteral(node.Name())
}

func comparePathsByRedirectAndNumberOfDirectorySeparators(a ModulePath, b ModulePath) int {
	if a.IsRedirect == b.IsRedirect {
		return strings.Count(a.Path, "/") - strings.Count(b.Path, "/")
	}
	if a.IsRedirect {
		return 1
	}
	return -1
}

func pathIsBareSpecifier(path string) bool {
	return !tspath.PathIsAbsolute(path) && !tspath.PathIsRelative(path)
}

func isExcludedByRegex(moduleSpecifier string, excludes []string) bool {
	for _, pattern := range excludes {
		compiled, err := regexp2.Compile(pattern, regexp2.None)
		if err != nil {
			continue
		}
		match, _ := compiled.MatchString(moduleSpecifier)
		if match {
			return true
		}
	}
	return false
}

/**
 * Ensures a path is either absolute (prefixed with `/` or `c:`) or dot-relative (prefixed
 * with `./` or `../`) so as not to be confused with an unprefixed module name.
 *
 * ```ts
 * ensurePathIsNonModuleName("/path/to/file.ext") === "/path/to/file.ext"
 * ensurePathIsNonModuleName("./path/to/file.ext") === "./path/to/file.ext"
 * ensurePathIsNonModuleName("../path/to/file.ext") === "../path/to/file.ext"
 * ensurePathIsNonModuleName("path/to/file.ext") === "./path/to/file.ext"
 * ```
 *
 */
func ensurePathIsNonModuleName(path string) string {
	if pathIsBareSpecifier(path) {
		return "./" + path
	}
	return path
}

func getJsExtensionForDeclarationFileExtension(ext string) string {
	switch ext {
	case tspath.ExtensionDts:
		return tspath.ExtensionJs
	case tspath.ExtensionDmts:
		return tspath.ExtensionMjs
	case tspath.ExtensionDcts:
		return tspath.ExtensionCjs
	default:
		// .d.json.ts and the like
		return ext[len(".d") : len(ext)-len(tspath.ExtensionTs)]
	}
}

func getJSExtensionForFile(fileName string, options *core.CompilerOptions) string {
	result := tryGetJSExtensionForFile(fileName, options)
	if len(result) == 0 {
		panic(fmt.Sprintf("Extension %s is unsupported:: FileName:: %s", extensionFromPath(fileName), fileName))
	}
	return result
}

/**
 * Gets the extension from a path.
 * Path must have a valid extension.
 */
func extensionFromPath(path string) string {
	ext := tspath.TryGetExtensionFromPath(path)
	if len(ext) == 0 {
		panic(fmt.Sprintf("File %s has unknown extension.", path))
	}
	return ext
}

func tryGetJSExtensionForFile(fileName string, options *core.CompilerOptions) string {
	ext := tspath.TryGetExtensionFromPath(fileName)
	switch ext {
	case tspath.ExtensionTs, tspath.ExtensionDts:
		return tspath.ExtensionJs
	case tspath.ExtensionTsx:
		if options.Jsx == core.JsxEmitPreserve {
			return tspath.ExtensionJsx
		}
		return tspath.ExtensionJs
	case tspath.ExtensionJs, tspath.ExtensionJsx, tspath.ExtensionJson:
		return ext
	case tspath.ExtensionDmts, tspath.ExtensionMts, tspath.ExtensionMjs:
		return tspath.ExtensionMjs
	case tspath.ExtensionDcts, tspath.ExtensionCts, tspath.ExtensionCjs:
		return tspath.ExtensionCjs
	default:
		return ""
	}
}

func tryGetAnyFileFromPath(host ModuleSpecifierGenerationHost, path string) bool {
	// !!! TODO: shouldn't this use readdir instead of fileexists for perf?
	// We check all js, `node` and `json` extensions in addition to TS, since node module resolution would also choose those over the directory
	extGroups := tsoptions.GetSupportedExtensions(
		&core.CompilerOptions{
			AllowJs: core.TSTrue,
		},
		[]tsoptions.FileExtensionInfo{
			tsoptions.FileExtensionInfo{
				Extension:      "node",
				IsMixedContent: false,
				ScriptKind:     core.ScriptKindExternal,
			},
			tsoptions.FileExtensionInfo{
				Extension:      "json",
				IsMixedContent: false,
				ScriptKind:     core.ScriptKindJSON,
			},
		},
	)
	for _, exts := range extGroups {
		for _, e := range exts {
			fullPath := path + e
			if host.FileExists(fullPath) {
				return true
			}
		}
	}
	return false
}

func getPathsRelativeToRootDirs(path string, rootDirs []string, useCaseSensitiveFileNames bool) []string {
	var results []string
	for _, rootDir := range rootDirs {
		relativePath := getRelativePathIfInSameVolume(path, rootDir, useCaseSensitiveFileNames)
		if len(relativePath) > 0 && isPathRelativeToParent(relativePath) {
			results = append(results, relativePath)
		}
	}
	return results
}

func isPathRelativeToParent(path string) bool {
	return strings.HasPrefix(path, "..")
}

func getRelativePathIfInSameVolume(path string, directoryPath string, useCaseSensitiveFileNames bool) string {
	relativePath := tspath.GetRelativePathToDirectoryOrUrl(directoryPath, path, false, tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: useCaseSensitiveFileNames,
		CurrentDirectory:          directoryPath,
	})
	if tspath.IsRootedDiskPath(relativePath) {
		return ""
	}
	return relativePath
}

func getNearestAncestorDirectoryWithPackageJson(host ModuleSpecifierGenerationHost, dir string) string {
	// no fallback impl, required on host
	return host.GetNearestAncestorDirectoryWithPackageJson(dir)
}

func packageJsonPathsAreEqual(a string, b string, options tspath.ComparePathsOptions) bool {
	if a == b {
		return true
	}
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	return tspath.ComparePaths(a, b, options) == 0
}

func prefersTsExtension(allowedEndings []ModuleSpecifierEnding) bool {
	jsPriority := slices.Index(allowedEndings, ModuleSpecifierEndingJsExtension)
	tsPriority := slices.Index(allowedEndings, ModuleSpecifierEndingTsExtension)
	if tsPriority > -1 {
		return tsPriority < jsPriority
	}
	return false
}
