package projectv2testutil

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

//go:generate go tool github.com/matryer/moq -stub -fmt goimports -pkg projectv2testutil -out clientmock_generated.go ../../projectv2 Client
//go:generate go tool mvdan.cc/gofumpt -lang=go1.24 -w clientmock_generated.go

const (
	TestTypingsLocation = "/home/src/Library/Caches/typescript"
)

type SessionUtils struct {
	fs     vfs.FS
	client *ClientMock
}

func (h *SessionUtils) Client() *ClientMock {
	return h.client
}

func (h *SessionUtils) ExpectWatchFilesCalls(count int) func(t *testing.T) {
	var actualCalls atomic.Int32
	var wg sync.WaitGroup
	wg.Add(count)
	saveFunc := h.client.WatchFilesFunc
	h.client.WatchFilesFunc = func(_ context.Context, id projectv2.WatcherID, _ []*lsproto.FileSystemWatcher) error {
		actualCalls.Add(1)
		wg.Done()
		return nil
	}
	return func(t *testing.T) {
		t.Helper()
		wg.Wait()
		assert.Equal(t, actualCalls.Load(), int32(count))
		h.client.WatchFilesFunc = saveFunc
	}
}

func (h *SessionUtils) ExpectUnwatchFilesCalls(count int) func(t *testing.T) {
	var actualCalls atomic.Int32
	var wg sync.WaitGroup
	wg.Add(count)
	saveFunc := h.client.UnwatchFilesFunc
	h.client.UnwatchFilesFunc = func(_ context.Context, id projectv2.WatcherID) error {
		actualCalls.Add(1)
		wg.Done()
		return nil
	}
	return func(t *testing.T) {
		t.Helper()
		wg.Wait()
		assert.Equal(t, actualCalls.Load(), int32(count))
		h.client.UnwatchFilesFunc = saveFunc
	}
}

func (h *SessionUtils) FS() vfs.FS {
	return h.fs
}

func Setup(files map[string]any) (*projectv2.Session, *SessionUtils) {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	clientMock := &ClientMock{}
	sessionHandle := &SessionUtils{
		fs:     fs,
		client: clientMock,
	}

	session := projectv2.NewSession(projectv2.SessionOptions{
		CurrentDirectory:   "/",
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    TestTypingsLocation,
		PositionEncoding:   lsproto.PositionEncodingKindUTF8,
		WatchEnabled:       true,
		LoggingEnabled:     true,
	}, fs, clientMock)

	return session, sessionHandle
}
