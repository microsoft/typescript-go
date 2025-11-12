package autoimport

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -type=ExportSyntax -output=parse_stringer_generated.go
//go:generate go tool mvdan.cc/gofumpt -w parse_stringer_generated.go

// ModuleID uniquely identifies a module across multiple declarations.
// If the export is from an ambient module declaration, this is the module name.
// If the export is from a module augmentation, this is the Path() of the resolved module file.
// Otherwise this is the Path() of the exporting source file.
type ModuleID string

type ExportID struct {
	ModuleID   ModuleID
	ExportName string
}

type ExportSyntax int

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

func (s ExportSyntax) IsAlias() bool {
	switch s {
	case ExportSyntaxNamed, ExportSyntaxEquals, ExportSyntaxDefaultDeclaration:
		return true
	default:
		return false
	}
}

type RawExport struct {
	ExportID
	Syntax ExportSyntax
	Flags  ast.SymbolFlags

	// Checker-set fields

	Target                     ExportID
	ScriptElementKind          lsutil.ScriptElementKind
	ScriptElementKindModifiers collections.Set[lsutil.ScriptElementKindModifier]

	// The file where the export was found.
	FileName string
	Path     tspath.Path

	NodeModulesDirectory tspath.Path
}

func (e *RawExport) Name() string {
	return e.ExportName
}

func Parse(file *ast.SourceFile, nodeModulesDirectory tspath.Path, getChecker func() (*checker.Checker, func())) []*RawExport {
	if file.Symbol != nil {
		return parseModule(file, nodeModulesDirectory, getChecker)
	}

	// !!!
	return nil
}

func parseModule(file *ast.SourceFile, nodeModulesDirectory tspath.Path, getChecker func() (*checker.Checker, func())) []*RawExport {
	exports := make([]*RawExport, 0, len(file.Symbol.Exports))
	for name, symbol := range file.Symbol.Exports {
		if strings.HasPrefix(name, ast.InternalSymbolNamePrefix) {
			// !!! resolve these and determine names
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
				// !!! this can probably happen in erroring code
				panic(fmt.Sprintf("mixed export syntaxes for symbol %s: %s", file.FileName(), name))
			}
			syntax = declSyntax
		}

		export := &RawExport{
			ExportID: ExportID{
				ExportName: name,
				ModuleID:   ModuleID(file.Path()),
			},
			Syntax:               syntax,
			Flags:                symbol.Flags,
			ScriptElementKind:    lsutil.GetSymbolKindSimple(symbol),
			FileName:             file.FileName(),
			Path:                 file.Path(),
			NodeModulesDirectory: nodeModulesDirectory,
		}

		if symbol.Flags&ast.SymbolFlagsAlias != 0 {
			checker, release := getChecker()
			targetSymbol := checker.GetAliasedSymbol(symbol)
			if !checker.IsUnknownSymbol(targetSymbol) {
				var decl *ast.Node
				if len(targetSymbol.Declarations) > 0 {
					decl = targetSymbol.Declarations[0]
				} else if len(symbol.Declarations) > 0 {
					decl = symbol.Declarations[0]
				}
				if decl == nil {
					panic("I want to know how this can happen")
				}
				export.ScriptElementKind = lsutil.GetSymbolKind(checker, targetSymbol, decl)
				export.ScriptElementKindModifiers = lsutil.GetSymbolModifiers(checker, targetSymbol)
				// !!! completely wrong
				// do we need this for anything other than grouping reexports?
				export.Target = ExportID{
					ExportName: targetSymbol.Name,
					ModuleID:   ModuleID(ast.GetSourceFileOfNode(decl).Path()),
				}
			}
			release()
		} else {
			export.ScriptElementKind = lsutil.GetSymbolKindSimple(symbol)
			export.ScriptElementKindModifiers = lsutil.GetSymbolModifiers(nil, symbol)
		}

		exports = append(exports, export)
	}
	// !!! handle module augmentations
	return exports
}
