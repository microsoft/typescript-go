package autoimport

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -type=ExportSyntax -output=parse_stringer_generated.go
//go:generate go tool mvdan.cc/gofumpt -w parse_stringer_generated.go

type ExportSyntax int
type ModuleID string

const (
	ExportSyntaxNone ExportSyntax = iota
	// export const x = {}
	ExportSyntaxModifier
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
	ModuleID             ModuleID
	NodeModulesDirectory tspath.Path
}

func (e *RawExport) Name() string {
	return e.ExportName
}

func Parse(file *ast.SourceFile, nodeModulesDirectory tspath.Path) []*RawExport {
	if file.Symbol != nil {
		return parseModule(file, nodeModulesDirectory)
	}

	// !!!
	return nil
}

func parseModule(file *ast.SourceFile, nodeModulesDirectory tspath.Path) []*RawExport {
	exports := make([]*RawExport, 0, len(file.Symbol.Exports))
	for name, symbol := range file.Symbol.Exports {
		if strings.HasPrefix(name, ast.InternalSymbolNamePrefix) {
			continue
		}
		var syntax ExportSyntax
		for _, decl := range symbol.Declarations {
			var declSyntax ExportSyntax
			switch decl.Kind {
			case ast.KindExportSpecifier:
				declSyntax = ExportSyntaxNamed
			case ast.KindExportAssignment:
				declSyntax = core.IfElse(
					decl.AsExportAssignment().IsExportEquals,
					ExportSyntaxEquals,
					ExportSyntaxDefaultDeclaration,
				)
			default:
				declSyntax = ExportSyntaxModifier
			}
			if syntax != ExportSyntaxNone && syntax != declSyntax {
				panic(fmt.Sprintf("mixed export syntaxes for symbol %s: %s", file.FileName(), name))
			}
			syntax = declSyntax
		}

		exports = append(exports, &RawExport{
			Syntax:               syntax,
			ExportName:           name,
			Flags:                symbol.Flags,
			FileName:             file.FileName(),
			Path:                 file.Path(),
			ModuleID:             ModuleID(file.Path()),
			NodeModulesDirectory: nodeModulesDirectory,
		})
	}
	// !!! handle module augmentations
	return exports
}
