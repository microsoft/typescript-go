package modulespecifiers

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func GetModuleSpecifiers(
	moduleSymbol *ast.Symbol,
	checker CheckerShape,
	compilerOptions *core.CompilerOptions,
	importingSourceFile *ast.SourceFile,
	host any,
	userPreferences *UserPreferences,
	options *ModuleSpecifierOptions,
) []string {
	return nil // !!!
}
