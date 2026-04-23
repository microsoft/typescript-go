package tspath

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetPathComponentsForWatching(t *testing.T) {
	t.Parallel()

	assert.DeepEqual(t, GetPathComponentsForWatching("/project", ""), []string{"/", "project"})
	assert.DeepEqual(t, GetPathComponentsForWatching("C:\\project", ""), []string{"C:/", "project"})
	assert.DeepEqual(t, GetPathComponentsForWatching("//server/share/project/tsconfig.json", ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, GetPathComponentsForWatching(`\\server\share\project\tsconfig.json`, ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, GetPathComponentsForWatching("C:\\Users", ""), []string{"C:/Users"})
	assert.DeepEqual(t, GetPathComponentsForWatching("C:\\Users\\andrew\\project", ""), []string{"C:/Users/andrew", "project"})
	assert.DeepEqual(t, GetPathComponentsForWatching("/home", ""), []string{"/home"})
	assert.DeepEqual(t, GetPathComponentsForWatching("/home/andrew/project", ""), []string{"/home/andrew", "project"})
}
