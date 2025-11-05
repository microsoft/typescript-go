package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// validateRegularExpressionLiteralNode validates a regular expression literal during semantic analysis.
// This function is called from checkGrammarRegularExpressionLiteral and performs deep validation
// of the regexp pattern and flags that would affect code generation based on language version.
func (c *Checker) validateRegularExpressionLiteralNode(node *ast.RegularExpressionLiteral) {
	sourceFile := ast.GetSourceFileOfNode(node.AsNode())
	// Check if the token has the Unterminated flag
	if c.hasParseDiagnostics(sourceFile) || node.AsNode().LiteralLikeData().TokenFlags&ast.TokenFlagsUnterminated != 0 {
		return
	}

	var lastError *ast.Diagnostic

	// Create an error callback that collects diagnostics
	onError := func(message *diagnostics.Message, start int, length int, args ...any) {
		// Adjust start position relative to the node position in the source file (skipping leading trivia)
		nodeStart := scanner.GetTokenPosOfNode(node.AsNode(), sourceFile, false)
		adjustedStart := nodeStart + start

		// For providing spelling suggestions - if this is a suggestion message,
		// add it as related info to the previous error
		if message.Category() == diagnostics.CategoryMessage && lastError != nil &&
			adjustedStart == lastError.Pos() && length == lastError.Len() {
			err := ast.NewDiagnostic(nil, core.NewTextRange(adjustedStart, adjustedStart+length), message, args...)
			lastError.AddRelatedInfo(err)
		} else if lastError == nil || adjustedStart != lastError.Pos() {
			lastError = ast.NewDiagnostic(sourceFile, core.NewTextRange(adjustedStart, adjustedStart+length), message, args...)
			c.diagnostics.Add(lastError)
		}
	}

	// Perform regexp validation
	scanner.ValidateRegularExpressionLiteral(node, sourceFile, c.languageVersion, onError)
}
