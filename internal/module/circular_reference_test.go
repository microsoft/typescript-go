package module_test

import (
	"io/fs"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// TestCircularModuleReference tests the scenario where a module imports itself
// through a relative path, which can cause a panic in the Pattern.Matches method
// when StarIndex is 0.
func TestCircularModuleReference(t *testing.T) {
	// Create a virtual file system with the problematic files
	fs := vfstest.FromMap(map[string]string{
		"/project/subdir/eslint.config.js": `export { default } from '../eslint.config.js'`,
		"/project/eslint.config.js":        `export default { rules: {} }`,
	}, false)

	// Create a host with the virtual file system
	host := newTestHost(fs, "/", false)

	// Create a resolver with default options
	resolver := module.NewResolver(host, &core.CompilerOptions{
		ModuleResolution: core.ModuleResolutionKindNode16,
	})
// Create a pattern that would have caused a panic before the fix
pattern := core.Pattern{
	Text:      "",
	StarIndex: 0,
}

// This should no longer panic after the fix
result := pattern.Matches("../eslint.config.js")

// Verify the result is as expected
// With our fix, this should return true if the candidate ends with the suffix
// which is everything after the star (p.Text[p.StarIndex+1:])
// In this case, the suffix is an empty string, so any candidate should match
if !result {
	t.Errorf("Expected pattern to match, but it didn't")
}

// Now try with the actual module resolution
resolvedModule := resolver.ResolveModuleName(
	"../eslint.config.js",
	"/project/subdir/file.js",
	core.ModuleKindCommonJS,
	nil,
)

// This should not panic now
t.Logf("Module resolution completed without panic: %v", resolvedModule)
}

// newTestHost creates a test host with the given file system
func newTestHost(fs vfs.FS, cwd string, useCaseSensitiveFileNames bool) *testHost {
	return &testHost{
		fs:                       fs,
		cwd:                      cwd,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
		traces:                   []string{},
	}
}

// testHost implements the module.Host interface for testing
type testHost struct {
	fs                       vfs.FS
	cwd                      string
	useCaseSensitiveFileNames bool
	traces                   []string
}

func (h *testHost) FS() vfs.FS {
	return h.fs
}

func (h *testHost) GetCurrentDirectory() string {
	return h.cwd
}

func (h *testHost) ReadFile(path string) ([]byte, error) {
	contents, ok := h.fs.ReadFile(path)
	if !ok {
		return nil, fs.ErrNotExist
	}
	return []byte(contents), nil
}

func (h *testHost) FileExists(path string) bool {
	return h.fs.FileExists(path)
}

func (h *testHost) DirectoryExists(path string) bool {
	return h.fs.DirectoryExists(path)
}

func (h *testHost) GetDirectories(path string) []string {
	// Simplified implementation for the test
	return []string{}
}

func (h *testHost) RealPath(path string) string {
	// Simplified implementation for the test
	return path
}

func (h *testHost) Trace(message string) {
	h.traces = append(h.traces, message)
}

func (h *testHost) GetEnvironmentVariable(name string) string {
	return ""
}

func (h *testHost) GetPathsBasedOnExtensions(extensions []string, path string) []string {
	// Simplified implementation for the test
	return []string{}
}