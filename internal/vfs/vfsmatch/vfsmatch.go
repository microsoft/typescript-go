package vfsmatch

import (
	"math"
	"sort"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -type=Usage -trimprefix=Usage -output=stringer_generated.go
//go:generate go tool mvdan.cc/gofumpt -w stringer_generated.go

type Usage int8

const (
	UsageFiles Usage = iota
	UsageDirectories
	UsageExclude
)

// UnlimitedDepth can be passed as the depth argument to indicate there is no depth limit.
const UnlimitedDepth = math.MaxInt

const newMatch = true

func ReadDirectory(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth int) []string {
	if newMatch {
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

// SpecMatcher is an interface for matching file paths against compiled glob patterns.
type SpecMatcher interface {
	// MatchString returns true if the given path matches the pattern.
	MatchString(path string) bool
}

// SpecMatchers is an interface for matching file paths against multiple compiled glob patterns.
// It can return the index of the matching pattern.
type SpecMatchers interface {
	// MatchIndex returns the index of the first matching pattern, or -1 if none match.
	MatchIndex(path string) int
}

// NewSpecMatcher creates a matcher for one or more glob specs.
// It returns a matcher that can test if paths match any of the patterns.
func NewSpecMatcher(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) SpecMatcher {
	if newMatch {
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
	if newMatch {
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
	if newMatch {
		if m := newGlobSpecMatcher(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
			return m
		}
		return nil
	}
	if m := newRegexSpecMatchers(specs, basePath, usage, useCaseSensitiveFileNames); m != nil {
		return m
	}
	return nil
}

var wildcardCharCodes = []rune{'*', '?'}

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
