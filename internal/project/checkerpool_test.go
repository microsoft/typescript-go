package project

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func setupCheckerPoolSession(t *testing.T, opts CheckerPoolOptions) (*Session, *checkerPool) {
	t.Helper()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/src/tsconfig.json": `{ "compilerOptions": { "noLib": true } }`,
		"/src/index.ts":      "export const x: number = 1;",
	}
	fs := bundled.WrapFS(vfstest.FromMap(files, false))
	session := NewSession(&SessionInit{
		BackgroundCtx: context.Background(),
		Options: &SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       false,
			LoggingEnabled:     true,
			CheckerPoolOptions: opts,
		},
		FS:     fs,
		Logger: logging.NewTestLogger(),
	})
	session.DidOpenFile(context.Background(), "file:///src/index.ts", 1, "export const x: number = 1;", lsproto.LanguageKindTypeScript)

	snapshot := session.Snapshot()
	project := snapshot.ProjectCollection.ConfiguredProject("/src/tsconfig.json")
	assert.Assert(t, project != nil, "expected configured project")
	assert.Assert(t, project.checkerPool != nil, "expected checker pool")
	return session, project.checkerPool
}

// newTestCheckerPool creates a checker pool inside the current goroutine context
// (suitable for use inside synctest.Test) using the given program.
func newTestCheckerPool(program *compiler.Program, opts CheckerPoolOptions) *checkerPool {
	return newCheckerPool(opts, program, func(string) {})
}

func TestCheckerPoolDiagnosticsRouting(t *testing.T) {
	t.Parallel()
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 10 * time.Second})

	// Diagnostics requests should get checker at index 0.
	ctx := core.WithRequestID(context.Background(), "diag-req-1")
	ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeDiagnostics)
	c, release := pool.GetChecker(ctx, nil)
	assert.Assert(t, c != nil)
	assert.Assert(t, pool.checkers[0] == c, "diagnostics should use checker index 0")
	release()
}

func TestCheckerPoolQueryRouting(t *testing.T) {
	t.Parallel()
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 10 * time.Second})

	// Query requests should get a checker at index > 0.
	ctx := core.WithRequestID(context.Background(), "query-req-1")
	ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
	c, release := pool.GetChecker(ctx, nil)
	assert.Assert(t, c != nil)

	// Verify it's not the diagnostics checker slot.
	assert.Assert(t, pool.checkers[0] != c, "query should not use checker index 0")
	release()
}

func TestCheckerPoolRequestAffinity(t *testing.T) {
	t.Parallel()
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 10 * time.Second})

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx = core.WithRequestID(ctx, "req-affinity")
	ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)

	// First call acquires.
	c1, release1 := pool.GetChecker(ctx, nil)

	// Second call with same request ID while still held returns same checker (noop release).
	c2, release2 := pool.GetChecker(ctx, nil)
	release2()
	release1()

	assert.Assert(t, c1 == c2, "same request ID should return the same checker while held")

	// After release, same request should still get the same checker (cross-release affinity).
	c3, release3 := pool.GetChecker(ctx, nil)
	release3()

	assert.Assert(t, c1 == c3, "same request ID should return the same checker after release")
}

func TestCheckerPoolIdleCleanup(t *testing.T) {
	t.Parallel()
	// Get a real program to use for checker creation, then test the pool
	// with fake time via synctest.
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 5 * time.Second})

		// Create a checker via a diagnostics request.
		ctx := core.WithRequestID(context.Background(), "diag-cleanup")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeDiagnostics)
		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil)
		release()
		synctest.Wait()

		// Create a query checker as well.
		ctx2 := core.WithRequestID(context.Background(), "query-cleanup")
		ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeQuery)
		c2, release2 := pool.GetChecker(ctx2, nil)
		assert.Assert(t, c2 != nil)
		release2()
		synctest.Wait()

		// Both checkers should exist.
		pool.mu.Lock()
		assert.Assert(t, pool.checkers[0] != nil, "diagnostics checker should exist")
		var queryIdx int
		for i := 1; i < len(pool.checkers); i++ {
			if pool.checkers[i] != nil {
				queryIdx = i
				break
			}
		}
		assert.Assert(t, queryIdx > 0, "query checker should exist")
		pool.mu.Unlock()

		// Advance past idle timeout.
		time.Sleep(5 * time.Second)
		synctest.Wait()

		// After cleanup, all checkers should be nil.
		pool.mu.Lock()
		assert.Assert(t, pool.checkers[0] == nil, "diagnostics checker should be disposed after idle timeout")
		assert.Assert(t, pool.checkers[queryIdx] == nil, "query checker should be disposed after idle timeout")
		pool.mu.Unlock()
	})
}

func TestCheckerPoolFileAssociationCleanup(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()
	sourceFile := program.GetSourceFile("/src/index.ts")
	assert.Assert(t, sourceFile != nil)

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 5 * time.Second})

		// Create a query checker with file affinity.
		ctx := core.WithRequestID(context.Background(), "file-assoc-req")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
		c, release := pool.GetChecker(ctx, sourceFile)
		assert.Assert(t, c != nil)
		release()
		synctest.Wait()

		// File association should exist.
		pool.mu.Lock()
		_, hasAssoc := pool.fileAssociations[sourceFile]
		pool.mu.Unlock()
		assert.Assert(t, hasAssoc, "file should have a checker association")

		// Advance past idle timeout.
		time.Sleep(5 * time.Second)
		synctest.Wait()

		// File association should be cleared.
		pool.mu.Lock()
		_, hasAssoc = pool.fileAssociations[sourceFile]
		pool.mu.Unlock()
		assert.Assert(t, !hasAssoc, "file association should be cleared after checker disposal")
	})
}

func TestCheckerPoolMinCheckers(t *testing.T) {
	t.Parallel()
	// Requesting maxCheckers=1 should be clamped to 2.
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 1, IdleTimeout: 10 * time.Second})
	assert.Equal(t, pool.opts.MaxCheckers, 2)
	assert.Equal(t, len(pool.checkers), 2)
}

func TestCheckerPoolDefaultIdleTimeout(t *testing.T) {
	t.Parallel()
	// Zero idle timeout should default to 30s.
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4})
	assert.Equal(t, pool.opts.IdleTimeout, 30*time.Second)
}

func TestCheckerPoolQueryContention(t *testing.T) {
	t.Parallel()
	// maxCheckers=2 means 1 diagnostics + 1 query checker slot.
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 30 * time.Second})

		// Acquire the only query checker slot.
		ctx1 := core.WithRequestID(context.Background(), "query-hold")
		ctx1 = core.WithCheckerPurpose(ctx1, core.CheckerPurposeQuery)
		c1, release1 := pool.GetChecker(ctx1, nil)
		assert.Assert(t, c1 != nil)

		// A second query request from a different request ID should block.
		var c2Got bool
		go func() {
			ctx2 := core.WithRequestID(context.Background(), "query-wait")
			ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeQuery)
			c2, release2 := pool.GetChecker(ctx2, nil)
			assert.Assert(t, c2 != nil)
			c2Got = true
			release2()
		}()

		// Wait for goroutine to reach the cond.Wait.
		synctest.Wait()
		assert.Assert(t, !c2Got, "second query should be blocked while first holds the checker")

		// Release the first checker — second should unblock.
		release1()
		synctest.Wait()
		assert.Assert(t, c2Got, "second query should have acquired the checker after release")
	})
}

func TestCheckerPoolDiagnosticsContention(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 30 * time.Second})

		// Acquire the diagnostics checker.
		ctx1 := core.WithRequestID(context.Background(), "diag-hold")
		ctx1 = core.WithCheckerPurpose(ctx1, core.CheckerPurposeDiagnostics)
		c1, release1 := pool.GetChecker(ctx1, nil)
		assert.Assert(t, c1 != nil)

		// A second diagnostics request should block since there's only one diag checker.
		var c2Got bool
		go func() {
			ctx2 := core.WithRequestID(context.Background(), "diag-wait")
			ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeDiagnostics)
			c2, release2 := pool.GetChecker(ctx2, nil)
			assert.Assert(t, c2 != nil)
			c2Got = true
			release2()
		}()

		synctest.Wait()
		assert.Assert(t, !c2Got, "second diagnostics request should be blocked")

		// A query request should NOT be blocked (separate slot).
		ctx3 := core.WithRequestID(context.Background(), "query-concurrent")
		ctx3 = core.WithCheckerPurpose(ctx3, core.CheckerPurposeQuery)
		c3, release3 := pool.GetChecker(ctx3, nil)
		assert.Assert(t, c3 != nil)
		assert.Assert(t, c3 != c1, "query checker should be different from diagnostics checker")
		release3()

		// Release the diagnostics checker — second diag request should unblock.
		release1()
		synctest.Wait()
		assert.Assert(t, c2Got, "second diagnostics request should have acquired the checker after release")
	})
}

func TestCheckerPoolCanceledCheckerDisposal(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()
	sourceFile := program.GetSourceFile("/src/index.ts")
	assert.Assert(t, sourceFile != nil)

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})

		// Acquire a query checker and cancel it.
		ctx := core.WithRequestID(context.Background(), "cancel-test")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil)

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		c.GetDiagnostics(canceledCtx, sourceFile)
		assert.Assert(t, c.WasCanceled())

		// Release should dispose the canceled checker.
		release()
		synctest.Wait()

		// Next request should get a fresh checker.
		ctx2 := core.WithRequestID(context.Background(), "after-cancel")
		ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeQuery)
		c2, release2 := pool.GetChecker(ctx2, nil)
		assert.Assert(t, c2 != c, "should get a new checker, not the canceled one")
		release2()
	})
}

func TestCheckerPoolRequestAssociationCleanupOnDisposal(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 5 * time.Second})

		// Create a query checker with a request association.
		reqCtx, reqCancel := context.WithCancel(context.Background())
		defer reqCancel()
		ctx := core.WithRequestID(reqCtx, "assoc-cleanup-req")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil)

		// Cancel the checker to trigger disposal on release.
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		sourceFile := program.GetSourceFile("/src/index.ts")
		c.GetDiagnostics(canceledCtx, sourceFile)
		assert.Assert(t, c.WasCanceled())

		release()
		synctest.Wait()

		// Request association should be cleared after checker disposal.
		pool.mu.Lock()
		_, hasReqAssoc := pool.requestAssociations["assoc-cleanup-req"]
		pool.mu.Unlock()
		assert.Assert(t, !hasReqAssoc, "request association should be cleared after checker disposal")
	})
}

func TestCheckerPoolRequestAssociationCleanupOnContextDone(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})

		// Create a cancellable context to simulate request lifecycle.
		reqCtx, reqCancel := context.WithCancel(context.Background())
		ctx := core.WithRequestID(reqCtx, "ctx-cleanup-req")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)

		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil)
		release()
		synctest.Wait()

		// Association should still exist after release.
		pool.mu.Lock()
		_, hasAssoc := pool.requestAssociations["ctx-cleanup-req"]
		pool.mu.Unlock()
		assert.Assert(t, hasAssoc, "request association should persist after release")

		// Cancel the request context — association should be cleaned up.
		reqCancel()
		synctest.Wait()

		pool.mu.Lock()
		_, hasAssoc = pool.requestAssociations["ctx-cleanup-req"]
		pool.mu.Unlock()
		assert.Assert(t, !hasAssoc, "request association should be cleaned up after context cancellation")
	})
}

func TestCheckerPoolDiagnosticsRecreatedAfterIdleDisposal(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 5 * time.Second})

		// Create and release diagnostics checker.
		ctx := core.WithRequestID(context.Background(), "diag-recreate-1")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeDiagnostics)
		c1, release1 := pool.GetChecker(ctx, nil)
		assert.Assert(t, c1 != nil)
		release1()
		synctest.Wait()

		// Advance past idle timeout so it gets disposed.
		time.Sleep(5 * time.Second)
		synctest.Wait()

		pool.mu.Lock()
		assert.Assert(t, pool.checkers[0] == nil, "diagnostics checker should be disposed")
		pool.mu.Unlock()

		// Request diagnostics checker again — should get a fresh one.
		ctx2 := core.WithRequestID(context.Background(), "diag-recreate-2")
		ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeDiagnostics)
		c2, release2 := pool.GetChecker(ctx2, nil)
		assert.Assert(t, c2 != nil, "diagnostics checker should be re-created")
		assert.Assert(t, c2 != c1, "should be a new checker instance")
		release2()
	})
}

func TestCheckerPoolCrossReleaseAffinityWithContention(t *testing.T) {
	t.Parallel()
	// maxCheckers=2: 1 diagnostics + 1 query slot.
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 30 * time.Second})

		reqCtx, reqCancel := context.WithCancel(context.Background())
		defer reqCancel()

		// Request A acquires the only query slot.
		ctxA := core.WithRequestID(reqCtx, "req-A")
		ctxA = core.WithCheckerPurpose(ctxA, core.CheckerPurposeQuery)
		cA, releaseA := pool.GetChecker(ctxA, nil)
		assert.Assert(t, cA != nil)
		releaseA()
		synctest.Wait()

		// Request B takes the query slot while A is released.
		ctxB := core.WithRequestID(context.Background(), "req-B")
		ctxB = core.WithCheckerPurpose(ctxB, core.CheckerPurposeQuery)
		cB, releaseB := pool.GetChecker(ctxB, nil)
		assert.Assert(t, cB != nil)

		// Request A reacquires — should block because B holds the slot.
		var cA2 *checker.Checker
		var reacquired bool
		go func() {
			c, release := pool.GetChecker(ctxA, nil)
			cA2 = c
			reacquired = true
			release()
		}()

		synctest.Wait()
		assert.Assert(t, !reacquired, "request A should block while B holds the slot")

		// Release B — A should unblock and get the same checker.
		releaseB()
		synctest.Wait()
		assert.Assert(t, reacquired, "request A should unblock after B releases")
		assert.Assert(t, cA2 == cA, "request A should get the same checker on reacquire")
	})
}

func TestCheckerPoolNoRequestID(t *testing.T) {
	t.Parallel()
	_, pool := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 10 * time.Second})

	// Calls without a request ID should still work (e.g., callhierarchy uses context.Background()).
	ctx := context.Background()

	c1, release1 := pool.GetChecker(ctx, nil)
	assert.Assert(t, c1 != nil)
	release1()

	c2, release2 := pool.GetChecker(ctx, nil)
	assert.Assert(t, c2 != nil)
	release2()

	// Without request ID, no affinity guarantee — just verify it doesn't crash.
}

func TestCheckerPoolDiagnosticsCrossReleaseAffinity(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})

		reqCtx, reqCancel := context.WithCancel(context.Background())
		defer reqCancel()
		ctx := core.WithRequestID(reqCtx, "diag-affinity")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeDiagnostics)

		c1, release1 := pool.GetChecker(ctx, nil)
		assert.Assert(t, c1 != nil)
		assert.Assert(t, pool.checkers[0] == c1, "should be the diagnostics checker")
		release1()
		synctest.Wait()

		// Same request reacquiring diagnostics should get the same checker.
		c2, release2 := pool.GetChecker(ctx, nil)
		assert.Assert(t, c2 == c1, "same diagnostics request should get the same checker after release")
		release2()
	})
}

func TestCheckerPoolDiscardDisposesIdleCheckers(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})

		// Create both a diagnostics and a query checker.
		ctx1 := core.WithRequestID(context.Background(), "obs-diag")
		ctx1 = core.WithCheckerPurpose(ctx1, core.CheckerPurposeDiagnostics)
		c1, release1 := pool.GetChecker(ctx1, nil)
		assert.Assert(t, c1 != nil)
		release1()
		synctest.Wait()

		ctx2 := core.WithRequestID(context.Background(), "obs-query")
		ctx2 = core.WithCheckerPurpose(ctx2, core.CheckerPurposeQuery)
		c2, release2 := pool.GetChecker(ctx2, nil)
		assert.Assert(t, c2 != nil)
		release2()
		synctest.Wait()

		// Both checkers should exist before Discard.
		pool.mu.Lock()
		assert.Assert(t, pool.checkers[0] != nil, "diagnostics checker should exist")
		pool.mu.Unlock()

		// Discard should dispose all idle checkers immediately.
		pool.Discard()

		pool.mu.Lock()
		allNil := true
		for _, c := range pool.checkers {
			if c != nil {
				allNil = false
				break
			}
		}
		pool.mu.Unlock()
		assert.Assert(t, allNil, "all idle checkers should be disposed after Discard")
	})
}

func TestCheckerPoolDiscardHeldCheckerDisposedOnRelease(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})

		// Acquire a checker and hold it.
		ctx := core.WithRequestID(context.Background(), "held-obs")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil)

		// Find which slot it's in.
		pool.mu.Lock()
		var heldIndex int
		for i := 1; i < len(pool.checkers); i++ {
			if pool.checkers[i] == c {
				heldIndex = i
				break
			}
		}
		pool.mu.Unlock()
		assert.Assert(t, heldIndex > 0, "should find the held checker")

		// Discard while checker is held — should NOT dispose it.
		pool.Discard()

		pool.mu.Lock()
		assert.Assert(t, pool.checkers[heldIndex] == c, "held checker should survive Discard")
		pool.mu.Unlock()

		// Release — checker should be disposed on next cleanup.
		release()
		synctest.Wait()

		pool.mu.Lock()
		assert.Assert(t, pool.checkers[heldIndex] == nil, "checker should be disposed after release on discarded pool")
		pool.mu.Unlock()
	})
}

func TestCheckerPoolDiscardStillFunctional(t *testing.T) {
	t.Parallel()
	session, _ := setupCheckerPoolSession(t, CheckerPoolOptions{MaxCheckers: 2, IdleTimeout: 10 * time.Second})
	ls, err := session.GetLanguageService(context.Background(), "file:///src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()

	synctest.Test(t, func(t *testing.T) {
		pool := newTestCheckerPool(program, CheckerPoolOptions{MaxCheckers: 4, IdleTimeout: 30 * time.Second})
		pool.Discard()

		// Pool should still work — GetChecker should create a fresh checker.
		ctx := core.WithRequestID(context.Background(), "post-obs")
		ctx = core.WithCheckerPurpose(ctx, core.CheckerPurposeQuery)
		c, release := pool.GetChecker(ctx, nil)
		assert.Assert(t, c != nil, "should still create checkers after Discard")

		// Find the slot.
		pool.mu.Lock()
		var idx int
		for i := 1; i < len(pool.checkers); i++ {
			if pool.checkers[i] == c {
				idx = i
				break
			}
		}
		pool.mu.Unlock()
		assert.Assert(t, idx > 0, "checker should be in a query slot")

		// Release — should dispose immediately (obsolete pool).
		release()
		synctest.Wait()

		pool.mu.Lock()
		assert.Assert(t, pool.checkers[idx] == nil, "checker should be disposed immediately after release on discarded pool")
		pool.mu.Unlock()
	})
}
