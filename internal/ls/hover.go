package ls

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
)

func (l *LanguageService) ProvideHover(fileName string, position int) (string, error) {
	program, file, err := l.getProgramAndFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to get program and file for hover: %w", err)
	}
	node := astnav.GetTouchingPropertyName(file, position)
	if node.Kind == ast.KindSourceFile {
		// Avoid giving quickInfo for the sourceFile as a whole.
		return "", nil
	}

	checker := program.GetTypeChecker()
	if symbol := checker.GetSymbolAtLocation(node); symbol != nil {
		if t := checker.GetTypeOfSymbolAtLocation(symbol, node); t != nil {
			return checker.TypeToString(t), nil
		}
	}
	return "", nil
}
