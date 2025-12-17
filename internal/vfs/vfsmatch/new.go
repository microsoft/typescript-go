package vfsmatch

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// globPattern represents a compiled glob pattern for matching file paths.
// It stores the pattern components for efficient matching without using regex.
type globPattern struct {
	// The original pattern specification
	spec string
	// The base path from which the pattern was derived
	basePath string
	// The normalized path components to match
	components []patternComponent
	// Whether this is an exclude pattern (affects matching behavior)
	isExclude bool
	// Whether pattern matching should be case-sensitive
	caseSensitive bool
	// For files patterns, exclude .min.js by default
	excludeMinJs bool
}

// patternComponent represents a single segment of a glob pattern
type patternComponent struct {
	// Whether this component is a ** wildcard
	isDoubleAsterisk bool
	// The literal text if not a wildcard pattern
	literal string
	// Whether this component contains wildcards
	hasWildcards bool
	// Parsed wildcard segments for matching
	segments []patternSegment
	// For include patterns (not exclude), implicitly exclude common package folders
	implicitlyExcludePackages bool
}

// patternSegment represents a parsed segment within a component
type patternSegment struct {
	kind    segmentKind
	literal string
}

type segmentKind int

const (
	segmentLiteral  segmentKind = iota
	segmentStar                 // * - matches any chars except /
	segmentQuestion             // ? - matches single char except /
)

// compileGlobPattern compiles a glob spec into a globPattern for matching.
func compileGlobPattern(spec string, basePath string, usage Usage, caseSensitive bool) *globPattern {
	components := tspath.GetNormalizedPathComponents(spec, basePath)
	lastComponent := core.LastOrNil(components)

	// If the last component is ** and this is not an exclude pattern, return nil
	// (such patterns match nothing)
	if usage != UsageExclude && lastComponent == "**" {
		return nil
	}

	// Remove trailing separator from root component
	components[0] = tspath.RemoveTrailingDirectorySeparator(components[0])

	// Handle implicit glob (directories become dir/**/*)
	if IsImplicitGlob(lastComponent) {
		components = append(components, "**", "*")
	}

	pattern := &globPattern{
		spec:          spec,
		basePath:      basePath,
		isExclude:     usage == UsageExclude,
		caseSensitive: caseSensitive,
		excludeMinJs:  usage == UsageFiles,
	}

	for _, comp := range components {
		pc := patternComponent{}

		if comp == "**" {
			pc.isDoubleAsterisk = true
		} else {
			pc.hasWildcards = strings.ContainsAny(comp, "*?")

			if pc.hasWildcards {
				pc.segments = parsePatternSegments(comp)
				// For non-exclude patterns with wildcards, implicitly exclude common package folders
				if usage != UsageExclude {
					pc.implicitlyExcludePackages = true
				}
			} else {
				pc.literal = comp
			}
		}

		pattern.components = append(pattern.components, pc)
	}

	return pattern
}

// parsePatternSegments breaks a component with wildcards into segments
func parsePatternSegments(comp string) []patternSegment {
	var segments []patternSegment
	var current strings.Builder

	for i := range len(comp) {
		switch comp[i] {
		case '*':
			if current.Len() > 0 {
				segments = append(segments, patternSegment{kind: segmentLiteral, literal: current.String()})
				current.Reset()
			}
			segments = append(segments, patternSegment{kind: segmentStar})
		case '?':
			if current.Len() > 0 {
				segments = append(segments, patternSegment{kind: segmentLiteral, literal: current.String()})
				current.Reset()
			}
			segments = append(segments, patternSegment{kind: segmentQuestion})
		default:
			current.WriteByte(comp[i])
		}
	}

	if current.Len() > 0 {
		segments = append(segments, patternSegment{kind: segmentLiteral, literal: current.String()})
	}

	return segments
}

// matches checks if the given path matches this glob pattern.
func (p *globPattern) matches(path string) bool {
	if p == nil {
		return false
	}
	return p.matchPathWorker(path, 0, 0, false, false)
}

// matchesPrefix checks if the given directory path could potentially match files under it.
// This is used for directory filtering during traversal.
func (p *globPattern) matchesPrefix(path string) bool {
	if p == nil {
		return false
	}
	return p.matchPathWorker(path, 0, 0, false, true)
}

// nextPathComponent extracts the next path component from path starting at offset.
// Returns the component, the offset after this component (pointing to char after '/' or len(path)), and whether a component was found.
func nextPathComponent(path string, offset int) (component string, nextOffset int, found bool) {
	if offset >= len(path) {
		return "", offset, false
	}

	// Handle leading slash for absolute paths - return empty string for root
	if offset == 0 && path[0] == '/' {
		return "", 1, true
	}

	// Skip any leading slashes (for cases like after root)
	for offset < len(path) && path[offset] == '/' {
		offset++
	}

	if offset >= len(path) {
		return "", offset, false
	}

	// Find the end of this component using optimized byte search
	remaining := path[offset:]
	idx := strings.IndexByte(remaining, '/')
	if idx < 0 {
		// No more slashes, rest of path is the component
		return remaining, len(path), true
	}
	return remaining[:idx], offset + idx, true
}

// matchPathWorker is the unified path matching function.
// When prefixMatch is true, it checks if the path could be a prefix of a matching path.
// When prefixMatch is false, it checks if the path fully matches the pattern.
func (p *globPattern) matchPathWorker(path string, pathOffset int, patternIdx int, inDoubleAsterisk bool, prefixMatch bool) bool {
	for {
		// Get the next path component
		pathComp, nextPathOffset, hasMore := nextPathComponent(path, pathOffset)

		// If we've consumed all path components
		if !hasMore {
			if prefixMatch {
				// For prefix matching, any prefix could match
				return true
			}
			// For full matching, check remaining pattern components
			if p.isExclude {
				return p.isImplicitGlobSuffix(patternIdx)
			}
			// Check if remaining pattern components are all optional (** only)
			for i := patternIdx; i < len(p.components); i++ {
				if !p.components[i].isDoubleAsterisk {
					return false
				}
			}
			return true
		}

		// If we've consumed all pattern components
		if patternIdx >= len(p.components) {
			if prefixMatch {
				// For prefix matching, no more matches possible
				return false
			}
			// For full matching with exclude patterns, we can match a prefix
			return p.isExclude
		}

		pc := p.components[patternIdx]

		if pc.isDoubleAsterisk {
			// ** can match zero or more directory levels
			// First, try matching zero directories (skip the **) - this requires recursion
			if p.matchPathWorker(path, pathOffset, patternIdx+1, true, prefixMatch) {
				return true
			}

			// For include patterns, ** should not match directories starting with . or common package folders
			if !p.isExclude {
				if len(pathComp) > 0 && pathComp[0] == '.' {
					return false
				}
				if isCommonPackageFolder(pathComp) {
					return false
				}
			}

			// Match current component with ** and continue (iterate instead of recurse)
			pathOffset = nextPathOffset
			inDoubleAsterisk = true
			continue
		}

		// Check implicit package folder exclusion
		if pc.implicitlyExcludePackages && !p.isExclude && isCommonPackageFolder(pathComp) {
			return false
		}

		// Match current component
		if !p.matchComponent(pc, pathComp, inDoubleAsterisk) {
			return false
		}

		// Continue to next components (iterate instead of recurse)
		pathOffset = nextPathOffset
		patternIdx++
		inDoubleAsterisk = false
	}
}

// matchComponent matches a single path component against a pattern component
func (p *globPattern) matchComponent(pc patternComponent, pathComp string, afterDoubleAsterisk bool) bool {
	if pc.isDoubleAsterisk {
		// Should not happen here, handled separately
		return true
	}

	// If the pattern component has no wildcards, do literal comparison
	if !pc.hasWildcards {
		return p.stringsEqual(pc.literal, pathComp)
	}

	// Match with wildcards
	// Note: The check for dotted names after ** is handled in matchWildcardComponent
	// where we only reject if the pattern itself starts with a wildcard
	return p.matchWildcardComponent(pc.segments, pathComp)
}

// matchWildcardComponent matches a path component against wildcard segments
func (p *globPattern) matchWildcardComponent(segments []patternSegment, s string) bool {
	// For non-exclude patterns, if the segments start with * or ?,
	// the matched string cannot start with '.'
	if !p.isExclude && len(segments) > 0 && len(s) > 0 && s[0] == '.' {
		firstSeg := segments[0]
		if firstSeg.kind == segmentStar || firstSeg.kind == segmentQuestion {
			// Pattern starts with wildcard, so it cannot match a string starting with '.'
			return false
		}
	}

	// Fast path for common pattern: * followed by literal suffix (e.g., "*.ts")
	if len(segments) == 2 && segments[0].kind == segmentStar && segments[1].kind == segmentLiteral {
		suffix := segments[1].literal
		if len(s) < len(suffix) {
			return false
		}
		// Check that there are no slashes in what * would match
		prefixLen := len(s) - len(suffix)
		for i := range prefixLen {
			if s[i] == '/' {
				return false
			}
		}
		// Check suffix match
		sSuffix := s[prefixLen:]
		if !p.stringsEqual(suffix, sSuffix) {
			return false
		}
		return p.checkMinJsExclusion(s, segments)
	}

	if !p.matchSegments(segments, 0, s, 0) {
		return false
	}
	return p.checkMinJsExclusion(s, segments)
}

// checkMinJsExclusion returns true if the match should be allowed (not excluded).
// Returns false if this is a .min.js file that should be excluded.
func (p *globPattern) checkMinJsExclusion(filename string, segments []patternSegment) bool {
	if !p.excludeMinJs {
		return true
	}
	if !p.wouldMatchMinJs(filename) {
		return true
	}
	// Exclude .min.js unless pattern explicitly includes it
	return p.patternExplicitlyIncludesMinJs(segments)
}

func (p *globPattern) matchSegments(segments []patternSegment, segIdx int, s string, sIdx int) bool {
	// If we've processed all segments
	if segIdx >= len(segments) {
		return sIdx >= len(s)
	}

	seg := segments[segIdx]

	switch seg.kind {
	case segmentLiteral:
		// Must match the literal exactly
		if sIdx+len(seg.literal) > len(s) {
			return false
		}
		substr := s[sIdx : sIdx+len(seg.literal)]
		if !p.stringsEqual(seg.literal, substr) {
			return false
		}
		return p.matchSegments(segments, segIdx+1, s, sIdx+len(seg.literal))

	case segmentQuestion:
		// Must match exactly one character (not /)
		if sIdx >= len(s) {
			return false
		}
		if s[sIdx] == '/' {
			return false
		}
		return p.matchSegments(segments, segIdx+1, s, sIdx+1)

	case segmentStar:
		// Match zero or more characters (not /)
		// Try matching zero characters first
		if p.matchSegments(segments, segIdx+1, s, sIdx) {
			return true
		}
		// Try matching more characters
		for i := sIdx; i < len(s); i++ {
			if s[i] == '/' {
				break
			}
			if p.matchSegments(segments, segIdx+1, s, i+1) {
				return true
			}
		}
		return false
	}

	return false
}

// wouldMatchMinJs checks if the filename ends with .min.js (case-insensitive)
func (p *globPattern) wouldMatchMinJs(filename string) bool {
	// Check length first to avoid string operations
	const suffix = ".min.js"
	if len(filename) < len(suffix) {
		return false
	}
	// Get the last 7 characters and compare case-insensitively
	end := filename[len(filename)-len(suffix):]
	return strings.EqualFold(end, suffix)
}

// patternExplicitlyIncludesMinJs checks if the pattern explicitly includes .min.js
func (p *globPattern) patternExplicitlyIncludesMinJs(segments []patternSegment) bool {
	// Look for .min.js in the literal segments
	for _, seg := range segments {
		if seg.kind == segmentLiteral && strings.Contains(strings.ToLower(seg.literal), ".min.js") {
			return true
		}
	}
	return false
}

// isImplicitGlobSuffix checks if the remaining pattern components from patternIdx
// are the implicit glob suffix (** followed by *) or all **
func (p *globPattern) isImplicitGlobSuffix(patternIdx int) bool {
	remaining := len(p.components) - patternIdx
	if remaining == 0 {
		return true
	}
	// All remaining must be ** (can match zero components)
	// OR it's exactly **/* (the implicit glob pattern added for directories)
	allDoubleAsterisk := true
	for i := patternIdx; i < len(p.components); i++ {
		if !p.components[i].isDoubleAsterisk {
			allDoubleAsterisk = false
			break
		}
	}
	if allDoubleAsterisk {
		return true
	}
	// Check for exactly **/* pattern (implicit glob suffix)
	if remaining == 2 {
		if p.components[patternIdx].isDoubleAsterisk {
			last := p.components[patternIdx+1]
			// The last component must be a pure * wildcard (matching any filename)
			if last.hasWildcards && len(last.segments) == 1 && last.segments[0].kind == segmentStar {
				return true
			}
		}
	}
	return false
}

// stringsEqual compares two strings with case sensitivity based on pattern settings
func (p *globPattern) stringsEqual(a, b string) bool {
	if p.caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// isCommonPackageFolder checks if a directory name is a common package folder
func isCommonPackageFolder(name string) bool {
	// Quick length check to avoid EqualFold for most cases
	switch len(name) {
	case len("node_modules"):
		return strings.EqualFold(name, "node_modules")
	case len("bower_components"):
		return strings.EqualFold(name, "bower_components")
	case len("jspm_packages"):
		return strings.EqualFold(name, "jspm_packages")
	default:
		return false
	}
}

// globMatcher holds compiled glob patterns for matching files.
type globMatcher struct {
	includePatterns []*globPattern
	excludePatterns []*globPattern
	caseSensitive   bool
	// hadIncludes tracks whether any include specs were provided (even if they compiled to nothing)
	hadIncludes bool
}

// newGlobMatcher creates a new globMatcher from include and exclude specs.
func newGlobMatcher(includes []string, excludes []string, basePath string, caseSensitive bool, usage Usage) *globMatcher {
	m := &globMatcher{
		caseSensitive: caseSensitive,
		hadIncludes:   len(includes) > 0,
	}

	for _, spec := range includes {
		if pattern := compileGlobPattern(spec, basePath, usage, caseSensitive); pattern != nil {
			m.includePatterns = append(m.includePatterns, pattern)
		}
	}

	for _, spec := range excludes {
		if pattern := compileGlobPattern(spec, basePath, UsageExclude, caseSensitive); pattern != nil {
			m.excludePatterns = append(m.excludePatterns, pattern)
		}
	}

	return m
}

// MatchesFile checks if a file path matches the include patterns and doesn't match exclude patterns.
// Returns the index of the matching include pattern, or -1 if no match.
func (m *globMatcher) MatchesFile(path string) int {
	// First check excludes
	for _, exc := range m.excludePatterns {
		if exc.matches(path) {
			return -1
		}
	}

	// If no valid include patterns but includes were specified, nothing matches
	if len(m.includePatterns) == 0 {
		if m.hadIncludes {
			return -1
		}
		return 0
	}

	// Check includes
	for i, inc := range m.includePatterns {
		if inc.matches(path) {
			return i
		}
	}

	return -1
}

// MatchesDirectory checks if a directory could contain matching files.
func (m *globMatcher) MatchesDirectory(path string) bool {
	// First check if excluded
	for _, exc := range m.excludePatterns {
		if exc.matches(path) {
			return false
		}
	}

	// If no valid include patterns but includes were specified, nothing matches
	if len(m.includePatterns) == 0 {
		if m.hadIncludes {
			return false
		}
		return true
	}

	// Check if any include pattern could match files in this directory
	for _, inc := range m.includePatterns {
		if inc.matchesPrefix(path) {
			return true
		}
	}

	return false
}

// visitorNoRegex is similar to visitor but uses globMatcher instead of regex
type visitorNoRegex struct {
	fileMatcher               *globMatcher
	directoryMatcher          *globMatcher
	extensions                []string
	useCaseSensitiveFileNames bool
	host                      vfs.FS
	visited                   collections.Set[string]
	results                   [][]string
	numIncludePatterns        int
}

func (v *visitorNoRegex) visitDirectory(
	path string,
	absolutePath string,
	depth *int,
) {
	// Use the real path for cycle detection
	realPath := v.host.Realpath(absolutePath)
	canonicalPath := tspath.GetCanonicalFileName(realPath, v.useCaseSensitiveFileNames)
	if v.visited.Has(canonicalPath) {
		return
	}
	v.visited.Add(canonicalPath)

	systemEntries := v.host.GetAccessibleEntries(absolutePath)

	// Pre-compute path suffixes to reduce allocations
	// We'll build paths by appending "/" + entry name
	pathPrefix := path
	absPathPrefix := absolutePath
	if len(path) > 0 && path[len(path)-1] != '/' {
		pathPrefix = path + "/"
	}
	if len(absolutePath) > 0 && absolutePath[len(absolutePath)-1] != '/' {
		absPathPrefix = absolutePath + "/"
	}

	for _, current := range systemEntries.Files {
		// Check extension first using just the filename (avoids path concatenation)
		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(current, v.extensions) {
			continue
		}

		// Build absolute name for pattern matching
		absoluteName := absPathPrefix + current

		matchIdx := v.fileMatcher.MatchesFile(absoluteName)
		if matchIdx >= 0 {
			// Only build the relative name if we have a match
			name := pathPrefix + current
			if v.numIncludePatterns == 0 {
				v.results[0] = append(v.results[0], name)
			} else {
				v.results[matchIdx] = append(v.results[matchIdx], name)
			}
		}
	}

	if depth != nil {
		newDepth := *depth - 1
		if newDepth == 0 {
			return
		}
		depth = &newDepth
	}

	for _, current := range systemEntries.Directories {
		name := pathPrefix + current
		absoluteName := absPathPrefix + current

		if v.directoryMatcher.MatchesDirectory(absoluteName) {
			v.visitDirectory(name, absoluteName, depth)
		}
	}
}

// matchFilesNoRegex is the regex-free version of matchFiles
func matchFilesNoRegex(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host vfs.FS) []string {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	// Build file matcher
	fileMatcher := newGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageFiles)

	// Build directory matcher
	directoryMatcher := newGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageDirectories)

	basePaths := getBasePaths(path, includes, useCaseSensitiveFileNames)

	numIncludePatterns := len(fileMatcher.includePatterns)

	var results [][]string
	if numIncludePatterns > 0 {
		results = make([][]string, numIncludePatterns)
		for i := range results {
			results[i] = []string{}
		}
	} else {
		results = [][]string{{}}
	}

	v := visitorNoRegex{
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		host:                      host,
		fileMatcher:               fileMatcher,
		directoryMatcher:          directoryMatcher,
		extensions:                extensions,
		results:                   results,
		numIncludePatterns:        numIncludePatterns,
	}

	for _, basePath := range basePaths {
		v.visitDirectory(basePath, tspath.CombinePaths(currentDirectory, basePath), depth)
	}

	return core.Flatten(results)
}

// globSpecMatcher wraps glob patterns for matching paths.
type globSpecMatcher struct {
	patterns []*globPattern
}

func (m *globSpecMatcher) MatchString(path string) bool {
	if m == nil {
		return false
	}
	for _, p := range m.patterns {
		if p.matches(path) {
			return true
		}
	}
	return false
}

func (m *globSpecMatcher) MatchIndex(path string) int {
	if m == nil {
		return -1
	}
	for i, p := range m.patterns {
		if p.matches(path) {
			return i
		}
	}
	return -1
}

func (m *globSpecMatcher) Len() int {
	if m == nil {
		return 0
	}
	return len(m.patterns)
}

// newGlobSpecMatcher creates a glob-based matcher for multiple specs.
func newGlobSpecMatcher(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *globSpecMatcher {
	if len(specs) == 0 {
		return nil
	}
	m := &globSpecMatcher{}
	for _, spec := range specs {
		if pattern := compileGlobPattern(spec, basePath, usage, useCaseSensitiveFileNames); pattern != nil {
			m.patterns = append(m.patterns, pattern)
		}
	}
	if len(m.patterns) == 0 {
		return nil
	}
	return m
}

// newGlobSingleSpecMatcher creates a glob-based matcher for a single spec.
func newGlobSingleSpecMatcher(spec string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *globSpecMatcher {
	pattern := compileGlobPattern(spec, basePath, usage, useCaseSensitiveFileNames)
	if pattern == nil {
		return nil
	}
	return &globSpecMatcher{patterns: []*globPattern{pattern}}
}
