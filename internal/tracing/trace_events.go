package tracing

import (
	"fmt"
	"maps"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// Phase represents the phase of compilation for trace events
type Phase string

const (
	PhaseParse      Phase = "parse"
	PhaseProgram    Phase = "program"
	PhaseBind       Phase = "bind"
	PhaseCheck      Phase = "check"      // Before we get into checking types (e.g. checkSourceFile)
	PhaseCheckTypes Phase = "checkTypes" // Type checking operations
	PhaseEmit       Phase = "emit"
	PhaseSession    Phase = "session"
)

// sampleInterval is 10ms - we sample events that cross this boundary
const sampleInterval = 10 * time.Millisecond

// Clock provides timestamps for tracing. For testing, use a deterministic clock.
type Clock interface {
	Now() time.Time
}

// realClock uses the system clock
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// MonotonicClock provides a deterministic clock that increments by a fixed amount each call.
// Useful for testing to get reproducible timestamps.
type MonotonicClock struct {
	counter atomic.Int64
	step    int64 // microseconds per call
}

// NewMonotonicClock creates a clock that increments by step microseconds each call.
func NewMonotonicClock(stepMicros int64) *MonotonicClock {
	return &MonotonicClock{step: stepMicros}
}

func (c *MonotonicClock) Now() time.Time {
	micros := c.counter.Add(c.step)
	return time.UnixMicro(micros)
}

// EventTracer handles Chrome DevTools trace format events
type EventTracer struct {
	fs         vfs.FS
	tracePath  string
	clock      Clock
	startTime  time.Time
	events     []string
	eventStack []stackEvent
	mu         sync.Mutex
}

type stackEvent struct {
	phase               Phase
	name                string
	args                map[string]any
	time                time.Duration
	separateBeginAndEnd bool
}

// NewEventTracer creates a new event tracer
func NewEventTracer(fs vfs.FS, traceDir string, traceSuffix string) (*EventTracer, error) {
	return NewEventTracerWithClock(fs, traceDir, traceSuffix, realClock{})
}

// NewEventTracerWithClock creates a new event tracer with a custom clock.
// Use NewMonotonicClock for deterministic timestamps in tests.
func NewEventTracerWithClock(fs vfs.FS, traceDir string, traceSuffix string, clock Clock) (*EventTracer, error) {
	tracePath := tspath.CombinePaths(traceDir, fmt.Sprintf("trace%s.json", traceSuffix))

	now := clock.Now()
	et := &EventTracer{
		fs:         fs,
		tracePath:  tracePath,
		clock:      clock,
		startTime:  now,
		events:     []string{},
		eventStack: []stackEvent{},
	}

	// Write initial metadata events
	ts := et.timestamp()
	meta := fmt.Sprintf(`{"cat":"__metadata","ph":"M","ts":%d,"pid":1,"tid":1`, ts)
	et.events = append(et.events,
		meta+`,"name":"process_name","args":{"name":"tsgo"}}`,
		meta+`,"name":"thread_name","args":{"name":"Main"}}`,
		meta+`,"name":"TracingStartedInBrowser","cat":"disabled-by-default-devtools.timeline"}`,
	)

	return et, nil
}

// TracePath returns the path to the trace file
func (et *EventTracer) TracePath() string {
	return et.tracePath
}

// timestamp returns microseconds since start
func (et *EventTracer) timestamp() int64 {
	return et.clock.Now().Sub(et.startTime).Microseconds()
}

// Instant writes an instant event (single point in time)
func (et *EventTracer) Instant(phase Phase, name string, args map[string]any) {
	et.writeEvent("I", phase, name, args, `"s":"g"`, et.timestamp())
}

// Push starts a new trace event that will be closed by Pop
// separateBeginAndEnd is used for special cases where we need the trace point even if the event
// never terminates (typically for reducing a scenario too big to trace to one that can be completed).
func (et *EventTracer) Push(phase Phase, name string, args map[string]any, separateBeginAndEnd bool) {
	et.mu.Lock()
	defer et.mu.Unlock()

	ts := et.timestamp()
	if separateBeginAndEnd {
		et.writeEventLocked("B", phase, name, args, "", ts)
	}
	et.eventStack = append(et.eventStack, stackEvent{
		phase:               phase,
		name:                name,
		args:                args,
		time:                time.Duration(ts) * time.Microsecond,
		separateBeginAndEnd: separateBeginAndEnd,
	})
}

// Pop closes the most recent Push event
func (et *EventTracer) Pop(results map[string]any) {
	et.mu.Lock()
	defer et.mu.Unlock()

	if len(et.eventStack) == 0 {
		return
	}

	endTime := time.Duration(et.timestamp()) * time.Microsecond
	et.writeStackEventLocked(len(et.eventStack)-1, endTime, results)
	et.eventStack = et.eventStack[:len(et.eventStack)-1]
}

// PopAll closes all open Push events
func (et *EventTracer) PopAll() {
	et.mu.Lock()
	defer et.mu.Unlock()

	endTime := time.Duration(et.timestamp()) * time.Microsecond
	for i := len(et.eventStack) - 1; i >= 0; i-- {
		et.writeStackEventLocked(i, endTime, nil)
	}
	et.eventStack = et.eventStack[:0]
}

func (et *EventTracer) writeStackEventLocked(index int, endTime time.Duration, results map[string]any) {
	ev := et.eventStack[index]
	if ev.separateBeginAndEnd {
		// For separateBeginAndEnd events, write an End event
		et.writeEventLocked("E", ev.phase, ev.name, ev.args, "", int64(endTime/time.Microsecond))
	} else {
		// Test if [time, endTime) straddles a sampling point
		startMicros := int64(ev.time / time.Microsecond)
		endMicros := int64(endTime / time.Microsecond)
		sampleMicros := int64(sampleInterval / time.Microsecond)

		if sampleMicros-(startMicros%sampleMicros) <= endMicros-startMicros {
			// Merge args and results
			mergedArgs := make(map[string]any)
			maps.Copy(mergedArgs, ev.args)
			if results != nil {
				mergedArgs["results"] = results
			}

			dur := endMicros - startMicros
			et.writeEventLocked("X", ev.phase, ev.name, mergedArgs, fmt.Sprintf(`"dur":%d`, dur), startMicros)
		}
	}
}

func (et *EventTracer) writeEvent(eventType string, phase Phase, name string, args map[string]any, extras string, ts int64) {
	et.mu.Lock()
	defer et.mu.Unlock()
	et.writeEventLocked(eventType, phase, name, args, extras, ts)
}

func (et *EventTracer) writeEventLocked(eventType string, phase Phase, name string, args map[string]any, extras string, ts int64) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`{"pid":1,"tid":1,"ph":"%s","cat":"%s","ts":%d,"name":"%s"`, eventType, phase, ts, name))

	if extras != "" {
		sb.WriteString(",")
		sb.WriteString(extras)
	}

	if len(args) > 0 {
		argsJSON, err := json.Marshal(args)
		if err == nil {
			sb.WriteString(`,"args":`)
			sb.Write(argsJSON)
		}
	}

	sb.WriteString("}")
	et.events = append(et.events, sb.String())
}

// Finish writes all events to the trace file
func (et *EventTracer) Finish() error {
	et.mu.Lock()
	defer et.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, event := range et.events {
		sb.WriteString(event)
		if i < len(et.events)-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n]\n")

	return et.fs.WriteFile(et.tracePath, sb.String(), false)
}
