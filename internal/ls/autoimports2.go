package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
)

func (l *LanguageService) getExportsForAutoImport(ctx context.Context, fromFile *ast.SourceFile) (*autoimport.View, error) {
	registry := l.host.AutoImportRegistry()
	if !registry.IsPreparedForImportingFile(fromFile.FileName(), l.projectPath) {
		return nil, ErrNeedsAutoImports
	}

	view := autoimport.NewView(registry, fromFile, l.projectPath)
	return view, nil
}
