package project

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestGetPathComponentsForWatching(t *testing.T) {
	t.Parallel()

	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("/project", ""), []string{"/", "project"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("C:\\project", ""), []string{"C:/", "project"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("//server/share/project/tsconfig.json", ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching(`\\server\share\project\tsconfig.json`, ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("C:\\Users", ""), []string{"C:/Users"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("C:\\Users\\andrew\\project", ""), []string{"C:/Users/andrew", "project"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("/home", ""), []string{"/home"})
	assert.DeepEqual(t, tspath.GetPathComponentsForWatching("/home/andrew/project", ""), []string{"/home/andrew", "project"})
}

func TestNilWatchedFilesClone(t *testing.T) {
	t.Parallel()

	var w *WatchedFiles[int]
	result := w.Clone(42)
	assert.Assert(t, result == nil, "clone on a nil `WatchedFiles` should return nil")
}
