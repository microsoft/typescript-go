package vfsmatch

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// Cache for normalized path components to avoid repeated allocations
type pathCache struct {
	cache map[string][]string
}

func newPathCache() *pathCache {
	return &pathCache{
		cache: make(map[string][]string),
	}
}

func (pc *pathCache) getNormalizedPathComponents(path string) []string {
	if components, exists := pc.cache[path]; exists {
		return components
	}

	components := tspath.GetNormalizedPathComponents(path, "")
	pc.cache[path] = components
	return components
}

// matchFilesNew is the regex-free implementation of file matching
func matchFilesNew(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host vfs.FS) []string {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	basePaths := getBasePaths(path, includes, useCaseSensitiveFileNames)

	// If no base paths found, return nil (consistent with original implementation)
	if len(basePaths) == 0 {
		return nil
	}

	// Create a shared path cache for this operation
	pathCache := newPathCache()

	// Prepare matchers for includes and excludes
	includeMatchers := make([]globMatcher, len(includes))
	for i, include := range includes {
		includeMatchers[i] = newGlobMatcher(include, absolutePath, useCaseSensitiveFileNames, pathCache)
	}

	excludeMatchers := make([]globMatcher, len(excludes))
	for i, exclude := range excludes {
		excludeMatchers[i] = newGlobMatcher(exclude, absolutePath, useCaseSensitiveFileNames, pathCache)
	}

	// Associate an array of results with each include matcher. This keeps results in order of the "include" order.
	// If there are no "includes", then just put everything in results[0].
	var results [][]string
	if len(includeMatchers) > 0 {
		tempResults := make([][]string, len(includeMatchers))
		for i := range includeMatchers {
			tempResults[i] = []string{}
		}
		results = tempResults
	} else {
		results = [][]string{{}}
	}

	visitor := newGlobVisitor{
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		host:                      host,
		includeMatchers:           includeMatchers,
		excludeMatchers:           excludeMatchers,
		extensions:                extensions,
		results:                   results,
		visited:                   *collections.NewSetWithSizeHint[string](0),
		pathCache:                 pathCache,
	}

	for _, basePath := range basePaths {
		visitor.visitDirectory(basePath, tspath.CombinePaths(currentDirectory, basePath), depth)
	}

	flattened := core.Flatten(results)
	if len(flattened) == 0 {
		return nil // Consistent with original implementation
	}
	return flattened
}

// globMatcher represents a glob pattern matcher without using regex
type globMatcher struct {
	pattern                   string
	basePath                  string
	useCaseSensitiveFileNames bool
	segments                  []string
	pathCache                 *pathCache
}

// newGlobMatcher creates a new glob matcher for the given pattern
func newGlobMatcher(pattern string, basePath string, useCaseSensitiveFileNames bool, pathCache *pathCache) globMatcher {
	// Convert pattern to absolute path if it's relative
	var absolutePattern string
	if tspath.IsRootedDiskPath(pattern) {
		absolutePattern = pattern
	} else {
		absolutePattern = tspath.NormalizePath(tspath.CombinePaths(basePath, pattern))
	}

	// Split into path segments - use cache to avoid repeated calls
	segments := pathCache.getNormalizedPathComponents(absolutePattern)
	// Remove the empty root component
	if len(segments) > 0 && segments[0] == "" {
		segments = segments[1:]
	}

	// Handle implicit glob - if the last component has no extension and no wildcards, add **/*
	if len(segments) > 0 {
		lastComponent := segments[len(segments)-1]
		if IsImplicitGlob(lastComponent) {
			segments = append(segments, "**", "*")
		}
	}

	return globMatcher{
		pattern:                   absolutePattern,
		basePath:                  basePath,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		segments:                  segments,
		pathCache:                 pathCache,
	}
}

// newGlobMatcherOld creates a new glob matcher for the given pattern (for backwards compatibility)
func newGlobMatcherOld(pattern string, basePath string, useCaseSensitiveFileNames bool) globMatcher {
	// Create a temporary path cache for old implementation
	tempCache := newPathCache()
	return newGlobMatcher(pattern, basePath, useCaseSensitiveFileNames, tempCache)
}

// matchesFile returns true if the given absolute file path matches the glob pattern
func (gm globMatcher) matchesFile(absolutePath string) bool {
	return gm.matchesPath(absolutePath, false)
}

// matchesDirectory returns true if the given absolute directory path matches the glob pattern
func (gm globMatcher) matchesDirectory(absolutePath string) bool {
	return gm.matchesPath(absolutePath, true)
}

// couldMatchInSubdirectory returns true if this pattern could match files within the given directory
func (gm globMatcher) couldMatchInSubdirectory(absolutePath string) bool {
	pathSegments := gm.pathCache.getNormalizedPathComponents(absolutePath)
	// Remove the empty root component
	if len(pathSegments) > 0 && pathSegments[0] == "" {
		pathSegments = pathSegments[1:]
	}

	return gm.couldMatchInSubdirectoryRecursive(gm.segments, pathSegments)
}

// couldMatchInSubdirectoryRecursive checks if the pattern could match files under the given path
func (gm globMatcher) couldMatchInSubdirectoryRecursive(patternSegments []string, pathSegments []string) bool {
	if len(patternSegments) == 0 {
		return false
	}

	pattern := patternSegments[0]
	remainingPattern := patternSegments[1:]

	if pattern == "**" {
		// Double asterisk can match anywhere
		return true
	}

	if len(pathSegments) == 0 {
		// We've run out of path but still have pattern segments
		// This means we could match files in the current directory
		return true
	}

	pathSegment := pathSegments[0]
	remainingPath := pathSegments[1:]

	// Check if this segment matches
	if gm.matchSegment(pattern, pathSegment) {
		// If we match and have more pattern segments, continue
		if len(remainingPattern) > 0 {
			return gm.couldMatchInSubdirectoryRecursive(remainingPattern, remainingPath)
		}
		// If no more pattern segments, we could match files in this directory
		return true
	}

	return false
}

// matchesPath performs the actual glob matching logic
func (gm globMatcher) matchesPath(absolutePath string, isDirectory bool) bool {
	pathSegments := gm.pathCache.getNormalizedPathComponents(absolutePath)
	// Remove the empty root component
	if len(pathSegments) > 0 && pathSegments[0] == "" {
		pathSegments = pathSegments[1:]
	}

	return gm.matchSegments(gm.segments, pathSegments, isDirectory)
}

// matchSegments recursively matches glob pattern segments against path segments
func (gm globMatcher) matchSegments(patternSegments []string, pathSegments []string, isDirectory bool) bool {
	pi, ti := 0, 0
	plen, tlen := len(patternSegments), len(pathSegments)
	for pi < plen {
		pattern := patternSegments[pi]
		if pattern == "**" {
			// Try matching remaining pattern at current position
			if gm.matchSegments(patternSegments[pi+1:], pathSegments[ti:], isDirectory) {
				return true
			}
			// Try consuming one path segment and continue with **
			for ti < tlen && (isDirectory || tlen-ti > 1) {
				ti++
				if gm.matchSegments(patternSegments[pi+1:], pathSegments[ti:], isDirectory) {
					return true
				}
			}
			return false
		}
		if ti >= tlen {
			return false
		}
		pathSegment := pathSegments[ti]
		isFinalSegment := (pi == plen-1) && (ti == tlen-1)
		isFileSegment := !isDirectory && isFinalSegment
		var segmentMatches bool
		if isFileSegment {
			segmentMatches = gm.matchSegmentForFile(pattern, pathSegment)
		} else {
			segmentMatches = gm.matchSegment(pattern, pathSegment)
		}
		if !segmentMatches {
			return false
		}
		pi++
		ti++
	}
	return ti == tlen
}

// matchSegment matches a single glob pattern segment against a path segment
func (gm globMatcher) matchSegment(pattern, segment string) bool {
	// Handle case sensitivity
	if !gm.useCaseSensitiveFileNames {
		pattern = strings.ToLower(pattern)
		segment = strings.ToLower(segment)
	}

	return gm.matchGlobPattern(pattern, segment, false)
}

func (gm globMatcher) matchSegmentForFile(pattern, segment string) bool {
	// Handle case sensitivity
	if !gm.useCaseSensitiveFileNames {
		pattern = strings.ToLower(pattern)
		segment = strings.ToLower(segment)
	}

	return gm.matchGlobPattern(pattern, segment, true)
}

// matchGlobPattern implements glob pattern matching for a single segment
func (gm globMatcher) matchGlobPattern(pattern, text string, isFileMatch bool) bool {
	pi, ti := 0, 0
	starIdx, match := -1, 0

	for ti < len(text) {
		if pi < len(pattern) && (pattern[pi] == '?' || pattern[pi] == text[ti]) {
			pi++
			ti++
		} else if pi < len(pattern) && pattern[pi] == '*' {
			// For file matching, * should not match .min.js files UNLESS the pattern explicitly ends with .min.js
			if isFileMatch && strings.HasSuffix(text, ".min.js") && !strings.HasSuffix(pattern, ".min.js") {
				return false
			}
			starIdx = pi
			match = ti
			pi++
		} else if starIdx != -1 {
			pi = starIdx + 1
			match++
			ti = match
		} else {
			return false
		}
	}

	// Handle remaining '*' in pattern
	for pi < len(pattern) && pattern[pi] == '*' {
		pi++
	}

	return pi == len(pattern)
}

type newGlobVisitor struct {
	includeMatchers           []globMatcher
	excludeMatchers           []globMatcher
	extensions                []string
	useCaseSensitiveFileNames bool
	host                      vfs.FS
	visited                   collections.Set[string]
	results                   [][]string
	pathCache                 *pathCache
}

func (v *newGlobVisitor) visitDirectory(path string, absolutePath string, depth *int) {
	canonicalPath := tspath.GetCanonicalFileName(absolutePath, v.useCaseSensitiveFileNames)
	if v.visited.Has(canonicalPath) {
		return
	}
	v.visited.Add(canonicalPath)

	systemEntries := v.host.GetAccessibleEntries(absolutePath)
	files := systemEntries.Files
	directories := systemEntries.Directories

	// Preallocate local buffers for results
	var localResults [][]string
	if len(v.includeMatchers) > 0 {
		localResults = make([][]string, len(v.includeMatchers))
		for i := range localResults {
			localResults[i] = make([]string, 0, len(files)/len(v.includeMatchers)+1)
		}
	} else {
		localResults = [][]string{make([]string, 0, len(files))}
	}
	for _, current := range files {
		if len(current) > 0 && current[0] == '.' {
			continue
		}
		var nameBuilder, absBuilder strings.Builder
		nameBuilder.Grow(len(path) + len(current) + 2)
		absBuilder.Grow(len(absolutePath) + len(current) + 2)
		if path == "" {
			nameBuilder.WriteString(current)
		} else {
			nameBuilder.WriteString(path)
			if path[len(path)-1] != '/' {
				nameBuilder.WriteByte('/')
			}
			nameBuilder.WriteString(current)
		}
		if absolutePath == "" {
			absBuilder.WriteString(current)
		} else {
			absBuilder.WriteString(absolutePath)
			if absolutePath[len(absolutePath)-1] != '/' {
				absBuilder.WriteByte('/')
			}
			absBuilder.WriteString(current)
		}
		name := nameBuilder.String()
		absoluteName := absBuilder.String()
		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(name, v.extensions) {
			continue
		}
		excluded := false
		for _, excludeMatcher := range v.excludeMatchers {
			if excludeMatcher.matchesFile(absoluteName) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		if len(v.includeMatchers) == 0 {
			localResults[0] = append(localResults[0], name)
		} else {
			for i, includeMatcher := range v.includeMatchers {
				if includeMatcher.matchesFile(absoluteName) {
					localResults[i] = append(localResults[i], name)
					break
				}
			}
		}
	}
	// Merge local buffers into main results
	for i := range localResults {
		v.results[i] = append(v.results[i], localResults[i]...)
	}
	if depth != nil {
		newDepth := *depth - 1
		if newDepth == 0 {
			return
		}
		depth = &newDepth
	}
	for _, current := range directories {
		if len(current) > 0 && current[0] == '.' {
			continue
		}
		isCommonPackageFolder := false
		for _, pkg := range commonPackageFolders {
			if current == pkg {
				isCommonPackageFolder = true
				break
			}
		}
		if isCommonPackageFolder {
			continue
		}
		var nameBuilder, absBuilder strings.Builder
		nameBuilder.Grow(len(path) + len(current) + 2)
		absBuilder.Grow(len(absolutePath) + len(current) + 2)
		if path == "" {
			nameBuilder.WriteString(current)
		} else {
			nameBuilder.WriteString(path)
			if path[len(path)-1] != '/' {
				nameBuilder.WriteByte('/')
			}
			nameBuilder.WriteString(current)
		}
		if absolutePath == "" {
			absBuilder.WriteString(current)
		} else {
			absBuilder.WriteString(absolutePath)
			if absolutePath[len(absolutePath)-1] != '/' {
				absBuilder.WriteByte('/')
			}
			absBuilder.WriteString(current)
		}
		name := nameBuilder.String()
		absoluteName := absBuilder.String()
		shouldInclude := len(v.includeMatchers) == 0
		if !shouldInclude {
			for _, includeMatcher := range v.includeMatchers {
				if includeMatcher.couldMatchInSubdirectory(absoluteName) {
					shouldInclude = true
					break
				}
			}
		}
		shouldExclude := false
		for _, excludeMatcher := range v.excludeMatchers {
			if excludeMatcher.matchesDirectory(absoluteName) {
				shouldExclude = true
				break
			}
		}
		if shouldInclude && !shouldExclude {
			v.visitDirectory(name, absoluteName, depth)
		}
	}
}

func readDirectoryNew(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return matchFilesNew(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// MatchesExclude checks if a file matches any of the exclude patterns using glob matching (no regexp2)
func MatchesExclude(fileName string, excludeSpecs []string, currentDirectory string, useCaseSensitiveFileNames bool) bool {
	if len(excludeSpecs) == 0 {
		return false
	}
	for _, excludeSpec := range excludeSpecs {
		matcher := GlobMatcherForPattern(excludeSpec, currentDirectory, useCaseSensitiveFileNames)
		if matcher.matchesFile(fileName) {
			return true
		}
		// Also check if it matches as a directory (for extensionless files)
		if !tspath.HasExtension(fileName) {
			if matcher.matchesDirectory(tspath.EnsureTrailingDirectorySeparator(fileName)) {
				return true
			}
		}
	}
	return false
}

// MatchesInclude checks if a file matches any of the include patterns using glob matching (no regexp2)
func MatchesInclude(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	if len(includeSpecs) == 0 {
		return false
	}
	for _, includeSpec := range includeSpecs {
		matcher := GlobMatcherForPattern(includeSpec, basePath, useCaseSensitiveFileNames)
		if matcher.matchesFile(fileName) {
			return true
		}
	}
	return false
}

// MatchesIncludeWithJsonOnly checks if a file matches any of the JSON-only include patterns using glob matching (no regexp2)
func MatchesIncludeWithJsonOnly(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	if len(includeSpecs) == 0 {
		return false
	}
	// Filter to only JSON include patterns
	jsonIncludes := core.Filter(includeSpecs, func(include string) bool {
		return strings.HasSuffix(include, tspath.ExtensionJson)
	})
	for _, includeSpec := range jsonIncludes {
		matcher := GlobMatcherForPattern(includeSpec, basePath, useCaseSensitiveFileNames)
		if matcher.matchesFile(fileName) {
			return true
		}
	}
	return false
}

// GlobMatcherForPattern is an exported wrapper for newGlobMatcher for use outside this file
func GlobMatcherForPattern(pattern string, basePath string, useCaseSensitiveFileNames bool) globMatcher {
	tempCache := newPathCache()
	return newGlobMatcher(pattern, basePath, useCaseSensitiveFileNames, tempCache)
}
