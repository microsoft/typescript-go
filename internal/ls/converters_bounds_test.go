package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type mockScript struct {
	fileName string
	text     string
}

func (m *mockScript) FileName() string {
	return m.fileName
}

func (m *mockScript) Text() string {
	return m.text
}

func TestPositionToLineAndCharacterBoundsCheck(t *testing.T) {
	t.Parallel()
	
	// Test case that reproduces the panic with multi-byte characters and trailing newlines
	text := "\"→\" ;\n\n\n"
	
	lineMap := ComputeLineStarts(text)
	
	converters := NewConverters(lsproto.PositionEncodingKindUTF16, func(fileName string) *LineMap {
		return lineMap
	})
	
	script := &mockScript{
		fileName: "test.ts",
		text:     text,
	}
	
	// This should not panic even if position is beyond text length
	textLen := len(text)
	position := core.TextPos(textLen) // position at end of text
	
	// This should work
	result := converters.PositionToLineAndCharacter(script, position)
	
	t.Logf("Text length: %d, Position: %d", textLen, position)
	t.Logf("Result: line=%d, char=%d", result.Line, result.Character)
	
	// This should also not panic, even though position is beyond text length
	beyondEndPosition := core.TextPos(textLen + 1)
	result2 := converters.PositionToLineAndCharacter(script, beyondEndPosition)
	
	t.Logf("Beyond end position: %d", beyondEndPosition)
	t.Logf("Result2: line=%d, char=%d", result2.Line, result2.Character)
}