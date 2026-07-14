package project

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestGetPathComponentsForWatching(t *testing.T) {
	t.Parallel()

	assert.DeepEqual(t, getPathComponentsForWatching("/project", ""), []string{"/", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\project", ""), []string{"C:/", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("//server/share/project/tsconfig.json", ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, getPathComponentsForWatching(`\\server\share\project\tsconfig.json`, ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\Users", ""), []string{"C:/Users"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\Users\\andrew\\project", ""), []string{"C:/Users/andrew", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("/home", ""), []string{"/home"})
	assert.DeepEqual(t, getPathComponentsForWatching("/home/andrew/project", ""), []string{"/home/andrew", "project"})
}

func TestNilWatchedFilesClone(t *testing.T) {
	t.Parallel()

	var w *WatchedFiles[int]
	result := w.Clone(42)
	assert.Assert(t, result == nil, "clone on a nil `WatchedFiles` should return nil")
}

func TestCreateResolutionLookupGlobMapperSkipsBundledLibs(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	data := &collections.SyncSet[tspath.Path]{}
	data.Add(tspath.Path("bundled:///libs/lib.d.ts"))

	mapper := createResolutionLookupGlobMapper("/workspace", bundled.LibPath(), "/workspace", true)
	result := mapper(data)

	assert.DeepEqual(t, result.patternsInsideWorkspace, []string(nil))
}

func TestCreateResolutionLookupGlobMapperWatchesRealLibDirectory(t *testing.T) {
	t.Parallel()

	data := &collections.SyncSet[tspath.Path]{}
	data.Add(tspath.Path("/ts/lib/lib.d.ts"))

	mapper := createResolutionLookupGlobMapper("/workspace", "/ts/lib", "/workspace", true)
	result := mapper(data)

	assert.DeepEqual(t, result.patternsInsideWorkspace, []string{"/ts/lib/**/*"})
}
