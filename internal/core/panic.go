package core

import (
	"bytes"
	"fmt"
	"runtime/debug"
)

// PanicWithStack wraps a recovered panic value along with the original stack trace
// captured at the site of the panic. This allows re-panicking while preserving the
// original stack context, rather than losing it at the re-panic site.
type PanicWithStack struct {
	Value any
	Stack []byte
}

func (p *PanicWithStack) String() string {
	return fmt.Sprintf("%v\n%s", p.Value, string(p.Stack))
}

// NewPanicWithStack creates a PanicWithStack from a recovered panic value.
// It captures the current stack trace via debug.Stack() and strips the
// recovery infrastructure frames (debug.Stack, the deferred recovery function,
// and the panic runtime frame) so that only the actual crash site frames remain.
func NewPanicWithStack(recovered any) *PanicWithStack {
	return &PanicWithStack{
		Value: recovered,
		Stack: trimPanicRecoveryFrames(debug.Stack()),
	}
}

// trimPanicRecoveryFrames strips recovery infrastructure frames from a stack
// trace captured by debug.Stack() inside a deferred recovery function.
// It removes everything up to and including the "panic(...)" frame and its
// file/line pair, leaving only the frames from the actual crash site onward.
func trimPanicRecoveryFrames(stack []byte) []byte {
	// Find the panic() frame. In Go stack traces, it appears as a line
	// starting with "panic(" (possibly with leading whitespace stripped).
	lines := bytes.Split(stack, []byte("\n"))
	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, []byte("panic(")) {
			// The panic frame is followed by its file/line pair on the next line.
			// Skip both to get to the actual crash site.
			startIdx := i + 2
			if startIdx < len(lines) {
				return bytes.Join(lines[startIdx:], []byte("\n"))
			}
		}
	}
	// If we couldn't find the panic frame, return the original stack.
	return stack
}
