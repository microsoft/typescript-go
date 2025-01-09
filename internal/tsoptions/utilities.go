package tsoptions

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type FileSystemEntries struct {
	files       []string
	directories []string
}

func getAccessibleFileSystemEntries(path string, host vfs.FS) FileSystemEntries {
	entries := host.GetEntries(path)
	var files []string
	var directories []string
	var entry string
	for _, dirent := range entries {
		entry = dirent.Name()

		// This is necessary because on some file system node fails to exclude
		// "." and "..". See https://github.com/nodejs/node/issues/4002
		if entry == "." || entry == ".." {
			continue
		}

		index := strings.Index(entry[1:], string("."))
		if index == -1 {
			directories = append(directories, entry)
		} else {
			files = append(files, entry)
		}
	}
	return FileSystemEntries{files: files, directories: directories}
}

type FileMatcherPatterns struct {
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
	return fmt.Sprintf("^%s%s", pattern, terminator)
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

// An "includes" path "foo" is implicitly a glob "foo/** /*" (without the space) if its last component has no extension,
// and does not contain any glob characters itself.
func isImplicitGlob(lastPathComponent string) bool {
	re := regexp.MustCompile(`[.*?]`)
	return !re.MatchString(lastPathComponent)
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

var filesMatcher = WildcardMatcher{
	// Matches any single directory segment unless it is the last segment and a .min.js file
	// Breakdown:
	//  [^./]                   # matches everything up to the first . character (excluding directory separators)
	//  (\\.(?!min\\.js$))?     # matches . characters but not if they are part of the .min.js file extension
	singleAsteriskRegexFragment: "([^./]|(\\.(?!min\\.js$))?)*",
	// Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	// files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, "([^./]|(\\.(?!min\\.js$))?)*")
	},
}

var directoriesMatcher = WildcardMatcher{
	singleAsteriskRegexFragment: "[^/]*",
	//Regex for the ** wildcard. Matches any number of subdirectories. When used for including
	//files or directories, does not match subdirectories that start with a . character
	doubleAsteriskRegexFragment: "(/" + implicitExcludePathRegexPattern + "[^/.][^/]*)*?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, "[^/]*")
	},
}

var excludeMatcher = WildcardMatcher{
	singleAsteriskRegexFragment: "[^/]*",
	doubleAsteriskRegexFragment: "(/.+?)?",
	replaceWildcardCharacter: func(match string) string {
		return replaceWildcardCharacter(match, "[^/]*")
	},
}

var wildcardMatchers = map[usage]WildcardMatcher{
	usageFiles:       filesMatcher,
	usageDirectories: directoriesMatcher,
	usageExclude:     excludeMatcher,
}

func getSubPatternFromSpec(
	spec string,
	basePath string,
	usage usage,
	matcher WildcardMatcher,
) string {
	matcher = wildcardMatchers[usage]

	replaceWildcardCharacter := matcher.replaceWildcardCharacter

	subpattern := ""
	hasWrittenComponent := false
	components := tspath.GetNormalizedPathComponents(spec, basePath)
	lastComponent := core.LastOrNil(components)
	if usage != "exclude" && lastComponent == "**" {
		return ""
	}

	// getNormalizedPathComponents includes the separator for the root component.
	// We need to remove to create our regex correctly.
	components[0] = tspath.RemoveTrailingDirectorySeparator(components[0])

	if isImplicitGlob(lastComponent) {
		components = append(components, "**", "*")
	}

	optionalCount := 0
	for _, component := range components {
		if component == "**" {
			subpattern += matcher.doubleAsteriskRegexFragment
		} else {
			if usage == "directories" {
				subpattern += "("
				optionalCount++
			}

			if hasWrittenComponent {
				subpattern += string(tspath.DirectorySeparator)
			}

			if usage != "exclude" {
				componentPattern := ""
				if component != "" && []rune(component)[0] == 0x2A {
					componentPattern += "([^./]" + matcher.singleAsteriskRegexFragment + ")?"
					component = component[1:]
				} else if component != "" && []rune(component)[0] == 0x3F {
					componentPattern += "[^./]"
					component = component[1:]
				}
				componentPattern += reservedCharacterPattern.ReplaceAllStringFunc(component, replaceWildcardCharacter)

				// Patterns should not include subfolders like node_modules unless they are
				// explicitly included as part of the path.
				//
				// As an optimization, if the component pattern is the same as the component,
				// then there definitely were no wildcard characters and we do not need to
				// add the exclusion pattern.
				if componentPattern != component {
					subpattern += implicitExcludePathRegexPattern
				}
				subpattern += componentPattern
			} else {
				subpattern += reservedCharacterPattern.ReplaceAllStringFunc(component, replaceWildcardCharacter)
			}
		}
		hasWrittenComponent = true
	}

	for optionalCount > 0 {
		subpattern += ")?"
		optionalCount--
	}

	return subpattern
}

func getIncludeBasePath(absolute string) string {
	wildcardOffset := core.IndexOfAnyCharCode(absolute, wildcardCharCodes, 0)
	if wildcardOffset < 0 {
		// No "*" or "?" in the path
		if !tspath.HasExtension(absolute) {
			return absolute
		} else {
			tspath.RemoveTrailingDirectorySeparator(tspath.GetDirectoryPath(absolute))
		}
	}
	return absolute[0:strings.LastIndex(absolute, string(tspath.DirectorySeparator))]
}

// getBasePaths computes the unique non-wildcard base paths amongst the provided include patterns.
func getBasePaths(path string, includes []string, useCaseSensitiveFileNames bool) []string {
	// Storage for our results in the form of literal paths (e.g. the paths as written by the user).
	basePaths := []string{path}

	if includes != nil {
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
			includeBasePaths = []string{getIncludeBasePath(absolute)}
		}

		// Sort the offsets array using either the literal or canonical path representations.
		sort.SliceStable(includeBasePaths, func(i, j int) bool {
			return stringutil.GetStringComparer(!useCaseSensitiveFileNames)(includeBasePaths[i], includeBasePaths[j]) < 0
		})

		// Iterate over each include base path and include unique base paths that are not a
		// subpath of an existing base path
		for _, includeBasePath := range includeBasePaths {
			core.Every(basePaths, func(basepath string) bool {
				if !tspath.ContainsPath(basepath, includeBasePath, tspath.ComparePathsOptions{CurrentDirectory: path, UseCaseSensitiveFileNames: !useCaseSensitiveFileNames}) {
					basePaths = append(basePaths, includeBasePath)
					return false
				}
				return true
			})
		}
	}

	return basePaths
}

// getFileMatcherPatterns generates file matching patterns based on the provided path,
// includes, excludes, and other parameters. path is the directory of the tsconfig.json file.
func getFileMatcherPatterns(path string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string) FileMatcherPatterns {
	path = tspath.NormalizePath(path)
	currentDirectory = tspath.NormalizePath(currentDirectory)
	absolutePath := tspath.CombinePaths(currentDirectory, path)

	return FileMatcherPatterns{
		includeFilePatterns:     core.Map(getRegularExpressionsForWildcards(includes, absolutePath, "files"), func(pattern string) string { return fmt.Sprintf("^%s$", pattern) }),
		includeFilePattern:      getRegularExpressionForWildcard(includes, absolutePath, "files"),
		includeDirectoryPattern: getRegularExpressionForWildcard(includes, absolutePath, "directories"),
		excludePattern:          getRegularExpressionForWildcard(excludes, absolutePath, "exclude"),
		basePaths:               getBasePaths(path, includes, useCaseSensitiveFileNames),
	}
}

func getRegexFromPattern(pattern string, useCaseSensitiveFileNames bool) *regexp2.Regexp {
	var re *regexp2.Regexp
	if useCaseSensitiveFileNames {
		re = regexp2.MustCompile(pattern, regexp2.ECMAScript)
	} else {
		re = regexp2.MustCompile(pattern, regexp2.ECMAScript|regexp2.IgnoreCase)
	}
	return re
}

func visitDirectory(visited map[string]bool, path string, absolutePath string, depth int, useCaseSensitiveFileNames bool, host vfs.FS, getFileSystemEntries func(path string, host vfs.FS) FileSystemEntries, includeFileRegexes []*regexp2.Regexp, excludeRegex *regexp2.Regexp, results [][]string, includeDirectoryRegex *regexp2.Regexp, extensions []string) {
	canonicalPath := tspath.GetCanonicalFileName(absolutePath, useCaseSensitiveFileNames)
	if visited[canonicalPath] {
		return
	}
	visited[canonicalPath] = true
	systemEntries := getFileSystemEntries(absolutePath, host)
	files := systemEntries.files
	directories := systemEntries.directories

	for _, current := range files {
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)
		if (extensions != nil && !tspath.FileExtensionIsOneOf(name, extensions)) || (excludeRegex != nil && core.Must(excludeRegex.MatchString(absoluteName))) {
			continue
		}
		if includeFileRegexes == nil {
			results[0] = append(results[0], name)
		} else {
			includeIndex := core.FindIndex(includeFileRegexes, func(re *regexp2.Regexp) bool { return core.Must(re.MatchString(absoluteName)) })
			if includeIndex != -1 {
				results[includeIndex] = append(results[includeIndex], name)
			}
		}
	}

	if depth >= 0 {
		depth--
		if depth == 0 {
			return
		}
	}

	for _, current := range directories {
		name := tspath.CombinePaths(path, current)
		absoluteName := tspath.CombinePaths(absolutePath, current)
		if (includeDirectoryRegex == nil || core.Must(includeDirectoryRegex.MatchString(absoluteName))) && (excludeRegex == nil || !core.Must(excludeRegex.MatchString(absoluteName))) {
			visitDirectory(visited, name, absoluteName, depth, useCaseSensitiveFileNames, host, getFileSystemEntries, includeFileRegexes, excludeRegex, results, includeDirectoryRegex, extensions)
		}
	}
}

// path is the directory of the tsconfig.json
func matchFiles(path string, extensions []string, excludes []string, includes []string, useCaseSensitiveFileNames bool, currentDirectory string, depth int, host vfs.FS, getFileSystemEntries func(path string, host vfs.FS) FileSystemEntries, realpath func(path string) string) []string {
	// path := tspath.NormalizePath(path)
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
	if includeFileRegexes[0] != nil {
		results = core.Map(includeFileRegexes, func(_ *regexp2.Regexp) []string { return []string{} })
	} else {
		results = [][]string{{}}
	}
	visited := make(map[string]bool)
	for _, basePath := range patterns.basePaths {
		visitDirectory(visited, basePath, tspath.CombinePaths(currentDirectory, basePath), depth, useCaseSensitiveFileNames, host, getFileSystemEntries, includeFileRegexes, excludeRegex, results, includeDirectoryRegex, extensions)
	}

	return core.Flatten(results)
}

func readDirectory(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth int) []string {
	return matchFiles(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host, getAccessibleFileSystemEntries, host.Realpath)
}
