package lsp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type shutdownTestReader struct{}

func (shutdownTestReader) Read() (*lsproto.Message, error) { return nil, io.EOF }

type shutdownTestWriter struct{}

func (shutdownTestWriter) Write(*lsproto.Message) error { return nil }

// TestServerShutdownNoDeadlock verifies that operations after shutdown
// don't block.
func TestServerShutdownNoDeadlock(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/test/tsconfig.json": "{}",
		"/test/index.ts":      "const x = 1;",
	}, false))

	server := NewServer(&ServerOptions{
		In:                 shutdownTestReader{},
		Out:                shutdownTestWriter{},
		Err:                io.Discard,
		Cwd:                "/test",
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
	})

	ctx, cancel := context.WithCancel(context.Background())
	server.backgroundCtx = ctx

	// Start write loop to drain queue
	writeLoopDone := make(chan struct{})
	go func() {
		_ = server.writeLoop(ctx)
		close(writeLoopDone)
	}()

	// Create session with the server's lifecycle context
	server.initStarted.Store(true)
	server.session = project.NewSession(&project.SessionInit{
		BackgroundCtx: ctx,
		Options: &project.SessionOptions{
			CurrentDirectory:   "/test",
			DefaultLibraryPath: bundled.LibPath(),
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       false,
			LoggingEnabled:     true,
		},
		FS:     fs,
		Logger: server.logger,
	})

	// Open a file to establish a project
	server.session.DidOpenFile(ctx, "file:///test/index.ts", 1, "const x = 1;", lsproto.LanguageKindTypeScript)
	server.session.WaitForBackgroundTasks()

	// Shutdown (cancel context and wait for write loop to exit)
	cancel()
	<-writeLoopDone

	// Fill the queue so any logging attempt would block
	dummyMsg := lsproto.WindowLogMessageInfo.NewNotificationMessage(&lsproto.LogMessageParams{
		Type:    lsproto.MessageTypeInfo,
		Message: "fill",
	}).Message()

	for range cap(server.outgoingQueue) {
		select {
		case server.outgoingQueue <- dummyMsg:
			// filled one slot
		default:
			// queue full
		}
	}

	// Trigger operations that would log (these should not block)
	server.session.DidChangeFile(ctx, "file:///test/index.ts", 2, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
		{
			WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{
				Text: "const x = 2;",
			},
		},
	})
	_, _ = server.session.GetLanguageService(ctx, "file:///test/index.ts")
	server.session.WaitForBackgroundTasks()

	server.session.Close()
}

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

	// Construct a JSON-RPC request with a position containing a number
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

	// Find the error response.
	var foundError bool
	for _, msg := range writer.messages {
		if msg.Kind != jsonrpc.MessageKindResponse {
			continue
		}
		resp := msg.AsResponse()
		if resp.Error != nil {
			foundError = true
			assert.Assert(t, resp.Error.Code == int32(lsproto.ErrorCodeInvalidParams),
				"expected error code %d, got %d", int32(lsproto.ErrorCodeInvalidParams), resp.Error.Code)
			assert.Assert(t, resp.ID != nil, "expected error response to have an ID")
			break
		}
	}
	assert.Assert(t, foundError, "expected an error response for the invalid position")
}

func TestServerInvalidNotificationParams(t *testing.T) {
	t.Parallel()

	// Construct a JSON-RPC notification with a position containing a number
	// too large to fit in a uint32 (Position.Line and Position.Character are uint32).
	// The value 99999999999999999999 exceeds both uint32 and int64 range,
	// so json.Unmarshal into uint32 should fail.
	body := `{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///test.ts","languageId":"typescript","version":99999999999999999999,"text":"const x = 1;"}}}`

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

	// The server should not have sent any error response for an invalid notification.
	writer.mu.Lock()
	defer writer.mu.Unlock()

	// Find any error response.
	var foundError bool
	for _, msg := range writer.messages {
		if msg.Kind != jsonrpc.MessageKindResponse {
			continue
		}
		resp := msg.AsResponse()
		if resp.Error != nil {
			foundError = true
			break
		}
	}
	assert.Assert(t, !foundError, "expected no error response for invalid notification")
}
