package ast_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestTokenCacheParentMismatch(t *testing.T) {
	// This test reproduces the scenario that caused the panic:
	// When AST structure changes, the same token position might have different parents.

	sourceText := "array?.at(0)"

	// Parse the source file
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/test.ts",
		Path:     "/test.ts",
	}, sourceText, core.ScriptKindTS)

	// Simulate getting tokens for different parent nodes at the same position
	// This would happen when the AST structure changes but the same token positions are accessed
	pos := 5 // Position of the "?" token
	end := 6

	// Create a token with one parent
	parent1 := &ast.Node{Kind: ast.KindCallExpression}
	token1 := sourceFile.GetOrCreateToken(ast.KindQuestionDotToken, pos, end, parent1)

	// Create a token with a different parent at the same position
	// This should not panic, but create a new token
	parent2 := &ast.Node{Kind: ast.KindAsExpression}
	token2 := sourceFile.GetOrCreateToken(ast.KindQuestionDotToken, pos, end, parent2)

	// Verify both tokens have the correct properties
	assert.Equal(t, token1.Kind, ast.KindQuestionDotToken)
	assert.Equal(t, token2.Kind, ast.KindQuestionDotToken)
	assert.Equal(t, token1.Parent, parent1)
	assert.Equal(t, token2.Parent, parent2)

	// The tokens should be different objects since they have different parents
	assert.Assert(t, token1 != token2, "tokens with different parents should be different objects")

	// But requesting the same token with the same parent should return the cached token
	token1Again := sourceFile.GetOrCreateToken(ast.KindQuestionDotToken, pos, end, parent1)
	assert.Equal(t, token1, token1Again, "token with same parent should be cached")
}
