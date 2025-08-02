package project

import (
	"context"
	"sync"
)

// BackgroundTask represents a task that can be executed asynchronously
type BackgroundTask func(ctx context.Context)

// BackgroundQueue manages background tasks execution
type BackgroundQueue struct {
	wg     sync.WaitGroup
	mu     sync.RWMutex
	closed bool
}

func newBackgroundQueue() *BackgroundQueue {
	return &BackgroundQueue{}
}

func (q *BackgroundQueue) Enqueue(ctx context.Context, task BackgroundTask) {
	q.mu.RLock()
	if q.closed {
		q.mu.RUnlock()
		return
	}

	q.wg.Add(1)
	q.mu.RUnlock()

	go func() {
		defer q.wg.Done()
		task(ctx)
	}()
}

// WaitForEmpty waits for all active tasks to complete.
func (q *BackgroundQueue) WaitForEmpty() {
	q.wg.Wait()
}

func (q *BackgroundQueue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
}
