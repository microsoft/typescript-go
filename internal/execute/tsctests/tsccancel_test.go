package tsctests

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
)

// TestTscNoEmitCancellation verifies that a canceled context (wired to SIGINT in
// cmd/tsgo/main.go) aborts a compile instead of running it to completion. A
// pre-canceled context deterministically hits the same checker cancellation polling
// a mid-flight SIGINT would.
func TestTscNoEmitCancellation(t *testing.T) {
	t.Parallel()

	// A type error lets us tell whether the checker ran to completion: if it did,
	// the diagnostic is reported; if aborted, it is not.
	const badSource = `const x: number = "not a number";`

	testCases := []struct {
		name  string
		args  []string
		files FileMap
	}{
		{
			// Top-level short-circuit in EmitFilesAndReportErrors.
			name: "single file",
			args: []string{"--noEmit"},
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
				"/home/src/workspaces/project/main.ts":       badSource,
			},
		},
		{
			// --singleThreaded funnels all files through one checker, so
			// forEachCheckerGroupDo reuses it across files. Pins the guard that stops
			// reuse after cancellation (else checkNotCanceled panics on the 2nd file).
			name: "multi file single checker",
			args: []string{"--noEmit", "--singleThreaded"},
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
				"/home/src/workspaces/project/a.ts":          badSource,
				"/home/src/workspaces/project/b.ts":          badSource,
				"/home/src/workspaces/project/c.ts":          badSource,
			},
		},
		{
			// `tsc -b`: exercises the context threaded through the build orchestrator.
			name: "build mode",
			args: []string{"-b"},
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "composite": true, "strict": true } }`,
				"/home/src/workspaces/project/main.ts":       badSource,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sys := newTestSys(&tscInput{
				commandLineArgs: tc.args,
				files:           tc.files,
			}, false)

			ctx, cancel := context.WithCancel(context.Background())
			cancel() // simulate SIGINT delivered before/at the start of the compile

			result := execute.CommandLine(ctx, sys, tc.args, sys)

			// Aborts with a distinct status (and must not panic reusing a canceled checker).
			if result.Status != tsc.ExitStatusCanceled {
				t.Errorf("status = %v, want ExitStatusCanceled (compile should abort on cancellation)", result.Status)
			}

			// Aborted checks must not report their incomplete diagnostics.
			if out := sys.getOutput(true); strings.Contains(out, "error TS") {
				t.Errorf("expected no diagnostics to be reported after cancellation; got output:\n%s", out)
			}
		})
	}
}
