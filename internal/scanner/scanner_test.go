package scanner

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"gotest.tools/v3/assert"
)

func TestScanStringPreservesLoneSurrogates(t *testing.T) {
	t.Parallel()
	s := NewScanner()
	s.SetText(`"🦀\ud7ff\ud800\ud801\uD83E\uDD80"`)
	assert.Equal(t, s.Scan(), ast.KindStringLiteral)
	assert.Equal(t, s.TokenValue(), "🦀"+
		stringutil.EncodeJSStringRune(0xD7FF)+
		stringutil.EncodeJSStringRune(0xD800)+
		stringutil.EncodeJSStringRune(0xD801)+
		"🦀")
}

func TestNormalizeJSDocTypeSourceText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		text          string
		expectedLines []string
	}{
		{name: "single line", text: " \t* \tFoo", expectedLines: []string{"Foo"}},
		{name: "ECMAScript line breaks", text: "Foo\r\n * Bar\r\t* Baz\u2028 * Qux\u2029* Quux", expectedLines: []string{"Foo", "Bar", "Baz", "Qux", "Quux"}},
		{name: "blank and trailing lines", text: "Foo\r\n *\r\n", expectedLines: []string{"Foo", "", ""}},
		{name: "line without marker", text: "Foo\n  Bar", expectedLines: []string{"Foo", "Bar"}},
		{name: "only leading marker", text: "**Foo", expectedLines: []string{"*Foo"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			expected := strings.Join(test.expectedLines, core.NewLineKindLF.GetNewLineCharacter())
			assert.Equal(t, normalizeJSDocTypeSourceText(test.text), expected)
		})
	}
}

func TestGetTextOfNodeFromJSDocTypePreservesAsteriskType(t *testing.T) {
	t.Parallel()

	sourceText := strings.Join([]string{"", " * *"}, core.NewLineKindLF.GetNewLineCharacter())
	node := &ast.Node{
		Kind:  ast.KindJSDocAllType,
		Flags: ast.NodeFlagsJSDoc | ast.NodeFlagsInJSDocType,
		Loc:   core.NewTextRange(0, len(sourceText)),
	}

	assert.Equal(t, GetTextOfNodeFromSourceText(sourceText, node, false /*includeTrivia*/), "*")
}
