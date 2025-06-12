package printer_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/testutil/emittestutil"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
)

func TestTemplateStringEscaping(t *testing.T) {
	t.Parallel()

	t.Run("ParsedTemplateLiteral", func(t *testing.T) {
		// Test case that parses template literal from source
		file := parsetestutil.ParseTypeScript("`\n`", false)
		emittestutil.CheckEmit(t, nil, file, "`\n`;")
	})

	t.Run("SyntheticTemplateLiteral", func(t *testing.T) {
		// This is the test case from TypeScript PR #60303
		// It should NOT escape the newline character in template literals
		var factory ast.NodeFactory
		
		// Create a template literal with just a newline character  
		templateLiteral := factory.NewNoSubstitutionTemplateLiteral("\n")
		
		// Create a synthetic source file containing the template literal
		file := factory.NewSourceFile("/test.ts", "/test.ts", "/test.ts", factory.NewNodeList([]*ast.Node{
			factory.NewExpressionStatement(templateLiteral),
		}))
		ast.SetParentInChildren(file)
		parsetestutil.MarkSyntheticRecursive(file)
		
		// The expected result should NOT escape the newline
		expected := "`\n`;"
		emittestutil.CheckEmit(t, nil, file.AsSourceFile(), expected)
	})

	t.Run("SyntheticTemplateLiteralExplicitTest", func(t *testing.T) {
		// More explicit test to see what's really happening
		var factory ast.NodeFactory
		
		templateLiteral := factory.NewNoSubstitutionTemplateLiteral("\n")
		
		file := factory.NewSourceFile("/test.ts", "/test.ts", "/test.ts", factory.NewNodeList([]*ast.Node{
			factory.NewExpressionStatement(templateLiteral),
		}))
		ast.SetParentInChildren(file)
		parsetestutil.MarkSyntheticRecursive(file)
		
		// Print the template literal directly using the same setup as emittestutil
		p := printer.NewPrinter(
			printer.PrinterOptions{
				NewLine: core.NewLineKindLF,
			},
			printer.PrintHandlers{},
			nil, // emitContext
		)
		result := p.EmitSourceFile(file.AsSourceFile())
		
		t.Logf("Full file result: %q", result)
		
		// Check if the newline is properly preserved (not escaped)
		expected1 := "`\n`;\n" // What we should get
		expected2 := "`\\n`;\n" // What would indicate the bug
		
		if result == expected1 {
			t.Log("✓ Newline is correctly preserved (not escaped)")
		} else if result == expected2 {
			t.Error("✗ Bug confirmed: newline is incorrectly escaped")
		} else {
			t.Errorf("✗ Unexpected result: %q", result)
		}
	})
}