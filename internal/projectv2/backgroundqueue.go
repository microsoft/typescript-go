package projectv2

import (
	"sync"
)

// BackgroundTask represents a task that can be executed asynchronously
type BackgroundTask func()

// BackgroundQueue manages background tasks execution
type BackgroundQueue struct {
	tasks chan BackgroundTask
	wg    sync.WaitGroup
	done  chan struct{}
}

func newBackgroundTaskQueue() *BackgroundQueue {
	queue := &BackgroundQueue{
		tasks: make(chan BackgroundTask, 10),
		done:  make(chan struct{}),
	}

	// Start the dispatcher goroutine
	go queue.dispatcher()
	return queue
}

func (q *BackgroundQueue) dispatcher() {
	for {
		select {
		case task := <-q.tasks:
			// Execute task in a new goroutine
			q.wg.Add(1)
			go func() {
				defer q.wg.Done()
				task()
			}()
		case <-q.done:
			return
		}
	}
}

func (q *BackgroundQueue) Enqueue(task BackgroundTask) {
	select {
	case q.tasks <- task:
	case <-q.done:
		// Queue is shutting down, don't enqueue
	}
}

// WaitForEmpty waits for all active tasks to complete.
func (q *BackgroundQueue) WaitForEmpty() {
	q.wg.Wait()
	for {
		if len(q.tasks) == 0 {
			break
		}
		q.wg.Wait()
	}
}

func (q *BackgroundQueue) Close() {
	close(q.done)
}
