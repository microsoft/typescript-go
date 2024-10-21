package compiler

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	// Test binaries always start in the package directory; grab this early.
	cwd = must(os.Getwd())

	srcCompilerCheckerTS = newTestFixture(cwd, "../../_submodules/TypeScript/src/compiler/checker.ts")
)

type testFixture struct {
	path     string
	contents func() (string, error)
}

func newTestFixture(elem ...string) *testFixture {
	p := filepath.Clean(filepath.Join(elem...))
	return &testFixture{
		path: p,
		// Cache the file contents and errors.
		contents: sync.OnceValues(func() (string, error) {
			b, err := os.ReadFile(p)
			return string(b), err
		}),
	}
}

func (f *testFixture) Path() string {
	return f.path
}

func (f *testFixture) SkipIfNotExist(t testing.TB) {
	t.Helper()

	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		t.Skipf("Test fixture %q does not exist", f.path)
	}
}

func (f *testFixture) ReadFile(t testing.TB) string {
	t.Helper()

	contents, err := f.contents()
	if err != nil {
		t.Fatalf("Failed to read test fixture %q: %v", f.path, err)
	}
	return contents
}
