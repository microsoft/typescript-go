package compiler_test

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// cancelAfterNPolls cancels itself after Err has been polled pollThreshold times,
// then stays canceled, so cancellation lands after checking has begun.
type cancelAfterNPolls struct {
	context.Context
	pollThreshold int32
	polls         atomic.Int32
	tripped       atomic.Bool
	done          chan struct{}
}

func newCancelAfterNPolls(pollThreshold int32) *cancelAfterNPolls {
	return &cancelAfterNPolls{Context: context.Background(), pollThreshold: pollThreshold, done: make(chan struct{})}
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

// TestGetGlobalDiagnosticsAfterCancellation pins the checker-pool behavior that
// GetGlobalDiagnostics skips a checker canceled mid-check rather than reusing it
// (which panics in checkNotCanceled). This guards every caller, including
// emitBuildInfo's error-state probe, which has no cancellation check of its own.
func TestGetGlobalDiagnosticsAfterCancellation(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	fs := bundled.WrapFS(vfstest.FromMap[any](nil, false /*useCaseSensitiveFileNames*/))

	// Many statements with type errors across a few files: global diagnostics exist,
	// and the checker polls often so cancellation lands mid-check.
	var src strings.Builder
	for i := range 50 {
		src.WriteString("export const v")
		src.WriteString(strings.Repeat("x", i+1))
		src.WriteString(`: number = "not a number";` + "\n")
	}
	for _, name := range []string{"/src/a.ts", "/src/b.ts", "/src/c.ts"} {
		_ = fs.WriteFile(name, src.String())
	}

	// One checker, so the mid-check cancellation marks the same checker that the
	// subsequent GetGlobalDiagnostics will visit.
	oneChecker := 1
	program := compiler.NewProgram(compiler.ProgramOptions{
		Config: &tsoptions.ParsedCommandLine{
			ParsedConfig: &core.ParsedOptions{
				FileNames:       []string{"/src/a.ts", "/src/b.ts", "/src/c.ts"},
				CompilerOptions: &core.CompilerOptions{Strict: core.TSTrue, Checkers: &oneChecker},
			},
		},
		Host: compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil),
	})

	ctx := newCancelAfterNPolls(5)

	// Drive checking only to leave the checker canceled; the result is discarded.
	_ = program.GetSemanticDiagnostics(ctx, nil)
	if !ctx.tripped.Load() {
		t.Fatal("expected cancellation to trip during checking, but it never did")
	}

	// The real assertion: this must not panic on the canceled checker.
	_ = program.GetGlobalDiagnostics(ctx)
}
