package parser_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// JSDocTest structure represents JSDoc parser tests
type JSDocTest struct {
	name     string
	source   string
	expected func(t *testing.T, file *ast.SourceFile)
}

// Test cases for JSDoc parser
var jsdocTests = []JSDocTest{
	{
		name: "Simple JSDoc comment",
		source: `
/**
 * This is a test description
 */
function test() {}
`,
		expected: func(t *testing.T, file *ast.SourceFile) {
			// Main function node
			functionNode := file.Statements.Nodes[0]

			// Check if JSDoc flag is set
			if functionNode.Flags&ast.NodeFlagsHasJSDoc == 0 {
				t.Errorf("JSDoc flag is not set")
			}
		},
	},
	{
		name: "JSDoc type description",
		source: `
/**
 * @param {string} name - User name
 * @returns {boolean} - Is operation successful?
 */
function validateName(name) {
  return name.length > 0;
}
`,
		expected: func(t *testing.T, file *ast.SourceFile) {
			functionNode := file.Statements.Nodes[0]

			// Check if JSDoc flag is set
			if functionNode.Flags&ast.NodeFlagsHasJSDoc == 0 {
				t.Errorf("JSDoc flag is not set")
			}
		},
	},
	
}

// Run JSDoc parser tests
func TestJSDocParser(t *testing.T) {
	for _, test := range jsdocTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel for better performance
			defer testutil.RecoverAndFail(t, "Panic during JSDoc parser test: "+test.name)

			// Parse TypeScript file with JSDoc comments
			fileName := "/main.ts"
			file := parser.ParseSourceFile(fileName, tspath.Path(fileName), test.source, core.ScriptTargetESNext, scanner.JSDocParsingModeParseAll)

			// There should be no errors
			if len(file.Diagnostics()) > 0 {
				t.Errorf("JSDoc parser error: %v", file.Diagnostics())
				return
			}

			// Run expected validations
			test.expected(t, file)
		})
	}
}

// Test that checks for JSDoc tags
func TestJSDocTags(t *testing.T) {
	t.Parallel() // Run test in parallel for better performance
	source := `
/**
 * @param {string} name - User name
 * @returns {boolean} - Is operation successful?
 * @deprecated This function is deprecated
 */
function validateName(name) {
  return name.length > 0;
}
`
	fileName := "/main.ts"
	file := parser.ParseSourceFile(fileName, tspath.Path(fileName), source, core.ScriptTargetESNext, scanner.JSDocParsingModeParseAll)

	// There should be no errors
	if len(file.Diagnostics()) > 0 {
		t.Errorf("JSDoc parser error: %v", file.Diagnostics())
		return
	}

	// Main function node
	functionNode := file.Statements.Nodes[0]

	// Check if JSDoc flag is set
	if functionNode.Flags&ast.NodeFlagsHasJSDoc == 0 {
		t.Errorf("JSDoc flag is not set")
	}

	// Deprecated flag should be set
	if functionNode.Flags&ast.NodeFlagsDeprecated == 0 {
		t.Errorf("Deprecated flag is not set")
	}
}
