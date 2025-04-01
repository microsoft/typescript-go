package ls

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
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
