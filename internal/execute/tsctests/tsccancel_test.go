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

// TestTscPreCanceledCompilation verifies that a context canceled before the compile
// starts aborts immediately with ExitStatusCanceled and reports no diagnostics,
// without panicking on checker reuse.
func TestTscPreCanceledCompilation(t *testing.T) {
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
			// forEachCheckerGroupDo reuses it across files.
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
			cancel() // simulate SIGINT delivered before the compile starts

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

// cancelAfterNPolls is a context that reports itself canceled only after Err has
// been polled pollThreshold times while still uncanceled. The checker polls
// ctx.Err() between top-level statements (checkSourceElements), so this lands the
// cancellation *after* checking has begun rather than before it starts -- the case
// a pre-canceled context cannot exercise. Once tripped it stays canceled.
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

func (c *cancelAfterNPolls) Done() <-chan struct{} {
	return c.done
}

// TestTscMidCheckCancellation cancels the context after type-checking has already
// begun (not before it starts). This exercises the paths a pre-canceled context
// skips: a checker actually runs, sets wasCanceled, and would panic in
// checkNotCanceled on the second GetGlobalDiagnostics / on reuse across files if
// the cancellation guards were missing.
func TestTscMidCheckCancellation(t *testing.T) {
	t.Parallel()

	// Many top-level statements per file so the checker polls ctx.Err() often and
	// cancellation reliably lands mid-check, across more than one file.
	var badStatements strings.Builder
	for i := range 50 {
		badStatements.WriteString("export const v")
		badStatements.WriteString(strings.Repeat("x", i+1))
		badStatements.WriteString(`: number = "not a number";` + "\n")
	}
	badSrc := badStatements.String()

	// Inferred return types force the checker to serialize types during declaration
	// emit (SerializeReturnTypeForSignature -> node reuse -> checkNotCanceled), a
	// distinct reuse path from plain semantic checking.
	var inferredSrc strings.Builder
	for i := range 50 {
		inferredSrc.WriteString("export function make")
		inferredSrc.WriteString(strings.Repeat("x", i+1))
		inferredSrc.WriteString("() { return { a: 1, b: 'x', deep: [1, 2, 3] as const }; }\n")
	}

	testCases := []struct {
		name     string
		args     []string
		tsconfig string
		src      string
	}{
		{
			// Single checker reused across files: after it cancels on an early file,
			// forEachCheckerGroupDo must stop feeding it later files.
			name:     "single checker",
			args:     []string{"--noEmit", "--singleThreaded"},
			tsconfig: `{ "compilerOptions": { "noEmit": true, "strict": true } }`,
			src:      badSrc,
		},
		{
			// tsc -b through the incremental program + orchestrator, which also reaches
			// GetGlobalDiagnostics from emitBuildInfo -> ensureHasErrorsForState.
			name:     "build mode",
			args:     []string{"-b"},
			tsconfig: `{ "compilerOptions": { "composite": true, "strict": true } }`,
			src:      badSrc,
		},
		{
			// Declaration emit reuses checkers to serialize types; a canceled checker
			// must not be handed to GetDeclarationDiagnostics.
			name:     "declaration emit",
			args:     []string{"--noEmit", "--declaration", "--singleThreaded"},
			tsconfig: `{ "compilerOptions": { "noEmit": true, "declaration": true, "strict": true } }`,
			src:      inferredSrc.String(),
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

			// Let checking start, then cancel. The threshold is small relative to the
			// number of statements so cancellation lands well before checking finishes.
			ctx := newCancelAfterNPolls(5)

			resultCh := make(chan tsc.CommandLineResult, 1)
			go func() {
				resultCh <- execute.CommandLine(ctx, sys, tc.args, sys)
			}()

			select {
			case result := <-resultCh:
				// The run must abort (not run to completion) once canceled mid-check,
				// and must not panic in checkNotCanceled.
				if result.Status != tsc.ExitStatusCanceled {
					t.Errorf("status = %v, want ExitStatusCanceled", result.Status)
				}
				// The cancellation must actually have landed mid-check, otherwise this
				// test is not exercising what it claims.
				if !ctx.tripped.Load() {
					t.Error("expected cancellation to trip during checking, but it never did")
				}
			case <-time.After(30 * time.Second):
				t.Fatal("compile did not abort after mid-check cancellation")
			}
		})
	}
}

// TestTscCancellationSweep cancels at every successive point in the compile by
// increasing the poll threshold one step at a time. It walks cancellation past the
// check phase and into emit, covering windows a single fixed threshold would miss:
// the no-emit-on-error recheck that runs during emit, the emit-returns-nil path, and
// declaration-emit type serialization. At no threshold may the run panic, and it
// must always end Canceled or Success (if it finished before the trip).
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
