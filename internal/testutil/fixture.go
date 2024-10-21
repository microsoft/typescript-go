package testutil

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

type Fixture struct {
	path     string
	contents func() (string, error)
}

func NewFixture(elem ...string) *Fixture {
	p := filepath.Clean(filepath.Join(elem...))
	return &Fixture{
		path: p,
		// Cache the file contents and errors.
		contents: sync.OnceValues(func() (string, error) {
			b, err := os.ReadFile(p)
			return string(b), err
		}),
	}
}

func (f *Fixture) Path() string {
	return f.path
}

func (f *Fixture) SkipIfNotExist(t testing.TB) {
	t.Helper()

	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		t.Skipf("Test fixture %q does not exist", f.path)
	}
}

func (f *Fixture) ReadFile(t testing.TB) string {
	t.Helper()

	contents, err := f.contents()
	if err != nil {
		t.Fatalf("Failed to read test fixture %q: %v", f.path, err)
	}
	return contents
}
