package tsctests

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
)

// cancelAfterNPolls is a context that reports itself canceled only after Err has
// been polled pollThreshold times while still uncanceled. The checker polls
// ctx.Err() between top-level statements (checkSourceElements), so a small threshold
// lands the cancellation *after* checking has begun rather than before it starts --
// the case a pre-canceled context cannot exercise. Once tripped it stays canceled.
type cancelAfterNPolls struct {
	context.Context
	pollThreshold int32
	polls         atomic.Int32
	tripped       atomic.Bool
	done          chan struct{}
}

func newCancelAfterNPolls(pollThreshold int32) *cancelAfterNPolls {
	return &cancelAfterNPolls{
		Context:       context.Background(),
		pollThreshold: pollThreshold,
		done:          make(chan struct{}),
	}
}

func (c *cancelAfterNPolls) Err() error {
	if c.tripped.Load() {
		return context.Canceled
	}
	if c.polls.Add(1) > c.pollThreshold {
		if c.tripped.CompareAndSwap(false, true) {
			close(c.done)
		}
		return context.Canceled
	}
	return nil
}

func (c *cancelAfterNPolls) Done() <-chan struct{} { return c.done }

// TestTscCancellationAborts verifies that a canceled compile aborts with
// ExitStatusCanceled and never reports its (incomplete) diagnostics -- both when the
// signal arrives before the compile starts and when it lands mid-check, where a
// checker actually runs and is marked canceled. The mid-check cases are what would
// panic in checkNotCanceled if the reuse guards were missing (a canceled checker fed
// more files, or asked for global/declaration diagnostics again).
func TestTscCancellationAborts(t *testing.T) {
	t.Parallel()

	// Many statements per file so a mid-check cancellation reliably lands while
	// checking, across more than one file. The type errors let us tell whether the
	// checker ran to completion: if it did the diagnostics are reported, if aborted
	// they are not. Distinct names per statement keep the checker busy.
	var bad, inferred strings.Builder
	for i := range 50 {
		x := strings.Repeat("x", i+1)
		bad.WriteString("export const v")
		bad.WriteString(x)
		bad.WriteString(`: number = "not a number";` + "\n")
		// Inferred return types force type serialization during declaration emit
		// (SerializeReturnTypeForSignature -> node reuse -> checkNotCanceled), a
		// distinct checker-reuse path from plain semantic checking.
		inferred.WriteString("export function make")
		inferred.WriteString(x)
		inferred.WriteString("() { return { a: 1, b: 'x', deep: [1, 2, 3] as const }; }\n")
	}
	badSrc, inferredSrc := bad.String(), inferred.String()

	testCases := []struct {
		name     string
		args     []string
		tsconfig string
		src      string
		midCheck bool // cancel during checking rather than before the compile starts
	}{
		{
			name:     "pre-canceled single file",
			args:     []string{"--noEmit"},
			tsconfig: `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
			src:      badSrc,
		},
		{
			name:     "pre-canceled build mode",
			args:     []string{"-b"},
			tsconfig: `{ "compilerOptions": { "composite": true, "strict": true } }`,
			src:      badSrc,
		},
		{
			// --singleThreaded funnels all files through one checker, so after it
			// cancels on an early file forEachCheckerGroupDo must stop feeding it later
			// files.
			name:     "mid-check single checker",
			args:     []string{"--noEmit", "--singleThreaded"},
			tsconfig: `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
			src:      badSrc,
			midCheck: true,
		},
		{
			// tsc -b through the incremental program + orchestrator, which also reaches
			// GetGlobalDiagnostics from emitBuildInfo -> ensureHasErrorsForState.
			name:     "mid-check build mode",
			args:     []string{"-b"},
			tsconfig: `{ "compilerOptions": { "composite": true, "strict": true } }`,
			src:      badSrc,
			midCheck: true,
		},
		{
			// A canceled checker must not be handed to GetDeclarationDiagnostics.
			name:     "mid-check declaration emit",
			args:     []string{"--noEmit", "--declaration", "--singleThreaded"},
			tsconfig: `{ "compilerOptions": { "noEmit": true, "declaration": true, "strict": true } }`,
			src:      inferredSrc,
			midCheck: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sys := newTestSys(&tscInput{
				commandLineArgs: tc.args,
				files: FileMap{
					"/home/src/workspaces/project/tsconfig.json": tc.tsconfig,
					"/home/src/workspaces/project/a.ts":          tc.src,
					"/home/src/workspaces/project/b.ts":          tc.src,
					"/home/src/workspaces/project/c.ts":          tc.src,
				},
			}, false)

			result, midChecked := runWithCancellation(t, sys, tc.args, tc.midCheck)

			// Aborts with a distinct status (and must not panic reusing a canceled checker).
			if result.Status != tsc.ExitStatusCanceled {
				t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
			}
			// Aborted checks must not report their incomplete diagnostics.
			if out := sys.getOutput(true); strings.Contains(out, "error TS") {
				t.Errorf("expected no diagnostics after cancellation; got output:\n%s", out)
			}
			// A mid-check case that never tripped during checking isn't testing what it claims.
			if tc.midCheck && !midChecked {
				t.Error("expected cancellation to trip during checking, but it never did")
			}
		})
	}
}

// runWithCancellation runs the command line under a canceled context and returns the
// result. When midCheck is false the context is canceled before the run starts; when
// true it is canceled after checking has begun, and the returned bool reports whether
// that mid-check cancellation actually tripped. The run is guarded by a timeout so a
// regression that ignores cancellation fails loudly instead of hanging.
func runWithCancellation(t *testing.T, sys *TestSys, args []string, midCheck bool) (tsc.CommandLineResult, bool) {
	t.Helper()
	var (
		ctx        context.Context
		midChecked func() bool
	)
	if midCheck {
		// Threshold small relative to the statement count so cancellation lands well
		// before checking finishes.
		c := newCancelAfterNPolls(5)
		ctx, midChecked = c, c.tripped.Load
	} else {
		canceled, cancel := context.WithCancel(context.Background())
		cancel()
		ctx, midChecked = canceled, func() bool { return false }
	}

	resultCh := make(chan tsc.CommandLineResult, 1)
	go func() {
		resultCh <- execute.CommandLine(ctx, sys, args, sys)
	}()
	select {
	case result := <-resultCh:
		return result, midChecked()
	case <-time.After(30 * time.Second):
		t.Fatal("compile did not abort after cancellation")
		return tsc.CommandLineResult{}, false
	}
}

// runPreCanceled runs the command line on an existing sys under a context that is
// already canceled before the run starts, guarded by a timeout so a regression that
// ignores cancellation fails loudly instead of hanging.
func runPreCanceled(t *testing.T, sys *TestSys, args []string) tsc.CommandLineResult {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resultCh := make(chan tsc.CommandLineResult, 1)
	go func() {
		resultCh <- execute.CommandLine(ctx, sys, args, sys)
	}()
	select {
	case result := <-resultCh:
		return result
	case <-time.After(30 * time.Second):
		t.Fatal("run did not abort after cancellation")
		return tsc.CommandLineResult{}
	}
}

// TestTscBuildCancellationUpToDate verifies that an interrupt is honored even when a
// `tsc -b` build has nothing to do. A root project that is already up to date has no
// upstream to wait on, so the no-build path must still observe a pre-canceled context
// and report ExitStatusCanceled rather than swallowing the interrupt as success.
func TestTscBuildCancellationUpToDate(t *testing.T) {
	t.Parallel()
	files := FileMap{
		"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "composite": true } }`,
		"/home/src/workspaces/project/a.ts":          "export const a = 1;\n",
	}
	sys := newTestSys(&tscInput{
		commandLineArgs: []string{"-b"},
		files:           files,
	}, false)

	// First build to success so the project is up to date on the next run.
	if result := execute.CommandLine(context.Background(), sys, []string{"-b"}, sys); result.Status != tsc.ExitStatusSuccess {
		t.Fatalf("initial build status = %v, want ExitStatusSuccess", result.Status)
	}

	// A second, pre-canceled build has nothing to build; it must still report canceled.
	result := runPreCanceled(t, sys, []string{"-b"})
	if result.Status != tsc.ExitStatusCanceled {
		t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
	}
}

// TestTscCleanCancellation verifies that `tsc -b --clean` interrupted before it runs
// does not delete outputs and reports ExitStatusCanceled instead of success.
func TestTscCleanCancellation(t *testing.T) {
	t.Parallel()
	const outFile = "/home/src/workspaces/project/a.js"
	files := FileMap{
		"/home/src/workspaces/project/tsconfig.json": `{ "compilerOptions": { "composite": true } }`,
		"/home/src/workspaces/project/a.ts":          "export const a = 1;\n",
	}
	sys := newTestSys(&tscInput{
		commandLineArgs: []string{"-b"},
		files:           files,
	}, false)

	// Build first so there is an output for clean to (potentially) delete.
	if result := execute.CommandLine(context.Background(), sys, []string{"-b"}, sys); result.Status != tsc.ExitStatusSuccess {
		t.Fatalf("initial build status = %v, want ExitStatusSuccess", result.Status)
	}
	if !sys.FS().FileExists(outFile) {
		t.Fatalf("expected %s to exist after build", outFile)
	}

	// A pre-canceled clean must abort before deleting outputs and report canceled.
	result := runPreCanceled(t, sys, []string{"-b", "--clean"})
	if result.Status != tsc.ExitStatusCanceled {
		t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
	}
	if !sys.FS().FileExists(outFile) {
		t.Errorf("expected %s to survive a canceled clean", outFile)
	}
}

// TestTscCancellationSweep steps the cancellation point across the whole compile by
// increasing the poll threshold one step at a time, walking past the check phase and
// into emit. This covers windows a single fixed threshold would miss -- the
// no-emit-on-error recheck that runs during emit, the emit-returns-nil path, and
// declaration-emit type serialization -- and asserts that at no point does the run
// panic or report a partial result as success.
func TestTscCancellationSweep(t *testing.T) {
	t.Parallel()

	configs := []struct {
		name     string
		tsconfig string
	}{
		{
			// noEmitOnError + incremental: emit performs the internal no-emit-on-error
			// recheck, the path where a mid-emit cancellation makes Emit return nil.
			name:     "incremental noEmitOnError",
			tsconfig: `{ "compilerOptions": { "outDir": "out", "incremental": true, "noEmitOnError": true, "strict": true } }`,
		},
		{
			// declaration emit serializes inferred types via the checker during emit.
			name:     "declaration",
			tsconfig: `{ "compilerOptions": { "outDir": "out", "declaration": true, "noEmitOnError": true, "strict": true } }`,
		},
	}

	for _, cfg := range configs {
		t.Run(cfg.name, func(t *testing.T) {
			t.Parallel()
			files := FileMap{
				"/home/src/workspaces/project/tsconfig.json": cfg.tsconfig,
				"/home/src/workspaces/project/a.ts":          "export const a = 1;\nexport function f() { return { x: 1, y: 'z' }; }\n",
				"/home/src/workspaces/project/b.ts":          "export const b = 2;\nexport function g() { return [1, 2, 3] as const; }\n",
			}

			// Upper bound comfortably exceeds a full clean run's poll count for this
			// project, so the sweep covers check, the second global-diagnostics pass,
			// and emit.
			for threshold := int32(1); threshold <= 150; threshold++ {
				sys := newTestSys(&tscInput{
					commandLineArgs: []string{"--singleThreaded"},
					files:           files,
				}, false)
				ctx := newCancelAfterNPolls(threshold)

				var result tsc.CommandLineResult
				func() {
					defer func() {
						if r := recover(); r != nil {
							t.Fatalf("threshold=%d: panicked (want clean abort): %v", threshold, r)
						}
					}()
					result = execute.CommandLine(ctx, sys, []string{"--singleThreaded"}, sys)
				}()

				// Success only if cancellation never tripped (the build outran it);
				// otherwise it must be a clean Canceled. Anything else means partial
				// state leaked out.
				if ctx.tripped.Load() {
					if result.Status != tsc.ExitStatusCanceled {
						t.Fatalf("threshold=%d: status = %v, want ExitStatusCanceled", threshold, result.Status)
					}
				} else if result.Status != tsc.ExitStatusSuccess {
					t.Fatalf("threshold=%d: status = %v, want ExitStatusSuccess (cancellation never tripped)", threshold, result.Status)
				}
			}
		})
	}
}
