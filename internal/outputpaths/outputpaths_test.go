package outputpaths_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/outputpaths"
	"gotest.tools/v3/assert"
)

func TestGetSourceFilePathInNewDirSourceMatchesCommonDirectory(t *testing.T) {
	t.Parallel()

	actual := outputpaths.GetSourceFilePathInNewDir("/project/src", "/project/out", "/project", "/project/src/", true)
	assert.Equal(t, actual, "/project/src")
}
