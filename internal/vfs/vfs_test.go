package vfs_test

import (
	"encoding/binary"
	"testing"
	"testing/fstest"
	"unicode/utf16"

	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestVFSTestMapFS(t *testing.T) {
	t.Parallel()

	fs := vfstest.FromMapFS(fstest.MapFS{
		"foo.ts": &fstest.MapFile{
			Data: []byte("hello, world"),
		},
		"dir1/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"dir1/file2.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"dir2/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
	}, false /*useCaseSensitiveFileNames*/)

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		content, ok := fs.ReadFile("/foo.ts")
		assert.Assert(t, ok)
		assert.Equal(t, content, "hello, world")

		content, ok = fs.ReadFile("/does/not/exist.ts")
		assert.Assert(t, !ok)
		assert.Equal(t, content, "")
	})

	t.Run("Realpath", func(t *testing.T) {
		t.Parallel()

		realpath := fs.Realpath("/foo.ts")
		assert.Equal(t, realpath, "/foo.ts")

		realpath = fs.Realpath("/Foo.ts")
		assert.Equal(t, realpath, "/foo.ts")

		realpath = fs.Realpath("/does/not/exist.ts")
		assert.Equal(t, realpath, "/does/not/exist.ts")
	})

	t.Run("UseCaseSensitiveFileNames", func(t *testing.T) {
		t.Parallel()

		assert.Assert(t, !fs.UseCaseSensitiveFileNames())
	})
}

func TestVFSTestMapFSWindows(t *testing.T) {
	t.Parallel()

	fs := vfstest.FromMapFS(fstest.MapFS{
		"c:/foo.ts": &fstest.MapFile{
			Data: []byte("hello, world"),
		},
		"c:/dir1/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"c:/dir1/file2.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"c:/dir2/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
	}, false)

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		content, ok := fs.ReadFile("c:/foo.ts")
		assert.Assert(t, ok)
		assert.Equal(t, content, "hello, world")

		content, ok = fs.ReadFile("c:/does/not/exist.ts")
		assert.Assert(t, !ok)
		assert.Equal(t, content, "")
	})

	t.Run("Realpath", func(t *testing.T) {
		t.Parallel()

		realpath := fs.Realpath("c:/foo.ts")
		assert.Equal(t, realpath, "c:/foo.ts")

		realpath = fs.Realpath("c:/Foo.ts")
		assert.Equal(t, realpath, "c:/foo.ts")

		realpath = fs.Realpath("c:/does/not/exist.ts")
		assert.Equal(t, realpath, "c:/does/not/exist.ts")
	})
}

func TestBOM(t *testing.T) {
	t.Parallel()

	const expected = "hello, world"

	tests := []struct {
		name  string
		order binary.ByteOrder
		bom   [2]byte
	}{
		{"BigEndian", binary.BigEndian, [2]byte{0xFE, 0xFF}},
		{"LittleEndian", binary.LittleEndian, [2]byte{0xFF, 0xFE}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var codePoints []uint16

			for _, r := range expected {
				codePoints = utf16.AppendRune(codePoints, r)
			}

			buf := tt.bom[:]

			for _, r := range codePoints {
				var err error
				buf, err = binary.Append(buf, tt.order, r)
				assert.NilError(t, err)
			}

			testfs := fstest.MapFS{
				"foo.ts": &fstest.MapFile{
					Data: buf,
				},
			}

			fs := vfs.FromIOFS(testfs, true)

			content, ok := fs.ReadFile("/foo.ts")
			assert.Assert(t, ok)
			assert.Equal(t, content, expected)
		})
	}

	t.Run("UTF8", func(t *testing.T) {
		t.Parallel()

		testfs := fstest.MapFS{
			"foo.ts": &fstest.MapFile{
				Data: []byte("\xEF\xBB\xBF" + expected),
			},
		}

		fs := vfs.FromIOFS(testfs, true)

		content, ok := fs.ReadFile("/foo.ts")
		assert.Assert(t, ok)
		assert.Equal(t, content, expected)
	})
}

func BenchmarkReadFile(b *testing.B) {
	type bench struct {
		name string
		fs   vfs.FS
		path string
	}

	osFS := vfs.FromOS()

	const smallData = "hello, world"
	tmpdir := tspath.NormalizeSlashes(b.TempDir())
	osSmallDataPath := tspath.CombinePaths(tmpdir, "foo.ts")
	err := osFS.WriteFile(osSmallDataPath, smallData, false)
	assert.NilError(b, err)

	tests := []bench{
		{"MapFS small", vfstest.FromMapFS(fstest.MapFS{
			"foo.ts": &fstest.MapFile{
				Data: []byte(smallData),
			},
		}, true), "/foo.ts"},
		{"OS small", osFS, osSmallDataPath},
	}

	if repo.TypeScriptSubmoduleExists() {
		checkerPath := tspath.CombinePaths(tspath.NormalizeSlashes(repo.TypeScriptSubmodulePath), "src", "compiler", "checker.ts")

		checkerContents, ok := osFS.ReadFile(checkerPath)
		assert.Assert(b, ok)

		tests = append(tests, bench{
			"MapFS checker.ts",
			vfstest.FromMapFS(fstest.MapFS{
				"checker.ts": &fstest.MapFile{
					Data: []byte(checkerContents),
				},
			}, true),
			"/checker.ts",
		})
		tests = append(tests, bench{"OS checker.ts", osFS, checkerPath})
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				_, _ = tt.fs.ReadFile(tt.path)
			}
		})
	}
}
