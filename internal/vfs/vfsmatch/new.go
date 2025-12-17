package vfsmatch

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// GlobPattern represents a compiled glob pattern for matching file paths.
// It stores the pattern components for efficient matching without using regex.
type GlobPattern struct {
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

// CompileGlobPattern compiles a glob spec into a GlobPattern for matching.
func CompileGlobPattern(spec string, basePath string, usage Usage, caseSensitive bool) *GlobPattern {
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

	pattern := &GlobPattern{
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

// Matches checks if the given path matches this glob pattern.
func (p *GlobPattern) Matches(path string) bool {
	if p == nil {
		return false
	}

	// Split the path into components
	pathComponents := splitPath(path)

	matched := p.matchComponents(pathComponents, 0, 0, false)

	return matched
}

// MatchesPrefix checks if the given directory path could potentially match files under it.
// This is used for directory filtering during traversal.
func (p *GlobPattern) MatchesPrefix(path string) bool {
	if p == nil {
		return false
	}

	pathComponents := splitPath(path)

	return p.matchComponentsPrefix(pathComponents, 0, 0)
}

// splitPath splits a path into its components
func splitPath(path string) []string {
	// Handle the case of an absolute path
	if len(path) > 0 && path[0] == '/' {
		rest := strings.Split(strings.TrimPrefix(path, "/"), "/")
		// Prepend empty string to represent root
		result := make([]string, 0, len(rest)+1)
		result = append(result, "")
		for _, s := range rest {
			if s != "" {
				result = append(result, s)
			}
		}
		return result
	}

	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))
	for _, s := range parts {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// matchComponents recursively matches path components against pattern components
func (p *GlobPattern) matchComponents(pathComps []string, pathIdx int, patternIdx int, inDoubleAsterisk bool) bool {
	// If we've consumed all pattern components, check if path is also fully consumed
	if patternIdx >= len(p.components) {
		if p.isExclude {
			// For exclude patterns, we can match a prefix
			return true
		}
		return pathIdx >= len(pathComps)
	}

	// If we've consumed all path components but still have pattern components
	if pathIdx >= len(pathComps) {
		// Check if remaining pattern components are all optional (** only)
		for i := patternIdx; i < len(p.components); i++ {
			if !p.components[i].isDoubleAsterisk {
				return false
			}
		}
		return true
	}

	pc := p.components[patternIdx]
	pathComp := pathComps[pathIdx]

	if pc.isDoubleAsterisk {
		// ** can match zero or more directory levels
		// First, try matching zero directories (skip the **)
		if p.matchComponents(pathComps, pathIdx, patternIdx+1, true) {
			return true
		}

		// For include patterns, ** should not match directories starting with . or common package folders
		// But we still try to skip those directories and continue matching
		if !p.isExclude {
			if len(pathComp) > 0 && pathComp[0] == '.' {
				// Don't match hidden directories in ** for includes - return false
				// The next pattern component (if any) might explicitly match it
				return false
			}
			if isCommonPackageFolder(pathComp) {
				// Don't match common package folders in ** for includes
				return false
			}
		}

		// Match current component with ** and continue
		return p.matchComponents(pathComps, pathIdx+1, patternIdx, true)
	}

	// Check implicit package folder exclusion
	if pc.implicitlyExcludePackages && !p.isExclude && isCommonPackageFolder(pathComp) {
		return false
	}

	// Match current component
	if !p.matchComponent(pc, pathComp, inDoubleAsterisk) {
		return false
	}

	// Continue to next components
	return p.matchComponents(pathComps, pathIdx+1, patternIdx+1, false)
}

// matchComponentsPrefix checks if the path could be a prefix of a matching path
func (p *GlobPattern) matchComponentsPrefix(pathComps []string, pathIdx int, patternIdx int) bool {
	// If we've consumed all path components, this prefix could match
	if pathIdx >= len(pathComps) {
		return true
	}

	// If we've consumed all pattern components, no more matches possible
	if patternIdx >= len(p.components) {
		return false
	}

	pc := p.components[patternIdx]
	pathComp := pathComps[pathIdx]

	if pc.isDoubleAsterisk {
		// ** can match any directory level
		// Try matching zero (skip **) or more directories
		if p.matchComponentsPrefix(pathComps, pathIdx, patternIdx+1) {
			return true
		}

		// For include patterns, ** should not match hidden or package directories
		if !p.isExclude {
			if len(pathComp) > 0 && pathComp[0] == '.' {
				return false
			}
			if isCommonPackageFolder(pathComp) {
				return false
			}
		}

		return p.matchComponentsPrefix(pathComps, pathIdx+1, patternIdx)
	}

	// Check implicit package folder exclusion
	if pc.implicitlyExcludePackages && !p.isExclude && isCommonPackageFolder(pathComp) {
		return false
	}

	// Match current component
	if !p.matchComponent(pc, pathComp, false) {
		return false
	}

	return p.matchComponentsPrefix(pathComps, pathIdx+1, patternIdx+1)
}

// matchComponent matches a single path component against a pattern component
func (p *GlobPattern) matchComponent(pc patternComponent, pathComp string, afterDoubleAsterisk bool) bool {
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
func (p *GlobPattern) matchWildcardComponent(segments []patternSegment, s string) bool {
	// For non-exclude patterns, if the segments start with * or ?,
	// the matched string cannot start with '.'
	if !p.isExclude && len(segments) > 0 && len(s) > 0 && s[0] == '.' {
		firstSeg := segments[0]
		if firstSeg.kind == segmentStar || firstSeg.kind == segmentQuestion {
			// Pattern starts with wildcard, so it cannot match a string starting with '.'
			return false
		}
	}

	return p.matchSegments(segments, 0, s, 0)
}

func (p *GlobPattern) matchSegments(segments []patternSegment, segIdx int, s string, sIdx int) bool {
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
		// For files usage, also need to handle .min.js exclusion

		// Try matching zero characters first
		if p.matchSegments(segments, segIdx+1, s, sIdx) {
			// Before returning true, check min.js exclusion
			if p.excludeMinJs && segIdx == 0 && segIdx+1 < len(segments) {
				// Check if this could result in matching a .min.js file
				if p.wouldMatchMinJs(s) {
					return false
				}
			}
			return true
		}

		// Try matching more characters
		for i := sIdx; i < len(s); i++ {
			if s[i] == '/' {
				break
			}
			if p.matchSegments(segments, segIdx+1, s, i+1) {
				// Check min.js exclusion
				if p.excludeMinJs && strings.HasSuffix(s, ".min.js") {
					// Only exclude if pattern doesn't explicitly include .min.js
					if !p.patternExplicitlyIncludesMinJs(segments) {
						return false
					}
				}
				return true
			}
		}
		return false
	}

	return false
}

// wouldMatchMinJs checks if the filename ends with .min.js
func (p *GlobPattern) wouldMatchMinJs(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".min.js")
}

// patternExplicitlyIncludesMinJs checks if the pattern explicitly includes .min.js
func (p *GlobPattern) patternExplicitlyIncludesMinJs(segments []patternSegment) bool {
	// Look for .min.js in the literal segments
	for _, seg := range segments {
		if seg.kind == segmentLiteral && strings.Contains(strings.ToLower(seg.literal), ".min.js") {
			return true
		}
	}
	return false
}

// stringsEqual compares two strings with case sensitivity based on pattern settings
func (p *GlobPattern) stringsEqual(a, b string) bool {
	if p.caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// isCommonPackageFolder checks if a directory name is a common package folder
func isCommonPackageFolder(name string) bool {
	lower := strings.ToLower(name)
	return lower == "node_modules" || lower == "bower_components" || lower == "jspm_packages"
}

// GlobMatcher holds compiled glob patterns for matching files.
type GlobMatcher struct {
	includePatterns []*GlobPattern
	excludePatterns []*GlobPattern
	caseSensitive   bool
	// hadIncludes tracks whether any include specs were provided (even if they compiled to nothing)
	hadIncludes bool
}

// NewGlobMatcher creates a new GlobMatcher from include and exclude specs.
func NewGlobMatcher(includes []string, excludes []string, basePath string, caseSensitive bool, usage Usage) *GlobMatcher {
	m := &GlobMatcher{
		caseSensitive: caseSensitive,
		hadIncludes:   len(includes) > 0,
	}

	for _, spec := range includes {
		if pattern := CompileGlobPattern(spec, basePath, usage, caseSensitive); pattern != nil {
			m.includePatterns = append(m.includePatterns, pattern)
		}
	}

	for _, spec := range excludes {
		if pattern := CompileGlobPattern(spec, basePath, UsageExclude, caseSensitive); pattern != nil {
			m.excludePatterns = append(m.excludePatterns, pattern)
		}
	}

	return m
}

// MatchesFile checks if a file path matches the include patterns and doesn't match exclude patterns.
// Returns the index of the matching include pattern, or -1 if no match.
func (m *GlobMatcher) MatchesFile(path string) int {
	// First check excludes
	for _, exc := range m.excludePatterns {
		if exc.Matches(path) {
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
		if inc.Matches(path) {
			return i
		}
	}

	return -1
}

// MatchesDirectory checks if a directory could contain matching files.
func (m *GlobMatcher) MatchesDirectory(path string) bool {
	// First check if excluded
	for _, exc := range m.excludePatterns {
		if exc.Matches(path) {
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
		if inc.MatchesPrefix(path) {
			return true
		}
	}

	return false
}

// visitorNoRegex is similar to visitor but uses GlobMatcher instead of regex
type visitorNoRegex struct {
	fileMatcher               *GlobMatcher
	directoryMatcher          *GlobMatcher
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

	for _, current := range systemEntries.Files {
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)

		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(name, v.extensions) {
			continue
		}

		matchIdx := v.fileMatcher.MatchesFile(absoluteName)
		if matchIdx >= 0 {
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
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)

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
	fileMatcher := NewGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageFiles)

	// Build directory matcher
	directoryMatcher := NewGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageDirectories)

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
