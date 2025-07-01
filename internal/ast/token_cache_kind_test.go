package ast_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
)

func TestTokenCachePanicOnKindMismatch(t *testing.T) {
	// This test verifies that kind mismatches still panic (as they should)
	// but parent mismatches don't panic anymore

	sourceText := "array?.at(0)"

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/test.ts",
		Path:     "/test.ts",
	}, sourceText, core.ScriptKindTS)

	pos := 5
	end := 6
	parent := &ast.Node{Kind: ast.KindCallExpression}

	// Create a token with one kind
	_ = sourceFile.GetOrCreateToken(ast.KindQuestionDotToken, pos, end, parent)

	// Try to create a token with a different kind at the same position
	// This should still panic
	defer func() {
		if r := recover(); r != nil {
			// Expected panic - test passes
			return
		}
		t.Fatal("Expected panic for kind mismatch")
	}()

	// This should panic
	_ = sourceFile.GetOrCreateToken(ast.KindDotToken, pos, end, parent)
}
