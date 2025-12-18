package vfsmatch

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type fileMatcherPatterns struct {
	// One pattern for each "include" spec.
	includeFilePatterns []string
	// One pattern matching one of any of the "include" specs.
	includeFilePattern      string
	includeDirectoryPattern string
	excludePattern          string
	basePaths               []string
}

func getRegularExpressionsForWildcards(specs []string, basePath string, usage Usage) []string {
	if len(specs) == 0 {
		return nil
	}
	return core.Map(specs, func(spec string) string {
		return getSubPatternFromSpec(spec, basePath, usage, wildcardMatchers[usage])
	})
}

func getRegularExpressionForWildcard(specs []string, basePath string, usage Usage) string {
	patterns := getRegularExpressionsForWildcards(specs, basePath, usage)
	if len(patterns) == 0 {
		return ""
	}

	mappedPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		mappedPatterns[i] = fmt.Sprintf("(%s)", pattern)
	}
	pattern := strings.Join(mappedPatterns, "|")

	// If excluding, match "foo/bar/baz...", but if including, only allow "foo".
	var terminator string
	if usage == UsageExclude {
		terminator = "($|/)"
	} else {
		terminator = "$"
	}
	return fmt.Sprintf("^(%s)%s", pattern, terminator)
}

func replaceWildcardCharacter(match string, singleAsteriskRegexFragment string) string {
	if match == "*" {
		return singleAsteriskRegexFragment
	} else {
		if match == "?" {
			return "[^/]"
		} else {
			return "\\" + match
		}
	}
}

// Reserved characters - only escape actual regex metacharacters.
// Go's regexp doesn't support \x escape sequences for arbitrary characters,
// so we only escape characters that have special meaning in regex.
var (
	reservedCharacterPattern *regexp.Regexp = regexp.MustCompile(`[\\.\+*?()\[\]{}^$|#]`)
)

var (
	commonPackageFolders            = []string{"node_modules", "bower_components", "jspm_packages"}
	implicitExcludePathRegexPattern = "(?!(" + strings.Join(commonPackageFolders, "|") + ")(/|$))"
)

type wildcardMatcher struct {
	singleAsteriskRegexFragment string
	doubleAsteriskRegexFragment string
	replaceWildcardCharacter    func(match string) string
}

const (
	// Matches any single directory segment unless it is the last segment and a .min.js file
	// Breakdown:
	//
	//	[^./]                   # matches everything up to the first . character (excluding directory separators)
	//	(\\.(?!min\\.js$))?     # matches . characters but not if they are part of the .min.js file extension
	singleAsteriskRegexFragmentFilesMatcher = "([^./]|(\\.(?!min\\.js$))?)*"
	singleAsteriskRegexFragment             = "[^/]*"
)

var filesMatcher = wildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragmentFilesMatcher,
	// Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	// files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragmentFilesMatcher)
	},
}

var directoriesMatcher = wildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragment,
	// Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	// files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragment)
	},
}

var excludeMatcher = wildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragment,
	doubleAsteriskRegexFragment: "(/.+?)?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragment)
	},
}

var wildcardMatchers = map[Usage]wildcardMatcher{
	UsageFiles:       filesMatcher,
	UsageDirectories: directoriesMatcher,
	UsageExclude:     excludeMatcher,
}

func getPatternFromSpec(
	spec string,
	basePath string,
	usage Usage,
) string {
	pattern := getSubPatternFromSpec(spec, basePath, usage, wildcardMatchers[usage])
	if pattern == "" {
		return ""
	}
	ending := core.IfElse(usage == UsageExclude, "($|/)", "$")
	return fmt.Sprintf("^(%s)%s", pattern, ending)
}

func getSubPatternFromSpec(
	spec string,
	basePath string,
	usage Usage,
	matcher wildcardMatcher,
) string {
	matcher = wildcardMatchers[usage]

	replaceWildcardCharacter := matcher.replaceWildcardCharacter

	var subpattern strings.Builder
	hasWrittenComponent := false
	components := tspath.GetNormalizedPathComponents(spec, basePath)
	lastComponent := core.LastOrNil(components)
	if usage != UsageExclude && lastComponent == "**" {
		return ""
	}

	// getNormalizedPathComponents includes the separator for the root component.
	// We need to remove to create our regex correctly.
	components[0] = tspath.RemoveTrailingDirectorySeparator(components[0])

	if IsImplicitGlob(lastComponent) {
		components = append(components, "**", "*")
	}

	optionalCount := 0
	for _, component := range components {
		if component == "**" {
			subpattern.WriteString(matcher.doubleAsteriskRegexFragment)
		} else {
			if usage == UsageDirectories {
				subpattern.WriteString("(")
				optionalCount++
			}

			if hasWrittenComponent {
				subpattern.WriteRune(tspath.DirectorySeparator)
			}

			if usage != UsageExclude {
				var componentPattern strings.Builder
				if strings.HasPrefix(component, "*") {
					componentPattern.WriteString("([^./]" + matcher.singleAsteriskRegexFragment + ")?")
					component = component[1:]
				} else if strings.HasPrefix(component, "?") {
					componentPattern.WriteString("[^./]")
					component = component[1:]
				}
				componentPattern.WriteString(reservedCharacterPattern.ReplaceAllStringFunc(component, replaceWildcardCharacter))

				// Patterns should not include subfolders like node_modules unless they are
				// explicitly included as part of the path.
				//
				// As an optimization, if the component pattern is the same as the component,
				// then there definitely were no wildcard characters and we do not need to
				// add the exclusion pattern.
				if componentPattern.String() != component {
					subpattern.WriteString(implicitExcludePathRegexPattern)
				}
				subpattern.WriteString(componentPattern.String())
			} else {
				subpattern.WriteString(reservedCharacterPattern.ReplaceAllStringFunc(component, replaceWildcardCharacter))
			}
		}
		hasWrittenComponent = true
	}

	for optionalCount > 0 {
		subpattern.WriteString(")?")
		optionalCount--
	}

	return subpattern.String()
}

// getFileMatcherPatterns generates file matching patterns based on the provided path,
// includes, excludes, and other parameters. path is the directory of the tsconfig.json file.
func getFileMatcherPatterns(path string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string) fileMatcherPatterns {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	return fileMatcherPatterns{
		includeFilePatterns:     core.Map(getRegularExpressionsForWildcards(includes, absolutePath, UsageFiles), func(pattern string) string { return "^" + pattern + "$" }),
		includeFilePattern:      getRegularExpressionForWildcard(includes, absolutePath, UsageFiles),
		includeDirectoryPattern: getRegularExpressionForWildcard(includes, absolutePath, UsageDirectories),
		excludePattern:          getRegularExpressionForWildcard(excludes, absolutePath, UsageExclude),
		basePaths:               getBasePaths(path, includes, useCaseSensitiveFileNames),
	}
}

type regexp2CacheKey struct {
	pattern string
	opts    regexp2.RegexOptions
}

var (
	regexp2CacheMu sync.RWMutex
	regexp2Cache   = make(map[regexp2CacheKey]*regexp2.Regexp)
)

func getRegexFromPattern(pattern string, useCaseSensitiveFileNames bool) *regexp2.Regexp {
	opts := regexp2.RegexOptions(regexp2.ECMAScript)
	if !useCaseSensitiveFileNames {
		opts |= regexp2.IgnoreCase
	}

	key := regexp2CacheKey{pattern, opts}

	regexp2CacheMu.RLock()
	re, ok := regexp2Cache[key]
	regexp2CacheMu.RUnlock()
	if ok {
		return re
	}

	regexp2CacheMu.Lock()
	defer regexp2CacheMu.Unlock()

	re, ok = regexp2Cache[key]
	if ok {
		return re
	}

	// Avoid infinite growth; may cause thrashing but no worse than not caching at all.
	if len(regexp2Cache) > 1000 {
		clear(regexp2Cache)
	}

	// Avoid holding onto the pattern string, since this may pin a full config file in memory.
	pattern = strings.Clone(pattern)
	key.pattern = pattern

	re = regexp2.MustCompile(pattern, opts)
	regexp2Cache[key] = re
	return re
}

type visitor struct {
	includeFileRegexes        []*regexp2.Regexp
	excludeRegex              *regexp2.Regexp
	includeDirectoryRegex     *regexp2.Regexp
	extensions                []string
	useCaseSensitiveFileNames bool
	host                      vfs.FS
	visited                   collections.Set[string]
	results                   [][]string
}

func (v *visitor) visitDirectory(
	path string,
	absolutePath string,
	depth *int,
) {
	// Use the real path (with symlinks resolved) for cycle detection.
	// This prevents infinite loops when symlinks create cycles.
	realPath := v.host.Realpath(absolutePath)
	canonicalPath := tspath.GetCanonicalFileName(realPath, v.useCaseSensitiveFileNames)
	if v.visited.Has(canonicalPath) {
		return
	}
	v.visited.Add(canonicalPath)
	systemEntries := v.host.GetAccessibleEntries(absolutePath)
	files := systemEntries.Files
	directories := systemEntries.Directories

	for _, current := range files {
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)
		if len(v.extensions) > 0 && !tspath.FileExtensionIsOneOf(name, v.extensions) {
			continue
		}
		if v.excludeRegex != nil && core.Must(v.excludeRegex.MatchString(absoluteName)) {
			continue
		}
		if v.includeFileRegexes == nil {
			(v.results)[0] = append((v.results)[0], name)
		} else {
			includeIndex := core.FindIndex(v.includeFileRegexes, func(re *regexp2.Regexp) bool { return core.Must(re.MatchString(absoluteName)) })
			if includeIndex != -1 {
				(v.results)[includeIndex] = append((v.results)[includeIndex], name)
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

	for _, current := range directories {
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)
		if (v.includeDirectoryRegex == nil || core.Must(v.includeDirectoryRegex.MatchString(absoluteName))) && (v.excludeRegex == nil || !core.Must(v.excludeRegex.MatchString(absoluteName))) {
			v.visitDirectory(name, absoluteName, depth)
		}
	}
}

// path is the directory of the tsconfig.json
func matchFiles(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host vfs.FS) []string {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)

	patterns := getFileMatcherPatterns(path, excludes, includes, useCaseSensitiveFileNames, currentDirectory)
	var includeFileRegexes []*regexp2.Regexp
	if patterns.includeFilePatterns != nil {
		includeFileRegexes = core.Map(patterns.includeFilePatterns, func(pattern string) *regexp2.Regexp { return getRegexFromPattern(pattern, useCaseSensitiveFileNames) })
	}
	var includeDirectoryRegex *regexp2.Regexp
	if patterns.includeDirectoryPattern != "" {
		includeDirectoryRegex = getRegexFromPattern(patterns.includeDirectoryPattern, useCaseSensitiveFileNames)
	}
	var excludeRegex *regexp2.Regexp
	if patterns.excludePattern != "" {
		excludeRegex = getRegexFromPattern(patterns.excludePattern, useCaseSensitiveFileNames)
	}

	// Associate an array of results with each include regex. This keeps results in order of the "include" order.
	// If there are no "includes", then just put everything in results[0].
	var results [][]string
	if len(includeFileRegexes) > 0 {
		tempResults := make([][]string, len(includeFileRegexes))
		for i := range includeFileRegexes {
			tempResults[i] = []string{}
		}
		results = tempResults
	} else {
		results = [][]string{{}}
	}
	v := visitor{
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		host:                      host,
		includeFileRegexes:        includeFileRegexes,
		excludeRegex:              excludeRegex,
		includeDirectoryRegex:     includeDirectoryRegex,
		extensions:                extensions,
		results:                   results,
	}
	for _, basePath := range patterns.basePaths {
		v.visitDirectory(basePath, tspath.CombinePaths(currentDirectory, basePath), depth)
	}

	return core.Flatten(results)
}

// regexSpecMatcher wraps a regexp2.Regexp for SpecMatcher interface.
type regexSpecMatcher struct {
	re *regexp2.Regexp
}

func (m *regexSpecMatcher) MatchString(path string) bool {
	if m == nil || m.re == nil {
		return false
	}
	matched, err := m.re.MatchString(path)
	return err == nil && matched
}

// newRegexSpecMatcher creates a regex-based matcher for multiple specs.
func newRegexSpecMatcher(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *regexSpecMatcher {
	pattern := getRegularExpressionForWildcard(specs, basePath, usage)
	if pattern == "" {
		return nil
	}
	return &regexSpecMatcher{re: getRegexFromPattern(pattern, useCaseSensitiveFileNames)}
}

// newRegexSingleSpecMatcher creates a regex-based matcher for a single spec.
func newRegexSingleSpecMatcher(spec string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *regexSpecMatcher {
	pattern := getPatternFromSpec(spec, basePath, usage)
	if pattern == "" {
		return nil
	}
	return &regexSpecMatcher{re: getRegexFromPattern(pattern, useCaseSensitiveFileNames)}
}

// regexSpecMatchers holds a list of individual regex matchers for index lookup.
type regexSpecMatchers struct {
	matchers []*regexp2.Regexp
}

func (m *regexSpecMatchers) MatchIndex(path string) int {
	for i, re := range m.matchers {
		if matched, err := re.MatchString(path); err == nil && matched {
			return i
		}
	}
	return -1
}

// newRegexSpecMatchers creates individual regex matchers for each spec.
func newRegexSpecMatchers(specs []string, basePath string, usage Usage, useCaseSensitiveFileNames bool) *regexSpecMatchers {
	patterns := getRegularExpressionsForWildcards(specs, basePath, usage)
	if len(patterns) == 0 {
		return nil
	}
	matchers := make([]*regexp2.Regexp, len(patterns))
	for i, pattern := range patterns {
		// Wrap pattern with ^ and $ for full match
		fullPattern := "^" + pattern + "$"
		matchers[i] = getRegexFromPattern(fullPattern, useCaseSensitiveFileNames)
	}
	return &regexSpecMatchers{matchers: matchers}
}
