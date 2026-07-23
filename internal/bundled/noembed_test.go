//go:build noembed

package bundled

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestArgv0Path(t *testing.T) {
	t.Parallel()

	argv0 := argv0Path()
	assert.Assert(t, argv0 != "")
	assert.Assert(t, sameAsExecutable(argv0))
}

func TestSameAsExecutable(t *testing.T) {
	t.Parallel()

	exe, err := os.Executable()
	assert.NilError(t, err)
	assert.Assert(t, sameAsExecutable(exe))

	other := filepath.Join(t.TempDir(), "other")
	assert.NilError(t, os.WriteFile(other, []byte("not the executable"), 0o666))
	assert.Assert(t, !sameAsExecutable(other))
}
