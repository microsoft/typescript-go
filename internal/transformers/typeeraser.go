package transformers

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

type TypeEraserTransformer struct {
	Transformer
}

func NewTypeEraserTransformer(emitContext *printer.EmitContext) *Transformer {
	tx := &TypeEraserTransformer{}
	return tx.newTransformer(tx.visit, emitContext)
}

func (tx *TypeEraserTransformer) visit(node *ast.Node) *ast.Node {
	// !!! TransformFlags were traditionally used here to skip over subtrees that contain no TypeScript syntax
	if ast.IsStatement(node) && ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) {
		// !!! Use NotEmittedStatement to preserve comments
		return nil
	}

	switch node.Kind {
	case
		// TypeScript accessibility and readonly modifiers are elided
		ast.KindPublicKeyword,
		ast.KindPrivateKeyword,
		ast.KindProtectedKeyword,
		ast.KindAbstractKeyword,
		ast.KindOverrideKeyword,
		ast.KindConstKeyword,
		ast.KindDeclareKeyword,
		ast.KindReadonlyKeyword,
		ast.KindInKeyword,
		ast.KindOutKeyword,
		// TypeScript type nodes are elided.
		ast.KindArrayType,
		ast.KindTupleType,
		ast.KindOptionalType,
		ast.KindRestType,
		ast.KindTypeLiteral,
		ast.KindTypePredicate,
		ast.KindTypeParameter,
		ast.KindAnyKeyword,
		ast.KindUnknownKeyword,
		ast.KindBooleanKeyword,
		ast.KindStringKeyword,
		ast.KindNumberKeyword,
		ast.KindNeverKeyword,
		ast.KindVoidKeyword,
		ast.KindSymbolKeyword,
		ast.KindConstructorType,
		ast.KindFunctionType,
		ast.KindTypeQuery,
		ast.KindTypeReference,
		ast.KindUnionType,
		ast.KindIntersectionType,
		ast.KindConditionalType,
		ast.KindParenthesizedType,
		ast.KindThisType,
		ast.KindTypeOperator,
		ast.KindIndexedAccessType,
		ast.KindMappedType,
		ast.KindLiteralType,
		// TypeScript index signatures are elided.
		ast.KindIndexSignature:
		return nil

	case ast.KindTypeAliasDeclaration,
		ast.KindInterfaceDeclaration,
		ast.KindNamespaceExportDeclaration:
		// TypeScript type-only declarations are elided.
		// !!! Use NotEmittedStatement to preserve comments
		return nil

	case ast.KindExpressionWithTypeArguments:
		n := node.AsExpressionWithTypeArguments()
		return tx.Factory.UpdateExpressionWithTypeArguments(n, tx.VisitNode(n.Expression), nil)

	case ast.KindPropertyDeclaration:
		if ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) {
			// TypeScript `declare` fields are elided
			return nil
		}
		n := node.AsPropertyDeclaration()
		return tx.Factory.UpdatePropertyDeclaration(n, tx.VisitModifiers(n.Modifiers()), tx.VisitNode(n.Name()), nil, nil, tx.VisitNode(n.Initializer))

	case ast.KindConstructor:
		n := node.AsConstructorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		return tx.Factory.UpdateConstructorDeclaration(n, nil, nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindMethodDeclaration:
		n := node.AsMethodDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		return tx.Factory.UpdateMethodDeclaration(n, tx.VisitModifiers(n.Modifiers()), n.AsteriskToken, tx.VisitNode(n.Name()), nil, nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindGetAccessor:
		n := node.AsGetAccessorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		return tx.Factory.UpdateGetAccessorDeclaration(n, tx.VisitModifiers(n.Modifiers()), tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindSetAccessor:
		n := node.AsSetAccessorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		return tx.Factory.UpdateSetAccessorDeclaration(n, tx.VisitModifiers(n.Modifiers()), tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindVariableDeclaration:
		n := node.AsVariableDeclaration()
		return tx.Factory.UpdateVariableDeclaration(n, tx.VisitNode(n.Name()), nil, nil, tx.VisitNode(n.Initializer))

	case ast.KindHeritageClause:
		n := node.AsHeritageClause()
		if n.Token == ast.KindImplementsKeyword {
			// TypeScript `implements` clauses are elided
			return nil
		}
		return tx.Factory.UpdateHeritageClause(n, tx.VisitNodes(n.Types))

	case ast.KindClassDeclaration:
		n := node.AsClassDeclaration()
		return tx.Factory.UpdateClassDeclaration(n, tx.VisitModifiers(n.Modifiers()), tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.HeritageClauses), tx.VisitNodes(n.Members))

	case ast.KindClassExpression:
		n := node.AsClassExpression()
		return tx.Factory.UpdateClassExpression(n, tx.VisitModifiers(n.Modifiers()), tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.HeritageClauses), tx.VisitNodes(n.Members))

	case ast.KindFunctionDeclaration:
		n := node.AsFunctionDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		return tx.Factory.UpdateFunctionDeclaration(n, tx.VisitModifiers(n.Modifiers()), n.AsteriskToken, tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindFunctionExpression:
		n := node.AsFunctionExpression()
		return tx.Factory.UpdateFunctionExpression(n, tx.VisitModifiers(n.Modifiers()), n.AsteriskToken, tx.VisitNode(n.Name()), nil, tx.VisitNodes(n.Parameters), nil, tx.VisitNode(n.Body))

	case ast.KindArrowFunction:
		n := node.AsArrowFunction()
		return tx.Factory.UpdateArrowFunction(n, tx.VisitModifiers(n.Modifiers()), nil, tx.VisitNodes(n.Parameters), nil, n.EqualsGreaterThanToken, tx.VisitNode(n.Body))

	case ast.KindParameter:
		if ast.IsThisParameter(node) {
			// TypeScript `this` parameters are elided
			return nil
		}
		n := node.AsParameterDeclaration()
		return tx.Factory.UpdateParameterDeclaration(n, nil, n.DotDotDotToken, tx.VisitNode(n.Name()), nil, nil, tx.VisitNode(n.Initializer))

	case ast.KindCallExpression:
		n := node.AsCallExpression()
		return tx.Factory.UpdateCallExpression(n, tx.VisitNode(n.Expression), n.QuestionDotToken, nil, tx.VisitNodes(n.Arguments))

	case ast.KindNewExpression:
		n := node.AsNewExpression()
		return tx.Factory.UpdateNewExpression(n, tx.VisitNode(n.Expression), nil, tx.VisitNodes(n.Arguments))

	case ast.KindTaggedTemplateExpression:
		n := node.AsTaggedTemplateExpression()
		return tx.Factory.UpdateTaggedTemplateExpression(n, tx.VisitNode(n.Tag), n.QuestionDotToken, nil, tx.VisitNode(n.Template))

	case ast.KindNonNullExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return tx.VisitNode(node.AsNonNullExpression().Expression)

	case ast.KindTypeAssertionExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return tx.VisitNode(node.AsTypeAssertion().Expression)

	case ast.KindAsExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return tx.VisitNode(node.AsAsExpression().Expression)

	case ast.KindSatisfiesExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return tx.VisitNode(node.AsSatisfiesExpression().Expression)

	case ast.KindJsxSelfClosingElement:
		n := node.AsJsxSelfClosingElement()
		return tx.Factory.UpdateJsxSelfClosingElement(n, tx.VisitNode(n.TagName), nil, tx.VisitNode(n.Attributes))

	case ast.KindJsxOpeningElement:
		n := node.AsJsxOpeningElement()
		return tx.Factory.UpdateJsxOpeningElement(n, tx.VisitNode(n.TagName), nil, tx.VisitNode(n.Attributes))

	default:
		return tx.VisitEachChild(node)
	}
}
