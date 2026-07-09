package nativepath

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

// F_GETPATH can return any hardlink sibling's name for an nlink > 1 file;
// Realpath must return the canonical path of the name it was asked about.
func TestRealpathHardlinkedFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	// expected values resolve /var -> /private/var etc. the same way
	// Realpath's result does
	canonical := func(path string) string {
		resolved, err := filepath.EvalSymlinks(path)
		assert.NilError(t, err)
		return resolved
	}

	// prime the kernel name cache for the inode under a specific name
	prime := func(path string) {
		f, err := os.Open(path)
		assert.NilError(t, err)
		assert.NilError(t, f.Close())
	}

	t.Run("returns the queried hardlink, not a cached sibling", func(t *testing.T) {
		t.Parallel()
		real := filepath.Join(tmp, "real.d.ts")
		alias := filepath.Join(tmp, "alias.d.ts")
		assert.NilError(t, os.WriteFile(real, []byte("export * from './lib';\n"), 0o666))
		assert.NilError(t, os.Link(real, alias))

		// v_name follows the latest lookup: the bug needs a concurrent
		// sibling open between Realpath's open() and F_GETPATH
		stop := make(chan struct{})
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				select {
				case <-stop:
					return
				default:
					prime(real)
				}
			}
		}()
		defer func() { close(stop); <-done }()

		// measured unpatched first-failure max is 15 iterations; 200 = >10x margin
		want := canonical(alias)
		for range 200 {
			got, err := Realpath(alias)
			assert.NilError(t, err)
			assert.Equal(t, got, want)
		}
	})

	t.Run("symlink to a hardlinked file resolves to the link target", func(t *testing.T) {
		t.Parallel()
		real := filepath.Join(tmp, "sym-real.d.ts")
		alias := filepath.Join(tmp, "sym-alias.d.ts")
		link := filepath.Join(tmp, "sym-link.d.ts")
		assert.NilError(t, os.WriteFile(real, []byte("export {};\n"), 0o666))
		assert.NilError(t, os.Link(real, alias))
		assert.NilError(t, os.Symlink(alias, link))

		prime(real)
		got, err := Realpath(link)
		assert.NilError(t, err)
		assert.Equal(t, got, canonical(alias))
	})

	t.Run("nlink == 1 file keeps the fast path result", func(t *testing.T) {
		t.Parallel()
		single := filepath.Join(tmp, "single.d.ts")
		assert.NilError(t, os.WriteFile(single, []byte("export {};\n"), 0o666))

		got, err := Realpath(single)
		assert.NilError(t, err)
		assert.Equal(t, got, canonical(single))
	})

	t.Run("directory keeps the fast path result", func(t *testing.T) {
		t.Parallel()
		dir := filepath.Join(tmp, "some-dir")
		assert.NilError(t, os.MkdirAll(dir, 0o777))

		got, err := Realpath(dir)
		assert.NilError(t, err)
		assert.Equal(t, got, canonical(dir))
	})
}
