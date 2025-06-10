package transformers

import (
	"maps"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

type JSXTransformer struct {
	Transformer
	compilerOptions *core.CompilerOptions
	parentNode      *ast.Node
	currentNode     *ast.Node

	importSpecifier                string
	filenameDeclaration            *ast.Node
	utilizedImplicitRuntimeImports map[string]map[string]*ast.Node
	inJsxChild                     bool

	currentSourceFile *ast.SourceFile
}

func NewJSXTransformer(emitContext *printer.EmitContext, compilerOptions *core.CompilerOptions) *Transformer {
	tx := &JSXTransformer{
		compilerOptions: compilerOptions,
	}
	return tx.NewTransformer(tx.visit, emitContext)
}

func (tx *JSXTransformer) getJsxFactoryCalleePrimitive(isStaticChildren bool) string {
	if tx.compilerOptions.Jsx == core.JsxEmitReactJSXDev {
		return "jsxDEV"
	}
	if isStaticChildren {
		return "jsxs"
	}
	return "jsx"
}

func (tx *JSXTransformer) getJsxFactoryCallee(isStaticChildren bool) *ast.Node {
	t := tx.getJsxFactoryCalleePrimitive(isStaticChildren)
	return tx.getImplicitImportForName(t)
}

func (tx *JSXTransformer) getImplicitJsxFragmentReference() *ast.Node {
	return tx.getImplicitImportForName("Fragment")
}

func (tx *JSXTransformer) getImplicitImportForName(name string) *ast.Node {
	importSource := tx.importSpecifier
	if name != "createElement" {
		importSource = ast.GetJSXRuntimeImport(importSource, tx.compilerOptions)
	}
	existing, ok := tx.utilizedImplicitRuntimeImports[importSource]
	if ok {
		elem, ok := existing[name]
		if ok {
			return elem.AsImportSpecifier().Name()
		}
	} else {
		tx.utilizedImplicitRuntimeImports[importSource] = make(map[string]*ast.Node)
	}

	generatedName := tx.factory.NewUniqueNameEx("_"+name, printer.AutoGenerateOptions{
		Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel | printer.GeneratedIdentifierFlagsAllowNameSubstitution,
	})
	specifier := tx.factory.NewImportSpecifier(false, tx.factory.NewIdentifier(name), generatedName)
	// setIdentifierGeneratedImportReference(generatedName, specifier); // !!!
	tx.utilizedImplicitRuntimeImports[importSource][name] = specifier
	return specifier
}

func (tx *JSXTransformer) setParentAndCurrentNode(parent *ast.Node, current *ast.Node) {
	tx.parentNode = parent
	tx.currentNode = current
}

func (tx *JSXTransformer) setInChild(v bool) {
	tx.inJsxChild = v
}

func (tx *JSXTransformer) visit(node *ast.Node) *ast.Node {
	if node.SubtreeFacts()&ast.SubtreeContainsJsx == 0 {
		return node
	}
	oldParent := tx.parentNode
	tx.setParentAndCurrentNode(tx.currentNode, node)
	defer tx.setParentAndCurrentNode(oldParent, tx.parentNode)
	switch node.Kind {
	case ast.KindSourceFile:
		tx.setInChild(false)
		return tx.visitSourceFile(node.AsSourceFile())
	case ast.KindJsxElement:
		return tx.visitJsxElement(node.AsJsxElement())
	case ast.KindJsxSelfClosingElement:
		return tx.visitJsxSelfClosingElement(node.AsJsxSelfClosingElement())
	case ast.KindJsxFragment:
		return tx.visitJsxFragment(node.AsJsxFragment())
	case ast.KindJsxOpeningElement:
		tx.setInChild(false)
		return tx.visitJsxOpeningElement(node.AsJsxOpeningElement())
	case ast.KindJsxOpeningFragment:
		tx.setInChild(false)
		return tx.visitJsxOpeningFragment(node.AsJsxOpeningFragment())
	case ast.KindJsxText:
		tx.setInChild(false)
		return tx.visitJsxText(node.AsJsxText())
	case ast.KindJsxExpression:
		tx.setInChild(false)
		return tx.visitJsxExpression(node.AsJsxExpression())
	}
	tx.setInChild(false)
	return tx.visitor.VisitEachChild(node) // by default, do nothing
}

/**
 * The react jsx/jsxs transform falls back to `createElement` when an explicit `key` argument comes after a spread
 */
func hasKeyAfterPropsSpread(node *ast.Node) bool {
	spread := false
	for _, elem := range node.Attributes().Properties() {
		if ast.IsJsxSpreadAttribute(elem) && (!ast.IsObjectLiteralExpression(elem.Expression()) || core.Some(elem.Expression().Properties(), ast.IsSpreadAssignment)) {
			spread = true
		} else if spread && ast.IsJsxAttribute(elem) && ast.IsIdentifier(elem.Name()) && elem.Name().AsIdentifier().Text == "key" {
			return true
		}
	}
	return false
}

func (tx *JSXTransformer) shouldUseCreateElement(node *ast.Node) bool {
	return len(tx.importSpecifier) == 0 || hasKeyAfterPropsSpread(node)
}

func insertStatementAfterPrologue[T any](to []*ast.Node, statement *ast.Node, isPrologueDirective func(callee T, node *ast.Node) bool, callee T) []*ast.Node {
	if statement == nil {
		return to
	}
	statementIdx := 0
	// skip all prologue directives to insert at the correct position
	for ; statementIdx < len(to); statementIdx++ {
		if !isPrologueDirective(callee, to[statementIdx]) {
			break
		}
	}
	return slices.Insert(to, statementIdx, statement)
}

func (tx *JSXTransformer) isAnyPrologueDirective(node *ast.Node) bool {
	return ast.IsPrologueDirective(node) || (tx.emitContext.EmitFlags(node)&printer.EFCustomPrologue != 0)
}

func (tx *JSXTransformer) insertStatementAfterCustomPrologue(to []*ast.Node, statement *ast.Node) []*ast.Node {
	return insertStatementAfterPrologue(to, statement, (*JSXTransformer).isAnyPrologueDirective, tx)
}

func sortByImportDeclarationSource(a *ast.Node, b *ast.Node) int {
	return stringutil.CompareStringsCaseSensitive(a.AsImportDeclaration().ModuleSpecifier.AsStringLiteral().Text, b.AsImportDeclaration().ModuleSpecifier.AsStringLiteral().Text)
}

func getSpecifierOfRequireCall(s *ast.Node) string {
	return s.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes[0].AsVariableDeclaration().Initializer.AsCallExpression().Arguments.Nodes[0].AsStringLiteral().Text
}

func sortByRequireSource(a *ast.Node, b *ast.Node) int {
	return stringutil.CompareStringsCaseSensitive(getSpecifierOfRequireCall(a), getSpecifierOfRequireCall(b))
}

func sortImportSpecifiers(a *ast.Node, b *ast.Node) int {
	res := stringutil.CompareStringsCaseSensitive(a.AsImportSpecifier().PropertyName.Text(), b.AsImportSpecifier().PropertyName.Text())
	if res != 0 {
		return res
	}
	return stringutil.CompareStringsCaseSensitive(a.AsImportSpecifier().Name().AsIdentifier().Text, b.AsImportSpecifier().Name().AsIdentifier().Text)
}

func getSortedSpecifiers(m map[string]*ast.Node) []*ast.Node {
	res := slices.Collect(maps.Values(m))
	slices.SortFunc(res, sortImportSpecifiers)
	return res
}

func (tx *JSXTransformer) visitSourceFile(file *ast.SourceFile) *ast.Node {
	if file.IsDeclarationFile {
		return file.AsNode()
	}

	tx.currentSourceFile = file
	tx.importSpecifier = ast.GetJSXImplicitImportBase(tx.compilerOptions, file)
	tx.filenameDeclaration = nil
	tx.utilizedImplicitRuntimeImports = make(map[string]map[string]*ast.Node)

	visited := tx.visitor.VisitEachChild(file.AsNode())
	tx.emitContext.AddEmitHelper(visited.AsNode(), tx.emitContext.ReadEmitHelpers()...)
	statements := visited.Statements()
	statementsUpdated := false
	if tx.filenameDeclaration != nil {
		statements = tx.insertStatementAfterCustomPrologue(statements, tx.factory.NewVariableStatement(nil, tx.factory.NewVariableDeclarationList(
			ast.NodeFlagsConst,
			tx.factory.NewNodeList([]*ast.Node{tx.filenameDeclaration}),
		)))
		statementsUpdated = true
	}

	if len(tx.utilizedImplicitRuntimeImports) > 0 {
		// A key difference from strada is that these imports are sorted in corsa, rather than appearing in a use-defined order
		if ast.IsExternalModule(file) {
			statementsUpdated = true
			newStatements := make([]*ast.Node, 0, len(tx.utilizedImplicitRuntimeImports))
			for importSource, importSpecifiersMap := range tx.utilizedImplicitRuntimeImports {
				s := tx.factory.NewImportDeclaration(
					nil,
					tx.factory.NewImportClause(false, nil, tx.factory.NewNamedImports(tx.factory.NewNodeList(getSortedSpecifiers(importSpecifiersMap)))),
					tx.factory.NewStringLiteral(importSource),
					nil,
				)
				ast.SetParentInChildren(s)
				newStatements = append(newStatements, s)

			}
			slices.SortFunc(newStatements, sortByImportDeclarationSource)
			for _, e := range newStatements {
				statements = tx.insertStatementAfterCustomPrologue(statements, e)
			}
		} else if ast.IsExternalOrCommonJSModule(file) {
			statementsUpdated = true
			newStatements := make([]*ast.Node, 0, len(tx.utilizedImplicitRuntimeImports))
			for importSource, importSpecifiersMap := range tx.utilizedImplicitRuntimeImports {
				sorted := getSortedSpecifiers(importSpecifiersMap)
				asBindingElems := make([]*ast.Node, 0, len(sorted))
				for _, elem := range sorted {
					asBindingElems = append(asBindingElems, tx.factory.NewBindingElement(nil, elem.AsImportSpecifier().PropertyName, elem.AsImportSpecifier().Name(), nil))
				}
				s := tx.factory.NewVariableStatement(nil, tx.factory.NewVariableDeclarationList(ast.NodeFlagsConst, tx.factory.NewNodeList([]*ast.Node{tx.factory.NewVariableDeclaration(
					tx.factory.NewBindingPattern(ast.KindObjectBindingPattern, tx.factory.NewNodeList(asBindingElems)),
					nil,
					nil,
					tx.factory.NewCallExpression(tx.factory.NewIdentifier("require"), nil, nil, tx.factory.NewNodeList([]*ast.Node{tx.factory.NewStringLiteral(importSource)}), ast.NodeFlagsNone),
				)})))
				ast.SetParentInChildren(s)
				newStatements = append(newStatements, s)
			}
			slices.SortFunc(newStatements, sortByRequireSource)
			for _, e := range newStatements {
				statements = tx.insertStatementAfterCustomPrologue(statements, e)
			}
		} else {
			// Do nothing (script file) - consider an error in the checker?
		}
	}

	if statementsUpdated {
		visited = tx.factory.UpdateSourceFile(file, tx.factory.NewNodeList(statements))
	}

	tx.currentSourceFile = nil
	tx.importSpecifier = ""
	tx.filenameDeclaration = nil
	tx.utilizedImplicitRuntimeImports = nil

	return visited
}

func (tx *JSXTransformer) visitJsxElement(element *ast.JsxElement) *ast.Node {
	tagTransform := (*JSXTransformer).visitJsxOpeningLikeElementJSX
	if tx.shouldUseCreateElement(element.AsNode()) {
		tagTransform = (*JSXTransformer).visitJsxOpeningLikeElementCreateElement
	}
	return tagTransform(tx, element)
}

func (tx *JSXTransformer) visitJsxSelfClosingElement(element *ast.JsxSelfClosingElement) *ast.Node {
	tagTransform := (*JSXTransformer).visitJsxOpeningLikeElementJSX
	if tx.shouldUseCreateElement(element.AsNode()) {
		tagTransform = (*JSXTransformer).visitJsxOpeningLikeElementCreateElement
	}
	return tagTransform(tx, element)
}

func (tx *JSXTransformer) visitJsxFragment(fragment *ast.JsxFragment) *ast.Node {
	tagTransform := (*JSXTransformer).visitJsxOpeningFragmentJSX
	if len(tx.importSpecifier) == 0 {
		tagTransform = (*JSXTransformer).visitJsxOpeningFragmentCreateElement
	}
	return tagTransform(tx, fragment)
}

func (tx *JSXTransformer) visitJsxOpeningLikeElementJSX(element *ast.Node, children *ast.NodeList, location *ast.Node) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxOpeningFragmentJSX(fragment *ast.JsxOpeningFragment, children *ast.NodeList, location *ast.Node) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxOpeningLikeElementCreateElement(element *ast.Node, children *ast.NodeList, location *ast.Node) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxOpeningFragmentCreateElement(fragment *ast.JsxOpeningFragment, children *ast.NodeList, location *ast.Node) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxOpeningElement(element *ast.JsxOpeningElement) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxOpeningFragment(fragment *ast.JsxOpeningFragment) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxText(text *ast.JsxText) *ast.Node {
	panic("unimplemented")
}

func (tx *JSXTransformer) visitJsxExpression(expression *ast.JsxExpression) *ast.Node {
	panic("unimplemented")
}
