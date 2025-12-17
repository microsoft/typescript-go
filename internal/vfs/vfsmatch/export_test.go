package vfsmatch

import "github.com/microsoft/typescript-go/internal/vfs"

// Test-only exports for functions and types that are not part of the public API.

// ReadDirectoryOld is a test-only export for the regex-based implementation.
func ReadDirectoryOld(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return matchFiles(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// ReadDirectoryNew is a test-only export for the regex-free implementation.
func ReadDirectoryNew(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return matchFilesNoRegex(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// GlobPatternWrapper is a test-only wrapper for the unexported globPattern type.
type GlobPatternWrapper struct {
	pattern *globPattern
}

// Matches calls the unexported matches method on the wrapped globPattern.
func (w *GlobPatternWrapper) Matches(path string) bool {
	if w == nil || w.pattern == nil {
		return false
	}
	return w.pattern.matches(path)
}

// CompileGlobPattern is a test-only export for compiling glob patterns.
func CompileGlobPattern(spec string, basePath string, usage Usage, caseSensitive bool) *GlobPatternWrapper {
	p := compileGlobPattern(spec, basePath, usage, caseSensitive)
	if p == nil {
		return nil
	}
	return &GlobPatternWrapper{pattern: p}
}

// GetRegularExpressionForWildcard is a test-only export for getting the regex for wildcard specs.
func GetRegularExpressionForWildcard(specs []string, basePath string, usage Usage) string {
	return getRegularExpressionForWildcard(specs, basePath, usage)
}

// GetRegularExpressionsForWildcards is a test-only export for getting regexes for wildcard specs.
func GetRegularExpressionsForWildcards(specs []string, basePath string, usage Usage) []string {
	return getRegularExpressionsForWildcards(specs, basePath, usage)
}

// GetPatternFromSpec is a test-only export for getting a pattern from a spec.
func GetPatternFromSpec(spec string, basePath string, usage Usage) string {
	return getPatternFromSpec(spec, basePath, usage)
}
