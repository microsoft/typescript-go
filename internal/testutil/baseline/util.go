package baseline

import (
	"regexp"
	"strings"

	"github.com/microsoft/typescript-go/internal/compiler"
)

var testPathPrefix = regexp.MustCompile(`(?:(file:\/{3})|\/)\.(?:ts|lib|src)\/`)
var testPathCharacters = regexp.MustCompile(`[\^<>:"|?*%]`)
var testPathDotDot = regexp.MustCompile(`\.\.\/`)

// This is done so tests work on windows _and_ linux
var canonicalizeForHarness = strings.ToLower

var libFolder = "built/local/"
var builtFolder = "/.ts"

func removeTestPathPrefixes(text string, retainTrailingDirectorySeparator bool) string {
	testPathPrefix.ReplaceAllStringFunc(text, func(scheme string) string {
		if scheme != "" {
			return scheme
		}
		if retainTrailingDirectorySeparator {
			return "/"
		}
		return ""
	})
	return testPathPrefix.ReplaceAllString(text, "/")
}

func isDefaultLibraryFile(filePath string) bool {
	fileName := compiler.GetBaseFileName(filePath, nil, false)
	return strings.HasPrefix(fileName, "lib.") && strings.HasSuffix(fileName, compiler.ExtensionDts)
}

func isBuiltFile(filePath string) bool {
	return strings.HasPrefix(filePath, libFolder) || strings.HasPrefix(filePath, compiler.EnsureTrailingDirectorySeparator(builtFolder))
}

func isTsConfigFile(path string) bool {
	return strings.Contains(path, "tsconfig") && strings.Contains(path, "json")
}

func sanitizeTestFilePath(name string) string {
	path := testPathCharacters.ReplaceAllString(name, "_")
	path = compiler.NormalizeSlashes(path)
	path = testPathDotDot.ReplaceAllString(path, "__dotdot/")
	path = string(compiler.ToPath(path, "", canonicalizeForHarness))
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}
