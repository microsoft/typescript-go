package api

import (
	"math"
	"time"

	"github.com/microsoft/typescript-go/internal/jsonrpc"
)

// Message is an alias for jsonrpc.Message for convenience.
type Message = jsonrpc.Message

// serverTimedResult wraps a handler result together with the server-side
// processing time for the request. Protocols that support timing metadata emit
// the processing time alongside the response: the MessagePack protocol appends
// a fixed-size footer, and the JSON-RPC protocol wraps the result in an
// envelope. Protocols that don't support it simply unwrap and send the inner
// result.
type serverTimedResult struct {
	// result is the underlying handler result (may be RawBinary, nil, or any
	// JSON-serializable value).
	result any
	// processingTimeMicros is the wall-clock time the server spent handling the
	// request, in microseconds, saturated to the uint32 range.
	processingTimeMicros uint32
}

// timingStart returns the current time when timing collection is enabled, or
// the zero time otherwise. This avoids reading the clock on the request hot
// path when timing is disabled (the default).
func timingStart(enabled bool) time.Time {
	if enabled {
		return time.Now()
	}
	return time.Time{}
}

// maybeTimed wraps a successful handler result with the server's processing
// time when timing collection is enabled. Otherwise the result is returned
// unchanged. start marks the moment request handling began.
func maybeTimed(result any, start time.Time, enabled bool) any {
	if !enabled {
		return result
	}
	return serverTimedResult{
		result:               result,
		processingTimeMicros: durationToMicros(time.Since(start)),
	}
}

// durationToMicros converts a duration to whole microseconds, clamped to the
// uint32 range so it fits in the fixed-size timing metadata.
func durationToMicros(d time.Duration) uint32 {
	micros := d.Microseconds()
	if micros < 0 {
		return 0
	}
	if micros > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(micros)
}

// Protocol defines the interface for reading and writing API messages.
type Protocol interface {
	// ReadMessage reads the next message from the connection.
	ReadMessage() (*Message, error)
	// WriteRequest writes a request message.
	WriteRequest(id *jsonrpc.ID, method string, params any) error
	// WriteNotification writes a notification message (no ID).
	WriteNotification(method string, params any) error
	// WriteResponse writes a successful response.
	WriteResponse(id *jsonrpc.ID, result any) error
	// WriteError writes an error response.
	WriteError(id *jsonrpc.ID, err *jsonrpc.ResponseError) error
}
