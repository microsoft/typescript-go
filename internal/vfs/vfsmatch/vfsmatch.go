package vfsmatch

import (
	"strings"

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
