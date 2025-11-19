package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
)

func (l *LanguageService) getAutoImportView(ctx context.Context, fromFile *ast.SourceFile, program *compiler.Program) (*autoimport.View, error) {
	registry := l.host.AutoImportRegistry()
	if !registry.IsPreparedForImportingFile(fromFile.FileName(), l.projectPath) {
		return nil, ErrNeedsAutoImports
	}

	view := autoimport.NewView(registry, fromFile, l.projectPath, program, l.UserPreferences().ModuleSpecifierPreferences())
	return view, nil
}
