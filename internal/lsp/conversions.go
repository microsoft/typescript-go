package lsp

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func languageKindToScriptKind(languageID lsproto.LanguageKind) core.ScriptKind {
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

func lineAndCharacterToPosition(lineAndCharacter lsproto.Position, lineMap []core.TextPos) int {
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

func documentUriToFileName(uri lsproto.DocumentUri) string {
	uriStr := string(uri)
	if strings.HasPrefix(uriStr, "file:///") {
		path := uriStr[7:]
		if len(path) >= 4 {
			if nextSlash := strings.IndexByte(path[1:], '/'); nextSlash != -1 {
				if possibleDrive, _ := url.PathUnescape(path[1 : nextSlash+2]); strings.HasSuffix(possibleDrive, ":/") {
					return possibleDrive + path[len(possibleDrive)+3:]
				}
			}
		}
		return path
	}
	if strings.HasPrefix(uriStr, "file://") {
		// UNC path
		return uriStr[5:]
	}
	parsed := core.Must(url.Parse(uriStr))
	authority := parsed.Host
	if authority == "" {
		authority = "ts-nul-authority"
	}
	path := parsed.Path
	if path == "" {
		path = parsed.Opaque
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fragment := parsed.Fragment
	if fragment != "" {
		fragment = "#" + fragment
	}
	return fmt.Sprintf("^/%s/%s%s%s", parsed.Scheme, authority, path, fragment)
}
