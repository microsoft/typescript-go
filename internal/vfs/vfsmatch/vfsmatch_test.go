package vfsmatch

import (
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestMatchFiles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		files                     map[string]string
		path                      string
		extensions                []string
		excludes                  []string
		includes                  []string
		useCaseSensitiveFileNames bool
		currentDirectory          string
		depth                     *int
		expected                  []string
	}{
		{
			name: "simple include all",
			files: map[string]string{
				"/project/src/index.ts":              "export {}",
				"/project/src/util.ts":               "export {}",
				"/project/src/sub/file.ts":           "export {}",
				"/project/tests/test.ts":             "export {}",
				"/project/node_modules/pkg/index.js": "module.exports = {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts", "/project/src/sub/file.ts", "/project/tests/test.ts"},
		},
		{
			name: "exclude node_modules",
			files: map[string]string{
				"/project/src/index.ts":              "export {}",
				"/project/src/util.ts":               "export {}",
				"/project/node_modules/pkg/index.ts": "export {}",
				"/project/node_modules/pkg/lib.d.ts": "declare module 'pkg'",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx", ".d.ts"},
			excludes:                  []string{"node_modules/**/*"},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts"},
		},
		{
			name: "specific include directory",
			files: map[string]string{
				"/project/src/index.ts":    "export {}",
				"/project/src/util.ts":     "export {}",
				"/project/tests/test.ts":   "export {}",
				"/project/docs/readme.md":  "# readme",
				"/project/build/output.js": "console.log('hello')",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"src/**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts"},
		},
		{
			name: "multiple include patterns",
			files: map[string]string{
				"/project/src/index.ts":     "export {}",
				"/project/src/util.ts":      "export {}",
				"/project/tests/test.ts":    "export {}",
				"/project/scripts/build.ts": "export {}",
				"/project/docs/readme.md":   "# readme",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"src/**/*", "tests/**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts", "/project/tests/test.ts"},
		},
		{
			name: "case insensitive matching",
			files: map[string]string{
				"/project/SRC/Index.TS":   "export {}",
				"/project/src/UTIL.ts":    "export {}",
				"/project/Docs/readme.md": "# readme",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"src/**/*"},
			useCaseSensitiveFileNames: false,
			currentDirectory:          "/",
			expected:                  []string{"/project/SRC/UTIL.ts"},
		},
		{
			name: "exclude with wildcards",
			files: map[string]string{
				"/project/src/index.ts":         "export {}",
				"/project/src/util.ts":          "export {}",
				"/project/src/types.d.ts":       "export {}",
				"/project/src/generated/api.ts": "export {}",
				"/project/src/generated/db.ts":  "export {}",
				"/project/tests/test.ts":        "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx", ".d.ts"},
			excludes:                  []string{"src/generated/**/*"},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/types.d.ts", "/project/src/util.ts", "/project/tests/test.ts"},
		},
		{
			name: "depth limit",
			files: map[string]string{
				"/project/index.ts":                "export {}",
				"/project/src/util.ts":             "export {}",
				"/project/src/deep/nested/file.ts": "export {}",
				"/project/src/other/file.ts":       "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			depth:                     func() *int { d := 2; return &d }(),
			expected:                  []string{"/project/index.ts", "/project/src/util.ts"},
		},
		{
			name: "relative excludes",
			files: map[string]string{
				"/project/src/index.ts":              "export {}",
				"/project/src/util.ts":               "export {}",
				"/project/build/output.js":           "console.log('hello')",
				"/project/node_modules/pkg/index.js": "module.exports = {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx", ".js"},
			excludes:                  []string{"./node_modules", "./build"},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts"},
		},
		{
			name: "empty includes and excludes",
			files: map[string]string{
				"/project/index.ts":      "export {}",
				"/project/src/util.ts":   "export {}",
				"/project/tests/test.ts": "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/index.ts", "/project/src/util.ts", "/project/tests/test.ts"},
		},
		{
			name: "star pattern matching",
			files: map[string]string{
				"/project/test.ts":      "export {}",
				"/project/test.spec.ts": "export {}",
				"/project/util.ts":      "export {}",
				"/project/other.js":     "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{"*.spec.ts"},
			includes:                  []string{"*.ts"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/test.ts", "/project/util.ts"},
		},
		{
			name: "mixed file extensions",
			files: map[string]string{
				"/project/component.tsx": "export {}",
				"/project/util.ts":       "export {}",
				"/project/types.d.ts":    "export {}",
				"/project/config.js":     "module.exports = {}",
				"/project/styles.css":    "body {}",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx", ".d.ts"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/component.tsx", "/project/types.d.ts", "/project/util.ts"},
		},
		{
			name:                      "empty filesystem",
			files:                     map[string]string{},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  nil, // ReadDirectory returns nil for empty results
		},
		{
			name: "no matching extensions",
			files: map[string]string{
				"/project/file.js":  "export {}",
				"/project/file.py":  "print('hello')",
				"/project/file.txt": "hello world",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  nil, // ReadDirectory returns nil for empty results
		},
		{
			name: "exclude everything",
			files: map[string]string{
				"/project/index.ts": "export {}",
				"/project/util.ts":  "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{"**/*"},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  nil, // ReadDirectory returns nil for empty results
		},
		{
			name: "zero depth",
			files: map[string]string{
				"/project/index.ts":    "export {}",
				"/project/src/util.ts": "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			depth:                     func() *int { d := 0; return &d }(),
			expected:                  []string{"/project/index.ts", "/project/src/util.ts"},
		},
		{
			name: "complex wildcard patterns",
			files: map[string]string{
				"/project/src/component.min.js": "console.log('minified')",
				"/project/src/component.js":     "console.log('normal')",
				"/project/src/util.ts":          "export {}",
				"/project/dist/build.min.js":    "console.log('built')",
			},
			path:                      "/project",
			extensions:                []string{".js", ".ts"},
			excludes:                  []string{"**/*.min.js"},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/component.js", "/project/src/util.ts"},
		},
		{
			name: "ignore dotted files and folders",
			files: map[string]string{
				"/apath/..c.ts":        "export {}",
				"/apath/.b.ts":         "export {}",
				"/apath/.git/a.ts":     "export {}",
				"/apath/test.ts":       "export {}",
				"/apath/tsconfig.json": "{}",
			},
			path:                      "/apath",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/apath/test.ts"},
		},
		{
			name: "implicitly exclude common package folders",
			files: map[string]string{
				"/bower_components/b.ts": "export {}",
				"/d.ts":                  "export {}",
				"/folder/e.ts":           "export {}",
				"/jspm_packages/c.ts":    "export {}",
				"/node_modules/a.ts":     "export {}",
				"/tsconfig.json":         "{}",
			},
			path:                      "/",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/d.ts", "/folder/e.ts"},
		},
		{
			name: "comprehensive test case",
			files: map[string]string{
				"/project/src/index.ts":                "export {}",
				"/project/src/util.ts":                 "export {}",
				"/project/src/components/App.tsx":      "export {}",
				"/project/src/types/index.d.ts":        "export {}",
				"/project/tests/unit.test.ts":          "export {}",
				"/project/tests/e2e.spec.ts":           "export {}",
				"/project/node_modules/react/index.js": "module.exports = {}",
				"/project/build/output.js":             "console.log('hello')",
				"/project/docs/readme.md":              "# Project",
				"/project/scripts/deploy.js":           "console.log('deploying')",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx", ".d.ts"},
			excludes:                  []string{"node_modules/**/*", "build/**/*"},
			includes:                  []string{"src/**/*", "tests/**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/components/App.tsx", "/project/src/index.ts", "/project/src/types/index.d.ts", "/project/src/util.ts", "/project/tests/e2e.spec.ts", "/project/tests/unit.test.ts"},
		},
		{
			name: "case insensitive comparison",
			files: map[string]string{
				"/project/SRC/Index.TS":       "export {}",
				"/project/src/Util.ts":        "export {}",
				"/project/Tests/Unit.test.ts": "export {}",
				"/project/BUILD/output.js":    "console.log('hello')",
			},
			path:                      "/project",
			extensions:                []string{".ts", ".tsx"},
			excludes:                  []string{"build/**/*"},
			includes:                  []string{"src/**/*", "tests/**/*"},
			useCaseSensitiveFileNames: false,
			currentDirectory:          "/",
			expected:                  []string{"/project/SRC/Index.TS", "/project/src/Util.ts", "/project/Tests/Unit.test.ts"},
		},
		{
			name: "depth limited comparison",
			files: map[string]string{
				"/project/index.ts":                "export {}",
				"/project/src/util.ts":             "export {}",
				"/project/src/deep/nested/file.ts": "export {}",
				"/project/src/other/file.ts":       "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{"**/*"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			depth:                     func() *int { d := 2; return &d }(),
			expected:                  []string{"/project/index.ts", "/project/src/util.ts"},
		},
		{
			name: "wildcard questions and asterisks",
			files: map[string]string{
				"/project/test1.ts":   "export {}",
				"/project/test2.ts":   "export {}",
				"/project/testAB.ts":  "export {}",
				"/project/another.ts": "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{"test?.ts"},
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			depth:                     func() *int { d := 2; return &d }(),
			expected:                  []string{"/project/test1.ts", "/project/test2.ts"},
		},
		{
			name: "implicit glob behavior",
			files: map[string]string{
				"/project/src/index.ts":    "export {}",
				"/project/src/util.ts":     "export {}",
				"/project/src/sub/file.ts": "export {}",
			},
			path:                      "/project",
			extensions:                []string{".ts"},
			excludes:                  []string{},
			includes:                  []string{"src"}, // Should be treated as src/**/*
			useCaseSensitiveFileNames: true,
			currentDirectory:          "/",
			expected:                  []string{"/project/src/index.ts", "/project/src/sub/file.ts", "/project/src/util.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := vfstest.FromMap(tt.files, tt.useCaseSensitiveFileNames)

			// Test new implementation
			newResult := matchFilesNew(
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.useCaseSensitiveFileNames,
				tt.currentDirectory,
				tt.depth,
				fs,
			)

			// Test old implementation for compatibility
			oldResult := matchFilesOld(
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.useCaseSensitiveFileNames,
				tt.currentDirectory,
				tt.depth,
				fs,
			)

			// Assert the new result matches expected
			assert.DeepEqual(t, newResult, tt.expected)

			// Compatibility check: both implementations should return the same result
			assert.DeepEqual(t, newResult, oldResult)
		})
	}
}

// Test that verifies MatchesExcludeNew and MatchesExcludeOld return the same data
func TestMatchesExclude(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		fileName                  string
		excludeSpecs              []string
		currentDirectory          string
		useCaseSensitiveFileNames bool
		expectExcluded            bool
	}{
		{
			name:                      "no exclude specs",
			fileName:                  "/project/src/index.ts",
			excludeSpecs:              []string{},
			currentDirectory:          "/",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "simple exclude match",
			fileName:                  "/project/node_modules/react/index.js",
			excludeSpecs:              []string{"node_modules/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "exclude does not match",
			fileName:                  "/project/src/index.ts",
			excludeSpecs:              []string{"node_modules/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "multiple exclude patterns",
			fileName:                  "/project/dist/output.js",
			excludeSpecs:              []string{"node_modules/**/*", "dist/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "case insensitive exclude",
			fileName:                  "/project/BUILD/output.js",
			excludeSpecs:              []string{"build/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: false,
			expectExcluded:            true,
		},
		{
			name:                      "extensionless file matches directory pattern",
			fileName:                  "/project/LICENSE",
			excludeSpecs:              []string{"LICENSE/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "complex patterns with wildcards",
			fileName:                  "/project/src/test.spec.ts",
			excludeSpecs:              []string{"**/*.spec.*", "build/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "relative path handling",
			fileName:                  "/project/src/index.ts",
			excludeSpecs:              []string{"./src/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "nested directory exclusion",
			fileName:                  "/project/src/deep/nested/file.ts",
			excludeSpecs:              []string{"src/deep/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "hidden files and directories",
			fileName:                  "/project/.git/config",
			excludeSpecs:              []string{".git/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "partial path match should not exclude",
			fileName:                  "/project/src/node_modules_util.ts",
			excludeSpecs:              []string{"node_modules/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "exact filename match",
			fileName:                  "/project/temp.log",
			excludeSpecs:              []string{"temp.log"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "wildcard in middle of pattern",
			fileName:                  "/project/src/components/Button.test.tsx",
			excludeSpecs:              []string{"src/**/test.*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "wildcard matching file extension",
			fileName:                  "/project/src/components/Button.test.tsx",
			excludeSpecs:              []string{"**/*.test.*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "case sensitive mismatch",
			fileName:                  "/project/BUILD/output.js",
			excludeSpecs:              []string{"build/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "absolute path vs relative exclude",
			fileName:                  "/usr/local/project/src/index.ts",
			excludeSpecs:              []string{"src/**/*"},
			currentDirectory:          "/usr/local/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "empty exclude specs array",
			fileName:                  "/any/path/file.ts",
			excludeSpecs:              []string{},
			currentDirectory:          "/any/path",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "multiple patterns, first matches",
			fileName:                  "/project/node_modules/pkg/index.js",
			excludeSpecs:              []string{"node_modules/**/*", "build/**/*", "dist/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "multiple patterns, last matches",
			fileName:                  "/project/dist/bundle.js",
			excludeSpecs:              []string{"node_modules/**/*", "build/**/*", "dist/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "multiple patterns, none match",
			fileName:                  "/project/src/index.ts",
			excludeSpecs:              []string{"node_modules/**/*", "build/**/*", "dist/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "question mark wildcard matching single char",
			fileName:                  "/project/test1.ts",
			excludeSpecs:              []string{"test?.ts"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "question mark wildcard not matching multiple chars",
			fileName:                  "/project/testAB.ts",
			excludeSpecs:              []string{"test?.ts"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            false,
		},
		{
			name:                      "dot in exclude pattern",
			fileName:                  "/project/.eslintrc.js",
			excludeSpecs:              []string{".eslintrc.*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "directory separator handling",
			fileName:                  "/project/src\\components\\Button.ts", // Backslash in filename
			excludeSpecs:              []string{"src/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
		{
			name:                      "deeply nested exclusion",
			fileName:                  "/project/src/very/deep/nested/directory/file.ts",
			excludeSpecs:              []string{"src/very/**/*"},
			currentDirectory:          "/project",
			useCaseSensitiveFileNames: true,
			expectExcluded:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test new implementation
			newResult := MatchesExclude(tt.fileName, tt.excludeSpecs, tt.currentDirectory, tt.useCaseSensitiveFileNames)

			// Test old implementation for compatibility
			oldResult := matchesExcludeOld(
				tt.fileName,
				tt.excludeSpecs,
				tt.currentDirectory,
				tt.useCaseSensitiveFileNames,
			)

			// Assert the new result matches expected
			assert.Equal(t, tt.expectExcluded, newResult)

			// Compatibility check: both implementations should return the same result
			assert.Equal(t, newResult, oldResult)
		})
	}
}

func TestMatchesInclude(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		fileName                  string
		includeSpecs              []string
		basePath                  string
		useCaseSensitiveFileNames bool
		expectIncluded            bool
	}{
		{
			name:                      "no include specs",
			fileName:                  "/project/src/index.ts",
			includeSpecs:              []string{},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "simple include match",
			fileName:                  "/project/src/index.ts",
			includeSpecs:              []string{"src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "include does not match",
			fileName:                  "/project/tests/unit.test.ts",
			includeSpecs:              []string{"src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "multiple include patterns",
			fileName:                  "/project/tests/unit.test.ts",
			includeSpecs:              []string{"src/**/*", "tests/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "case insensitive include",
			fileName:                  "/project/SRC/Index.ts",
			includeSpecs:              []string{"src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: false,
			expectIncluded:            true,
		},
		{
			name:                      "specific file pattern",
			fileName:                  "/project/package.json",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "min.js files with explicit pattern",
			fileName:                  "/dev/js/d.min.js",
			includeSpecs:              []string{"js/*.min.js"},
			basePath:                  "/dev",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "min.js files should not match generic * pattern",
			fileName:                  "/dev/js/d.min.js",
			includeSpecs:              []string{"js/*"},
			basePath:                  "/dev",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "complex nested patterns",
			fileName:                  "/project/src/components/button/index.tsx",
			includeSpecs:              []string{"src/**/*.tsx", "tests/**/*.test.*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "wildcard file extensions",
			fileName:                  "/project/src/util.ts",
			includeSpecs:              []string{"src/**/*.{ts,tsx,js}"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "question mark pattern",
			fileName:                  "/project/test1.ts",
			includeSpecs:              []string{"test?.ts"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "relative base path",
			fileName:                  "/project/src/index.ts",
			includeSpecs:              []string{"./src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "question mark not matching multiple chars",
			fileName:                  "/project/testAB.ts",
			includeSpecs:              []string{"test?.ts"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "case sensitive mismatch",
			fileName:                  "/project/SRC/Index.ts",
			includeSpecs:              []string{"src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "nested directory pattern",
			fileName:                  "/project/src/deep/nested/file.ts",
			includeSpecs:              []string{"src/deep/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple patterns, first matches",
			fileName:                  "/project/src/index.ts",
			includeSpecs:              []string{"src/**/*", "tests/**/*", "docs/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple patterns, last matches",
			fileName:                  "/project/docs/readme.md",
			includeSpecs:              []string{"src/**/*", "tests/**/*", "docs/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple patterns, none match",
			fileName:                  "/project/build/output.js",
			includeSpecs:              []string{"src/**/*", "tests/**/*", "docs/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "exact filename pattern",
			fileName:                  "/project/README.md",
			includeSpecs:              []string{"README.md"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "star pattern at root",
			fileName:                  "/project/package.json",
			includeSpecs:              []string{"*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "double star pattern",
			fileName:                  "/project/src/components/deep/nested/component.tsx",
			includeSpecs:              []string{"**/component.tsx"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "pattern with brackets",
			fileName:                  "/project/src/util.ts",
			includeSpecs:              []string{"src/**/*.{ts,tsx}"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "pattern with brackets no match",
			fileName:                  "/project/src/util.js",
			includeSpecs:              []string{"src/**/*.{ts,tsx}"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "hidden file inclusion",
			fileName:                  "/project/.gitignore",
			includeSpecs:              []string{".git*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "complex wildcard with negation-like pattern",
			fileName:                  "/project/src/components/!important.ts",
			includeSpecs:              []string{"src/**/*.ts"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "empty base path",
			fileName:                  "/file.ts",
			includeSpecs:              []string{"*.ts"},
			basePath:                  "/",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "file in subdirectory with star pattern",
			fileName:                  "/project/subdir/file.ts",
			includeSpecs:              []string{"*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false, // * doesn't match subdirectories
		},
		{
			name:                      "file with special characters in path",
			fileName:                  "/project/src/file with spaces.ts",
			includeSpecs:              []string{"src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test new implementation
			newResult := MatchesInclude(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)

			// Test old implementation for compatibility
			oldResult := matchesIncludeOld(
				tt.fileName,
				tt.includeSpecs,
				tt.basePath,
				tt.useCaseSensitiveFileNames,
			)

			// Assert the new result matches expected
			assert.Equal(t, tt.expectIncluded, newResult)

			// Compatibility check: both implementations should return the same result
			assert.Equal(t, newResult, oldResult)
		})
	}
}

func TestMatchesIncludeWithJsonOnly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		fileName                  string
		includeSpecs              []string
		basePath                  string
		useCaseSensitiveFileNames bool
		expectIncluded            bool
	}{
		{
			name:                      "no include specs",
			fileName:                  "/project/package.json",
			includeSpecs:              []string{},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "json file matches json pattern",
			fileName:                  "/project/package.json",
			includeSpecs:              []string{"*.json", "src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "non-json file does not match",
			fileName:                  "/project/src/index.ts",
			includeSpecs:              []string{"*.json", "src/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "json file does not match non-json pattern",
			fileName:                  "/project/config.json",
			includeSpecs:              []string{"src/**/*", "tests/**/*"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "nested json file",
			fileName:                  "/project/src/config/app.json",
			includeSpecs:              []string{"src/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple json patterns",
			fileName:                  "/project/tsconfig.json",
			includeSpecs:              []string{"*.json", "config/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "case insensitive json matching",
			fileName:                  "/project/CONFIG.JSON",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: false,
			expectIncluded:            true,
		},
		{
			name:                      "json file with complex pattern",
			fileName:                  "/project/src/data/users.json",
			includeSpecs:              []string{"src/**/*.json", "test/**/*.spec.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "non-json extension ignored",
			fileName:                  "/project/src/util.ts",
			includeSpecs:              []string{"src/**/*.json", "src/**/*.ts"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "json pattern with wildcards",
			fileName:                  "/project/config/dev.json",
			includeSpecs:              []string{"config/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file case sensitive mismatch",
			fileName:                  "/project/CONFIG.JSON",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "json file in deep nested structure",
			fileName:                  "/project/src/components/forms/config/validation.json",
			includeSpecs:              []string{"src/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file with question mark pattern",
			fileName:                  "/project/config1.json",
			includeSpecs:              []string{"config?.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file question mark no match",
			fileName:                  "/project/configAB.json",
			includeSpecs:              []string{"config?.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "multiple json patterns first matches",
			fileName:                  "/project/package.json",
			includeSpecs:              []string{"*.json", "config/**/*.json", "src/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple json patterns last matches",
			fileName:                  "/project/src/assets/manifest.json",
			includeSpecs:              []string{"*.xml", "config/**/*.json", "src/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "multiple json patterns none match",
			fileName:                  "/project/build/config.json",
			includeSpecs:              []string{"*.xml", "config/**/*.json", "src/**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "json file with relative path pattern",
			fileName:                  "/project/src/config.json",
			includeSpecs:              []string{"./src/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file with double star pattern",
			fileName:                  "/project/any/deep/nested/path/settings.json",
			includeSpecs:              []string{"**/settings.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file with bracket pattern",
			fileName:                  "/project/config.json",
			includeSpecs:              []string{"*.{json,xml,yaml}"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "non-json file with bracket pattern",
			fileName:                  "/project/config.txt",
			includeSpecs:              []string{"*.{json,xml,yaml}"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
		{
			name:                      "hidden json file",
			fileName:                  "/project/.config.json",
			includeSpecs:              []string{".*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json file with special characters",
			fileName:                  "/project/config with spaces.json",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "json extension in middle of filename",
			fileName:                  "/project/config.json.backup",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false, // Should not match .json in middle
		},
		{
			name:                      "json file at root with empty base path",
			fileName:                  "/config.json",
			includeSpecs:              []string{"*.json"},
			basePath:                  "/",
			useCaseSensitiveFileNames: true,
			expectIncluded:            true,
		},
		{
			name:                      "typescript definition file should not match json pattern",
			fileName:                  "/project/types.d.ts",
			includeSpecs:              []string{"*.json", "**/*.json"},
			basePath:                  "/project",
			useCaseSensitiveFileNames: true,
			expectIncluded:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test new implementation
			newResult := MatchesIncludeWithJsonOnly(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)

			// Test old implementation for compatibility
			oldResult := matchesIncludeWithJsonOnlyOld(
				tt.fileName,
				tt.includeSpecs,
				tt.basePath,
				tt.useCaseSensitiveFileNames,
			)

			// Assert the new result matches expected
			assert.Equal(t, tt.expectIncluded, newResult)

			// Compatibility check: both implementations should return the same result
			assert.Equal(t, newResult, oldResult)
		})
	}
}

func BenchmarkMatchFiles(b *testing.B) {
	currentDirectory := "/"
	var depth *int = nil

	benchCases := []struct {
		name     string
		path     string
		exts     []string
		excludes []string
		includes []string
		useFS    func(bool) vfs.FS
	}{
		{
			name:     "CommonPattern",
			path:     "/",
			exts:     []string{".ts", ".tsx"},
			excludes: []string{"**/node_modules/**", "**/dist/**", "**/.hidden/**", "**/*.min.js"},
			includes: []string{"src/**/*", "test/**/*.spec.*"},
			useFS:    setupComplexTestFS,
		},
		{
			name:     "SimpleInclude",
			path:     "/src",
			exts:     []string{".ts", ".tsx"},
			excludes: nil,
			includes: []string{"**/*.ts"},
			useFS:    setupComplexTestFS,
		},
		{
			name:     "EmptyIncludes",
			path:     "/src",
			exts:     []string{".ts", ".tsx"},
			excludes: []string{"**/node_modules/**"},
			includes: []string{},
			useFS:    setupComplexTestFS,
		},
		{
			name:     "HiddenDirectories",
			path:     "/",
			exts:     []string{".json"},
			excludes: nil,
			includes: []string{"**/*", ".vscode/*.json"},
			useFS:    setupComplexTestFS,
		},
		{
			name:     "NodeModulesSearch",
			path:     "/",
			exts:     []string{".ts", ".tsx", ".js"},
			excludes: []string{"**/node_modules/m2/**/*"},
			includes: []string{"**/*", "**/node_modules/**/*"},
			useFS:    setupComplexTestFS,
		},
		{
			name:     "LargeFileSystem",
			path:     "/",
			exts:     []string{".ts", ".tsx", ".js"},
			excludes: []string{"**/node_modules/**", "**/dist/**", "**/.hidden/**"},
			includes: []string{"src/**/*", "tests/**/*.spec.*"},
			useFS:    setupLargeTestFS,
		},
	}

	for _, bc := range benchCases {
		// Create the appropriate file system for this benchmark case
		fs := bc.useFS(true)
		// Wrap with cached FS for the benchmark
		fs = cachedvfs.From(fs)

		b.Run(bc.name+"/Original", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				matchFilesOld(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})

		b.Run(bc.name+"/New", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				matchFilesNew(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})
	}
}

// setupTestFS creates a test file system with a specific structure for testing glob patterns
func setupTestFS(useCaseSensitiveFileNames bool) vfs.FS {
	return vfstest.FromMap(map[string]any{
		"/src/foo.ts":                   "export const foo = 1;",
		"/src/bar.ts":                   "export const bar = 2;",
		"/src/baz.tsx":                  "export const baz = 3;",
		"/src/subfolder/qux.ts":         "export const qux = 4;",
		"/src/subfolder/quux.tsx":       "export const quux = 5;",
		"/src/node_modules/lib.ts":      "export const lib = 6;",
		"/src/.hidden/secret.ts":        "export const secret = 7;",
		"/src/test.min.js":              "console.log('minified');",
		"/dist/output.js":               "console.log('output');",
		"/build/temp.ts":                "export const temp = 8;",
		"/test/test1.spec.ts":           "describe('test1', () => {});",
		"/test/test2.spec.tsx":          "describe('test2', () => {});",
		"/test/subfolder/test3.spec.ts": "describe('test3', () => {});",
	}, useCaseSensitiveFileNames)
}

// setupComplexTestFS creates a more complex test file system for additional pattern testing
func setupComplexTestFS(useCaseSensitiveFileNames bool) vfs.FS {
	return vfstest.FromMap(map[string]any{
		// Regular source files
		"/src/index.ts":          "export * from './utils';",
		"/src/utils.ts":          "export function add(a: number, b: number): number { return a + b; }",
		"/src/utils.d.ts":        "export declare function add(a: number, b: number): number;",
		"/src/models/user.ts":    "export interface User { id: string; name: string; }",
		"/src/models/product.ts": "export interface Product { id: string; price: number; }",

		// Nested directories
		"/src/components/button/index.tsx": "export const Button = () => <button>Click me</button>;",
		"/src/components/input/index.tsx":  "export const Input = () => <input />;",
		"/src/components/form/index.tsx":   "export const Form = () => <form></form>;",

		// Test files
		"/tests/unit/utils.test.ts":      "import { add } from '../../src/utils';",
		"/tests/integration/app.test.ts": "import { app } from '../../src/app';",

		// Node modules
		"/node_modules/lodash/index.js":              "// lodash package",
		"/node_modules/react/index.js":               "// react package",
		"/node_modules/typescript/lib/typescript.js": "// typescript package",
		"/node_modules/@types/react/index.d.ts":      "// react types",

		// Various file types
		"/build/index.js":           "console.log('built')",
		"/assets/logo.png":          "binary content",
		"/assets/images/banner.jpg": "binary content",
		"/assets/fonts/roboto.ttf":  "binary content",
		"/.git/HEAD":                "ref: refs/heads/main",
		"/.vscode/settings.json":    "{ \"typescript.enable\": true }",
		"/package.json":             "{ \"name\": \"test-project\" }",
		"/README.md":                "# Test Project",

		// Files with special characters
		"/src/special-case.ts": "export const special = 'case';",
		"/src/[id].ts":         "export const dynamic = (id) => id;",
		"/src/weird.name.ts":   "export const weird = 'name';",
		"/src/problem?.ts":     "export const problem = 'maybe';",
		"/src/with space.ts":   "export const withSpace = 'test';",
	}, useCaseSensitiveFileNames)
}

// setupLargeTestFS creates a test file system with thousands of files for benchmarking
func setupLargeTestFS(useCaseSensitiveFileNames bool) vfs.FS {
	// Create a map to hold all the files
	files := make(map[string]any)

	// Add some standard structure
	files["/src/index.ts"] = "export * from './lib';"
	files["/src/lib.ts"] = "export const VERSION = '1.0.0';"
	files["/package.json"] = "{ \"name\": \"large-test-project\" }"
	files["/.vscode/settings.json"] = "{ \"typescript.enable\": true }"
	files["/node_modules/typescript/package.json"] = "{ \"name\": \"typescript\", \"version\": \"5.0.0\" }"

	// Add 1000 TypeScript files in src/components
	for i := range 1000 {
		files[fmt.Sprintf("/src/components/component%d.ts", i)] = fmt.Sprintf("export const Component%d = () => null;", i)
	}

	// Add 500 TypeScript files in src/utils with nested structure
	for i := range 500 {
		folder := i % 10 // Create 10 different folders
		files[fmt.Sprintf("/src/utils/folder%d/util%d.ts", folder, i)] = fmt.Sprintf("export function util%d() { return %d; }", i, i)
	}

	// Add 500 test files
	for i := range 500 {
		files[fmt.Sprintf("/tests/unit/test%d.spec.ts", i)] = fmt.Sprintf("describe('test%d', () => { it('works', () => {}) });", i)
	}

	// Add 200 files in node_modules with various extensions
	for i := range 200 {
		pkg := i % 20 // Create 20 different packages
		files[fmt.Sprintf("/node_modules/pkg%d/file%d.js", pkg, i)] = fmt.Sprintf("module.exports = { value: %d };", i)

		// Add some .d.ts files
		if i < 50 {
			files[fmt.Sprintf("/node_modules/pkg%d/types/file%d.d.ts", pkg, i)] = "export declare const value: number;"
		}
	}

	// Add 100 files in dist directory (build output)
	for i := range 100 {
		files[fmt.Sprintf("/dist/file%d.js", i)] = fmt.Sprintf("console.log(%d);", i)
	}

	// Add some hidden files
	for i := range 50 {
		files[fmt.Sprintf("/.hidden/file%d.ts", i)] = fmt.Sprintf("// Hidden file %d", i)
	}

	return vfstest.FromMap(files, useCaseSensitiveFileNames)
}

func BenchmarkMatchFilesLarge(b *testing.B) {
	fs := setupLargeTestFS(true)
	// Wrap with cached FS for the benchmark
	fs = cachedvfs.From(fs)
	currentDirectory := "/"
	var depth *int = nil

	benchCases := []struct {
		name     string
		path     string
		exts     []string
		excludes []string
		includes []string
	}{
		{
			name:     "AllFiles",
			path:     "/",
			exts:     []string{".ts", ".tsx", ".js"},
			excludes: []string{"**/node_modules/**", "**/dist/**"},
			includes: []string{"**/*"},
		},
		{
			name:     "Components",
			path:     "/src/components",
			exts:     []string{".ts"},
			excludes: nil,
			includes: []string{"**/*.ts"},
		},
		{
			name:     "TestFiles",
			path:     "/tests",
			exts:     []string{".ts"},
			excludes: nil,
			includes: []string{"**/*.spec.ts"},
		},
		{
			name:     "NestedUtilsWithPattern",
			path:     "/src/utils",
			exts:     []string{".ts"},
			excludes: nil,
			includes: []string{"**/folder*/*.ts"},
		},
	}

	for _, bc := range benchCases {
		b.Run(bc.name+"/Original", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				matchFilesOld(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})

		b.Run(bc.name+"/New", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				matchFilesNew(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})
	}
}

func TestIsImplicitGlob(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		lastPathComponent  string
		expectImplicitGlob bool
	}{
		{
			name:               "simple directory name",
			lastPathComponent:  "src",
			expectImplicitGlob: true,
		},
		{
			name:               "file with extension",
			lastPathComponent:  "index.ts",
			expectImplicitGlob: false,
		},
		{
			name:               "pattern with asterisk",
			lastPathComponent:  "*.ts",
			expectImplicitGlob: false,
		},
		{
			name:               "pattern with question mark",
			lastPathComponent:  "test?.ts",
			expectImplicitGlob: false,
		},
		{
			name:               "hidden file",
			lastPathComponent:  ".hidden",
			expectImplicitGlob: false,
		},
		{
			name:               "empty string",
			lastPathComponent:  "",
			expectImplicitGlob: true,
		},
		{
			name:               "multiple dots",
			lastPathComponent:  "file.min.js",
			expectImplicitGlob: false,
		},
		{
			name:               "directory with special chars but no glob",
			lastPathComponent:  "my-folder_name",
			expectImplicitGlob: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsImplicitGlob(tt.lastPathComponent)
			assert.Equal(t, tt.expectImplicitGlob, result)
		})
	}
}
