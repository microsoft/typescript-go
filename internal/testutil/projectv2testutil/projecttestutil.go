package projectv2testutil

import (
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

const (
	TestTypingsLocation = "/home/src/Library/Caches/typescript"
)

func Setup(files map[string]any) *projectv2.Session {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	session := projectv2.NewSession(projectv2.SessionOptions{
		CurrentDirectory:   "/",
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    TestTypingsLocation,
		PositionEncoding:   lsproto.PositionEncodingKindUTF8,
		WatchEnabled:       false,
		LoggingEnabled:     true,
	}, fs)
	return session
}
