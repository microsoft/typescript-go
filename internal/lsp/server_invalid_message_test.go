package lsp

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

type collectingWriter struct {
	mu       sync.Mutex
	messages []*lsproto.Message
}

func (w *collectingWriter) Write(msg *lsproto.Message) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.messages = append(w.messages, msg)
	return nil
}

func TestServerInvalidRequestParams(t *testing.T) {
	t.Parallel()

	// Construct a JSON-RPC message with a position containing a number
	// too large to fit in a uint32 (Position.Line and Position.Character are uint32).
	// The value 99999999999999999999 exceeds both uint32 and int64 range,
	// so json.Unmarshal into uint32 should fail.
	body := `{"jsonrpc":"2.0","id":1,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///test.ts"},"position":{"line":99999999999999999999,"character":0}}}`

	rawMessage := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)

	reader := ToReader(io.NopCloser(bytes.NewReader([]byte(rawMessage))))
	writer := &collectingWriter{}

	server := NewServer(&ServerOptions{
		In:  reader,
		Out: writer,
		Err: io.Discard,
		Cwd: "/test",
	})

	err := server.Run(t.Context())
	assert.NilError(t, err)

	// The server should have sent an error response for the invalid request.
	writer.mu.Lock()
	defer writer.mu.Unlock()

	assert.Assert(t, len(writer.messages) > 0, "expected at least one error response message")

	// Find the error response.
	var foundError bool
	for _, msg := range writer.messages {
		resp := msg.AsResponse()
		if resp != nil && resp.Error != nil {
			foundError = true
			assert.Assert(t, resp.Error.Code == int32(lsproto.ErrorCodeInvalidParams),
				"expected error code %d, got %d", int32(lsproto.ErrorCodeInvalidParams), resp.Error.Code)
			assert.Assert(t, resp.ID != nil, "expected error response to have an ID")
			break
		}
	}
	assert.Assert(t, foundError, "expected an error response for the invalid position")
}
