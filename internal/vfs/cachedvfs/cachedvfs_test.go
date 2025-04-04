package cachedvfs_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfsmock"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func createMockFS() *vfsmock.FSMock {
	return vfsmock.Wrap(vfstest.FromMap(map[string]string{
		"/some/path/file.txt": "hello world",
	}, true))
}

func TestDirectoryExists(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.DirectoryExists("/some/path")
	assert.Equal(t, 1, len(underlying.DirectoryExistsCalls()))

	cached.DirectoryExists("/some/path")
	assert.Equal(t, 1, len(underlying.DirectoryExistsCalls()), "DirectoryExists called underlying FS more than once")

	cached.DirectoryExists("/other/path")
	assert.Equal(t, 2, len(underlying.DirectoryExistsCalls()))
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.FileExists("/some/path/file.txt")
	assert.Equal(t, 1, len(underlying.FileExistsCalls()))

	cached.FileExists("/some/path/file.txt")
	assert.Equal(t, 1, len(underlying.FileExistsCalls()), "FileExists called underlying FS more than once")

	cached.FileExists("/other/path/file.txt")
	assert.Equal(t, 2, len(underlying.FileExistsCalls()))
}

func TestGetAccessibleEntries(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.GetAccessibleEntries("/some/path")
	assert.Equal(t, 1, len(underlying.GetAccessibleEntriesCalls()))

	cached.GetAccessibleEntries("/some/path")
	assert.Equal(t, 1, len(underlying.GetAccessibleEntriesCalls()), "GetAccessibleEntries called underlying FS more than once")

	cached.GetAccessibleEntries("/other/path")
	assert.Equal(t, 2, len(underlying.GetAccessibleEntriesCalls()))
}

func TestRealpath(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.Realpath("/some/path")
	assert.Equal(t, 1, len(underlying.RealpathCalls()))

	cached.Realpath("/some/path")
	assert.Equal(t, 1, len(underlying.RealpathCalls()), "Realpath called underlying FS more than once")

	cached.Realpath("/other/path")
	assert.Equal(t, 2, len(underlying.RealpathCalls()))
}

func TestStat(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.Stat("/some/path")
	assert.Equal(t, 1, len(underlying.StatCalls()))

	cached.Stat("/some/path")
	assert.Equal(t, 1, len(underlying.StatCalls()), "Stat called underlying FS more than once")

	cached.Stat("/other/path")
	assert.Equal(t, 2, len(underlying.StatCalls()))
}

func TestReadFileNotCached(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.ReadFile("/some/path/file.txt")
	assert.Equal(t, 1, len(underlying.ReadFileCalls()))

	cached.ReadFile("/some/path/file.txt")
	assert.Equal(t, 2, len(underlying.ReadFileCalls()), "ReadFile should not be cached")
}

func TestUseCaseSensitiveFileNamesNotCached(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.UseCaseSensitiveFileNames()
	assert.Equal(t, 1, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.UseCaseSensitiveFileNames()
	assert.Equal(t, 2, len(underlying.UseCaseSensitiveFileNamesCalls()), "UseCaseSensitiveFileNames should not be cached")
}

func TestWalkDirNotCached(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	walkFn := vfs.WalkDirFunc(func(path string, info vfs.DirEntry, err error) error {
		return nil
	})

	cached.WalkDir("/some/path", walkFn)
	assert.Equal(t, 1, len(underlying.WalkDirCalls()))

	cached.WalkDir("/some/path", walkFn)
	assert.Equal(t, 2, len(underlying.WalkDirCalls()), "WalkDir should not be cached")
}
