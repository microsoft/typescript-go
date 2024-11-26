package vfstest_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestInsensitive(t *testing.T) {
	t.Parallel()

	contents := []byte("bar")

	vfs := vfstest.ToMapFS(fstest.MapFS{
		"foo/bar/baz": &fstest.MapFile{
			Data: contents,
		},
	}, false /*useCaseSensitiveFileNames*/)

	sensitive, err := fs.ReadFile(vfs, "foo/bar/baz")
	assert.NilError(t, err)
	assert.DeepEqual(t, sensitive, contents)

	assert.NilError(t, fstest.TestFS(vfs, "foo/bar/baz"))

	insensitive, err := fs.ReadFile(vfs, "Foo/Bar/Baz")
	assert.NilError(t, err)
	assert.DeepEqual(t, insensitive, contents)

	assert.NilError(t, fstest.TestFS(vfs, "Foo/Bar/Baz"))
}

func TestSensitive(t *testing.T) {
	t.Parallel()

	contents := []byte("bar")

	vfs := vfstest.ToMapFS(fstest.MapFS{
		"foo/bar/baz": &fstest.MapFile{
			Data: contents,
		},
	}, true /*useCaseSensitiveFileNames*/)

	sensitive, err := fs.ReadFile(vfs, "foo/bar/baz")
	assert.NilError(t, err)
	assert.DeepEqual(t, sensitive, contents)

	_, err = fs.ReadFile(vfs, "Foo/Bar/Baz")
	assert.ErrorContains(t, err, "file does not exist")
}
