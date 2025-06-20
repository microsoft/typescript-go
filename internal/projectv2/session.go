package projectv2

import "github.com/microsoft/typescript-go/internal/lsp/lsproto"

type SessionOptions struct {
	DefaultLibraryPath string
	TypingsLocation    string
	PositionEncoding   lsproto.PositionEncodingKind
	WatchEnabled       bool
}
