package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureOutput(f func()) (string, string) {
	oldStdout, oldStderr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout, os.Stderr = wOut, wErr

	f()

	// Close write ends before reading
	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr = oldStdout, oldStderr

	// Read output
	var outBuf, errBuf bytes.Buffer
	io.Copy(&outBuf, rOut)
	io.Copy(&errBuf, rErr)

	return outBuf.String(), errBuf.String()
}

func TestRunLSP_StderrClosed(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		closeStderr bool
		wantExit    int
		wantOutput  string
		wantErr     string
	}{
		{
			name:        "Stdio mode enabled (stderr open)",
			args:        []string{"-stdio"},
			closeStderr: false,
			wantExit:    0,
			wantOutput:  "",
			wantErr:     "",
		},
		{
			name:        "Stdio mode disabled (stderr open)",
			args:        []string{},
			closeStderr: false,
			wantExit:    1,
			wantOutput:  "",
			wantErr:     "stdio not supported\n",
		},
		{
			name:        "Stdio mode disabled (stderr closed)",
			args:        []string{},
			closeStderr: true,
			wantExit:    1,
			wantOutput:  "stderr unavailable, exiting with error stdio not supported\n",
			wantErr:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotExit int
			var out, err string

			if tt.closeStderr {
				// Simulate closed stderr
				oldStderr := os.Stderr
				_, wErr, _ := os.Pipe()
				os.Stderr = wErr
				wErr.Close() // Close stderr to trigger a write error

				out, err = captureOutput(func() {
					gotExit = runLSP(tt.args)
				})

				// Restore os.Stderr after the test
				os.Stderr = oldStderr
			} else {
				out, err = captureOutput(func() {
					gotExit = runLSP(tt.args)
				})
			}

			// Check exit code
			if gotExit != tt.wantExit {
				t.Errorf("runLSP(%v) exit code = %d; want %d", tt.args, gotExit, tt.wantExit)
			}

			// Check stdout output
			if out != tt.wantOutput {
				t.Errorf("runLSP(%v) stdout = %q; want %q", tt.args, out, tt.wantOutput)
			}

			// Check stderr output
			if err != tt.wantErr {
				t.Errorf("runLSP(%v) stderr = %q; want %q", tt.args, err, tt.wantErr)
			}
		})
	}
}
