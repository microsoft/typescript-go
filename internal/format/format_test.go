package format_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestFormatNoTrailingSpace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		text string
	}{
		{"simple statement without trailing newline", "1;"},
		{"function call without trailing newline", "console.log('hello');"},
		{"if block on single line", "if (true) { }"},
		{"class declaration", "class A {\n    // Class Contents Go Here\n}"},
		{"class declaration with trailing newline", "class A {\n    // Class Contents Go Here\n}\n"},
		{"empty block", "if (true) {}"},
		{"module declaration", "module M { }"},
		{"enum declaration", "enum E { A, B }"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := format.WithFormatCodeSettings(t.Context(), &lsutil.FormatCodeSettings{
				EditorSettings: lsutil.EditorSettings{
					TabSize:                4,
					IndentSize:             4,
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
			// Formatting should not add trailing whitespace at end of file
			for i, line := range strings.Split(newText, "\n") {
				trimmed := strings.TrimRight(line, " \t")
				assert.Equal(t, line, trimmed, "Formatter should not add trailing whitespace on line %d", i+1)
			}
		})
	}
}

func TestTernaryWithTabs(t *testing.T) {
	t.Parallel()

	text := "const test = (a: string) => (\n\ta === '1' ? (\n\t\t10\n\t) : (\n\t\t12\n\t)\n)"
	expected := "const test=(a: string) => (\n\ta==='1'? (\n\t\t10\n\t):(\n\t\t12\n\t)\n)"

	ctx := format.WithFormatCodeSettings(t.Context(), &format.FormatCodeSettings{
		EditorSettings: format.EditorSettings{
			TabSize:                4,
			IndentSize:             4,
			BaseIndentSize:         0,
			NewLineCharacter:       "\n",
			ConvertTabsToSpaces:    false,
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

	assert.Equal(t, expected, newText)
}
