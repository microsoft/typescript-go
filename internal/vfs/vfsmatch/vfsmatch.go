package vfsmatch

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/vfs"
)

// newNewMatch controls whether to use the regex-free glob matching implementation.
const newNewMatch = true

func ReadDirectory(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	if newNewMatch {
		return matchFilesNoRegex(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
	}
	return matchFiles(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// IsImplicitGlob checks if a path component is implicitly a glob.
// An "includes" path "foo" is implicitly a glob "foo/** /*" (without the space) if its last component has no extension,
// and does not contain any glob characters itself.
func IsImplicitGlob(lastPathComponent string) bool {
	return !strings.ContainsAny(lastPathComponent, ".*?")
}

type Usage string

const (
	UsageFiles       Usage = "files"
	UsageDirectories Usage = "directories"
	UsageExclude     Usage = "exclude"
)

// SpecMatcher is an interface for matching file paths against compiled glob patterns.
// It abstracts over both regex-based and regex-free implementations.
type SpecMatcher interface {
	// MatchString returns true if the given path matches the pattern.
	MatchString(path string) bool
}

// SpecMatchers is an interface for matching file paths against multiple compiled glob patterns.
// It can return the index of the matching pattern.
type SpecMatchers interface {
	// MatchIndex returns the index of the first matching pattern, or -1 if none match.
	MatchIndex(path string) int
	// Len returns the number of patterns.
	Len() int
}

// NewSpecMatcher creates a matcher for one or more glob specs.
// It returns a matcher that can test if paths match any of the patterns.
func NewSpecMatcher(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) SpecMatcher {
	if newNewMatch {
		if m := newGlobSpecMatcher(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
			return m
		}
		return nil
	}
	if m := newRegexSpecMatcher(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
		return m
	}
	return nil
}

// NewSingleSpecMatcher creates a matcher for a single glob spec.
// Returns nil if the spec compiles to an empty pattern (e.g., trailing ** for non-exclude).
func NewSingleSpecMatcher(spec string, basePath string, usage Usage, useCaseSensitiveFileNames bool) SpecMatcher {
	if newNewMatch {
		if m := newGlobSingleSpecMatcher(spec, basePath, usage, useCaseSensitiveFileNames); m != nil {
			return m
		}
		return nil
	}
	if m := newRegexSingleSpecMatcher(spec, basePath, usage, useCaseSensitiveFileNames); m != nil {
		return m
	}
	return nil
}

// NewSpecMatchers creates individual matchers for each spec, allowing lookup of which spec matched.
// Returns nil if no valid patterns could be compiled from the specs.
func NewSpecMatchers(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) SpecMatchers {
	if newNewMatch {
		if m := newGlobSpecMatchers(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
			return m
		}
		return nil
	}
	if m := newRegexSpecMatchers(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
		return m
	}
	return nil
}
