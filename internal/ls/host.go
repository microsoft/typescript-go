package ls

import (
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"golang.org/x/text/language"
)

type Host interface {
	GetProgram() *compiler.Program
	GetPositionEncoding() lsproto.PositionEncodingKind
	GetLineMap(fileName string) *LineMap
	GetLocale() language.Tag
}
