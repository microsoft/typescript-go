package transformers

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/printer"
)

func copyIdentifier(emitContext *printer.EmitContext, node *ast.IdentifierNode) *ast.IdentifierNode {
	var nodeCopy *ast.IdentifierNode
	if emitContext.HasAutoGenerateInfo(node) {
		nodeCopy = emitContext.NewGeneratedNameForNode(node, printer.AutoGenerateOptions{})
	} else {
		nodeCopy = emitContext.Factory.NewIdentifier(node.Text())
		nodeCopy.Flags = node.Flags
		nodeCopy.Loc = node.Loc
	}
	emitContext.SetOriginal(nodeCopy, node)
	return nodeCopy
}

func getName(emitContext *printer.EmitContext, node *ast.Declaration, allowComments bool, allowSourceMaps bool, emitFlags printer.EmitFlags, ignoreAssignedName bool) *ast.IdentifierNode {
	var nodeName *ast.IdentifierNode
	if node != nil {
		if ignoreAssignedName {
			nodeName = ast.GetNonAssignedNameOfDeclaration(node)
		} else {
			nodeName = ast.GetNameOfDeclaration(node)
		}
	}

	if nodeName != nil {
		name := copyIdentifier(emitContext, nodeName)
		if !allowComments {
			emitContext.AddEmitFlags(name, printer.EFNoComments)
		}
		if !allowSourceMaps {
			emitContext.AddEmitFlags(name, printer.EFNoSourceMap)
		}
		return name
	}

	return emitContext.NewGeneratedNameForNode(node, printer.AutoGenerateOptions{})
}

// Gets the local name of a declaration. This is primarily used for declarations that can be referred to by name in the
// declaration's immediate scope (classes, enums, namespaces). A local name will *never* be prefixed with a module or
// namespace export modifier like "exports." when emitted as an expression.
//
// The value of the allowComments parameter indicates whether comments may be emitted for the name.
// The value of the allowSourceMaps parameter indicates whether source maps may be emitted for the name.
// The value of the ignoreAssignedName parameter indicates whether the assigned name of a declaration shouldn't be considered.
func getLocalName(emitContext *printer.EmitContext, node *ast.Declaration, allowComments bool, allowSourceMaps bool, ignoreAssignedName bool) *ast.IdentifierNode {
	return getName(emitContext, node, allowComments, allowSourceMaps, printer.EFLocalName, ignoreAssignedName)
}

// Gets the name of a declaration to use during emit.
//
// The value of the allowComments parameter indicates whether comments may be emitted for the name.
// The value of the allowSourceMaps parameter indicates whether source maps may be emitted for the name.
func getDeclarationName(emitContext *printer.EmitContext, node *ast.Declaration, allowComments bool, allowSourceMaps bool) *ast.IdentifierNode {
	return getName(emitContext, node, allowComments, allowSourceMaps, printer.EFNone, false /*ignoreAssignedName*/)
}

func isIdentifierReference(name *ast.IdentifierNode, parent *ast.Node) bool {
	switch parent.Kind {
	case ast.KindBinaryExpression,
		ast.KindPrefixUnaryExpression,
		ast.KindPostfixUnaryExpression,
		ast.KindYieldExpression,
		ast.KindAsExpression,
		ast.KindSatisfiesExpression,
		ast.KindElementAccessExpression,
		ast.KindNonNullExpression,
		ast.KindSpreadElement,
		ast.KindSpreadAssignment,
		ast.KindParenthesizedExpression,
		ast.KindArrayLiteralExpression,
		ast.KindDeleteExpression,
		ast.KindTypeOfExpression,
		ast.KindVoidExpression,
		ast.KindAwaitExpression,
		ast.KindTypeAssertionExpression,
		ast.KindExpressionWithTypeArguments,
		ast.KindJsxSelfClosingElement,
		ast.KindJsxSpreadAttribute,
		ast.KindJsxExpression,
		ast.KindCommaListExpression,
		ast.KindPartiallyEmittedExpression:
		// all immediate children that can be `Identifier` would be instances of `IdentifierReference`
		return true
	case ast.KindComputedPropertyName,
		ast.KindDecorator,
		ast.KindIfStatement,
		ast.KindDoStatement,
		ast.KindWhileStatement,
		ast.KindWithStatement,
		ast.KindReturnStatement,
		ast.KindSwitchStatement,
		ast.KindCaseClause,
		ast.KindThrowStatement,
		ast.KindExpressionStatement,
		ast.KindExportAssignment,
		ast.KindPropertyAccessExpression:
		// only an `Expression()` child that can be `Identifier` would be an instance of `IdentifierReference`
		return parent.Expression() == name
	case ast.KindVariableDeclaration,
		ast.KindParameter,
		ast.KindBindingElement,
		ast.KindPropertyDeclaration,
		ast.KindPropertySignature,
		ast.KindPropertyAssignment,
		ast.KindEnumMember,
		ast.KindJsxAttribute:
		// only an `Initializer()` child that can be `Identifier` would be an instance of `IdentifierReference`
		return parent.Initializer() == name
	case ast.KindForStatement:
		return parent.AsForStatement().Initializer == name ||
			parent.AsForStatement().Condition == name ||
			parent.AsForStatement().Incrementor == name
	case ast.KindForInStatement,
		ast.KindForOfStatement:
		return parent.AsForInOrOfStatement().Initializer == name ||
			parent.AsForInOrOfStatement().Expression == name
	case ast.KindImportEqualsDeclaration:
		return parent.AsImportEqualsDeclaration().ModuleReference == name
	case ast.KindArrowFunction:
		return parent.AsArrowFunction().Body == name
	case ast.KindConditionalExpression:
		return parent.AsConditionalExpression().Condition == name ||
			parent.AsConditionalExpression().WhenTrue == name ||
			parent.AsConditionalExpression().WhenFalse == name
	case ast.KindCallExpression:
		return parent.AsCallExpression().Expression == name ||
			slices.Contains(parent.AsCallExpression().Arguments.Nodes, name)
	case ast.KindNewExpression:
		return parent.AsNewExpression().Expression == name ||
			parent.AsNewExpression().Arguments.Nodes != nil &&
				slices.Contains(parent.AsNewExpression().Arguments.Nodes, name)
	case ast.KindTaggedTemplateExpression:
		return parent.AsTaggedTemplateExpression().Tag == name
	case ast.KindImportAttribute:
		return parent.AsImportAttribute().Value == name
	case ast.KindJsxOpeningElement:
		return parent.AsJsxOpeningElement().TagName == name
	default:
		return false
	}
}

func constantValue(node *ast.Expression) any {
	node = ast.SkipOuterExpressions(node, ast.OEKAll)
	if ast.IsStringLiteralLike(node) {
		return node.Text()
	}
	if ast.IsPrefixUnaryExpression(node) {
		prefixUnary := node.AsPrefixUnaryExpression()
		if value, ok := constantValue(prefixUnary.Operand).(jsnum.Number); ok {
			switch prefixUnary.Operator {
			case ast.KindPlusToken:
				return value
			case ast.KindMinusToken:
				return -value
			case ast.KindTildeToken:
				return value.BitwiseNOT()
			}
		}
	}
	if ast.IsNumericLiteral(node) {
		return jsnum.FromString(node.Text())
	}
	return nil
}

func constantExpression(value any, factory *ast.NodeFactory) *ast.Expression {
	switch value := value.(type) {
	case string:
		return factory.NewStringLiteral(value)
	case jsnum.Number:
		if value.IsInf() || value.IsNaN() {
			return nil
		}
		if value < 0 {
			return factory.NewPrefixUnaryExpression(ast.KindMinusToken, constantExpression(-value, factory))
		}
		return factory.NewNumericLiteral(value.String())
	}
	return nil
}
