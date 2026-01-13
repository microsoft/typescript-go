package lsp

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type shutdownTestReader struct{}

func (shutdownTestReader) Read() (*lsproto.Message, error) { return nil, io.EOF }

type shutdownTestWriter struct{}

func (shutdownTestWriter) Write(*lsproto.Message) error { return nil }

func newShutdownTestServer(t *testing.T) *Server {
	t.Helper()
	return NewServer(&ServerOptions{
		In:  shutdownTestReader{},
		Out: shutdownTestWriter{},
		Err: io.Discard,
		Cwd: "/",
		FS:  vfstest.FromMap(map[string]string{}, true),
	})
}

// Documents the deadlock observed when the server keeps logging after
// the write loop has exited during shutdown. Currently fails because the
// outgoing queue fills and Logger.Info blocks forever.
func TestLoggerBlocksAfterShutdown(t *testing.T) {
	t.Parallel()

	server := newShutdownTestServer(t)
	server.initStarted.Store(true)

	ctx, cancel := context.WithCancel(context.Background())
	writeLoopDone := make(chan struct{})
	go func() {
		_ = server.writeLoop(ctx)
		close(writeLoopDone)
	}()

	cancel()
	<-writeLoopDone

	msg := lsproto.WindowLogMessageInfo.NewNotificationMessage(&lsproto.LogMessageParams{
		Type:    lsproto.MessageTypeInfo,
		Message: "pre-shutdown",
	}).Message()

	for i := 0; i < cap(server.outgoingQueue); i++ {
		select {
		case server.outgoingQueue <- msg:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timeout filling outgoing queue at %d", i)
		}
	}

	logDone := make(chan struct{})
	go func() {
		server.logger.Info("log after shutdown")
		close(logDone)
	}()

	select {
	case <-logDone:
		// Expected once the shutdown handling is fixed.
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("log send blocked after shutdown")
	}
}
