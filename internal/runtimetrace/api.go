//nolint:depguard
package runtimetrace

import (
	"context"
	"runtime/trace"
	"sync/atomic"
)

// IsEnabled reports whether the runtime execution tracer is enabled. It is a
// thin wrapper over runtime/trace.IsEnabled. Hot-path callers should gate
// expensive payload construction (e.g. fmt.Sprintf, slice/string building)
// on this so the disabled path is allocation-free:
//
//	if runtimetrace.IsEnabled() {
//	    runtimetrace.LogSafef(ctx, "phase", "n=%d", expensive())
//	}
func IsEnabled() bool { return trace.IsEnabled() }

// noopEnd is the no-op closure returned by Region when tracing is off, so
// the fast path does not allocate a method-value closure per call.
func noopEnd() {}

// Region starts a region in the calling goroutine and returns a function that
// ends it. Designed to be used with defer:
//
//	defer runtimetrace.Region(ctx, "parse")()
//
// When the runtime tracer is disabled, Region returns a shared no-op closure
// without allocating, making the deferred call free other than the defer
// itself. When enabled, Region calls trace.StartRegion and returns its End
// method value (a small closure allocation).
func Region(ctx context.Context, name string) func() {
	if !trace.IsEnabled() {
		return noopEnd
	}
	return trace.StartRegion(ctx, name).End
}

// NewTask creates a new task and returns a derived context that carries it
// along with a function that ends the task. Designed to be used with defer:
//
//	ctx, end := runtimetrace.NewTask(ctx, "lsp.textDocument/completion")
//	defer end()
//
// Tasks are higher-level than regions: they group regions and log events
// across goroutines and produce a latency entry in the trace's task table.
// Note that runtime/trace.NewTask is not gated on IsEnabled (tasks may span
// trace enable/disable boundaries), so this always allocates a Task.
func NewTask(ctx context.Context, name string) (context.Context, func()) {
	ctx, task := trace.NewTask(ctx, name)
	return ctx, task.End
}

// --- Logging --------------------------------------------------------------
//
// Two logging variants are provided to make it explicit at the call site
// whether the payload is safe to share:
//
//   - LogSafe / LogSafef: for payloads known to contain no user data.
//     Examples: counts, sizes, durations, protocol method names, JSON-RPC
//     request IDs, enum values. Always emitted when tracing is on.
//
//   - LogUnsafe / LogUnsafef: for payloads that may contain user data such
//     as file paths, module specifiers, or identifier names. Only emitted
//     when the user has opted in via TS_GO_RUNTIME_TRACE_DETAIL.
//
// All four helpers fast-path with a single IsEnabled / unsafeLogging check
// before doing anything else, so calls on the disabled path are cheap.
// However, callers using *f variants still pay the cost of boxing variadic
// arguments before the call. Hot-path callers should additionally gate on
// IsEnabled to skip the call entirely.

// unsafeLogging is set by Start when TS_GO_RUNTIME_TRACE_DETAIL is truthy.
var unsafeLogging atomic.Bool

// LogSafe emits a one-off event to the execution trace, attached to the task
// in ctx (if any). The payload must not contain user data.
func LogSafe(ctx context.Context, category, message string) {
	if trace.IsEnabled() {
		trace.Log(ctx, category, message)
	}
}

// LogSafef is like LogSafe but formats the message with fmt.Sprintf-style
// arguments. The formatted payload must not contain user data.
func LogSafef(ctx context.Context, category, format string, args ...any) {
	if trace.IsEnabled() {
		trace.Logf(ctx, category, format, args...)
	}
}

// LogUnsafe emits a one-off event to the execution trace only when the user
// has opted in via TS_GO_RUNTIME_TRACE_DETAIL. Use it for payloads that may
// include file paths, identifier names, or other user data.
func LogUnsafe(ctx context.Context, category, message string) {
	if unsafeLogging.Load() && trace.IsEnabled() {
		trace.Log(ctx, category, message)
	}
}

// LogUnsafef is like LogUnsafe but formats the message with fmt.Sprintf-style
// arguments.
func LogUnsafef(ctx context.Context, category, format string, args ...any) {
	if unsafeLogging.Load() && trace.IsEnabled() {
		trace.Logf(ctx, category, format, args...)
	}
}
