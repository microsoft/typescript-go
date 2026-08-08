package core

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// GoroutinePool runs functions on a fixed set of long-lived workers so that
// deeply recursive, stack-growing work does not allocate a fresh grown stack
// per submitted task. Use it only when bounding peak stack memory matters; for
// plain fan-out, prefer WorkGroup or sync.WaitGroup.
//
// Submission blocks while all workers are busy. Submitting from a function
// running on the same pool will deadlock.
type GoroutinePool interface {
	// Run submits fn and blocks until it returns. Panics in fn re-panic on the caller.
	Run(fn func())
	// Queue submits fn without waiting on worker-backed pools. The single-threaded
	// immediatePool implementation runs fn inline on the caller instead. Panics in
	// fn crash the program.
	Queue(fn func())
	// Close drains submitted work and stops the workers. Must be called once;
	// after it returns, Run and Queue panic.
	Close()
}

// NewGOMAXPROCSPool returns a pool with GOMAXPROCS workers, or a degenerate
// pool that runs everything on the caller when singleThreaded is true.
func NewGOMAXPROCSPool(singleThreaded bool) GoroutinePool {
	if singleThreaded {
		return &immediatePool{}
	}
	return newGoroutinePool(runtime.GOMAXPROCS(0))
}

type goroutinePool struct {
	work   chan func()
	wg     sync.WaitGroup
	closed atomic.Bool
}

func newGoroutinePool(size int) *goroutinePool {
	p := &goroutinePool{
		work: make(chan func()),
	}
	p.wg.Add(size)
	for range size {
		go func() {
			defer p.wg.Done()
			for fn := range p.work {
				fn()
			}
		}()
	}
	return p
}

func (p *goroutinePool) Run(fn func()) {
	if p.closed.Load() {
		panic("GoroutinePool: Run called after Close")
	}
	// Recover on the worker so we can re-panic on the caller's goroutine,
	// keeping Run transparent to deferred recovers.
	var pv any
	done := make(chan struct{})
	p.work <- func() {
		defer close(done)
		defer func() { pv = recover() }()
		fn()
	}
	<-done
	if pv != nil {
		panic(pv)
	}
}

func (p *goroutinePool) Queue(fn func()) {
	if p.closed.Load() {
		panic("GoroutinePool: Queue called after Close")
	}
	p.work <- fn
}

func (p *goroutinePool) Close() {
	// Set closed before close(p.work) so racing submitters see a clear panic
	// rather than the runtime's "send on closed channel".
	p.closed.Store(true)
	close(p.work)
	p.wg.Wait()
}

type immediatePool struct {
	closed atomic.Bool
}

var _ GoroutinePool = (*immediatePool)(nil)

func (p *immediatePool) Run(fn func()) {
	if p.closed.Load() {
		panic("GoroutinePool: Run called after Close")
	}
	fn()
}

func (p *immediatePool) Queue(fn func()) {
	if p.closed.Load() {
		panic("GoroutinePool: Queue called after Close")
	}
	fn()
}

func (p *immediatePool) Close() {
	p.closed.Store(true)
}
