package autoimport

import (
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type exportExtractor struct {
	nodeModulesDirectory tspath.Path
	packageName          string

	localNameResolver *binder.NameResolver
	moduleResolver    *module.Resolver
	getChecker        func() (*checker.Checker, func())
	toPath            func(fileName string) tspath.Path
}

type checkerLease struct {
	getChecker func() (*checker.Checker, func())
	checker    *checker.Checker
	release    func()
}

func (l *checkerLease) GetChecker() *checker.Checker {
	if l.checker == nil {
		l.checker, l.release = l.getChecker()
	}
	return l.checker
}

func (l *checkerLease) TryChecker() *checker.Checker {
	return l.checker
}

func (l *checkerLease) Done() {
	if l.release != nil {
		l.release()
		l.release = nil
	}
}

func (b *registryBuilder) newExportExtractor(nodeModulesDirectory tspath.Path, packageName string, getChecker func() (*checker.Checker, func())) *exportExtractor {
	return &exportExtractor{
		nodeModulesDirectory: nodeModulesDirectory,
		packageName:          packageName,
		moduleResolver:       b.resolver,
		getChecker:           getChecker,
		toPath:               b.base.toPath,
		localNameResolver: &binder.NameResolver{
			CompilerOptions: core.EmptyCompilerOptions,
		},
	}
}

func (e *exportExtractor) extractFromFile(file *ast.SourceFile) []*Export {
	if file.Symbol != nil {
		return e.extractFromModule(file)
	}
	if len(file.AmbientModuleNames) > 0 {
		moduleDeclarations := core.Filter(file.Statements.Nodes, ast.IsModuleWithStringLiteralName)
		var exportCount int
		for _, decl := range moduleDeclarations {
			exportCount += len(decl.AsModuleDeclaration().Symbol.Exports)
		}
		exports := make([]*Export, 0, exportCount)
		for _, decl := range moduleDeclarations {
			e.extractFromModuleDeclaration(decl.AsModuleDeclaration(), file, ModuleID(decl.Name().Text()), &exports)
		}
		return exports
	}
	return nil
}

func (e *exportExtractor) extractFromModule(file *ast.SourceFile) []*Export {
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
	exports := make([]*Export, 0, len(file.Symbol.Exports)+augmentationExportCount)
	for name, symbol := range file.Symbol.Exports {
		e.extractFromSymbol(name, symbol, ModuleID(file.Path()), file, &exports)
	}
	for _, decl := range moduleAugmentations {
		name := decl.Name().AsStringLiteral().Text
		moduleID := ModuleID(name)
		if tspath.IsExternalModuleNameRelative(name) {
			// !!! need to resolve non-relative names in separate pass
			if resolved, _ := e.moduleResolver.ResolveModuleName(name, file.FileName(), core.ModuleKindCommonJS, nil); resolved.IsResolved() {
				moduleID = ModuleID(e.toPath(resolved.ResolvedFileName))
			} else {
				// :shrug:
				moduleID = ModuleID(e.toPath(tspath.ResolvePath(tspath.GetDirectoryPath(file.FileName()), name)))
			}
		}
		e.extractFromModuleDeclaration(decl, file, moduleID, &exports)
	}
	return exports
}

func (e *exportExtractor) extractFromModuleDeclaration(decl *ast.ModuleDeclaration, file *ast.SourceFile, moduleID ModuleID, exports *[]*Export) {
	for name, symbol := range decl.Symbol.Exports {
		e.extractFromSymbol(name, symbol, moduleID, file, exports)
	}
}

func (e *exportExtractor) extractFromSymbol(name string, symbol *ast.Symbol, moduleID ModuleID, file *ast.SourceFile, exports *[]*Export) {
	if shouldIgnoreSymbol(symbol) {
		return
	}

	if name == ast.InternalSymbolNameExportStar {
		checkerLease := &checkerLease{getChecker: e.getChecker}
		defer checkerLease.Done()
		checker := checkerLease.GetChecker()
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
			export, _ := e.createExport(reexportedSymbol, moduleID, ExportSyntaxStar, file, checkerLease)
			if export != nil {
				export.through = ast.InternalSymbolNameExportStar
				*exports = append(*exports, export)
			}
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
		case ast.KindJSExportAssignment:
			declSyntax = ExportSyntaxCommonJSModuleExports
		case ast.KindCommonJSExport:
			declSyntax = ExportSyntaxCommonJSExportsProperty
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

	checkerLease := &checkerLease{getChecker: e.getChecker}
	export, target := e.createExport(symbol, moduleID, syntax, file, checkerLease)
	defer checkerLease.Done()
	if export == nil {
		return
	}

	export.localName = localName
	*exports = append(*exports, export)

	if target != nil {
		if syntax == ExportSyntaxEquals && target.Flags&ast.SymbolFlagsNamespace != 0 {
			// !!! what is the right boundary for recursion? we never need to expand named exports into another level of named
			//     exports, but for getting flags/kinds, we should resolve each named export as an alias
			*exports = slices.Grow(*exports, len(target.Exports))
			for _, namedExport := range target.Exports {
				export, _ := e.createExport(namedExport, moduleID, syntax, file, checkerLease)
				if export != nil {
					export.through = name
					*exports = append(*exports, export)
				}
			}
		}
	} else if syntax == ExportSyntaxCommonJSModuleExports {
		expression := symbol.Declarations[0].AsExportAssignment().Expression
		if expression.Kind == ast.KindObjectLiteralExpression {
			// what is actually desirable here? I think it would be reasonable to only treat these as exports
			// if *every* property is a shorthand property or identifier: identifier
			// At least, it would be sketchy if there were any methods, computed properties...
			*exports = slices.Grow(*exports, len(expression.AsObjectLiteralExpression().Properties.Nodes))
			for _, prop := range expression.AsObjectLiteralExpression().Properties.Nodes {
				if ast.IsShorthandPropertyAssignment(prop) || ast.IsPropertyAssignment(prop) && prop.AsPropertyAssignment().Name().Kind == ast.KindIdentifier {
					export, _ := e.createExport(expression.Symbol().Members[prop.Name().Text()], moduleID, syntax, file, checkerLease)
					if export != nil {
						export.through = name
						*exports = append(*exports, export)
					}
				}
			}
		}
	}
}

// createExport creates an Export for the given symbol, returning the Export and the target symbol if the export is an alias.
func (e *exportExtractor) createExport(symbol *ast.Symbol, moduleID ModuleID, syntax ExportSyntax, file *ast.SourceFile, checkerLease *checkerLease) (*Export, *ast.Symbol) {
	if shouldIgnoreSymbol(symbol) {
		return nil, nil
	}

	export := &Export{
		ExportID: ExportID{
			ExportName: symbol.Name,
			ModuleID:   moduleID,
		},
		Syntax:               syntax,
		Flags:                symbol.Flags,
		Path:                 file.Path(),
		NodeModulesDirectory: e.nodeModulesDirectory,
		PackageName:          e.packageName,
	}

	var targetSymbol *ast.Symbol
	if symbol.Flags&ast.SymbolFlagsAlias != 0 {
		checker := checkerLease.GetChecker()
		// !!! try localNameResolver first?
		targetSymbol = checker.GetAliasedSymbol(symbol)
		if !checker.IsUnknownSymbol(targetSymbol) {
			var decl *ast.Node
			if len(targetSymbol.Declarations) > 0 {
				decl = targetSymbol.Declarations[0]
			} else if targetSymbol.CheckFlags&ast.CheckFlagsMapped != 0 {
				if mappedDecl := checker.GetMappedTypeSymbolOfProperty(targetSymbol); mappedDecl != nil && len(mappedDecl.Declarations) > 0 {
					decl = mappedDecl.Declarations[0]
				}
			}
			if decl == nil {
				// !!! consider GetImmediateAliasedSymbol to go as far as we can
				decl = symbol.Declarations[0]
			}
			if decl == nil {
				panic("no declaration for aliased symbol")
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
	} else {
		export.ScriptElementKind = lsutil.GetSymbolKind(checkerLease.TryChecker(), symbol, symbol.Declarations[0])
		export.ScriptElementKindModifiers = lsutil.GetSymbolModifiers(checkerLease.TryChecker(), symbol)
	}

	return export, targetSymbol
}

func shouldIgnoreSymbol(symbol *ast.Symbol) bool {
	if symbol.Flags&ast.SymbolFlagsPrototype != 0 {
		return true
	}
	return false
}
