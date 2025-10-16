package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) getExportsForAutoImport(ctx context.Context, fromFile *ast.SourceFile) (*autoimport.View, error) {
	// !!! snapshot integration
	registry, err := (&autoimport.Registry{}).Clone(ctx, autoimport.RegistryChange{
		WithProject: &autoimport.Project{
			Key:     "!!! TODO",
			Program: l.GetProgram(),
		},
	})
	if err != nil {
		return nil, err
	}

	view := autoimport.NewView(registry, fromFile, "!!! TODO")
	return view, nil
}

func (l *LanguageService) getAutoImportSourceFile(path tspath.Path) *ast.SourceFile {
	// !!! other sources
	return l.GetProgram().GetSourceFileByPath(path)
}
