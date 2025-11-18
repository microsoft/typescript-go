package autoimport

import (
	"fmt"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/module"
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
	// export * from "module"
	ExportSyntaxStar
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
	Syntax    ExportSyntax
	Flags     ast.SymbolFlags
	localName string

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
	if e.Syntax == ExportSyntaxStar {
		return e.Target.ExportName
	}
	if e.localName != "" {
		return e.localName
	}
	if strings.HasPrefix(e.ExportName, ast.InternalSymbolNamePrefix) {
		return "!!! TODO"
	}
	return e.ExportName
}

func parseFile(file *ast.SourceFile, nodeModulesDirectory tspath.Path, moduleResolver *module.Resolver, getChecker func() (*checker.Checker, func()), toPath func(string) tspath.Path) []*RawExport {
	if file.Symbol != nil {
		return parseModule(file, nodeModulesDirectory, moduleResolver, getChecker, toPath)
	}
	if len(file.AmbientModuleNames) > 0 {
		moduleDeclarations := core.Filter(file.Statements.Nodes, ast.IsModuleWithStringLiteralName)
		var exportCount int
		for _, decl := range moduleDeclarations {
			exportCount += len(decl.AsModuleDeclaration().Symbol.Exports)
		}
		exports := make([]*RawExport, 0, exportCount)
		for _, decl := range moduleDeclarations {
			parseModuleDeclaration(decl.AsModuleDeclaration(), file, ModuleID(decl.Name().Text()), nodeModulesDirectory, getChecker, &exports)
		}
		return exports
	}
	return nil
}

func parseModule(file *ast.SourceFile, nodeModulesDirectory tspath.Path, moduleResolver *module.Resolver, getChecker func() (*checker.Checker, func()), toPath func(string) tspath.Path) []*RawExport {
	moduleAugmentations := core.MapNonNil(file.ModuleAugmentations, func(name *ast.ModuleName) *ast.ModuleDeclaration {
		decl := name.Parent
		if ast.IsGlobalScopeAugmentation(decl) {
			return nil
		}
		return decl.AsModuleDeclaration()
	})
	var augmentationExportCount int
	for _, decl := range moduleAugmentations {
		augmentationExportCount += len(decl.Symbol.Exports)
	}
	exports := make([]*RawExport, 0, len(file.Symbol.Exports)+augmentationExportCount)
	for name, symbol := range file.Symbol.Exports {
		parseExport(name, symbol, ModuleID(file.Path()), file, nodeModulesDirectory, getChecker, &exports)
	}
	for _, decl := range moduleAugmentations {
		name := decl.Name().AsStringLiteral().Text
		moduleID := ModuleID(name)
		if tspath.IsExternalModuleNameRelative(name) {
			// !!! need to resolve non-relative names in separate pass
			if resolved, _ := moduleResolver.ResolveModuleName(name, file.FileName(), core.ModuleKindCommonJS, nil); resolved.IsResolved() {
				moduleID = ModuleID(toPath(resolved.ResolvedFileName))
			} else {
				// :shrug:
				moduleID = ModuleID(toPath(tspath.ResolvePath(tspath.GetDirectoryPath(file.FileName()), name)))
			}
		}
		parseModuleDeclaration(decl, file, moduleID, nodeModulesDirectory, getChecker, &exports)
	}
	return exports
}

func parseExport(name string, symbol *ast.Symbol, moduleID ModuleID, file *ast.SourceFile, nodeModulesDirectory tspath.Path, getChecker func() (*checker.Checker, func()), exports *[]*RawExport) {
	if name == ast.InternalSymbolNameExportStar {
		checker, release := getChecker()
		defer release()
		allExports := checker.GetExportsOfModule(symbol.Parent)
		// allExports includes named exports from the file that will be processed separately;
		// we want to add only the ones that come from the star
		for name, namedExport := range symbol.Parent.Exports {
			if name != ast.InternalSymbolNameExportStar {
				idx := slices.Index(allExports, namedExport)
				if idx >= 0 {
					allExports = slices.Delete(allExports, idx, idx+1)
				}
			}
		}
		for _, reexportedSymbol := range allExports {
			var scriptElementKind lsutil.ScriptElementKind
			var targetModuleID ModuleID
			if len(reexportedSymbol.Declarations) > 0 {
				scriptElementKind = lsutil.GetSymbolKind(checker, reexportedSymbol, reexportedSymbol.Declarations[0])
				// !!!
				targetModuleID = ModuleID(ast.GetSourceFileOfNode(reexportedSymbol.Declarations[0]).Path())
			}

			*exports = append(*exports, &RawExport{
				ExportID: ExportID{
					// !!! these are overlapping, what do I even want with this
					//     overlapping actually useful for merging later
					ExportName: name,
					ModuleID:   moduleID,
				},
				Syntax: ExportSyntaxStar,
				Flags:  reexportedSymbol.Flags,
				Target: ExportID{
					ExportName: reexportedSymbol.Name,
					ModuleID:   targetModuleID,
				},
				ScriptElementKind:          scriptElementKind,
				ScriptElementKindModifiers: lsutil.GetSymbolModifiers(checker, reexportedSymbol),
				FileName:                   file.FileName(),
				Path:                       file.Path(),
				NodeModulesDirectory:       nodeModulesDirectory,
			})
		}
		return
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
			if ast.GetCombinedModifierFlags(decl)&ast.ModifierFlagsDefault != 0 {
				declSyntax = ExportSyntaxDefaultModifier
			} else {
				declSyntax = ExportSyntaxModifier
			}
		}
		if syntax != ExportSyntaxNone && syntax != declSyntax {
			// !!! this can probably happen in erroring code
			panic(fmt.Sprintf("mixed export syntaxes for symbol %s: %s", file.FileName(), name))
		}
		syntax = declSyntax
	}

	var localName string
	if symbol.Name == ast.InternalSymbolNameDefault || symbol.Name == ast.InternalSymbolNameExportEquals {
		namedSymbol := symbol
		if s := binder.GetLocalSymbolForExportDefault(symbol); s != nil {
			namedSymbol = s
		}
		localName = getDefaultLikeExportNameFromDeclaration(namedSymbol)
		if localName == "" {
			localName = lsutil.ModuleSpecifierToValidIdentifier(string(moduleID), core.ScriptTargetESNext, false)
		}
	}

	export := &RawExport{
		ExportID: ExportID{
			ExportName: name,
			ModuleID:   moduleID,
		},
		Syntax:               syntax,
		localName:            localName,
		Flags:                symbol.Flags,
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
		export.ScriptElementKind = lsutil.GetSymbolKind(nil, symbol, symbol.Declarations[0])
		export.ScriptElementKindModifiers = lsutil.GetSymbolModifiers(nil, symbol)
	}

	*exports = append(*exports, export)
}

func parseModuleDeclaration(decl *ast.ModuleDeclaration, file *ast.SourceFile, moduleID ModuleID, nodeModulesDirectory tspath.Path, getChecker func() (*checker.Checker, func()), exports *[]*RawExport) {
	for name, symbol := range decl.Symbol.Exports {
		parseExport(name, symbol, moduleID, file, nodeModulesDirectory, getChecker, exports)
	}
}
