package ls

import (
	"context"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) getExportsForAutoImport(ctx context.Context) (*autoimport.Registry, error) {
	// !!! snapshot integration
	return autoimport.Collect(ctx, l.GetProgram().GetSourceFiles())
}

func (l *LanguageService) getAutoImportSourceFile(path tspath.Path) *ast.SourceFile {
	// !!! other sources
	return l.GetProgram().GetSourceFileByPath(path)
}

func isInUnreachableNodeModules(from, to string) bool {
	nodeModulesIndexTo := strings.Index(to, "/node_modules/")
	if nodeModulesIndexTo == -1 {
		return false
	}

}
