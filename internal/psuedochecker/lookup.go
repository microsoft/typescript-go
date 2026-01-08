package psuedochecker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
)

func (ch *PsuedoChecker) GetReturnTypeOfSignature(signatureNode *ast.Node) *PsuedoType {
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

func (ch *PsuedoChecker) GetTypeOfAccessor(accessor *ast.Node) *PsuedoType {
	annotated := ch.typeFromAccessor(accessor)
	if annotated.Kind == PsuedoTypeKindNoResult {
		return ch.inferAccessorType(accessor)
	}
	return annotated
}

func (ch *PsuedoChecker) GetTypeOfExpression(node *ast.Node) *PsuedoType {
	return ch.typeFromExpression(node)
}

func (ch *PsuedoChecker) GetTypeOfDeclaration(node *ast.Node) *PsuedoType {
	switch node.Kind {
	case ast.KindParameter:
		return ch.typeFromParameter(node.AsParameterDeclaration())
	case ast.KindVariableDeclaration:
		return ch.typeFromVariable(node.AsVariableDeclaration())
	case ast.KindPropertySignature, ast.KindPropertyDeclaration, ast.KindJSDocPropertyTag:
		return ch.typeFromProperty(node)
	case ast.KindBindingElement:
		return NewPsuedoTypeNoResult(node)
	case ast.KindExportAssignment, ast.KindJSExportAssignment:
		return ch.typeFromExpression(node.AsExportAssignment().Expression)
	case ast.KindPropertyAccessExpression, ast.KindElementAccessExpression, ast.KindBinaryExpression:
		return ch.typeFromExpandoProperty(node)
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment:
		return ch.typeFromPropertyAssignment(node)
	default:
		debug.FailBadSyntaxKind(node, "node needs to be an inferrable node")
		return nil
	}
}

func (ch *PsuedoChecker) typeFromPropertyAssignment(node *ast.Node) *PsuedoType {
	annotation := node.Type()
	if annotation != nil {
		return NewPsuedoTypeDirect(annotation)
	}
	if node.Kind == ast.KindPropertyAssignment {
		init := node.Initializer()
		if init != nil {
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PsuedoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PsuedoTypeKindInferred
		}
	}
	return NewPsuedoTypeNoResult(node)
}

// This is _not_ redundant with the reparser; see how expandoFunctionSymbolProperty.ts and similar behaves
func (ch *PsuedoChecker) typeFromExpandoProperty(node *ast.Node) *PsuedoType {
	declaredType := node.Type()
	if declaredType != nil {
		return NewPsuedoTypeDirect(declaredType)
	}
	// While `node` is an expression, as an expando, it should also always be a
	// declaration with a `.Symbol()` which requires declaration fallback handling
	return NewPsuedoTypeNoResult(node)
}

func (ch *PsuedoChecker) typeFromProperty(node *ast.Node) *PsuedoType {
	t := node.Type()
	if t != nil {
		return NewPsuedoTypeDirect(t)
	}
	if ast.IsPropertyDeclaration(node) {
		init := node.Initializer()
		if init != nil && !isContextuallyTyped(node) {
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PsuedoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PsuedoTypeKindInferred
		}
	}
	return NewPsuedoTypeNoResult(node)
}

func (ch *PsuedoChecker) typeFromVariable(declaration *ast.VariableDeclaration) *PsuedoType {
	t := declaration.Type
	if t != nil {
		return NewPsuedoTypeDirect(t)
	}
	init := declaration.Initializer
	if init != nil && (len(declaration.Symbol.Declarations) == 1 || core.CountWhere(declaration.Symbol.Declarations, ast.IsVariableDeclaration) == 1) {
		if !isContextuallyTyped(declaration.AsNode()) { // TODO: also should bail on expando declarations; reuse syntactic expando check used in declaration emit
			expr := ch.typeFromExpression(init)
			if expr != nil && expr.Kind != PsuedoTypeKindInferred {
				return expr
			}
			// fallback to NoResult if PsuedoTypeKindInferred
		}
	}
	return NewPsuedoTypeNoResult(declaration.AsNode())
}

func (ch *PsuedoChecker) typeFromAccessor(accessor *ast.Node) *PsuedoType {
	accessorDeclarations := ast.GetAllAccessorDeclarationsForDeclaration(accessor, accessor.DeclarationData().Symbol.Declarations)
	accessorType := ch.getTypeAnnotationFromAllAccessorDeclarations(accessor, accessorDeclarations)
	if accessorType != nil && !ast.IsTypePredicateNode(accessorType) {
		return NewPsuedoTypeDirect(accessorType)
	}
	if accessorDeclarations.GetAccessor != nil {
		return ch.createReturnFromSignature(accessorDeclarations.GetAccessor.AsNode())
	}
	return NewPsuedoTypeNoResult(accessor)
}

func (ch *PsuedoChecker) inferAccessorType(node *ast.Node) *PsuedoType {
	if node.Kind == ast.KindGetAccessor {
		return ch.createReturnFromSignature(node)
	}
	return NewPsuedoTypeNoResult(node)
}

func (ch *PsuedoChecker) getTypeAnnotationFromAllAccessorDeclarations(node *ast.Node, accessors ast.AllAccessorDeclarations) *ast.Node {
	accessorType := ch.getTypeAnnotationFromAccessor(node)
	if accessorType == nil && node != accessors.FirstAccessor {
		accessorType = ch.getTypeAnnotationFromAccessor(accessors.FirstAccessor)
	}
	if accessorType == nil && accessors.SecondAccessor != nil && node != accessors.SecondAccessor {
		accessorType = ch.getTypeAnnotationFromAccessor(accessors.SecondAccessor)
	}
	return accessorType
}

func (ch *PsuedoChecker) getTypeAnnotationFromAccessor(node *ast.Node) *ast.Node {
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

// does not return `nil`, returns a `NoResult` psuedotype instead
func (ch *PsuedoChecker) createReturnFromSignature(fn *ast.Node) *PsuedoType {
	if ast.IsFunctionLike(fn) {
		d := fn.FunctionLikeData()
		// !!! TODO: support ripping return type off of .FullSignature
		r := d.Type
		if r != nil {
			return NewPsuedoTypeDirect(r)
		}
	}
	if isValueSignatureDeclaration(fn) {
		return ch.typeFromSingleReturnExpression(fn)
	}
	return NewPsuedoTypeNoResult(fn)
}

func (ch *PsuedoChecker) typeFromSingleReturnExpression(fn *ast.Node) *PsuedoType {
	var candidateExpr *ast.Node
	if fn != nil && !ast.NodeIsMissing(fn.Body()) {
		flags := ast.GetFunctionFlags(fn)
		if flags&ast.FunctionFlagsAsyncGenerator != 0 {
			return NewPsuedoTypeNoResult(fn)
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
				return NewPsuedoTypeDirect(t)
			}
		} else {
			return ch.typeFromExpression(candidateExpr)
		}
	}
	return NewPsuedoTypeNoResult(fn)
}

// This is basically `checkExpression` for psuedotypes
func (ch *PsuedoChecker) typeFromExpression(node *ast.Node) *PsuedoType {
	switch node.Kind {
	case ast.KindOmittedExpression:
		return PsuedoTypeUndefined
	case ast.KindParenthesizedExpression:
		// assertions transformed on reparse, just unwrap
		return ch.typeFromExpression(node.AsParenthesizedExpression().Expression)
	case ast.KindIdentifier:
		// !!! TODO: in strada, this uses symbol information to ensure `node` refers to the global `undefined` symbol instead
		// we should probably import `resolveName` and use it here to check for the same; but we have to setup some barebones psuedoglobals for that to work!
		if node.AsIdentifier().Text == "undefined" {
			return PsuedoTypeUndefined
		}
	case ast.KindNullKeyword:
		return PsuedoTypeNull
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
		return NewPsuedoTypeInferred(node) // No possible annotation/directly mappable syntax
	case ast.KindTemplateExpression:
		// templateLitWithHoles as const, not supported
		return NewPsuedoTypeMaybeConstLocation(node, NewPsuedoTypeInferred(node), PsuedoTypeString)
	case ast.KindNumericLiteral:
		return NewPsuedoTypeMaybeConstLocation(node, NewPsuedoTypeNumericLiteral(node), PsuedoTypeNumber)
	case ast.KindNoSubstitutionTemplateLiteral:
		return NewPsuedoTypeMaybeConstLocation(node, NewPsuedoTypeStringLiteral(node), PsuedoTypeString)
	case ast.KindStringLiteral:
		return NewPsuedoTypeMaybeConstLocation(node, NewPsuedoTypeStringLiteral(node), PsuedoTypeString)
	case ast.KindBigIntLiteral:
		return NewPsuedoTypeMaybeConstLocation(node, NewPsuedoTypeBigIntLiteral(node), PsuedoTypeBigInt)
	case ast.KindTrueKeyword:
		return NewPsuedoTypeMaybeConstLocation(node, PsuedoTypeTrue, PsuedoTypeBoolean)
	case ast.KindFalseKeyword:
		return NewPsuedoTypeMaybeConstLocation(node, PsuedoTypeFalse, PsuedoTypeBoolean)
	}
	return NewPsuedoTypeInferred(node)
}

func (ch *PsuedoChecker) typeFromObjectLiteral(node *ast.ObjectLiteralExpression) *PsuedoType {
	if !ch.canGetTypeFromObjectLiteral(node) {
		return NewPsuedoTypeInferred(node.AsNode())
	}
	// we are in a const context producing an object literal type, there are no shorthand or spread assignments
	if node.Properties == nil || len(node.Properties.Nodes) == 0 {
		return NewPsuedoTypeObjectLiteral(nil)
	}
	results := make([]*PsuedoObjectElement, 0, len(node.Properties.Nodes))
	for _, e := range node.Properties.Nodes {
		switch e.Kind {
		case ast.KindMethodDeclaration:
			optional := e.AsMethodDeclaration().PostfixToken != nil && e.AsMethodDeclaration().PostfixToken.Kind == ast.KindQuestionToken
			if e.FunctionLikeData().FullSignature != nil {
				results = append(results, NewPsuedoPropertyAssignment(
					false,
					e.Name(),
					optional,
					NewPsuedoTypeDirect(e.FunctionLikeData().FullSignature),
				))
			} else {
				results = append(results, NewPsuedoObjectMethod(
					e.Name(),
					optional,
					ch.cloneParameters(e.ParameterList()),
					ch.createReturnFromSignature(e),
				))
			}
		case ast.KindPropertyAssignment:
			results = append(results, NewPsuedoPropertyAssignment(
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
	return NewPsuedoTypeObjectLiteral(results)
}

// rougly analogous to typeFromObjectLiteralAccessor in strada
func (ch *PsuedoChecker) getAccessorMember(accessor *ast.Node, name *ast.Node) *PsuedoObjectElement {
	allAccessors := ast.GetAllAccessorDeclarationsForDeclaration(accessor, accessor.Symbol().Declarations) // TODO: node preservation for late-bound accessor pairs?

	// TODO: handle psuedo-annotations from get accessor return positions?
	if allAccessors.GetAccessor != nil && allAccessors.GetAccessor.Type != nil &&
		allAccessors.SetAccessor != nil && len(allAccessors.SetAccessor.Parameters.Nodes) > 0 && allAccessors.SetAccessor.Parameters.Nodes[0].AsParameterDeclaration().Type != nil {
		// We have possible types for both accessors, we can't know if they are the same type so we keep both accessors

		if ast.IsGetAccessorDeclaration(accessor) {
			return NewPsuedoGetAccessor(
				name,
				false,
				ch.typeFromAccessor(accessor),
			)
		} else {
			return NewPsuedoSetAccessor(
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
		return NewPsuedoPropertyAssignment(
			readonly,
			name,
			false,
			accessorType,
		)
	}
	return nil
}

func (ch *PsuedoChecker) canGetTypeFromObjectLiteral(node *ast.ObjectLiteralExpression) bool {
	if node.Properties == nil || len(node.Properties.Nodes) == 0 {
		return true // empty object
	}
	// !!! TODO: strada reports errors on multiple non-inferrable props
	// via calling reportInferenceFallback multiple times here before returning.
	// Does that logic need to be included in this checker? Or can it
	// be kept to the `PsuedoType` -> `Node` mapping logic, so this
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

func (ch *PsuedoChecker) typeFromArrayLiteral(node *ast.ArrayLiteralExpression) *PsuedoType {
	if !ch.canGetTypeFromArrayLiteral(node) {
		return NewPsuedoTypeInferred(node.AsNode())
	}
	// we are in a const context producing a tuple type, there are no spread elements
	results := make([]*PsuedoType, 0, len(node.Elements.Nodes))
	for _, e := range node.Elements.Nodes {
		results = append(results, ch.typeFromExpression(e))
	}
	return NewPsuedoTypeTuple(results)
}

func (ch *PsuedoChecker) canGetTypeFromArrayLiteral(node *ast.ArrayLiteralExpression) bool {
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
func (ch *PsuedoChecker) isInConstContext(node *ast.Node) bool {
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

func (ch *PsuedoChecker) typeFromPrimitiveLiteralPrefix(node *ast.PrefixUnaryExpression) *PsuedoType {
	inner := node.Operand
	if inner.Kind == ast.KindBigIntLiteral {
		return NewPsuedoTypeMaybeConstLocation(node.AsNode(), NewPsuedoTypeBigIntLiteral(node.AsNode()), PsuedoTypeBigInt)
	}
	if inner.Kind == ast.KindNumericLiteral {
		return NewPsuedoTypeMaybeConstLocation(node.AsNode(), NewPsuedoTypeNumericLiteral(node.AsNode()), PsuedoTypeNumber)
	}
	debug.FailBadSyntaxKind(inner)
	return nil
}

func (ch *PsuedoChecker) typeFromTypeAssertion(expression *ast.Node, typeNode *ast.Node) *PsuedoType {
	if ast.IsConstTypeReference(typeNode) {
		return ch.typeFromExpression(expression)
	}
	return NewPsuedoTypeDirect(typeNode)
}

func (ch *PsuedoChecker) typeFromFunctionLikeExpression(node *ast.Node) *PsuedoType {
	if node.FunctionLikeData().FullSignature != nil {
		return NewPsuedoTypeDirect(node.FunctionLikeData().FullSignature)
	}
	returnType := ch.createReturnFromSignature(node)
	if returnType.Kind == PsuedoTypeKindNoResult {
		// no result for the return type can just be an inferred result for the whole expression
		return NewPsuedoTypeInferred(node.AsNode())
	}
	typeParameters := ch.cloneTypeParameters(node.FunctionLikeData().TypeParameters)
	parameters := ch.cloneParameters(node.FunctionLikeData().Parameters)
	return NewPsuedoTypeSingleCallSignature(
		parameters,
		typeParameters,
		returnType,
	)
}

func (ch *PsuedoChecker) cloneTypeParameters(nodes *ast.NodeList) []*ast.TypeParameterDeclaration {
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

func (ch *PsuedoChecker) typeFromParameter(node *ast.ParameterDeclaration) *PsuedoType {
	parent := node.Parent
	if parent.Kind == ast.KindSetAccessor {
		return ch.GetTypeOfAccessor(parent)
	}
	declaredType := node.Type
	if declaredType != nil {
		return NewPsuedoTypeDirect(declaredType)
	}
	if node.Initializer != nil && ast.IsIdentifier(node.Name()) && !isContextuallyTyped(node.AsNode()) {
		return ch.typeFromExpression(node.Initializer)
	}
	// TODO: In strada, the ID checker doesn't infer a parameter type from binding pattern names, but the real checker _does_!
	// This means ID won't let you write, say, `({elem}) => false` without an annotation, even though it's trivially of type
	// `(p0: {elem: any}) => boolean` and error-free under `noImplicitAny: false`!
	// That limitation is retained here.
	return NewPsuedoTypeNoResult(node.AsNode())
}

func (ch *PsuedoChecker) cloneParameters(nodes *ast.NodeList) []*PsuedoParameter {
	if nodes == nil {
		return nil
	}
	if len(nodes.Nodes) == 0 {
		return nil
	}
	result := make([]*PsuedoParameter, 0, len(nodes.Nodes))
	for _, e := range nodes.Nodes {
		result = append(result, NewPsuedoParameter(
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
