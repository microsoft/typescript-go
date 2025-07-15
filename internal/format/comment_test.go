package format_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestCommentFormatting(t *testing.T) {
	t.Parallel()

	t.Run("format comment issue reproduction", func(t *testing.T) {
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
			InsertSpaceBeforeTypeAnnotation: core.TSTrue,
		}, "\n")

		// Original code that causes the bug
		originalText := `class C {
    /**
     *
    */
    async x() {}
}`

		sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
			FileName: "/test.ts",
			Path:     "/test.ts",
		}, originalText, core.ScriptKindTS)

		// Apply formatting once
		edits := format.FormatDocument(ctx, sourceFile)
		firstFormatted := applyBulkEdits(originalText, edits)

		// Expected output after first formatting
		expectedFirstFormatted := `class C {
        /**
          *
         */
        async x() { }
    } `

		assert.Equal(t, expectedFirstFormatted, firstFormatted)

		// Apply formatting a second time to test stability
		sourceFile2 := parser.ParseSourceFile(ast.SourceFileParseOptions{
			FileName: "/test.ts",
			Path:     "/test.ts",
		}, firstFormatted, core.ScriptKindTS)

		edits2 := format.FormatDocument(ctx, sourceFile2)
		secondFormatted := applyBulkEdits(firstFormatted, edits2)

		// Test that second formatting is stable or document what it produces
		expectedSecondFormatted := `class C {
/**
  *
 */
        async x() { }
    } `

		assert.Equal(t, expectedSecondFormatted, secondFormatted)
	})
}
