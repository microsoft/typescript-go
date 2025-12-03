package format_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/parser"
)

func TestTernaryWithTabs(t *testing.T) {
	t.Parallel()

	// This is the original code with tabs for indentation
	// Using literal tab characters
	text := "const test = (a: string) => (\n\ta === '1' ? (\n\t\t10\n\t) : (\n\t\t12\n\t)\n)"

	ctx := format.WithFormatCodeSettings(t.Context(), &format.FormatCodeSettings{
		EditorSettings: format.EditorSettings{
			TabSize:                4,
			IndentSize:             4,
			BaseIndentSize:         0,
			NewLineCharacter:       "\n",
			ConvertTabsToSpaces:    false, // Use tabs
			IndentStyle:            format.IndentStyleSmart,
			TrimTrailingWhitespace: true,
		},
	}, "\n")

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/test.ts",
		Path:     "/test.ts",
	}, text, core.ScriptKindTS)

	edits := format.FormatDocument(ctx, sourceFile)
	newText := applyBulkEdits(text, edits)

	// Check that we don't have mixed tabs and spaces
	// The formatted text should only use tabs for indentation
	lines := splitLines(newText)
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		// Check that the line doesn't have spaces for indentation followed by tabs
		hasLeadingSpaces := false
		for j := range len(line) {
			if line[j] == ' ' {
				hasLeadingSpaces = true
			} else if line[j] == '\t' {
				// If we already saw spaces and now see a tab, that's a problem
				if hasLeadingSpaces {
					t.Errorf("Line %d has mixed tabs and spaces: %q", i, line)
				}
			} else {
				// Hit non-whitespace
				break
			}
		}
	}
}

func splitLines(text string) []string {
	var lines []string
	start := 0
	for i := range len(text) {
		if text[i] == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}
