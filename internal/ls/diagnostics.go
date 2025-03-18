package ls

import (
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
)

func (l *LanguageService) GetDocumentDiagnostics(fileName string) ([]*ast.Diagnostic, error) {
	program, file, err := l.getProgramAndFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get program and file for diagnostics: %w", err)
	}
	syntaxDiagnostics := program.GetSyntacticDiagnostics(file)
	semanticDiagnostics := program.GetSemanticDiagnostics(file)
	return slices.Concat(syntaxDiagnostics, semanticDiagnostics), nil
}
