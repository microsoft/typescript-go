package parser

import (
	"github.com/microsoft/typescript-go/internal/ast"
)

type jsDeclarationKind int

const (
	jsDeclarationKindNone jsDeclarationKind = iota
	/// module.exports = expr
	jsDeclarationKindModuleExports
	/// exports.name = expr
	/// module.exports.name = expr
	jsDeclarationKindExportsProperty
	/// className.prototype.name = expr
	jsDeclarationKindPrototypeProperty
	/// this.name = expr
	jsDeclarationKindThisProperty
	// F.name = expr
	jsDeclarationKindProperty
)

func (p *Parser) withCommonJS(node *ast.Node) {
	if node.Kind == ast.KindVariableStatement && node.AsVariableStatement().DeclarationList != nil {
		for _, declaration := range node.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes {
			if isRequireCall(declaration.Initializer()) {
				ref := p.factory.NewExternalModuleReference(declaration.Initializer())
				ref.Flags = ast.NodeFlagsReparsed
				ref.Loc = declaration.Initializer().Loc
				importeq := p.factory.NewJSImportEqualsDeclaration(nil /*modifiers*/, declaration.Name(), ref)
				importeq.Flags = ast.NodeFlagsReparsed
				importeq.Loc = declaration.Loc
				p.reparseList = append(p.reparseList, importeq)
			}
		}
		return
	}
	if node.Kind != ast.KindExpressionStatement || node.AsExpressionStatement().Expression.Kind != ast.KindBinaryExpression {
		return
	}
	bin := node.AsExpressionStatement().Expression.AsBinaryExpression()
	kind := getAssignmentDeclarationKind(bin)
	switch kind {
	case jsDeclarationKindModuleExports:
		export := p.factory.NewJSExportAssignment(nil /*modifiers*/, bin.Right)
		export.Flags = ast.NodeFlagsReparsed
		export.Loc = bin.Loc
		p.reparseList = append(p.reparseList, export)
	}

	// TODO: Duplicate all the places that reference either kind, or either predicate, or .. check the others too.
	// TODO: mark the file as a (commonjs) module if either is found
	// TODO: maybe remove modifiers, isTypeOnly, isExportEquals, since they're not used (but makes identical handling way easier)
}

func getAssignmentDeclarationKind(bin *ast.BinaryExpression) jsDeclarationKind {
	if bin.OperatorToken.Kind != ast.KindEqualsToken || !ast.IsAccessExpression(bin.Left) {
		return jsDeclarationKindNone
	}
	if isModuleExportsAccessExpression(bin.Left) {
		return jsDeclarationKindModuleExports
	}
	// !!! module.exports property, this.property, expando.property
	return jsDeclarationKindNone
}

func isModuleExportsAccessExpression(node *ast.Node) bool {
	return (ast.IsPropertyAccessExpression(node) || isLiteralLikeElementAccess(node)) &&
		isModuleIdentifier(node.Expression()) &&
		ast.GetElementOrPropertyAccessName(node) == "exports"
}

func isModuleIdentifier(node *ast.Node) bool {
	return ast.IsIdentifier(node) && node.AsIdentifier().Text == "module"
}

func isLiteralLikeElementAccess(node *ast.Node) bool {
	return node.Kind == ast.KindElementAccessExpression && ast.IsStringOrNumericLiteralLike(node.AsElementAccessExpression().ArgumentExpression)
}

func isRequireCall(node *ast.Node) bool {
	if node.Kind != ast.KindCallExpression {
		return false
	}
	expr := node.AsCallExpression().Expression
	args := node.AsCallExpression().Arguments
	return expr.Name().Kind == ast.KindIdentifier && expr.Name().AsIdentifier().Text == "require" &&
		len(args.Nodes) == 1 && ast.IsStringLiteralLike(args.Nodes[0])
}
