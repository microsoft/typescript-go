package projectv2

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type SessionOptions struct {
	DefaultLibraryPath string
	TypingsLocation    string
	PositionEncoding   lsproto.PositionEncodingKind
	WatchEnabled       bool
	CurrentDirectory   string
	NewLine            string
}

type Session struct {
	options    SessionOptions
	fs         *overlayFS
	parseCache *parseCache
	snapshot   *Snapshot
}

func (s *Session) DidOpenFile(ctx context.Context, uri lsproto.DocumentUri, version int32, content string, languageKind lsproto.LanguageKind) {
	s.fs.replaceOverlay(uri, content, version, ls.LanguageKindToScriptKind(languageKind))

}
