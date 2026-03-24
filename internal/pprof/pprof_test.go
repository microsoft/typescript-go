package pprof

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil"
	"gotest.tools/v3/assert"
)

func TestBeginProfilingCleansUpCPUFileWhenStartFails(t *testing.T) {
	wantErr := errors.New("failed to start profile")
	originalStartCPUProfile := startCPUProfile
	startCPUProfile = func(io.Writer) error {
		return wantErr
	}
	t.Cleanup(func() {
		startCPUProfile = originalStartCPUProfile
	})

	profileDir := t.TempDir()
	cpuProfilePath := filepath.Join(profileDir, fmt.Sprintf("%d-cpuprofile.pb.gz", os.Getpid()))

	testutil.AssertPanics(t, func() {
		BeginProfiling(profileDir, io.Discard)
	}, wantErr)

	_, err := os.Stat(cpuProfilePath)
	assert.Assert(t, os.IsNotExist(err), "expected CPU profile file to be removed after start failure, got %v", err)
}
