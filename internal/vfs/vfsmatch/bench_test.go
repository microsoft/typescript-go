package vfsmatch

import (
	"fmt"
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

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
