package vfs_test

import (
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/vfs"
	"gotest.tools/v3/assert"
)

func TestFS(t *testing.T) {
	t.Parallel()

	testfs := fstest.MapFS{
		"foo": &fstest.MapFile{
			Data: []byte("hello, world"),
		},
	}

	fs := vfs.FromIOFS("/", true, testfs)

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		p := fs.ToPath("/foo")

		content, ok := fs.ReadFile(p)
		assert.Assert(t, ok)
		assert.Equal(t, content, "hello, world")
	})
}
