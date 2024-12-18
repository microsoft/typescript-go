package transformers

import (
	"github.com/microsoft/typescript-go/internal/ast"
)

type TypeEraserTransformer struct {
	ast.NodeVisitor
}

func NewTypeEraserTransformer() *TypeEraserTransformer {
	visitor := &TypeEraserTransformer{}
	visitor.Visit = visitor.visit
	return visitor
}

func (v *TypeEraserTransformer) visit(node *ast.Node) *ast.Node {
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
		expression := v.VisitNode(n.Expression)
		if expression != n.Expression || nil != n.TypeArguments {
			return v.UpdateNode(v.Factory.NewExpressionWithTypeArguments(expression, nil), node)
		}

	case ast.KindPropertyDeclaration:
		if ast.HasSyntacticModifier(node, ast.ModifierFlagsAmbient) {
			// TypeScript `declare` fields are elided
			return nil
		}
		n := node.AsPropertyDeclaration()
		modifiers := v.VisitModifiers(n.Modifiers())
		name := v.VisitNode(n.Name())
		initializer := v.VisitNode(n.Initializer)
		if modifiers != n.Modifiers() || name != n.Name() || nil != n.PostfixToken || nil != n.Type || initializer != n.Initializer {
			return v.UpdateNode(v.Factory.NewPropertyDeclaration(modifiers, name, nil, nil, initializer), node)
		}

	case ast.KindConstructor:
		n := node.AsConstructorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if nil != n.Modifiers() || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewConstructorDeclaration(nil, nil, parameters, nil, body), node)
		}

	case ast.KindMethodDeclaration:
		n := node.AsMethodDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		modifiers := v.VisitModifiers(n.Modifiers())
		asteriskToken := v.VisitToken(n.AsteriskToken)
		name := v.VisitNode(n.Name())
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || asteriskToken != n.AsteriskToken || name != n.Name() || nil != n.PostfixToken || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewMethodDeclaration(modifiers, asteriskToken, name, nil, nil, parameters, nil, body), node)
		}

	case ast.KindGetAccessor:
		n := node.AsGetAccessorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		modifiers := v.VisitModifiers(n.Modifiers())
		name := v.VisitNode(n.Name())
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || name != n.Name() || nil != n.PostfixToken || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewGetAccessorDeclaration(modifiers, name, nil, parameters, nil, body), node)
		}

	case ast.KindSetAccessor:
		n := node.AsSetAccessorDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		modifiers := v.VisitModifiers(n.Modifiers())
		name := v.VisitNode(n.Name())
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || name != n.Name() || nil != n.PostfixToken || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewSetAccessorDeclaration(modifiers, name, nil, parameters, nil, body), node)
		}

	case ast.KindVariableDeclaration:
		n := node.AsVariableDeclaration()
		name := v.VisitNode(n.Name())
		initializer := v.VisitNode(n.Initializer)
		if name != n.Name() || nil != n.ExclamationToken || nil != n.Type || initializer != n.Initializer {
			return v.UpdateNode(v.Factory.NewVariableDeclaration(name, nil, nil, initializer), node)
		}

	case ast.KindHeritageClause:
		n := node.AsHeritageClause()
		if n.Token == ast.KindImplementsKeyword {
			// TypeScript `implements` clauses are elided
			return nil
		}
		types := v.VisitNodes(n.Types)
		if types != n.Types {
			return v.UpdateNode(v.Factory.NewHeritageClause(n.Token, types), node)
		}

	case ast.KindClassDeclaration:
		n := node.AsClassDeclaration()
		modifiers := v.VisitModifiers(n.Modifiers())
		name := v.VisitNode(n.Name())
		heritageClauses := v.VisitNodes(n.HeritageClauses)
		if heritageClauses != nil && len(heritageClauses.Nodes) == 0 {
			heritageClauses = nil
		}
		members := v.VisitNodes(n.Members)
		if modifiers != n.Modifiers() || name != n.Name() || nil != n.TypeParameters || heritageClauses != n.HeritageClauses || members != n.Members {
			return v.UpdateNode(v.Factory.NewClassDeclaration(modifiers, name, nil, heritageClauses, members), node)
		}

	case ast.KindClassExpression:
		n := node.AsClassExpression()
		modifiers := v.VisitModifiers(n.Modifiers())
		name := v.VisitNode(n.Name())
		heritageClauses := v.VisitNodes(n.HeritageClauses)
		if heritageClauses != nil && len(heritageClauses.Nodes) == 0 {
			heritageClauses = nil
		}
		members := v.VisitNodes(n.Members)
		if modifiers != n.Modifiers() || name != n.Name() || nil != n.TypeParameters || heritageClauses != n.HeritageClauses || members != n.Members {
			return v.UpdateNode(v.Factory.NewClassExpression(modifiers, name, nil, heritageClauses, members), node)
		}

	case ast.KindFunctionDeclaration:
		n := node.AsFunctionDeclaration()
		if n.Body == nil {
			// TypeScript overloads are elided
			return nil
		}
		modifiers := v.VisitModifiers(n.Modifiers())
		asteriskToken := v.VisitToken(n.AsteriskToken)
		name := v.VisitNode(n.Name())
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || asteriskToken != n.AsteriskToken || name != n.Name() || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewFunctionDeclaration(modifiers, asteriskToken, name, nil, parameters, nil, body), node)
		}

	case ast.KindFunctionExpression:
		n := node.AsFunctionExpression()
		modifiers := v.VisitModifiers(n.Modifiers())
		asteriskToken := v.VisitToken(n.AsteriskToken)
		name := v.VisitNode(n.Name())
		parameters := v.VisitNodes(n.Parameters)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || asteriskToken != n.AsteriskToken || name != n.Name() || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || body != n.Body {
			return v.UpdateNode(v.Factory.NewFunctionExpression(modifiers, asteriskToken, name, nil, parameters, nil, body), node)
		}

	case ast.KindArrowFunction:
		n := node.AsArrowFunction()
		modifiers := v.VisitModifiers(n.Modifiers())
		parameters := v.VisitNodes(n.Parameters)
		equalsGreaterThanToken := v.VisitToken(n.EqualsGreaterThanToken)
		body := v.VisitNode(n.Body)
		if modifiers != n.Modifiers() || nil != n.TypeParameters || parameters != n.Parameters || nil != n.Type || equalsGreaterThanToken != n.EqualsGreaterThanToken || body != n.Body {
			return v.UpdateNode(v.Factory.NewArrowFunction(modifiers, nil, parameters, nil, equalsGreaterThanToken, body), node)
		}

	case ast.KindParameter:
		if ast.IsThisParameter(node) {
			// TypeScript `this` parameters are elided
			return nil
		}
		n := node.AsParameterDeclaration()
		dotDotDotToken := v.VisitToken(n.DotDotDotToken)
		name := v.VisitNode(n.Name())
		initializer := v.VisitNode(n.Initializer)
		if nil != n.Modifiers() || dotDotDotToken != n.DotDotDotToken || name != n.Name() || nil != n.QuestionToken || nil != n.Type || initializer != n.Initializer {
			return v.UpdateNode(v.Factory.NewParameterDeclaration(nil, dotDotDotToken, name, nil, nil, initializer), node)
		}

	case ast.KindCallExpression:
		n := node.AsCallExpression()
		expression := v.VisitNode(n.Expression)
		questionDotToken := v.VisitToken(n.QuestionDotToken)
		arguments := v.VisitNodes(n.Arguments)
		if expression != n.Expression || questionDotToken != n.QuestionDotToken || nil != n.TypeArguments || arguments != n.Arguments {
			return v.UpdateNode(v.Factory.NewCallExpression(expression, questionDotToken, nil, arguments, node.Flags), node)
		}

	case ast.KindNewExpression:
		n := node.AsNewExpression()
		expression := v.VisitNode(n.Expression)
		arguments := v.VisitNodes(n.Arguments)
		if expression != n.Expression || nil != n.TypeArguments || arguments != n.Arguments {
			return v.UpdateNode(v.Factory.NewNewExpression(expression, nil, arguments), node)
		}

	case ast.KindTaggedTemplateExpression:
		n := node.AsTaggedTemplateExpression()
		tag := v.VisitNode(n.Tag)
		questionDotToken := v.VisitToken(n.QuestionDotToken)
		template := v.VisitNode(n.Template)
		if tag != n.Tag || questionDotToken != n.QuestionDotToken || nil != n.TypeArguments || template != n.Template {
			return v.UpdateNode(v.Factory.NewTaggedTemplateExpression(tag, questionDotToken, nil, template, node.Flags), node)
		}

	case ast.KindNonNullExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return v.VisitNode(node.AsNonNullExpression().Expression)

	case ast.KindTypeAssertionExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return v.VisitNode(node.AsTypeAssertion().Expression)

	case ast.KindAsExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return v.VisitNode(node.AsAsExpression().Expression)

	case ast.KindSatisfiesExpression:
		// !!! Use PartiallyEmittedExpression to preserve comments
		return v.VisitNode(node.AsSatisfiesExpression().Expression)

	case ast.KindJsxSelfClosingElement:
		n := node.AsJsxSelfClosingElement()
		tagName := v.VisitNode(n.TagName)
		attributes := v.VisitNode(n.Attributes)
		if tagName != n.TagName || nil != n.TypeArguments || attributes != n.Attributes {
			return v.UpdateNode(v.Factory.NewJsxSelfClosingElement(tagName, nil, attributes), node)
		}

	case ast.KindJsxOpeningElement:
		n := node.AsJsxOpeningElement()
		tagName := v.VisitNode(n.TagName)
		attributes := v.VisitNode(n.Attributes)
		if tagName != n.TagName || nil != n.TypeArguments || attributes != n.Attributes {
			return v.UpdateNode(v.Factory.NewJsxOpeningElement(tagName, nil, attributes), node)
		}

	default:
		return v.VisitEachChild(node)
	}

	// Reuse subtree
	return node
}
