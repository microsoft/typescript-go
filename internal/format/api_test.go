package format_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/scanner"
	"gotest.tools/v3/assert"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	t.Run("format checker.ts", func(t *testing.T) {
		ctx := format.NewContext(context.Background(), &format.FormatCodeSettings{
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
		sourceFile := parser.ParseSourceFile(
			"/checker.ts",
			"/checker.ts",
			text,
			core.ScriptTargetESNext,
			scanner.JSDocParsingModeParseAll,
		)
		ast.SetParentInChildren(sourceFile.AsNode())
		edits := format.FormatDocument(ctx, sourceFile)
		newText := text
		for _, e := range edits {
			newText = e.ApplyTo(newText)
		}
	})
}

func BenchmarkFormat(b *testing.B) {
	b.Run("format checker.ts", func(b *testing.B) {
		ctx := format.NewContext(context.Background(), &format.FormatCodeSettings{
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
		sourceFile := parser.ParseSourceFile(
			"/checker.ts",
			"/checker.ts",
			text,
			core.ScriptTargetESNext,
			scanner.JSDocParsingModeParseAll,
		)
		ast.SetParentInChildren(sourceFile.AsNode())
		for b.Loop() {
			edits := format.FormatDocument(ctx, sourceFile)
			newText := text
			for _, e := range edits {
				newText = e.ApplyTo(newText)
			}
		}
	})
}
