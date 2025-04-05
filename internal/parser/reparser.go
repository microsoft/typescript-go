package parser

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
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
	if p.scriptKind != core.ScriptKindJS && p.scriptKind != core.ScriptKindJSX {
		return
	}
	if node.Kind != ast.KindExpressionStatement || node.AsExpressionStatement().Expression.Kind != ast.KindBinaryExpression {
		return
	}
	bin := node.AsExpressionStatement().Expression.AsBinaryExpression()
	kind := getAssignmentDeclarationKind(bin)
	switch kind {
	case jsDeclarationKindModuleExports:
		export := p.factory.NewJSExportAssignment(bin.Right)
		export.Flags = ast.NodeFlagsReparsed
		export.Loc = bin.Loc
		p.reparseList = append(p.reparseList, export)
	case jsDeclarationKindExportsProperty:
		nodes := p.nodeSlicePool.NewSlice(1)
		nodes[0] = p.factory.NewModifier(ast.KindExportKeyword)
		nodes[0].Flags = ast.NodeFlagsReparsed
		nodes[0].Loc = bin.Loc
		// TODO: Name can sometimes be a string literal, so downstream code needs to handle this
		export := p.factory.NewCommonJSExport(p.newModifierList(bin.Loc, nodes), ast.GetElementOrPropertyAccessArgumentExpressionOrName(bin.Left), bin.Right, node.AsExpressionStatement().Expression)
		export.Flags = ast.NodeFlagsReparsed
		export.Loc = bin.Loc
		p.reparseList = append(p.reparseList, export)
	}

	// TODO: mark the file as a (commonjs) module if either is found
}

func getAssignmentDeclarationKind(bin *ast.BinaryExpression) jsDeclarationKind {
	if bin.OperatorToken.Kind != ast.KindEqualsToken || !ast.IsAccessExpression(bin.Left) {
		return jsDeclarationKindNone
	}
	if isModuleExportsAccessExpression(bin.Left) {
		return jsDeclarationKindModuleExports
	} else if ast.IsAccessExpression(bin.Left) &&
		(isModuleExportsAccessExpression(bin.Left.Expression()) || isExport(bin.Left.Expression())) &&
		
		(ast.IsIdentifier(ast.GetElementOrPropertyAccessArgumentExpressionOrName(bin.Left)) || ast.IsStringLiteralLike(ast.GetElementOrPropertyAccessArgumentExpressionOrName(bin.Left))) {
		return jsDeclarationKindExportsProperty
	}
	// !!! module.exports property, this.property, expando.property
	return jsDeclarationKindNone
}

func isExport(node *ast.Node) bool {
	return ast.IsIdentifier(node) && node.AsIdentifier().Text == "exports"
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
