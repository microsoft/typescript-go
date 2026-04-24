package testrunner

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCompilerVaryByIncludesStrictArrayVariance(t *testing.T) {
	t.Parallel()
	_, ok := compilerVaryBy["strictarrayvariance"]
	assert.Assert(t, ok, "strictArrayVariance must stay in compilerVaryBy so // @strictArrayVariance: ... works in compiler baselines")
}
