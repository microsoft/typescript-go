package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

// TestIssue1389Fix tests the specific scenario from issue #1389
func TestIssue1389Fix(t *testing.T) {
	t.Parallel()
	
	// This is the exact content from the issue report
	content := "\"→\" ;\n\n\n"
	
	lineMap := ComputeLineStarts(content)
	
	converters := NewConverters(lsproto.PositionEncodingKindUTF16, func(fileName string) *LineMap {
		return lineMap
	})
	
	script := &mockScript{
		fileName: "test.ts",
		text:     content,
	}
	
	// Test the exact scenario that was causing the panic
	textLen := len(content)
	
	// These should not panic
	result1 := converters.PositionToLineAndCharacter(script, core.TextPos(textLen))
	t.Logf("Position at end of text (%d): line=%d, char=%d", textLen, result1.Line, result1.Character)
	
	result2 := converters.PositionToLineAndCharacter(script, core.TextPos(textLen+1))
	t.Logf("Position beyond end of text (%d): line=%d, char=%d", textLen+1, result2.Line, result2.Character)
	
	// Also test converting ranges, which is what the formatter does
	textRange := core.NewTextRange(0, textLen)
	lspRange := converters.ToLSPRange(script, textRange)
	
	t.Logf("Range conversion successful: start=(%d,%d), end=(%d,%d)", 
		lspRange.Start.Line, lspRange.Start.Character,
		lspRange.End.Line, lspRange.End.Character)
	
	// Test with a range that extends beyond the text (this was the original issue)
	textRangeBeyondEnd := core.NewTextRange(0, textLen+1)
	lspRangeBeyondEnd := converters.ToLSPRange(script, textRangeBeyondEnd)
	
	t.Logf("Range beyond end conversion successful: start=(%d,%d), end=(%d,%d)", 
		lspRangeBeyondEnd.Start.Line, lspRangeBeyondEnd.Start.Character,
		lspRangeBeyondEnd.End.Line, lspRangeBeyondEnd.End.Character)
}