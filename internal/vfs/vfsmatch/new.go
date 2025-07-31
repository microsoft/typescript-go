package vfsmatch

import (
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// isImplicitlyExcluded checks if a file or directory should be implicitly excluded
// based on TypeScript's default behavior (dotted files/folders and common package folders)
func isImplicitlyExcluded(name string, isDirectory bool) bool {
	// Exclude files/directories that start with a dot
	if strings.HasPrefix(name, ".") {
		return true
	}

	// For directories, exclude common package folders
	if isDirectory && slices.Contains(commonPackageFolders, name) {
		return true
	}

	return false
}

// shouldImplicitlyExcludeRelativePath checks if a relative path should be implicitly excluded
func shouldImplicitlyExcludeRelativePath(relativePath string) bool {
	if relativePath == "" {
		return false
	}

	// Split path into segments and check each segment
	segments := strings.Split(relativePath, "/")
	for _, segment := range segments {
		if isImplicitlyExcluded(segment, true) { // Check as directory since it's a path segment
			return true
		}
	}

	return false
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

	// Create relative pattern matchers
	includeMatchers := make([]globMatcher, len(includes))
	for i, include := range includes {
		includeMatchers[i] = globMatcherForPatternRelative(include, useCaseSensitiveFileNames)
	}

	excludeMatchers := make([]globMatcher, len(excludes))
	for i, exclude := range excludes {
		excludeMatchers[i] = globMatcherForPatternAbsolute(exclude, absolutePath, useCaseSensitiveFileNames)
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

	visitor := newRelativeGlobVisitor{
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		host:                      host,
		includeMatchers:           includeMatchers,
		excludeMatchers:           excludeMatchers,
		extensions:                extensions,
		results:                   results,
		visited:                   *collections.NewSetWithSizeHint[string](0),
		basePath:                  absolutePath,
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

// matchSegments recursively matches glob pattern segments against path segments
func (gm globMatcher) matchSegments(patternSegments []string, pathSegments []string, isDirectory bool) bool {
	pi, ti := 0, 0
	plen, tlen := len(patternSegments), len(pathSegments)

	// Special case for directory matching: if the path is a prefix of the pattern (ignoring final wildcards),
	// then it matches. This handles cases like pattern "LICENSE/**/*" matching directory "LICENSE/"
	if isDirectory && tlen < plen {
		// Check if path segments match the beginning of pattern segments
		matchesPrefix := true
		for i := 0; i < tlen && i < plen; i++ {
			if patternSegments[i] == "**" {
				// If we hit ** in the pattern, we're done - this directory could contain matching files
				break
			}
			if !gm.matchSegment(patternSegments[i], pathSegments[i]) {
				matchesPrefix = false
				break
			}
		}

		if matchesPrefix && tlen < plen {
			// Check if the remaining pattern segments are wildcards that could match files in this directory
			remainingPattern := patternSegments[tlen:]
			if len(remainingPattern) > 0 && (remainingPattern[0] == "**" ||
				(len(remainingPattern) >= 2 && remainingPattern[0] == "**" && remainingPattern[1] == "*")) {
				return true
			}
		}
	}

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

type newRelativeGlobVisitor struct {
	includeMatchers           []globMatcher
	excludeMatchers           []globMatcher
	extensions                []string
	useCaseSensitiveFileNames bool
	host                      vfs.FS
	visited                   collections.Set[string]
	results                   [][]string
	basePath                  string // The absolute base path for the search
}

func (v *newRelativeGlobVisitor) visitDirectory(path string, absolutePath string, depth *int) {
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

		// Convert to relative path for matching
		relativePath := absoluteName
		if strings.HasPrefix(absoluteName, v.basePath) {
			relativePath = absoluteName[len(v.basePath):]
			if strings.HasPrefix(relativePath, "/") {
				relativePath = relativePath[1:]
			}
		}

		// Apply implicit exclusions (dotted files and common package folders)
		if shouldImplicitlyExcludeRelativePath(relativePath) || isImplicitlyExcluded(current, false) {
			continue
		}

		excluded := false
		for _, excludeMatcher := range v.excludeMatchers {
			if excludeMatcher.matchesFileAbsolute(absoluteName) {
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
				if includeMatcher.matchesFileRelative(relativePath) {
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

		// Convert to relative path for matching
		relativePath := absoluteName
		if strings.HasPrefix(absoluteName, v.basePath) {
			relativePath = absoluteName[len(v.basePath):]
			if strings.HasPrefix(relativePath, "/") {
				relativePath = relativePath[1:]
			}
		}

		// Apply implicit exclusions (dotted directories and common package folders)
		if shouldImplicitlyExcludeRelativePath(relativePath) || isImplicitlyExcluded(current, true) {
			continue
		}

		shouldInclude := len(v.includeMatchers) == 0
		if !shouldInclude {
			for _, includeMatcher := range v.includeMatchers {
				if includeMatcher.couldMatchInSubdirectoryRelative(relativePath) {
					shouldInclude = true
					break
				}
			}
		}

		shouldExclude := false
		for _, excludeMatcher := range v.excludeMatchers {
			if excludeMatcher.matchesDirectoryAbsolute(absoluteName) {
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

func matchesExcludeNew(fileName string, excludeSpecs []string, currentDirectory string, useCaseSensitiveFileNames bool) bool {
	if len(excludeSpecs) == 0 {
		return false
	}

	// Convert fileName to relative path from currentDirectory for matching
	relativePath := fileName
	if strings.HasPrefix(fileName, currentDirectory) {
		relativePath = fileName[len(currentDirectory):]
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}
	}

	for _, excludeSpec := range excludeSpecs {
		// Special case: empty pattern matches everything (consistent with TypeScript behavior)
		if excludeSpec == "" {
			return true
		}

		matcher := globMatcherForPatternRelative(excludeSpec, useCaseSensitiveFileNames)
		if matcher.matchesFileRelative(relativePath) {
			return true
		}
		// Also check if it matches as a directory (for extensionless files)
		if !tspath.HasExtension(fileName) {
			relativePathWithSlash := relativePath
			if relativePathWithSlash != "" && !strings.HasSuffix(relativePathWithSlash, "/") {
				relativePathWithSlash += "/"
			}
			// Check if the file with trailing slash matches the pattern
			if matcher.matchesDirectoryRelative(relativePathWithSlash) {
				return true
			}
		}
	}
	return false
}

func matchesIncludeNew(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	if len(includeSpecs) == 0 {
		return false
	}

	// Convert fileName to relative path from basePath for matching
	relativePath := fileName
	if strings.HasPrefix(fileName, basePath) {
		relativePath = fileName[len(basePath):]
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}
	}

	for _, includeSpec := range includeSpecs {
		// Special case: empty pattern matches everything (consistent with TypeScript behavior)
		if includeSpec == "" {
			return true
		}

		matcher := globMatcherForPatternRelative(includeSpec, useCaseSensitiveFileNames)
		if matcher.matchesFileRelative(relativePath) {
			return true
		}
	}
	return false
}

func matchesIncludeWithJsonOnlyNew(fileName string, includeSpecs []string, basePath string, useCaseSensitiveFileNames bool) bool {
	if len(includeSpecs) == 0 {
		return false
	}

	// Convert fileName to relative path from basePath for matching
	relativePath := fileName
	if strings.HasPrefix(fileName, basePath) {
		relativePath = fileName[len(basePath):]
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}
	}

	// Special case: empty pattern matches everything (consistent with TypeScript behavior)
	for _, includeSpec := range includeSpecs {
		if includeSpec == "" {
			return true
		}
	}

	// Filter to only JSON include patterns
	jsonIncludes := core.Filter(includeSpecs, func(include string) bool {
		return strings.HasSuffix(include, tspath.ExtensionJson)
	})
	for _, includeSpec := range jsonIncludes {
		matcher := globMatcherForPatternRelative(includeSpec, useCaseSensitiveFileNames)
		if matcher.matchesFileRelative(relativePath) {
			return true
		}
	}
	return false
}

// globMatcherForPatternAbsolute creates a matcher for absolute pattern matching
// This is used for exclude patterns which are resolved against the absolutePath
func globMatcherForPatternAbsolute(pattern string, absolutePath string, useCaseSensitiveFileNames bool) globMatcher {
	// Resolve the pattern against the absolute path, similar to how getSubPatternFromSpec works
	// in the old implementation
	resolvedPattern := tspath.CombinePaths(absolutePath, pattern)
	resolvedPattern = tspath.NormalizePath(resolvedPattern)

	// Convert to relative pattern from the absolute path
	var relativePart string
	if strings.HasPrefix(resolvedPattern, absolutePath) {
		relativePart = resolvedPattern[len(absolutePath):]
		if strings.HasPrefix(relativePart, "/") {
			relativePart = relativePart[1:]
		}
	} else {
		// If the pattern doesn't start with absolutePath, use it as-is
		relativePart = pattern
		if strings.HasPrefix(relativePart, "/") {
			relativePart = relativePart[1:]
		}
	}

	// Parse the pattern as a relative path
	var segments []string
	if relativePart == "" {
		segments = []string{}
	} else {
		segments = strings.Split(relativePart, "/")
		// Remove empty segments
		filteredSegments := segments[:0]
		for _, seg := range segments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		segments = filteredSegments
	}

	// Handle implicit glob - if the last component has no extension and no wildcards, add **/*
	if len(segments) > 0 {
		lastComponent := segments[len(segments)-1]
		if IsImplicitGlob(lastComponent) {
			segments = append(segments, "**", "*")
		}
	}

	return globMatcher{
		pattern:                   pattern,
		basePath:                  absolutePath,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		segments:                  segments,
	}
}

// globMatcherForPatternRelative creates a matcher for relative pattern matching
func globMatcherForPatternRelative(pattern string, useCaseSensitiveFileNames bool) globMatcher {
	// Handle patterns starting with "./" - remove the leading "./"
	if strings.HasPrefix(pattern, "./") {
		pattern = pattern[2:]
	}

	// Parse the pattern as a relative path
	var segments []string
	if pattern == "" {
		segments = []string{}
	} else {
		segments = strings.Split(pattern, "/")
		// Remove empty segments
		filteredSegments := segments[:0]
		for _, seg := range segments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		segments = filteredSegments
	}

	// Handle implicit glob - if the last component has no extension and no wildcards, add **/*
	if len(segments) > 0 {
		lastComponent := segments[len(segments)-1]
		if IsImplicitGlob(lastComponent) {
			segments = append(segments, "**", "*")
		}
	}

	return globMatcher{
		pattern:                   pattern,
		basePath:                  "", // No base path for relative matching
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		segments:                  segments,
	}
}

// matchesFileAbsolute returns true if the given absolute file path matches the glob pattern
func (gm globMatcher) matchesFileAbsolute(absolutePath string) bool {
	// Special case for exclude patterns: if the pattern exactly matches the base path,
	// then it should exclude everything under that path (like "/apath" excluding "/apath/*")
	if gm.basePath != "" && len(gm.segments) == 0 {
		// Empty segments means the pattern exactly matched the base path
		// For excludes, this should match anything under the base path
		return strings.HasPrefix(absolutePath, gm.basePath) &&
			(absolutePath == gm.basePath || strings.HasPrefix(absolutePath, gm.basePath+"/"))
	}

	// Convert absolute path to relative path from the matcher's base path
	var relativePath string
	if gm.basePath != "" && strings.HasPrefix(absolutePath, gm.basePath) {
		relativePath = absolutePath[len(gm.basePath):]
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}
	} else {
		relativePath = absolutePath
	}

	return gm.matchesFileRelative(relativePath)
}

// matchesFileRelative returns true if the given relative file path matches the glob pattern
func (gm globMatcher) matchesFileRelative(relativePath string) bool {
	// Special case: empty pattern matches everything (consistent with TypeScript behavior)
	if gm.pattern == "" {
		return true
	}

	// Split the relative path into segments
	var pathSegments []string
	if relativePath == "" {
		pathSegments = []string{}
	} else {
		pathSegments = strings.Split(relativePath, "/")
		// Remove empty segments
		filteredSegments := pathSegments[:0]
		for _, seg := range pathSegments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		pathSegments = filteredSegments
	}

	return gm.matchSegments(gm.segments, pathSegments, false)
}

// matchesDirectoryAbsolute returns true if the given absolute directory path matches the glob pattern
func (gm globMatcher) matchesDirectoryAbsolute(absolutePath string) bool {
	// Special case for exclude patterns: if the pattern exactly matches the base path,
	// then it should exclude everything under that path (like "/apath" excluding "/apath/*")
	if gm.basePath != "" && len(gm.segments) == 0 {
		// Empty segments means the pattern exactly matched the base path
		// For excludes, this should match anything under the base path
		return strings.HasPrefix(absolutePath, gm.basePath) &&
			(absolutePath == gm.basePath || strings.HasPrefix(absolutePath, gm.basePath+"/"))
	}

	// Convert absolute path to relative path from the matcher's base path
	var relativePath string
	if gm.basePath != "" && strings.HasPrefix(absolutePath, gm.basePath) {
		relativePath = absolutePath[len(gm.basePath):]
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}
	} else {
		relativePath = absolutePath
	}

	return gm.matchesDirectoryRelative(relativePath)
}

// matchesDirectoryRelative returns true if the given relative directory path matches the glob pattern
func (gm globMatcher) matchesDirectoryRelative(relativePath string) bool {
	// Special case: empty pattern matches everything (consistent with TypeScript behavior)
	if gm.pattern == "" {
		return true
	}

	// Split the relative path into segments
	var pathSegments []string
	if relativePath == "" {
		pathSegments = []string{}
	} else {
		pathSegments = strings.Split(relativePath, "/")
		// Remove empty segments
		filteredSegments := pathSegments[:0]
		for _, seg := range pathSegments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		pathSegments = filteredSegments
	}

	return gm.matchSegments(gm.segments, pathSegments, true)
}

// couldMatchInSubdirectoryRelative returns true if this pattern could match files within the given relative directory
func (gm globMatcher) couldMatchInSubdirectoryRelative(relativePath string) bool {
	// Split the relative path into segments
	var pathSegments []string
	if relativePath == "" {
		pathSegments = []string{}
	} else {
		pathSegments = strings.Split(relativePath, "/")
		// Remove empty segments
		filteredSegments := pathSegments[:0]
		for _, seg := range pathSegments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		pathSegments = filteredSegments
	}

	return gm.couldMatchInSubdirectoryRecursive(gm.segments, pathSegments)
}
