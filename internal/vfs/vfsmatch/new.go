package vfsmatch

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// globPattern is a compiled glob pattern for matching file paths without regex.
type globPattern struct {
	components    []component // path segments to match (e.g., ["src", "**", "*.ts"])
	isExclude     bool        // exclude patterns have different matching rules
	caseSensitive bool
	excludeMinJs  bool // for "files" patterns, exclude .min.js by default
}

// component is a single path segment in a glob pattern.
// Examples: "src" (literal), "*" (wildcard), "*.ts" (wildcard), "**" (recursive)
type component struct {
	kind     componentKind
	literal  string    // for kindLiteral: the exact string to match
	segments []segment // for kindWildcard: parsed wildcard pattern
	// Include patterns with wildcards skip common package folders (node_modules, etc.)
	skipPackageFolders bool
}

type componentKind int

const (
	kindLiteral        componentKind = iota // exact match (e.g., "src")
	kindWildcard                            // contains * or ? (e.g., "*.ts")
	kindDoubleAsterisk                      // ** matches zero or more directories
)

// segment is a piece of a wildcard component.
// Example: "*.ts" becomes [segStar, segLiteral(".ts")]
type segment struct {
	kind    segmentKind
	literal string // only for segLiteral
}

type segmentKind int

const (
	segLiteral  segmentKind = iota // exact text
	segStar                        // * matches any chars except /
	segQuestion                    // ? matches single char except /
)

// compileGlobPattern compiles a glob spec (e.g., "src/**/*.ts") into a pattern.
// Returns nil if the pattern would match nothing.
func compileGlobPattern(spec string, basePath string, usage Usage, caseSensitive bool) *globPattern {
	parts := tspath.GetNormalizedPathComponents(spec, basePath)

	// "src/**" without a filename matches nothing (for include patterns)
	if usage != UsageExclude && core.LastOrNil(parts) == "**" {
		return nil
	}

	// Normalize root: "/home/" -> "/home"
	parts[0] = tspath.RemoveTrailingDirectorySeparator(parts[0])

	// Directories implicitly match all files: "src" -> "src/**/*"
	if IsImplicitGlob(core.LastOrNil(parts)) {
		parts = append(parts, "**", "*")
	}

	p := &globPattern{
		isExclude:     usage == UsageExclude,
		caseSensitive: caseSensitive,
		excludeMinJs:  usage == UsageFiles,
	}

	for _, part := range parts {
		p.components = append(p.components, parseComponent(part, usage != UsageExclude))
	}
	return p
}

// parseComponent converts a path segment string into a component.
func parseComponent(s string, isInclude bool) component {
	if s == "**" {
		return component{kind: kindDoubleAsterisk}
	}
	if !strings.ContainsAny(s, "*?") {
		return component{kind: kindLiteral, literal: s}
	}
	return component{
		kind:               kindWildcard,
		segments:           parseSegments(s),
		skipPackageFolders: isInclude,
	}
}

// parseSegments breaks "*.ts" into [segStar, segLiteral(".ts")]
func parseSegments(s string) []segment {
	var result []segment
	start := 0
	for i := range len(s) {
		switch s[i] {
		case '*', '?':
			if i > start {
				result = append(result, segment{kind: segLiteral, literal: s[start:i]})
			}
			if s[i] == '*' {
				result = append(result, segment{kind: segStar})
			} else {
				result = append(result, segment{kind: segQuestion})
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, segment{kind: segLiteral, literal: s[start:]})
	}
	return result
}

// matches returns true if path matches this pattern.
func (p *globPattern) matches(path string) bool {
	return p.matchPath(path, 0, 0, false)
}

// matchesPrefix returns true if files under this directory path could match.
// Used to skip directories during traversal.
func (p *globPattern) matchesPrefix(path string) bool {
	return p.matchPath(path, 0, 0, true)
}

// matchPath checks if path matches the pattern starting from the given offsets.
// If prefixOnly is true, returns true when path is exhausted (prefix matching for directories).
func (p *globPattern) matchPath(path string, pathOffset, compIdx int, prefixOnly bool) bool {
	for {
		pathPart, nextOffset, ok := nextPathPart(path, pathOffset)
		if !ok {
			if prefixOnly {
				return true // Path exhausted - could potentially match
			}
			return p.patternSatisfied(compIdx)
		}

		if compIdx >= len(p.components) {
			// Exclude patterns match prefixes (e.g., "node_modules" excludes "node_modules/foo")
			return p.isExclude && !prefixOnly
		}

		comp := p.components[compIdx]

		switch comp.kind {
		case kindDoubleAsterisk:
			// ** can match zero directories: try skipping it
			if p.matchPath(path, pathOffset, compIdx+1, prefixOnly) {
				return true
			}
			// ** should not match hidden dirs or package folders (for includes)
			if !p.isExclude && (isHiddenPath(pathPart) || isPackageFolder(pathPart)) {
				return false
			}
			// ** matches this directory, try next path part with same **
			pathOffset = nextOffset
			continue

		case kindLiteral:
			if comp.skipPackageFolders && isPackageFolder(pathPart) {
				panic("unreachable: literal components never have skipPackageFolders")
			}
			if !p.stringsEqual(comp.literal, pathPart) {
				return false
			}

		case kindWildcard:
			if comp.skipPackageFolders && isPackageFolder(pathPart) {
				return false
			}
			if !p.matchWildcard(comp.segments, pathPart) {
				return false
			}
		}

		pathOffset = nextOffset
		compIdx++
	}
}

// patternSatisfied checks if remaining pattern components can match empty input.
func (p *globPattern) patternSatisfied(compIdx int) bool {
	// A pattern is satisfied when remaining components can match empty input.
	// For both include and exclude patterns, only trailing "**" components may match nothing.
	for _, c := range p.components[compIdx:] {
		if c.kind != kindDoubleAsterisk {
			return false
		}
	}
	return true
}

// nextPathPart extracts the next path component from path starting at offset.
func nextPathPart(path string, offset int) (part string, nextOffset int, ok bool) {
	if offset >= len(path) {
		return "", offset, false
	}

	// Handle leading slash (root of absolute path)
	if offset == 0 && path[0] == '/' {
		return "", 1, true
	}

	// Skip consecutive slashes
	for offset < len(path) && path[offset] == '/' {
		offset++
	}
	if offset >= len(path) {
		return "", offset, false
	}

	// Find end of this component
	rest := path[offset:]
	if idx := strings.IndexByte(rest, '/'); idx >= 0 {
		return rest[:idx], offset + idx, true
	}
	return rest, len(path), true
}

// matchWildcard matches a path component against wildcard segments.
func (p *globPattern) matchWildcard(segs []segment, s string) bool {
	// Include patterns: wildcards at start cannot match hidden files
	if !p.isExclude && len(segs) > 0 && isHiddenPath(s) && (segs[0].kind == segStar || segs[0].kind == segQuestion) {
		return false
	}

	// Fast path: single * followed by literal suffix (e.g., "*.ts")
	if len(segs) == 2 && segs[0].kind == segStar && segs[1].kind == segLiteral {
		suffix := segs[1].literal
		if len(s) < len(suffix) || !p.stringsEqual(suffix, s[len(s)-len(suffix):]) {
			return false
		}
		return p.checkMinJsExclusion(s, segs)
	}

	return p.matchSegments(segs, 0, s, 0) && p.checkMinJsExclusion(s, segs)
}

// matchSegments recursively matches segments against string s.
func (p *globPattern) matchSegments(segs []segment, segIdx int, s string, sIdx int) bool {
	if segIdx >= len(segs) {
		return sIdx >= len(s)
	}

	seg := segs[segIdx]

	switch seg.kind {
	case segLiteral:
		end := sIdx + len(seg.literal)
		if end > len(s) {
			return false
		}
		if !p.stringsEqual(seg.literal, s[sIdx:end]) {
			return false
		}
		return p.matchSegments(segs, segIdx+1, s, end)

	case segQuestion:
		if sIdx >= len(s) || s[sIdx] == '/' {
			return false
		}
		return p.matchSegments(segs, segIdx+1, s, sIdx+1)

	case segStar:
		// Try matching 0, 1, 2, ... characters (but not /)
		if p.matchSegments(segs, segIdx+1, s, sIdx) {
			return true
		}
		for i := sIdx; i < len(s) && s[i] != '/'; i++ {
			if p.matchSegments(segs, segIdx+1, s, i+1) {
				return true
			}
		}
		return false
	default:
		panic("unreachable: unknown segment kind")
	}
}

// checkMinJsExclusion returns false if this is a .min.js file that should be excluded.
func (p *globPattern) checkMinJsExclusion(filename string, segs []segment) bool {
	if !p.excludeMinJs || !strings.HasSuffix(strings.ToLower(filename), ".min.js") {
		return true
	}
	// Allow if pattern explicitly includes .min.js
	for _, seg := range segs {
		if seg.kind == segLiteral && strings.Contains(strings.ToLower(seg.literal), ".min.js") {
			return true
		}
	}
	return false
}

// stringsEqual compares strings with appropriate case sensitivity.
func (p *globPattern) stringsEqual(a, b string) bool {
	if p.caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// isHiddenPath checks if a path component is hidden (starts with dot).
func isHiddenPath(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// isPackageFolder checks if name is a common package folder (node_modules, etc.)
func isPackageFolder(name string) bool {
	switch len(name) {
	case len("node_modules"):
		return strings.EqualFold(name, "node_modules")
	case len("jspm_packages"):
		return strings.EqualFold(name, "jspm_packages")
	case len("bower_components"):
		return strings.EqualFold(name, "bower_components")
	}
	return false
}

func ensureTrailingSlash(s string) string {
	if len(s) > 0 && s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

// globMatcher combines include and exclude patterns for file matching.
type globMatcher struct {
	includes    []*globPattern
	excludes    []*globPattern
	hadIncludes bool // true if include specs were provided (even if none compiled)
}

func newGlobMatcher(includeSpecs, excludeSpecs []string, basePath string, caseSensitive bool, usage Usage) *globMatcher {
	m := &globMatcher{hadIncludes: len(includeSpecs) > 0}

	for _, spec := range includeSpecs {
		if p := compileGlobPattern(spec, basePath, usage, caseSensitive); p != nil {
			m.includes = append(m.includes, p)
		}
	}
	for _, spec := range excludeSpecs {
		if p := compileGlobPattern(spec, basePath, UsageExclude, caseSensitive); p != nil {
			m.excludes = append(m.excludes, p)
		}
	}
	return m
}

// MatchesFile returns the index of the matching include pattern, or -1 if excluded/no match.
func (m *globMatcher) MatchesFile(path string) int {
	for _, exc := range m.excludes {
		if exc.matches(path) {
			return -1
		}
	}
	if len(m.includes) == 0 {
		if m.hadIncludes {
			return -1
		}
		return 0
	}
	for i, inc := range m.includes {
		if inc.matches(path) {
			return i
		}
	}
	return -1
}

// MatchesDirectory returns true if this directory could contain matching files.
func (m *globMatcher) MatchesDirectory(path string) bool {
	for _, exc := range m.excludes {
		if exc.matches(path) {
			return false
		}
	}

	if len(m.includes) == 0 {
		return !m.hadIncludes
	}

	for _, inc := range m.includes {
		if inc.matchesPrefix(path) {
			return true
		}
	}
	return false
}

// globVisitor traverses directories matching files against glob patterns.
type globVisitor struct {
	host                      vfs.FS
	fileMatcher               *globMatcher
	directoryMatcher          *globMatcher
	extensions                []string
	useCaseSensitiveFileNames bool
	visited                   collections.Set[string]
	results                   [][]string
}

func (v *globVisitor) visit(path, absolutePath string, depth *int) {
	// Detect symlink cycles
	realPath := v.host.Realpath(absolutePath)
	canonicalPath := tspath.GetCanonicalFileName(realPath, v.useCaseSensitiveFileNames)
	if v.visited.Has(canonicalPath) {
		return
	}
	v.visited.Add(canonicalPath)

	entries := v.host.GetAccessibleEntries(absolutePath)

	pathPrefix := ensureTrailingSlash(path)
	absPrefix := ensureTrailingSlash(absolutePath)

	for _, file := range entries.Files {
		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(file, v.extensions) {
			continue
		}
		if idx := v.fileMatcher.MatchesFile(absPrefix + file); idx >= 0 {
			v.results[idx] = append(v.results[idx], pathPrefix+file)
		}
	}

	if depth != nil {
		newDepth := *depth - 1
		if newDepth == 0 {
			return
		}
		depth = &newDepth
	}

	for _, dir := range entries.Directories {
		absDir := absPrefix + dir
		if v.directoryMatcher.MatchesDirectory(absDir) {
			v.visit(pathPrefix+dir, absDir, depth)
		}
	}
}

// matchFilesNoRegex matches files using compiled glob patterns (no regex).
func matchFilesNoRegex(path string, extensions, excludes, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host vfs.FS) []string {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	fileMatcher := newGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageFiles)
	directoryMatcher := newGlobMatcher(includes, excludes, absolutePath, useCaseSensitiveFileNames, UsageDirectories)

	v := globVisitor{
		host:                      host,
		fileMatcher:               fileMatcher,
		directoryMatcher:          directoryMatcher,
		extensions:                extensions,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		results:                   make([][]string, max(len(fileMatcher.includes), 1)),
	}

	for _, basePath := range getBasePaths(path, includes, useCaseSensitiveFileNames) {
		v.visit(basePath, tspath.CombinePaths(currentDirectory, basePath), depth)
	}

	return core.Flatten(v.results)
}

// globSpecMatcher wraps multiple glob patterns for matching paths.
type globSpecMatcher struct {
	patterns []*globPattern
}

// MatchString returns true if any pattern matches the path.
func (m *globSpecMatcher) MatchString(path string) bool {
	for _, p := range m.patterns {
		if p.matches(path) {
			return true
		}
	}
	return false
}

// MatchIndex returns the index of the first matching pattern, or -1.
func (m *globSpecMatcher) MatchIndex(path string) int {
	for i, p := range m.patterns {
		if p.matches(path) {
			return i
		}
	}
	return -1
}

// newGlobSpecMatcher creates a matcher for multiple glob specs.
func newGlobSpecMatcher(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *globSpecMatcher {
	if len(specs) == 0 {
		return nil
	}
	var patterns []*globPattern
	for _, spec := range specs {
		if p := compileGlobPattern(spec, basePath, usage, useCaseSensitiveFileNames); p != nil {
			patterns = append(patterns, p)
		}
	}
	if len(patterns) == 0 {
		return nil
	}
	return &globSpecMatcher{patterns: patterns}
}

// newGlobSingleSpecMatcher creates a matcher for a single glob spec.
func newGlobSingleSpecMatcher(spec, basePath string, usage Usage, useCaseSensitiveFileNames bool) *globSpecMatcher {
	return newGlobSpecMatcher([]string{spec}, basePath, usage, useCaseSensitiveFileNames)
}
