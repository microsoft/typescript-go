package scanner

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
)

// TestTypePredicate_Issue1092 tests the specific scenario that caused the panic in issue #1092
// This test ensures that type predicates with declaration maps enabled do not cause panics
func TestTypePredicate_Issue1092(t *testing.T) {
	// The original problematic text that caused the slice bounds out of range panic
	// This simulates the content from a declaration map that references a position beyond text bounds
	
	shortText := "import { foo } from './export';\nexport const x = foo();"
	textLen := len(shortText) // Should be 55 characters
	
	// The original issue had pos=158 (or 167) but text length was 55/58
	problematicPos := 158
	
	// This should not panic - before the fix, this would cause:
	// panic: runtime error: slice bounds out of range [158:55]
	result := SkipTriviaEx(shortText, problematicPos, nil)
	
	// The function should return the original position when it's beyond text bounds
	if result != problematicPos {
		t.Errorf("Expected position %d, got %d", problematicPos, result)
	}
	
	// Also test the GetLineAndCharacterOfPosition function that had a similar issue
	sourceFile := &mockSourceFile{
		text:    shortText,
		lineMap: []core.TextPos{0, 27}, // Approximate line starts
	}
	
	// This should also not panic
	line, char := GetLineAndCharacterOfPosition(sourceFile, problematicPos)
	
	// We should get valid results (not negative numbers)
	if line < 0 {
		t.Errorf("Got negative line number: %d", line)
	}
	if char < 0 {
		t.Errorf("Got negative character number: %d", char)
	}
	
	t.Logf("Successfully handled position %d beyond text length %d", problematicPos, textLen)
	t.Logf("Computed line: %d, character: %d", line, char)
}