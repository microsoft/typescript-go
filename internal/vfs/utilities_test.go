package vfs_test

import (
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
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

func TestMatchFilesImplicitExclusions(t *testing.T) {
	t.Parallel()

	t.Run("ignore dotted files and folders", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"/apath/..c.ts":        "export {}",
			"/apath/.b.ts":         "export {}",
			"/apath/.git/a.ts":     "export {}",
			"/apath/test.ts":       "export {}",
			"/apath/tsconfig.json": "{}",
		}
		fs := vfstest.FromMap(files, true)

		// This should only return test.ts, not the dotted files
		result := vfs.MatchFilesNew(
			"/apath",
			[]string{".ts"},
			[]string{}, // no explicit excludes
			[]string{}, // no explicit includes - should include all
			true,
			"/",
			nil,
			fs,
		)

		expected := []string{"/apath/test.ts"}
		assert.DeepEqual(t, result, expected)
	})

	t.Run("implicitly exclude common package folders", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"/bower_components/b.ts": "export {}",
			"/d.ts":                  "export {}",
			"/folder/e.ts":           "export {}",
			"/jspm_packages/c.ts":    "export {}",
			"/node_modules/a.ts":     "export {}",
			"/tsconfig.json":         "{}",
		}
		fs := vfstest.FromMap(files, true)

		// This should only return d.ts and folder/e.ts, not the package folders
		result := vfs.MatchFilesNew(
			"/",
			[]string{".ts"},
			[]string{}, // no explicit excludes
			[]string{}, // no explicit includes - should include all
			true,
			"/",
			nil,
			fs,
		)

		expected := []string{"/d.ts", "/folder/e.ts"}
		assert.DeepEqual(t, result, expected)
	})
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

// Test specific patterns that were originally in debug test
func TestDottedFilesAndPackageFolders(t *testing.T) {
	t.Parallel()

	t.Run("ignore dotted files and folders", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"/apath/..c.ts":        "export {}",
			"/apath/.b.ts":         "export {}",
			"/apath/.git/a.ts":     "export {}",
			"/apath/test.ts":       "export {}",
			"/apath/tsconfig.json": "{}",
		}
		fs := vfstest.FromMap(files, true)

		// Test the new implementation
		result := vfs.MatchFilesNew(
			"/apath",
			[]string{".ts"},
			[]string{}, // no explicit excludes
			[]string{}, // no explicit includes - should include all
			true,
			"/",
			nil,
			fs,
		)

		// Based on TypeScript behavior, dotted files should be excluded
		expected := []string{"/apath/test.ts"}
		assert.DeepEqual(t, result, expected)
	})

	t.Run("implicitly exclude common package folders", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"/bower_components/b.ts": "export {}",
			"/d.ts":                  "export {}",
			"/folder/e.ts":           "export {}",
			"/jspm_packages/c.ts":    "export {}",
			"/node_modules/a.ts":     "export {}",
			"/tsconfig.json":         "{}",
		}
		fs := vfstest.FromMap(files, true)

		// This should only return d.ts and folder/e.ts, not the package folders
		result := vfs.MatchFilesNew(
			"/",
			[]string{".ts"},
			[]string{}, // no explicit excludes
			[]string{}, // no explicit includes - should include all
			true,
			"/",
			nil,
			fs,
		)

		expected := []string{"/d.ts", "/folder/e.ts"}
		assert.DeepEqual(t, result, expected)
	})
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
				vfs.MatchFiles(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})

		b.Run(bc.name+"/New", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				vfs.MatchFilesNew(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
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
	for i := 0; i < 1000; i++ {
		files[fmt.Sprintf("/src/components/component%d.ts", i)] = fmt.Sprintf("export const Component%d = () => null;", i)
	}

	// Add 500 TypeScript files in src/utils with nested structure
	for i := 0; i < 500; i++ {
		folder := i % 10 // Create 10 different folders
		files[fmt.Sprintf("/src/utils/folder%d/util%d.ts", folder, i)] = fmt.Sprintf("export function util%d() { return %d; }", i, i)
	}

	// Add 500 test files
	for i := 0; i < 500; i++ {
		files[fmt.Sprintf("/tests/unit/test%d.spec.ts", i)] = fmt.Sprintf("describe('test%d', () => { it('works', () => {}) });", i)
	}

	// Add 200 files in node_modules with various extensions
	for i := 0; i < 200; i++ {
		pkg := i % 20 // Create 20 different packages
		files[fmt.Sprintf("/node_modules/pkg%d/file%d.js", pkg, i)] = fmt.Sprintf("module.exports = { value: %d };", i)

		// Add some .d.ts files
		if i < 50 {
			files[fmt.Sprintf("/node_modules/pkg%d/types/file%d.d.ts", pkg, i)] = "export declare const value: number;"
		}
	}

	// Add 100 files in dist directory (build output)
	for i := 0; i < 100; i++ {
		files[fmt.Sprintf("/dist/file%d.js", i)] = fmt.Sprintf("console.log(%d);", i)
	}

	// Add some hidden files
	for i := 0; i < 50; i++ {
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
				vfs.MatchFiles(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})

		b.Run(bc.name+"/New", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				vfs.MatchFilesNew(bc.path, bc.exts, bc.excludes, bc.includes, fs.UseCaseSensitiveFileNames(), currentDirectory, depth, fs)
			}
		})
	}
}
