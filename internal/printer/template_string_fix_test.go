package printer_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/testutil/emittestutil"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
)

// TestTemplateStringEscapingFix tests the fix for TypeScript issue #59150
// This test mirrors the test case added in TypeScript PR #60303
func TestTemplateStringEscapingFix(t *testing.T) {
	t.Parallel()

	// https://github.com/microsoft/TypeScript/issues/59150
	// Replicates: printsCorrectly("template string", {}, printer =>
	//   printer.printNode(
	//     ts.EmitHint.Unspecified,
	//     ts.factory.createNoSubstitutionTemplateLiteral("\n"),
	//     ts.createSourceFile("source.ts", "", ts.ScriptTarget.ESNext),
	//   ));

	var factory ast.NodeFactory

	// Create a synthetic NoSubstitutionTemplateLiteral with just a newline character
	templateLiteral := factory.NewNoSubstitutionTemplateLiteral("\n")

	// Create a synthetic source file 
	file := factory.NewSourceFile("/source.ts", "/source.ts", "/source.ts", factory.NewNodeList([]*ast.Node{
		factory.NewExpressionStatement(templateLiteral),
	}))
	ast.SetParentInChildren(file)
	parsetestutil.MarkSyntheticRecursive(file)

	// The fix ensures that LF newlines are NOT escaped in template literals
	// Expected: `\n` (with literal newline character, not escaped)
	// Bug would produce: `\\n` (with escaped newline)
	emittestutil.CheckEmit(t, nil, file.AsSourceFile(), "`\n`;")
}