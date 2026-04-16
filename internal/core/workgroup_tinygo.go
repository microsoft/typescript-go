//go:build tinygo

package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type WorkGroup interface {
	Queue(fn func())
	RunAndWait()
}

func NewWorkGroup(_ bool) WorkGroup {
	return &singleThreadedWorkGroup{}
}

type singleThreadedWorkGroup struct {
	done  atomic.Bool
	fnsMu sync.Mutex
	fns   []func()
}

var _ WorkGroup = (*singleThreadedWorkGroup)(nil)

func (w *singleThreadedWorkGroup) Queue(fn func()) {
	if w.done.Load() {
		panic("Queue called after RunAndWait returned")
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
	w.fns[end] = nil
	w.fns = w.fns[:end]
	return fn
}

// ThrottleGroup runs functions sequentially under tinygo.
type ThrottleGroup struct {
	fns []func() error
}

func NewThrottleGroup(_ context.Context, _ chan struct{}) *ThrottleGroup {
	return &ThrottleGroup{}
}

func (tg *ThrottleGroup) Go(fn func() error) {
	tg.fns = append(tg.fns, fn)
}

func (tg *ThrottleGroup) Wait() error {
	for _, fn := range tg.fns {
		if err := fn(); err != nil {
			return err
		}
	}
	tg.fns = nil
	return nil
}

