package compiler

import (
	"fmt"
	"maps"
	"math"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Links store

type LinkStore[K comparable, V any] struct {
	entries map[K]*V
	pool    core.Pool[V]
}

func (s *LinkStore[K, V]) get(key K) *V {
	value := s.entries[key]
	if value != nil {
		return value
	}
	if s.entries == nil {
		s.entries = make(map[K]*V)
	}
	value = s.pool.New()
	s.entries[key] = value
	return value
}

// Atomic ids

var nextNodeId atomic.Uint32
var nextSymbolId atomic.Uint32
var nextMergeId atomic.Uint32

func getNodeId(node *Node) NodeId {
	if node.Id == 0 {
		node.Id = NodeId(nextNodeId.Add(1))
	}
	return node.Id
}

func getSymbolId(symbol *Symbol) SymbolId {
	if symbol.Id == 0 {
		symbol.Id = SymbolId(nextSymbolId.Add(1))
	}
	return symbol.Id
}

func getMergeId(symbol *Symbol) MergeId {
	if symbol.MergeId == 0 {
		symbol.MergeId = MergeId(nextMergeId.Add(1))
	}
	return symbol.MergeId
}

func NewDiagnostic(file *SourceFile, loc core.TextRange, message *diagnostics.Message, args ...any) *Diagnostic {
	text := message.Message()
	if len(args) != 0 {
		text = formatStringFromArgs(text, args)
	}
	return &Diagnostic{
		File_:     file,
		Loc_:      loc,
		Code_:     message.Code(),
		Category_: message.Category(),
		Message_:  text,
	}
}

func NewDiagnosticForNode(node *Node, message *diagnostics.Message, args ...any) *Diagnostic {
	var file *SourceFile
	var loc core.TextRange
	if node != nil {
		file = getSourceFileOfNode(node)
		loc = getErrorRangeForNode(file, node)
	}
	return NewDiagnostic(file, loc, message, args...)
}

func NewDiagnosticFromMessageChain(file *SourceFile, loc core.TextRange, messageChain *MessageChain) *Diagnostic {
	return &Diagnostic{
		File_:         file,
		Loc_:          loc,
		Code_:         messageChain.Code_,
		Category_:     messageChain.Category_,
		Message_:      messageChain.Message_,
		MessageChain_: messageChain.MessageChain_,
	}
}

func NewDiagnosticForNodeFromMessageChain(node *Node, messageChain *MessageChain) *Diagnostic {
	var file *SourceFile
	var loc core.TextRange
	if node != nil {
		file = getSourceFileOfNode(node)
		loc = getErrorRangeForNode(file, node)
	}
	return NewDiagnosticFromMessageChain(file, loc, messageChain)
}

func NewMessageChain(message *diagnostics.Message, args ...any) *MessageChain {
	text := message.Message()
	if len(args) != 0 {
		text = formatStringFromArgs(text, args)
	}
	return &MessageChain{
		Code_:     message.Code(),
		Category_: message.Category(),
		Message_:  text,
	}
}

func chainDiagnosticMessages(details *MessageChain, message *diagnostics.Message, args ...any) *MessageChain {
	return NewMessageChain(message, args...).AddMessageChain(details)
}

type OperatorPrecedence int

const (
	// Expression:
	//     AssignmentExpression
	//     Expression `,` AssignmentExpression
	OperatorPrecedenceComma OperatorPrecedence = iota
	// NOTE: `Spread` is higher than `Comma` due to how it is parsed in |ElementList|
	// SpreadElement:
	//     `...` AssignmentExpression
	OperatorPrecedenceSpread
	// AssignmentExpression:
	//     ConditionalExpression
	//     YieldExpression
	//     ArrowFunction
	//     AsyncArrowFunction
	//     LeftHandSideExpression `=` AssignmentExpression
	//     LeftHandSideExpression AssignmentOperator AssignmentExpression
	//
	// NOTE: AssignmentExpression is broken down into several precedences due to the requirements
	//       of the parenthesizer rules.
	// AssignmentExpression: YieldExpression
	// YieldExpression:
	//     `yield`
	//     `yield` AssignmentExpression
	//     `yield` `*` AssignmentExpression
	OperatorPrecedenceYield
	// AssignmentExpression: LeftHandSideExpression `=` AssignmentExpression
	// AssignmentExpression: LeftHandSideExpression AssignmentOperator AssignmentExpression
	// AssignmentOperator: one of
	//     `*=` `/=` `%=` `+=` `-=` `<<=` `>>=` `>>>=` `&=` `^=` `|=` `**=`
	OperatorPrecedenceAssignment
	// NOTE: `Conditional` is considered higher than `Assignment` here, but in reality they have
	//       the same precedence.
	// AssignmentExpression: ConditionalExpression
	// ConditionalExpression:
	//     ShortCircuitExpression
	//     ShortCircuitExpression `?` AssignmentExpression `:` AssignmentExpression
	// ShortCircuitExpression:
	//     LogicalORExpression
	//     CoalesceExpression
	OperatorPrecedenceConditional
	// LogicalORExpression:
	//     LogicalANDExpression
	//     LogicalORExpression `||` LogicalANDExpression
	OperatorPrecedenceLogicalOR
	// LogicalANDExpression:
	//     BitwiseORExpression
	//     LogicalANDExprerssion `&&` BitwiseORExpression
	OperatorPrecedenceLogicalAND
	// BitwiseORExpression:
	//     BitwiseXORExpression
	//     BitwiseORExpression `^` BitwiseXORExpression
	OperatorPrecedenceBitwiseOR
	// BitwiseXORExpression:
	//     BitwiseANDExpression
	//     BitwiseXORExpression `^` BitwiseANDExpression
	OperatorPrecedenceBitwiseXOR
	// BitwiseANDExpression:
	//     EqualityExpression
	//     BitwiseANDExpression `^` EqualityExpression
	OperatorPrecedenceBitwiseAND
	// EqualityExpression:
	//     RelationalExpression
	//     EqualityExpression `==` RelationalExpression
	//     EqualityExpression `!=` RelationalExpression
	//     EqualityExpression `===` RelationalExpression
	//     EqualityExpression `!==` RelationalExpression
	OperatorPrecedenceEquality
	// RelationalExpression:
	//     ShiftExpression
	//     RelationalExpression `<` ShiftExpression
	//     RelationalExpression `>` ShiftExpression
	//     RelationalExpression `<=` ShiftExpression
	//     RelationalExpression `>=` ShiftExpression
	//     RelationalExpression `instanceof` ShiftExpression
	//     RelationalExpression `in` ShiftExpression
	//     [+TypeScript] RelationalExpression `as` Type
	OperatorPrecedenceRelational
	// ShiftExpression:
	//     AdditiveExpression
	//     ShiftExpression `<<` AdditiveExpression
	//     ShiftExpression `>>` AdditiveExpression
	//     ShiftExpression `>>>` AdditiveExpression
	OperatorPrecedenceShift
	// AdditiveExpression:
	//     MultiplicativeExpression
	//     AdditiveExpression `+` MultiplicativeExpression
	//     AdditiveExpression `-` MultiplicativeExpression
	OperatorPrecedenceAdditive
	// MultiplicativeExpression:
	//     ExponentiationExpression
	//     MultiplicativeExpression MultiplicativeOperator ExponentiationExpression
	// MultiplicativeOperator: one of `*`, `/`, `%`
	OperatorPrecedenceMultiplicative
	// ExponentiationExpression:
	//     UnaryExpression
	//     UpdateExpression `**` ExponentiationExpression
	OperatorPrecedenceExponentiation
	// UnaryExpression:
	//     UpdateExpression
	//     `delete` UnaryExpression
	//     `void` UnaryExpression
	//     `typeof` UnaryExpression
	//     `+` UnaryExpression
	//     `-` UnaryExpression
	//     `~` UnaryExpression
	//     `!` UnaryExpression
	//     AwaitExpression
	// UpdateExpression:            // TODO: Do we need to investigate the precedence here?
	//     `++` UnaryExpression
	//     `--` UnaryExpression
	OperatorPrecedenceUnary
	// UpdateExpression:
	//     LeftHandSideExpression
	//     LeftHandSideExpression `++`
	//     LeftHandSideExpression `--`
	OperatorPrecedenceUpdate
	// LeftHandSideExpression:
	//     NewExpression
	//     CallExpression
	// NewExpression:
	//     MemberExpression
	//     `new` NewExpression
	OperatorPrecedenceLeftHandSide
	// CallExpression:
	//     CoverCallExpressionAndAsyncArrowHead
	//     SuperCall
	//     ImportCall
	//     CallExpression Arguments
	//     CallExpression `[` Expression `]`
	//     CallExpression `.` IdentifierName
	//     CallExpression TemplateLiteral
	// MemberExpression:
	//     PrimaryExpression
	//     MemberExpression `[` Expression `]`
	//     MemberExpression `.` IdentifierName
	//     MemberExpression TemplateLiteral
	//     SuperProperty
	//     MetaProperty
	//     `new` MemberExpression Arguments
	OperatorPrecedenceMember
	// TODO: JSXElement?
	// PrimaryExpression:
	//     `this`
	//     IdentifierReference
	//     Literal
	//     ArrayLiteral
	//     ObjectLiteral
	//     FunctionExpression
	//     ClassExpression
	//     GeneratorExpression
	//     AsyncFunctionExpression
	//     AsyncGeneratorExpression
	//     RegularExpressionLiteral
	//     TemplateLiteral
	//     CoverParenthesizedExpressionAndArrowParameterList
	OperatorPrecedencePrimary
	// CoalesceExpression:
	//     CoalesceExpressionHead `??` BitwiseORExpression
	// CoalesceExpressionHead:
	//     CoalesceExpression
	//     BitwiseORExpression
	OperatorPrecedenceCoalesce = OperatorPrecedenceConditional // NOTE: This is wrong
	OperatorPrecedenceLowest   = OperatorPrecedenceComma
	OperatorPrecedenceHighest  = OperatorPrecedencePrimary
	// -1 is lower than all other precedences. Returning it will cause binary expression
	// parsing to stop.
	OperatorPrecedenceInvalid OperatorPrecedence = -1
)

func getOperatorPrecedence(nodeKind ast.Kind, operatorKind ast.Kind, hasArguments bool) OperatorPrecedence {
	switch nodeKind {
	case ast.KindCommaListExpression:
		return OperatorPrecedenceComma
	case ast.KindSpreadElement:
		return OperatorPrecedenceSpread
	case ast.KindYieldExpression:
		return OperatorPrecedenceYield
	case ast.KindConditionalExpression:
		return OperatorPrecedenceConditional
	case ast.KindBinaryExpression:
		switch operatorKind {
		case ast.KindCommaToken:
			return OperatorPrecedenceComma
		case ast.KindEqualsToken, ast.KindPlusEqualsToken, ast.KindMinusEqualsToken, ast.KindAsteriskAsteriskEqualsToken,
			ast.KindAsteriskEqualsToken, ast.KindSlashEqualsToken, ast.KindPercentEqualsToken, ast.KindLessThanLessThanEqualsToken,
			ast.KindGreaterThanGreaterThanEqualsToken, ast.KindGreaterThanGreaterThanGreaterThanEqualsToken, ast.KindAmpersandEqualsToken,
			ast.KindCaretEqualsToken, ast.KindBarEqualsToken, ast.KindBarBarEqualsToken, ast.KindAmpersandAmpersandEqualsToken,
			ast.KindQuestionQuestionEqualsToken:
			return OperatorPrecedenceAssignment
		}
		return getBinaryOperatorPrecedence(operatorKind)
	// TODO: Should prefix `++` and `--` be moved to the `Update` precedence?
	case ast.KindTypeAssertionExpression, ast.KindNonNullExpression, ast.KindPrefixUnaryExpression, ast.KindTypeOfExpression,
		ast.KindVoidExpression, ast.KindDeleteExpression, ast.KindAwaitExpression:
		return OperatorPrecedenceUnary
	case ast.KindPostfixUnaryExpression:
		return OperatorPrecedenceUpdate
	case ast.KindCallExpression:
		return OperatorPrecedenceLeftHandSide
	case ast.KindNewExpression:
		if hasArguments {
			return OperatorPrecedenceMember
		}
		return OperatorPrecedenceLeftHandSide
	case ast.KindTaggedTemplateExpression, ast.KindPropertyAccessExpression, ast.KindElementAccessExpression, ast.KindMetaProperty:
		return OperatorPrecedenceMember
	case ast.KindAsExpression, ast.KindSatisfiesExpression:
		return OperatorPrecedenceRelational
	case ast.KindThisKeyword, ast.KindSuperKeyword, ast.KindIdentifier, ast.KindPrivateIdentifier, ast.KindNullKeyword,
		ast.KindTrueKeyword, ast.KindFalseKeyword, ast.KindNumericLiteral, ast.KindBigIntLiteral, ast.KindStringLiteral,
		ast.KindArrayLiteralExpression, ast.KindObjectLiteralExpression, ast.KindFunctionExpression, ast.KindArrowFunction,
		ast.KindClassExpression, ast.KindRegularExpressionLiteral, ast.KindNoSubstitutionTemplateLiteral, ast.KindTemplateExpression,
		ast.KindParenthesizedExpression, ast.KindOmittedExpression, ast.KindJsxElement, ast.KindJsxSelfClosingElement, ast.KindJsxFragment:
		return OperatorPrecedencePrimary
	}
	return OperatorPrecedenceInvalid
}

func getBinaryOperatorPrecedence(kind ast.Kind) OperatorPrecedence {
	switch kind {
	case ast.KindQuestionQuestionToken:
		return OperatorPrecedenceCoalesce
	case ast.KindBarBarToken:
		return OperatorPrecedenceLogicalOR
	case ast.KindAmpersandAmpersandToken:
		return OperatorPrecedenceLogicalAND
	case ast.KindBarToken:
		return OperatorPrecedenceBitwiseOR
	case ast.KindCaretToken:
		return OperatorPrecedenceBitwiseXOR
	case ast.KindAmpersandToken:
		return OperatorPrecedenceBitwiseAND
	case ast.KindEqualsEqualsToken, ast.KindExclamationEqualsToken, ast.KindEqualsEqualsEqualsToken, ast.KindExclamationEqualsEqualsToken:
		return OperatorPrecedenceEquality
	case ast.KindLessThanToken, ast.KindGreaterThanToken, ast.KindLessThanEqualsToken, ast.KindGreaterThanEqualsToken,
		ast.KindInstanceOfKeyword, ast.KindInKeyword, ast.KindAsKeyword, ast.KindSatisfiesKeyword:
		return OperatorPrecedenceRelational
	case ast.KindLessThanLessThanToken, ast.KindGreaterThanGreaterThanToken, ast.KindGreaterThanGreaterThanGreaterThanToken:
		return OperatorPrecedenceShift
	case ast.KindPlusToken, ast.KindMinusToken:
		return OperatorPrecedenceAdditive
	case ast.KindAsteriskToken, ast.KindSlashToken, ast.KindPercentToken:
		return OperatorPrecedenceMultiplicative
	case ast.KindAsteriskAsteriskToken:
		return OperatorPrecedenceExponentiation
	}
	// -1 is lower than all other precedences.  Returning it will cause binary expression
	// parsing to stop.
	return OperatorPrecedenceInvalid
}

func formatStringFromArgs(text string, args []any) string {
	return core.MakeRegexp(`{(\d+)}`).ReplaceAllStringFunc(text, func(match string) string {
		index, err := strconv.ParseInt(match[1:len(match)-1], 10, 0)
		if err != nil || int(index) >= len(args) {
			panic("Invalid formatting placeholder")
		}
		return fmt.Sprintf("%v", args[int(index)])
	})
}

func formatMessage(message *diagnostics.Message, args ...any) string {
	text := message.Message()
	if len(args) != 0 {
		text = formatStringFromArgs(text, args)
	}
	return text
}

func findInMap[K comparable, V any](m map[K]V, predicate func(V) bool) V {
	for _, value := range m {
		if predicate(value) {
			return value
		}
	}
	return *new(V)
}

func boolToTristate(b bool) core.Tristate {
	if b {
		return core.TSTrue
	}
	return core.TSFalse
}

func modifierToFlag(token ast.Kind) ast.ModifierFlags {
	switch token {
	case ast.KindStaticKeyword:
		return ast.ModifierFlagsStatic
	case ast.KindPublicKeyword:
		return ast.ModifierFlagsPublic
	case ast.KindProtectedKeyword:
		return ast.ModifierFlagsProtected
	case ast.KindPrivateKeyword:
		return ast.ModifierFlagsPrivate
	case ast.KindAbstractKeyword:
		return ast.ModifierFlagsAbstract
	case ast.KindAccessorKeyword:
		return ast.ModifierFlagsAccessor
	case ast.KindExportKeyword:
		return ast.ModifierFlagsExport
	case ast.KindDeclareKeyword:
		return ast.ModifierFlagsAmbient
	case ast.KindConstKeyword:
		return ast.ModifierFlagsConst
	case ast.KindDefaultKeyword:
		return ast.ModifierFlagsDefault
	case ast.KindAsyncKeyword:
		return ast.ModifierFlagsAsync
	case ast.KindReadonlyKeyword:
		return ast.ModifierFlagsReadonly
	case ast.KindOverrideKeyword:
		return ast.ModifierFlagsOverride
	case ast.KindInKeyword:
		return ast.ModifierFlagsIn
	case ast.KindOutKeyword:
		return ast.ModifierFlagsOut
	case ast.KindImmediateKeyword:
		return ast.ModifierFlagsImmediate
	case ast.KindDecorator:
		return ast.ModifierFlagsDecorator
	}
	return ast.ModifierFlagsNone
}

func modifiersToFlags(modifierList *Node) ast.ModifierFlags {
	flags := ast.ModifierFlagsNone
	if modifierList != nil {
		for _, modifier := range modifierList.AsModifierList().Modifiers_ {
			flags |= modifierToFlag(modifier.Kind)
		}
	}
	return flags
}

func nodeIsMissing(node *Node) bool {
	return node == nil || node.Loc.Pos_ == node.Loc.End_ && node.Loc.Pos_ >= 0 && node.Kind != ast.KindEndOfFile
}

func nodeIsPresent(node *Node) bool {
	return !nodeIsMissing(node)
}

func isLeftHandSideExpression(node *Node) bool {
	return isLeftHandSideExpressionKind(node.Kind)
}

func isLeftHandSideExpressionKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindPropertyAccessExpression, ast.KindElementAccessExpression, ast.KindNewExpression, ast.KindCallExpression,
		ast.KindJsxElement, ast.KindJsxSelfClosingElement, ast.KindJsxFragment, ast.KindTaggedTemplateExpression, ast.KindArrayLiteralExpression,
		ast.KindParenthesizedExpression, ast.KindObjectLiteralExpression, ast.KindClassExpression, ast.KindFunctionExpression, ast.KindIdentifier,
		ast.KindPrivateIdentifier, ast.KindRegularExpressionLiteral, ast.KindNumericLiteral, ast.KindBigIntLiteral, ast.KindStringLiteral,
		ast.KindNoSubstitutionTemplateLiteral, ast.KindTemplateExpression, ast.KindFalseKeyword, ast.KindNullKeyword, ast.KindThisKeyword,
		ast.KindTrueKeyword, ast.KindSuperKeyword, ast.KindNonNullExpression, ast.KindExpressionWithTypeArguments, ast.KindMetaProperty,
		ast.KindImportKeyword, ast.KindMissingDeclaration:
		return true
	}
	return false
}

func isUnaryExpression(node *Node) bool {
	return isUnaryExpressionKind(node.Kind)
}

func isUnaryExpressionKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindPrefixUnaryExpression, ast.KindPostfixUnaryExpression, ast.KindDeleteExpression, ast.KindTypeOfExpression,
		ast.KindVoidExpression, ast.KindAwaitExpression, ast.KindTypeAssertionExpression:
		return true
	}
	return isLeftHandSideExpressionKind(kind)
}

/**
 * Determines whether a node is an expression based only on its kind.
 */
func isExpression(node *Node) bool {
	return isExpressionKind(node.Kind)
}

func isExpressionKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindConditionalExpression, ast.KindYieldExpression, ast.KindArrowFunction, ast.KindBinaryExpression,
		ast.KindSpreadElement, ast.KindAsExpression, ast.KindOmittedExpression, ast.KindCommaListExpression,
		ast.KindPartiallyEmittedExpression, ast.KindSatisfiesExpression:
		return true
	}
	return isUnaryExpressionKind(kind)
}

func isAssignmentOperator(token ast.Kind) bool {
	return token >= ast.KindFirstAssignment && token <= ast.KindLastAssignment
}

func isExpressionWithTypeArguments(node *Node) bool {
	return node.Kind == ast.KindExpressionWithTypeArguments
}

func isNonNullExpression(node *Node) bool {
	return node.Kind == ast.KindNonNullExpression
}

func isStringLiteralLike(node *Node) bool {
	return node.Kind == ast.KindStringLiteral || node.Kind == ast.KindNoSubstitutionTemplateLiteral
}

func isNumericLiteral(node *Node) bool {
	return node.Kind == ast.KindNumericLiteral
}

func isStringOrNumericLiteralLike(node *Node) bool {
	return isStringLiteralLike(node) || isNumericLiteral(node)
}

func isSignedNumericLiteral(node *Node) bool {
	if node.Kind == ast.KindPrefixUnaryExpression {
		node := node.AsPrefixUnaryExpression()
		return (node.Operator == ast.KindPlusToken || node.Operator == ast.KindMinusToken) && isNumericLiteral(node.Operand)
	}
	return false
}

func ifElse[T any](b bool, whenTrue T, whenFalse T) T {
	if b {
		return whenTrue
	}
	return whenFalse
}

func tokenIsIdentifierOrKeyword(token ast.Kind) bool {
	return token >= ast.KindIdentifier
}

func tokenIsIdentifierOrKeywordOrGreaterThan(token ast.Kind) bool {
	return token == ast.KindGreaterThanToken || tokenIsIdentifierOrKeyword(token)
}

func getTextOfNode(node *Node) string {
	return getSourceTextOfNodeFromSourceFile(getSourceFileOfNode(node), node)
}

func getSourceTextOfNodeFromSourceFile(sourceFile *SourceFile, node *Node) string {
	return getTextOfNodeFromSourceText(sourceFile.Text, node)
}

func getTextOfNodeFromSourceText(sourceText string, node *Node) string {
	if nodeIsMissing(node) {
		return ""
	}
	text := sourceText[skipTrivia(sourceText, node.Pos()):node.End()]
	// if (isJSDocTypeExpressionOrChild(node)) {
	//     // strip space + asterisk at line start
	//     text = text.split(/\r\n|\n|\r/).map(line => line.replace(/^\s*\*/, "").trimStart()).join("\n");
	// }
	return text
}

func isAssignmentDeclaration(decl *Node) bool {
	return isBinaryExpression(decl) || isAccessExpression(decl) || IsIdentifier(decl) || IsCallExpression(decl)
}

func isBinaryExpression(node *Node) bool {
	return node.Kind == ast.KindBinaryExpression
}

func isAccessExpression(node *Node) bool {
	return node.Kind == ast.KindPropertyAccessExpression || node.Kind == ast.KindElementAccessExpression
}

func isInJSFile(node *Node) bool {
	return node != nil && node.Flags&ast.NodeFlagsJavaScriptFile != 0
}

func isEffectiveModuleDeclaration(node *Node) bool {
	return IsModuleDeclaration(node) || IsIdentifier(node)
}

func isObjectLiteralOrClassExpressionMethodOrAccessor(node *Node) bool {
	kind := node.Kind
	return (kind == ast.KindMethodDeclaration || kind == ast.KindGetAccessor || kind == ast.KindSetAccessor) &&
		(node.Parent.Kind == ast.KindObjectLiteralExpression || node.Parent.Kind == ast.KindClassExpression)
}

func isFunctionLike(node *Node) bool {
	return node != nil && isFunctionLikeKind(node.Kind)
}

func isFunctionLikeKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindMethodSignature, ast.KindCallSignature, ast.KindJSDocSignature, ast.KindConstructSignature, ast.KindIndexSignature,
		ast.KindFunctionType, ast.KindJSDocFunctionType, ast.KindConstructorType:
		return true
	}
	return isFunctionLikeDeclarationKind(kind)
}

func isFunctionLikeDeclaration(node *Node) bool {
	return node != nil && isFunctionLikeDeclarationKind(node.Kind)
}

func isFunctionLikeDeclarationKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindFunctionDeclaration, ast.KindMethodDeclaration, ast.KindConstructor, ast.KindGetAccessor, ast.KindSetAccessor,
		ast.KindFunctionExpression, ast.KindArrowFunction:
		return true
	}
	return false
}

type OuterExpressionKinds int16

const (
	OEKParentheses                  OuterExpressionKinds = 1 << 0
	OEKTypeAssertions               OuterExpressionKinds = 1 << 1
	OEKNonNullAssertions            OuterExpressionKinds = 1 << 2
	OEKExpressionsWithTypeArguments OuterExpressionKinds = 1 << 3
	OEKExcludeJSDocTypeAssertion                         = 1 << 4
	OEKAssertions                                        = OEKTypeAssertions | OEKNonNullAssertions
	OEKAll                                               = OEKParentheses | OEKAssertions | OEKExpressionsWithTypeArguments
)

func isOuterExpression(node *Node, kinds OuterExpressionKinds) bool {
	switch node.Kind {
	case ast.KindParenthesizedExpression:
		return kinds&OEKParentheses != 0 && !(kinds&OEKExcludeJSDocTypeAssertion != 0 && isJSDocTypeAssertion(node))
	case ast.KindTypeAssertionExpression, ast.KindAsExpression, ast.KindSatisfiesExpression:
		return kinds&OEKTypeAssertions != 0
	case ast.KindExpressionWithTypeArguments:
		return kinds&OEKExpressionsWithTypeArguments != 0
	case ast.KindNonNullExpression:
		return kinds&OEKNonNullAssertions != 0
	}
	return false
}

func skipOuterExpressions(node *Node, kinds OuterExpressionKinds) *Node {
	for isOuterExpression(node, kinds) {
		node = node.Expression()
	}
	return node
}

func skipParentheses(node *Node) *Node {
	return skipOuterExpressions(node, OEKParentheses)
}

func walkUpParenthesizedTypes(node *Node) *Node {
	for node != nil && node.Kind == ast.KindParenthesizedType {
		node = node.Parent
	}
	return node
}

func walkUpParenthesizedExpressions(node *Node) *Node {
	for node != nil && node.Kind == ast.KindParenthesizedExpression {
		node = node.Parent
	}
	return node
}

func isJSDocTypeAssertion(node *Node) bool {
	return false // !!!
}

// Return true if the given identifier is classified as an IdentifierName
func isIdentifierName(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindPropertyDeclaration, ast.KindPropertySignature, ast.KindMethodDeclaration, ast.KindMethodSignature, ast.KindGetAccessor,
		ast.KindSetAccessor, ast.KindEnumMember, ast.KindPropertyAssignment, ast.KindPropertyAccessExpression:
		return parent.Name() == node
	case ast.KindQualifiedName:
		return parent.AsQualifiedName().Right == node
	case ast.KindBindingElement:
		return parent.AsBindingElement().PropertyName == node
	case ast.KindImportSpecifier:
		return parent.AsImportSpecifier().PropertyName == node
	case ast.KindExportSpecifier, ast.KindJsxAttribute, ast.KindJsxSelfClosingElement, ast.KindJsxOpeningElement, ast.KindJsxClosingElement:
		return true
	}
	return false
}

func getSourceFileOfNode(node *Node) *SourceFile {
	for {
		if node == nil {
			return nil
		}
		if node.Kind == ast.KindSourceFile {
			return node.Data.(*SourceFile)
		}
		node = node.Parent
	}
}

/** @internal */
func getErrorRangeForNode(sourceFile *SourceFile, node *Node) core.TextRange {
	errorNode := node
	switch node.Kind {
	case ast.KindSourceFile:
		pos := skipTrivia(sourceFile.Text, 0)
		if pos == len(sourceFile.Text) {
			return core.NewTextRange(0, 0)
		}
		return getRangeOfTokenAtPosition(sourceFile, pos)
	// This list is a work in progress. Add missing node kinds to improve their error spans
	case ast.KindVariableDeclaration, ast.KindBindingElement, ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration,
		ast.KindModuleDeclaration, ast.KindEnumDeclaration, ast.KindEnumMember, ast.KindFunctionDeclaration, ast.KindFunctionExpression,
		ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindTypeAliasDeclaration, ast.KindPropertyDeclaration,
		ast.KindPropertySignature, ast.KindNamespaceImport:
		errorNode = getNameOfDeclaration(node)
	case ast.KindArrowFunction:
		return getErrorRangeForArrowFunction(sourceFile, node)
	case ast.KindCaseClause:
	case ast.KindDefaultClause:
		start := skipTrivia(sourceFile.Text, node.Pos())
		end := node.End()
		statements := node.Data.(*CaseOrDefaultClause).Statements
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return core.NewTextRange(start, end)
	case ast.KindReturnStatement, ast.KindYieldExpression:
		pos := skipTrivia(sourceFile.Text, node.Pos())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindSatisfiesExpression:
		pos := skipTrivia(sourceFile.Text, node.AsSatisfiesExpression().Expression.End())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindConstructor:
		scanner := getScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.tokenStart
		for scanner.token != ast.KindConstructorKeyword && scanner.token != ast.KindStringLiteral && scanner.token != ast.KindEndOfFile {
			scanner.Scan()
		}
		return core.NewTextRange(start, scanner.pos)
		// !!!
		// case ast.KindJSDocSatisfiesTag:
		// 	pos := skipTrivia(sourceFile.text, node.tagName.pos)
		// 	return getRangeOfTokenAtPosition(sourceFile, pos)
	}
	if errorNode == nil {
		// If we don't have a better node, then just set the error on the first token of
		// construct.
		return getRangeOfTokenAtPosition(sourceFile, node.Pos())
	}
	pos := errorNode.Pos()
	if !nodeIsMissing(errorNode) {
		pos = skipTrivia(sourceFile.Text, pos)
	}
	return core.NewTextRange(pos, errorNode.End())
}

func getErrorRangeForArrowFunction(sourceFile *SourceFile, node *Node) core.TextRange {
	pos := skipTrivia(sourceFile.Text, node.Pos())
	body := node.AsArrowFunction().Body
	if body != nil && body.Kind == ast.KindBlock {
		startLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.Pos())
		endLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.End())
		if startLine < endLine {
			// The arrow function spans multiple lines,
			// make the error span be the first line, inclusive.
			return core.NewTextRange(pos, getEndLinePosition(sourceFile, startLine))
		}
	}
	return core.NewTextRange(pos, node.End())
}

func getContainingClass(node *Node) *Node {
	return findAncestor(node.Parent, isClassLike)
}

func findAncestor(node *Node, callback func(*Node) bool) *Node {
	for node != nil {
		result := callback(node)
		if result {
			return node
		}
		node = node.Parent
	}
	return nil
}

type FindAncestorResult int32

const (
	FindAncestorFalse FindAncestorResult = iota
	FindAncestorTrue
	FindAncestorQuit
)

func findAncestorOrQuit(node *Node, callback func(*Node) FindAncestorResult) *Node {
	for node != nil {
		switch callback(node) {
		case FindAncestorQuit:
			return nil
		case FindAncestorTrue:
			return node
		}
		node = node.Parent
	}
	return nil
}

func isClassLike(node *Node) bool {
	return node != nil && (node.Kind == ast.KindClassDeclaration || node.Kind == ast.KindClassExpression)
}

func declarationNameToString(name *Node) string {
	if name == nil || name.Pos() == name.End() {
		return "(Missing)"
	}
	return getTextOfNode(name)
}

func isExternalModule(file *SourceFile) bool {
	return file.ExternalModuleIndicator != nil
}

func isInTopLevelContext(node *Node) bool {
	// The name of a class or function declaration is a BindingIdentifier in its surrounding scope.
	if IsIdentifier(node) {
		parent := node.Parent
		if (IsClassDeclaration(parent) || IsFunctionDeclaration(parent)) && parent.Name() == node {
			node = parent
		}
	}
	container := getThisContainer(node, true /*includeArrowFunctions*/, false /*includeClassComputedPropertyName*/)
	return IsSourceFile(container)
}

func getThisContainer(node *Node, includeArrowFunctions bool, includeClassComputedPropertyName bool) *Node {
	for {
		node = node.Parent
		if node == nil {
			panic("nil parent in getThisContainer")
		}
		switch node.Kind {
		case ast.KindComputedPropertyName:
			if includeClassComputedPropertyName && isClassLike(node.Parent.Parent) {
				return node
			}
			node = node.Parent.Parent
		case ast.KindDecorator:
			if node.Parent.Kind == ast.KindParameter && isClassElement(node.Parent.Parent) {
				// If the decorator's parent is a Parameter, we resolve the this container from
				// the grandparent class declaration.
				node = node.Parent.Parent
			} else if isClassElement(node.Parent) {
				// If the decorator's parent is a class element, we resolve the 'this' container
				// from the parent class declaration.
				node = node.Parent
			}
		case ast.KindArrowFunction:
			if includeArrowFunctions {
				return node
			}
		case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindModuleDeclaration, ast.KindClassStaticBlockDeclaration,
			ast.KindPropertyDeclaration, ast.KindPropertySignature, ast.KindMethodDeclaration, ast.KindMethodSignature, ast.KindConstructor,
			ast.KindGetAccessor, ast.KindSetAccessor, ast.KindCallSignature, ast.KindConstructSignature, ast.KindIndexSignature,
			ast.KindEnumDeclaration, ast.KindSourceFile:
			return node
		}
	}
}

func isClassElement(node *Node) bool {
	switch node.Kind {
	case ast.KindConstructor, ast.KindPropertyDeclaration, ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor,
		ast.KindIndexSignature, ast.KindClassStaticBlockDeclaration, ast.KindSemicolonClassElement:
		return true
	}
	return false
}

func isPartOfTypeQuery(node *Node) bool {
	for node.Kind == ast.KindQualifiedName || node.Kind == ast.KindIdentifier {
		node = node.Parent
	}
	return node.Kind == ast.KindTypeQuery
}

func getModifierFlags(node *Node) ast.ModifierFlags {
	modifiers := node.Modifiers()
	if modifiers != nil {
		return modifiers.AsModifierList().ModifierFlags
	}
	return ast.ModifierFlagsNone
}

func getNodeFlags(node *Node) ast.NodeFlags {
	return node.Flags
}

func hasSyntacticModifier(node *Node, flags ast.ModifierFlags) bool {
	return getModifierFlags(node)&flags != 0
}

func hasAccessorModifier(node *Node) bool {
	return hasSyntacticModifier(node, ast.ModifierFlagsAccessor)
}

func hasStaticModifier(node *Node) bool {
	return hasSyntacticModifier(node, ast.ModifierFlagsStatic)
}

func getEffectiveModifierFlags(node *Node) ast.ModifierFlags {
	return getModifierFlags(node) // !!! Handle JSDoc
}

func hasEffectiveModifier(node *Node, flags ast.ModifierFlags) bool {
	return getEffectiveModifierFlags(node)&flags != 0
}

func hasEffectiveReadonlyModifier(node *Node) bool {
	return hasEffectiveModifier(node, ast.ModifierFlagsReadonly)
}

func getImmediatelyInvokedFunctionExpression(fn *Node) *Node {
	if fn.Kind == ast.KindFunctionExpression || fn.Kind == ast.KindArrowFunction {
		prev := fn
		parent := fn.Parent
		for parent.Kind == ast.KindParenthesizedExpression {
			prev = parent
			parent = parent.Parent
		}
		if parent.Kind == ast.KindCallExpression && parent.AsCallExpression().Expression == prev {
			return parent
		}
	}
	return nil
}

// Does not handle signed numeric names like `a[+0]` - handling those would require handling prefix unary expressions
// throughout late binding handling as well, which is awkward (but ultimately probably doable if there is demand)
func getElementOrPropertyAccessArgumentExpressionOrName(node *Node) *Node {
	switch node.Kind {
	case ast.KindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().Name_
	case ast.KindElementAccessExpression:
		arg := skipParentheses(node.AsElementAccessExpression().ArgumentExpression)
		if isStringOrNumericLiteralLike(arg) {
			return arg
		}
		return node
	}
	panic("Unhandled case in getElementOrPropertyAccessArgumentExpressionOrName")
}

func getQuestionDotToken(node *Node) *Node {
	switch node.Kind {
	case ast.KindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().QuestionDotToken
	case ast.KindElementAccessExpression:
		return node.AsElementAccessExpression().QuestionDotToken
	case ast.KindCallExpression:
		return node.AsCallExpression().QuestionDotToken
	}
	panic("Unhandled case in getQuestionDotToken")
}

/**
 * A declaration has a dynamic name if all of the following are true:
 *   1. The declaration has a computed property name.
 *   2. The computed name is *not* expressed as a StringLiteral.
 *   3. The computed name is *not* expressed as a NumericLiteral.
 *   4. The computed name is *not* expressed as a PlusToken or MinusToken
 *      immediately followed by a NumericLiteral.
 */
func hasDynamicName(declaration *Node) bool {
	name := getNameOfDeclaration(declaration)
	return name != nil && isDynamicName(name)
}

func isDynamicName(name *Node) bool {
	var expr *Node
	switch name.Kind {
	case ast.KindComputedPropertyName:
		expr = name.AsComputedPropertyName().Expression
	case ast.KindElementAccessExpression:
		expr = skipParentheses(name.AsElementAccessExpression().ArgumentExpression)
	default:
		return false
	}
	return !isStringOrNumericLiteralLike(expr) && !isSignedNumericLiteral(expr)
}

func getNameOfDeclaration(declaration *Node) *Node {
	if declaration == nil {
		return nil
	}
	nonAssignedName := getNonAssignedNameOfDeclaration(declaration)
	if nonAssignedName != nil {
		return nonAssignedName
	}
	if IsFunctionExpression(declaration) || IsArrowFunction(declaration) || IsClassExpression(declaration) {
		return getAssignedName(declaration)
	}
	return nil
}

func getNonAssignedNameOfDeclaration(declaration *Node) *Node {
	switch declaration.Kind {
	case ast.KindBinaryExpression:
		if isFunctionPropertyAssignment(declaration) {
			return getElementOrPropertyAccessArgumentExpressionOrName(declaration.AsBinaryExpression().Left)
		}
		return nil
	case ast.KindExportAssignment:
		expr := declaration.AsExportAssignment().Expression
		if IsIdentifier(expr) {
			return expr
		}
		return nil
	}
	return declaration.Name()
}

func getAssignedName(node *Node) *Node {
	parent := node.Parent
	if parent != nil {
		switch parent.Kind {
		case ast.KindPropertyAssignment:
			return parent.AsPropertyAssignment().Name_
		case ast.KindBindingElement:
			return parent.AsBindingElement().Name_
		case ast.KindBinaryExpression:
			if node == parent.AsBinaryExpression().Right {
				left := parent.AsBinaryExpression().Left
				switch left.Kind {
				case ast.KindIdentifier:
					return left
				case ast.KindPropertyAccessExpression:
					return left.AsPropertyAccessExpression().Name_
				case ast.KindElementAccessExpression:
					arg := skipParentheses(left.AsElementAccessExpression().ArgumentExpression)
					if isStringOrNumericLiteralLike(arg) {
						return arg
					}
				}
			}
		case ast.KindVariableDeclaration:
			name := parent.AsVariableDeclaration().Name_
			if IsIdentifier(name) {
				return name
			}
		}
	}
	return nil
}

func isFunctionPropertyAssignment(node *Node) bool {
	if node.Kind == ast.KindBinaryExpression {
		expr := node.AsBinaryExpression()
		if expr.OperatorToken.Kind == ast.KindEqualsToken {
			switch expr.Left.Kind {
			case ast.KindPropertyAccessExpression:
				// F.id = expr
				return IsIdentifier(expr.Left.AsPropertyAccessExpression().Expression) && IsIdentifier(expr.Left.AsPropertyAccessExpression().Name_)
			case ast.KindElementAccessExpression:
				// F[xxx] = expr
				return IsIdentifier(expr.Left.AsElementAccessExpression().Expression)
			}
		}
	}
	return false
}

func isAssignmentExpression(node *Node, excludeCompoundAssignment bool) bool {
	if node.Kind == ast.KindBinaryExpression {
		expr := node.AsBinaryExpression()
		return (expr.OperatorToken.Kind == ast.KindEqualsToken || !excludeCompoundAssignment && isAssignmentOperator(expr.OperatorToken.Kind)) &&
			isLeftHandSideExpression(expr.Left)
	}
	return false
}

func isBlockOrCatchScoped(declaration *Node) bool {
	return getCombinedNodeFlags(declaration)&ast.NodeFlagsBlockScoped != 0 || isCatchClauseVariableDeclarationOrBindingElement(declaration)
}

func isCatchClauseVariableDeclarationOrBindingElement(declaration *Node) bool {
	node := getRootDeclaration(declaration)
	return node.Kind == ast.KindVariableDeclaration && node.Parent.Kind == ast.KindCatchClause
}

func isAmbientModule(node *Node) bool {
	return IsModuleDeclaration(node) && (node.AsModuleDeclaration().Name_.Kind == ast.KindStringLiteral || isGlobalScopeAugmentation(node))
}

func isGlobalScopeAugmentation(node *Node) bool {
	return node.Flags&ast.NodeFlagsGlobalAugmentation != 0
}

func isPropertyNameLiteral(node *Node) bool {
	switch node.Kind {
	case ast.KindIdentifier, ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral, ast.KindNumericLiteral:
		return true
	}
	return false
}

func isMemberName(node *Node) bool {
	return node.Kind == ast.KindIdentifier || node.Kind == ast.KindPrivateIdentifier
}

func setParent(child *Node, parent *Node) {
	if child != nil {
		child.Parent = parent
	}
}

func setParentInChildren(node *Node) {
	node.ForEachChild(func(child *Node) bool {
		child.Parent = node
		setParentInChildren(child)
		return false
	})
}

func getCombinedFlags[T ~uint32](node *Node, getFlags func(*Node) T) T {
	node = getRootDeclaration(node)
	flags := getFlags(node)
	if node.Kind == ast.KindVariableDeclaration {
		node = node.Parent
	}
	if node != nil && node.Kind == ast.KindVariableDeclarationList {
		flags |= getFlags(node)
		node = node.Parent
	}
	if node != nil && node.Kind == ast.KindVariableStatement {
		flags |= getFlags(node)
	}
	return flags
}

func getCombinedModifierFlags(node *Node) ast.ModifierFlags {
	return getCombinedFlags(node, getModifierFlags)
}

func getCombinedNodeFlags(node *Node) ast.NodeFlags {
	return getCombinedFlags(node, getNodeFlags)
}

func isBindingPattern(node *Node) bool {
	return node != nil && (node.Kind == ast.KindArrayBindingPattern || node.Kind == ast.KindObjectBindingPattern)
}

func isParameterPropertyDeclaration(node *Node, parent *Node) bool {
	return IsParameter(node) && hasSyntacticModifier(node, ast.ModifierFlagsParameterPropertyModifier) && parent.Kind == ast.KindConstructor
}

/**
 * Like {@link isVariableDeclarationInitializedToRequire} but allows things like `require("...").foo.bar` or `require("...")["baz"]`.
 */
func isVariableDeclarationInitializedToBareOrAccessedRequire(node *Node) bool {
	return isVariableDeclarationInitializedWithRequireHelper(node, true /*allowAccessedRequire*/)
}

func isVariableDeclarationInitializedWithRequireHelper(node *Node, allowAccessedRequire bool) bool {
	if node.Kind == ast.KindVariableDeclaration && node.AsVariableDeclaration().Initializer != nil {
		initializer := node.AsVariableDeclaration().Initializer
		if allowAccessedRequire {
			initializer = getLeftmostAccessExpression(initializer)
		}
		return isRequireCall(initializer, true /*requireStringLiteralLikeArgument*/)
	}
	return false
}

func getLeftmostAccessExpression(expr *Node) *Node {
	for isAccessExpression(expr) {
		expr = expr.Expression()
	}
	return expr
}

func isRequireCall(node *Node, requireStringLiteralLikeArgument bool) bool {
	if IsCallExpression(node) {
		callExpression := node.AsCallExpression()
		if len(callExpression.Arguments) == 1 {
			if IsIdentifier(callExpression.Expression) && callExpression.Expression.AsIdentifier().Text == "require" {
				return !requireStringLiteralLikeArgument || isStringLiteralLike(callExpression.Arguments[0])
			}
		}
	}
	return false
}

/**
 * This function returns true if the this node's root declaration is a parameter.
 * For example, passing a `ParameterDeclaration` will return true, as will passing a
 * binding element that is a child of a `ParameterDeclaration`.
 *
 * If you are looking to test that a `Node` is a `ParameterDeclaration`, use `isParameter`.
 */
func isPartOfParameterDeclaration(node *Node) bool {
	return getRootDeclaration(node).Kind == ast.KindParameter
}

func getRootDeclaration(node *Node) *Node {
	for node.Kind == ast.KindBindingElement {
		node = node.Parent.Parent
	}
	return node
}

func isExternalOrCommonJsModule(file *SourceFile) bool {
	return file.ExternalModuleIndicator != nil
}

func isAutoAccessorPropertyDeclaration(node *Node) bool {
	return IsPropertyDeclaration(node) && hasAccessorModifier(node)
}

func isAsyncFunction(node *Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindMethodDeclaration:
		data := node.BodyData()
		return data.Body != nil && data.AsteriskToken == nil && hasSyntacticModifier(node, ast.ModifierFlagsAsync)
	}
	return false
}

func isObjectLiteralMethod(node *Node) bool {
	return node != nil && node.Kind == ast.KindMethodDeclaration && node.Parent.Kind == ast.KindObjectLiteralExpression
}

func symbolName(symbol *Symbol) string {
	if symbol.ValueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(symbol.ValueDeclaration) {
		return symbol.ValueDeclaration.Name().AsPrivateIdentifier().Text
	}
	return symbol.Name
}

func isStaticPrivateIdentifierProperty(s *Symbol) bool {
	return s.ValueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(s.ValueDeclaration) && isStatic(s.ValueDeclaration)
}

func isPrivateIdentifierClassElementDeclaration(node *Node) bool {
	return (IsPropertyDeclaration(node) || isMethodOrAccessor(node)) && IsPrivateIdentifier(node.Name())
}

func isMethodOrAccessor(node *Node) bool {
	switch node.Kind {
	case ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor:
		return true
	}
	return false
}

func isFunctionLikeOrClassStaticBlockDeclaration(node *Node) bool {
	return node != nil && (isFunctionLikeKind(node.Kind) || IsClassStaticBlockDeclaration(node))
}

func isModuleAugmentationExternal(node *Node) bool {
	// external module augmentation is a ambient module declaration that is either:
	// - defined in the top level scope and source file is an external module
	// - defined inside ambient module declaration located in the top level scope and source file not an external module
	switch node.Parent.Kind {
	case ast.KindSourceFile:
		return isExternalModule(node.Parent.AsSourceFile())
	case ast.KindModuleBlock:
		grandParent := node.Parent.Parent
		return isAmbientModule(grandParent) && IsSourceFile(grandParent.Parent) && !isExternalModule(grandParent.Parent.AsSourceFile())
	}
	return false
}

func isValidPattern(pattern Pattern) bool {
	return pattern.StarIndex == -1 || pattern.StarIndex < len(pattern.Text)
}

func tryParsePattern(pattern string) Pattern {
	starIndex := strings.Index(pattern, "*")
	if starIndex == -1 || !strings.Contains(pattern[starIndex+1:], "*") {
		return Pattern{Text: pattern, StarIndex: starIndex}
	}
	return Pattern{}
}

func positionIsSynthesized(pos int) bool {
	return pos < 0
}
func isDeclarationStatementKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindFunctionDeclaration, ast.KindMissingDeclaration, ast.KindClassDeclaration, ast.KindInterfaceDeclaration,
		ast.KindTypeAliasDeclaration, ast.KindEnumDeclaration, ast.KindModuleDeclaration, ast.KindImportDeclaration,
		ast.KindImportEqualsDeclaration, ast.KindExportDeclaration, ast.KindExportAssignment, ast.KindNamespaceExportDeclaration:
		return true
	}
	return false
}

func isDeclarationStatement(node *Node) bool {
	return isDeclarationStatementKind(node.Kind)
}

func isStatementKindButNotDeclarationKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindBreakStatement, ast.KindContinueStatement, ast.KindDebuggerStatement, ast.KindDoStatement, ast.KindExpressionStatement,
		ast.KindEmptyStatement, ast.KindForInStatement, ast.KindForOfStatement, ast.KindForStatement, ast.KindIfStatement,
		ast.KindLabeledStatement, ast.KindReturnStatement, ast.KindSwitchStatement, ast.KindThrowStatement, ast.KindTryStatement,
		ast.KindVariableStatement, ast.KindWhileStatement, ast.KindWithStatement, ast.KindNotEmittedStatement:
		return true
	}
	return false
}

func isStatementButNotDeclaration(node *Node) bool {
	return isStatementKindButNotDeclarationKind(node.Kind)
}

func isStatement(node *Node) bool {
	kind := node.Kind
	return isStatementKindButNotDeclarationKind(kind) || isDeclarationStatementKind(kind) || isBlockStatement(node)
}

func isBlockStatement(node *Node) bool {
	if node.Kind != ast.KindBlock {
		return false
	}
	if node.Parent != nil && (node.Parent.Kind == ast.KindTryStatement || node.Parent.Kind == ast.KindCatchClause) {
		return false
	}
	return !isFunctionBlock(node)
}

func isFunctionBlock(node *Node) bool {
	return node != nil && node.Kind == ast.KindBlock && isFunctionLike(node.Parent)
}

func shouldPreserveConstEnums(options *CompilerOptions) bool {
	return options.PreserveConstEnums == core.TSTrue || options.IsolatedModules == core.TSTrue
}

func exportAssignmentIsAlias(node *Node) bool {
	return isAliasableExpression(getExportAssignmentExpression(node))
}

func getExportAssignmentExpression(node *Node) *Node {
	switch node.Kind {
	case ast.KindExportAssignment:
		return node.AsExportAssignment().Expression
	case ast.KindBinaryExpression:
		return node.AsBinaryExpression().Right
	}
	panic("Unhandled case in getExportAssignmentExpression")
}

func isAliasableExpression(e *Node) bool {
	return isEntityNameExpression(e) || IsClassExpression(e)
}

func isEmptyObjectLiteral(expression *Node) bool {
	return expression.Kind == ast.KindObjectLiteralExpression && len(expression.AsObjectLiteralExpression().Properties) == 0
}

func isFunctionSymbol(symbol *Symbol) bool {
	d := symbol.ValueDeclaration
	return d != nil && (IsFunctionDeclaration(d) || IsVariableDeclaration(d) && isFunctionLike(d.AsVariableDeclaration().Initializer))
}

func isLogicalOrCoalescingAssignmentOperator(token ast.Kind) bool {
	return token == ast.KindBarBarEqualsToken || token == ast.KindAmpersandAmpersandEqualsToken || token == ast.KindQuestionQuestionEqualsToken
}

func isLogicalOrCoalescingAssignmentExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingAssignmentOperator(expr.AsBinaryExpression().OperatorToken.Kind)
}

func isLogicalOrCoalescingBinaryOperator(token ast.Kind) bool {
	return isBinaryLogicalOperator(token) || token == ast.KindQuestionQuestionToken
}

func isLogicalOrCoalescingBinaryExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingBinaryOperator(expr.AsBinaryExpression().OperatorToken.Kind)
}

func isBinaryLogicalOperator(token ast.Kind) bool {
	return token == ast.KindBarBarToken || token == ast.KindAmpersandAmpersandToken
}

/**
 * Determines whether a node is the outermost `OptionalChain` in an ECMAScript `OptionalExpression`:
 *
 * 1. For `a?.b.c`, the outermost chain is `a?.b.c` (`c` is the end of the chain starting at `a?.`)
 * 2. For `a?.b!`, the outermost chain is `a?.b` (`b` is the end of the chain starting at `a?.`)
 * 3. For `(a?.b.c).d`, the outermost chain is `a?.b.c` (`c` is the end of the chain starting at `a?.` since parens end the chain)
 * 4. For `a?.b.c?.d`, both `a?.b.c` and `a?.b.c?.d` are outermost (`c` is the end of the chain starting at `a?.`, and `d` is
 *   the end of the chain starting at `c?.`)
 * 5. For `a?.(b?.c).d`, both `b?.c` and `a?.(b?.c)d` are outermost (`c` is the end of the chain starting at `b`, and `d` is
 *   the end of the chain starting at `a?.`)
 */
func isOutermostOptionalChain(node *Node) bool {
	parent := node.Parent
	return !isOptionalChain(parent) || // cases 1, 2, and 3
		isOptionalChainRoot(parent) || // case 4
		node != parent.Expression() // case 5
}

func isNullishCoalesce(node *Node) bool {
	return node.Kind == ast.KindBinaryExpression && node.AsBinaryExpression().OperatorToken.Kind == ast.KindQuestionQuestionToken
}

func isDottedName(node *Node) bool {
	switch node.Kind {
	case ast.KindIdentifier, ast.KindThisKeyword, ast.KindSuperKeyword, ast.KindMetaProperty:
		return true
	case ast.KindPropertyAccessExpression, ast.KindParenthesizedExpression:
		return isDottedName(node.Expression())
	}
	return false
}

func unusedLabelIsError(options *CompilerOptions) bool {
	return options.AllowUnusedLabels == core.TSFalse
}

func unreachableCodeIsError(options *CompilerOptions) bool {
	return options.AllowUnreachableCode == core.TSFalse
}

func isDestructuringAssignment(node *Node) bool {
	if isAssignmentExpression(node, true /*excludeCompoundAssignment*/) {
		kind := node.AsBinaryExpression().Left.Kind
		return kind == ast.KindObjectLiteralExpression || kind == ast.KindArrayLiteralExpression
	}
	return false
}

func isTopLevelLogicalExpression(node *Node) bool {
	for IsParenthesizedExpression(node.Parent) || IsPrefixUnaryExpression(node.Parent) && node.Parent.AsPrefixUnaryExpression().Operator == ast.KindExclamationToken {
		node = node.Parent
	}
	return !isStatementCondition(node) && !isLogicalExpression(node.Parent) && !(isOptionalChain(node.Parent) && node.Parent.Expression() == node)
}

func isStatementCondition(node *Node) bool {
	switch node.Parent.Kind {
	case ast.KindIfStatement:
		return node.Parent.AsIfStatement().Expression == node
	case ast.KindWhileStatement:
		return node.Parent.AsWhileStatement().Expression == node
	case ast.KindDoStatement:
		return node.Parent.AsDoStatement().Expression == node
	case ast.KindForStatement:
		return node.Parent.AsForStatement().Condition == node
	case ast.KindConditionalExpression:
		return node.Parent.AsConditionalExpression().Condition == node
	}
	return false
}

type AssignmentKind int32

const (
	AssignmentKindNone AssignmentKind = iota
	AssignmentKindDefinite
	AssignmentKindCompound
)

type AssignmentTarget = Node // BinaryExpression | PrefixUnaryExpression | PostfixUnaryExpression | ForInOrOfStatement

func getAssignmentTargetKind(node *Node) AssignmentKind {
	target := getAssignmentTarget(node)
	if target == nil {
		return AssignmentKindNone
	}
	switch target.Kind {
	case ast.KindBinaryExpression:
		binaryOperator := target.AsBinaryExpression().OperatorToken.Kind
		if binaryOperator == ast.KindEqualsToken || isLogicalOrCoalescingAssignmentOperator(binaryOperator) {
			return AssignmentKindDefinite
		}
		return AssignmentKindCompound
	case ast.KindPrefixUnaryExpression, ast.KindPostfixUnaryExpression:
		return AssignmentKindCompound
	case ast.KindForInStatement, ast.KindForOfStatement:
		return AssignmentKindDefinite
	}
	panic("Unhandled case in getAssignmentTargetKind")
}

// A node is an assignment target if it is on the left hand side of an '=' token, if it is parented by a property
// assignment in an object literal that is an assignment target, or if it is parented by an array literal that is
// an assignment target. Examples include 'a = xxx', '{ p: a } = xxx', '[{ a }] = xxx'.
// (Note that `p` is not a target in the above examples, only `a`.)
func isAssignmentTarget(node *Node) bool {
	return getAssignmentTarget(node) != nil
}

// Returns the BinaryExpression, PrefixUnaryExpression, PostfixUnaryExpression, or ForInOrOfStatement that references
// the given node as an assignment target
func getAssignmentTarget(node *Node) *Node {
	for {
		parent := node.Parent
		switch parent.Kind {
		case ast.KindBinaryExpression:
			if isAssignmentOperator(parent.AsBinaryExpression().OperatorToken.Kind) && parent.AsBinaryExpression().Left == node {
				return parent
			}
			return nil
		case ast.KindPrefixUnaryExpression:
			if parent.AsPrefixUnaryExpression().Operator == ast.KindPlusPlusToken || parent.AsPrefixUnaryExpression().Operator == ast.KindMinusMinusToken {
				return parent
			}
			return nil
		case ast.KindPostfixUnaryExpression:
			if parent.AsPostfixUnaryExpression().Operator == ast.KindPlusPlusToken || parent.AsPostfixUnaryExpression().Operator == ast.KindMinusMinusToken {
				return parent
			}
			return nil
		case ast.KindForInStatement, ast.KindForOfStatement:
			if parent.AsForInOrOfStatement().Initializer == node {
				return parent
			}
			return nil
		case ast.KindParenthesizedExpression, ast.KindArrayLiteralExpression, ast.KindSpreadElement, ast.KindNonNullExpression:
			node = parent
		case ast.KindSpreadAssignment:
			node = parent.Parent
		case ast.KindShorthandPropertyAssignment:
			if parent.AsShorthandPropertyAssignment().Name_ != node {
				return nil
			}
			node = parent.Parent
		case ast.KindPropertyAssignment:
			if parent.AsPropertyAssignment().Name_ == node {
				return nil
			}
			node = parent.Parent
		default:
			return nil
		}
	}
}

func isDeleteTarget(node *Node) bool {
	if !isAccessExpression(node) {
		return false
	}
	node = walkUpParenthesizedExpressions(node.Parent)
	return node != nil && node.Kind == ast.KindDeleteExpression
}

func isInCompoundLikeAssignment(node *Node) bool {
	target := getAssignmentTarget(node)
	return target != nil && isAssignmentExpression(target /*excludeCompoundAssignment*/, true) && isCompoundLikeAssignment(target)
}

func isCompoundLikeAssignment(assignment *Node) bool {
	right := skipParentheses(assignment.AsBinaryExpression().Right)
	return right.Kind == ast.KindBinaryExpression && isShiftOperatorOrHigher(right.AsBinaryExpression().OperatorToken.Kind)
}

func isPushOrUnshiftIdentifier(node *Node) bool {
	text := node.AsIdentifier().Text
	return text == "push" || text == "unshift"
}

func isBooleanLiteral(node *Node) bool {
	return node.Kind == ast.KindTrueKeyword || node.Kind == ast.KindFalseKeyword
}

func isOptionalChain(node *Node) bool {
	kind := node.Kind
	return node.Flags&ast.NodeFlagsOptionalChain != 0 && (kind == ast.KindPropertyAccessExpression ||
		kind == ast.KindElementAccessExpression || kind == ast.KindCallExpression || kind == ast.KindNonNullExpression)
}

func isOptionalChainRoot(node *Node) bool {
	return isOptionalChain(node) && !isNonNullExpression(node) && getQuestionDotToken(node) != nil
}

/**
 * Determines whether a node is the expression preceding an optional chain (i.e. `a` in `a?.b`).
 */
func isExpressionOfOptionalChainRoot(node *Node) bool {
	return isOptionalChainRoot(node.Parent) && node.Parent.Expression() == node
}

func isEntityNameExpression(node *Node) bool {
	return node.Kind == ast.KindIdentifier || isPropertyAccessEntityNameExpression(node)
}

func isPropertyAccessEntityNameExpression(node *Node) bool {
	if node.Kind == ast.KindPropertyAccessExpression {
		expr := node.AsPropertyAccessExpression()
		return expr.Name_.Kind == ast.KindIdentifier && isEntityNameExpression(expr.Expression)
	}
	return false
}

func isPrologueDirective(node *Node) bool {
	return node.Kind == ast.KindExpressionStatement && node.AsExpressionStatement().Expression.Kind == ast.KindStringLiteral
}

func getStatementsOfBlock(block *Node) []*Statement {
	switch block.Kind {
	case ast.KindBlock:
		return block.AsBlock().Statements
	case ast.KindModuleBlock:
		return block.AsModuleBlock().Statements
	case ast.KindSourceFile:
		return block.AsSourceFile().Statements
	}
	panic("Unhandled case in getStatementsOfBlock")
}

func nodeHasName(statement *Node, id *Node) bool {
	name := statement.Name()
	if name != nil {
		return IsIdentifier(name) && name.AsIdentifier().Text == id.AsIdentifier().Text
	}
	if IsVariableStatement(statement) {
		declarations := statement.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations
		return core.Some(declarations, func(d *Node) bool { return nodeHasName(d, id) })
	}
	return false
}

func isImportMeta(node *Node) bool {
	if node.Kind == ast.KindMetaProperty {
		return node.AsMetaProperty().KeywordToken == ast.KindImportKeyword && node.AsMetaProperty().Name_.AsIdentifier().Text == "meta"
	}
	return false
}

func ensureScriptKind(fileName string, scriptKind core.ScriptKind) core.ScriptKind {
	// Using scriptKind as a condition handles both:
	// - 'scriptKind' is unspecified and thus it is `undefined`
	// - 'scriptKind' is set and it is `Unknown` (0)
	// If the 'scriptKind' is 'undefined' or 'Unknown' then we attempt
	// to get the core.ScriptKind from the file name. If it cannot be resolved
	// from the file name then the default 'TS' script kind is returned.
	if scriptKind == core.ScriptKindUnknown {
		scriptKind = getScriptKindFromFileName(fileName)
	}
	if scriptKind == core.ScriptKindUnknown {
		scriptKind = core.ScriptKindTS
	}
	return scriptKind
}

const (
	ExtensionTs          = ".ts"
	ExtensionTsx         = ".tsx"
	ExtensionDts         = ".d.ts"
	ExtensionJs          = ".js"
	ExtensionJsx         = ".jsx"
	ExtensionJson        = ".json"
	ExtensionTsBuildInfo = ".tsbuildinfo"
	ExtensionMjs         = ".mjs"
	ExtensionMts         = ".mts"
	ExtensionDmts        = ".d.mts"
	ExtensionCjs         = ".cjs"
	ExtensionCts         = ".cts"
	ExtensionDcts        = ".d.cts"
)

var supportedDeclarationExtensions = []string{ExtensionDts, ExtensionDcts, ExtensionDmts}

func getScriptKindFromFileName(fileName string) core.ScriptKind {
	dotPos := strings.LastIndex(fileName, ".")
	if dotPos >= 0 {
		switch strings.ToLower(fileName[dotPos:]) {
		case ExtensionJs, ExtensionCjs, ExtensionMjs:
			return core.ScriptKindJS
		case ExtensionJsx:
			return core.ScriptKindJSX
		case ExtensionTs, ExtensionCts, ExtensionMts:
			return core.ScriptKindTS
		case ExtensionTsx:
			return core.ScriptKindTSX
		case ExtensionJson:
			return core.ScriptKindJSON
		}
	}
	return core.ScriptKindUnknown
}

func getLanguageVariant(scriptKind core.ScriptKind) core.LanguageVariant {
	switch scriptKind {
	case core.ScriptKindTSX, core.ScriptKindJSX, core.ScriptKindJS, core.ScriptKindJSON:
		// .tsx and .jsx files are treated as jsx language variant.
		return core.LanguageVariantJSX
	}
	return core.LanguageVariantStandard
}

func getEmitScriptTarget(options *CompilerOptions) core.ScriptTarget {
	if options.Target != core.ScriptTargetNone {
		return options.Target
	}
	return core.ScriptTargetES5
}

func getEmitModuleKind(options *CompilerOptions) ModuleKind {
	if options.ModuleKind != ModuleKindNone {
		return options.ModuleKind
	}
	if options.Target >= core.ScriptTargetES2015 {
		return ModuleKindES2015
	}
	return ModuleKindCommonJS
}

func getEmitModuleResolutionKind(options *CompilerOptions) ModuleResolutionKind {
	if options.ModuleResolution != ModuleResolutionKindUnknown {
		return options.ModuleResolution
	}
	switch getEmitModuleKind(options) {
	case ModuleKindCommonJS:
		return ModuleResolutionKindBundler
	case ModuleKindNode16:
		return ModuleResolutionKindNode16
	case ModuleKindNodeNext:
		return ModuleResolutionKindNodeNext
	case ModuleKindPreserve:
		return ModuleResolutionKindBundler
	default:
		panic("Unhandled case in getEmitModuleResolutionKind")
	}
}

func getESModuleInterop(options *CompilerOptions) bool {
	if options.ESModuleInterop != core.TSUnknown {
		return options.ESModuleInterop == core.TSTrue
	}
	switch getEmitModuleKind(options) {
	case ModuleKindNode16:
	case ModuleKindNodeNext:
	case ModuleKindPreserve:
		return true
	}
	return false

}
func getAllowSyntheticDefaultImports(options *CompilerOptions) bool {
	if options.AllowSyntheticDefaultImports != core.TSUnknown {
		return options.AllowSyntheticDefaultImports == core.TSTrue
	}
	return getESModuleInterop(options) ||
		getEmitModuleKind(options) == ModuleKindSystem ||
		getEmitModuleResolutionKind(options) == ModuleResolutionKindBundler
}

type DiagnosticsCollection struct {
	fileDiagnostics    map[string][]*Diagnostic
	nonFileDiagnostics []*Diagnostic
}

func (c *DiagnosticsCollection) add(diagnostic *Diagnostic) {
	if diagnostic.File_ != nil {
		fileName := diagnostic.File_.FileName_
		if c.fileDiagnostics == nil {
			c.fileDiagnostics = make(map[string][]*Diagnostic)
		}
		c.fileDiagnostics[fileName] = core.InsertSorted(c.fileDiagnostics[fileName], diagnostic, compareDiagnostics)
	} else {
		c.nonFileDiagnostics = core.InsertSorted(c.nonFileDiagnostics, diagnostic, compareDiagnostics)
	}
}

func (c *DiagnosticsCollection) lookup(diagnostic *Diagnostic) *Diagnostic {
	var diagnostics []*Diagnostic
	if diagnostic.File_ != nil {
		diagnostics = c.fileDiagnostics[diagnostic.File_.FileName_]
	} else {
		diagnostics = c.nonFileDiagnostics
	}
	if i, ok := slices.BinarySearchFunc(diagnostics, diagnostic, compareDiagnostics); ok {
		return diagnostics[i]
	}
	return nil
}

func (c *DiagnosticsCollection) GetGlobalDiagnostics() []*Diagnostic {
	return c.nonFileDiagnostics
}

func (c *DiagnosticsCollection) GetDiagnosticsForFile(fileName string) []*Diagnostic {
	return c.fileDiagnostics[fileName]
}

func (c *DiagnosticsCollection) GetDiagnostics() []*Diagnostic {
	fileNames := slices.Collect(maps.Keys(c.fileDiagnostics))
	slices.Sort(fileNames)
	diagnostics := c.nonFileDiagnostics
	for _, fileName := range fileNames {
		diagnostics = append(diagnostics, c.fileDiagnostics[fileName]...)
	}
	return diagnostics
}

func sortAndDeduplicateDiagnostics(diagnostics []*Diagnostic) []*Diagnostic {
	result := slices.Clone(diagnostics)
	slices.SortFunc(result, compareDiagnostics)
	return slices.CompactFunc(result, equalDiagnostics)
}

func equalDiagnostics(d1, d2 *Diagnostic) bool {
	return getDiagnosticPath(d1) == getDiagnosticPath(d2) &&
		d1.Loc_ == d2.Loc_ &&
		d1.Code_ == d2.Code_ &&
		d1.Message_ == d2.Message_ &&
		slices.EqualFunc(d1.MessageChain_, d2.MessageChain_, equalMessageChain) &&
		slices.EqualFunc(d1.RelatedInformation_, d2.RelatedInformation_, equalDiagnostics)
}

func equalMessageChain(c1, c2 *MessageChain) bool {
	return c1.Code_ == c2.Code_ &&
		c1.Message_ == c2.Message_ &&
		slices.EqualFunc(c1.MessageChain_, c2.MessageChain_, equalMessageChain)
}

func compareDiagnostics(d1, d2 *Diagnostic) int {
	c := strings.Compare(getDiagnosticPath(d1), getDiagnosticPath(d2))
	if c != 0 {
		return c
	}
	c = int(d1.Loc_.Pos_) - int(d2.Loc_.Pos_)
	if c != 0 {
		return c
	}
	c = int(d1.Loc_.End_) - int(d2.Loc_.End_)
	if c != 0 {
		return c
	}
	c = int(d1.Code_) - int(d2.Code_)
	if c != 0 {
		return c
	}
	c = strings.Compare(d1.Message_, d2.Message_)
	if c != 0 {
		return c
	}
	c = compareMessageChainSize(d1.MessageChain_, d2.MessageChain_)
	if c != 0 {
		return c
	}
	c = compareMessageChainContent(d1.MessageChain_, d2.MessageChain_)
	if c != 0 {
		return c
	}
	return compareRelatedInfo(d1.RelatedInformation_, d2.RelatedInformation_)
}

func compareMessageChainSize(c1, c2 []*MessageChain) int {
	c := len(c2) - len(c1)
	if c != 0 {
		return c
	}
	for i := range c1 {
		c = compareMessageChainSize(c1[i].MessageChain_, c2[i].MessageChain_)
		if c != 0 {
			return c
		}
	}
	return 0
}

func compareMessageChainContent(c1, c2 []*MessageChain) int {
	for i := range c1 {
		c := strings.Compare(c1[i].Message_, c2[i].Message_)
		if c != 0 {
			return c
		}
		if c1[i].MessageChain_ != nil {
			c = compareMessageChainContent(c1[i].MessageChain_, c2[i].MessageChain_)
			if c != 0 {
				return c
			}
		}
	}
	return 0
}

func compareRelatedInfo(r1, r2 []*Diagnostic) int {
	c := len(r2) - len(r1)
	if c != 0 {
		return c
	}
	for i := range r1 {
		c = compareDiagnostics(r1[i], r2[i])
		if c != 0 {
			return c
		}
	}
	return 0
}

func getDiagnosticPath(d *Diagnostic) string {
	if d.File_ != nil {
		return d.File_.Path_
	}
	return ""
}

func isConstAssertion(location *Node) bool {
	switch location.Kind {
	case ast.KindAsExpression:
		return isConstTypeReference(location.AsAsExpression().TypeNode)
	case ast.KindTypeAssertionExpression:
		return isConstTypeReference(location.AsTypeAssertion().TypeNode)
	}
	return false
}

func isConstTypeReference(node *Node) bool {
	if node.Kind == ast.KindTypeReference {
		ref := node.AsTypeReference()
		return ref.TypeArguments == nil && IsIdentifier(ref.TypeName) && ref.TypeName.AsIdentifier().Text == "const"
	}
	return false
}

func isModuleOrEnumDeclaration(node *Node) bool {
	return node.Kind == ast.KindModuleDeclaration || node.Kind == ast.KindEnumDeclaration
}

func getLocalsOfNode(node *Node) SymbolTable {
	data := node.LocalsContainerData()
	if data != nil {
		return data.Locals_
	}
	return nil
}

func getBodyOfNode(node *Node) *Node {
	bodyData := node.BodyData()
	if bodyData != nil {
		return bodyData.Body
	}
	return nil
}

func getFlowNodeOfNode(node *Node) *FlowNode {
	flowNodeData := node.FlowNodeData()
	if flowNodeData != nil {
		return flowNodeData.FlowNode
	}
	return nil
}

func isGlobalSourceFile(node *Node) bool {
	return node.Kind == ast.KindSourceFile && !isExternalOrCommonJsModule(node.AsSourceFile())
}

func isParameterLikeOrReturnTag(node *Node) bool {
	switch node.Kind {
	case ast.KindParameter, ast.KindTypeParameter, ast.KindJSDocParameterTag, ast.KindJSDocReturnTag:
		return true
	}
	return false
}

func getEmitStandardClassFields(options *CompilerOptions) bool {
	return options.UseDefineForClassFields != core.TSFalse && getEmitScriptTarget(options) >= core.ScriptTargetES2022
}

func isTypeNodeKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindAnyKeyword, ast.KindUnknownKeyword, ast.KindNumberKeyword, ast.KindBigIntKeyword, ast.KindObjectKeyword,
		ast.KindBooleanKeyword, ast.KindStringKeyword, ast.KindSymbolKeyword, ast.KindVoidKeyword, ast.KindUndefinedKeyword,
		ast.KindNeverKeyword, ast.KindIntrinsicKeyword, ast.KindExpressionWithTypeArguments, ast.KindJSDocAllType, ast.KindJSDocUnknownType,
		ast.KindJSDocNullableType, ast.KindJSDocNonNullableType, ast.KindJSDocOptionalType, ast.KindJSDocFunctionType, ast.KindJSDocVariadicType:
		return true
	}
	return kind >= ast.KindFirstTypeNode && kind <= ast.KindLastTypeNode
}

func isTypeNode(node *Node) bool {
	return isTypeNodeKind(node.Kind)
}

func getLocalSymbolForExportDefault(symbol *Symbol) *Symbol {
	if !isExportDefaultSymbol(symbol) || len(symbol.Declarations) == 0 {
		return nil
	}
	for _, decl := range symbol.Declarations {
		localSymbol := decl.LocalSymbol()
		if localSymbol != nil {
			return localSymbol
		}
	}
	return nil
}

func isExportDefaultSymbol(symbol *Symbol) bool {
	return symbol != nil && len(symbol.Declarations) > 0 && hasSyntacticModifier(symbol.Declarations[0], ast.ModifierFlagsDefault)
}

func getDeclarationOfKind(symbol *Symbol, kind ast.Kind) *Node {
	for _, declaration := range symbol.Declarations {
		if declaration.Kind == kind {
			return declaration
		}
	}
	return nil
}

func getIsolatedModules(options *CompilerOptions) bool {
	return options.IsolatedModules == core.TSTrue || options.VerbatimModuleSyntax == core.TSTrue
}

func findConstructorDeclaration(node *Node) *Node {
	for _, member := range node.ClassLikeData().Members {
		if IsConstructorDeclaration(member) && nodeIsPresent(member.AsConstructorDeclaration().Body) {
			return member
		}
	}
	return nil
}

type NameResolver struct {
	compilerOptions                  *CompilerOptions
	getSymbolOfDeclaration           func(node *Node) *Symbol
	error                            func(location *Node, message *diagnostics.Message, args ...any) *Diagnostic
	globals                          SymbolTable
	argumentsSymbol                  *Symbol
	requireSymbol                    *Symbol
	lookup                           func(symbols SymbolTable, name string, meaning ast.SymbolFlags) *Symbol
	setRequiresScopeChangeCache      func(node *Node, value core.Tristate)
	getRequiresScopeChangeCache      func(node *Node) core.Tristate
	onPropertyWithInvalidInitializer func(location *Node, name string, declaration *Node, result *Symbol) bool
	onFailedToResolveSymbol          func(location *Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message)
	onSuccessfullyResolvedSymbol     func(location *Node, result *Symbol, meaning ast.SymbolFlags, lastLocation *Node, associatedDeclarationForContainingInitializerOrBindingName *Node, withinDeferredContext bool)
}

func (r *NameResolver) resolve(location *Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *Symbol {
	var result *Symbol
	var lastLocation *Node
	var lastSelfReferenceLocation *Node
	var propertyWithInvalidInitializer *Node
	var associatedDeclarationForContainingInitializerOrBindingName *Node
	var withinDeferredContext bool
	var grandparent *Node
	originalLocation := location // needed for did-you-mean error reporting, which gathers candidates starting from the original location
	nameIsConst := name == "const"
loop:
	for location != nil {
		if nameIsConst && isConstAssertion(location) {
			// `const` in an `as const` has no symbol, but issues no error because there is no *actual* lookup of the type
			// (it refers to the constant type of the expression instead)
			return nil
		}
		if isModuleOrEnumDeclaration(location) && lastLocation != nil && location.Name() == lastLocation {
			// If lastLocation is the name of a namespace or enum, skip the parent since it will have is own locals that could
			// conflict.
			lastLocation = location
			location = location.Parent
		}
		locals := getLocalsOfNode(location)
		// Locals of a source file are not in scope (because they get merged into the global symbol table)
		if locals != nil && !isGlobalSourceFile(location) {
			result = r.lookup(locals, name, meaning)
			if result != nil {
				useResult := true
				if isFunctionLike(location) && lastLocation != nil && lastLocation != getBodyOfNode(location) {
					// symbol lookup restrictions for function-like declarations
					// - Type parameters of a function are in scope in the entire function declaration, including the parameter
					//   list and return type. However, local types are only in scope in the function body.
					// - parameters are only in the scope of function body
					// This restriction does not apply to JSDoc comment types because they are parented
					// at a higher level than type parameters would normally be
					if meaning&result.Flags&ast.SymbolFlagsType != 0 && lastLocation.Kind != ast.KindJSDoc {
						useResult = result.Flags&ast.SymbolFlagsTypeParameter != 0 && (lastLocation.Flags&ast.NodeFlagsSynthesized != 0 ||
							lastLocation == location.ReturnType() ||
							isParameterLikeOrReturnTag(lastLocation))
					}
					if meaning&result.Flags&ast.SymbolFlagsVariable != 0 {
						// expression inside parameter will lookup as normal variable scope when targeting es2015+
						if r.useOuterVariableScopeInParameter(result, location, lastLocation) {
							useResult = false
						} else if result.Flags&ast.SymbolFlagsFunctionScopedVariable != 0 {
							// parameters are visible only inside function body, parameter list and return type
							// technically for parameter list case here we might mix parameters and variables declared in function,
							// however it is detected separately when checking initializers of parameters
							// to make sure that they reference no variables declared after them.
							useResult = lastLocation.Kind == ast.KindParameter ||
								lastLocation.Flags&ast.NodeFlagsSynthesized != 0 ||
								lastLocation == location.ReturnType() && findAncestor(result.ValueDeclaration, IsParameter) != nil
						}
					}
				} else if location.Kind == ast.KindConditionalType {
					// A type parameter declared using 'infer T' in a conditional type is visible only in
					// the true branch of the conditional type.
					useResult = lastLocation == location.AsConditionalTypeNode().TrueType
				}
				if useResult {
					break loop
				}
				result = nil
			}
		}
		withinDeferredContext = withinDeferredContext || getIsDeferredContext(location, lastLocation)
		switch location.Kind {
		case ast.KindSourceFile:
			if !isExternalOrCommonJsModule(location.AsSourceFile()) {
				break
			}
			fallthrough
		case ast.KindModuleDeclaration:
			moduleExports := r.getSymbolOfDeclaration(location).Exports
			if IsSourceFile(location) || (IsModuleDeclaration(location) && location.Flags&ast.NodeFlagsAmbient != 0 && !isGlobalScopeAugmentation(location)) {
				// It's an external module. First see if the module has an export default and if the local
				// name of that export default matches.
				result = moduleExports[InternalSymbolNameDefault]
				if result != nil {
					localSymbol := getLocalSymbolForExportDefault(result)
					if localSymbol != nil && result.Flags&meaning != 0 && localSymbol.Name == name {
						break loop
					}
					result = nil
				}
				// Because of module/namespace merging, a module's exports are in scope,
				// yet we never want to treat an export specifier as putting a member in scope.
				// Therefore, if the name we find is purely an export specifier, it is not actually considered in scope.
				// Two things to note about this:
				//     1. We have to check this without calling getSymbol. The problem with calling getSymbol
				//        on an export specifier is that it might find the export specifier itself, and try to
				//        resolve it as an alias. This will cause the checker to consider the export specifier
				//        a circular alias reference when it might not be.
				//     2. We check === ast.SymbolFlags.Alias in order to check that the symbol is *purely*
				//        an alias. If we used &, we'd be throwing out symbols that have non alias aspects,
				//        which is not the desired behavior.
				moduleExport := moduleExports[name]
				if moduleExport != nil && moduleExport.Flags == ast.SymbolFlagsAlias && (getDeclarationOfKind(moduleExport, ast.KindExportSpecifier) != nil || getDeclarationOfKind(moduleExport, ast.KindNamespaceExport) != nil) {
					break
				}
			}
			if name != InternalSymbolNameDefault {
				result = r.lookup(moduleExports, name, meaning&ast.SymbolFlagsModuleMember)
				if result != nil {
					break loop
				}
			}
		case ast.KindEnumDeclaration:
			result = r.lookup(r.getSymbolOfDeclaration(location).Exports, name, meaning&ast.SymbolFlagsEnumMember)
			if result != nil {
				if nameNotFoundMessage != nil && getIsolatedModules(r.compilerOptions) && location.Flags&ast.NodeFlagsAmbient == 0 && getSourceFileOfNode(location) != getSourceFileOfNode(result.ValueDeclaration) {
					isolatedModulesLikeFlagName := ifElse(r.compilerOptions.VerbatimModuleSyntax == core.TSTrue, "verbatimModuleSyntax", "isolatedModules")
					r.error(originalLocation, diagnostics.Cannot_access_0_from_another_file_without_qualification_when_1_is_enabled_Use_2_instead,
						name, isolatedModulesLikeFlagName, r.getSymbolOfDeclaration(location).Name+"."+name)
				}
				break loop
			}
		case ast.KindPropertyDeclaration:
			if !isStatic(location) {
				ctor := findConstructorDeclaration(location.Parent)
				if ctor != nil && ctor.AsConstructorDeclaration().Locals_ != nil {
					if r.lookup(ctor.AsConstructorDeclaration().Locals_, name, meaning&ast.SymbolFlagsValue) != nil {
						// Remember the property node, it will be used later to report appropriate error
						propertyWithInvalidInitializer = location
					}
				}
			}
		case ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration:
			result = r.lookup(r.getSymbolOfDeclaration(location).Members, name, meaning&ast.SymbolFlagsType)
			if result != nil {
				if !isTypeParameterSymbolDeclaredInContainer(result, location) {
					// ignore type parameters not declared in this container
					result = nil
					break
				}
				if lastLocation != nil && isStatic(lastLocation) {
					// TypeScript 1.0 spec (April 2014): 3.4.1
					// The scope of a type parameter extends over the entire declaration with which the type
					// parameter list is associated, with the exception of static member declarations in classes.
					if nameNotFoundMessage != nil {
						r.error(originalLocation, diagnostics.Static_members_cannot_reference_class_type_parameters)
					}
					return nil
				}
				break loop
			}
			if IsClassExpression(location) && meaning&ast.SymbolFlagsClass != 0 {
				className := location.AsClassExpression().Name_
				if className != nil && name == className.AsIdentifier().Text {
					result = location.AsClassExpression().Symbol
					break loop
				}
			}
		case ast.KindExpressionWithTypeArguments:
			if lastLocation == location.AsExpressionWithTypeArguments().Expression && IsHeritageClause(location.Parent) && location.Parent.AsHeritageClause().Token == ast.KindExtendsKeyword {
				container := location.Parent.Parent
				if isClassLike(container) {
					result = r.lookup(r.getSymbolOfDeclaration(container).Members, name, meaning&ast.SymbolFlagsType)
					if result != nil {
						if nameNotFoundMessage != nil {
							r.error(originalLocation, diagnostics.Base_class_expressions_cannot_reference_class_type_parameters)
						}
						return nil
					}
				}
			}
		// It is not legal to reference a class's own type parameters from a computed property name that
		// belongs to the class. For example:
		//
		//   function foo<T>() { return '' }
		//   class C<T> { // <-- Class's own type parameter T
		//       [foo<T>()]() { } // <-- Reference to T from class's own computed property
		//   }
		case ast.KindComputedPropertyName:
			grandparent = location.Parent.Parent
			if isClassLike(grandparent) || IsInterfaceDeclaration(grandparent) {
				// A reference to this grandparent's type parameters would be an error
				result = r.lookup(r.getSymbolOfDeclaration(grandparent).Members, name, meaning&ast.SymbolFlagsType)
				if result != nil {
					if nameNotFoundMessage != nil {
						r.error(originalLocation, diagnostics.A_computed_property_name_cannot_reference_a_type_parameter_from_its_containing_type)
					}
					return nil
				}
			}
		case ast.KindArrowFunction:
			// when targeting ES6 or higher there is no 'arguments' in an arrow function
			// for lower compile targets the resolved symbol is used to emit an error
			if getEmitScriptTarget(r.compilerOptions) >= core.ScriptTargetES2015 {
				break
			}
			fallthrough
		case ast.KindMethodDeclaration, ast.KindConstructor, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindFunctionDeclaration:
			if meaning&ast.SymbolFlagsVariable != 0 && name == "arguments" {
				result = r.argumentsSymbol
				break loop
			}
		case ast.KindFunctionExpression:
			if meaning&ast.SymbolFlagsVariable != 0 && name == "arguments" {
				result = r.argumentsSymbol
				break loop
			}
			if meaning&ast.SymbolFlagsFunction != 0 {
				functionName := location.AsFunctionExpression().Name_
				if functionName != nil && name == functionName.AsIdentifier().Text {
					result = location.AsFunctionExpression().Symbol
					break loop
				}
			}
		case ast.KindDecorator:
			// Decorators are resolved at the class declaration. Resolving at the parameter
			// or member would result in looking up locals in the method.
			//
			//   function y() {}
			//   class C {
			//       method(@y x, y) {} // <-- decorator y should be resolved at the class declaration, not the parameter.
			//   }
			//
			if location.Parent != nil && location.Parent.Kind == ast.KindParameter {
				location = location.Parent
			}
			//   function y() {}
			//   class C {
			//       @y method(x, y) {} // <-- decorator y should be resolved at the class declaration, not the method.
			//   }
			//
			// class Decorators are resolved outside of the class to avoid referencing type parameters of that class.
			//
			//   type T = number;
			//   declare function y(x: T): any;
			//   @param(1 as T) // <-- T should resolve to the type alias outside of class C
			//   class C<T> {}
			if location.Parent != nil && (isClassElement(location.Parent) || location.Parent.Kind == ast.KindClassDeclaration) {
				location = location.Parent
			}
		case ast.KindParameter:
			parameterDeclaration := location.AsParameterDeclaration()
			if lastLocation != nil && (lastLocation == parameterDeclaration.Initializer ||
				lastLocation == parameterDeclaration.Name_ && isBindingPattern(lastLocation)) {
				if associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case ast.KindBindingElement:
			bindingElement := location.AsBindingElement()
			if lastLocation != nil && (lastLocation == bindingElement.Initializer ||
				lastLocation == bindingElement.Name_ && isBindingPattern(lastLocation)) {
				if isPartOfParameterDeclaration(location) && associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case ast.KindInferType:
			if meaning&ast.SymbolFlagsTypeParameter != 0 {
				parameterName := location.AsInferTypeNode().TypeParameter.AsTypeParameter().Name_
				if parameterName != nil && name == parameterName.AsIdentifier().Text {
					result = location.AsInferTypeNode().TypeParameter.AsTypeParameter().Symbol
					break loop
				}
			}
		case ast.KindExportSpecifier:
			exportSpecifier := location.AsExportSpecifier()
			if lastLocation != nil && lastLocation == exportSpecifier.PropertyName && location.Parent.Parent.AsExportDeclaration().ModuleSpecifier != nil {
				location = location.Parent.Parent.Parent
			}
		}
		if isSelfReferenceLocation(location, lastLocation) {
			lastSelfReferenceLocation = location
		}
		lastLocation = location
		switch {
		// case isJSDocTemplateTag(location):
		// 	location = getEffectiveContainerForJSDocTemplateTag(location.(*JSDocTemplateTag))
		// 	if location == nil {
		// 		location = location.parent
		// 	}
		// case isJSDocParameterTag(location) || isJSDocReturnTag(location):
		// 	location = getHostSignatureFromJSDoc(location)
		// 	if location == nil {
		// 		location = location.parent
		// 	}
		default:
			location = location.Parent
		}
	}
	// We just climbed up parents looking for the name, meaning that we started in a descendant node of `lastLocation`.
	// If `result === lastSelfReferenceLocation.symbol`, that means that we are somewhere inside `lastSelfReferenceLocation` looking up a name, and resolving to `lastLocation` itself.
	// That means that this is a self-reference of `lastLocation`, and shouldn't count this when considering whether `lastLocation` is used.
	if isUse && result != nil && (lastSelfReferenceLocation == nil || result != lastSelfReferenceLocation.Symbol()) {
		// !!! result.isReferenced |= meaning
	}
	if result == nil {
		if !excludeGlobals {
			result = r.lookup(r.globals, name, meaning)
		}
	}
	if nameNotFoundMessage != nil {
		if propertyWithInvalidInitializer != nil && r.onPropertyWithInvalidInitializer(originalLocation, name, propertyWithInvalidInitializer, result) {
			return nil
		}
		if result == nil {
			r.onFailedToResolveSymbol(originalLocation, name, meaning, nameNotFoundMessage)
		} else {
			r.onSuccessfullyResolvedSymbol(originalLocation, result, meaning, lastLocation, associatedDeclarationForContainingInitializerOrBindingName, withinDeferredContext)
		}
	}
	return result
}

func (r *NameResolver) useOuterVariableScopeInParameter(result *Symbol, location *Node, lastLocation *Node) bool {
	if IsParameter(lastLocation) {
		body := getBodyOfNode(location)
		if body != nil && result.ValueDeclaration != nil && result.ValueDeclaration.Pos() >= body.Pos() && result.ValueDeclaration.End() <= body.End() {
			// check for several cases where we introduce temporaries that require moving the name/initializer of the parameter to the body
			// - static field in a class expression
			// - optional chaining pre-es2020
			// - nullish coalesce pre-es2020
			// - spread assignment in binding pattern pre-es2017
			target := getEmitScriptTarget(r.compilerOptions)
			if target >= core.ScriptTargetES2015 {
				functionLocation := location
				declarationRequiresScopeChange := r.getRequiresScopeChangeCache(functionLocation)
				if declarationRequiresScopeChange == core.TSUnknown {
					declarationRequiresScopeChange = boolToTristate(core.Some(functionLocation.Parameters(), r.requiresScopeChange))
					r.setRequiresScopeChangeCache(functionLocation, declarationRequiresScopeChange)
				}
				return declarationRequiresScopeChange == core.TSTrue
			}
		}
	}
	return false
}

func (r *NameResolver) requiresScopeChange(node *Node) bool {
	d := node.AsParameterDeclaration()
	return r.requiresScopeChangeWorker(d.Name_) || d.Initializer != nil && r.requiresScopeChangeWorker(d.Initializer)
}

func (r *NameResolver) requiresScopeChangeWorker(node *Node) bool {
	switch node.Kind {
	case ast.KindArrowFunction, ast.KindFunctionExpression, ast.KindFunctionDeclaration, ast.KindConstructor:
		return false
	case ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindPropertyAssignment:
		return r.requiresScopeChangeWorker(node.Name())
	case ast.KindPropertyDeclaration:
		if hasStaticModifier(node) {
			return !getEmitStandardClassFields(r.compilerOptions)
		}
		return r.requiresScopeChangeWorker(node.AsPropertyDeclaration().Name_)
	default:
		if isNullishCoalesce(node) || isOptionalChain(node) {
			return getEmitScriptTarget(r.compilerOptions) < core.ScriptTargetES2020
		}
		if IsBindingElement(node) && node.AsBindingElement().DotDotDotToken != nil && IsObjectBindingPattern(node.Parent) {
			return getEmitScriptTarget(r.compilerOptions) < core.ScriptTargetES2017
		}
		if isTypeNode(node) {
			return false
		}
		return node.ForEachChild(r.requiresScopeChangeWorker)
	}
}

func getIsDeferredContext(location *Node, lastLocation *Node) bool {
	if location.Kind != ast.KindArrowFunction && location.Kind != ast.KindFunctionExpression {
		// initializers in instance property declaration of class like entities are executed in constructor and thus deferred
		// A name is evaluated within the enclosing scope - so it shouldn't count as deferred
		return IsTypeQueryNode(location) ||
			(isFunctionLikeDeclaration(location) || location.Kind == ast.KindPropertyDeclaration && !isStatic(location)) &&
				(lastLocation == nil || lastLocation != location.Name())
	}
	if lastLocation != nil && lastLocation == location.Name() {
		return false
	}
	// generator functions and async functions are not inlined in control flow when immediately invoked
	if location.BodyData().AsteriskToken != nil || hasSyntacticModifier(location, ast.ModifierFlagsAsync) {
		return true
	}
	return getImmediatelyInvokedFunctionExpression(location) == nil
}

func isTypeParameterSymbolDeclaredInContainer(symbol *Symbol, container *Node) bool {
	for _, decl := range symbol.Declarations {
		if decl.Kind == ast.KindTypeParameter {
			parent := decl.Parent.Parent
			if parent == container {
				return true
			}
		}
	}
	return false
}

func isSelfReferenceLocation(node *Node, lastLocation *Node) bool {
	switch node.Kind {
	case ast.KindParameter:
		return lastLocation != nil && lastLocation == node.AsParameterDeclaration().Name_
	case ast.KindFunctionDeclaration, ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindEnumDeclaration,
		ast.KindTypeAliasDeclaration, ast.KindModuleDeclaration: // For `namespace N { N; }`
		return true
	}
	return false
}

func isTypeReferenceIdentifier(node *Node) bool {
	for node.Parent.Kind == ast.KindQualifiedName {
		node = node.Parent
	}
	return IsTypeReferenceNode(node.Parent)
}

func isInTypeQuery(node *Node) bool {
	// TypeScript 1.0 spec (April 2014): 3.6.3
	// A type query consists of the keyword typeof followed by an expression.
	// The expression is restricted to a single identifier or a sequence of identifiers separated by periods
	return findAncestorOrQuit(node, func(n *Node) FindAncestorResult {
		switch n.Kind {
		case ast.KindTypeQuery:
			return FindAncestorTrue
		case ast.KindIdentifier, ast.KindQualifiedName:
			return FindAncestorFalse
		}
		return FindAncestorQuit
	}) != nil
}

func nodeKindIs(node *Node, kinds ...ast.Kind) bool {
	return slices.Contains(kinds, node.Kind)
}

func isTypeOnlyImportDeclaration(node *Node) bool {
	switch node.Kind {
	case ast.KindImportSpecifier:
		return node.AsImportSpecifier().IsTypeOnly || node.Parent.Parent.AsImportClause().IsTypeOnly
	case ast.KindNamespaceImport:
		return node.Parent.AsImportClause().IsTypeOnly
	case ast.KindImportClause:
		return node.AsImportClause().IsTypeOnly
	case ast.KindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().IsTypeOnly
	}
	return false
}

func isTypeOnlyExportDeclaration(node *Node) bool {
	switch node.Kind {
	case ast.KindExportSpecifier:
		return node.AsExportSpecifier().IsTypeOnly || node.Parent.Parent.AsExportDeclaration().IsTypeOnly
	case ast.KindExportDeclaration:
		d := node.AsExportDeclaration()
		return d.IsTypeOnly && d.ModuleSpecifier != nil && d.ExportClause == nil
	case ast.KindNamespaceExport:
		return node.Parent.AsExportDeclaration().IsTypeOnly
	}
	return false
}

func isTypeOnlyImportOrExportDeclaration(node *Node) bool {
	return isTypeOnlyImportDeclaration(node) || isTypeOnlyExportDeclaration(node)
}

func getNameFromImportDeclaration(node *Node) *Node {
	switch node.Kind {
	case ast.KindImportSpecifier:
		return node.AsImportSpecifier().Name_
	case ast.KindNamespaceImport:
		return node.AsNamespaceImport().Name_
	case ast.KindImportClause:
		return node.AsImportClause().Name_
	case ast.KindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().Name_
	}
	return nil
}

func isValidTypeOnlyAliasUseSite(useSite *Node) bool {
	return useSite.Flags&ast.NodeFlagsAmbient != 0 ||
		isPartOfTypeQuery(useSite) ||
		isIdentifierInNonEmittingHeritageClause(useSite) ||
		isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(useSite) ||
		!(isExpressionNode(useSite) || isShorthandPropertyNameUseSite(useSite))
}

func isIdentifierInNonEmittingHeritageClause(node *Node) bool {
	if node.Kind != ast.KindIdentifier {
		return false
	}
	heritageClause := findAncestorOrQuit(node.Parent, func(parent *Node) FindAncestorResult {
		switch parent.Kind {
		case ast.KindHeritageClause:
			return FindAncestorTrue
		case ast.KindPropertyAccessExpression, ast.KindExpressionWithTypeArguments:
			return FindAncestorFalse
		default:
			return FindAncestorQuit
		}
	})
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().Token == ast.KindImmediateKeyword || heritageClause.Parent.Kind == ast.KindInterfaceDeclaration
	}
	return false
}

func isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(node *Node) bool {
	for nodeKindIs(node, ast.KindIdentifier, ast.KindPropertyAccessExpression) {
		node = node.Parent
	}
	if node.Kind != ast.KindComputedPropertyName {
		return false
	}
	if hasSyntacticModifier(node.Parent, ast.ModifierFlagsAbstract) {
		return true
	}
	return nodeKindIs(node.Parent.Parent, ast.KindInterfaceDeclaration, ast.KindTypeLiteral)
}

func isExpressionNode(node *Node) bool {
	switch node.Kind {
	case ast.KindSuperKeyword, ast.KindNullKeyword, ast.KindTrueKeyword, ast.KindFalseKeyword, ast.KindRegularExpressionLiteral,
		ast.KindArrayLiteralExpression, ast.KindObjectLiteralExpression, ast.KindPropertyAccessExpression, ast.KindElementAccessExpression,
		ast.KindCallExpression, ast.KindNewExpression, ast.KindTaggedTemplateExpression, ast.KindAsExpression, ast.KindTypeAssertionExpression,
		ast.KindSatisfiesExpression, ast.KindNonNullExpression, ast.KindParenthesizedExpression, ast.KindFunctionExpression,
		ast.KindClassExpression, ast.KindArrowFunction, ast.KindVoidExpression, ast.KindDeleteExpression, ast.KindTypeOfExpression,
		ast.KindPrefixUnaryExpression, ast.KindPostfixUnaryExpression, ast.KindBinaryExpression, ast.KindConditionalExpression,
		ast.KindSpreadElement, ast.KindTemplateExpression, ast.KindOmittedExpression, ast.KindJsxElement, ast.KindJsxSelfClosingElement,
		ast.KindJsxFragment, ast.KindYieldExpression, ast.KindAwaitExpression, ast.KindMetaProperty:
		return true
	case ast.KindExpressionWithTypeArguments:
		return !IsHeritageClause(node.Parent)
	case ast.KindQualifiedName:
		for node.Parent.Kind == ast.KindQualifiedName {
			node = node.Parent
		}
		return IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node)
	case ast.KindJSDocMemberName:
		return IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node)
	case ast.KindPrivateIdentifier:
		return isBinaryExpression(node.Parent) && node.Parent.AsBinaryExpression().Left == node && node.Parent.AsBinaryExpression().OperatorToken.Kind == ast.KindInKeyword
	case ast.KindIdentifier:
		if IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node) {
			return true
		}
		fallthrough
	case ast.KindNumericLiteral, ast.KindBigIntLiteral, ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral, ast.KindThisKeyword:
		return isInExpressionContext(node)
	default:
		return false
	}
}

func isInExpressionContext(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindVariableDeclaration:
		return parent.AsVariableDeclaration().Initializer == node
	case ast.KindParameter:
		return parent.AsParameterDeclaration().Initializer == node
	case ast.KindPropertyDeclaration:
		return parent.AsPropertyDeclaration().Initializer == node
	case ast.KindPropertySignature:
		return parent.AsPropertySignatureDeclaration().Initializer == node
	case ast.KindEnumMember:
		return parent.AsEnumMember().Initializer == node
	case ast.KindPropertyAssignment:
		return parent.AsPropertyAssignment().Initializer == node
	case ast.KindBindingElement:
		return parent.AsBindingElement().Initializer == node
	case ast.KindExpressionStatement:
		return parent.AsExpressionStatement().Expression == node
	case ast.KindIfStatement:
		return parent.AsIfStatement().Expression == node
	case ast.KindDoStatement:
		return parent.AsDoStatement().Expression == node
	case ast.KindWhileStatement:
		return parent.AsWhileStatement().Expression == node
	case ast.KindReturnStatement:
		return parent.AsReturnStatement().Expression == node
	case ast.KindWithStatement:
		return parent.AsWithStatement().Expression == node
	case ast.KindSwitchStatement:
		return parent.AsSwitchStatement().Expression == node
	case ast.KindCaseClause, ast.KindDefaultClause:
		return parent.AsCaseOrDefaultClause().Expression == node
	case ast.KindThrowStatement:
		return parent.AsThrowStatement().Expression == node
	case ast.KindForStatement:
		s := parent.AsForStatement()
		return s.Initializer == node && s.Initializer.Kind != ast.KindVariableDeclarationList || s.Condition == node || s.Incrementor == node
	case ast.KindForInStatement, ast.KindForOfStatement:
		s := parent.AsForInOrOfStatement()
		return s.Initializer == node && s.Initializer.Kind != ast.KindVariableDeclarationList || s.Expression == node
	case ast.KindTypeAssertionExpression:
		return parent.AsTypeAssertion().Expression == node
	case ast.KindAsExpression:
		return parent.AsAsExpression().Expression == node
	case ast.KindTemplateSpan:
		return parent.AsTemplateSpan().Expression == node
	case ast.KindComputedPropertyName:
		return parent.AsComputedPropertyName().Expression == node
	case ast.KindDecorator, ast.KindJsxExpression, ast.KindJsxSpreadAttribute, ast.KindSpreadAssignment:
		return true
	case ast.KindExpressionWithTypeArguments:
		return parent.AsExpressionWithTypeArguments().Expression == node && !isPartOfTypeNode(parent)
	case ast.KindShorthandPropertyAssignment:
		return parent.AsShorthandPropertyAssignment().ObjectAssignmentInitializer == node
	case ast.KindSatisfiesExpression:
		return parent.AsSatisfiesExpression().Expression == node
	default:
		return isExpressionNode(parent)
	}
}

func isPartOfTypeNode(node *Node) bool {
	kind := node.Kind
	if kind >= ast.KindFirstTypeNode && kind <= ast.KindLastTypeNode {
		return true
	}
	switch node.Kind {
	case ast.KindAnyKeyword, ast.KindUnknownKeyword, ast.KindNumberKeyword, ast.KindBigIntKeyword, ast.KindStringKeyword,
		ast.KindBooleanKeyword, ast.KindSymbolKeyword, ast.KindObjectKeyword, ast.KindUndefinedKeyword, ast.KindNullKeyword,
		ast.KindNeverKeyword:
		return true
	case ast.KindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(node)
	case ast.KindTypeParameter:
		return node.Parent.Kind == ast.KindMappedType || node.Parent.Kind == ast.KindInferType
	case ast.KindIdentifier:
		parent := node.Parent
		if IsQualifiedName(parent) && parent.AsQualifiedName().Right == node {
			return isPartOfTypeNodeInParent(parent)
		}
		if IsPropertyAccessExpression(parent) && parent.AsPropertyAccessExpression().Name_ == node {
			return isPartOfTypeNodeInParent(parent)
		}
		return isPartOfTypeNodeInParent(node)
	case ast.KindQualifiedName, ast.KindPropertyAccessExpression, ast.KindThisKeyword:
		return isPartOfTypeNodeInParent(node)
	}
	return false
}

func isPartOfTypeNodeInParent(node *Node) bool {
	parent := node.Parent
	// Do not recursively call isPartOfTypeNode on the parent. In the example:
	//
	//     let a: A.B.C;
	//
	// Calling isPartOfTypeNode would consider the qualified name A.B a type node.
	// Only C and A.B.C are type nodes.
	if parent.Kind >= ast.KindFirstTypeNode && parent.Kind <= ast.KindLastTypeNode {
		return true
	}
	switch parent.Kind {
	case ast.KindTypeQuery:
		return false
	case ast.KindImportType:
		return !parent.AsImportTypeNode().IsTypeOf
	case ast.KindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(parent)
	case ast.KindTypeParameter:
		return node == parent.AsTypeParameter().Constraint
	case ast.KindPropertyDeclaration:
		return node == parent.AsPropertyDeclaration().TypeNode
	case ast.KindPropertySignature:
		return node == parent.AsPropertySignatureDeclaration().TypeNode
	case ast.KindParameter:
		return node == parent.AsParameterDeclaration().TypeNode
	case ast.KindVariableDeclaration:
		return node == parent.AsVariableDeclaration().TypeNode
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindConstructor, ast.KindMethodDeclaration,
		ast.KindMethodSignature, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindCallSignature, ast.KindConstructSignature,
		ast.KindIndexSignature:
		return node == parent.ReturnType()
	case ast.KindTypeAssertionExpression:
		return node == parent.AsTypeAssertion().TypeNode
	case ast.KindCallExpression:
		return typeArgumentListContains(parent.AsCallExpression().TypeArguments, node)
	case ast.KindNewExpression:
		return typeArgumentListContains(parent.AsNewExpression().TypeArguments, node)
	case ast.KindTaggedTemplateExpression:
		return typeArgumentListContains(parent.AsTaggedTemplateExpression().TypeArguments, node)
	}
	return false
}

func isPartOfTypeExpressionWithTypeArguments(node *Node) bool {
	parent := node.Parent
	return IsHeritageClause(parent) && (!isClassLike(parent.Parent) || parent.AsHeritageClause().Token == ast.KindImplementsKeyword)
}

func typeArgumentListContains(list *Node, node *Node) bool {
	if list != nil {
		return slices.Contains(list.AsTypeArgumentList().Arguments, node)
	}
	return false
}

func isJSDocLinkLike(node *Node) bool {
	return nodeKindIs(node, ast.KindJSDocLink, ast.KindJSDocLinkCode, ast.KindJSDocLinkPlain)
}

func isJSXTagName(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindJsxOpeningElement:
		return parent.AsJsxOpeningElement().TagName == node
	case ast.KindJsxSelfClosingElement:
		return parent.AsJsxSelfClosingElement().TagName == node
	case ast.KindJsxClosingElement:
		return parent.AsJsxClosingElement().TagName == node
	}
	return false
}

func isShorthandPropertyNameUseSite(useSite *Node) bool {
	return IsIdentifier(useSite) && IsShorthandPropertyAssignment(useSite.Parent) && useSite.Parent.AsShorthandPropertyAssignment().Name_ == useSite
}

func isTypeDeclaration(node *Node) bool {
	switch node.Kind {
	case ast.KindTypeParameter, ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindTypeAliasDeclaration, ast.KindEnumDeclaration:
		return true
	case ast.KindImportClause:
		return node.AsImportClause().IsTypeOnly
	case ast.KindImportSpecifier:
		return node.Parent.Parent.AsImportClause().IsTypeOnly
	case ast.KindExportSpecifier:
		return node.Parent.Parent.AsExportDeclaration().IsTypeOnly
	default:
		return false
	}
}

func canHaveSymbol(node *Node) bool {
	switch node.Kind {
	case ast.KindArrowFunction, ast.KindBinaryExpression, ast.KindBindingElement, ast.KindCallExpression, ast.KindCallSignature,
		ast.KindClassDeclaration, ast.KindClassExpression, ast.KindClassStaticBlockDeclaration, ast.KindConstructor, ast.KindConstructorType,
		ast.KindConstructSignature, ast.KindElementAccessExpression, ast.KindEnumDeclaration, ast.KindEnumMember, ast.KindExportAssignment,
		ast.KindExportDeclaration, ast.KindExportSpecifier, ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindFunctionType,
		ast.KindGetAccessor, ast.KindIdentifier, ast.KindImportClause, ast.KindImportEqualsDeclaration, ast.KindImportSpecifier,
		ast.KindIndexSignature, ast.KindInterfaceDeclaration, ast.KindJSDocCallbackTag, ast.KindJSDocEnumTag, ast.KindJSDocFunctionType,
		ast.KindJSDocParameterTag, ast.KindJSDocPropertyTag, ast.KindJSDocSignature, ast.KindJSDocTypedefTag, ast.KindJSDocTypeLiteral,
		ast.KindJsxAttribute, ast.KindJsxAttributes, ast.KindJsxSpreadAttribute, ast.KindMappedType, ast.KindMethodDeclaration,
		ast.KindMethodSignature, ast.KindModuleDeclaration, ast.KindNamedTupleMember, ast.KindNamespaceExport, ast.KindNamespaceExportDeclaration,
		ast.KindNamespaceImport, ast.KindNewExpression, ast.KindNoSubstitutionTemplateLiteral, ast.KindNumericLiteral, ast.KindObjectLiteralExpression,
		ast.KindParameter, ast.KindPropertyAccessExpression, ast.KindPropertyAssignment, ast.KindPropertyDeclaration, ast.KindPropertySignature,
		ast.KindSetAccessor, ast.KindShorthandPropertyAssignment, ast.KindSourceFile, ast.KindSpreadAssignment, ast.KindStringLiteral,
		ast.KindTypeAliasDeclaration, ast.KindTypeLiteral, ast.KindTypeParameter, ast.KindVariableDeclaration:
		return true
	}
	return false
}

func canHaveLocals(node *Node) bool {
	switch node.Kind {
	case ast.KindArrowFunction, ast.KindBlock, ast.KindCallSignature, ast.KindCaseBlock, ast.KindCatchClause,
		ast.KindClassStaticBlockDeclaration, ast.KindConditionalType, ast.KindConstructor, ast.KindConstructorType,
		ast.KindConstructSignature, ast.KindForStatement, ast.KindForInStatement, ast.KindForOfStatement, ast.KindFunctionDeclaration,
		ast.KindFunctionExpression, ast.KindFunctionType, ast.KindGetAccessor, ast.KindIndexSignature, ast.KindJSDocCallbackTag,
		ast.KindJSDocEnumTag, ast.KindJSDocFunctionType, ast.KindJSDocSignature, ast.KindJSDocTypedefTag, ast.KindMappedType,
		ast.KindMethodDeclaration, ast.KindMethodSignature, ast.KindModuleDeclaration, ast.KindSetAccessor, ast.KindSourceFile,
		ast.KindTypeAliasDeclaration:
		return true
	}
	return false
}

func isAnyImportOrReExport(node *Node) bool {
	return isAnyImportSyntax(node) || IsExportDeclaration(node)
}

func isAnyImportSyntax(node *Node) bool {
	return nodeKindIs(node, ast.KindImportDeclaration, ast.KindImportEqualsDeclaration)
}

func getExternalModuleName(node *Node) *Node {
	switch node.Kind {
	case ast.KindImportDeclaration:
		return node.AsImportDeclaration().ModuleSpecifier
	case ast.KindExportDeclaration:
		return node.AsExportDeclaration().ModuleSpecifier
	case ast.KindImportEqualsDeclaration:
		if node.AsImportEqualsDeclaration().ModuleReference.Kind == ast.KindExternalModuleReference {
			return node.AsImportEqualsDeclaration().ModuleReference.AsExternalModuleReference().expression
		}
		return nil
	case ast.KindImportType:
		return getImportTypeNodeLiteral(node)
	case ast.KindCallExpression:
		return node.AsCallExpression().Arguments[0]
	case ast.KindModuleDeclaration:
		if IsStringLiteral(node.AsModuleDeclaration().Name_) {
			return node.AsModuleDeclaration().Name_
		}
		return nil
	}
	panic("Unhandled case in getExternalModuleName")
}

func getImportTypeNodeLiteral(node *Node) *Node {
	if IsImportTypeNode(node) {
		importTypeNode := node.AsImportTypeNode()
		if IsLiteralTypeNode(importTypeNode.Argument) {
			literalTypeNode := importTypeNode.Argument.AsLiteralTypeNode()
			if IsStringLiteral(literalTypeNode.Literal) {
				return literalTypeNode.Literal
			}
		}
	}
	return nil
}

func isExternalModuleNameRelative(moduleName string) bool {
	// TypeScript 1.0 spec (April 2014): 11.2.1
	// An external module name is "relative" if the first term is "." or "..".
	// Update: We also consider a path like `C:\foo.ts` "relative" because we do not search for it in `node_modules` or treat it as an ambient module.
	return pathIsRelative(moduleName) || tspath.IsRootedDiskPath(moduleName)
}

func pathIsRelative(path string) bool {
	return core.MakeRegexp(`^\.\.?(?:$|[\\/])`).MatchString(path)
}

func extensionIsTs(ext string) bool {
	return ext == ExtensionTs || ext == ExtensionTsx || ext == ExtensionDts || ext == ExtensionMts || ext == ExtensionDmts || ext == ExtensionCts || ext == ExtensionDcts || len(ext) >= 7 && ext[:3] == ".d." && ext[len(ext)-3:] == ".ts"
}

func isShorthandAmbientModuleSymbol(moduleSymbol *Symbol) bool {
	return isShorthandAmbientModule(moduleSymbol.ValueDeclaration)
}

func isShorthandAmbientModule(node *Node) bool {
	// The only kind of module that can be missing a body is a shorthand ambient module.
	return node != nil && node.Kind == ast.KindModuleDeclaration && node.AsModuleDeclaration().Body == nil
}

func isEntityName(node *Node) bool {
	return node.Kind == ast.KindIdentifier || node.Kind == ast.KindQualifiedName
}

func nodeIsSynthesized(node *Node) bool {
	return node.Loc.Pos_ < 0 || node.Loc.End_ < 0
}

func getFirstIdentifier(node *Node) *Node {
	switch node.Kind {
	case ast.KindIdentifier:
		return node
	case ast.KindQualifiedName:
		return getFirstIdentifier(node.AsQualifiedName().Left)
	case ast.KindPropertyAccessExpression:
		return getFirstIdentifier(node.AsPropertyAccessExpression().Expression)
	}
	panic("Unhandled case in getFirstIdentifier")
}

func getAliasDeclarationFromName(node *Node) *Node {
	switch node.Kind {
	case ast.KindImportClause, ast.KindImportSpecifier, ast.KindNamespaceImport, ast.KindExportSpecifier, ast.KindExportAssignment,
		ast.KindImportEqualsDeclaration, ast.KindNamespaceExport:
		return node.Parent
	case ast.KindQualifiedName:
		return getAliasDeclarationFromName(node.Parent)
	}
	return nil
}

func entityNameToString(name *Node) string {
	switch name.Kind {
	case ast.KindThisKeyword:
		return "this"
	case ast.KindIdentifier, ast.KindPrivateIdentifier:
		return getTextOfNode(name)
	case ast.KindQualifiedName:
		return entityNameToString(name.AsQualifiedName().Left) + "." + entityNameToString(name.AsQualifiedName().Right)
	case ast.KindPropertyAccessExpression:
		return entityNameToString(name.AsPropertyAccessExpression().Expression) + "." + entityNameToString(name.AsPropertyAccessExpression().Name_)
	case ast.KindJsxNamespacedName:
		return entityNameToString(name.AsJsxNamespacedName().Namespace) + ":" + entityNameToString(name.AsJsxNamespacedName().Name_)
	}
	panic("Unhandled case in entityNameToString")
}

func getContainingQualifiedNameNode(node *Node) *Node {
	for IsQualifiedName(node.Parent) {
		node = node.Parent
	}
	return node
}

var extensionsToRemove = []string{ExtensionDts, ExtensionDmts, ExtensionDcts, ExtensionMjs, ExtensionMts, ExtensionCjs, ExtensionCts, ExtensionTs, ExtensionJs, ExtensionTsx, ExtensionJsx, ExtensionJson}

func removeFileExtension(path string) string {
	// Remove any known extension even if it has more than one dot
	for _, ext := range extensionsToRemove {
		if strings.HasSuffix(path, ext) {
			return path[:len(path)-len(ext)]
		}
	}
	// Otherwise just remove single dot extension, if any
	return path[:len(path)-len(filepath.Ext(path))]
}

func isSideEffectImport(node *Node) bool {
	ancestor := findAncestor(node, IsImportDeclaration)
	return ancestor != nil && ancestor.AsImportDeclaration().ImportClause == nil
}

func getExternalModuleRequireArgument(node *Node) *Node {
	if isVariableDeclarationInitializedToBareOrAccessedRequire(node) {
		return getLeftmostAccessExpression(node.AsVariableDeclaration().Initializer).AsCallExpression().Arguments[0]
	}
	return nil
}

func getExternalModuleImportEqualsDeclarationExpression(node *Node) *Node {
	//Debug.assert(isExternalModuleImportEqualsDeclaration(node))
	return node.AsImportEqualsDeclaration().ModuleReference.AsExternalModuleReference().expression
}

func isRightSideOfQualifiedNameOrPropertyAccess(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindQualifiedName:
		return parent.AsQualifiedName().Right == node
	case ast.KindPropertyAccessExpression:
		return parent.AsPropertyAccessExpression().Name_ == node
	case ast.KindMetaProperty:
		return parent.AsMetaProperty().Name_ == node
	}
	return false
}

func getNamespaceDeclarationNode(node *Node) *Node {
	switch node.Kind {
	case ast.KindImportDeclaration:
		importClause := node.AsImportDeclaration().ImportClause
		if importClause != nil && IsNamespaceImport(importClause.AsImportClause().NamedBindings) {
			return importClause.AsImportClause().NamedBindings
		}
	case ast.KindImportEqualsDeclaration:
		return node
	case ast.KindExportDeclaration:
		exportClause := node.AsExportDeclaration().ExportClause
		if exportClause != nil && IsNamespaceExport(exportClause) {
			return exportClause
		}
	default:
		panic("Unhandled case in getNamespaceDeclarationNode")
	}
	return nil
}

func isImportCall(node *Node) bool {
	return IsCallExpression(node) && node.AsCallExpression().Expression.Kind == ast.KindImportKeyword
}

func getSourceFileOfModule(module *Symbol) *SourceFile {
	declaration := module.ValueDeclaration
	if declaration == nil {
		declaration = getNonAugmentationDeclaration(module)
	}
	return getSourceFileOfNode(declaration)
}

func getNonAugmentationDeclaration(symbol *Symbol) *Node {
	return core.Find(symbol.Declarations, func(d *Node) bool {
		return !isExternalModuleAugmentation(d) && !(IsModuleDeclaration(d) && isGlobalScopeAugmentation(d))
	})
}

func isExternalModuleAugmentation(node *Node) bool {
	return isAmbientModule(node) && isModuleAugmentationExternal(node)
}

func isJsonSourceFile(file *SourceFile) bool {
	return file.ScriptKind == core.ScriptKindJSON
}

func isSyntacticDefault(node *Node) bool {
	return (IsExportAssignment(node) && !node.AsExportAssignment().IsExportEquals) ||
		hasSyntacticModifier(node, ast.ModifierFlagsDefault) ||
		IsExportSpecifier(node) ||
		IsNamespaceExport(node)
}

func hasExportAssignmentSymbol(moduleSymbol *Symbol) bool {
	return moduleSymbol.Exports[InternalSymbolNameExportEquals] != nil
}

func isImportOrExportSpecifier(node *Node) bool {
	return IsImportSpecifier(node) || IsExportSpecifier(node)
}

func parsePseudoBigInt(stringValue string) string {
	return stringValue // !!!
}

func isTypeAlias(node *Node) bool {
	return IsTypeAliasDeclaration(node)
}

/**
 * Gets the effective type parameters. If the node was parsed in a
 * JavaScript file, gets the type parameters from the `@template` tag from JSDoc.
 *
 * This does *not* return type parameters from a jsdoc reference to a generic type, eg
 *
 * type Id = <T>(x: T) => T
 * /** @type {Id} /
 * function id(x) { return x }
 */

func getEffectiveTypeParameterDeclarations(node *Node) []*Node {
	// if isJSDocSignature(node) {
	// 	if isJSDocOverloadTag(node.parent) {
	// 		jsDoc := getJSDocRoot(node.parent)
	// 		if jsDoc && length(jsDoc.tags) {
	// 			return flatMap(jsDoc.tags, func(tag JSDocTag) *NodeArray[TypeParameterDeclaration] {
	// 				if isJSDocTemplateTag(tag) {
	// 					return tag.typeParameters
	// 				} else {
	// 					return nil
	// 				}
	// 			})
	// 		}
	// 	}
	// 	return emptyArray
	// }
	// if isJSDocTypeAlias(node) {
	// 	Debug.assert(node.parent.kind == ast.KindJSDoc)
	// 	return flatMap(node.parent.tags, func(tag JSDocTag) *NodeArray[TypeParameterDeclaration] {
	// 		if isJSDocTemplateTag(tag) {
	// 			return tag.typeParameters
	// 		} else {
	// 			return nil
	// 		}
	// 	})
	// }
	typeParameters := node.TypeParameters()
	if typeParameters != nil {
		return typeParameters.AsTypeParameterList().Parameters
	}
	// if isInJSFile(node) {
	// 	decls := getJSDocTypeParameterDeclarations(node)
	// 	if decls.length {
	// 		return decls
	// 	}
	// 	typeTag := getJSDocType(node)
	// 	if typeTag && isFunctionTypeNode(typeTag) && typeTag.typeParameters {
	// 		return typeTag.typeParameters
	// 	}
	// }
	return nil
}

func getTypeParameterNodesFromNode(node *Node) []*Node {
	typeParameterList := node.TypeParameters()
	if typeParameterList != nil {
		return typeParameterList.AsTypeParameterList().Parameters
	}
	return nil
}

func getTypeArgumentNodesFromNode(node *Node) []*Node {
	typeArgumentList := getTypeArgumentListFromNode(node)
	if typeArgumentList != nil {
		return typeArgumentList.AsTypeArgumentList().Arguments
	}
	return nil
}

func getTypeArgumentListFromNode(node *Node) *Node {
	switch node.Kind {
	case ast.KindCallExpression:
		return node.AsCallExpression().TypeArguments
	case ast.KindNewExpression:
		return node.AsNewExpression().TypeArguments
	case ast.KindTaggedTemplateExpression:
		return node.AsTaggedTemplateExpression().TypeArguments
	case ast.KindTypeReference:
		return node.AsTypeReference().TypeArguments
	case ast.KindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().TypeArguments
	case ast.KindImportType:
		return node.AsImportTypeNode().TypeArguments
	case ast.KindTypeQuery:
		return node.AsTypeQueryNode().TypeArguments
	}
	panic("Unhandled case in getTypeArgumentListFromNode")
}

func getInitializerFromNode(node *Node) *Node {
	switch node.Kind {
	case ast.KindVariableDeclaration:
		return node.AsVariableDeclaration().Initializer
	case ast.KindParameter:
		return node.AsParameterDeclaration().Initializer
	case ast.KindBindingElement:
		return node.AsBindingElement().Initializer
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().Initializer
	case ast.KindPropertyAssignment:
		return node.AsPropertyAssignment().Initializer
	case ast.KindEnumMember:
		return node.AsEnumMember().Initializer
	case ast.KindForStatement:
		return node.AsForStatement().Initializer
	case ast.KindForInStatement, ast.KindForOfStatement:
		return node.AsForInOrOfStatement().Initializer
	case ast.KindJsxAttribute:
		return node.AsJsxAttribute().Initializer
	}
	return nil
}

/**
 * Gets the effective type annotation of a variable, parameter, or property. If the node was
 * parsed in a JavaScript file, gets the type annotation from JSDoc.  Also gets the type of
 * functions only the JSDoc case.
 */
func getEffectiveTypeAnnotationNode(node *Node) *Node {
	switch node.Kind {
	case ast.KindVariableDeclaration:
		return node.AsVariableDeclaration().TypeNode
	case ast.KindParameter:
		return node.AsParameterDeclaration().TypeNode
	case ast.KindPropertySignature:
		return node.AsPropertySignatureDeclaration().TypeNode
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().TypeNode
	case ast.KindTypePredicate:
		return node.AsTypePredicateNode().TypeNode
	case ast.KindParenthesizedType:
		return node.AsParenthesizedTypeNode().TypeNode
	case ast.KindTypeOperator:
		return node.AsTypeOperatorNode().TypeNode
	case ast.KindMappedType:
		return node.AsMappedTypeNode().TypeNode
	case ast.KindTypeAssertionExpression:
		return node.AsTypeAssertion().TypeNode
	case ast.KindAsExpression:
		return node.AsAsExpression().TypeNode
	default:
		if isFunctionLike(node) {
			return node.ReturnType()
		}
	}
	return nil
}

func isTypeAny(t *Type) bool {
	return t != nil && t.flags&TypeFlagsAny != 0
}

func isJSDocOptionalParameter(node *ParameterDeclaration) bool {
	return false // !!!
}

func isQuestionToken(node *Node) bool {
	return node != nil && node.Kind == ast.KindQuestionToken
}

func isOptionalDeclaration(declaration *Node) bool {
	switch declaration.Kind {
	case ast.KindParameter:
		return declaration.AsParameterDeclaration().QuestionToken != nil
	case ast.KindPropertyDeclaration:
		return isQuestionToken(declaration.AsPropertyDeclaration().PostfixToken)
	case ast.KindPropertySignature:
		return isQuestionToken(declaration.AsPropertySignatureDeclaration().PostfixToken)
	case ast.KindMethodDeclaration:
		return isQuestionToken(declaration.AsMethodDeclaration().PostfixToken)
	case ast.KindMethodSignature:
		return isQuestionToken(declaration.AsMethodSignatureDeclaration().PostfixToken)
	case ast.KindPropertyAssignment:
		return isQuestionToken(declaration.AsPropertyAssignment().PostfixToken)
	case ast.KindShorthandPropertyAssignment:
		return isQuestionToken(declaration.AsShorthandPropertyAssignment().PostfixToken)
	}
	return false
}

func isEmptyArrayLiteral(expression *Node) bool {
	return expression.Kind == ast.KindArrayLiteralExpression && len(expression.AsArrayLiteralExpression().Elements) == 0
}

func declarationBelongsToPrivateAmbientMember(declaration *Node) bool {
	root := getRootDeclaration(declaration)
	memberDeclaration := root
	if root.Kind == ast.KindParameter {
		memberDeclaration = root.Parent
	}
	return isPrivateWithinAmbient(memberDeclaration)
}

func isPrivateWithinAmbient(node *Node) bool {
	return (hasEffectiveModifier(node, ast.ModifierFlagsPrivate) || isPrivateIdentifierClassElementDeclaration(node)) && node.Flags&ast.NodeFlagsAmbient != 0
}

func identifierToKeywordKind(node *Identifier) ast.Kind {
	return textToKeyword[node.Text]
}

func isAssertionExpression(node *Node) bool {
	kind := node.Kind
	return kind == ast.KindTypeAssertionExpression || kind == ast.KindAsExpression
}

func isTypeAssertion(node *Node) bool {
	return isAssertionExpression(skipParentheses(node))
}

func createSymbolTable(symbols []*Symbol) SymbolTable {
	if len(symbols) == 0 {
		return nil
	}
	result := make(SymbolTable)
	for _, symbol := range symbols {
		result[symbol.Name] = symbol
	}
	return result
}

func sortSymbols(symbols []*Symbol) {
	slices.SortFunc(symbols, compareSymbols)
}

func compareSymbols(s1, s2 *Symbol) int {
	if s1 == s2 {
		return 0
	}
	if s1.ValueDeclaration != nil && s2.ValueDeclaration != nil {
		if s1.Parent != nil && s2.Parent != nil {
			// Symbols with the same unmerged parent are always in the same file
			if s1.Parent != s2.Parent {
				f1 := getSourceFileOfNode(s1.ValueDeclaration)
				f2 := getSourceFileOfNode(s2.ValueDeclaration)
				if f1 != f2 {
					// In different files, first compare base filename
					r := strings.Compare(filepath.Base(f1.Path_), filepath.Base(f2.Path_))
					if r == 0 {
						// Same base filename, compare the full paths (no two files should have the same full path)
						r = strings.Compare(f1.Path_, f2.Path_)
					}
					return r
				}
			}
			// In the same file, compare source positions
			return s1.ValueDeclaration.Pos() - s2.ValueDeclaration.Pos()
		}
	}
	// Sort by name
	r := strings.Compare(s1.Name, s2.Name)
	if r == 0 {
		// Same name, sort by symbol id
		r = int(getSymbolId(s1)) - int(getSymbolId(s2))
	}
	return r
}

func getClassLikeDeclarationOfSymbol(symbol *Symbol) *Node {
	return core.Find(symbol.Declarations, isClassLike)
}

func isThisInTypeQuery(node *Node) bool {
	if !isThisIdentifier(node) {
		return false
	}
	for IsQualifiedName(node.Parent) && node.Parent.AsQualifiedName().Left == node {
		node = node.Parent
	}
	return node.Parent.Kind == ast.KindTypeQuery
}

func isThisIdentifier(node *Node) bool {
	return node != nil && node.Kind == ast.KindIdentifier && identifierIsThisKeyword(node)
}

func identifierIsThisKeyword(id *Node) bool {
	return id.AsIdentifier().Text == "this"
}

func getDeclarationModifierFlagsFromSymbol(s *Symbol) ast.ModifierFlags {
	return getDeclarationModifierFlagsFromSymbolEx(s, false /*isWrite*/)
}

func getDeclarationModifierFlagsFromSymbolEx(s *Symbol, isWrite bool) ast.ModifierFlags {
	if s.ValueDeclaration != nil {
		var declaration *Node
		if isWrite {
			declaration = core.Find(s.Declarations, IsSetAccessorDeclaration)
		}
		if declaration == nil && s.Flags&ast.SymbolFlagsGetAccessor != 0 {
			declaration = core.Find(s.Declarations, IsGetAccessorDeclaration)
		}
		if declaration == nil {
			declaration = s.ValueDeclaration
		}
		flags := getCombinedModifierFlags(declaration)
		if s.Parent != nil && s.Parent.Flags&ast.SymbolFlagsClass != 0 {
			return flags
		}
		return flags & ^ast.ModifierFlagsAccessibilityModifier
	}
	if s.CheckFlags&ast.CheckFlagsSynthetic != 0 {
		var accessModifier ast.ModifierFlags
		switch {
		case s.CheckFlags&ast.CheckFlagsContainsPrivate != 0:
			accessModifier = ast.ModifierFlagsPrivate
		case s.CheckFlags&ast.CheckFlagsContainsPublic != 0:
			accessModifier = ast.ModifierFlagsPublic
		default:
			accessModifier = ast.ModifierFlagsProtected
		}
		var staticModifier ast.ModifierFlags
		if s.CheckFlags&ast.CheckFlagsContainsStatic != 0 {
			staticModifier = ast.ModifierFlagsStatic
		}
		return accessModifier | staticModifier
	}
	if s.Flags&ast.SymbolFlagsPrototype != 0 {
		return ast.ModifierFlagsPublic | ast.ModifierFlagsStatic
	}
	return ast.ModifierFlagsNone
}

func isExponentiationOperator(kind ast.Kind) bool {
	return kind == ast.KindAsteriskAsteriskToken
}

func isMultiplicativeOperator(kind ast.Kind) bool {
	return kind == ast.KindAsteriskToken || kind == ast.KindSlashToken || kind == ast.KindPercentToken
}

func isMultiplicativeOperatorOrHigher(kind ast.Kind) bool {
	return isExponentiationOperator(kind) || isMultiplicativeOperator(kind)
}

func isAdditiveOperator(kind ast.Kind) bool {
	return kind == ast.KindPlusToken || kind == ast.KindMinusToken
}

func isAdditiveOperatorOrHigher(kind ast.Kind) bool {
	return isAdditiveOperator(kind) || isMultiplicativeOperatorOrHigher(kind)
}

func isShiftOperator(kind ast.Kind) bool {
	return kind == ast.KindLessThanLessThanToken || kind == ast.KindGreaterThanGreaterThanToken ||
		kind == ast.KindGreaterThanGreaterThanGreaterThanToken
}

func isShiftOperatorOrHigher(kind ast.Kind) bool {
	return isShiftOperator(kind) || isAdditiveOperatorOrHigher(kind)
}

func isRelationalOperator(kind ast.Kind) bool {
	return kind == ast.KindLessThanToken || kind == ast.KindLessThanEqualsToken || kind == ast.KindGreaterThanToken ||
		kind == ast.KindGreaterThanEqualsToken || kind == ast.KindInstanceOfKeyword || kind == ast.KindInKeyword
}

func isRelationalOperatorOrHigher(kind ast.Kind) bool {
	return isRelationalOperator(kind) || isShiftOperatorOrHigher(kind)
}

func isEqualityOperator(kind ast.Kind) bool {
	return kind == ast.KindEqualsEqualsToken || kind == ast.KindEqualsEqualsEqualsToken ||
		kind == ast.KindExclamationEqualsToken || kind == ast.KindExclamationEqualsEqualsToken
}

func isEqualityOperatorOrHigher(kind ast.Kind) bool {
	return isEqualityOperator(kind) || isRelationalOperatorOrHigher(kind)
}

func isBitwiseOperator(kind ast.Kind) bool {
	return kind == ast.KindAmpersandToken || kind == ast.KindBarToken || kind == ast.KindCaretToken
}

func isBitwiseOperatorOrHigher(kind ast.Kind) bool {
	return isBitwiseOperator(kind) || isEqualityOperatorOrHigher(kind)
}

// NOTE: The version in utilities includes ExclamationToken, which is not a binary operator.
func isLogicalOperator(kind ast.Kind) bool {
	return kind == ast.KindAmpersandAmpersandToken || kind == ast.KindBarBarToken
}

func isLogicalOperatorOrHigher(kind ast.Kind) bool {
	return isLogicalOperator(kind) || isBitwiseOperatorOrHigher(kind)
}

func isAssignmentOperatorOrHigher(kind ast.Kind) bool {
	return kind == ast.KindQuestionQuestionToken || isLogicalOperatorOrHigher(kind) || isAssignmentOperator(kind)
}

func isBinaryOperator(kind ast.Kind) bool {
	return isAssignmentOperatorOrHigher(kind) || kind == ast.KindCommaToken
}

func isObjectLiteralType(t *Type) bool {
	return t.objectFlags&ObjectFlagsObjectLiteral != 0
}

func isDeclarationReadonly(declaration *Node) bool {
	return getCombinedModifierFlags(declaration)&ast.ModifierFlagsReadonly != 0 && !isParameterPropertyDeclaration(declaration, declaration.Parent)
}

func getPostfixTokenFromNode(node *Node) *Node {
	switch node.Kind {
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().PostfixToken
	case ast.KindPropertySignature:
		return node.AsPropertySignatureDeclaration().PostfixToken
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration().PostfixToken
	case ast.KindMethodSignature:
		return node.AsMethodSignatureDeclaration().PostfixToken
	}
	panic("Unhandled case in getPostfixTokenFromNode")
}

func isStatic(node *Node) bool {
	// https://tc39.es/ecma262/#sec-static-semantics-isstatic
	return isClassElement(node) && hasStaticModifier(node) || IsClassStaticBlockDeclaration(node)
}

func isLogicalExpression(node *Node) bool {
	for {
		if node.Kind == ast.KindParenthesizedExpression {
			node = node.AsParenthesizedExpression().Expression
		} else if node.Kind == ast.KindPrefixUnaryExpression && node.AsPrefixUnaryExpression().Operator == ast.KindExclamationToken {
			node = node.AsPrefixUnaryExpression().Operand
		} else {
			return isLogicalOrCoalescingBinaryExpression(node)
		}
	}
}

type orderedMap[K comparable, V any] struct {
	valuesByKey map[K]V
	values      []V
}

func (m *orderedMap[K, V]) contains(key K) bool {
	_, ok := m.valuesByKey[key]
	return ok
}

func (m *orderedMap[K, V]) add(key K, value V) {
	if m.valuesByKey == nil {
		m.valuesByKey = make(map[K]V)
	}
	m.valuesByKey[key] = value
	m.values = append(m.values, value)
}

type Set[T comparable] struct {
	m map[T]struct{}
}

func (s *Set[T]) Has(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *Set[T]) Add(key T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[key] = struct{}{}
}

func (s *Set[T]) Delete(key T) {
	delete(s.m, key)
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) Keys() map[T]struct{} {
	return s.m
}

func getContainingFunction(node *Node) *Node {
	return findAncestor(node.Parent, isFunctionLike)
}

func isTypeReferenceType(node *Node) bool {
	return node.Kind == ast.KindTypeReference || node.Kind == ast.KindExpressionWithTypeArguments
}

func isNodeDescendantOf(node *Node, ancestor *Node) bool {
	for node != nil {
		if node == ancestor {
			return true
		}
		node = node.Parent
	}
	return false
}

func isTypeUsableAsPropertyName(t *Type) bool {
	return t.flags&TypeFlagsStringOrNumberLiteralOrUnique != 0
}

/**
 * Gets the symbolic name for a member from its type.
 */
func getPropertyNameFromType(t *Type) string {
	switch {
	case t.flags&TypeFlagsStringLiteral != 0:
		return t.AsLiteralType().value.(string)
	case t.flags&TypeFlagsNumberLiteral != 0:
		return numberToString(t.AsLiteralType().value.(float64))
	case t.flags&TypeFlagsUniqueESSymbol != 0:
		return t.AsUniqueESSymbolType().name
	}
	panic("Unhandled case in getPropertyNameFromType")
}

func isNumericLiteralName(name string) bool {
	// The intent of numeric names is that
	//     - they are names with text in a numeric form, and that
	//     - setting properties/indexing with them is always equivalent to doing so with the numeric literal 'numLit',
	//         acquired by applying the abstract 'ToNumber' operation on the name's text.
	//
	// The subtlety is in the latter portion, as we cannot reliably say that anything that looks like a numeric literal is a numeric name.
	// In fact, it is the case that the text of the name must be equal to 'ToString(numLit)' for this to hold.
	//
	// Consider the property name '"0xF00D"'. When one indexes with '0xF00D', they are actually indexing with the value of 'ToString(0xF00D)'
	// according to the ECMAScript specification, so it is actually as if the user indexed with the string '"61453"'.
	// Thus, the text of all numeric literals equivalent to '61543' such as '0xF00D', '0xf00D', '0170015', etc. are not valid numeric names
	// because their 'ToString' representation is not equal to their original text.
	// This is motivated by ECMA-262 sections 9.3.1, 9.8.1, 11.1.5, and 11.2.1.
	//
	// Here, we test whether 'ToString(ToNumber(name))' is exactly equal to 'name'.
	// The '+' prefix operator is equivalent here to applying the abstract ToNumber operation.
	// Applying the 'toString()' method on a number gives us the abstract ToString operation on a number.
	//
	// Note that this accepts the values 'Infinity', '-Infinity', and 'NaN', and that this is intentional.
	// This is desired behavior, because when indexing with them as numeric entities, you are indexing
	// with the strings '"Infinity"', '"-Infinity"', and '"NaN"' respectively.
	return numberToString(stringToNumber(name)) == name
}

func isPropertyName(node *Node) bool {
	switch node.Kind {
	case ast.KindIdentifier, ast.KindPrivateIdentifier, ast.KindStringLiteral, ast.KindNumericLiteral, ast.KindComputedPropertyName:
		return true
	}
	return false
}

func getPropertyNameForPropertyNameNode(name *Node) string {
	switch name.Kind {
	case ast.KindIdentifier, ast.KindPrivateIdentifier, ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral,
		ast.KindNumericLiteral, ast.KindBigIntLiteral, ast.KindJsxNamespacedName:
		return name.Text()
	case ast.KindComputedPropertyName:
		nameExpression := name.AsComputedPropertyName().Expression
		if isStringOrNumericLiteralLike(nameExpression) {
			return nameExpression.Text()
		}
		if isSignedNumericLiteral(nameExpression) {
			text := nameExpression.AsPrefixUnaryExpression().Operand.Text()
			if nameExpression.AsPrefixUnaryExpression().Operator == ast.KindMinusToken {
				text = "-" + text
			}
			return text
		}
		return InternalSymbolNameMissing
	}
	panic("Unhandled case in getPropertyNameForPropertyNameNode")
}

func isThisProperty(node *Node) bool {
	return (IsPropertyAccessExpression(node) || IsElementAccessExpression(node)) && node.Expression().Kind == ast.KindThisKeyword
}

func numberToString(f float64) string {
	// !!! This function should behave identically to the expression `"" + f` in JS
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func stringToNumber(s string) float64 {
	// !!! This function should behave identically to the expression `+s` in JS
	// This includes parsing binary, octal, and hex numeric strings
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return value
}

func isValidESSymbolDeclaration(node *Node) bool {
	if IsVariableDeclaration(node) {
		return isVarConst(node) && IsIdentifier(node.AsVariableDeclaration().Name_) && isVariableDeclarationInVariableStatement(node)
	}
	if IsPropertyDeclaration(node) {
		return hasEffectiveReadonlyModifier(node) && hasStaticModifier(node)
	}
	return IsPropertySignatureDeclaration(node) && hasEffectiveReadonlyModifier(node)
}

func isVarConst(node *Node) bool {
	return getCombinedNodeFlags(node)&ast.NodeFlagsBlockScoped == ast.NodeFlagsConst
}

func isVariableDeclarationInVariableStatement(node *Node) bool {
	return IsVariableDeclarationList(node.Parent) && IsVariableStatement(node.Parent.Parent)
}

func isKnownSymbol(symbol *Symbol) bool {
	return isLateBoundName(symbol.Name)
}

func isLateBoundName(name string) bool {
	return len(name) >= 2 && name[0] == '\xfe' && name[1] == '@'
}

func getSymbolTable(data *SymbolTable) SymbolTable {
	if *data == nil {
		*data = make(SymbolTable)
	}
	return *data
}

func getMembers(symbol *Symbol) SymbolTable {
	return getSymbolTable(&symbol.Members)
}

func getExports(symbol *Symbol) SymbolTable {
	return getSymbolTable(&symbol.Exports)
}

func getLocals(container *Node) SymbolTable {
	return getSymbolTable(&container.LocalsContainerData().Locals_)
}

func getThisParameter(signature *Node) *Node {
	// callback tags do not currently support this parameters
	if len(signature.Parameters()) != 0 {
		thisParameter := signature.Parameters()[0]
		if parameterIsThisKeyword(thisParameter) {
			return thisParameter
		}
	}
	return nil
}

func parameterIsThisKeyword(parameter *Node) bool {
	return isThisIdentifier(parameter.Name())
}

func getInterfaceBaseTypeNodes(node *Node) []*Node {
	heritageClause := getHeritageClause(node.AsInterfaceDeclaration().HeritageClauses, ast.KindExtendsKeyword)
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().Types
	}
	return nil
}

func getHeritageClause(clauses []*Node, kind ast.Kind) *Node {
	for _, clause := range clauses {
		if clause.AsHeritageClause().Token == kind {
			return clause
		}
	}
	return nil
}

func getClassExtendsHeritageElement(node *Node) *Node {
	heritageClause := getHeritageClause(node.ClassLikeData().HeritageClauses, ast.KindExtendsKeyword)
	if heritageClause != nil && len(heritageClause.AsHeritageClause().Types) > 0 {
		return heritageClause.AsHeritageClause().Types[0]
	}
	return nil
}

func concatenateDiagnosticMessageChains(headChain *MessageChain, tailChain *MessageChain) {
	lastChain := headChain
	for len(lastChain.MessageChain_) != 0 {
		lastChain = lastChain.MessageChain_[0]
	}
	lastChain.MessageChain_ = []*MessageChain{tailChain}
}

func isObjectOrArrayLiteralType(t *Type) bool {
	return t.objectFlags&(ObjectFlagsObjectLiteral|ObjectFlagsArrayLiteral) != 0
}

func getContainingClassExcludingClassDecorators(node *Node) *ClassLikeDeclaration {
	decorator := findAncestorOrQuit(node.Parent, func(n *Node) FindAncestorResult {
		if isClassLike(n) {
			return FindAncestorQuit
		}
		if IsDecorator(n) {
			return FindAncestorTrue
		}
		return FindAncestorFalse
	})
	if decorator != nil && isClassLike(decorator.Parent) {
		return getContainingClass(decorator.Parent)
	}
	if decorator != nil {
		return getContainingClass(decorator)
	}
	return getContainingClass(node)
}

func isThisTypeParameter(t *Type) bool {
	return t.flags&TypeFlagsTypeParameter != 0 && t.AsTypeParameter().isThisType
}

func isCallLikeExpression(node *Node) bool {
	switch node.Kind {
	case ast.KindJsxOpeningElement, ast.KindJsxSelfClosingElement, ast.KindCallExpression, ast.KindNewExpression,
		ast.KindTaggedTemplateExpression, ast.KindDecorator:
		return true
	}
	return false
}

func isCallOrNewExpression(node *Node) bool {
	return IsCallExpression(node) || IsNewExpression(node)
}

func isClassInstanceProperty(node *Node) bool {
	return node.Parent != nil && isClassLike(node.Parent) && IsPropertyDeclaration(node) && !hasAccessorModifier(node)
}

func isThisInitializedObjectBindingExpression(node *Node) bool {
	return node != nil && (IsShorthandPropertyAssignment(node) || IsPropertyAssignment(node)) && isBinaryExpression(node.Parent.Parent) &&
		node.Parent.Parent.AsBinaryExpression().OperatorToken.Kind == ast.KindEqualsToken &&
		node.Parent.Parent.AsBinaryExpression().Right.Kind == ast.KindThisKeyword
}

func isThisInitializedDeclaration(node *Node) bool {
	return node != nil && IsVariableDeclaration(node) && node.AsVariableDeclaration().Initializer != nil && node.AsVariableDeclaration().Initializer.Kind == ast.KindThisKeyword
}

func isWriteOnlyAccess(node *Node) bool {
	return accessKind(node) == AccessKindWrite
}

func isWriteAccess(node *Node) bool {
	return accessKind(node) != AccessKindRead
}

type AccessKind int32

const (
	AccessKindRead      AccessKind = iota // Only reads from a variable
	AccessKindWrite                       // Only writes to a variable without ever reading it. E.g.: `x=1;`.
	AccessKindReadWrite                   // Reads from and writes to a variable. E.g.: `f(x++);`, `x/=1`.
)

func accessKind(node *Node) AccessKind {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindParenthesizedExpression:
		return accessKind(parent)
	case ast.KindPrefixUnaryExpression:
		operator := parent.AsPrefixUnaryExpression().Operator
		if operator == ast.KindPlusPlusToken || operator == ast.KindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case ast.KindPostfixUnaryExpression:
		operator := parent.AsPostfixUnaryExpression().Operator
		if operator == ast.KindPlusPlusToken || operator == ast.KindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case ast.KindBinaryExpression:
		if parent.AsBinaryExpression().Left == node {
			operator := parent.AsBinaryExpression().OperatorToken
			if isAssignmentOperator(operator.Kind) {
				if operator.Kind == ast.KindEqualsToken {
					return AccessKindWrite
				}
				return AccessKindReadWrite
			}
		}
		return AccessKindRead
	case ast.KindPropertyAccessExpression:
		if parent.AsPropertyAccessExpression().Name_ != node {
			return AccessKindRead
		}
		return accessKind(parent)
	case ast.KindPropertyAssignment:
		parentAccess := accessKind(parent.Parent)
		// In `({ x: varname }) = { x: 1 }`, the left `x` is a read, the right `x` is a write.
		if node == parent.AsPropertyAssignment().Name_ {
			return reverseAccessKind(parentAccess)
		}
		return parentAccess
	case ast.KindShorthandPropertyAssignment:
		// Assume it's the local variable being accessed, since we don't check public properties for --noUnusedLocals.
		if node == parent.AsShorthandPropertyAssignment().ObjectAssignmentInitializer {
			return AccessKindRead
		}
		return accessKind(parent.Parent)
	case ast.KindArrayLiteralExpression:
		return accessKind(parent)
	case ast.KindForInStatement, ast.KindForOfStatement:
		if node == parent.AsForInOrOfStatement().Initializer {
			return AccessKindWrite
		}
		return AccessKindRead
	}
	return AccessKindRead
}

func reverseAccessKind(a AccessKind) AccessKind {
	switch a {
	case AccessKindRead:
		return AccessKindWrite
	case AccessKindWrite:
		return AccessKindRead
	case AccessKindReadWrite:
		return AccessKindReadWrite
	}
	panic("Unhandled case in reverseAccessKind")
}

func isJsxOpeningLikeElement(node *Node) bool {
	return IsJsxOpeningElement(node) || IsJsxSelfClosingElement(node)
}

func isObjectLiteralElementLike(node *Node) bool {
	switch node.Kind {
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment, ast.KindSpreadAssignment,
		ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor:
		return true
	}
	return false
}
