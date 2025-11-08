package format_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestFormatNoTrailingNewline(t *testing.T) {
	t.Parallel()
	// Issue: Formatter adds extra space at end of line
	// When formatting a file that has content "1;" with no trailing newline,
	// an extra space should NOT be added at the end of the line

	testCases := []struct {
		name string
		text string
	}{
		{"simple statement without trailing newline", "1;"},
		{"function call without trailing newline", "console.log('hello');"},
		{"variable declaration without trailing newline", "const x = 1;"},
		{"multiple statements without trailing newline", "const x = 1;\nconst y = 2;"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := format.WithFormatCodeSettings(t.Context(), &format.FormatCodeSettings{
				EditorSettings: format.EditorSettings{
					TabSize:                4,
					IndentSize:             4,
					BaseIndentSize:         4,
					NewLineCharacter:       "\n",
					ConvertTabsToSpaces:    true,
					IndentStyle:            format.IndentStyleSmart,
					TrimTrailingWhitespace: true,
				},
			}, "\n")
			sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
				FileName: "/test.ts",
				Path:     "/test.ts",
			}, tc.text, core.ScriptKindTS)
			edits := format.FormatDocument(ctx, sourceFile)
			newText := applyBulkEdits(tc.text, edits)

			// The formatted text should not add extra space at the end
			// It may add proper spacing within the code, but not after the last character
			assert.Assert(t, !strings.HasSuffix(newText, " "), "Formatter should not add trailing space")
			// Also check that no space was added at EOF position if text didn't end with newline
			if !strings.HasSuffix(tc.text, "\n") {
				assert.Assert(t, !strings.HasSuffix(newText, " "), "Formatter should not add space before EOF")
			}
		})
	}
}

// Test for panic in childStartsOnTheSameLineWithElseInIfStatement
// when FindPrecedingToken returns nil (Issue: panic handling request textDocument/onTypeFormatting)
func TestFormatOnEnter_NilPrecedingToken(t *testing.T) {
	t.Parallel()

	// Test case where else statement is at the beginning of the file
	// which can cause FindPrecedingToken to return nil
	testCases := []struct {
		name     string
		text     string
		position int // position where enter is pressed
	}{
		{
			name:     "else at file start - edge case",
			text:     "if(a){}\nelse{}",
			position: 9, // After the newline, before 'else'
		},
		{
			name:     "simple if-else with enter after if block",
			text:     "if (true) {\n}\nelse {\n}",
			position: 13, // After "}\n", before "else"
		},
		{
			name:     "if-else with enter in else block",
			text:     "if (true) {\n} else {\n}",
			position: 21, // Inside else block
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
				FileName: "/test.ts",
				Path:     "/test.ts",
			}, tc.text, core.ScriptKindTS)

			ctx := format.WithFormatCodeSettings(context.Background(), &format.FormatCodeSettings{
				EditorSettings: format.EditorSettings{
					TabSize:             4,
					IndentSize:          4,
					NewLineCharacter:    "\n",
					ConvertTabsToSpaces: true,
					IndentStyle:         format.IndentStyleSmart,
				},
			}, "\n")

			// This should not panic
			edits := format.FormatOnEnter(ctx, sourceFile, tc.position)
			_ = edits // Just ensuring no panic
		})
	}
}
