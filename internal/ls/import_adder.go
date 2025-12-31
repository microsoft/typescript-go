package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

type ImportAdder interface {
	HasFixes() bool
	AddImportFromExportedSymbol(symbol *ast.Symbol, isValidTypeOnlyUseSite bool)
	WriteFixes() // !!!
}

// !!!
type importAdder struct {
	// Context
	ctx     context.Context
	ls      *LanguageService
	program *compiler.Program
	// !!! do we need to use a checker?
	checker *checker.Checker
	file    *ast.SourceFile

	// State
	addToNamespace []*ImportFix
	importType     []*ImportFix
	addToExisting  map[*ast.ImportClauseOrBindingPattern]*addToExistingState
}

type addToExistingState struct {
	importClauseOrBindingPattern *ast.ImportClauseOrBindingPattern
	defaultImport                any // !!!
	namedImports                     // !!!
}

func NewImportAdder() ImportAdder {
	return &importAdder{}
}

func (adder *importAdder) HasFixes() bool {
	// !!!
}

func (adder *importAdder) AddImportFromExportedSymbol(exportedSymbol *ast.Symbol, isValidTypeOnlyUseSite bool) {
	moduleSymbol := debug.CheckDefined(exportedSymbol.Parent, "Expected exported symbol to have module symbol as parent")
	symbolName := getNameForExportedSymbol(exportedSymbol, false /*preferCapitalized*/)
	symbol := adder.checker.GetMergedSymbol(adder.checker.SkipAlias(exportedSymbol))
	exportInfo := getAllExportInfoForSymbol(adder.file, symbol, symbolName, moduleSymbol) // !!! args
	if exportInfo == nil {
		// If no exportInfo is found, this means export could not be resolved when we have filtered for autoImportFileExcludePatterns,
		//     so we should not generate an import.
		debug.Assert(len(adder.ls.UserPreferences().AutoImportFileExcludePatterns) > 0)
		return
	}
	useRequire := shouldUseRequire(adder.file, adder.program)
	fix := getImportFixForSymbol(adder.file, exportInfo) // !!!
	if fix != nil {
		// !!! referenceImport
		localName := symbolName
		var addAsTypeOnly AddAsTypeOnly
		var propertyName string
		// !!! referenceImport
		if exportedSymbol.Name != localName {
			// checks if the symbol was aliased at the referenced import
			propertyName = exportedSymbol.Name
		}
		if addAsTypeOnly != 0 {
			fix.addAsTypeOnly = addAsTypeOnly
		}
		if propertyName != "" {
			fix.propertyName = propertyName
		}
		addImport(fix) // !!! HERE
	}
}

func (adder *importAdder) addImport(fix *ImportFix) {
	switch fix.kind {
	case ImportFixKindUseNamespace:
		adder.addToNamespace = append(adder.addToNamespace, fix)
	case ImportFixKindJsdocTypeImport:
		adder.importType = append(adder.importType, fix)
	case ImportFixKindAddToExisting:
		entry := adder.addToExisting[importClauseOrBindingPattern]
	}
}

func (adder *importAdder) WriteFixes() {
	// !!!
}

func getAllExportInfoForSymbol(
	importingFile *ast.SourceFile,
	symbol *ast.Symbol,
	symbolName string,
	moduleSymbol *ast.Symbol, // !!! other params
) []*SymbolExportInfo {
	// !!! getExportInfoMap
	// !!!
}

func typeToAutoImportableTypeNode(
	c *checker.Checker,
	importAdder ImportAdder,
	t *checker.Type,
	contextNode *ast.Node, // !!! flags
) *ast.TypeNode {
	idToSymbol := make(map[*ast.IdentifierNode]*ast.Symbol)
	typeNode := c.TypeToTypeNode(t, contextNode, nodebuilder.FlagsNone, idToSymbol)
	if typeNode == nil {
		return nil
	}
	return typeNodeToAutoImportableTypeNode(typeNode, importAdder, idToSymbol)
}

func typeNodeToAutoImportableTypeNode(
	typeNode *ast.TypeNode,
	importAdder ImportAdder,
	idToSymbol map[*ast.IdentifierNode]*ast.Symbol,
) *ast.TypeNode {
	referenceTypeNode, importableSymbols := tryGetAutoImportableReferenceFromTypeNode(typeNode, idToSymbol)
	if referenceTypeNode != nil {
		importSymbols(importAdder, importableSymbols)
		typeNode = referenceTypeNode
	}

	// Ensure nodes are fresh so they can have different positions when going through formatting.
	// !!! we may not need this if we can disable type reuse when producing the type node?
	return getSynthesizedTypeNode(typeNode)
}

func importSymbols(importAdder ImportAdder, symbols []*ast.Symbol) {
	for _, symbol := range symbols {
		importAdder.AddImportFromExportedSymbol(symbol, true /*isValidTypeOnlyUseSite*/)
	}
}

// Given a type node containing 'import("./a").SomeType<import("./b").OtherType<...>>',
// returns an equivalent type reference node with any nested ImportTypeNodes also replaced
// with type references, and a list of symbols that must be imported to use the type reference.
func tryGetAutoImportableReferenceFromTypeNode(importTypeNode *ast.TypeNode, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) (*ast.TypeNode, []*ast.Symbol) {
	var symbols []*ast.Symbol
	var visit func(node *ast.Node) *ast.Node
	factory := ast.NewNodeFactory(ast.NodeFactoryHooks{})
	visitor := ast.NewNodeVisitor(visit, factory, ast.NodeVisitorHooks{})
	visit = func(node *ast.Node) *ast.Node {
		if ast.IsLiteralImportTypeNode(node) && node.AsImportTypeNode().Qualifier != nil {
			importTypeNode := node.AsImportTypeNode()
			// Symbol for the left-most thing after the dot
			firstIdentifier := ast.GetFirstIdentifier(importTypeNode.Qualifier)
			// !!! HERE: will this work in Corsa?
			symbol := idToSymbol[firstIdentifier]
			if symbol == nil {
				// if symbol is missing then this doesn't come from a synthesized import type node
				// it has to be an import type node authored by the user and thus it has to be valid
				// it can't refer to reserved internal symbol names and such
				return node.VisitEachChild(visitor)
			}
			name := getNameForExportedSymbol(symbol, false /*preferCapitalized*/)
			var qualifier *ast.EntityName
			if name != firstIdentifier.Text() {
				qualifier = replaceFirstIdentifierOfEntityName(factory, importTypeNode.Qualifier, factory.NewIdentifier(name))
			} else {
				qualifier = importTypeNode.Qualifier
			}
			symbols = append(symbols, symbol)
			typeArguments := visitor.VisitNodes(importTypeNode.TypeArguments)
			return factory.NewTypeReferenceNode(qualifier, typeArguments)
		}
		return visitor.VisitEachChild(node)
	}

	typeNode := visitor.VisitNode(importTypeNode)
	debug.Assert(typeNode == nil || ast.IsTypeNode(typeNode), "expected a type node")
	return typeNode, symbols
}

// If a type checker and multiple files are available, consider using `forEachNameOfDefaultExport`
// instead, which searches for names of re-exported defaults/namespaces in target files.
func getNameForExportedSymbol(symbol *ast.Symbol, preferCapitalized bool) string {
	if symbol.Name == ast.InternalSymbolNameExportEquals || symbol.Name == ast.InternalSymbolNameDefault {
		// Names for default exports:
		// - export default foo => foo
		// - export { foo as default } => foo
		// - export default 0 => filename converted to camelCase
		name := getDefaultLikeExportNameFromDeclaration(symbol)
		if name != "" {
			return name
		}
		debug.AssertIsDefined(symbol.Parent, "Expected exported symbol to have module symbol as parent")
		return moduleSymbolToValidIdentifier(symbol.Parent, preferCapitalized)
	}
	return symbol.Name
}

func replaceFirstIdentifierOfEntityName(factory *ast.NodeFactory, name *ast.EntityName, newIdentifier *ast.IdentifierNode) *ast.EntityName {
	if name.Kind == ast.KindIdentifier {
		return newIdentifier
	}
	return factory.NewQualifiedName(
		replaceFirstIdentifierOfEntityName(factory, name.AsQualifiedName().Left, newIdentifier),
		name.AsQualifiedName().Right,
	)
}
