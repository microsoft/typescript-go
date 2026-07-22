package nativepath

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func mustRealpath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := Realpath(path)
	assert.NilError(t, err)
	return resolved
}

func TestRealpath(t *testing.T) {
	t.Parallel()

	// Canonicalize the temp dir itself so expectations are not polluted by
	// symlinks higher up in the temp path.
	tmp := mustRealpath(t, t.TempDir())

	t.Run("regular file resolves to itself", func(t *testing.T) {
		t.Parallel()
		file := filepath.Join(tmp, "regular.txt")
		assert.NilError(t, os.WriteFile(file, []byte("hello"), 0o666))
		assert.Equal(t, mustRealpath(t, file), file)
	})

	t.Run("directory symlink resolves, file name preserved", func(t *testing.T) {
		t.Parallel()
		// Mirrors the pnpm layout: the directory components are symlinked,
		// the final component keeps its name.
		target := filepath.Join(tmp, "store", "pkg")
		assert.NilError(t, os.MkdirAll(target, 0o777))
		file := filepath.Join(target, "index.d.ts")
		assert.NilError(t, os.WriteFile(file, []byte("export {};"), 0o666))
		link := filepath.Join(tmp, "linked-pkg")
		assert.NilError(t, os.Symlink(target, link))

		assert.Equal(t, mustRealpath(t, filepath.Join(link, "index.d.ts")), file)
	})

	t.Run("symlinked final component resolves to its target", func(t *testing.T) {
		t.Parallel()
		// A genuine symlink may rename the final component; the kernel's
		// answer must be kept.
		target := filepath.Join(tmp, "real-name.txt")
		assert.NilError(t, os.WriteFile(target, []byte("hello"), 0o666))
		link := filepath.Join(tmp, "other-name.txt")
		assert.NilError(t, os.Symlink(target, link))

		assert.Equal(t, mustRealpath(t, link), target)
	})

	t.Run("hardlink keeps the name it was opened by", func(t *testing.T) {
		t.Parallel()
		// Hardlinks are equal names for one inode: canonicalization must not
		// swap one for the other. This is the invariant the masked-filesystem
		// fallback relies on (see the package comment).
		first := filepath.Join(tmp, "first-name.txt")
		assert.NilError(t, os.WriteFile(first, []byte("hello"), 0o666))
		second := filepath.Join(tmp, "second-name.txt")
		assert.NilError(t, os.Link(first, second))

		assert.Equal(t, mustRealpath(t, second), second)
	})

	t.Run("missing file errors", func(t *testing.T) {
		t.Parallel()
		_, err := Realpath(filepath.Join(tmp, "does-not-exist.txt"))
		assert.Assert(t, err != nil)
	})
}
