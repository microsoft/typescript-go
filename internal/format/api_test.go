package format_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/repo"
	"gotest.tools/v3/assert"
)

func applyBulkEdits(text string, edits []core.TextChange) string {
	b := strings.Builder{}
	b.Grow(len(text))
	lastEnd := 0
	for _, e := range edits {
		start := e.TextRange.Pos()
		if start != lastEnd {
			b.WriteString(text[lastEnd:e.TextRange.Pos()])
		}
		b.WriteString(e.NewText)

		lastEnd = e.TextRange.End()
	}
	b.WriteString(text[lastEnd:])

	return b.String()
}

func TestFormat(t *testing.T) {
	t.Parallel()

	t.Run("format checker.ts", func(t *testing.T) {
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
		repo.SkipIfNoTypeScriptSubmodule(t)
		filePath := filepath.Join(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts")
		fileContent, err := os.ReadFile(filePath)
		assert.NilError(t, err)
		text := string(fileContent)
		sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
			FileName: "/checker.ts",
			Path:     "/checker.ts",
		}, text, core.ScriptKindTS)
		edits := format.FormatDocument(ctx, sourceFile)
		newText := applyBulkEdits(text, edits)
		assert.Assert(t, len(newText) > 0)
		assert.Assert(t, text != newText)
	})
}

func TestFormatNoTrailingNewline(t *testing.T) {
	t.Parallel()
	// Issue: Formatter adds extra space at end of line
	// When formatting a file that has content "1;" with no trailing newline,
	// an extra space should NOT be added at the end of the line
	
	t.Run("simple statement without trailing newline", func(t *testing.T) {
		t.Parallel()
		text := "1;"
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
		}, text, core.ScriptKindTS)
		edits := format.FormatDocument(ctx, sourceFile)
		
		// Debug: print edits to understand what's happening
		for i, edit := range edits {
			t.Logf("Edit %d: pos=%d, end=%d, newText=%q", i, edit.Pos(), edit.End(), edit.NewText)
		}
		
		newText := applyBulkEdits(text, edits)
		
		// The formatted text should be the same as the input - no extra space added
		assert.Equal(t, text, newText, "Formatter should not add extra space at end of line")
	})
}

func BenchmarkFormat(b *testing.B) {
	ctx := format.WithFormatCodeSettings(b.Context(), &format.FormatCodeSettings{
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
	repo.SkipIfNoTypeScriptSubmodule(b)
	filePath := filepath.Join(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts")
	fileContent, err := os.ReadFile(filePath)
	assert.NilError(b, err)
	text := string(fileContent)
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/checker.ts",
		Path:     "/checker.ts",
	}, text, core.ScriptKindTS)

	b.Run("format checker.ts", func(b *testing.B) {
		for b.Loop() {
			edits := format.FormatDocument(ctx, sourceFile)
			newText := applyBulkEdits(text, edits)
			assert.Assert(b, len(newText) > 0)
		}
	})

	b.Run("format checker.ts (no edit application)", func(b *testing.B) { // for comparison (how long does applying many edits take?)
		for b.Loop() {
			edits := format.FormatDocument(ctx, sourceFile)
			assert.Assert(b, len(edits) > 0)
		}
	})

	p := printer.NewPrinter(printer.PrinterOptions{}, printer.PrintHandlers{}, printer.NewEmitContext())
	b.Run("pretty print checker.ts", func(b *testing.B) { // for comparison
		for b.Loop() {
			newText := p.EmitSourceFile(sourceFile)
			assert.Assert(b, len(newText) > 0)
		}
	})
}
