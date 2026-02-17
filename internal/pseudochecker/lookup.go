package pseudochecker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
)

func (ch *PseudoChecker) GetReturnTypeOfSignature(signatureNode *ast.Node) *PseudoType {
	switch signatureNode.Kind {
	case ast.KindGetAccessor:
		return ch.GetTypeOfAccessor(signatureNode)
	case ast.KindMethodDeclaration, ast.KindFunctionDeclaration, ast.KindConstructor,
		ast.KindMethodSignature, ast.KindCallSignature, ast.KindConstructSignature,
		ast.KindSetAccessor, ast.KindIndexSignature, ast.KindFunctionType, ast.KindConstructorType,
		ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindJSDocSignature:
		return ch.createReturnFromSignature(signatureNode)
	default:
		debug.FailBadSyntaxKind(signatureNode, "Node needs to be an inferrable node")
		return nil
	}
}

func (ch *PseudoChecker) GetTypeOfAccessor(accessor *ast.Node) *PseudoType {
	annotated := ch.typeFromAccessor(accessor)
	if annotated.Kind == PseudoTypeKindNoResult {
		return ch.inferAccessorType(accessor)
	}
	return annotated
}

func (ch *PseudoChecker) GetTypeOfExpression(node *ast.Node) *PseudoType {
	return ch.typeFromExpression(node)
}

func (ch *PseudoChecker) GetTypeOfDeclaration(node *ast.Node) *PseudoType {
	switch node.Kind {
	case ast.KindParameter:
		return ch.typeFromParameter(node.AsParameterDeclaration())
	case ast.KindVariableDeclaration:
		return ch.typeFromVariable(node.AsVariableDeclaration())
	case ast.KindPropertySignature, ast.KindPropertyDeclaration, ast.KindJSDocPropertyTag:
		return ch.typeFromProperty(node)
	case ast.KindBindingElement:
		return NewPseudoTypeNoResult(node)
	case ast.KindExportAssignment, ast.KindJSExportAssignment:
		return ch.typeFromExpression(node.AsExportAssignment().Expression)
	case ast.KindPropertyAccessExpression, ast.KindElementAccessExpression, ast.KindBinaryExpression:
		return ch.typeFromExpandoProperty(node)
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment:
		return ch.typeFromPropertyAssignment(node)
	case ast.KindCommonJSExport:
		t := node.AsCommonJSExport().Type
		if t != nil {
			return NewPseudoTypeDirect(t)
		}
		return ch.typeFromExpression(node.AsCommonJSExport().Initializer)
	case ast.KindCallExpression:
		switch ast.GetAssignmentDeclarationKind(node) {
		// TODO: How much of the checker's getTypeFromPropertyDescriptor is worth trying to emulate over ASTs?
		case ast.JSDeclarationKindObjectDefinePropertyValue:
			{
				// !!!
			}
		case ast.JSDeclarationKindObjectDefinePropertyExports:
			{
				// !!!
			}
		}
		return NewPseudoTypeNoResult(node)
	default:
		debug.FailBadSyntaxKind(node, "node needs to be an inferrable node")
		return nil
	}
}

func (ch *PseudoChecker) typeFromPropertyAssignment(node *ast.Node) *PseudoType {
	annotation := node.Type()
	if annotation != nil {
		return NewPseudoTypeDirect(annotation)
	}
	if node.Kind == ast.KindPropertyAssignment {
		init := node.Initializer()
		if init != nil {
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PseudoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PseudoTypeKindInferred
		}
	}
	return NewPseudoTypeNoResult(node)
}

// This is _not_ redundant with the reparser; see how expandoFunctionSymbolProperty.ts and similar behaves
func (ch *PseudoChecker) typeFromExpandoProperty(node *ast.Node) *PseudoType {
	declaredType := node.Type()
	if declaredType != nil {
		return NewPseudoTypeDirect(declaredType)
	}
	// While `node` is an expression, as an expando, it should also always be a
	// declaration with a `.Symbol()` which requires declaration fallback handling
	return NewPseudoTypeNoResult(node)
}

func (ch *PseudoChecker) typeFromProperty(node *ast.Node) *PseudoType {
	t := node.Type()
	if t != nil {
		return NewPseudoTypeDirect(t)
	}
	if ast.IsPropertyDeclaration(node) {
		init := node.Initializer()
		if init != nil && !isContextuallyTyped(node) {
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PseudoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PseudoTypeKindInferred
		}
	}
	return NewPseudoTypeNoResult(node)
}

func (ch *PseudoChecker) typeFromVariable(declaration *ast.VariableDeclaration) *PseudoType {
	t := declaration.Type
	if t != nil {
		return NewPseudoTypeDirect(t)
	}
	init := declaration.Initializer
	if init != nil && (len(declaration.Symbol.Declarations) == 1 || core.CountWhere(declaration.Symbol.Declarations, ast.IsVariableDeclaration) == 1) {
		if !isContextuallyTyped(declaration.AsNode()) { // TODO: also should bail on expando declarations; reuse syntactic expando check used in declaration emit
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PseudoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PseudoTypeKindInferred
		}
	}
	return NewPseudoTypeNoResult(declaration.AsNode())
}

func (ch *PseudoChecker) typeFromAccessor(accessor *ast.Node) *PseudoType {
	accessorDeclarations := ast.GetAllAccessorDeclarationsForDeclaration(accessor, accessor.DeclarationData().Symbol.Declarations)
	accessorType := ch.getTypeAnnotationFromAllAccessorDeclarations(accessor, accessorDeclarations)
	if accessorType != nil && !ast.IsTypePredicateNode(accessorType) {
		return NewPseudoTypeDirect(accessorType)
	}
	if accessorDeclarations.GetAccessor != nil {
		return ch.createReturnFromSignature(accessorDeclarations.GetAccessor.AsNode())
	}
	return NewPseudoTypeNoResult(accessor)
}

func (ch *PseudoChecker) inferAccessorType(node *ast.Node) *PseudoType {
	if node.Kind == ast.KindGetAccessor {
		return ch.createReturnFromSignature(node)
	}
	return NewPseudoTypeNoResult(node)
}

func (ch *PseudoChecker) getTypeAnnotationFromAllAccessorDeclarations(node *ast.Node, accessors ast.AllAccessorDeclarations) *ast.Node {
	accessorType := ch.getTypeAnnotationFromAccessor(node)
	if accessorType == nil && node != accessors.FirstAccessor {
		accessorType = ch.getTypeAnnotationFromAccessor(accessors.FirstAccessor)
	}
	if accessorType == nil && accessors.SecondAccessor != nil && node != accessors.SecondAccessor {
		accessorType = ch.getTypeAnnotationFromAccessor(accessors.SecondAccessor)
	}
	return accessorType
}

func (ch *PseudoChecker) getTypeAnnotationFromAccessor(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}
	// !!! TODO: support ripping return type off of .FullSignature
	if node.Kind == ast.KindGetAccessor {
		return node.AsGetAccessorDeclaration().Type
	}
	set := node.AsSetAccessorDeclaration()
	if set.Parameters == nil || len(set.Parameters.Nodes) < 1 {
		return nil
	}
	p := set.Parameters.Nodes[0]
	if !ast.IsParameter(p) {
		return nil
	}
	return p.AsParameterDeclaration().Type
}

func isValueSignatureDeclaration(node *ast.Node) bool {
	return ast.IsFunctionExpression(node) || ast.IsArrowFunction(node) || ast.IsMethodDeclaration(node) || ast.IsAccessor(node) || ast.IsFunctionDeclaration(node) || ast.IsConstructorDeclaration(node)
}

// does not return `nil`, returns a `NoResult` pseudotype instead
func (ch *PseudoChecker) createReturnFromSignature(fn *ast.Node) *PseudoType {
	if ast.IsFunctionLike(fn) {
		d := fn.FunctionLikeData()
		// !!! TODO: support ripping return type off of .FullSignature
		r := d.Type
		if r != nil {
			return NewPseudoTypeDirect(r)
		}
	}
	if isValueSignatureDeclaration(fn) {
		return ch.typeFromSingleReturnExpression(fn)
	}
	return NewPseudoTypeNoResult(fn)
}

func (ch *PseudoChecker) typeFromSingleReturnExpression(fn *ast.Node) *PseudoType {
	var candidateExpr *ast.Node
	if fn != nil && !ast.NodeIsMissing(fn.Body()) {
		flags := ast.GetFunctionFlags(fn)
		if flags&ast.FunctionFlagsAsyncGenerator != 0 {
			return NewPseudoTypeNoResult(fn)
		}

		body := fn.Body()
		if ast.IsBlock(body) {
			ast.ForEachReturnStatement(body, func(stmt *ast.Node) bool {
				if stmt.Parent != body { // Why bail on nested return statements?
					candidateExpr = nil
					return true
				}
				if candidateExpr == nil {
					candidateExpr = stmt.AsReturnStatement().Expression
				} else {
					candidateExpr = nil
					return true
				}
				return false
			})
		} else {
			candidateExpr = body
		}
	}
	if candidateExpr != nil {
		if isContextuallyTyped(candidateExpr) {
			var t *ast.Node
			if candidateExpr.Kind == ast.KindTypeAssertionExpression {
				t = candidateExpr.AsTypeAssertion().Type
			} else if candidateExpr.Kind == ast.KindAsExpression {
				t = candidateExpr.AsAsExpression().Type
			}
			if t != nil && !ast.IsConstTypeReference(t) {
				return NewPseudoTypeDirect(t)
			}
		} else {
			return ch.typeFromExpression(candidateExpr)
		}
	}
	return NewPseudoTypeNoResult(fn)
}

// This is basically `checkExpression` for pseudotypes
func (ch *PseudoChecker) typeFromExpression(node *ast.Node) *PseudoType {
	switch node.Kind {
	case ast.KindOmittedExpression:
		return PseudoTypeUndefined
	case ast.KindParenthesizedExpression:
		// assertions transformed on reparse, just unwrap
		return ch.typeFromExpression(node.AsParenthesizedExpression().Expression)
	case ast.KindIdentifier:
		// !!! TODO: in strada, this uses symbol information to ensure `node` refers to the global `undefined` symbol instead
		// we should probably import `resolveName` and use it here to check for the same; but we have to setup some barebones pseudoglobals for that to work!
		if node.AsIdentifier().Text == "undefined" {
			return PseudoTypeUndefined
		}
	case ast.KindNullKeyword:
		return PseudoTypeNull
	case ast.KindArrowFunction, ast.KindFunctionExpression:
		return ch.typeFromFunctionLikeExpression(node)
	case ast.KindTypeAssertionExpression:
		return ch.typeFromTypeAssertion(node.AsTypeAssertion().Expression, node.AsTypeAssertion().Type)
	case ast.KindAsExpression:
		return ch.typeFromTypeAssertion(node.AsAsExpression().Expression, node.AsAsExpression().Type)
	case ast.KindPrefixUnaryExpression:
		if ast.IsPrimitiveLiteralValue(node, true) {
			return ch.typeFromPrimitiveLiteralPrefix(node.AsPrefixUnaryExpression())
		}
	case ast.KindArrayLiteralExpression:
		return ch.typeFromArrayLiteral(node.AsArrayLiteralExpression())
	case ast.KindObjectLiteralExpression:
		return ch.typeFromObjectLiteral(node.AsObjectLiteralExpression())
	case ast.KindClassExpression:
		return NewPseudoTypeInferred(node) // No possible annotation/directly mappable syntax
	case ast.KindTemplateExpression:
		// templateLitWithHoles as const, not supported
		return NewPseudoTypeMaybeConstLocation(node, NewPseudoTypeInferred(node), PseudoTypeString)
	case ast.KindNumericLiteral:
		return NewPseudoTypeMaybeConstLocation(node, NewPseudoTypeNumericLiteral(node), PseudoTypeNumber)
	case ast.KindNoSubstitutionTemplateLiteral:
		return NewPseudoTypeMaybeConstLocation(node, NewPseudoTypeStringLiteral(node), PseudoTypeString)
	case ast.KindStringLiteral:
		return NewPseudoTypeMaybeConstLocation(node, NewPseudoTypeStringLiteral(node), PseudoTypeString)
	case ast.KindBigIntLiteral:
		return NewPseudoTypeMaybeConstLocation(node, NewPseudoTypeBigIntLiteral(node), PseudoTypeBigInt)
	case ast.KindTrueKeyword:
		return NewPseudoTypeMaybeConstLocation(node, PseudoTypeTrue, PseudoTypeBoolean)
	case ast.KindFalseKeyword:
		return NewPseudoTypeMaybeConstLocation(node, PseudoTypeFalse, PseudoTypeBoolean)
	}
	return NewPseudoTypeInferred(node)
}

func (ch *PseudoChecker) typeFromObjectLiteral(node *ast.ObjectLiteralExpression) *PseudoType {
	if !ch.canGetTypeFromObjectLiteral(node) {
		return NewPseudoTypeInferred(node.AsNode())
	}
	// we are in a const context producing an object literal type, there are no shorthand or spread assignments
	if node.Properties == nil || len(node.Properties.Nodes) == 0 {
		return NewPseudoTypeObjectLiteral(nil)
	}
	results := make([]*PseudoObjectElement, 0, len(node.Properties.Nodes))
	for _, e := range node.Properties.Nodes {
		switch e.Kind {
		case ast.KindMethodDeclaration:
			optional := e.AsMethodDeclaration().PostfixToken != nil && e.AsMethodDeclaration().PostfixToken.Kind == ast.KindQuestionToken
			if e.FunctionLikeData().FullSignature != nil {
				results = append(results, NewPseudoPropertyAssignment(
					false,
					e.Name(),
					optional,
					NewPseudoTypeDirect(e.FunctionLikeData().FullSignature),
				))
			} else {
				results = append(results, NewPseudoObjectMethod(
					e.Name(),
					optional,
					ch.cloneParameters(e.ParameterList()),
					ch.createReturnFromSignature(e),
				))
			}
		case ast.KindPropertyAssignment:
			results = append(results, NewPseudoPropertyAssignment(
				false,
				e.Name(),
				e.AsPropertyAssignment().PostfixToken != nil && e.AsPropertyAssignment().PostfixToken.Kind == ast.KindQuestionToken,
				ch.typeFromExpression(e.Initializer()),
			))
		case ast.KindSetAccessor, ast.KindGetAccessor:
			member := ch.getAccessorMember(e, e.Name())
			if member != nil {
				results = append(results, member)
			}
		}
	}
	return NewPseudoTypeObjectLiteral(results)
}

// roughly analogous to typeFromObjectLiteralAccessor in strada
func (ch *PseudoChecker) getAccessorMember(accessor *ast.Node, name *ast.Node) *PseudoObjectElement {
	allAccessors := ast.GetAllAccessorDeclarationsForDeclaration(accessor, accessor.Symbol().Declarations) // TODO: node preservation for late-bound accessor pairs?

	// TODO: handle pseudo-annotations from get accessor return positions?
	if allAccessors.GetAccessor != nil && allAccessors.GetAccessor.Type != nil &&
		allAccessors.SetAccessor != nil && len(allAccessors.SetAccessor.Parameters.Nodes) > 0 && allAccessors.SetAccessor.Parameters.Nodes[0].AsParameterDeclaration().Type != nil {
		// We have possible types for both accessors, we can't know if they are the same type so we keep both accessors

		if ast.IsGetAccessorDeclaration(accessor) {
			return NewPseudoGetAccessor(
				name,
				false,
				ch.typeFromAccessor(accessor),
			)
		} else {
			return NewPseudoSetAccessor(
				name,
				false,
				ch.cloneParameters(accessor.AsSetAccessorDeclaration().Parameters)[0],
			)
		}
	}

	if accessor == allAccessors.FirstAccessor {
		// only one annotated accessor; output a property - `readonly` for a single `get` accessor

		accessorType := ch.typeFromAccessor(accessor)
		readonly := ast.IsGetAccessorDeclaration(accessor) && allAccessors.SecondAccessor == nil
		return NewPseudoPropertyAssignment(
			readonly,
			name,
			false,
			accessorType,
		)
	}
	return nil
}

func (ch *PseudoChecker) canGetTypeFromObjectLiteral(node *ast.ObjectLiteralExpression) bool {
	if node.Properties == nil || len(node.Properties.Nodes) == 0 {
		return true // empty object
	}
	// !!! TODO: strada reports errors on multiple non-inferrable props
	// via calling reportInferenceFallback multiple times here before returning.
	// Does that logic need to be included in this checker? Or can it
	// be kept to the `PseudoType` -> `Node` mapping logic, so this
	// checker can avoid needing any error reporting logic?
	for _, e := range node.Properties.Nodes {
		if e.Flags&ast.NodeFlagsThisNodeHasError != 0 {
			return false
		}
		if e.Kind == ast.KindShorthandPropertyAssignment || e.Kind == ast.KindSpreadAssignment {
			return false
		}
		if e.Name().Flags&ast.NodeFlagsThisNodeHasError != 0 {
			return false
		}
		if e.Name().Kind == ast.KindPrivateIdentifier {
			return false
		}
		if e.Name().Kind == ast.KindComputedPropertyName {
			expression := e.Name().Expression()
			if !ast.IsPrimitiveLiteralValue(expression, false) {
				return false
			}
		}
	}
	return true
}

func (ch *PseudoChecker) typeFromArrayLiteral(node *ast.ArrayLiteralExpression) *PseudoType {
	if !ch.canGetTypeFromArrayLiteral(node) {
		return NewPseudoTypeInferred(node.AsNode())
	}
	// we are in a const context producing a tuple type, there are no spread elements
	results := make([]*PseudoType, 0, len(node.Elements.Nodes))
	for _, e := range node.Elements.Nodes {
		results = append(results, ch.typeFromExpression(e))
	}
	return NewPseudoTypeTuple(results)
}

func (ch *PseudoChecker) canGetTypeFromArrayLiteral(node *ast.ArrayLiteralExpression) bool {
	if !ch.isInConstContext(node.AsNode()) {
		return false
	}
	for _, e := range node.Elements.Nodes {
		if e.Kind == ast.KindSpreadElement {
			return false
		}
	}
	return true
}

// Traverses up the parent chain to determine if the node is within a const context without needing any
// persistent traversal scope tracking (which could be unreliable in the presence of `typeof` queries anyway!)
func (ch *PseudoChecker) isInConstContext(node *ast.Node) bool {
	// An expression is in a const context if an ancestor is a const type maybeAssertion expression
	maybeAssertion := ast.FindAncestor(
		node,
		func(n *ast.Node) bool {
			// stop traversing up at assertions, new scopes, and anything not an expression - they're contextual barriers
			return ast.IsAssertionExpression(n) || ast.IsFunctionLike(n) || !ast.IsExpressionNode(n)
		},
	)
	return ast.IsConstAssertion(maybeAssertion)
}

func (ch *PseudoChecker) typeFromPrimitiveLiteralPrefix(node *ast.PrefixUnaryExpression) *PseudoType {
	inner := node.Operand
	if inner.Kind == ast.KindBigIntLiteral {
		return NewPseudoTypeMaybeConstLocation(node.AsNode(), NewPseudoTypeBigIntLiteral(node.AsNode()), PseudoTypeBigInt)
	}
	if inner.Kind == ast.KindNumericLiteral {
		return NewPseudoTypeMaybeConstLocation(node.AsNode(), NewPseudoTypeNumericLiteral(node.AsNode()), PseudoTypeNumber)
	}
	debug.FailBadSyntaxKind(inner)
	return nil
}

func (ch *PseudoChecker) typeFromTypeAssertion(expression *ast.Node, typeNode *ast.Node) *PseudoType {
	if ast.IsConstTypeReference(typeNode) {
		return ch.typeFromExpression(expression)
	}
	return NewPseudoTypeDirect(typeNode)
}

func (ch *PseudoChecker) typeFromFunctionLikeExpression(node *ast.Node) *PseudoType {
	if node.FunctionLikeData().FullSignature != nil {
		return NewPseudoTypeDirect(node.FunctionLikeData().FullSignature)
	}
	returnType := ch.createReturnFromSignature(node)
	if returnType.Kind == PseudoTypeKindNoResult {
		// no result for the return type can just be an inferred result for the whole expression
		return NewPseudoTypeInferred(node.AsNode())
	}
	typeParameters := ch.cloneTypeParameters(node.FunctionLikeData().TypeParameters)
	parameters := ch.cloneParameters(node.FunctionLikeData().Parameters)
	return NewPseudoTypeSingleCallSignature(
		parameters,
		typeParameters,
		returnType,
	)
}

func (ch *PseudoChecker) cloneTypeParameters(nodes *ast.NodeList) []*ast.TypeParameterDeclaration {
	if nodes == nil {
		return nil
	}
	if len(nodes.Nodes) == 0 {
		return nil
	}
	result := make([]*ast.TypeParameterDeclaration, 0, len(nodes.Nodes))
	for _, e := range nodes.Nodes {
		result = append(result, e.AsTypeParameter())
	}
	return result
}

func (ch *PseudoChecker) typeFromParameter(node *ast.ParameterDeclaration) *PseudoType {
	parent := node.Parent
	if parent.Kind == ast.KindSetAccessor {
		return ch.GetTypeOfAccessor(parent)
	}
	declaredType := node.Type
	if declaredType != nil {
		return NewPseudoTypeDirect(declaredType)
	}
	if node.Initializer != nil && ast.IsIdentifier(node.Name()) && !isContextuallyTyped(node.Parent.AsNode()) {
		return ch.typeFromExpression(node.Initializer)
	}
	// TODO: In strada, the ID checker doesn't infer a parameter type from binding pattern names, but the real checker _does_!
	// This means ID won't let you write, say, `({elem}) => false` without an annotation, even though it's trivially of type
	// `(p0: {elem: any}) => boolean` and error-free under `noImplicitAny: false`!
	// That limitation is retained here.
	return NewPseudoTypeNoResult(node.AsNode())
}

func (ch *PseudoChecker) cloneParameters(nodes *ast.NodeList) []*PseudoParameter {
	if nodes == nil {
		return nil
	}
	if len(nodes.Nodes) == 0 {
		return nil
	}
	result := make([]*PseudoParameter, 0, len(nodes.Nodes))
	for _, e := range nodes.Nodes {
		result = append(result, NewPseudoParameter(
			e.AsParameterDeclaration().DotDotDotToken != nil,
			e.Name(),
			e.AsParameterDeclaration().QuestionToken != nil,
			ch.typeFromParameter(e.AsParameterDeclaration()),
		))
	}
	return result
}

func isContextuallyTyped(node *ast.Node) bool {
	return ast.FindAncestor(node.Parent, func(n *ast.Node) bool {
		// Functions calls or parent type annotations (but not the return type of a function expression) may impact the inferred type and local inference is unreliable
		if ast.IsCallExpression(n) {
			return true
		}
		if ast.IsFunctionLikeDeclaration(n) {
			return n.FunctionLikeData().Type != nil || n.FunctionLikeData().FullSignature != nil
		}
		return ast.IsJsxElement(n) || ast.IsJsxExpression(n)
	}) != nil
}
