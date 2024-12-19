package vfs

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestSymlinkRealpath(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	target := filepath.Join(tmp, "target")
	targetFile := filepath.Join(target, "file")

	link := filepath.Join(tmp, "link")
	linkFile := filepath.Join(link, "file")

	const expectedContents = "hello"

	assert.NilError(t, os.MkdirAll(target, 0o777))
	assert.NilError(t, os.WriteFile(targetFile, []byte(expectedContents), 0o666))

	if runtime.GOOS == "windows" {
		// Don't use os.Symlink on Windows, as it creates a "real" symlink, not a junction.
		assert.NilError(t, exec.Command("cmd", "/c", "mklink", "/J", link, target).Run())
	} else {
		assert.NilError(t, os.Symlink(target, link))
	}

	gotContents, err := os.ReadFile(linkFile)
	assert.NilError(t, err)
	assert.Equal(t, string(gotContents), expectedContents)

	fs := FromOS()

	targetRealpath := tspath.NormalizePath(targetFile)
	if runtime.GOOS == "darwin" {
		// macOS makes its temporary directory a symlink itself.
		targetRealpath = fs.Realpath(targetRealpath)
	}
	linkRealpath := fs.Realpath(tspath.NormalizePath(linkFile))
	assert.Equal(t, targetRealpath, linkRealpath)
}
