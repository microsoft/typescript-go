package vfsmatch

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
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
			// !!! This seems wrong, but the old implementation behaved this way.
			expected: []string{
				"/apath/..c.ts",
				"/apath/.b.ts",
				"/apath/test.ts",
				"/apath/.git/a.ts",
			},
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
			// !!! This seems wrong, but the old implementation behaved this way.
			expected: []string{
				"/d.ts",
				"/bower_components/b.ts",
				"/folder/e.ts",
				"/jspm_packages/c.ts",
				"/node_modules/a.ts",
			},
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
			expected: []string{
				"/project/src/index.ts",
				"/project/src/util.ts",
				"/project/src/components/App.tsx",
				"/project/src/types/index.d.ts",
				"/project/tests/e2e.spec.ts",
				"/project/tests/unit.test.ts",
			},
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
			expected:                  []string{"/project/SRC/Util.ts", "/project/Tests/Unit.test.ts"},
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
			expected:                  []string{"/project/src/index.ts", "/project/src/util.ts", "/project/src/sub/file.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := vfstest.FromMap(tt.files, tt.useCaseSensitiveFileNames)

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
			assert.Check(t, cmp.DeepEqual(oldResult, tt.expected))

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
			assert.Check(t, cmp.DeepEqual(newResult, tt.expected))

			mainResult := ReadDirectory(
				fs,
				tt.currentDirectory,
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.depth,
			)
			assert.Check(t, cmp.DeepEqual(mainResult, tt.expected))
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
			expectExcluded:            false, // Surprising, but Strada's code does not match this.
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

			oldResult := matchesExcludeOld(tt.fileName, tt.excludeSpecs, tt.currentDirectory, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(oldResult, tt.expectExcluded))

			newResult := matchesExcludeNew(tt.fileName, tt.excludeSpecs, tt.currentDirectory, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(newResult, tt.expectExcluded))

			mainResult := MatchesExclude(tt.fileName, tt.excludeSpecs, tt.currentDirectory, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(mainResult, tt.expectExcluded))
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
			expectIncluded:            false,
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
			expectIncluded:            false,
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

			oldResult := matchesIncludeOld(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(oldResult, tt.expectIncluded))

			newResult := matchesIncludeNew(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(newResult, tt.expectIncluded))

			mainResult := MatchesInclude(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(mainResult, tt.expectIncluded))
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
			expectIncluded:            false,
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

			oldResult := matchesIncludeWithJsonOnlyOld(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(oldResult, tt.expectIncluded))

			newResult := matchesIncludeWithJsonOnlyNew(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(newResult, tt.expectIncluded))

			mainResult := MatchesIncludeWithJsonOnly(tt.fileName, tt.includeSpecs, tt.basePath, tt.useCaseSensitiveFileNames)
			assert.Check(t, cmp.Equal(mainResult, tt.expectIncluded))
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
			assert.Equal(t, tt.expectImplicitGlob, IsImplicitGlob(tt.lastPathComponent))
		})
	}
}
