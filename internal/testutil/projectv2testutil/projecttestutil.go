package projectv2testutil

import (
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
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
