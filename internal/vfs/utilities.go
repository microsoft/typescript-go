package vfs

import (
	"strings"
)

// An "includes" path "foo" is implicitly a glob "foo/** /*" (without the space) if its last component has no extension,
// and does not contain any glob characters itself.
func IsImplicitGlob(lastPathComponent string) bool {
	return !strings.ContainsAny(lastPathComponent, ".*?")
}

func ReadDirectory(host FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return MatchFilesNew(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}
