package ls

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func IsInString(sourceFile *ast.SourceFile, position int, previousToken *ast.Node) bool {
	if previousToken != nil && ast.IsStringTextContainingNode(previousToken) {
		start := previousToken.Pos()
		end := previousToken.End()

		// To be "in" one of these literals, the position has to be:
		//   1. entirely within the token text.
		//   2. at the end position of an unterminated token.
		//   3. at the end of a regular expression (due to trailing flags like '/foo/g').
		if start < position && position < end {
			return true
		}

		if position == end {
			return true
			//return !!(previousToken as LiteralExpression).isUnterminated; tbd
		}
	}
	return false
}

func RangeContainsRange(r1 core.TextRange, r2 core.TextRange) bool {
	return startEndContainsRange(r1.Pos(), r1.End(), r2)
}

func startEndContainsRange(start int, end int, textRange core.TextRange) bool {
	return start <= textRange.Pos() && end >= textRange.End()
}

func getPossibleGenericSignatures(called *ast.Expression, typeArgumentCount int, c *checker.Checker) []*checker.Signature {
	typeAtLocation := c.GetTypeAtLocation(called)
	if ast.IsOptionalChain(called.Parent) {
		typeAtLocation = removeOptionality(typeAtLocation, ast.IsOptionalChainRoot(called.Parent), true /*isOptionalChain*/, c)
	}
	var signatures []*checker.Signature
	if ast.IsNewExpression(called.Parent) {
		signatures = c.GetSignaturesOfType(typeAtLocation, checker.SignatureKindConstruct)
	} else {
		signatures = c.GetSignaturesOfType(typeAtLocation, checker.SignatureKindCall)
	}
	return core.Filter(signatures, func(s *checker.Signature) bool {
		return s.TypeParameters() != nil && len(s.TypeParameters()) >= typeArgumentCount
	})
}

func removeOptionality(t *checker.Type, isOptionalExpression bool, isOptionalChain bool, c *checker.Checker) *checker.Type {
	if isOptionalExpression {
		return c.GetNonNullableType(t)
	} else if isOptionalChain {
		return c.GetNonOptionalType(t)
	}
	return t
}

// nodeTests.ts
func isNoSubstitutionTemplateLiteral(node *ast.Node) bool {
	return node.Kind == ast.KindNoSubstitutionTemplateLiteral
}

func isTaggedTemplateExpression(node *ast.Node) bool {
	return node.Kind == ast.KindTaggedTemplateExpression
}

func isInsideTemplateLiteral(node *ast.Node, position int, sourceFile *ast.SourceFile) bool {
	return ast.IsTemplateLiteralKind(node.Kind) && (scanner.GetTokenPosOfNode(node, sourceFile, false) < position && position < node.End() || (ast.IsUnterminatedLiteral(node) && position == node.End()))
}

// Pseudo-literals
func isTemplateHead(node *ast.Node) bool {
	return node.Kind == ast.KindTemplateHead
}

func isTemplateTail(node *ast.Node) bool {
	return node.Kind == ast.KindTemplateTail
}

//

func findPrecedingMatchingToken(token *ast.Node, matchingTokenKind ast.Kind, sourceFile *ast.SourceFile) *ast.Node {
	closeTokenText := scanner.TokenToString(token.Kind)
	matchingTokenText := scanner.TokenToString(matchingTokenKind)
	tokenFullStart := token.Loc.Pos()
	// Text-scan based fast path - can be bamboozled by comments and other trivia, but often provides
	// a good, fast approximation without too much extra work in the cases where it fails.
	bestGuessIndex := strings.LastIndex(sourceFile.Text(), matchingTokenText)
	if bestGuessIndex == -1 {
		return nil // if the token text doesn't appear in the file, there can't be a match - super fast bail
	}
	// we can only use the textual result directly if we didn't have to count any close tokens within the range
	if strings.LastIndex(sourceFile.Text(), closeTokenText) < bestGuessIndex {
		nodeAtGuess := astnav.FindPrecedingToken(sourceFile, bestGuessIndex+1)
		if nodeAtGuess != nil && nodeAtGuess.Kind == matchingTokenKind {
			return nodeAtGuess
		}
	}
	tokenKind := token.Kind
	remainingMatchingTokens := 0
	for true {
		preceding := astnav.FindPrecedingToken(sourceFile, tokenFullStart)
		if preceding == nil {
			return nil
		}
		token = preceding
		if token.Kind == matchingTokenKind {
			if remainingMatchingTokens == 0 {
				return token
			}
			remainingMatchingTokens--
		} else if token.Kind == tokenKind {
			remainingMatchingTokens++
		}
	}
	return nil
}

// Display-part writer helpers
// var defaultMaximumTruncationLength = 160
// type displayPartsSymbolWriter struct {
// 	printer.EmitTextWriter
// 	lineStart 		bool
// 	indent 			int
// 	length 			int
// 	defaultMaximumTruncationLength int
// }
// func newDisplayPartsSymbolWriter() *displayPartsSymbolWriter {
// 	return &displayPartsSymbolWriter{
// 		EmitTextWriter: printer.NewTextWriter(""), //tbd
// 		lineStart: false,
// 		indent: 0,
// 		length: 0,
// 		defaultMaximumTruncationLength: defaultMaximumTruncationLength * 10, // A hard cutoff to avoid overloading the messaging channel in worst-case scenarios
// 	}
// }
// func (d *displayPartsSymbolWriter) displayParts() []symbolDisplayPart {
// 	finalText := len(d.String())
// 	return nil
// }

// type displayPartsSymbolWriterImpl struct {

// 	lineStart          bool
//     indent             int
//     length             int
//     absoluteMaxLength  int
//     defaultIndentation string
// }

// func (w *displayPartsSymbolWriterImpl) displayParts() []symbolDisplayPart {
// 	return nil
// }

// func NewDisplayPartsSymbolWriter() displayPartsSymbolWriter {
// 	writer := &displayPartsSymbolWriterImpl{
// 		lineStart:         true,
//         indent:            0,
//         length:            0,
//         absoluteMaxLength: 1000,
// 	}
// 	writer.Clear()
// 	return writer
// }

// func (w *displayPartsSymbolWriterImpl) Clear() {
//     w.displayParts = []symbolDisplayPart{}
//     w.lineStart = true
//     w.indent = 0
//     w.length = 0
// }
