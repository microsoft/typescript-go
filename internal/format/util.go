package format

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func rangeIsOnOneLine(node core.TextRange, file *ast.SourceFile) bool {
	startLine, _ := scanner.GetLineAndCharacterOfPosition(file, node.Pos())
	endLine, _ := scanner.GetLineAndCharacterOfPosition(file, node.End())
	return startLine == endLine
}
