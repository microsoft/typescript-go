package ls

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
)

type ImpExpKind int32

const (
	ImpExpKindUnknown ImpExpKind = iota
	ImpExpKindImport
	ImpExpKindExport
)

type ExportKind int32

const (
	ExportKindDefault ExportKind = iota
	ExportKindNamed
	ExportKindExportEquals
)

type ImportExportSymbol struct {
	kind       ImpExpKind
	symbol     *ast.Symbol
	exportInfo *ExportInfo
}

type ExportInfo struct {
	exportingModuleSymbol *ast.Symbol
	exportKind            ExportKind
}

func getImportOrExportSymbol(node *ast.Node, symbol *ast.Symbol, checker *checker.Checker, comingFromExport bool) *ImportExportSymbol {
	exportInfo := func(symbol *ast.Symbol, kind ExportKind) *ImportExportSymbol {
		if exportInfo := getExportInfo(symbol, kind, checker); exportInfo != nil {
			return &ImportExportSymbol{
				kind:       ImpExpKindExport,
				symbol:     symbol,
				exportInfo: exportInfo,
			}
		}
		return nil
	}

	getExport := func() *ImportExportSymbol {
		getExportAssignmentExport := func(ex *ast.Node) *ImportExportSymbol {
			// Get the symbol for the `export =` node; its parent is the module it's the export of.
			if ex.Symbol().Parent == nil {
				return nil
			}
			exportKind := core.IfElse(ex.AsExportAssignment().IsExportEquals, ExportKindExportEquals, ExportKindDefault)
			return &ImportExportSymbol{
				kind:   ImpExpKindExport,
				symbol: symbol,
				exportInfo: &ExportInfo{
					exportingModuleSymbol: ex.Symbol().Parent,
					exportKind:            exportKind,
				},
			}
		}

		// Not meant for use with export specifiers or export assignment.
		getExportKindForDeclaration := func(node *ast.Node) ExportKind {
			if ast.HasSyntacticModifier(node, ast.ModifierFlagsDefault) {
				return ExportKindDefault
			}
			return ExportKindNamed
		}

		getSpecialPropertyExport := func(node *ast.Node, useLhsSymbol bool) *ImportExportSymbol {
			var kind ExportKind
			switch ast.GetAssignmentDeclarationKind(node.AsBinaryExpression()) {
			case ast.JSDeclarationKindExportsProperty:
				kind = ExportKindNamed
			case ast.JSDeclarationKindModuleExports:
				kind = ExportKindExportEquals
			default:
				return nil
			}
			sym := symbol
			if useLhsSymbol {
				sym = checker.GetSymbolAtLocation(ast.GetElementOrPropertyAccessName(node.AsBinaryExpression().Left))
			}
			if sym == nil {
				return nil
			}
			return exportInfo(sym, kind)
		}

		parent := node.Parent
		grandparent := parent.Parent
		if symbol.ExportSymbol != nil {
			if ast.IsPropertyAccessExpression(parent) {
				// When accessing an export of a JS module, there's no alias. The symbol will still be flagged as an export even though we're at the use.
				// So check that we are at the declaration.
				if ast.IsBinaryExpression(grandparent) && slices.Contains(symbol.Declarations, parent) {
					return getSpecialPropertyExport(grandparent, false /*useLhsSymbol*/)
				}
				return nil
			}
			return exportInfo(symbol.ExportSymbol, getExportKindForDeclaration(parent))
		} else {
			exportNode := getExportNode(parent, node)
			switch {
			case exportNode != nil && ast.HasSyntacticModifier(exportNode, ast.ModifierFlagsExport):
				if ast.IsImportEqualsDeclaration(exportNode) && exportNode.AsImportEqualsDeclaration().ModuleReference == node {
					// We're at `Y` in `export import X = Y`. This is not the exported symbol, the left-hand-side is. So treat this as an import statement.
					if comingFromExport {
						return nil
					}
					lhsSymbol := checker.GetSymbolAtLocation(exportNode.Name())
					return &ImportExportSymbol{
						kind:   ImpExpKindImport,
						symbol: lhsSymbol,
					}
				}
				return exportInfo(symbol, getExportKindForDeclaration(exportNode))
			case ast.IsNamespaceExport(parent):
				return exportInfo(symbol, ExportKindNamed)
			case ast.IsExportAssignment(parent):
				return getExportAssignmentExport(parent)
			case ast.IsExportAssignment(grandparent):
				return getExportAssignmentExport(grandparent)
			case ast.IsBinaryExpression(parent):
				return getSpecialPropertyExport(parent, true /*useLhsSymbol*/)
			case ast.IsBinaryExpression(grandparent):
				return getSpecialPropertyExport(grandparent, true /*useLhsSymbol*/)
			case ast.IsJSDocTypedefTag(parent) || ast.IsJSDocCallbackTag(parent):
				return exportInfo(symbol, ExportKindNamed)
			}
		}
		return nil
	}

	getImport := func() *ImportExportSymbol {
		if !isNodeImport(node) {
			return nil
		}
		// A symbol being imported is always an alias. So get what that aliases to find the local symbol.
		importedSymbol := checker.GetImmediateAliasedSymbol(symbol)
		if importedSymbol == nil {
			return nil
		}
		// Search on the local symbol in the exporting module, not the exported symbol.
		importedSymbol = skipExportSpecifierSymbol(importedSymbol, checker)
		// Similarly, skip past the symbol for 'export ='
		if importedSymbol.Name == "export=" {
			importedSymbol = getExportEqualsLocalSymbol(importedSymbol, checker)
			if importedSymbol == nil {
				return nil
			}
		}
		// If the import has a different name than the export, do not continue searching.
		// If `importedName` is undefined, do continue searching as the export is anonymous.
		// (All imports returned from this function will be ignored anyway if we are in rename and this is a not a named export.)
		importedName := symbolNameNoDefault(importedSymbol)
		if importedName == "" || importedName == ast.InternalSymbolNameDefault || importedName == symbol.Name {
			return &ImportExportSymbol{
				kind:   ImpExpKindImport,
				symbol: importedSymbol,
			}
		}
		return nil
	}

	result := getExport()
	if result == nil && !comingFromExport {
		result = getImport()
	}
	return result
}

func getExportInfo(exportSymbol *ast.Symbol, exportKind ExportKind, c *checker.Checker) *ExportInfo {
	// Parent can be nil if an `export` is not at the top-level (which is a compile error).
	if exportSymbol.Parent != nil {
		exportingModuleSymbol := c.GetMergedSymbol(exportSymbol.Parent)
		// `export` may appear in a namespace. In that case, just rely on global search.
		if checker.IsExternalModuleSymbol(exportingModuleSymbol) {
			return &ExportInfo{
				exportingModuleSymbol: exportingModuleSymbol,
				exportKind:            exportKind,
			}
		}
	}
	return nil
}

// If a reference is a class expression, the exported node would be its parent.
// If a reference is a variable declaration, the exported node would be the variable statement.
func getExportNode(parent *ast.Node, node *ast.Node) *ast.Node {
	var declaration *ast.Node
	switch {
	case ast.IsVariableDeclaration(parent):
		declaration = parent
	case ast.IsBindingElement(parent):
		declaration = ast.WalkUpBindingElementsAndPatterns(parent)
	}
	if declaration != nil {
		if parent.Name() == node && !ast.IsCatchClause(declaration.Parent) && ast.IsVariableStatement(declaration.Parent.Parent) {
			return declaration.Parent.Parent
		}
		return nil
	}
	return parent
}

func isNodeImport(node *ast.Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindImportEqualsDeclaration:
		return parent.Name() == node && isExternalModuleImportEquals(parent)
	case ast.KindImportSpecifier:
		// For a rename import `{ foo as bar }`, don't search for the imported symbol. Just find local uses of `bar`.
		return parent.PropertyName() == nil
	case ast.KindImportClause, ast.KindNamespaceImport:
		debug.Assert(parent.Name() == node)
		return true
	case ast.KindBindingElement:
		return ast.IsInJSFile(node) && ast.IsVariableDeclarationInitializedToRequire(parent.Parent.Parent)
	}
	return false
}

func isExternalModuleImportEquals(node *ast.Node) bool {
	moduleReference := node.AsImportEqualsDeclaration().ModuleReference
	return ast.IsExternalModuleReference(moduleReference) && moduleReference.Expression().Kind == ast.KindStringLiteral
}

// If at an export specifier, go to the symbol it refers to. */
func skipExportSpecifierSymbol(symbol *ast.Symbol, checker *checker.Checker) *ast.Symbol {
	// For `export { foo } from './bar", there's nothing to skip, because it does not create a new alias. But `export { foo } does.
	for _, declaration := range symbol.Declarations {
		switch {
		case ast.IsExportSpecifier(declaration) && declaration.PropertyName() == nil && declaration.Parent.Parent.ModuleSpecifier() == nil:
			return core.OrElse(checker.GetExportSpecifierLocalTargetSymbol(declaration), symbol)
		case ast.IsPropertyAccessExpression(declaration) && ast.IsModuleExportsAccessExpression(declaration.Expression()) && !ast.IsPrivateIdentifier(declaration.Name()):
			// Export of form 'module.exports.propName = expr';
			return checker.GetSymbolAtLocation(declaration)
		case ast.IsShorthandPropertyAssignment(declaration) && ast.IsBinaryExpression(declaration.Parent.Parent) && ast.GetAssignmentDeclarationKind(declaration.Parent.Parent.AsBinaryExpression()) == ast.JSDeclarationKindModuleExports:
			return checker.GetExportSpecifierLocalTargetSymbol(declaration.Name())
		}
	}
	return symbol
}

func getExportEqualsLocalSymbol(importedSymbol *ast.Symbol, checker *checker.Checker) *ast.Symbol {
	if importedSymbol.Flags&ast.SymbolFlagsAlias != 0 {
		return checker.GetImmediateAliasedSymbol(importedSymbol)
	}
	decl := debug.CheckDefined(importedSymbol.ValueDeclaration)
	switch {
	case ast.IsExportAssignment(decl):
		return decl.Expression().Symbol()
	case ast.IsBinaryExpression(decl):
		return decl.AsBinaryExpression().Right.Symbol()
	case ast.IsSourceFile(decl):
		return decl.Symbol()
	}
	return nil
}

func symbolNameNoDefault(symbol *ast.Symbol) string {
	if symbol.Name != ast.InternalSymbolNameDefault {
		return symbol.Name
	}
	for _, decl := range symbol.Declarations {
		name := ast.GetNameOfDeclaration(decl)
		if name != nil && ast.IsIdentifier(name) {
			return name.Text()
		}
	}
	return ""
}
