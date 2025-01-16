package ast

import (
	"sync"
)

// ReentrantLock is a lock which allows the same ref to lock multiple times.
type ReentrantLock struct {
	mu     sync.Mutex
	ref    any
	count  int
	notify chan struct{}
}

// Lock locks the lock for the given ref. If the lock is already held by the lock
// the lock count is incremented. If the lock is held by a different ref, Lock
// waits until the lock is released and retries.
func (l *ReentrantLock) Lock(ref any) {
	for {
		notify := l.tryLock(ref)
		if notify == nil {
			return
		}
		<-notify
	}
}

func (l *ReentrantLock) tryLock(ref any) <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ref == ref {
		// Already locked.
		l.count++
		return nil
	}
	if l.ref == nil {
		// New locker
		l.ref = ref
		l.notify = make(chan struct{})
		return nil
	}
	// Locker from a different ref; wait for unlock and retry.
	return l.notify
}

// Unlock unlocks the lock for the given ref. If the lock is held multiple times,
// the lock count is decremented. If the ref is not the current locker, Unlock panics.
func (l *ReentrantLock) Unlock(ref any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ref != ref {
		panic("ReentrantLock: unlock with the wrong ref")
	}
	l.count--
	if l.count > 0 {
		return
	}
	close(l.notify)
	l.notify = nil
	l.ref = nil
}
