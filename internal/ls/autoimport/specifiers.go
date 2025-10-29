package autoimport

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
)

func GetModuleSpecifier(
	fromFile *ast.SourceFile,
	export *RawExport,
	userPreferences modulespecifiers.UserPreferences,
	host modulespecifiers.ModuleSpecifierGenerationHost,
	compilerOptions *core.CompilerOptions,
) string {
	// !!! try using existing import
	specifiers, _ := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
		fromFile,
		export.FileName,
		compilerOptions,
		host,
		userPreferences,
		modulespecifiers.ModuleSpecifierOptions{},
		true,
	)
	if len(specifiers) > 0 {
		return specifiers[0]
	}
	return ""
}
