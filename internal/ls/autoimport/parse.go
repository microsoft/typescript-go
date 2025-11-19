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
	PackageName          string
}

func (e *RawExport) Name() string {
	if e.Syntax == ExportSyntaxStar {
		return e.Target.ExportName
	}
	if e.localName != "" {
		return e.localName
	}
	if e.ExportName == ast.InternalSymbolNameExportEquals {
		return e.Target.ExportName
	}
	if strings.HasPrefix(e.ExportName, ast.InternalSymbolNamePrefix) {
		panic("unexpected internal symbol name in export")
	}
	return e.ExportName
}

func (e *RawExport) AmbientModuleName() string {
	if !tspath.IsExternalModuleNameRelative(string(e.ModuleID)) {
		return string(e.ModuleID)
	}
	return ""
}

func (e *RawExport) ModuleFileName() string {
	if e.AmbientModuleName() == "" {
		return string(e.ModuleID)
	}
	return ""
}

func (b *registryBuilder) parseFile(file *ast.SourceFile, nodeModulesDirectory tspath.Path, packageName string, getChecker func() (*checker.Checker, func())) []*RawExport {
	if file.Symbol != nil {
		return b.parseModule(file, nodeModulesDirectory, packageName, getChecker)
	}
	if len(file.AmbientModuleNames) > 0 {
		moduleDeclarations := core.Filter(file.Statements.Nodes, ast.IsModuleWithStringLiteralName)
		var exportCount int
		for _, decl := range moduleDeclarations {
			exportCount += len(decl.AsModuleDeclaration().Symbol.Exports)
		}
		exports := make([]*RawExport, 0, exportCount)
		for _, decl := range moduleDeclarations {
			parseModuleDeclaration(decl.AsModuleDeclaration(), file, ModuleID(decl.Name().Text()), nodeModulesDirectory, packageName, getChecker, &exports)
		}
		return exports
	}
	return nil
}

func (b *registryBuilder) parseModule(file *ast.SourceFile, nodeModulesDirectory tspath.Path, packageName string, getChecker func() (*checker.Checker, func())) []*RawExport {
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
		parseExport(name, symbol, ModuleID(file.Path()), file, nodeModulesDirectory, packageName, getChecker, &exports)
	}
	for _, decl := range moduleAugmentations {
		name := decl.Name().AsStringLiteral().Text
		moduleID := ModuleID(name)
		if tspath.IsExternalModuleNameRelative(name) {
			// !!! need to resolve non-relative names in separate pass
			if resolved, _ := b.resolver.ResolveModuleName(name, file.FileName(), core.ModuleKindCommonJS, nil); resolved.IsResolved() {
				moduleID = ModuleID(b.base.toPath(resolved.ResolvedFileName))
			} else {
				// :shrug:
				moduleID = ModuleID(b.base.toPath(tspath.ResolvePath(tspath.GetDirectoryPath(file.FileName()), name)))
			}
		}
		parseModuleDeclaration(decl, file, moduleID, nodeModulesDirectory, packageName, getChecker, &exports)
	}
	return exports
}

func parseExport(name string, symbol *ast.Symbol, moduleID ModuleID, file *ast.SourceFile, nodeModulesDirectory tspath.Path, packageName string, getChecker func() (*checker.Checker, func()), exports *[]*RawExport) {
	if shouldIgnoreSymbol(symbol) {
		return
	}

	if name == ast.InternalSymbolNameExportStar {
		checker, release := getChecker()
		defer release()
		allExports := checker.GetExportsOfModule(symbol.Parent)
		// allExports includes named exports from the file that will be processed separately;
		// we want to add only the ones that come from the star
		for name, namedExport := range symbol.Parent.Exports {
			if name != ast.InternalSymbolNameExportStar {
				idx := slices.Index(allExports, namedExport)
				if idx >= 0 || shouldIgnoreSymbol(namedExport) {
					allExports = slices.Delete(allExports, idx, idx+1)
				}
			}
		}

		*exports = slices.Grow(*exports, len(allExports))
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
				PackageName:                packageName,
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
		PackageName:          packageName,
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

			if syntax == ExportSyntaxEquals && targetSymbol.Flags&ast.SymbolFlagsNamespace != 0 {
				// !!! what is the right boundary for recursion? we never need to expand named exports into another level of named
				//     exports, but for getting flags/kinds, we should resolve each named export as an alias
				*exports = slices.Grow(*exports, len(targetSymbol.Exports))
				for _, namedExport := range targetSymbol.Exports {
					resolved := checker.SkipAlias(namedExport)
					if shouldIgnoreSymbol(resolved) {
						continue
					}
					*exports = append(*exports, &RawExport{
						ExportID: ExportID{
							ExportName: name,
							ModuleID:   moduleID,
						},
						// !!! decide what this means for reexports
						Syntax: ExportSyntaxNamed,
						Flags:  resolved.Flags,
						Target: ExportID{
							ExportName: namedExport.Name,
							// !!!
							ModuleID: ModuleID(ast.GetSourceFileOfNode(resolved.Declarations[0]).Path()),
						},
						ScriptElementKind:          lsutil.GetSymbolKind(checker, resolved, resolved.Declarations[0]),
						ScriptElementKindModifiers: lsutil.GetSymbolModifiers(checker, resolved),
						FileName:                   file.FileName(),
						Path:                       file.Path(),
						NodeModulesDirectory:       nodeModulesDirectory,
						PackageName:                packageName,
					})
				}
			}
		}
		release()
	} else {
		export.ScriptElementKind = lsutil.GetSymbolKind(nil, symbol, symbol.Declarations[0])
		export.ScriptElementKindModifiers = lsutil.GetSymbolModifiers(nil, symbol)
	}

	*exports = append(*exports, export)
}

func parseModuleDeclaration(decl *ast.ModuleDeclaration, file *ast.SourceFile, moduleID ModuleID, nodeModulesDirectory tspath.Path, packageName string, getChecker func() (*checker.Checker, func()), exports *[]*RawExport) {
	for name, symbol := range decl.Symbol.Exports {
		parseExport(name, symbol, moduleID, file, nodeModulesDirectory, packageName, getChecker, exports)
	}
}

func shouldIgnoreSymbol(symbol *ast.Symbol) bool {
	if symbol.Flags&ast.SymbolFlagsPrototype != 0 {
		return true
	}
	return false
}
