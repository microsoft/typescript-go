package incremental

import (
	"encoding/binary"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func computeDeclarationInputSignature(text string) string {
	scan := scanner.NewScanner()
	scan.SetText(text)
	scan.SetSkipTrivia(false)

	var normalized strings.Builder
	var header [9]byte
	hasPrecedingLineBreak := false
	hasSignificantToken := false
	for {
		kind := scan.Scan()
		switch kind {
		case ast.KindEndOfFile:
			return ComputeHash(normalized.String(), false)
		case ast.KindNewLineTrivia:
			hasPrecedingLineBreak = true
			continue
		case ast.KindWhitespaceTrivia:
			continue
		case ast.KindSingleLineCommentTrivia, ast.KindMultiLineCommentTrivia:
			if !isDeclarationInputComment(kind, scan.TokenText()) {
				if kind == ast.KindMultiLineCommentTrivia && strings.ContainsAny(scan.TokenText(), "\r\n") {
					hasPrecedingLineBreak = true
				}
				continue
			}
		}

		if hasSignificantToken && hasPrecedingLineBreak {
			header[0] = 1
		} else {
			header[0] = 0
		}
		hasPrecedingLineBreak = false
		hasSignificantToken = true
		binary.LittleEndian.PutUint32(header[1:5], uint32(kind))
		binary.LittleEndian.PutUint32(header[5:9], uint32(len(scan.TokenText())))
		normalized.Write(header[:])
		normalized.WriteString(scan.TokenText())
	}
}

func isDeclarationInputComment(kind ast.Kind, text string) bool {
	switch kind {
	case ast.KindMultiLineCommentTrivia:
		return strings.HasPrefix(text, "/**") ||
			strings.HasPrefix(strings.TrimLeft(text[2:], " \t"), "@")
	case ast.KindSingleLineCommentTrivia:
		return strings.HasPrefix(text, "///") ||
			strings.HasPrefix(strings.TrimLeft(text[2:], " \t"), "@")
	default:
		return false
	}
}
