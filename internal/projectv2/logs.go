package projectv2

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var seq atomic.Uint64

type dispatcher struct {
	closed bool
	ch     chan func()
}

func newDispatcher() (*dispatcher, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &dispatcher{
		ch: make(chan func(), 1024),
	}

	go func() {
		for {
			select {
			// Drain the queue before checking for cancellation to avoid dropping logs
			case fn := <-d.ch:
				fn()
			case <-ctx.Done():
				return
			}
		}
	}()

	return d, func() {
		done := make(chan struct{})
		d.Dispatch(func() {
			close(done)
		})
		<-done
		cancel()
		close(d.ch)
		d.closed = true
	}
}

func (d *dispatcher) Dispatch(fn func()) {
	if d.closed {
		panic("tried to log after logger was closed")
	}
	d.ch <- fn
}

type log struct {
	seq     uint64
	time    time.Time
	message string
	child   *logCollector
}

func newLog(child *logCollector, message string) log {
	return log{
		seq:     seq.Add(1),
		time:    time.Now(),
		message: message,
		child:   child,
	}
}

type logCollector struct {
	name       string
	logs       []log
	dispatcher *dispatcher
}

func NewLogCollector(name string) (*logCollector, func()) {
	dispatcher, close := newDispatcher()
	return &logCollector{
		name:       name,
		dispatcher: dispatcher,
	}, close
}

func (c *logCollector) Log(message string) {
	log := newLog(nil, message)
	c.dispatcher.Dispatch(func() {
		c.logs = append(c.logs, log)
	})
}

func (c *logCollector) Logf(format string, args ...any) {
	log := newLog(nil, fmt.Sprintf(format, args...))
	c.dispatcher.Dispatch(func() {
		c.logs = append(c.logs, log)
	})
}

func (c *logCollector) Fork(name string, message string) *logCollector {
	child := &logCollector{name: name, dispatcher: c.dispatcher}
	log := newLog(child, message)
	c.dispatcher.Dispatch(func() {
		c.logs = append(c.logs, log)
	})
	return child
}
