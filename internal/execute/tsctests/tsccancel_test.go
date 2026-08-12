package tsctests

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
)

// cancelAfterNPolls is a context that cancels itself after Err has been polled
// pollThreshold times, then stays canceled. The checker polls between top-level
// statements (checkSourceElements), so the threshold selects how far into checking
// the cancellation lands -- something a pre-canceled context cannot exercise.
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
// ExitStatusCanceled and never reports its incomplete diagnostics. Without the
// checker-reuse guards, the mid-check cases panic in checkNotCanceled -- a canceled
// checker fed more files, or asked for global/declaration diagnostics again.
func TestTscCancellationAborts(t *testing.T) {
	t.Parallel()

	// Many distinctly-named statements per file, so a mid-check cancellation reliably
	// lands while checking. The type errors are the signal: reported means the checker
	// ran to completion, absent means it aborted.
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
			// One checker for all files, so forEachCheckerGroupDo must stop feeding it
			// files after it cancels on an early one.
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

			if result.Status != tsc.ExitStatusCanceled {
				t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
			}
			if out := sys.getOutput(true); strings.Contains(out, "error TS") {
				t.Errorf("expected no diagnostics after cancellation; got output:\n%s", out)
			}
			// A mid-check case that never tripped isn't testing what it claims.
			if tc.midCheck && !midChecked {
				t.Error("expected cancellation to trip during checking, but it never did")
			}
		})
	}
}

// runWithCancellation runs the command line under a context canceled either before
// the run starts (midCheck false) or after checking has begun (midCheck true). The
// returned bool reports whether a mid-check cancellation actually tripped.
func runWithCancellation(t *testing.T, sys *TestSys, args []string, midCheck bool) (tsc.CommandLineResult, bool) {
	t.Helper()
	var (
		ctx        context.Context
		midChecked func() bool
	)
	if midCheck {
		// Small relative to the statement count, so cancellation lands well before
		// checking finishes.
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

// runPreCanceled runs the command line on an existing sys under an already-canceled
// context.
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

// TestTscBuildCancellationUpToDate verifies that `tsc -b` honors an interrupt even
// with nothing to build. An up-to-date root has no upstream to wait on, so the
// no-build path is where the interrupt would otherwise be swallowed as success.
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

	result := runPreCanceled(t, sys, []string{"-b"})
	if result.Status != tsc.ExitStatusCanceled {
		t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
	}
}

// TestTscBuildCancellationDuringGraph verifies that an interrupt during graph
// construction stops resolving configs and reports ExitStatusCanceled, rather than
// building from the partial graph or panicking on a task that was never created.
func TestTscBuildCancellationDuringGraph(t *testing.T) {
	t.Parallel()

	// Each project references the next, so graph construction has several levels to
	// resolve before it completes.
	files := FileMap{}
	const projects = 8
	for i := range projects {
		dir := fmt.Sprintf("/home/src/workspaces/project/p%d", i)
		references := ""
		if i+1 < projects {
			references = fmt.Sprintf(`, "references": [{ "path": "../p%d" }]`, i+1)
		}
		files[dir+"/tsconfig.json"] = fmt.Sprintf(
			`{ "compilerOptions": { "composite": true }%s }`, references)
		files[dir+fmt.Sprintf("/p%d.ts", i)] = "export {}\n"
	}

	args := []string{"-b", "/home/src/workspaces/project/p0"}

	// Sweep rather than fix one threshold, so the graph phase keeps being exercised
	// as the number of polls before it shifts.
	for threshold := int32(0); threshold <= 20; threshold++ {
		sys := newTestSys(&tscInput{commandLineArgs: args, files: files}, false)
		ctx := newCancelAfterNPolls(threshold)

		var result tsc.CommandLineResult
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("threshold=%d: panicked (want clean abort): %v", threshold, r)
				}
			}()
			result = execute.CommandLine(ctx, sys, args, sys)
		}()

		if ctx.tripped.Load() {
			if result.Status != tsc.ExitStatusCanceled {
				t.Fatalf("threshold=%d: status = %v, want ExitStatusCanceled", threshold, result.Status)
			}
		} else if result.Status != tsc.ExitStatusSuccess {
			t.Fatalf("threshold=%d: status = %v, want ExitStatusSuccess (cancellation never tripped)", threshold, result.Status)
		}
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

	// Build first so there is an output for clean to delete.
	if result := execute.CommandLine(context.Background(), sys, []string{"-b"}, sys); result.Status != tsc.ExitStatusSuccess {
		t.Fatalf("initial build status = %v, want ExitStatusSuccess", result.Status)
	}
	if !sys.FS().FileExists(outFile) {
		t.Fatalf("expected %s to exist after build", outFile)
	}

	result := runPreCanceled(t, sys, []string{"-b", "--clean"})
	if result.Status != tsc.ExitStatusCanceled {
		t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
	}
	if !sys.FS().FileExists(outFile) {
		t.Errorf("expected %s to survive a canceled clean", outFile)
	}
}

// TestTscCancellationSweep steps the cancellation point across the whole compile,
// past checking and into emit, asserting the run never panics or reports a partial
// result as success. A single fixed threshold would miss the narrow windows: the
// no-emit-on-error recheck, the emit-returns-nil path, and declaration-emit type
// serialization.
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

			// The bound exceeds a full run's poll count for this project, so the sweep
			// reaches check, the second global-diagnostics pass, and emit.
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

				// Success is only legitimate if the build outran the cancellation;
				// anything else means partial state leaked out.
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
