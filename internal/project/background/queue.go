package background

import (
	"context"
	"sync"
	"sync/atomic"
)

// Queue manages background tasks execution
type Queue struct {
	wg     sync.WaitGroup
	closed atomic.Bool
}

// NewQueue creates a new background queue for managing background tasks execution.
func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(ctx context.Context, fn func(context.Context)) {
	if q.closed.Load() {
		return
	}

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		fn(ctx)
	}()
}

// Wait waits for all active tasks to complete.
// It does not prevent new tasks from being enqueued while waiting.
func (q *Queue) Wait() {
	q.wg.Wait()
}

func (q *Queue) Close() {
	q.closed.Store(true)
}
