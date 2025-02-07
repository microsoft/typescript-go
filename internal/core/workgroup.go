package core

import (
	"sync"
	"sync/atomic"
)

type WorkGroup interface {
	Queue(fn func())
	RunAndWait()
}

func NewWorkGroup(singleThreaded bool) WorkGroup {
	if singleThreaded {
		return &singleThreadedWorkGroup{}
	}
	return &parallelWorkGroup{}
}

type parallelWorkGroup struct {
	done atomic.Bool
	wg   sync.WaitGroup
}

var _ WorkGroup = (*parallelWorkGroup)(nil)

func (w *parallelWorkGroup) Queue(fn func()) {
	if w.done.Load() {
		panic("Queue called after Wait returned")
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		fn()
	}()
}

func (w *parallelWorkGroup) RunAndWait() {
	defer w.done.Store(true)
	w.wg.Wait()
}

type singleThreadedWorkGroup struct {
	done  atomic.Bool
	fnsMu sync.Mutex
	fns   []func()
}

var _ WorkGroup = (*singleThreadedWorkGroup)(nil)

func (w *singleThreadedWorkGroup) Queue(fn func()) {
	if w.done.Load() {
		panic("Queue called after Wait returned")
	}

	w.fnsMu.Lock()
	defer w.fnsMu.Unlock()
	w.fns = append(w.fns, fn)
}

func (w *singleThreadedWorkGroup) RunAndWait() {
	defer w.done.Store(true)
	for {
		fn := w.pop()
		if fn == nil {
			return
		}
		fn()
	}
}

func (w *singleThreadedWorkGroup) pop() func() {
	w.fnsMu.Lock()
	defer w.fnsMu.Unlock()
	if len(w.fns) == 0 {
		return nil
	}
	end := len(w.fns) - 1
	fn := w.fns[end]
	w.fns[end] = nil // Allow GC
	w.fns = w.fns[:end]
	return fn
}
