package lsp

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func LanguageIDToScriptKind(languageID lsproto.LanguageKind) core.ScriptKind {
	switch languageID {
	case "typescript":
		return core.ScriptKindTS
	case "typescriptreact":
		return core.ScriptKindTSX
	case "javascript":
		return core.ScriptKindJS
	case "javascriptreact":
		return core.ScriptKindJSX
	default:
		return core.ScriptKindUnknown
	}
}

func LineAndCharacterToPosition(lineAndCharacter lsproto.Position, lineMap []core.TextPos) int {
	line := int(lineAndCharacter.Line)
	offset := int(lineAndCharacter.Character)

	if line < 0 || line >= len(lineMap) {
		panic(fmt.Sprintf("Bad line number. Line: %d, lineMap length: %d", line, len(lineMap)))
	}

	res := int(lineMap[line]) + offset
	if line < len(lineMap)-1 && res >= int(lineMap[line+1]) {
		panic("resulting position is out of bounds")
	}
	return res
}
