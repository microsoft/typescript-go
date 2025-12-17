package vfsmatch

import (
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// Test cases modeled after TypeScript's matchFiles tests in
// _submodules/TypeScript/src/testRunner/unittests/config/matchFiles.ts

func ptrTo[T any](v T) *T {
	return &v
}

// readDirectoryFunc is a function type for ReadDirectory implementations
type readDirectoryFunc func(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string

// readDirectoryOld wraps matchFiles with the expected test signature
func readDirectoryOld(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return matchFiles(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// readDirectoryNew wraps matchFilesNoRegex with the expected test signature
func readDirectoryNew(host vfs.FS, currentDir string, path string, extensions []string, excludes []string, includes []string, depth *int) []string {
	return matchFilesNoRegex(path, extensions, excludes, includes, host.UseCaseSensitiveFileNames(), currentDir, depth, host)
}

// readDirectoryImplementations contains all implementations to test
var readDirectoryImplementations = []struct {
	name string
	fn   readDirectoryFunc
}{
	{"Old", readDirectoryOld},
	{"New", readDirectoryNew},
}

// caseInsensitiveHost simulates a Windows-like file system
func caseInsensitiveHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/a.ts":         "",
		"/dev/a.d.ts":       "",
		"/dev/a.js":         "",
		"/dev/b.ts":         "",
		"/dev/b.js":         "",
		"/dev/c.d.ts":       "",
		"/dev/z/a.ts":       "",
		"/dev/z/abz.ts":     "",
		"/dev/z/aba.ts":     "",
		"/dev/z/b.ts":       "",
		"/dev/z/bbz.ts":     "",
		"/dev/z/bba.ts":     "",
		"/dev/x/a.ts":       "",
		"/dev/x/aa.ts":      "",
		"/dev/x/b.ts":       "",
		"/dev/x/y/a.ts":     "",
		"/dev/x/y/b.ts":     "",
		"/dev/js/a.js":      "",
		"/dev/js/b.js":      "",
		"/dev/js/d.min.js":  "",
		"/dev/js/ab.min.js": "",
		"/ext/ext.ts":       "",
		"/ext/b/a..b.ts":    "",
	}, false)
}

// caseSensitiveHost simulates a Unix-like case-sensitive file system
func caseSensitiveHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/a.ts":         "",
		"/dev/a.d.ts":       "",
		"/dev/a.js":         "",
		"/dev/b.ts":         "",
		"/dev/b.js":         "",
		"/dev/A.ts":         "",
		"/dev/B.ts":         "",
		"/dev/c.d.ts":       "",
		"/dev/z/a.ts":       "",
		"/dev/z/abz.ts":     "",
		"/dev/z/aba.ts":     "",
		"/dev/z/b.ts":       "",
		"/dev/z/bbz.ts":     "",
		"/dev/z/bba.ts":     "",
		"/dev/x/a.ts":       "",
		"/dev/x/b.ts":       "",
		"/dev/x/y/a.ts":     "",
		"/dev/x/y/b.ts":     "",
		"/dev/q/a/c/b/d.ts": "",
		"/dev/js/a.js":      "",
		"/dev/js/b.js":      "",
	}, true)
}

// commonFoldersHost includes node_modules, bower_components, jspm_packages
func commonFoldersHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/a.ts":                  "",
		"/dev/a.d.ts":                "",
		"/dev/a.js":                  "",
		"/dev/b.ts":                  "",
		"/dev/x/a.ts":                "",
		"/dev/node_modules/a.ts":     "",
		"/dev/bower_components/a.ts": "",
		"/dev/jspm_packages/a.ts":    "",
	}, false)
}

// dottedFoldersHost includes files and folders starting with a dot
func dottedFoldersHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/x/d.ts":           "",
		"/dev/x/y/d.ts":         "",
		"/dev/x/y/.e.ts":        "",
		"/dev/x/.y/a.ts":        "",
		"/dev/.z/.b.ts":         "",
		"/dev/.z/c.ts":          "",
		"/dev/w/.u/e.ts":        "",
		"/dev/g.min.js/.g/g.ts": "",
	}, false)
}

// mixedExtensionHost has various file extensions
func mixedExtensionHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/a.ts":    "",
		"/dev/a.d.ts":  "",
		"/dev/a.js":    "",
		"/dev/b.tsx":   "",
		"/dev/b.d.ts":  "",
		"/dev/b.jsx":   "",
		"/dev/c.tsx":   "",
		"/dev/c.js":    "",
		"/dev/d.js":    "",
		"/dev/e.jsx":   "",
		"/dev/f.other": "",
	}, false)
}

// sameNamedDeclarationsHost has files with same names but different extensions
func sameNamedDeclarationsHost() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/dev/a.tsx":  "",
		"/dev/a.d.ts": "",
		"/dev/b.tsx":  "",
		"/dev/b.ts":   "",
		"/dev/c.tsx":  "",
		"/dev/m.ts":   "",
		"/dev/m.d.ts": "",
		"/dev/n.tsx":  "",
		"/dev/n.ts":   "",
		"/dev/n.d.ts": "",
		"/dev/o.ts":   "",
		"/dev/x.d.ts": "",
	}, false)
}

type readDirTestCase struct {
	name       string
	host       func() vfs.FS
	currentDir string
	path       string
	extensions []string
	excludes   []string
	includes   []string
	depth      *int
	expect     func(t *testing.T, got []string)
}

func runReadDirectoryCase(t *testing.T, tc readDirTestCase, readDir readDirectoryFunc) {
	currentDir := tc.currentDir
	if currentDir == "" {
		currentDir = "/"
	}
	path := tc.path
	if path == "" {
		path = "/dev"
	}
	got := ReadDirectory(tc.host(), currentDir, path, tc.extensions, tc.excludes, tc.includes, tc.depth)
	tc.expect(t, got)
}

func TestReadDirectory(t *testing.T) {
	t.Parallel()

	cases := []readDirTestCase{
		{
			name:       "defaults include common package folders",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/node_modules/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/bower_components/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/jspm_packages/a.ts"))
			},
		},
		{
			name:       "literal includes without exclusions",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"a.ts", "b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/a.ts", "/dev/b.ts"})
			},
		},
		{
			name:       "literal includes with non ts extensions excluded",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"a.js", "b.js"},
			expect: func(t *testing.T, got []string) {
				assert.Equal(t, len(got), 0)
			},
		},
		{
			name:       "literal includes missing files excluded",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"z.ts", "x.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Equal(t, len(got), 0)
			},
		},
		{
			name:       "literal includes with literal excludes",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"b.ts"},
			includes:   []string{"a.ts", "b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/a.ts"})
			},
		},
		{
			name:       "literal includes with wildcard excludes",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"*.ts", "z/??z.ts", "*/b.ts"},
			includes:   []string{"a.ts", "b.ts", "z/a.ts", "z/abz.ts", "z/aba.ts", "x/b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/z/a.ts", "/dev/z/aba.ts"})
			},
		},
		{
			name:       "literal includes with recursive excludes",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**/b.ts"},
			includes:   []string{"a.ts", "b.ts", "x/a.ts", "x/b.ts", "x/y/a.ts", "x/y/b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/a.ts", "/dev/x/a.ts", "/dev/x/y/a.ts"})
			},
		},
		{
			name:       "case sensitive exclude is respected",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**/b.ts"},
			includes:   []string{"B.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/B.ts"})
			},
		},
		{
			name:       "explicit includes keep common package folders",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"a.ts", "b.ts", "node_modules/a.ts", "bower_components/a.ts", "jspm_packages/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/node_modules/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/bower_components/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/jspm_packages/a.ts"))
			},
		},
		{
			name:       "wildcard include sorted order",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"z/*.ts", "x/*.ts"},
			expect: func(t *testing.T, got []string) {
				expected := []string{
					"/dev/z/a.ts", "/dev/z/aba.ts", "/dev/z/abz.ts", "/dev/z/b.ts", "/dev/z/bba.ts", "/dev/z/bbz.ts",
					"/dev/x/a.ts", "/dev/x/aa.ts", "/dev/x/b.ts",
				}
				assert.DeepEqual(t, got, expected)
			},
		},
		{
			name:       "wildcard include same named declarations excluded",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/a.d.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/c.d.ts"))
			},
		},
		{
			name:       "wildcard star matches only ts files",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, contains(f, ".ts") || contains(f, ".tsx") || contains(f, ".d.ts"), "unexpected file: %s", f)
				}
				assert.Assert(t, !slices.Contains(got, "/dev/a.js"))
				assert.Assert(t, !slices.Contains(got, "/dev/b.js"))
			},
		},
		{
			name:       "wildcard question mark single character",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/?.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/x/a.ts", "/dev/x/b.ts"})
			},
		},
		{
			name:       "wildcard recursive directory",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/z/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/a.ts"))
			},
		},
		{
			name:       "wildcard multiple recursive directories",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/y/**/a.ts", "x/**/a.ts", "z/**/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, len(got) > 0)
			},
		},
		{
			name:       "wildcard case sensitive matching",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/A.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/A.ts"})
			},
		},
		{
			name:       "wildcard missing files excluded",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*/z.ts"},
			expect:     func(t *testing.T, got []string) { assert.Equal(t, len(got), 0) },
		},
		{
			name:       "exclude folders with wildcards",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"z", "x"},
			includes:   []string{"**/*"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, !contains(f, "/z/") && !contains(f, "/x/"), "should not contain z or x: %s", f)
				}
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
			},
		},
		{
			name:       "include paths outside project absolute",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*", "/ext/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/ext/ext.ts"))
			},
		},
		{
			name:       "include paths outside project relative",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**"},
			includes:   []string{"*", "../ext/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/ext/ext.ts"))
			},
		},
		{
			name:       "include files containing double dots",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**"},
			includes:   []string{"/ext/b/a..b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/ext/b/a..b.ts"))
			},
		},
		{
			name:       "exclude files containing double dots",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"/ext/b/a..b.ts"},
			includes:   []string{"/ext/**/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/ext/ext.ts"))
				assert.Assert(t, !slices.Contains(got, "/ext/b/a..b.ts"))
			},
		},
		{
			name:       "common package folders implicitly excluded",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/node_modules/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/bower_components/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/jspm_packages/a.ts"))
			},
		},
		{
			name:       "common package folders explicit recursive include",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/a.ts", "**/node_modules/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/node_modules/a.ts"))
			},
		},
		{
			name:       "common package folders wildcard include",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/node_modules/a.ts"))
			},
		},
		{
			name:       "common package folders explicit wildcard include",
			host:       commonFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*/a.ts", "node_modules/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/node_modules/a.ts"))
			},
		},
		{
			name:       "dotted folders not implicitly included",
			host:       dottedFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/**/*", "w/*/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/d.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/d.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/.y/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/y/.e.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/w/.u/e.ts"))
			},
		},
		{
			name:       "dotted folders explicitly included",
			host:       dottedFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/.y/a.ts", "/dev/.z/.b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/.y/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/.z/.b.ts"))
			},
		},
		{
			name:       "dotted folders recursive wildcard matches directories",
			host:       dottedFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/.*/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/.y/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/.z/c.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/w/.u/e.ts"))
			},
		},
		{
			name:       "trailing recursive include returns empty",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**"},
			expect:     func(t *testing.T, got []string) { assert.Equal(t, len(got), 0) },
		},
		{
			name:       "trailing recursive exclude removes everything",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**"},
			includes:   []string{"**/*"},
			expect:     func(t *testing.T, got []string) { assert.Equal(t, len(got), 0) },
		},
		{
			name:       "multiple recursive directory patterns in includes",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/x/**/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/a.ts"))
			},
		},
		{
			name:       "multiple recursive directory patterns in excludes",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**/x/**"},
			includes:   []string{"**/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/z/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/y/a.ts"))
			},
		},
		{
			name:       "implicit globbification expands directory",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"z"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/z/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/z/aba.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/z/b.ts"))
			},
		},
		{
			name:       "exclude patterns starting with starstar",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**/x"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, !contains(f, "/x/"), "should not contain /x/: %s", f)
				}
			},
		},
		{
			name:       "include patterns starting with starstar",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/x", "**/a/**/b"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/q/a/c/b/d.ts"))
			},
		},
		{
			name:       "depth limit one",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			depth:      ptrTo(1),
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					suffix := f[len("/dev/"):]
					assert.Assert(t, !contains(suffix, "/"), "depth 1 should not include nested files: %s", f)
				}
			},
		},
		{
			name:       "depth limit two",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			depth:      ptrTo(2),
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/z/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/y/a.ts"))
			},
		},
		{
			name:       "mixed extensions only ts",
			host:       mixedExtensionHost,
			extensions: []string{".ts"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, hasSuffix(f, ".ts"), "should only have .ts files: %s", f)
				}
			},
		},
		{
			name:       "mixed extensions ts and tsx",
			host:       mixedExtensionHost,
			extensions: []string{".ts", ".tsx"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, hasSuffix(f, ".ts") || hasSuffix(f, ".tsx"), "should only have .ts or .tsx files: %s", f)
				}
			},
		},
		{
			name:       "mixed extensions js and jsx",
			host:       mixedExtensionHost,
			extensions: []string{".js", ".jsx"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, hasSuffix(f, ".js") || hasSuffix(f, ".jsx"), "should only have .js or .jsx files: %s", f)
				}
			},
		},
		{
			name:       "min js files excluded by wildcard",
			host:       caseInsensitiveHost,
			extensions: []string{".js"},
			includes:   []string{"js/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/js/a.js"))
				assert.Assert(t, slices.Contains(got, "/dev/js/b.js"))
				assert.Assert(t, !slices.Contains(got, "/dev/js/d.min.js"))
				assert.Assert(t, !slices.Contains(got, "/dev/js/ab.min.js"))
			},
		},
		{
			name:       "min js files explicitly included",
			host:       caseInsensitiveHost,
			extensions: []string{".js"},
			includes:   []string{"js/*.min.js"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/js/d.min.js"))
				assert.Assert(t, slices.Contains(got, "/dev/js/ab.min.js"))
			},
		},
		{
			name:       "same named declarations include ts",
			host:       sameNamedDeclarationsHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, len(got) > 0) },
		},
		{
			name:       "same named declarations include tsx",
			host:       sameNamedDeclarationsHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*.tsx"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, hasSuffix(f, ".tsx"), "should only have .tsx files: %s", f)
				}
			},
		},
		{
			name:       "empty includes returns all matching files",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, len(got) > 0)
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
			},
		},
		{
			name: "nil extensions returns all files",
			host: caseInsensitiveHost,
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/a.js"))
			},
		},
		{
			name:       "empty extensions slice returns all files",
			host:       caseInsensitiveHost,
			extensions: []string{},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, len(got) > 0, "expected files to be returned") },
		},
	}

	for _, tc := range cases {
		for _, impl := range readDirectoryImplementations {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				runReadDirectoryCase(t, tc, impl.fn)
			})
		}
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// Additional tests for helper functions

func TestIsImplicitGlob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "simple", input: "foo", expected: true},
		{name: "folder", input: "src", expected: true},
		{name: "with extension", input: "foo.ts", expected: false},
		{name: "trailing dot", input: "foo.", expected: false},
		{name: "star", input: "*", expected: false},
		{name: "question", input: "?", expected: false},
		{name: "star suffix", input: "foo*", expected: false},
		{name: "question suffix", input: "foo?", expected: false},
		{name: "dot name", input: "foo.bar", expected: false},
		{name: "empty", input: "", expected: true},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := IsImplicitGlob(tc.input)
			assert.Equal(t, result, tc.expected)
		})
	}
}

func TestGetRegularExpressionForWildcard(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		specs    []string
		usage    Usage
		expected string
		assertFn func(t *testing.T, got string)
	}{
		{name: "nil specs", specs: nil, usage: UsageFiles, expected: "", assertFn: func(t *testing.T, got string) { assert.Equal(t, got, "") }},
		{name: "empty specs", specs: []string{}, usage: UsageFiles, expected: "", assertFn: func(t *testing.T, got string) { assert.Equal(t, got, "") }},
		{name: "single spec", specs: []string{"*.ts"}, usage: UsageFiles, assertFn: func(t *testing.T, got string) { assert.Assert(t, got != "") }},
		{name: "multiple specs", specs: []string{"*.ts", "*.tsx"}, usage: UsageFiles, assertFn: func(t *testing.T, got string) { assert.Assert(t, got != "") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := getRegularExpressionForWildcard(tc.specs, "/", tc.usage)
			if tc.assertFn != nil {
				tc.assertFn(t, result)
			} else {
				assert.Equal(t, result, tc.expected)
			}
		})
	}
}

func TestGetRegularExpressionsForWildcards(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		specs    []string
		usage    Usage
		assertFn func(t *testing.T, got []string)
	}{
		{name: "nil specs", specs: nil, usage: UsageFiles, assertFn: func(t *testing.T, got []string) { assert.Assert(t, got == nil) }},
		{name: "empty specs", specs: []string{}, usage: UsageFiles, assertFn: func(t *testing.T, got []string) { assert.Assert(t, got == nil) }},
		{name: "two specs", specs: []string{"*.ts", "*.tsx"}, usage: UsageFiles, assertFn: func(t *testing.T, got []string) { assert.Equal(t, len(got), 2) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := getRegularExpressionsForWildcards(tc.specs, "/", tc.usage)
			tc.assertFn(t, result)
		})
	}
}

func TestGetPatternFromSpec(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		spec     string
		usage    Usage
		assertFn func(t *testing.T, got string)
	}{
		{name: "files usage", spec: "*.ts", usage: UsageFiles, assertFn: func(t *testing.T, got string) {
			assert.Assert(t, got != "")
			assert.Assert(t, hasSuffix(got, "$"))
		}},
		{name: "directories usage", spec: "src", usage: UsageDirectories, assertFn: func(t *testing.T, got string) { assert.Assert(t, got != "") }},
		{name: "exclude usage", spec: "node_modules", usage: UsageExclude, assertFn: func(t *testing.T, got string) {
			assert.Assert(t, got != "")
			assert.Assert(t, contains(got, "($|/)"))
		}},
		{name: "trailing starstar non exclude", spec: "**", usage: UsageFiles, assertFn: func(t *testing.T, got string) { assert.Equal(t, got, "") }},
		{name: "trailing starstar exclude allowed", spec: "**", usage: UsageExclude, assertFn: func(t *testing.T, got string) { assert.Assert(t, got != "") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := getPatternFromSpec(tc.spec, "/", tc.usage)
			tc.assertFn(t, result)
		})
	}
}

// Edge case tests for various pattern scenarios
func TestReadDirectoryEdgeCases(t *testing.T) {
	t.Parallel()

	cases := []readDirTestCase{
		{
			name:       "rooted include path",
			host:       caseInsensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"/dev/a.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/dev/a.ts")) },
		},
		{
			name:       "include with extension in path",
			host:       caseInsensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"a.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/dev/a.ts")) },
		},
		{
			name: "special regex characters in path",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/file+test.ts":  "",
					"/dev/file[0].ts":    "",
					"/dev/file(1).ts":    "",
					"/dev/file$money.ts": "",
					"/dev/file^start.ts": "",
					"/dev/file|pipe.ts":  "",
					"/dev/file#hash.ts":  "",
				}, false)
			},
			extensions: []string{".ts"},
			includes:   []string{"file+test.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/dev/file+test.ts")) },
		},
		{
			name:       "include pattern starting with question mark",
			host:       caseInsensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"?.ts"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
			},
		},
		{
			name:       "include pattern starting with star",
			host:       caseInsensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"*b.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/dev/b.ts")) },
		},
		{
			name: "case insensitive file matching",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/File.ts": "",
					"/dev/FILE.ts": "",
				}, true)
			},
			extensions: []string{".ts"},
			includes:   []string{"*.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, len(got) == 2) },
		},
		{
			name:       "nested subdirectory base path",
			host:       caseSensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"q/a/c/b/d.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/dev/q/a/c/b/d.ts")) },
		},
		{
			name:       "current directory differs from path",
			host:       caseInsensitiveHost,
			extensions: []string{".ts"},
			includes:   []string{"z/*.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, len(got) > 0) },
		},
	}

	for _, tc := range cases {
		for _, impl := range readDirectoryImplementations {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				runReadDirectoryCase(t, tc, impl.fn)
			})
		}
	}
}

func TestReadDirectoryEmptyIncludes(t *testing.T) {
	t.Parallel()
	cases := []readDirTestCase{
		{
			name: "empty includes slice behavior",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/root/a.ts": "",
				}, true)
			},
			path:       "/root",
			currentDir: "/",
			extensions: []string{".ts"},
			includes:   []string{},
			expect: func(t *testing.T, got []string) {
				if len(got) == 0 {
					return
				}
				assert.Assert(t, slices.Contains(got, "/root/a.ts"))
			},
		},
	}

	for _, tc := range cases {
		for _, impl := range readDirectoryImplementations {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				runReadDirectoryCase(t, tc, impl.fn)
			})
		}
	}
}

// TestReadDirectorySymlinkCycle tests that cyclic symlinks don't cause infinite loops.
// The cycle is detected by the vfs package using Realpath for cycle detection.
// This means directories with cyclic symlinks will be skipped during traversal.
func TestReadDirectorySymlinkCycle(t *testing.T) {
	t.Parallel()
	cases := []readDirTestCase{
		{
			name: "detects and skips symlink cycles",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]any{
					"/root/file.ts":   "",
					"/root/a/file.ts": "",
					"/root/a/b":       vfstest.Symlink("/root/a"),
				}, true)
			},
			path:       "/root",
			currentDir: "/",
			extensions: []string{".ts"},
			includes:   []string{"**/*"},
			expect: func(t *testing.T, got []string) {
				expected := []string{"/root/file.ts", "/root/a/file.ts"}
				assert.DeepEqual(t, got, expected)
			},
		},
	}

	for _, tc := range cases {
		for _, impl := range readDirectoryImplementations {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				runReadDirectoryCase(t, tc, impl.fn)
			})
		}
	}
}

// TestReadDirectoryMatchesTypeScriptBaselines contains tests that verify the Go implementation
// matches the TypeScript baseline outputs from _submodules/TypeScript/tests/baselines/reference/config/matchFiles/
func TestReadDirectoryMatchesTypeScriptBaselines(t *testing.T) {
	t.Parallel()

	cases := []readDirTestCase{
		{
			name: "sorted in include order then alphabetical",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/z/a.ts":   "",
					"/dev/z/aba.ts": "",
					"/dev/z/abz.ts": "",
					"/dev/z/b.ts":   "",
					"/dev/z/bba.ts": "",
					"/dev/z/bbz.ts": "",
					"/dev/x/a.ts":   "",
					"/dev/x/aa.ts":  "",
					"/dev/x/b.ts":   "",
				}, false)
			},
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"z/*.ts", "x/*.ts"},
			expect: func(t *testing.T, got []string) {
				expected := []string{
					"/dev/z/a.ts", "/dev/z/aba.ts", "/dev/z/abz.ts", "/dev/z/b.ts", "/dev/z/bba.ts", "/dev/z/bbz.ts",
					"/dev/x/a.ts", "/dev/x/aa.ts", "/dev/x/b.ts",
				}
				assert.DeepEqual(t, got, expected)
			},
		},
		{
			name: "recursive wildcards match dotted directories",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/x/d.ts":           "",
					"/dev/x/y/d.ts":         "",
					"/dev/x/y/.e.ts":        "",
					"/dev/x/.y/a.ts":        "",
					"/dev/.z/.b.ts":         "",
					"/dev/.z/c.ts":          "",
					"/dev/w/.u/e.ts":        "",
					"/dev/g.min.js/.g/g.ts": "",
				}, false)
			},
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/.*/*"},
			expect: func(t *testing.T, got []string) {
				expected := []string{"/dev/.z/c.ts", "/dev/g.min.js/.g/g.ts", "/dev/w/.u/e.ts", "/dev/x/.y/a.ts"}
				assert.Equal(t, len(got), len(expected))
				for _, want := range expected {
					assert.Assert(t, slices.Contains(got, want))
				}
			},
		},
		{
			name: "common package folders implicitly excluded with wildcard",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/a.ts":                  "",
					"/dev/a.d.ts":                "",
					"/dev/a.js":                  "",
					"/dev/b.ts":                  "",
					"/dev/x/a.ts":                "",
					"/dev/node_modules/a.ts":     "",
					"/dev/bower_components/a.ts": "",
					"/dev/jspm_packages/a.ts":    "",
				}, false)
			},
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/a.ts"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/a.ts", "/dev/x/a.ts"}) },
		},
		{
			name: "js wildcard excludes min js files",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/js/a.js":      "",
					"/dev/js/b.js":      "",
					"/dev/js/d.min.js":  "",
					"/dev/js/ab.min.js": "",
				}, false)
			},
			extensions: []string{".js"},
			includes:   []string{"js/*"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/js/a.js", "/dev/js/b.js"}) },
		},
		{
			name: "explicit min js pattern includes min files",
			host: func() vfs.FS {
				return vfstest.FromMap(map[string]string{
					"/dev/js/a.js":      "",
					"/dev/js/b.js":      "",
					"/dev/js/d.min.js":  "",
					"/dev/js/ab.min.js": "",
				}, false)
			},
			extensions: []string{".js"},
			includes:   []string{"js/*.min.js"},
			expect: func(t *testing.T, got []string) {
				expected := []string{"/dev/js/ab.min.js", "/dev/js/d.min.js"}
				assert.Equal(t, len(got), len(expected))
				for _, want := range expected {
					assert.Assert(t, slices.Contains(got, want))
				}
			},
		},
		{
			name:       "literal excludes baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"b.ts"},
			includes:   []string{"a.ts", "b.ts"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/a.ts"}) },
		},
		{
			name:       "wildcard excludes baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"*.ts", "z/??z.ts", "*/b.ts"},
			includes:   []string{"a.ts", "b.ts", "z/a.ts", "z/abz.ts", "z/aba.ts", "x/b.ts"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/z/a.ts", "/dev/z/aba.ts"}) },
		},
		{
			name:       "recursive excludes baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**/b.ts"},
			includes:   []string{"a.ts", "b.ts", "x/a.ts", "x/b.ts", "x/y/a.ts", "x/y/b.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/a.ts", "/dev/x/a.ts", "/dev/x/y/a.ts"})
			},
		},
		{
			name:       "question mark baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/?.ts"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/x/a.ts", "/dev/x/b.ts"}) },
		},
		{
			name:       "recursive directory pattern baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/a.ts"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/a.ts", "/dev/x/a.ts", "/dev/x/y/a.ts", "/dev/z/a.ts"})
			},
		},
		{
			name:       "case sensitive baseline",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/A.ts"},
			expect:     func(t *testing.T, got []string) { assert.DeepEqual(t, got, []string{"/dev/A.ts"}) },
		},
		{
			name:       "exclude folders baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"z", "x"},
			includes:   []string{"**/*"},
			expect: func(t *testing.T, got []string) {
				for _, f := range got {
					assert.Assert(t, !contains(f, "/z/") && !contains(f, "/x/"), "should not contain z or x: %s", f)
				}
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/b.ts"))
			},
		},
		{
			name:       "implicit glob expansion baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"z"},
			expect: func(t *testing.T, got []string) {
				assert.DeepEqual(t, got, []string{"/dev/z/a.ts", "/dev/z/aba.ts", "/dev/z/abz.ts", "/dev/z/b.ts", "/dev/z/bba.ts", "/dev/z/bbz.ts"})
			},
		},
		{
			name:       "trailing recursive directory baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**"},
			expect:     func(t *testing.T, got []string) { assert.Equal(t, len(got), 0) },
		},
		{
			name:       "exclude trailing recursive directory baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**"},
			includes:   []string{"**/*"},
			expect:     func(t *testing.T, got []string) { assert.Equal(t, len(got), 0) },
		},
		{
			name:       "multiple recursive directory patterns baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/x/**/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/aa.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/b.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/b.ts"))
			},
		},
		{
			name:       "include dirs with starstar prefix baseline",
			host:       caseSensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"**/x", "**/a/**/b"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/a.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/b.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/q/a/c/b/d.ts"))
			},
		},
		{
			name:       "dotted folders not implicitly included baseline",
			host:       dottedFoldersHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"x/**/*", "w/*/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/x/d.ts"))
				assert.Assert(t, slices.Contains(got, "/dev/x/y/d.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/.y/a.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/x/y/.e.ts"))
				assert.Assert(t, !slices.Contains(got, "/dev/w/.u/e.ts"))
			},
		},
		{
			name:       "include paths outside project baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			includes:   []string{"*", "/ext/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/dev/a.ts"))
				assert.Assert(t, slices.Contains(got, "/ext/ext.ts"))
			},
		},
		{
			name:       "include files with double dots baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"**"},
			includes:   []string{"/ext/b/a..b.ts"},
			expect:     func(t *testing.T, got []string) { assert.Assert(t, slices.Contains(got, "/ext/b/a..b.ts")) },
		},
		{
			name:       "exclude files with double dots baseline",
			host:       caseInsensitiveHost,
			extensions: []string{".ts", ".tsx", ".d.ts"},
			excludes:   []string{"/ext/b/a..b.ts"},
			includes:   []string{"/ext/**/*"},
			expect: func(t *testing.T, got []string) {
				assert.Assert(t, slices.Contains(got, "/ext/ext.ts"))
				assert.Assert(t, !slices.Contains(got, "/ext/b/a..b.ts"))
			},
		},
	}

	for _, tc := range cases {
		for _, impl := range readDirectoryImplementations {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				runReadDirectoryCase(t, tc, impl.fn)
			})
		}
	}
}

// TestSpecMatcher tests the SpecMatcher API
func TestSpecMatcher(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                      string
		specs                     []string
		basePath                  string
		usage                     Usage
		useCaseSensitiveFileNames bool
		matchingPaths             []string
		nonMatchingPaths          []string
	}{
		{
			name:                      "simple wildcard",
			specs:                     []string{"*.ts"},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/a.ts", "/project/b.ts", "/project/foo.ts"},
			nonMatchingPaths:          []string{"/project/a.js", "/project/sub/a.ts"},
		},
		{
			name:                      "recursive wildcard",
			specs:                     []string{"**/*.ts"},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/a.ts", "/project/sub/a.ts", "/project/sub/deep/a.ts"},
			nonMatchingPaths:          []string{"/project/a.js"},
		},
		{
			name:                      "exclude pattern",
			specs:                     []string{"node_modules"},
			basePath:                  "/project",
			usage:                     UsageExclude,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/node_modules", "/project/node_modules/foo"},
			nonMatchingPaths:          []string{"/project/src"},
		},
		{
			name:                      "case insensitive",
			specs:                     []string{"*.ts"},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: false,
			matchingPaths:             []string{"/project/A.TS", "/project/B.Ts"},
			nonMatchingPaths:          []string{"/project/a.js"},
		},
		{
			name:                      "multiple specs",
			specs:                     []string{"*.ts", "*.tsx"},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/a.ts", "/project/b.tsx"},
			nonMatchingPaths:          []string{"/project/a.js"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			matcher := NewSpecMatcher(tc.specs, tc.basePath, tc.usage, tc.useCaseSensitiveFileNames)
			if matcher == nil {
				t.Fatal("matcher should not be nil")
			}
			for _, path := range tc.matchingPaths {
				assert.Assert(t, matcher.MatchString(path), "should match: %s", path)
			}
			for _, path := range tc.nonMatchingPaths {
				assert.Assert(t, !matcher.MatchString(path), "should not match: %s", path)
			}
		})
	}
}

func TestSingleSpecMatcher(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                      string
		spec                      string
		basePath                  string
		usage                     Usage
		useCaseSensitiveFileNames bool
		expectNil                 bool
		matchingPaths             []string
		nonMatchingPaths          []string
	}{
		{
			name:                      "simple spec",
			spec:                      "*.ts",
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/a.ts"},
			nonMatchingPaths:          []string{"/project/a.js"},
		},
		{
			name:                      "trailing ** non-exclude returns nil",
			spec:                      "**",
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			expectNil:                 true,
		},
		{
			name:                      "trailing ** exclude works",
			spec:                      "**",
			basePath:                  "/project",
			usage:                     UsageExclude,
			useCaseSensitiveFileNames: true,
			matchingPaths:             []string{"/project/anything", "/project/deep/path"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			matcher := NewSingleSpecMatcher(tc.spec, tc.basePath, tc.usage, tc.useCaseSensitiveFileNames)
			if tc.expectNil {
				assert.Assert(t, matcher == nil, "should be nil")
				return
			}
			if matcher == nil {
				t.Fatal("matcher should not be nil")
			}
			for _, path := range tc.matchingPaths {
				assert.Assert(t, matcher.MatchString(path), "should match: %s", path)
			}
			for _, path := range tc.nonMatchingPaths {
				assert.Assert(t, !matcher.MatchString(path), "should not match: %s", path)
			}
		})
	}
}

func TestSpecMatchers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                      string
		specs                     []string
		basePath                  string
		usage                     Usage
		useCaseSensitiveFileNames bool
		expectNil                 bool
		pathToIndex               map[string]int
	}{
		{
			name:                      "multiple specs return correct index",
			specs:                     []string{"*.ts", "*.tsx", "*.js"},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			pathToIndex: map[string]int{
				"/project/a.ts":  0,
				"/project/b.tsx": 1,
				"/project/c.js":  2,
				"/project/d.css": -1, // no match
			},
		},
		{
			name:                      "empty specs returns nil",
			specs:                     []string{},
			basePath:                  "/project",
			usage:                     UsageFiles,
			useCaseSensitiveFileNames: true,
			expectNil:                 true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			matchers := NewSpecMatchers(tc.specs, tc.basePath, tc.usage, tc.useCaseSensitiveFileNames)
			if tc.expectNil {
				assert.Assert(t, matchers == nil, "should be nil")
				return
			}
			if matchers == nil {
				t.Fatal("matchers should not be nil")
			}
			for path, expectedIndex := range tc.pathToIndex {
				gotIndex := matchers.MatchIndex(path)
				assert.Equal(t, gotIndex, expectedIndex, "path: %s", path)
			}
		})
	}
}
