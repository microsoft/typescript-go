package testutil

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

type FileFixture interface {
	Name() string
	Path() string
	SkipIfNotExist(t testing.TB)
	ReadFile(t testing.TB) string
}

type fixtureFromFile struct {
	name     string
	path     string
	contents func() (string, error)
}

func NewFileFixtureFromFile(name string, pathParts []string) FileFixture {
	p := filepath.Clean(filepath.Join(pathParts...))
	return &fixtureFromFile{
		name: name,
		path: p,
		// Cache the file contents and errors.
		contents: sync.OnceValues(func() (string, error) {
			b, err := os.ReadFile(p)
			return string(b), err
		}),
	}
}

func (f *fixtureFromFile) Name() string { return f.name }
func (f *fixtureFromFile) Path() string { return f.path }

func (f *fixtureFromFile) SkipIfNotExist(t testing.TB) {
	t.Helper()

	if _, err := os.Stat(f.path); err != nil {
		t.Skipf("Test fixture %q does not exist", f.path)
	}
}

func (f *fixtureFromFile) ReadFile(t testing.TB) string {
	t.Helper()

	contents, err := f.contents()
	if err != nil {
		t.Fatalf("Failed to read test fixture %q: %v", f.path, err)
	}
	return contents
}

type fixtureFromString struct {
	name     string
	path     string
	contents string
}

func NewFileFixtureFromString(name string, path string, contents string) FileFixture {
	return &fixtureFromString{
		name:     name,
		path:     path,
		contents: contents,
	}
}

func (f *fixtureFromString) Name() string { return f.name }
func (f *fixtureFromString) Path() string { return f.path }

func (f *fixtureFromString) SkipIfNotExist(t testing.TB) {}

func (f *fixtureFromString) ReadFile(t testing.TB) string { return f.contents }
