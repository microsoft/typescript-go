package ls

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
)

func (l *LanguageService) GetSymbolAtPosition(fileName string, position int) *ast.Symbol {
	program, file := l.tryGetProgramAndFile(fileName)
	if file == nil {
		return nil
	}
	node := astnav.GetTokenAtPosition(file, position)
	if node == nil {
		return nil
	}
	checker := program.GetTypeChecker()
	return checker.GetSymbolAtLocation(node)
}
