package autoimport

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ExportSyntax int

const (
	// export const x = {}
	ExportSyntaxModifier ExportSyntax = iota
	// export { x }
	ExportSyntaxNamed
	// export default function f() {}
	ExportSyntaxDefaultModifier
	// export default f
	ExportSyntaxDefaultDeclaration
	// export = x
	ExportSyntaxEquals
)

type RawExport struct {
	Syntax ExportSyntax
	Name   string
	// !!! other kinds of names
	Path                  tspath.Path
	ModuleDeclarationName string
}

func Parse(file *ast.SourceFile) []*RawExport {
	if file.Symbol != nil {
		return parseModule(file)
	}

	// !!!
	return nil
}

func parseModule(file *ast.SourceFile) []*RawExport {
	exports := make([]*RawExport, 0, len(file.Symbol.Exports))
	for name, symbol := range file.Symbol.Exports {
		if len(symbol.Declarations) != 1 {
			// !!! for debugging
			panic(fmt.Sprintf("unexpected number of declarations at %s: %s", file.Path(), name))
		}
		var syntax ExportSyntax
		switch symbol.Declarations[0].Kind {
		case ast.KindExportSpecifier:
			syntax = ExportSyntaxNamed
		case ast.KindExportAssignment:
			syntax = core.IfElse(
				symbol.Declarations[0].AsExportAssignment().IsExportEquals,
				ExportSyntaxEquals,
				ExportSyntaxDefaultDeclaration,
			)
		default:
			syntax = ExportSyntaxModifier
		}

		exports = append(exports, &RawExport{
			Syntax: syntax,
			Name:   name,
			Path:   file.Path(),
		})
	}
	return exports
}
