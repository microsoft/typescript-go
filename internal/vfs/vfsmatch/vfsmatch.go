package vfsmatch

import (
	"sort"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// An "includes" path "foo" is implicitly a glob "foo/** /*" (without the space) if its last component has no extension,
// and does not contain any glob characters itself.
func IsImplicitGlob(lastPathComponent string) bool {
	return !strings.ContainsAny(lastPathComponent, ".*?")
}

func ReadDirectory(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return readDirectoryNew(host, currentDir, path, extensions, excludes, includes, depth)
}

func MatchesExclude(fileName string, excludeSpecs []string, currentDirectory string, useCaseSensitiveFileNames bool) bool {
	return matchesExcludeNew(fileName, excludeSpecs, currentDirectory, useCaseSensitiveFileNames)
}

func MatchesInclude(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	return matchesIncludeNew(fileName, includeSpecs, basePath, useCaseSensitiveFileNames)
}

func MatchesIncludeWithJsonOnly(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	return matchesIncludeWithJsonOnlyNew(fileName, includeSpecs, basePath, useCaseSensitiveFileNames)
}

var (
	commonPackageFolders = []string{"node_modules", "bower_components", "jspm_packages"}
	wildcardCharCodes    = []rune{'*', '?'}
)

func getIncludeBasePath(absolute string) string {
	wildcardOffset := strings.IndexAny(absolute, string(wildcardCharCodes))
	if wildcardOffset < 0 {
		// No "*" or "?" in the path
		if !tspath.HasExtension(absolute) {
			return absolute
		} else {
			return tspath.RemoveTrailingDirectorySeparator(tspath.GetDirectoryPath(absolute))
		}
	}
	return absolute[:max(strings.LastIndex(absolute[:wildcardOffset], string(tspath.DirectorySeparator)), 0)]
}

// getBasePaths computes the unique non-wildcard base paths amongst the provided include patterns.
func getBasePaths(path string, includes []string, useCaseSensitiveFileNames bool) []string {
	// Storage for our results in the form of literal paths (e.g. the paths as written by the user).
	basePaths := []string{path}

	if len(includes) > 0 {
		// Storage for literal base paths amongst the include patterns.
		includeBasePaths := []string{}
		for _, include := range includes {
			// We also need to check the relative paths by converting them to absolute and normalizing
			// in case they escape the base path (e.g "..\somedirectory")
			var absolute string
			if tspath.IsRootedDiskPath(include) {
				absolute = include
			} else {
				absolute = tspath.NormalizePath(tspath.CombinePaths(path, include))
			}
			// Append the literal and canonical candidate base paths.
			includeBasePaths = append(includeBasePaths, getIncludeBasePath(absolute))
		}

		// Sort the offsets array using either the literal or canonical path representations.
		stringComparer := stringutil.GetStringComparer(!useCaseSensitiveFileNames)
		sort.SliceStable(includeBasePaths, func(i, j int) bool {
			return stringComparer(includeBasePaths[i], includeBasePaths[j]) < 0
		})

		// Iterate over each include base path and include unique base paths that are not a
		// subpath of an existing base path
		for _, includeBasePath := range includeBasePaths {
			if core.Every(basePaths, func(basepath string) bool {
				return !tspath.ContainsPath(basepath, includeBasePath, tspath.ComparePathsOptions{CurrentDirectory: path, UseCaseSensitiveFileNames: !useCaseSensitiveFileNames})
			}) {
				basePaths = append(basePaths, includeBasePath)
			}
		}
	}

	return basePaths
}
