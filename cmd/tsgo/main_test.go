package main

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"gotest.tools/v3/assert"
)

func TestExitCode(t *testing.T) {
	t.Parallel()

	assert.Equal(t, exitCode(tsc.ExitStatusSuccess), 0)
	assert.Equal(t, exitCode(tsc.ExitStatusDiagnosticsPresent_OutputsSkipped), 2)
	assert.Equal(t, exitCode(tsc.ExitStatusDiagnosticsPresent_OutputsGenerated), 2)
	assert.Equal(t, exitCode(tsc.ExitStatusInvalidProject_OutputsSkipped), 3)
	assert.Equal(t, exitCode(tsc.ExitStatusProjectReferenceCycle_OutputsSkipped), 4)
	assert.Equal(t, exitCode(tsc.ExitStatusNotImplemented), 5)
}
