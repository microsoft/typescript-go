package vfs_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// Test cases based on real-world patterns found in the TypeScript codebase
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := vfstest.FromMap(tt.files, tt.useCaseSensitiveFileNames)

			result := vfs.ReadDirectory(
				fs,
				tt.currentDirectory,
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.depth,
			)

			assert.DeepEqual(t, result, tt.expected)
		})
	}
}

// Test edge cases and error conditions
func TestMatchFilesEdgeCases(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := vfstest.FromMap(tt.files, tt.useCaseSensitiveFileNames)

			result := vfs.ReadDirectory(
				fs,
				tt.currentDirectory,
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.depth,
			)

			assert.DeepEqual(t, result, tt.expected)
		})
	}
}

// Test that verifies matchFiles and matchFilesNew return the same data
func TestMatchFilesCompatibility(t *testing.T) {
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
	}{
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := vfstest.FromMap(tt.files, tt.useCaseSensitiveFileNames)

			// Get results from original implementation
			originalResult := vfs.ReadDirectory(
				fs,
				tt.currentDirectory,
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.depth,
			)

			// Get results from new implementation
			newResult := vfs.MatchFilesNew(
				tt.path,
				tt.extensions,
				tt.excludes,
				tt.includes,
				tt.useCaseSensitiveFileNames,
				tt.currentDirectory,
				tt.depth,
				fs,
			)

			assert.DeepEqual(t, originalResult, newResult)

			// For now, just verify the original implementation works
			assert.Assert(t, originalResult != nil, "original implementation should not return nil")
		})
	}
}
