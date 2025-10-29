package autoimport

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -type=ExportSyntax -output=parse_stringer_generated.go
//go:generate go tool mvdan.cc/gofumpt -w parse_stringer_generated.go

type ExportSyntax int
type ModuleID string

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
	Syntax     ExportSyntax
	ExportName string
	Flags      ast.SymbolFlags
	// !!! other kinds of names

	// The file where the export was found.
	FileName string
	Path     tspath.Path

	// ModuleID uniquely identifies a module across multiple declarations.
	// If the export is from an ambient module declaration, this is the module name.
	// If the export is from a module augmentation, this is the Path() of the resolved module file.
	// Otherwise this is the Path() of the exporting source file.
	ModuleID ModuleID
}

func (e *RawExport) Name() string {
	return e.ExportName
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
			Syntax:     syntax,
			ExportName: name,
			Flags:      symbol.Flags,
			FileName:   file.FileName(),
			Path:       file.Path(),
			ModuleID:   ModuleID(file.Path()),
		})
	}
	// !!! handle module augmentations
	return exports
}
