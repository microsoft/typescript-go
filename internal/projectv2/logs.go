package projectv2

import (
	"context"
	"fmt"
	"strings"
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
		if d.closed {
			return
		}
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

func newLog(child *logCollector, message string) *log {
	return &log{
		seq:     seq.Add(1),
		time:    time.Now(),
		message: message,
		child:   child,
	}
}

type logCollector struct {
	name       string
	logs       []*log
	dispatcher *dispatcher
	root       *logCollector
	level      int

	// Only set on root
	count        atomic.Int32
	stringLength atomic.Int32
	close        func()
}

func NewLogCollector(name string) *logCollector {
	dispatcher, close := newDispatcher()
	lc := &logCollector{
		name:       name,
		dispatcher: dispatcher,
		close:      close,
	}
	lc.root = lc
	return lc
}

func (c *logCollector) add(log *log) {
	// indent + header + message + newline
	c.root.stringLength.Add(int32(c.level + 15 + len(log.message) + 1))
	c.root.count.Add(1)
	c.dispatcher.Dispatch(func() {
		c.logs = append(c.logs, log)
	})
}

func (c *logCollector) Log(message ...any) {
	if c == nil {
		return
	}
	log := newLog(nil, fmt.Sprint(message...))
	c.add(log)
}

func (c *logCollector) Logf(format string, args ...any) {
	if c == nil {
		return
	}
	log := newLog(nil, fmt.Sprintf(format, args...))
	c.add(log)
}

func (c *logCollector) Embed(logs *logCollector) {
	logs.Close()
	count := logs.count.Load()
	c.root.stringLength.Add(logs.stringLength.Load() + count*int32(c.level))
	c.root.count.Add(count)
	log := newLog(logs, logs.name)
	c.add(log)
}

func (c *logCollector) Fork(message string) *logCollector {
	if c == nil {
		return nil
	}
	child := &logCollector{dispatcher: c.dispatcher, level: c.level + 1, root: c.root}
	log := newLog(child, message)
	c.add(log)
	return child
}

func (c *logCollector) Close() {
	if c == nil {
		return
	}
	c.close()
}

type Logger interface {
	Log(msg ...any)
}

func (c *logCollector) String() string {
	if c.root != c {
		panic("can only call String on root logCollector")
	}
	c.Close()
	var builder strings.Builder
	header := fmt.Sprintf("======== %s ========\n", c.name)
	builder.Grow(int(c.stringLength.Load()) + len(header))
	builder.WriteString(header)
	c.writeLogsRecursive(&builder, "")
	return builder.String()
}

func (c *logCollector) writeLogsRecursive(builder *strings.Builder, indent string) {
	for _, log := range c.logs {
		builder.WriteString(indent)
		builder.WriteString("[")
		builder.WriteString(log.time.Format("15:04:05.000"))
		builder.WriteString("] ")
		builder.WriteString(log.message)
		builder.WriteString("\n")
		if log.child != nil {
			log.child.writeLogsRecursive(builder, indent+"\t")
		}
	}
}
