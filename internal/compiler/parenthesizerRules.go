package compiler

type ParenthesizerRules interface {
	GetParenthesizeLeftSideOfBinaryForOperator(binaryOperator SyntaxKind) func(leftSide *Expression) *Expression
	GetParenthesizeRightSideOfBinaryForOperator(binaryOperator SyntaxKind) func(rightSide *Expression) *Expression
	ParenthesizeLeftSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression) *Expression
	ParenthesizeRightSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression, rightSide *Expression) *Expression
	ParenthesizeExpressionOfComputedPropertyName(expression *Expression) *Expression
	ParenthesizeConditionOfConditionalExpression(condition *Expression) *Expression
	ParenthesizeBranchOfConditionalExpression(branch *Expression) *Expression
	ParenthesizeExpressionOfExportDefault(expression *Expression) *Expression
	ParenthesizeExpressionOfNew(expression *Expression) *Expression
	ParenthesizeLeftSideOfAccess(expression *Expression, optionalChain bool) *Expression
	ParenthesizeOperandOfPostfixUnary(operand *Expression) *Expression
	ParenthesizeOperandOfPrefixUnary(operand *Expression) *Expression
	ParenthesizeExpressionsOfCommaDelimitedList(elements []*Expression) []*Expression
	ParenthesizeExpressionForDisallowedComma(expression *Expression) *Expression
	ParenthesizeExpressionOfExpressionStatement(expression *Expression) *Expression
	ParenthesizeConciseBodyOfArrowFunction(body *BlockOrExpression) *BlockOrExpression
	ParenthesizeCheckTypeOfConditionalType(typeNode *TypeNode) *TypeNode
	ParenthesizeExtendsTypeOfConditionalType(typeNode *TypeNode) *TypeNode
	ParenthesizeOperandOfTypeOperator(typeNode *TypeNode) *TypeNode
	ParenthesizeOperandOfReadonlyTypeOperator(typeNode *TypeNode) *TypeNode
	ParenthesizeNonArrayTypeOfPostfixType(typeNode *TypeNode) *TypeNode
	ParenthesizeElementTypesOfTupleType(types []*Node) []*Node
	ParenthesizeElementTypeOfTupleType(typeNode *Node) *Node
	ParenthesizeTypeOfOptionalType(typeNode *TypeNode) *TypeNode
	ParenthesizeConstituentTypeOfUnionType(typeNode *TypeNode) *TypeNode
	ParenthesizeConstituentTypesOfUnionType(constituents []*TypeNode) []*TypeNode
	ParenthesizeConstituentTypeOfIntersectionType(typeNode *TypeNode) *TypeNode
	ParenthesizeConstituentTypesOfIntersectionType(constituents []*TypeNode) []*TypeNode
	ParenthesizeLeadingTypeArgument(typeNode *TypeNode) *TypeNode
	ParenthesizeTypeArguments(typeArguments *TypeArgumentListNode) *TypeArgumentListNode
}

type parenthesizerRules struct {
	factory                                       *NodeFactory
	parenthesizeExpressionForDisallowedComma      func(expression *Expression) *Expression
	parenthesizeElementTypeOfTupleType            func(typeNode *Node) *Node
	parenthesizeConstituentTypeOfUnionType        func(typeNode *TypeNode) *TypeNode
	parenthesizeConstituentTypeOfIntersectionType func(typeNode *TypeNode) *TypeNode
	parenthesizeOrdinalTypeArgument               func(node *TypeNode, i int) *TypeNode
	binaryLeftOperandParenthesizerCache           map[SyntaxKind]func(node *Expression) *Expression
	binaryRightOperandParenthesizerCache          map[SyntaxKind]func(node *Expression) *Expression
}

func NewParenthesizerRules(factory *NodeFactory) ParenthesizerRules {
	rules := &parenthesizerRules{}
	rules.factory = factory
	rules.parenthesizeExpressionForDisallowedComma = rules.ParenthesizeExpressionForDisallowedComma
	rules.parenthesizeElementTypeOfTupleType = rules.ParenthesizeElementTypeOfTupleType
	rules.parenthesizeConstituentTypeOfUnionType = rules.ParenthesizeConstituentTypeOfUnionType
	rules.parenthesizeConstituentTypeOfIntersectionType = rules.ParenthesizeConstituentTypeOfIntersectionType
	rules.parenthesizeOrdinalTypeArgument = rules.parenthesizeOrdinalTypeArgumentWorker
	return rules
}

func (p *parenthesizerRules) GetParenthesizeLeftSideOfBinaryForOperator(operatorKind SyntaxKind) func(leftSide *Expression) *Expression {
	if len(p.binaryLeftOperandParenthesizerCache) == 0 {
		p.binaryLeftOperandParenthesizerCache = make(map[SyntaxKind]func(node *Node) *Node)
	}

	parenthesizerRule, ok := p.binaryLeftOperandParenthesizerCache[operatorKind]
	if !ok {
		parenthesizerRule = func(node *Expression) *Expression {
			return p.ParenthesizeLeftSideOfBinary(operatorKind, node)
		}
		p.binaryLeftOperandParenthesizerCache[operatorKind] = parenthesizerRule
	}
	return parenthesizerRule
}

func (p *parenthesizerRules) GetParenthesizeRightSideOfBinaryForOperator(operatorKind SyntaxKind) func(rightSide *Expression) *Expression {
	if len(p.binaryRightOperandParenthesizerCache) == 0 {
		p.binaryRightOperandParenthesizerCache = make(map[SyntaxKind]func(node *Node) *Node)
	}

	parenthesizerRule, ok := p.binaryRightOperandParenthesizerCache[operatorKind]
	if !ok {
		parenthesizerRule = func(node *Expression) *Expression {
			return p.ParenthesizeRightSideOfBinary(operatorKind, nil /*leftSide*/, node)
		}
		p.binaryRightOperandParenthesizerCache[operatorKind] = parenthesizerRule
	}
	return parenthesizerRule
}

// Determines whether the operand to a BinaryExpression needs to be parenthesized.
//   - binaryOperator - The operator for the BinaryExpression.
//   - operand - The operand for the BinaryExpression.
//   - isLeftSideOfBinary - A value indicating whether the operand is the left side of the BinaryExpression.
func (p *parenthesizerRules) binaryOperandNeedsParentheses(binaryOperator SyntaxKind, operand *Expression, isLeftSideOfBinary bool, leftOperand *Expression) bool {
	// If the operand has lower precedence, then it needs to be parenthesized to preserve the
	// intent of the expression. For example, if the operand is `a + b` and the operator is
	// `*`, then we need to parenthesize the operand to preserve the intended order of
	// operations: `(a + b) * x`.
	//
	// If the operand has higher precedence, then it does not need to be parenthesized. For
	// example, if the operand is `a * b` and the operator is `+`, then we do not need to
	// parenthesize to preserve the intended order of operations: `a * b + x`.
	//
	// If the operand has the same precedence, then we need to check the associativity of
	// the operator based on whether this is the left or right operand of the expression.
	//
	// For example, if `a / d` is on the right of operator `*`, we need to parenthesize
	// to preserve the intended order of operations: `x * (a / d)`
	//
	// If `a ** d` is on the left of operator `**`, we need to parenthesize to preserve
	// the intended order of operations: `(a ** b) ** c`
	binaryOperatorPrecedence := getOperatorPrecedence(SyntaxKindBinaryExpression, binaryOperator, false /*hasArguments*/)
	binaryOperatorAssociativity := getOperatorAssociativity(SyntaxKindBinaryExpression, binaryOperator, false /*hasArguments*/)
	emittedOperand := skipPartiallyEmittedExpressions(operand)
	if !isLeftSideOfBinary && operand.kind == SyntaxKindArrowFunction && binaryOperatorPrecedence > OperatorPrecedenceAssignment {
		// We need to parenthesize arrow functions on the right side to avoid it being
		// parsed as parenthesized expression: `a && (() => {})`
		return true
	}
	operandPrecedence := getExpressionPrecedence(emittedOperand)

	switch {
	case operandPrecedence < binaryOperatorPrecedence:
		// If the operand is the right side of a right-associative binary operation
		// and is a yield expression, then we do not need parentheses.
		if !isLeftSideOfBinary &&
			binaryOperatorAssociativity == AssociativityRight &&
			operand.kind == SyntaxKindYieldExpression {
			return false
		}

		return true

	case operandPrecedence > binaryOperatorPrecedence:
		return false

	default:
		if isLeftSideOfBinary {
			// No need to parenthesize the left operand when the binary operator is
			// left associative:
			//  (a*b)/x    -> a*b/x
			//  (a**b)/x   -> a**b/x
			//
			// Parentheses are needed for the left operand when the binary operator is
			// right associative:
			//  (a/b)**x   -> (a/b)**x
			//  (a**b)**x  -> (a**b)**x
			return binaryOperatorAssociativity == AssociativityRight
		} else {
			if isBinaryExpression(emittedOperand) &&
				emittedOperand.AsBinaryExpression().operatorToken.kind == binaryOperator {
				// No need to parenthesize the right operand when the binary operator and
				// operand are the same and one of the following:
				//  x*(a*b)     => x*a*b
				//  x|(a|b)     => x|a|b
				//  x&(a&b)     => x&a&b
				//  x^(a^b)     => x^a^b
				if p.operatorHasAssociativeProperty(binaryOperator) {
					return false
				}

				// No need to parenthesize the right operand when the binary operator
				// is plus (+) if both the left and right operands consist solely of either
				// literals of the same kind or binary plus (+) expressions for literals of
				// the same kind (recursively).
				//  "a"+(1+2)       => "a"+(1+2)
				//  "a"+("b"+"c")   => "a"+"b"+"c"
				if binaryOperator == SyntaxKindPlusToken {
					leftKind := SyntaxKindUnknown
					if leftOperand != nil {
						leftKind = p.getLiteralKindOfBinaryPlusOperand(leftOperand)
					}
					if isLiteralKind(leftKind) && leftKind == p.getLiteralKindOfBinaryPlusOperand(emittedOperand) {
						return false
					}
				}
			}

			// No need to parenthesize the right operand when the operand is right
			// associative:
			//  x/(a**b)    -> x/a**b
			//  x**(a**b)   -> x**a**b
			//
			// Parentheses are needed for the right operand when the operand is left
			// associative:
			//  x/(a*b)     -> x/(a*b)
			//  x**(a/b)    -> x**(a/b)
			operandAssociativity := getExpressionAssociativity(emittedOperand)
			return operandAssociativity == AssociativityLeft
		}
	}
}

// Determines whether a binary operator is mathematically associative.
func (p *parenthesizerRules) operatorHasAssociativeProperty(binaryOperator SyntaxKind) bool {
	// The following operators are associative in JavaScript:
	//  (a*b)*c     -> a*(b*c)  -> a*b*c
	//  (a|b)|c     -> a|(b|c)  -> a|b|c
	//  (a&b)&c     -> a&(b&c)  -> a&b&c
	//  (a^b)^c     -> a^(b^c)  -> a^b^c
	//  (a,b),c     -> a,(b,c)  -> a,b,c
	//
	// While addition is associative in mathematics, JavaScript's `+` is not
	// guaranteed to be associative as it is overloaded with string concatenation.
	return binaryOperator == SyntaxKindAsteriskToken ||
		binaryOperator == SyntaxKindBarToken ||
		binaryOperator == SyntaxKindAmpersandToken ||
		binaryOperator == SyntaxKindCaretToken ||
		binaryOperator == SyntaxKindCommaToken
}

// This function determines whether an expression consists of a homogeneous set of
// literal expressions or binary plus expressions that all share the same literal kind.
// It is used to determine whether the right-hand operand of a binary plus expression can be
// emitted without parentheses.
func (p *parenthesizerRules) getLiteralKindOfBinaryPlusOperand(node *Expression) SyntaxKind {
	node = skipPartiallyEmittedExpressions(node)

	if isLiteralKind(node.kind) {
		return node.kind
	}

	if node.kind == SyntaxKindBinaryExpression {
		if n := node.AsBinaryExpression(); n.operatorToken.kind == SyntaxKindPlusToken {
			// TODO(rbuckton): Determine if caching this is worthwhile over recomputing
			// if n.cachedLiteralKind != SyntaxKindUnknown {
			// 	return n.cachedLiteralKind;
			// }

			leftKind := p.getLiteralKindOfBinaryPlusOperand(n.left)
			literalKind := SyntaxKindUnknown
			if isLiteralKind(leftKind) && leftKind == p.getLiteralKindOfBinaryPlusOperand(n.right) {
				literalKind = leftKind
			}

			// n.cachedLiteralKind = literalKind;
			return literalKind
		}
	}

	return SyntaxKindUnknown
}

// Wraps the operand to a BinaryExpression in parentheses if they are needed to preserve the intended order of operations.
//   - binaryOperator - The operator for the BinaryExpression.
//   - operand - The operand for the BinaryExpression.
//   - isLeftSideOfBinary - A value indicating whether the operand is the left side of the BinaryExpression.
func (p *parenthesizerRules) parenthesizeBinaryOperand(binaryOperator SyntaxKind, operand *Expression, isLeftSideOfBinary bool, leftOperand *Expression) *Expression {
	skipped := skipPartiallyEmittedExpressions(operand)

	// If the resulting expression is already parenthesized, we do not need to do any further processing.
	if skipped.kind == SyntaxKindParenthesizedExpression {
		return operand
	}

	if p.binaryOperandNeedsParentheses(binaryOperator, operand, isLeftSideOfBinary, leftOperand) {
		return p.factory.NewParenthesizedExpression(operand)
	}

	return operand
}

func (p *parenthesizerRules) ParenthesizeLeftSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression) *Expression {
	return p.parenthesizeBinaryOperand(binaryOperator, leftSide, true /*isLeftSideOfBinary*/, nil /*leftOperand*/)
}

func (p *parenthesizerRules) ParenthesizeRightSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression, rightSide *Expression) *Expression {
	return p.parenthesizeBinaryOperand(binaryOperator, rightSide, false /*isLeftSideOfBinary*/, leftSide)
}

func (p *parenthesizerRules) ParenthesizeExpressionOfComputedPropertyName(expression *Expression) *Expression {
	if isCommaSequence(expression) {
		return p.factory.NewParenthesizedExpression(expression)
	}
	return expression
}

func (p *parenthesizerRules) ParenthesizeConditionOfConditionalExpression(condition *Expression) *Expression {
	conditionalPrecedence := getOperatorPrecedence(SyntaxKindConditionalExpression, SyntaxKindQuestionToken, false /*hasArguments*/)
	emittedCondition := skipPartiallyEmittedExpressions(condition)
	conditionPrecedence := getExpressionPrecedence(emittedCondition)
	if conditionPrecedence <= conditionalPrecedence {
		return p.factory.NewParenthesizedExpression(condition)
	}
	return condition
}

func (p *parenthesizerRules) ParenthesizeBranchOfConditionalExpression(branch *Expression) *Expression {
	// per ES grammar both 'whenTrue' and 'whenFalse' parts of conditional expression are assignment expressions
	// so in case when comma expression is introduced as a part of previous transformations
	// if should be wrapped in parens since comma operator has the lowest precedence
	emittedExpression := skipPartiallyEmittedExpressions(branch)
	if isCommaSequence(emittedExpression) {
		p.factory.NewParenthesizedExpression(branch)
	}
	return branch
}

// [Per the spec](https://tc39.github.io/ecma262/#prod-ExportDeclaration), `export default` accepts _AssigmentExpression_ but
// has a lookahead restriction for `function`, `async function`, and `class`.
//
// Basically, that means we need to parenthesize in the following cases:
//
// - BinaryExpression of CommaToken
// - CommaList (synthetic list of multiple comma expressions)
// - FunctionExpression
// - ClassExpression
func (p *parenthesizerRules) ParenthesizeExpressionOfExportDefault(expression *Expression) *Expression {
	check := skipPartiallyEmittedExpressions(expression)
	needsParens := isCommaSequence(check)
	if !needsParens {
		switch getLeftmostExpression(check, false /*stopAtCallExpressions*/).kind {
		case SyntaxKindClassExpression, SyntaxKindFunctionExpression:
			needsParens = true
		}
	}
	if needsParens {
		return p.factory.NewParenthesizedExpression(expression)
	}
	return expression
}

// Wraps an expression in parentheses if it is needed in order to use the expression
// as the expression of a `NewExpression` node.
func (p *parenthesizerRules) ParenthesizeExpressionOfNew(expression *Expression) *Expression {
	leftmostExpr := getLeftmostExpression(expression /*stopAtCallExpressions*/, true)
	switch leftmostExpr.kind {
	case SyntaxKindCallExpression:
		return p.factory.NewParenthesizedExpression(expression)

	case SyntaxKindNewExpression:
		if leftmostExpr.AsNewExpression().arguments == nil {
			return p.factory.NewParenthesizedExpression(expression)
		}
		return expression
	}

	return p.ParenthesizeLeftSideOfAccess(expression, false /*optionalChain*/)
}

// Wraps an expression in parentheses if it is needed in order to use the expression for
// property or element access.
func (p *parenthesizerRules) ParenthesizeLeftSideOfAccess(expression *Expression, optionalChain bool) *Expression {
	// isLeftHandSideExpression is almost the correct criterion for when it is not necessary
	// to parenthesize the expression before a dot. The known exception is:
	//
	//    NewExpression:
	//       new C.x        -> not the same as (new C).x
	//
	emittedExpression := skipPartiallyEmittedExpressions(expression)
	if isLeftHandSideExpression(emittedExpression) &&
		(emittedExpression.kind != SyntaxKindNewExpression || emittedExpression.AsNewExpression().arguments != nil) &&
		(optionalChain || !isOptionalChain(emittedExpression)) {
		return expression
	}

	result := p.factory.NewParenthesizedExpression(expression)
	result.loc = expression.loc
	return result
}

func (p *parenthesizerRules) ParenthesizeOperandOfPostfixUnary(operand *Expression) *Expression {
	if isLeftHandSideExpression(operand) {
		return operand
	}
	result := p.factory.NewParenthesizedExpression(operand)
	result.loc = operand.loc
	return result
}

func (p *parenthesizerRules) ParenthesizeOperandOfPrefixUnary(operand *Expression) *Expression {
	if isUnaryExpression(operand) {
		return operand
	}
	result := p.factory.NewParenthesizedExpression(operand)
	result.loc = operand.loc
	return result
}

func (p *parenthesizerRules) ParenthesizeExpressionsOfCommaDelimitedList(elements []*Expression) []*Expression {
	return sameMap(elements, p.parenthesizeExpressionForDisallowedComma)
}

func (p *parenthesizerRules) ParenthesizeExpressionForDisallowedComma(expression *Expression) *Expression {
	emittedExpression := skipPartiallyEmittedExpressions(expression)
	expressionPrecedence := getExpressionPrecedence(emittedExpression)
	commaPrecedence := getOperatorPrecedence(SyntaxKindBinaryExpression, SyntaxKindCommaToken, false /*hasArguments*/)
	if expressionPrecedence > commaPrecedence {
		return expression
	}
	result := p.factory.NewParenthesizedExpression(expression)
	result.loc = expression.loc
	return result
}

func (p *parenthesizerRules) ParenthesizeExpressionOfExpressionStatement(expression *Expression) *Expression {
	emittedExpression := skipPartiallyEmittedExpressions(expression)
	if isCallExpression(emittedExpression) {
		callExpression := emittedExpression.AsCallExpression()
		callee := callExpression.expression
		kind := skipPartiallyEmittedExpressions(callee).kind
		if kind == SyntaxKindFunctionExpression || kind == SyntaxKindArrowFunction {
			parenthesizedCallee := p.factory.NewParenthesizedExpression(callee)
			parenthesizedCallee.loc = callee.loc
			updated := p.factory.UpdateCallExpression(
				emittedExpression,
				parenthesizedCallee,
				nil, /*questionDotToken*/
				callExpression.typeArguments,
				callExpression.arguments)
			return p.factory.RestoreOuterExpressions(expression, updated, OEKPartiallyEmittedExpressions)
		}
	}

	leftmostExpressionKind := getLeftmostExpression(emittedExpression, false /*stopAtCallExpressions*/).kind
	if leftmostExpressionKind == SyntaxKindObjectLiteralExpression || leftmostExpressionKind == SyntaxKindFunctionExpression {
		result := p.factory.NewParenthesizedExpression(expression)
		result.loc = expression.loc
		return result
	}

	return expression
}

func (p *parenthesizerRules) ParenthesizeConciseBodyOfArrowFunction(body *BlockOrExpression) *BlockOrExpression {
	if !isBlock(body) && (isCommaSequence(body) || getLeftmostExpression(body, false /*stopAtCallExpressions*/).kind == SyntaxKindObjectLiteralExpression) {
		result := p.factory.NewParenthesizedExpression(body)
		result.loc = body.loc
		return result
	}
	return body
}

// Type[Extends] :
//     FunctionOrConstructorType
//     ConditionalType[?Extends]

// ConditionalType[Extends] :
//
//	UnionType[?Extends]
//	[~Extends] UnionType[~Extends] `extends` Type[+Extends] `?` Type[~Extends] `:` Type[~Extends]
//
// - The check type (the `UnionType`, above) does not allow function, constructor, or conditional types (they must be parenthesized)
// - The extends type (the first `Type`, above) does not allow conditional types (they must be parenthesized). Function and constructor types are fine.
// - The true and false branch types (the second and third `Type` non-terminals, above) allow any type
func (p *parenthesizerRules) ParenthesizeCheckTypeOfConditionalType(checkType *TypeNode) *TypeNode {
	switch checkType.kind {
	case SyntaxKindFunctionType:
	case SyntaxKindConstructorType:
	case SyntaxKindConditionalType:
		return p.factory.NewParenthesizedTypeNode(checkType)
	}
	return checkType
}

func (p *parenthesizerRules) ParenthesizeExtendsTypeOfConditionalType(extendsType *TypeNode) *TypeNode {
	switch extendsType.kind {
	case SyntaxKindConditionalType:
		return p.factory.NewParenthesizedTypeNode(extendsType)
	}
	return extendsType
}

// TypeOperator[Extends] :
//
//	PostfixType
//	InferType[?Extends]
//	`keyof` TypeOperator[?Extends]
//	`unique` TypeOperator[?Extends]
//	`readonly` TypeOperator[?Extends]
func (p *parenthesizerRules) ParenthesizeOperandOfTypeOperator(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindIntersectionType:
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.ParenthesizeConstituentTypeOfIntersectionType(typeNode)
}

func (p *parenthesizerRules) ParenthesizeOperandOfReadonlyTypeOperator(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindTypeOperator:
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.ParenthesizeOperandOfTypeOperator(typeNode)
}

// PostfixType :
//
//	NonArrayType
//	NonArrayType [no LineTerminator here] `!` // JSDoc
//	NonArrayType [no LineTerminator here] `?` // JSDoc
//	IndexedAccessType
//	ArrayType
//
// IndexedAccessType :
//
//	NonArrayType `[` Type[~Extends] `]`
//
// ArrayType :
//
//	NonArrayType `[` `]`
func (p *parenthesizerRules) ParenthesizeNonArrayTypeOfPostfixType(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindInferType:
	case SyntaxKindTypeOperator:
	case SyntaxKindTypeQuery: // Not strictly necessary, but makes generated output more readable and avoids breaks in DT tests
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.ParenthesizeOperandOfTypeOperator(typeNode)
}

// TupleType :
//
//	`[` Elision? `]`
//	`[` NamedTupleElementTypes `]`
//	`[` NamedTupleElementTypes `,` Elision? `]`
//	`[` TupleElementTypes `]`
//	`[` TupleElementTypes `,` Elision? `]`
//
// NamedTupleElementTypes :
//
//	Elision? NamedTupleMember
//	NamedTupleElementTypes `,` Elision? NamedTupleMember
//
// NamedTupleMember :
//
//	Identifier `?`? `:` Type[~Extends]
//	`...` Identifier `:` Type[~Extends]
//
// TupleElementTypes :
//
//	Elision? TupleElementType
//	TupleElementTypes `,` Elision? TupleElementType
//
// TupleElementType :
//
//	Type[~Extends] // NOTE: Needs cover grammar to disallow JSDoc postfix-optional
//	OptionalType
//	RestType
//
// OptionalType :
//
//	Type[~Extends] `?` // NOTE: Needs cover grammar to disallow JSDoc postfix-optional
//
// RestType :
//
//	`...` Type[~Extends]
func (p *parenthesizerRules) ParenthesizeElementTypesOfTupleType(types []*Node) []*Node {
	return sameMap(types, p.parenthesizeElementTypeOfTupleType)
}

func (p *parenthesizerRules) hasJSDocPostfixQuestion(typeNode *Node) bool {
	if typeNode != nil {
		switch typeNode.kind {
		case SyntaxKindJSDocNullableType:
			return typeNode.AsJSDocNullableType().postfix
		case SyntaxKindNamedTupleMember:
			return p.hasJSDocPostfixQuestion(typeNode.AsNamedTupleMember().typeNode)
		case SyntaxKindFunctionType, SyntaxKindConstructorType:
			return p.hasJSDocPostfixQuestion(typeNode.ReturnType())
		case SyntaxKindTypeOperator:
			return p.hasJSDocPostfixQuestion(typeNode.AsTypeOperatorNode().typeNode)
		case SyntaxKindConditionalType:
			return p.hasJSDocPostfixQuestion(typeNode.AsConditionalTypeNode().falseType)
		case SyntaxKindUnionType:
			return p.hasJSDocPostfixQuestion(lastElement(typeNode.AsUnionTypeNode().types))
		case SyntaxKindIntersectionType:
			return p.hasJSDocPostfixQuestion(lastElement(typeNode.AsIntersectionTypeNode().types))
		case SyntaxKindInferType:
			return p.hasJSDocPostfixQuestion(typeNode.AsInferTypeNode().typeParameter.AsTypeParameter().constraint)
		}
	}
	return false
}

func (p *parenthesizerRules) ParenthesizeElementTypeOfTupleType(typeNode *Node) *Node {
	if p.hasJSDocPostfixQuestion(typeNode) {
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return typeNode
}

func (p *parenthesizerRules) ParenthesizeTypeOfOptionalType(typeNode *TypeNode) *TypeNode {
	if p.hasJSDocPostfixQuestion(typeNode) {
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.ParenthesizeNonArrayTypeOfPostfixType(typeNode)
}

// UnionType[Extends] :
//
//	`|`? IntersectionType[?Extends]
//	UnionType[?Extends] `|` IntersectionType[?Extends]
//
// - A union type constituent has the same precedence as the check type of a conditional type
func (p *parenthesizerRules) ParenthesizeConstituentTypeOfUnionType(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindUnionType: // Not strictly necessary, but a union containing a union should have been flattened
	case SyntaxKindIntersectionType: // Not strictly necessary, but makes generated output more readable and avoids breaks in DT tests
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.ParenthesizeCheckTypeOfConditionalType(typeNode)
}

func (p *parenthesizerRules) ParenthesizeConstituentTypesOfUnionType(constituents []*TypeNode) []*TypeNode {
	return sameMap(constituents, p.parenthesizeConstituentTypeOfUnionType)
}

// IntersectionType[Extends] :
//
//	`&`? TypeOperator[?Extends]
//	IntersectionType[?Extends] `&` TypeOperator[?Extends]
//
// - An intersection type constituent does not allow function, constructor, conditional, or union types (they must be parenthesized)
func (p *parenthesizerRules) ParenthesizeConstituentTypeOfIntersectionType(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindUnionType:
	case SyntaxKindIntersectionType: // Not strictly necessary, but an intersection containing an intersection should have been flattened
		return p.factory.NewParenthesizedTypeNode(typeNode)
	}
	return p.parenthesizeConstituentTypeOfUnionType(typeNode)
}

func (p *parenthesizerRules) ParenthesizeConstituentTypesOfIntersectionType(constituents []*TypeNode) []*TypeNode {
	return sameMap(constituents, p.parenthesizeConstituentTypeOfIntersectionType)
}

func (p *parenthesizerRules) ParenthesizeLeadingTypeArgument(typeNode *TypeNode) *TypeNode {
	switch typeNode.kind {
	case SyntaxKindFunctionType, SyntaxKindConstructorType:
		if typeNode.TypeParameters() != nil {
			return p.factory.NewParenthesizedTypeNode(typeNode)
		}
	}
	return typeNode
}

func (p *parenthesizerRules) parenthesizeOrdinalTypeArgumentWorker(node *TypeNode, i int) *TypeNode {
	if i == 0 {
		return p.ParenthesizeLeadingTypeArgument(node)
	}
	return node
}

func (p *parenthesizerRules) ParenthesizeTypeArguments(typeArguments *TypeArgumentListNode) *TypeArgumentListNode {
	if typeArguments != nil {
		typeArgumentList := typeArguments.AsTypeArgumentList()
		arguments, _ := sameMapIndex(typeArgumentList.arguments, p.parenthesizeOrdinalTypeArgument)
		return p.factory.UpdateTypeArgumentList(typeArguments, arguments)
	}
	return nil
}

type nullParenthesizerRules struct {
	identity func(node *Expression) *Expression
}

func NewNullParenthesizerRules() ParenthesizerRules {
	rules := &nullParenthesizerRules{}
	rules.identity = func(node *Expression) *Expression { return node }
	return rules
}

func (p *nullParenthesizerRules) GetParenthesizeLeftSideOfBinaryForOperator(operatorKind SyntaxKind) func(leftSide *Expression) *Expression {
	return p.identity
}

func (p *nullParenthesizerRules) GetParenthesizeRightSideOfBinaryForOperator(operatorKind SyntaxKind) func(rightSide *Expression) *Expression {
	return p.identity
}

func (p *nullParenthesizerRules) ParenthesizeLeftSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression) *Expression {
	return leftSide
}

func (p *nullParenthesizerRules) ParenthesizeRightSideOfBinary(binaryOperator SyntaxKind, leftSide *Expression, rightSide *Expression) *Expression {
	return rightSide
}

func (p *nullParenthesizerRules) ParenthesizeExpressionOfComputedPropertyName(expression *Expression) *Expression {
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeConditionOfConditionalExpression(condition *Expression) *Expression {
	return condition
}

func (p *nullParenthesizerRules) ParenthesizeBranchOfConditionalExpression(branch *Expression) *Expression {
	return branch
}

func (p *nullParenthesizerRules) ParenthesizeExpressionOfExportDefault(expression *Expression) *Expression {
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeExpressionOfNew(expression *Expression) *Expression {
	if !isLeftHandSideExpression(expression) {
		panic("Expected expression to be a valid LeftHandSideExpression")
	}
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeLeftSideOfAccess(expression *Expression, optionalChain bool) *Expression {
	if !isLeftHandSideExpression(expression) {
		panic("Expected expression to be a valid LeftHandSideExpression")
	}
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeOperandOfPostfixUnary(operand *Expression) *Expression {
	if !isLeftHandSideExpression(operand) {
		panic("Expected expression to be a valid LeftHandSideExpression")
	}
	return operand
}

func (p *nullParenthesizerRules) ParenthesizeOperandOfPrefixUnary(operand *Expression) *Expression {
	if !isUnaryExpression(operand) {
		panic("Expected expression to be a valid UnaryExpression")
	}
	return operand
}

func (p *nullParenthesizerRules) ParenthesizeExpressionsOfCommaDelimitedList(elements []*Expression) []*Expression {
	return elements
}

func (p *nullParenthesizerRules) ParenthesizeExpressionForDisallowedComma(expression *Expression) *Expression {
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeExpressionOfExpressionStatement(expression *Expression) *Expression {
	return expression
}

func (p *nullParenthesizerRules) ParenthesizeConciseBodyOfArrowFunction(body *BlockOrExpression) *BlockOrExpression {
	return body
}

func (p *nullParenthesizerRules) ParenthesizeCheckTypeOfConditionalType(checkType *TypeNode) *TypeNode {
	return checkType
}

func (p *nullParenthesizerRules) ParenthesizeExtendsTypeOfConditionalType(extendsType *TypeNode) *TypeNode {
	return extendsType
}

func (p *nullParenthesizerRules) ParenthesizeOperandOfTypeOperator(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeOperandOfReadonlyTypeOperator(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeNonArrayTypeOfPostfixType(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeElementTypesOfTupleType(types []*Node) []*Node {
	return types
}

func (p *nullParenthesizerRules) ParenthesizeElementTypeOfTupleType(typeNode *Node) *Node {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeTypeOfOptionalType(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeConstituentTypeOfUnionType(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeConstituentTypesOfUnionType(constituents []*TypeNode) []*TypeNode {
	return constituents
}

func (p *nullParenthesizerRules) ParenthesizeConstituentTypeOfIntersectionType(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeConstituentTypesOfIntersectionType(constituents []*TypeNode) []*TypeNode {
	return constituents
}

func (p *nullParenthesizerRules) ParenthesizeLeadingTypeArgument(typeNode *TypeNode) *TypeNode {
	return typeNode
}

func (p *nullParenthesizerRules) ParenthesizeTypeArguments(typeArguments *TypeArgumentListNode) *TypeArgumentListNode {
	return typeArguments
}
