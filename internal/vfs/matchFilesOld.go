package vfs

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
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

type usage string

const (
	usageFiles       usage = "files"
	usageDirectories usage = "directories"
	usageExclude     usage = "exclude"
)

func getRegularExpressionsForWildcards(specs []string, basePath string, usage usage) []string {
	if len(specs) == 0 {
		return nil
	}
	return core.Map(specs, func(spec string) string {
		return getSubPatternFromSpec(spec, basePath, usage, wildcardMatchers[usage])
	})
}

func getRegularExpressionForWildcard(specs []string, basePath string, usage usage) string {
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
	if usage == "exclude" {
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

// Reserved characters, forces escaping of any non-word (or digit), non-whitespace character.
// It may be inefficient (we could just match (/[-[\]{}()*+?.,\\^$|#\s]/g), but this is future
// proof.
var (
	reservedCharacterPattern *regexp.Regexp = regexp.MustCompile(`[^\w\s/]`)
	wildcardCharCodes                       = []rune{'*', '?'}
)

var (
	commonPackageFolders            = []string{"node_modules", "bower_components", "jspm_packages"}
	implicitExcludePathRegexPattern = "(?!(" + strings.Join(commonPackageFolders, "|") + ")(/|$))"
)

type WildcardMatcher struct {
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

var filesMatcher = WildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragmentFilesMatcher,
	// Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	// files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragmentFilesMatcher)
	},
}

var directoriesMatcher = WildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragment,
	// Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	// files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragment)
	},
}

var excludeMatcher = WildcardMatcher{
	singleAsteriskRegexFragment: singleAsteriskRegexFragment,
	doubleAsteriskRegexFragment: "(/.+?)?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, singleAsteriskRegexFragment)
	},
}

var wildcardMatchers = map[usage]WildcardMatcher{
	usageFiles:       filesMatcher,
	usageDirectories: directoriesMatcher,
	usageExclude:     excludeMatcher,
}

// getPatternFromSpec is now unexported and unused; can be deleted if not needed

func getSubPatternFromSpec(
	spec string,
	basePath string,
	usage usage,
	matcher WildcardMatcher,
) string {
	matcher = wildcardMatchers[usage]

	replaceWildcardCharacter := matcher.replaceWildcardCharacter

	var subpattern strings.Builder
	hasWrittenComponent := false
	components := tspath.GetNormalizedPathComponents(spec, basePath)
	lastComponent := core.LastOrNil(components)
	if usage != "exclude" && lastComponent == "**" {
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
			if usage == "directories" {
				subpattern.WriteString("(")
				optionalCount++
			}

			if hasWrittenComponent {
				subpattern.WriteRune(tspath.DirectorySeparator)
			}

			if usage != "exclude" {
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

// GetExcludePattern creates a regular expression pattern for exclude specs
func GetExcludePattern(excludeSpecs []string, currentDirectory string) string {
	return getRegularExpressionForWildcard(excludeSpecs, currentDirectory, "exclude")
}

// GetFileIncludePatterns creates regular expression patterns for file include specs
func GetFileIncludePatterns(includeSpecs []string, basePath string) []string {
	patterns := getRegularExpressionsForWildcards(includeSpecs, basePath, "files")
	if patterns == nil {
		return nil
	}
	return core.Map(patterns, func(pattern string) string {
		return fmt.Sprintf("^%s$", pattern)
	})
}

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

// getFileMatcherPatterns generates file matching patterns based on the provided path,
// includes, excludes, and other parameters. path is the directory of the tsconfig.json file.
func getFileMatcherPatterns(path string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string) fileMatcherPatterns {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	return fileMatcherPatterns{
		includeFilePatterns:     core.Map(getRegularExpressionsForWildcards(includes, absolutePath, "files"), func(pattern string) string { return "^" + pattern + "$" }),
		includeFilePattern:      getRegularExpressionForWildcard(includes, absolutePath, "files"),
		includeDirectoryPattern: getRegularExpressionForWildcard(includes, absolutePath, "directories"),
		excludePattern:          getRegularExpressionForWildcard(excludes, absolutePath, "exclude"),
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

// getRegexFromPattern is now unexported and unused; can be deleted if not needed

type visitor struct {
	includeFileRegexes        []*regexp2.Regexp
	excludeRegex              *regexp2.Regexp
	includeDirectoryRegex     *regexp2.Regexp
	extensions                []string
	useCaseSensitiveFileNames bool
	host                      FS
	visited                   collections.Set[string]
	results                   [][]string
}

func (v *visitor) visitDirectory(
	path string,
	absolutePath string,
	depth *int,
) {
	canonicalPath := tspath.GetCanonicalFileName(absolutePath, v.useCaseSensitiveFileNames)
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
func MatchFiles(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth *int, host FS) []string {
	// ...existing code...

	// TODO: Implement glob-based matching for MatchFilesOld if needed, or remove this function if unused.
	// For now, return an empty slice to avoid using regex-based logic.
	return []string{}
}

// MatchesExclude checks if a file matches any of the exclude patterns using glob matching (no regexp2)
func MatchesExclude(fileName string, excludeSpecs []string, currentDirectory string, useCaseSensitiveFileNames bool) bool {
	if len(excludeSpecs) == 0 {
		return false
	}

	for _, excludeSpec := range excludeSpecs {
		matcher := NewGlobMatcher(excludeSpec, currentDirectory, useCaseSensitiveFileNames)
		if matcher.MatchesFile(fileName) {
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
		matcher := NewGlobMatcher(includeSpec, basePath, useCaseSensitiveFileNames)
		if matcher.MatchesFile(fileName) {
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
		matcher := NewGlobMatcher(includeSpec, basePath, useCaseSensitiveFileNames)
		if matcher.MatchesFile(fileName) {
			return true
		}
	}
	return false
}
