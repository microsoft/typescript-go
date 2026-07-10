package tsctests

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
)

// TestTscNoEmitCancellation verifies that interrupting a `--noEmit` compile via
// the context passed to execute.CommandLine (which the CLI wires to SIGINT/SIGTERM
// in cmd/tsgo/main.go) aborts the compile promptly rather than running it to
// completion.
//
// A pre-canceled context deterministically exercises the same isCanceled() polling
// the checker uses for a mid-flight SIGINT (see internal/checker/utilities.go), so
// it covers both the "already canceled" and "canceled during check" cases.
func TestTscNoEmitCancellation(t *testing.T) {
	t.Parallel()

	// Each file contains a type error so we can prove whether the checker ran to
	// completion: if it did, the diagnostic is reported; if the compile was
	// abandoned on cancellation, it is not.
	const badSource = `const x: number = "not a number";`

	testCases := []struct {
		name  string
		args  []string
		files FileMap
	}{
		{
			// Single file: exercises the top-level cancellation short-circuit in
			// EmitFilesAndReportErrors / GetDiagnosticsOfAnyProgram.
			name: "single file",
			args: []string{"--noEmit"},
			files: FileMap{
				"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
				"/home/src/workspaces/project/main.ts":       badSource,
			},
		},
		{
			// Multiple files under --singleThreaded funnel through a single checker,
			// so the per-file loop in checkerPool.forEachCheckerGroupDo reuses the
			// same checker across files. Once the first file cancels it, reusing it
			// for the next file would panic in checkNotCanceled without the guard
			// there. This case pins that guard.
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
			// Build mode (`tsc -b`): exercises the context threaded through the build
			// orchestrator (Start -> buildOrClean -> rangeTask -> buildProject ->
			// compileAndEmit). A canceled build must abort rather than run the project's
			// compile to completion.
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

			// The compile should short-circuit with a distinct canceled status
			// instead of running the checker to completion (and it must not panic
			// reusing a canceled checker).
			if result.Status != tsc.ExitStatusCanceled {
				t.Errorf("status = %v, want ExitStatusCanceled (compile should abort on cancellation)", result.Status)
			}

			// Because the check was abandoned, its (incomplete) diagnostics must not
			// be reported: the type errors should be absent from the output.
			if out := sys.getOutput(true); strings.Contains(out, "error TS") {
				t.Errorf("expected no diagnostics to be reported after cancellation; got output:\n%s", out)
			}
		})
	}
}
