package project

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
)

// checkerHeldAnonymous is a sentinel stored in heldBy when a checker is held
// by a caller that has no request ID (e.g., context.Background()). This
// distinguishes "held without ID" from "not held" (empty string).
const checkerHeldAnonymous = "<anonymous>"

type CheckerPoolOptions struct {
	// MaxCheckers controls the total number of checker slots per project
	// (1 dedicated diagnostics checker + N-1 query checkers). Minimum 2.
	// Zero uses the default (4).
	MaxCheckers int
	// IdleTimeout controls how long an idle checker is kept
	// before being disposed. Zero uses the default (30s).
	IdleTimeout time.Duration
}

// checkerPool manages a set of type checkers for a project. It maintains a
// dedicated diagnostics checker (index 0) that provides consistent walk-order
// for diagnostic operations, plus a set of ephemeral query checkers (indices 1+)
// for language-service operations like hover, completions, go-to-definition, etc.
//
// All checkers are created lazily and automatically disposed after an idle
// timeout. Each pool manages its own cleanup timer, so orphaned pools (from
// replaced programs) still clean up after themselves. The diagnostics checker
// is separated from query checkers so that diagnostic walk order is not
// influenced by other operations, but it is still ephemeral.
//
// Concurrency is managed via two buffered channels used as semaphores:
// diagSem (capacity 1) and querySem (capacity N-1). A goroutine sends
// a value to claim a slot, and receives to release it.
type checkerPool struct {
	opts    CheckerPoolOptions
	program *compiler.Program

	mu sync.Mutex

	// checkers[0] is the dedicated diagnostics checker.
	// checkers[1:] are ephemeral query checkers.
	checkers            []*checker.Checker
	heldBy              []string                // heldBy[i] is the requestID holding checker i, checkerHeldAnonymous, or "" if not held
	fileAssociations    map[*ast.SourceFile]int // file → query checker index (1+)
	requestAssociations map[string]int          // requestID → checker index

	// lastReleased tracks when each checker was last released.
	lastReleased []time.Time

	// cleanupTimer is reset each time a checker is released.
	// When it fires, idle checkers are disposed.
	cleanupTimer *time.Timer

	// diagSem has capacity 1 — the diagnostics checker slot.
	// querySem has capacity opts.MaxCheckers-1 — one slot per query checker.
	// Send to acquire a slot, receive to release it.
	diagSem  chan struct{}
	querySem chan struct{}

	log                    func(msg string)
	globalDiagAccumulated  []*ast.Diagnostic
	globalDiagChanged      bool
	globalDiagCheckerCount []int // per-checker count of globals last seen
}

var _ compiler.CheckerPool = (*checkerPool)(nil)

func newCheckerPool(opts CheckerPoolOptions, program *compiler.Program, log func(msg string)) *checkerPool {
	if opts.MaxCheckers <= 0 {
		opts.MaxCheckers = 4
	} else if opts.MaxCheckers < 2 {
		opts.MaxCheckers = 2 // at least 1 diagnostics + 1 query checker
	}
	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = 30 * time.Second
	}
	querySlots := opts.MaxCheckers - 1
	pool := &checkerPool{
		program:                program,
		opts:                   opts,
		checkers:               make([]*checker.Checker, opts.MaxCheckers),
		heldBy:                 make([]string, opts.MaxCheckers),
		fileAssociations:       make(map[*ast.SourceFile]int),
		requestAssociations:    make(map[string]int),
		lastReleased:           make([]time.Time, opts.MaxCheckers),
		diagSem:                make(chan struct{}, 1),
		querySem:               make(chan struct{}, querySlots),
		log:                    log,
		globalDiagCheckerCount: make([]int, opts.MaxCheckers),
	}
	return pool
}

// holdTag returns the value to store in heldBy for the given request ID.
func holdTag(requestID string) string {
	if requestID == "" {
		return checkerHeldAnonymous
	}
	return requestID
}

func (p *checkerPool) GetChecker(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	purpose := core.GetCheckerPurpose(ctx)
	requestID := core.GetRequestID(ctx)

	if purpose == core.CheckerPurposeDiagnostics {
		return p.getDiagnosticsChecker(ctx, requestID)
	}
	return p.getQueryChecker(ctx, requestID, file)
}

// tryReacquireForRequest checks whether the given request already has an
// associated checker. If so, it either returns the checker directly (still held)
// or reacquires it by claiming a semaphore slot. The caller must provide the
// appropriate semaphore channel.
//
// Returns (checker, release, true) if the request was served (either still held
// or reclaimed). Returns (nil, nil, false) if the caller must proceed with
// normal acquisition — in this case, a semaphore slot has already been claimed.
// Must NOT be called with p.mu held.
func (p *checkerPool) tryReacquireForRequest(requestID string, sem chan<- struct{}) (*checker.Checker, func(), bool) {
	if requestID == "" {
		sem <- struct{}{}
		return nil, nil, false
	}

	p.mu.Lock()
	index, ok := p.requestAssociations[requestID]
	if !ok {
		p.mu.Unlock()
		sem <- struct{}{}
		return nil, nil, false
	}

	c := p.checkers[index]
	if c == nil {
		delete(p.requestAssociations, requestID)
		p.mu.Unlock()
		sem <- struct{}{}
		return nil, nil, false
	}

	held := p.heldBy[index]
	if held == requestID {
		// Same request, checker still held — return without claiming a slot.
		p.mu.Unlock()
		return c, noop, true
	}

	if held == "" {
		// Same request reacquiring after release — need a semaphore slot.
		p.mu.Unlock()
		sem <- struct{}{}
		p.mu.Lock()
		// Re-check: checker may have been disposed while waiting for the slot.
		if cc := p.checkers[index]; cc == c && p.heldBy[index] == "" {
			p.heldBy[index] = requestID
			p.mu.Unlock()
			return c, p.createRelease(requestID, index, c), true
		}
		p.mu.Unlock()
		// Checker was replaced/disposed while waiting for the slot.
		// The slot is still claimed; the caller will use it for normal acquisition.
		return nil, nil, false
	}

	// Checker held by another request — claim a slot normally.
	p.mu.Unlock()
	sem <- struct{}{}
	return nil, nil, false
}

// getDiagnosticsChecker returns the dedicated diagnostics checker (index 0).
// Creates it on first use. Blocks on diagSem if it's currently in use.
func (p *checkerPool) getDiagnosticsChecker(ctx context.Context, requestID string) (*checker.Checker, func()) {
	const diagIndex = 0

	if c, release, ok := p.tryReacquireForRequest(requestID, p.diagSem); ok {
		return c, release
	}

	// Token consumed — proceed with normal acquisition.
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.checkers[diagIndex] == nil {
		p.log("checkerpool: Creating diagnostics checker")
		c, _ := checker.NewChecker(p.program)
		p.checkers[diagIndex] = c
	}

	c := p.checkers[diagIndex]
	p.heldBy[diagIndex] = holdTag(requestID)
	p.log(fmt.Sprintf("checkerpool: Acquired diagnostics checker for request %s", holdTag(requestID)))
	if requestID != "" {
		if _, alreadyRegistered := p.requestAssociations[requestID]; !alreadyRegistered {
			p.requestAssociations[requestID] = diagIndex
			p.registerRequestCleanup(ctx, requestID)
		}
	}
	return c, p.createRelease(requestID, diagIndex, c)
}

// getQueryChecker returns an ephemeral query checker from indices 1+.
// Uses request affinity, then file affinity, then finds/creates.
// Blocks on querySem if all query slots are in use.
func (p *checkerPool) getQueryChecker(ctx context.Context, requestID string, file *ast.SourceFile) (*checker.Checker, func()) {
	if c, release, ok := p.tryReacquireForRequest(requestID, p.querySem); ok {
		return c, release
	}

	// Token consumed — proceed with normal acquisition.
	p.mu.Lock()
	defer p.mu.Unlock()

	// Try file affinity.
	if file != nil {
		if index, ok := p.fileAssociations[file]; ok && index > 0 {
			if c := p.checkers[index]; c != nil && p.heldBy[index] == "" {
				p.heldBy[index] = holdTag(requestID)
				if requestID != "" {
					if _, alreadyRegistered := p.requestAssociations[requestID]; !alreadyRegistered {
						p.requestAssociations[requestID] = index
						p.registerRequestCleanup(ctx, requestID)
					}
				}
				return c, p.createRelease(requestID, index, c)
			}
		}
	}

	// Find any available query checker or create one.
	c, index := p.findOrCreateQueryCheckerLocked()
	p.heldBy[index] = holdTag(requestID)
	p.log(fmt.Sprintf("checkerpool: Acquired query checker %d for request %s", index, holdTag(requestID)))
	if requestID != "" {
		if _, alreadyRegistered := p.requestAssociations[requestID]; !alreadyRegistered {
			p.requestAssociations[requestID] = index
			p.registerRequestCleanup(ctx, requestID)
		}
	}
	if file != nil {
		p.fileAssociations[file] = index
	}
	return c, p.createRelease(requestID, index, c)
}

// findOrCreateQueryCheckerLocked returns an idle query checker or creates one
// in the first empty slot. The semaphore guarantees at least one slot is
// available. Must be called with p.mu held.
func (p *checkerPool) findOrCreateQueryCheckerLocked() (*checker.Checker, int) {
	// Prefer an existing idle checker.
	for i := 1; i < len(p.checkers); i++ {
		if c := p.checkers[i]; c != nil && p.heldBy[i] == "" {
			return c, i
		}
	}
	// Create in the first empty slot.
	for i := 1; i < len(p.checkers); i++ {
		if p.checkers[i] == nil {
			p.log(fmt.Sprintf("checkerpool: Creating query checker %d", i))
			c, _ := checker.NewChecker(p.program)
			p.checkers[i] = c
			return c, i
		}
	}
	panic("checkerpool: no available query slot despite holding semaphore token")
}

func (p *checkerPool) createRelease(requestID string, index int, c *checker.Checker) func() {
	return sync.OnceFunc(func() {
		p.mu.Lock()

		if c.WasCanceled() {
			// Canceled checkers must be disposed.
			p.log(fmt.Sprintf("checkerpool: Checker %d for request %s was canceled, disposing", index, holdTag(requestID)))
			p.disposeCheckerLocked(index, c)
		} else {
			p.mergeGlobalDiagnosticsFromCheckerLocked(index, c)
			if p.opts.IdleTimeout == 0 {
				// Pool is discarded — dispose immediately instead of caching.
				p.log(fmt.Sprintf("checkerpool: Pool discarded, disposing checker %d for request %s on release", index, holdTag(requestID)))
				p.disposeCheckerLocked(index, c)
			} else {
				p.heldBy[index] = ""
				// Track release time and schedule cleanup only for live checkers.
				p.lastReleased[index] = time.Now()
				p.scheduleCleanupLocked()
			}
		}

		// Unlock before releasing the semaphore slot. If we received from
		// the channel while holding p.mu, a woken goroutine could immediately
		// try to acquire p.mu, risking priority inversion or unnecessary
		// contention.
		p.mu.Unlock()

		// Release the semaphore slot.
		if index == 0 {
			<-p.diagSem
		} else {
			<-p.querySem
		}
	})
}

// registerRequestCleanup uses context.AfterFunc to delete the request
// association when the request context is done. This prevents the map
// from growing unboundedly with completed request IDs.
// Must be called with p.mu held; the cleanup runs asynchronously.
func (p *checkerPool) registerRequestCleanup(ctx context.Context, requestID string) {
	context.AfterFunc(ctx, func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		delete(p.requestAssociations, requestID)
	})
}

// scheduleCleanupLocked resets (or starts) the cleanup timer so it fires
// idleTimeout after the most recent checker release. When the timer fires,
// it disposes any checkers that have been idle long enough.
// Must be called with p.mu held.
func (p *checkerPool) scheduleCleanupLocked() {
	if p.cleanupTimer != nil {
		p.cleanupTimer.Reset(p.opts.IdleTimeout)
	} else {
		p.cleanupTimer = time.AfterFunc(p.opts.IdleTimeout, p.cleanupIdleCheckers)
	}
}

// cleanupIdleCheckers disposes all checkers (diagnostics and query) that have
// been idle for longer than the idle timeout. If any checkers are still alive
// but not yet idle enough, the timer is rescheduled.
func (p *checkerPool) cleanupIdleCheckers() {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now()
	var earliestRemaining time.Time
	for i := range len(p.checkers) {
		c := p.checkers[i]
		if c == nil || p.heldBy[i] != "" {
			continue
		}
		if p.lastReleased[i].IsZero() {
			continue
		}
		idle := now.Sub(p.lastReleased[i])
		if idle >= p.opts.IdleTimeout {
			p.log(fmt.Sprintf("checkerpool: Disposing idle checker %d (idle %v)", i, idle))
			p.disposeCheckerLocked(i, c)
		} else if earliestRemaining.IsZero() || p.lastReleased[i].Before(earliestRemaining) {
			earliestRemaining = p.lastReleased[i]
		}
	}
	// If there are remaining checkers not yet idle enough, reschedule.
	if !earliestRemaining.IsZero() {
		remaining := max(p.opts.IdleTimeout-now.Sub(earliestRemaining), time.Second)
		p.cleanupTimer = time.AfterFunc(remaining, p.cleanupIdleCheckers)
	} else {
		p.cleanupTimer = nil
	}
}

// disposeCheckerLocked removes a checker from the pool and clears all associations
// (file and request) that reference it. Must be called with p.mu held.
func (p *checkerPool) disposeCheckerLocked(index int, c *checker.Checker) {
	debug.Assert(p.checkers[index] == c)
	p.checkers[index] = nil
	p.heldBy[index] = ""
	p.globalDiagCheckerCount[index] = 0
	p.lastReleased[index] = time.Time{}
	for file, idx := range p.fileAssociations {
		if idx == index {
			delete(p.fileAssociations, file)
		}
	}
	for req, idx := range p.requestAssociations {
		if idx == index {
			delete(p.requestAssociations, req)
		}
	}
}

// mergeGlobalDiagnosticsFromCheckerLocked checks if the given checker has produced new global
// diagnostics since the last time we looked, and if so merges them into the accumulated set.
// Must be called with p.mu held.
func (p *checkerPool) mergeGlobalDiagnosticsFromCheckerLocked(index int, c *checker.Checker) {
	globals := c.GetGlobalDiagnostics()
	if len(globals) == p.globalDiagCheckerCount[index] {
		return
	}
	p.globalDiagCheckerCount[index] = len(globals)
	before := len(p.globalDiagAccumulated)
	p.globalDiagAccumulated = compiler.SortAndDeduplicateDiagnostics(append(p.globalDiagAccumulated, globals...))
	if len(p.globalDiagAccumulated) != before {
		p.globalDiagChanged = true
	}
}

// GetGlobalDiagnostics returns the accumulated global diagnostics collected from
// all checkers that have been used so far in this pool's lifetime.
func (p *checkerPool) GetGlobalDiagnostics() []*ast.Diagnostic {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Clone(p.globalDiagAccumulated)
}

// TakeNewGlobalDiagnostics reports whether new global diagnostics have been
// accumulated since the last call, and resets the flag.
func (p *checkerPool) TakeNewGlobalDiagnostics() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	changed := p.globalDiagChanged
	p.globalDiagChanged = false
	return changed
}

// Discard signals that this pool's program has been replaced by a newer
// version. The pool remains fully functional (old snapshot handlers can
// still use it) but stops caching idle checkers, disposing them immediately
// after each use. This prevents checker accumulation during rapid typing,
// where each keystroke can produce a new program and pool.
func (p *checkerPool) Discard() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.opts.IdleTimeout == 0 {
		return // already discarded
	}
	p.log("checkerpool: Discarding pool, disposing idle checkers")
	p.opts.IdleTimeout = 0
	if p.cleanupTimer != nil {
		p.cleanupTimer.Stop()
		p.cleanupTimer = nil
	}
	// Dispose all currently idle checkers.
	for i, c := range p.checkers {
		if c != nil && p.heldBy[i] == "" {
			p.disposeCheckerLocked(i, c)
		}
	}
}

func noop() {}
