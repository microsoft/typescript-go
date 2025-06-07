package scanner

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
)

func TestSkipTriviaEx_BoundsCheck(t *testing.T) {
	// Test with trivia (spaces and tabs)
	text := "  \t hello"
	
	// Test normal functionality - should skip leading whitespace
	result := SkipTriviaEx(text, 0, nil)
	if result != 4 { // Should skip 2 spaces and 1 tab to position 4 ('h')
		t.Errorf("Expected position 4, got %d", result)
	}

	textLen := len(text) // 8 characters

	// Test position at end of text
	result = SkipTriviaEx(text, textLen, nil)
	if result != textLen {
		t.Errorf("Expected position %d, got %d", textLen, result)
	}

	// Test position beyond text bounds (should not panic)
	result = SkipTriviaEx(text, textLen+10, nil)
	if result != textLen+10 {
		t.Errorf("Expected position %d, got %d", textLen+10, result)
	}

	// Test position way beyond text bounds (reproduces the original issue)
	result = SkipTriviaEx(text, 158, nil)
	if result != 158 {
		t.Errorf("Expected position 158, got %d", result)
	}
}

func TestGetLineAndCharacterOfPosition_BoundsCheck(t *testing.T) {
	// Create a mock source file with a small text
	sourceFile := &mockSourceFile{
		text:    "hello\nworld",
		lineMap: []core.TextPos{0, 6}, // Line 0 starts at 0, line 1 starts at 6
	}

	textLen := len(sourceFile.text) // 11 characters

	// Test position within bounds
	line, char := GetLineAndCharacterOfPosition(sourceFile, 7)
	if line != 1 || char != 1 {
		t.Errorf("Expected line 1, char 1, got line %d, char %d", line, char)
	}

	// Test position beyond text bounds (should not panic)
	line, char = GetLineAndCharacterOfPosition(sourceFile, textLen+10)
	if line < 0 {
		t.Errorf("Expected valid line, got %d", line)
	}

	// Test the exact scenario from the bug report
	line, char = GetLineAndCharacterOfPosition(sourceFile, 158)
	if line < 0 {
		t.Errorf("Expected valid line, got %d", line)
	}
}

// Mock source file for testing
type mockSourceFile struct {
	text    string
	lineMap []core.TextPos
}

func (m *mockSourceFile) Text() string {
	return m.text
}

func (m *mockSourceFile) LineMap() []core.TextPos {
	return m.lineMap
}