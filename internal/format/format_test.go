package format_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func makeFormatCtx(t *testing.T) context.Context {
	t.Helper()
	return format.WithFormatCodeSettings(t.Context(), &lsutil.FormatCodeSettings{
		EditorSettings: lsutil.EditorSettings{
			TabSize:                4,
			IndentSize:             4,
			NewLineCharacter:       "\n",
			ConvertTabsToSpaces:    true,
			IndentStyle:            lsutil.IndentStyleSmart,
			TrimTrailingWhitespace: true,
		},
		InsertSpaceAfterCommaDelimiter:                        core.TSTrue,
		InsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces: core.TSTrue,
	}, "\n")
}

func formatText(t *testing.T, text string) string {
	t.Helper()
	ctx := makeFormatCtx(t)
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/test.ts",
		Path:     "/test.ts",
	}, text, core.ScriptKindTS)
	edits := format.FormatDocument(ctx, sourceFile)
	return applyBulkEdits(text, edits)
}

// TestFormatImportBraces tests that the formatter correctly adds a newline before
// a closing brace when the import's opening brace is on a different line.
// Regression test for: tsgo formatter doesn't format imports the same as TS
func TestFormatImportBraces(t *testing.T) {
	t.Parallel()
	t.Run("multiline import gets closing brace on new line", func(t *testing.T) {
		t.Parallel()
		input := "import {\n\tbasename,\n\textname, joinPath } from '../base/resources.js';\n"
		got := formatText(t, input)
		want := "import {\n    basename,\n    extname, joinPath\n} from '../base/resources.js';\n"
		assert.Equal(t, want, got)
	})
	t.Run("single-line import is unchanged", func(t *testing.T) {
		t.Parallel()
		input := "import { basename } from '../base/resources.js';\n"
		got := formatText(t, input)
		assert.Equal(t, input, got)
	})
	t.Run("already-formatted multiline import is unchanged", func(t *testing.T) {
		t.Parallel()
		input := "import {\n    basename,\n    extname,\n    joinPath,\n} from '../base/resources.js';\n"
		got := formatText(t, input)
		assert.Equal(t, input, got)
	})
}

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
			ctx := format.WithFormatCodeSettings(t.Context(), &lsutil.FormatCodeSettings{
				EditorSettings: lsutil.EditorSettings{
					TabSize:                4,
					IndentSize:             4,
					BaseIndentSize:         4,
					NewLineCharacter:       "\n",
					ConvertTabsToSpaces:    true,
					IndentStyle:            lsutil.IndentStyleSmart,
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
