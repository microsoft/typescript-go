package ls

import "github.com/microsoft/typescript-go/internal/ast"

func (l *LanguageService) ProvideHover(fileName string, position int) string {
	program := l.GetProgram()
	file := program.GetSourceFile(fileName)
	if file == nil {
		panic("file not found")
	}

	node := getTouchingPropertyName(file, position)
	if node.Kind == ast.KindSourceFile {
		// Avoid giving quickInfo for the sourceFile as a whole.
		return ""
	}

	checker := program.GetTypeChecker()
	symbol := checker.GetSymbolAtLocation(node)
	if t := checker.GetTypeOfSymbolAtLocation(symbol, node); t != nil {
		return t.Flags().String()
	}
	return ""
}
