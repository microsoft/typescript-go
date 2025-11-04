package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
)

func (l *LanguageService) getExportsForAutoImport(ctx context.Context, fromFile *ast.SourceFile) (*autoimport.View, error) {
	registry := l.host.AutoImportRegistry()
	view := autoimport.NewView(registry, fromFile, "!!! TODO")
	return view, nil
}
