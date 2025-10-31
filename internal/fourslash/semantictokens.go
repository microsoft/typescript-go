package fourslash

import (
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type SemanticToken struct {
	Type string
	Text string
}

func (f *FourslashTest) VerifySemanticTokens(t *testing.T, expected []SemanticToken) {
	t.Helper()

	// Get capabilities for semantic tokens
	tokenTypes := defaultTokenTypes()
	tokenModifiers := defaultTokenModifiers()

	trueVal := true
	caps := &lsproto.SemanticTokensClientCapabilities{
		Requests: &lsproto.ClientSemanticTokensRequestOptions{
			Full: &lsproto.BooleanOrClientSemanticTokensRequestFullDelta{
				Boolean: &trueVal,
			},
		},
		TokenTypes:     tokenTypes,
		TokenModifiers: tokenModifiers,
		Formats:        []lsproto.TokenFormat{lsproto.TokenFormatRelative},
	}

	params := &lsproto.SemanticTokensParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsconv.FileNameToDocumentURI(f.activeFilename),
		},
	}

	resMsg, result, resultOk := sendRequest(t, f, lsproto.TextDocumentSemanticTokensFullInfo, params)
	if resMsg == nil {
		t.Fatal("Nil response received for semantic tokens request")
	}
	if !resultOk {
		t.Fatalf("Unexpected response type for semantic tokens request: %T", resMsg.AsResponse().Result)
	}

	if result.SemanticTokens == nil {
		if len(expected) == 0 {
			return
		}
		t.Fatal("Expected semantic tokens but got nil")
	}

	// Decode the semantic tokens
	actual := decodeSemanticTokens(f, result.SemanticTokens.Data, caps)

	// Compare with expected
	if len(actual) != len(expected) {
		t.Fatalf("Expected %d semantic tokens, got %d\n\nExpected:\n%s\n\nActual:\n%s",
			len(expected), len(actual),
			formatSemanticTokens(expected),
			formatSemanticTokens(actual))
	}

	for i, exp := range expected {
		act := actual[i]
		if exp.Type != act.Type || exp.Text != act.Text {
			t.Errorf("Token %d mismatch:\n  Expected: {Type: %q, Text: %q}\n  Actual:   {Type: %q, Text: %q}",
				i, exp.Type, exp.Text, act.Type, act.Text)
		}
	}
}

func decodeSemanticTokens(f *FourslashTest, data []uint32, caps *lsproto.SemanticTokensClientCapabilities) []SemanticToken {
	if len(data)%5 != 0 {
		panic(fmt.Sprintf("Invalid semantic tokens data length: %d", len(data)))
	}

	scriptInfo := f.scriptInfos[f.activeFilename]
	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(_ string) *lsconv.LSPLineMap {
		return scriptInfo.lineMap
	})

	var tokens []SemanticToken
	prevLine := uint32(0)
	prevChar := uint32(0)

	for i := 0; i < len(data); i += 5 {
		deltaLine := data[i]
		deltaChar := data[i+1]
		length := data[i+2]
		tokenTypeIdx := data[i+3]
		tokenModifiers := data[i+4]

		// Calculate absolute position
		line := prevLine + deltaLine
		var char uint32
		if deltaLine == 0 {
			char = prevChar + deltaChar
		} else {
			char = deltaChar
		}

		// Get token type
		if int(tokenTypeIdx) >= len(caps.TokenTypes) {
			panic(fmt.Sprintf("Token type index out of range: %d", tokenTypeIdx))
		}
		tokenType := caps.TokenTypes[tokenTypeIdx]

		// Get modifiers
		var modifiers []string
		for i, mod := range caps.TokenModifiers {
			if tokenModifiers&(1<<i) != 0 {
				modifiers = append(modifiers, string(mod))
			}
		}

		// Build full type string (type.modifier1.modifier2)
		typeStr := string(tokenType)
		if len(modifiers) > 0 {
			typeStr = typeStr + "." + strings.Join(modifiers, ".")
		}

		// Get the text
		startPos := lsproto.Position{Line: line, Character: char}
		endPos := lsproto.Position{Line: line, Character: char + length}
		startOffset := int(converters.LineAndCharacterToPosition(scriptInfo, startPos))
		endOffset := int(converters.LineAndCharacterToPosition(scriptInfo, endPos))
		text := scriptInfo.content[startOffset:endOffset]

		tokens = append(tokens, SemanticToken{
			Type: typeStr,
			Text: text,
		})

		prevLine = line
		prevChar = char
	}

	return tokens
}

func formatSemanticTokens(tokens []SemanticToken) string {
	var lines []string
	for i, tok := range tokens {
		lines = append(lines, fmt.Sprintf("  [%d] {Type: %q, Text: %q}", i, tok.Type, tok.Text))
	}
	return strings.Join(lines, "\n")
}

func defaultTokenTypes() []string {
	return []string{
		string(lsproto.SemanticTokenTypesnamespace),
		string(lsproto.SemanticTokenTypesclass),
		string(lsproto.SemanticTokenTypesenum),
		string(lsproto.SemanticTokenTypesinterface),
		string(lsproto.SemanticTokenTypesstruct),
		string(lsproto.SemanticTokenTypestypeParameter),
		string(lsproto.SemanticTokenTypestype),
		string(lsproto.SemanticTokenTypesparameter),
		string(lsproto.SemanticTokenTypesvariable),
		string(lsproto.SemanticTokenTypesproperty),
		string(lsproto.SemanticTokenTypesenumMember),
		string(lsproto.SemanticTokenTypesdecorator),
		string(lsproto.SemanticTokenTypesevent),
		string(lsproto.SemanticTokenTypesfunction),
		string(lsproto.SemanticTokenTypesmethod),
		string(lsproto.SemanticTokenTypesmacro),
		string(lsproto.SemanticTokenTypeslabel),
		string(lsproto.SemanticTokenTypescomment),
		string(lsproto.SemanticTokenTypesstring),
		string(lsproto.SemanticTokenTypeskeyword),
		string(lsproto.SemanticTokenTypesnumber),
		string(lsproto.SemanticTokenTypesregexp),
		string(lsproto.SemanticTokenTypesoperator),
	}
}

func defaultTokenModifiers() []string {
	return []string{
		string(lsproto.SemanticTokenModifiersdeclaration),
		string(lsproto.SemanticTokenModifiersdefinition),
		string(lsproto.SemanticTokenModifiersreadonly),
		string(lsproto.SemanticTokenModifiersstatic),
		string(lsproto.SemanticTokenModifiersdeprecated),
		string(lsproto.SemanticTokenModifiersabstract),
		string(lsproto.SemanticTokenModifiersasync),
		string(lsproto.SemanticTokenModifiersmodification),
		string(lsproto.SemanticTokenModifiersdocumentation),
		string(lsproto.SemanticTokenModifiersdefaultLibrary),
	}
}
