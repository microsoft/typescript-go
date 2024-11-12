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
	"github.com/microsoft/typescript-go/internal/compiler/textpos"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// TextRange

type TextRange struct {
	pos textpos.TextPos
	end textpos.TextPos
}

func NewTextRange(pos int, end int) TextRange {
	return TextRange{pos: textpos.TextPos(pos), end: textpos.TextPos(end)}
}

func (t TextRange) Pos() int {
	return int(t.pos)
}

func (t TextRange) End() int {
	return int(t.end)
}

func (t TextRange) Len() int {
	return int(t.end - t.pos)
}

func (t TextRange) ContainsInclusive(pos int) bool {
	return pos >= int(t.pos) && pos <= int(t.end)
}

// Pool allocator

type Pool[T any] struct {
	data []T
}

func (p *Pool[T]) New() *T {
	if len(p.data) == cap(p.data) {
		p.data = make([]T, 0, nextPoolSize(len(p.data)))
	}
	index := len(p.data)
	p.data = p.data[:index+1]
	return &p.data[index]
}

// Links store

type LinkStore[K comparable, V any] struct {
	entries map[K]*V
	pool    Pool[V]
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
	if node.id == 0 {
		node.id = NodeId(nextNodeId.Add(1))
	}
	return node.id
}

func getSymbolId(symbol *Symbol) SymbolId {
	if symbol.id == 0 {
		symbol.id = SymbolId(nextSymbolId.Add(1))
	}
	return symbol.id
}

func getMergeId(symbol *Symbol) MergeId {
	if symbol.mergeId == 0 {
		symbol.mergeId = MergeId(nextMergeId.Add(1))
	}
	return symbol.mergeId
}

// Diagnostic

type Diagnostic struct {
	file               *SourceFile
	loc                TextRange
	code               int32
	category           diagnostics.Category
	message            string
	messageChain       []*MessageChain
	relatedInformation []*Diagnostic
}

func (d *Diagnostic) File() *SourceFile                 { return d.file }
func (d *Diagnostic) Pos() int                          { return d.loc.Pos() }
func (d *Diagnostic) End() int                          { return d.loc.End() }
func (d *Diagnostic) Len() int                          { return d.loc.Len() }
func (d *Diagnostic) Loc() TextRange                    { return d.loc }
func (d *Diagnostic) Code() int32                       { return d.code }
func (d *Diagnostic) Category() diagnostics.Category    { return d.category }
func (d *Diagnostic) Message() string                   { return d.message }
func (d *Diagnostic) MessageChain() []*MessageChain     { return d.messageChain }
func (d *Diagnostic) RelatedInformation() []*Diagnostic { return d.relatedInformation }

func (d *Diagnostic) SetCategory(category diagnostics.Category) { d.category = category }

func NewDiagnostic(file *SourceFile, loc TextRange, message *diagnostics.Message, args ...any) *Diagnostic {
	text := message.Message()
	if len(args) != 0 {
		text = formatStringFromArgs(text, args)
	}
	return &Diagnostic{
		file:     file,
		loc:      loc,
		code:     message.Code(),
		category: message.Category(),
		message:  text,
	}
}

func NewDiagnosticForNode(node *Node, message *diagnostics.Message, args ...any) *Diagnostic {
	var file *SourceFile
	var loc TextRange
	if node != nil {
		file = getSourceFileOfNode(node)
		loc = getErrorRangeForNode(file, node)
	}
	return NewDiagnostic(file, loc, message, args...)
}

func NewDiagnosticFromMessageChain(file *SourceFile, loc TextRange, messageChain *MessageChain) *Diagnostic {
	return &Diagnostic{
		file:         file,
		loc:          loc,
		code:         messageChain.code,
		category:     messageChain.category,
		message:      messageChain.message,
		messageChain: messageChain.messageChain,
	}
}

func NewDiagnosticForNodeFromMessageChain(node *Node, messageChain *MessageChain) *Diagnostic {
	var file *SourceFile
	var loc TextRange
	if node != nil {
		file = getSourceFileOfNode(node)
		loc = getErrorRangeForNode(file, node)
	}
	return NewDiagnosticFromMessageChain(file, loc, messageChain)
}

func (d *Diagnostic) setMessageChain(messageChain []*MessageChain) *Diagnostic {
	d.messageChain = messageChain
	return d
}

func (d *Diagnostic) addMessageChain(messageChain *MessageChain) *Diagnostic {
	if messageChain != nil {
		d.messageChain = append(d.messageChain, messageChain)
	}
	return d
}

func (d *Diagnostic) setRelatedInfo(relatedInformation []*Diagnostic) *Diagnostic {
	d.relatedInformation = relatedInformation
	return d
}

func (d *Diagnostic) addRelatedInfo(relatedInformation *Diagnostic) *Diagnostic {
	if relatedInformation != nil {
		d.relatedInformation = append(d.relatedInformation, relatedInformation)
	}
	return d
}

// MessageChain

type MessageChain struct {
	code         int32
	category     diagnostics.Category
	message      string
	messageChain []*MessageChain
}

func NewMessageChain(message *diagnostics.Message, args ...any) *MessageChain {
	text := message.Message()
	if len(args) != 0 {
		text = formatStringFromArgs(text, args)
	}
	return &MessageChain{
		code:     message.Code(),
		category: message.Category(),
		message:  text,
	}
}

func (m *MessageChain) Code() int32                    { return m.code }
func (m *MessageChain) Category() diagnostics.Category { return m.category }
func (m *MessageChain) Message() string                { return m.message }
func (m *MessageChain) MessageChain() []*MessageChain  { return m.messageChain }

func (m *MessageChain) addMessageChain(messageChain *MessageChain) *MessageChain {
	if messageChain != nil {
		m.messageChain = append(m.messageChain, messageChain)
	}
	return m
}

func chainDiagnosticMessages(details *MessageChain, message *diagnostics.Message, args ...any) *MessageChain {
	return NewMessageChain(message, args...).addMessageChain(details)
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

func modifierToFlag(token ast.Kind) ModifierFlags {
	switch token {
	case ast.KindStaticKeyword:
		return ModifierFlagsStatic
	case ast.KindPublicKeyword:
		return ModifierFlagsPublic
	case ast.KindProtectedKeyword:
		return ModifierFlagsProtected
	case ast.KindPrivateKeyword:
		return ModifierFlagsPrivate
	case ast.KindAbstractKeyword:
		return ModifierFlagsAbstract
	case ast.KindAccessorKeyword:
		return ModifierFlagsAccessor
	case ast.KindExportKeyword:
		return ModifierFlagsExport
	case ast.KindDeclareKeyword:
		return ModifierFlagsAmbient
	case ast.KindConstKeyword:
		return ModifierFlagsConst
	case ast.KindDefaultKeyword:
		return ModifierFlagsDefault
	case ast.KindAsyncKeyword:
		return ModifierFlagsAsync
	case ast.KindReadonlyKeyword:
		return ModifierFlagsReadonly
	case ast.KindOverrideKeyword:
		return ModifierFlagsOverride
	case ast.KindInKeyword:
		return ModifierFlagsIn
	case ast.KindOutKeyword:
		return ModifierFlagsOut
	case ast.KindImmediateKeyword:
		return ModifierFlagsImmediate
	case ast.KindDecorator:
		return ModifierFlagsDecorator
	}
	return ModifierFlagsNone
}

func modifiersToFlags(modifierList *Node) ModifierFlags {
	flags := ModifierFlagsNone
	if modifierList != nil {
		for _, modifier := range modifierList.AsModifierList().modifiers {
			flags |= modifierToFlag(modifier.kind)
		}
	}
	return flags
}

func nodeIsMissing(node *Node) bool {
	return node == nil || node.loc.pos == node.loc.end && node.loc.pos >= 0 && node.kind != ast.KindEndOfFile
}

func nodeIsPresent(node *Node) bool {
	return !nodeIsMissing(node)
}

func isLeftHandSideExpression(node *Node) bool {
	return isLeftHandSideExpressionKind(node.kind)
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
	return isUnaryExpressionKind(node.kind)
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
	return isExpressionKind(node.kind)
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
	return node.kind == ast.KindExpressionWithTypeArguments
}

func isNonNullExpression(node *Node) bool {
	return node.kind == ast.KindNonNullExpression
}

func isStringLiteralLike(node *Node) bool {
	return node.kind == ast.KindStringLiteral || node.kind == ast.KindNoSubstitutionTemplateLiteral
}

func isNumericLiteral(node *Node) bool {
	return node.kind == ast.KindNumericLiteral
}

func isStringOrNumericLiteralLike(node *Node) bool {
	return isStringLiteralLike(node) || isNumericLiteral(node)
}

func isSignedNumericLiteral(node *Node) bool {
	if node.kind == ast.KindPrefixUnaryExpression {
		node := node.AsPrefixUnaryExpression()
		return (node.operator == ast.KindPlusToken || node.operator == ast.KindMinusToken) && isNumericLiteral(node.operand)
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
	return getTextOfNodeFromSourceText(sourceFile.text, node)
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
	return isBinaryExpression(decl) || isAccessExpression(decl) || isIdentifier(decl) || isCallExpression(decl)
}

func isBinaryExpression(node *Node) bool {
	return node.kind == ast.KindBinaryExpression
}

func isAccessExpression(node *Node) bool {
	return node.kind == ast.KindPropertyAccessExpression || node.kind == ast.KindElementAccessExpression
}

func isInJSFile(node *Node) bool {
	return node != nil && node.flags&ast.NodeFlagsJavaScriptFile != 0
}

func isEffectiveModuleDeclaration(node *Node) bool {
	return isModuleDeclaration(node) || isIdentifier(node)
}

func isObjectLiteralOrClassExpressionMethodOrAccessor(node *Node) bool {
	kind := node.kind
	return (kind == ast.KindMethodDeclaration || kind == ast.KindGetAccessor || kind == ast.KindSetAccessor) &&
		(node.parent.kind == ast.KindObjectLiteralExpression || node.parent.kind == ast.KindClassExpression)
}

func isFunctionLike(node *Node) bool {
	return node != nil && isFunctionLikeKind(node.kind)
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
	return node != nil && isFunctionLikeDeclarationKind(node.kind)
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
	switch node.kind {
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
	for node != nil && node.kind == ast.KindParenthesizedType {
		node = node.parent
	}
	return node
}

func walkUpParenthesizedExpressions(node *Node) *Node {
	for node != nil && node.kind == ast.KindParenthesizedExpression {
		node = node.parent
	}
	return node
}

func isJSDocTypeAssertion(node *Node) bool {
	return false // !!!
}

// Return true if the given identifier is classified as an IdentifierName
func isIdentifierName(node *Node) bool {
	parent := node.parent
	switch parent.kind {
	case ast.KindPropertyDeclaration, ast.KindPropertySignature, ast.KindMethodDeclaration, ast.KindMethodSignature, ast.KindGetAccessor,
		ast.KindSetAccessor, ast.KindEnumMember, ast.KindPropertyAssignment, ast.KindPropertyAccessExpression:
		return parent.Name() == node
	case ast.KindQualifiedName:
		return parent.AsQualifiedName().right == node
	case ast.KindBindingElement:
		return parent.AsBindingElement().propertyName == node
	case ast.KindImportSpecifier:
		return parent.AsImportSpecifier().propertyName == node
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
		if node.kind == ast.KindSourceFile {
			return node.data.(*SourceFile)
		}
		node = node.parent
	}
}

/** @internal */
func getErrorRangeForNode(sourceFile *SourceFile, node *Node) TextRange {
	errorNode := node
	switch node.kind {
	case ast.KindSourceFile:
		pos := skipTrivia(sourceFile.text, 0)
		if pos == len(sourceFile.text) {
			return NewTextRange(0, 0)
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
		start := skipTrivia(sourceFile.text, node.Pos())
		end := node.End()
		statements := node.data.(*CaseOrDefaultClause).statements
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return NewTextRange(start, end)
	case ast.KindReturnStatement, ast.KindYieldExpression:
		pos := skipTrivia(sourceFile.text, node.Pos())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindSatisfiesExpression:
		pos := skipTrivia(sourceFile.text, node.AsSatisfiesExpression().expression.End())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindConstructor:
		scanner := getScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.tokenStart
		for scanner.token != ast.KindConstructorKeyword && scanner.token != ast.KindStringLiteral && scanner.token != ast.KindEndOfFile {
			scanner.Scan()
		}
		return NewTextRange(start, scanner.pos)
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
		pos = skipTrivia(sourceFile.text, pos)
	}
	return NewTextRange(pos, errorNode.End())
}

func getErrorRangeForArrowFunction(sourceFile *SourceFile, node *Node) TextRange {
	pos := skipTrivia(sourceFile.text, node.Pos())
	body := node.AsArrowFunction().body
	if body != nil && body.kind == ast.KindBlock {
		startLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.Pos())
		endLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.End())
		if startLine < endLine {
			// The arrow function spans multiple lines,
			// make the error span be the first line, inclusive.
			return NewTextRange(pos, getEndLinePosition(sourceFile, startLine))
		}
	}
	return NewTextRange(pos, node.End())
}

func getContainingClass(node *Node) *Node {
	return findAncestor(node.parent, isClassLike)
}

func findAncestor(node *Node, callback func(*Node) bool) *Node {
	for node != nil {
		result := callback(node)
		if result {
			return node
		}
		node = node.parent
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
		node = node.parent
	}
	return nil
}

func isClassLike(node *Node) bool {
	return node != nil && (node.kind == ast.KindClassDeclaration || node.kind == ast.KindClassExpression)
}

func declarationNameToString(name *Node) string {
	if name == nil || name.Pos() == name.End() {
		return "(Missing)"
	}
	return getTextOfNode(name)
}

func isExternalModule(file *SourceFile) bool {
	return file.externalModuleIndicator != nil
}

func isInTopLevelContext(node *Node) bool {
	// The name of a class or function declaration is a BindingIdentifier in its surrounding scope.
	if isIdentifier(node) {
		parent := node.parent
		if (isClassDeclaration(parent) || isFunctionDeclaration(parent)) && parent.Name() == node {
			node = parent
		}
	}
	container := getThisContainer(node, true /*includeArrowFunctions*/, false /*includeClassComputedPropertyName*/)
	return isSourceFile(container)
}

func getThisContainer(node *Node, includeArrowFunctions bool, includeClassComputedPropertyName bool) *Node {
	for {
		node = node.parent
		if node == nil {
			panic("nil parent in getThisContainer")
		}
		switch node.kind {
		case ast.KindComputedPropertyName:
			if includeClassComputedPropertyName && isClassLike(node.parent.parent) {
				return node
			}
			node = node.parent.parent
		case ast.KindDecorator:
			if node.parent.kind == ast.KindParameter && isClassElement(node.parent.parent) {
				// If the decorator's parent is a Parameter, we resolve the this container from
				// the grandparent class declaration.
				node = node.parent.parent
			} else if isClassElement(node.parent) {
				// If the decorator's parent is a class element, we resolve the 'this' container
				// from the parent class declaration.
				node = node.parent
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
	switch node.kind {
	case ast.KindConstructor, ast.KindPropertyDeclaration, ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor,
		ast.KindIndexSignature, ast.KindClassStaticBlockDeclaration, ast.KindSemicolonClassElement:
		return true
	}
	return false
}

func isPartOfTypeQuery(node *Node) bool {
	for node.kind == ast.KindQualifiedName || node.kind == ast.KindIdentifier {
		node = node.parent
	}
	return node.kind == ast.KindTypeQuery
}

func getModifierFlags(node *Node) ModifierFlags {
	modifiers := node.Modifiers()
	if modifiers != nil {
		return modifiers.AsModifierList().modifierFlags
	}
	return ModifierFlagsNone
}

func getNodeFlags(node *Node) ast.NodeFlags {
	return node.flags
}

func hasSyntacticModifier(node *Node, flags ModifierFlags) bool {
	return getModifierFlags(node)&flags != 0
}

func hasAccessorModifier(node *Node) bool {
	return hasSyntacticModifier(node, ModifierFlagsAccessor)
}

func hasStaticModifier(node *Node) bool {
	return hasSyntacticModifier(node, ModifierFlagsStatic)
}

func getEffectiveModifierFlags(node *Node) ModifierFlags {
	return getModifierFlags(node) // !!! Handle JSDoc
}

func hasEffectiveModifier(node *Node, flags ModifierFlags) bool {
	return getEffectiveModifierFlags(node)&flags != 0
}

func hasEffectiveReadonlyModifier(node *Node) bool {
	return hasEffectiveModifier(node, ModifierFlagsReadonly)
}

func getImmediatelyInvokedFunctionExpression(fn *Node) *Node {
	if fn.kind == ast.KindFunctionExpression || fn.kind == ast.KindArrowFunction {
		prev := fn
		parent := fn.parent
		for parent.kind == ast.KindParenthesizedExpression {
			prev = parent
			parent = parent.parent
		}
		if parent.kind == ast.KindCallExpression && parent.AsCallExpression().expression == prev {
			return parent
		}
	}
	return nil
}

// Does not handle signed numeric names like `a[+0]` - handling those would require handling prefix unary expressions
// throughout late binding handling as well, which is awkward (but ultimately probably doable if there is demand)
func getElementOrPropertyAccessArgumentExpressionOrName(node *Node) *Node {
	switch node.kind {
	case ast.KindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().name
	case ast.KindElementAccessExpression:
		arg := skipParentheses(node.AsElementAccessExpression().argumentExpression)
		if isStringOrNumericLiteralLike(arg) {
			return arg
		}
		return node
	}
	panic("Unhandled case in getElementOrPropertyAccessArgumentExpressionOrName")
}

func getQuestionDotToken(node *Node) *Node {
	switch node.kind {
	case ast.KindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().questionDotToken
	case ast.KindElementAccessExpression:
		return node.AsElementAccessExpression().questionDotToken
	case ast.KindCallExpression:
		return node.AsCallExpression().questionDotToken
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
	switch name.kind {
	case ast.KindComputedPropertyName:
		expr = name.AsComputedPropertyName().expression
	case ast.KindElementAccessExpression:
		expr = skipParentheses(name.AsElementAccessExpression().argumentExpression)
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
	if isFunctionExpression(declaration) || isArrowFunction(declaration) || isClassExpression(declaration) {
		return getAssignedName(declaration)
	}
	return nil
}

func getNonAssignedNameOfDeclaration(declaration *Node) *Node {
	switch declaration.kind {
	case ast.KindBinaryExpression:
		if isFunctionPropertyAssignment(declaration) {
			return getElementOrPropertyAccessArgumentExpressionOrName(declaration.AsBinaryExpression().left)
		}
		return nil
	case ast.KindExportAssignment:
		expr := declaration.AsExportAssignment().expression
		if isIdentifier(expr) {
			return expr
		}
		return nil
	}
	return declaration.Name()
}

func getAssignedName(node *Node) *Node {
	parent := node.parent
	if parent != nil {
		switch parent.kind {
		case ast.KindPropertyAssignment:
			return parent.AsPropertyAssignment().name
		case ast.KindBindingElement:
			return parent.AsBindingElement().name
		case ast.KindBinaryExpression:
			if node == parent.AsBinaryExpression().right {
				left := parent.AsBinaryExpression().left
				switch left.kind {
				case ast.KindIdentifier:
					return left
				case ast.KindPropertyAccessExpression:
					return left.AsPropertyAccessExpression().name
				case ast.KindElementAccessExpression:
					arg := skipParentheses(left.AsElementAccessExpression().argumentExpression)
					if isStringOrNumericLiteralLike(arg) {
						return arg
					}
				}
			}
		case ast.KindVariableDeclaration:
			name := parent.AsVariableDeclaration().name
			if isIdentifier(name) {
				return name
			}
		}
	}
	return nil
}

func isFunctionPropertyAssignment(node *Node) bool {
	if node.kind == ast.KindBinaryExpression {
		expr := node.AsBinaryExpression()
		if expr.operatorToken.kind == ast.KindEqualsToken {
			switch expr.left.kind {
			case ast.KindPropertyAccessExpression:
				// F.id = expr
				return isIdentifier(expr.left.AsPropertyAccessExpression().expression) && isIdentifier(expr.left.AsPropertyAccessExpression().name)
			case ast.KindElementAccessExpression:
				// F[xxx] = expr
				return isIdentifier(expr.left.AsElementAccessExpression().expression)
			}
		}
	}
	return false
}

func isAssignmentExpression(node *Node, excludeCompoundAssignment bool) bool {
	if node.kind == ast.KindBinaryExpression {
		expr := node.AsBinaryExpression()
		return (expr.operatorToken.kind == ast.KindEqualsToken || !excludeCompoundAssignment && isAssignmentOperator(expr.operatorToken.kind)) &&
			isLeftHandSideExpression(expr.left)
	}
	return false
}

func isBlockOrCatchScoped(declaration *Node) bool {
	return getCombinedNodeFlags(declaration)&ast.NodeFlagsBlockScoped != 0 || isCatchClauseVariableDeclarationOrBindingElement(declaration)
}

func isCatchClauseVariableDeclarationOrBindingElement(declaration *Node) bool {
	node := getRootDeclaration(declaration)
	return node.kind == ast.KindVariableDeclaration && node.parent.kind == ast.KindCatchClause
}

func isAmbientModule(node *Node) bool {
	return isModuleDeclaration(node) && (node.AsModuleDeclaration().name.kind == ast.KindStringLiteral || isGlobalScopeAugmentation(node))
}

func isGlobalScopeAugmentation(node *Node) bool {
	return node.flags&ast.NodeFlagsGlobalAugmentation != 0
}

func isPropertyNameLiteral(node *Node) bool {
	switch node.kind {
	case ast.KindIdentifier, ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral, ast.KindNumericLiteral:
		return true
	}
	return false
}

func isMemberName(node *Node) bool {
	return node.kind == ast.KindIdentifier || node.kind == ast.KindPrivateIdentifier
}

func setParent(child *Node, parent *Node) {
	if child != nil {
		child.parent = parent
	}
}

func setParentInChildren(node *Node) {
	node.ForEachChild(func(child *Node) bool {
		child.parent = node
		setParentInChildren(child)
		return false
	})
}

func getCombinedFlags[T ~uint32](node *Node, getFlags func(*Node) T) T {
	node = getRootDeclaration(node)
	flags := getFlags(node)
	if node.kind == ast.KindVariableDeclaration {
		node = node.parent
	}
	if node != nil && node.kind == ast.KindVariableDeclarationList {
		flags |= getFlags(node)
		node = node.parent
	}
	if node != nil && node.kind == ast.KindVariableStatement {
		flags |= getFlags(node)
	}
	return flags
}

func getCombinedModifierFlags(node *Node) ModifierFlags {
	return getCombinedFlags(node, getModifierFlags)
}

func getCombinedNodeFlags(node *Node) ast.NodeFlags {
	return getCombinedFlags(node, getNodeFlags)
}

func isBindingPattern(node *Node) bool {
	return node != nil && (node.kind == ast.KindArrayBindingPattern || node.kind == ast.KindObjectBindingPattern)
}

func isParameterPropertyDeclaration(node *Node, parent *Node) bool {
	return isParameter(node) && hasSyntacticModifier(node, ModifierFlagsParameterPropertyModifier) && parent.kind == ast.KindConstructor
}

/**
 * Like {@link isVariableDeclarationInitializedToRequire} but allows things like `require("...").foo.bar` or `require("...")["baz"]`.
 */
func isVariableDeclarationInitializedToBareOrAccessedRequire(node *Node) bool {
	return isVariableDeclarationInitializedWithRequireHelper(node, true /*allowAccessedRequire*/)
}

func isVariableDeclarationInitializedWithRequireHelper(node *Node, allowAccessedRequire bool) bool {
	if node.kind == ast.KindVariableDeclaration && node.AsVariableDeclaration().initializer != nil {
		initializer := node.AsVariableDeclaration().initializer
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
	if isCallExpression(node) {
		callExpression := node.AsCallExpression()
		if len(callExpression.arguments) == 1 {
			if isIdentifier(callExpression.expression) && callExpression.expression.AsIdentifier().text == "require" {
				return !requireStringLiteralLikeArgument || isStringLiteralLike(callExpression.arguments[0])
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
	return getRootDeclaration(node).kind == ast.KindParameter
}

func getRootDeclaration(node *Node) *Node {
	for node.kind == ast.KindBindingElement {
		node = node.parent.parent
	}
	return node
}

func isExternalOrCommonJsModule(file *SourceFile) bool {
	return file.externalModuleIndicator != nil
}

func isAutoAccessorPropertyDeclaration(node *Node) bool {
	return isPropertyDeclaration(node) && hasAccessorModifier(node)
}

func isAsyncFunction(node *Node) bool {
	switch node.kind {
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindMethodDeclaration:
		data := node.BodyData()
		return data.body != nil && data.asteriskToken == nil && hasSyntacticModifier(node, ModifierFlagsAsync)
	}
	return false
}

func isObjectLiteralMethod(node *Node) bool {
	return node != nil && node.kind == ast.KindMethodDeclaration && node.parent.kind == ast.KindObjectLiteralExpression
}

func symbolName(symbol *Symbol) string {
	if symbol.valueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(symbol.valueDeclaration) {
		return symbol.valueDeclaration.Name().AsPrivateIdentifier().text
	}
	return symbol.name
}

func isStaticPrivateIdentifierProperty(s *Symbol) bool {
	return s.valueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(s.valueDeclaration) && isStatic(s.valueDeclaration)
}

func isPrivateIdentifierClassElementDeclaration(node *Node) bool {
	return (isPropertyDeclaration(node) || isMethodOrAccessor(node)) && isPrivateIdentifier(node.Name())
}

func isMethodOrAccessor(node *Node) bool {
	switch node.kind {
	case ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor:
		return true
	}
	return false
}

func isFunctionLikeOrClassStaticBlockDeclaration(node *Node) bool {
	return node != nil && (isFunctionLikeKind(node.kind) || isClassStaticBlockDeclaration(node))
}

func isModuleAugmentationExternal(node *Node) bool {
	// external module augmentation is a ambient module declaration that is either:
	// - defined in the top level scope and source file is an external module
	// - defined inside ambient module declaration located in the top level scope and source file not an external module
	switch node.parent.kind {
	case ast.KindSourceFile:
		return isExternalModule(node.parent.AsSourceFile())
	case ast.KindModuleBlock:
		grandParent := node.parent.parent
		return isAmbientModule(grandParent) && isSourceFile(grandParent.parent) && !isExternalModule(grandParent.parent.AsSourceFile())
	}
	return false
}

type Pattern struct {
	text      string
	starIndex int // -1 for exact match
}

func isValidPattern(pattern Pattern) bool {
	return pattern.starIndex == -1 || pattern.starIndex < len(pattern.text)
}

func tryParsePattern(pattern string) Pattern {
	starIndex := strings.Index(pattern, "*")
	if starIndex == -1 || !strings.Contains(pattern[starIndex+1:], "*") {
		return Pattern{text: pattern, starIndex: starIndex}
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
	return isDeclarationStatementKind(node.kind)
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
	return isStatementKindButNotDeclarationKind(node.kind)
}

func isStatement(node *Node) bool {
	kind := node.kind
	return isStatementKindButNotDeclarationKind(kind) || isDeclarationStatementKind(kind) || isBlockStatement(node)
}

func isBlockStatement(node *Node) bool {
	if node.kind != ast.KindBlock {
		return false
	}
	if node.parent != nil && (node.parent.kind == ast.KindTryStatement || node.parent.kind == ast.KindCatchClause) {
		return false
	}
	return !isFunctionBlock(node)
}

func isFunctionBlock(node *Node) bool {
	return node != nil && node.kind == ast.KindBlock && isFunctionLike(node.parent)
}

func shouldPreserveConstEnums(options *CompilerOptions) bool {
	return options.PreserveConstEnums == core.TSTrue || options.IsolatedModules == core.TSTrue
}

func exportAssignmentIsAlias(node *Node) bool {
	return isAliasableExpression(getExportAssignmentExpression(node))
}

func getExportAssignmentExpression(node *Node) *Node {
	switch node.kind {
	case ast.KindExportAssignment:
		return node.AsExportAssignment().expression
	case ast.KindBinaryExpression:
		return node.AsBinaryExpression().right
	}
	panic("Unhandled case in getExportAssignmentExpression")
}

func isAliasableExpression(e *Node) bool {
	return isEntityNameExpression(e) || isClassExpression(e)
}

func isEmptyObjectLiteral(expression *Node) bool {
	return expression.kind == ast.KindObjectLiteralExpression && len(expression.AsObjectLiteralExpression().properties) == 0
}

func isFunctionSymbol(symbol *Symbol) bool {
	d := symbol.valueDeclaration
	return d != nil && (isFunctionDeclaration(d) || isVariableDeclaration(d) && isFunctionLike(d.AsVariableDeclaration().initializer))
}

func isLogicalOrCoalescingAssignmentOperator(token ast.Kind) bool {
	return token == ast.KindBarBarEqualsToken || token == ast.KindAmpersandAmpersandEqualsToken || token == ast.KindQuestionQuestionEqualsToken
}

func isLogicalOrCoalescingAssignmentExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingAssignmentOperator(expr.AsBinaryExpression().operatorToken.kind)
}

func isLogicalOrCoalescingBinaryOperator(token ast.Kind) bool {
	return isBinaryLogicalOperator(token) || token == ast.KindQuestionQuestionToken
}

func isLogicalOrCoalescingBinaryExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingBinaryOperator(expr.AsBinaryExpression().operatorToken.kind)
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
	parent := node.parent
	return !isOptionalChain(parent) || // cases 1, 2, and 3
		isOptionalChainRoot(parent) || // case 4
		node != parent.Expression() // case 5
}

func isNullishCoalesce(node *Node) bool {
	return node.kind == ast.KindBinaryExpression && node.AsBinaryExpression().operatorToken.kind == ast.KindQuestionQuestionToken
}

func isDottedName(node *Node) bool {
	switch node.kind {
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
		kind := node.AsBinaryExpression().left.kind
		return kind == ast.KindObjectLiteralExpression || kind == ast.KindArrayLiteralExpression
	}
	return false
}

func isTopLevelLogicalExpression(node *Node) bool {
	for isParenthesizedExpression(node.parent) || isPrefixUnaryExpression(node.parent) && node.parent.AsPrefixUnaryExpression().operator == ast.KindExclamationToken {
		node = node.parent
	}
	return !isStatementCondition(node) && !isLogicalExpression(node.parent) && !(isOptionalChain(node.parent) && node.parent.Expression() == node)
}

func isStatementCondition(node *Node) bool {
	switch node.parent.kind {
	case ast.KindIfStatement:
		return node.parent.AsIfStatement().expression == node
	case ast.KindWhileStatement:
		return node.parent.AsWhileStatement().expression == node
	case ast.KindDoStatement:
		return node.parent.AsDoStatement().expression == node
	case ast.KindForStatement:
		return node.parent.AsForStatement().condition == node
	case ast.KindConditionalExpression:
		return node.parent.AsConditionalExpression().condition == node
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
	switch target.kind {
	case ast.KindBinaryExpression:
		binaryOperator := target.AsBinaryExpression().operatorToken.kind
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
		parent := node.parent
		switch parent.kind {
		case ast.KindBinaryExpression:
			if isAssignmentOperator(parent.AsBinaryExpression().operatorToken.kind) && parent.AsBinaryExpression().left == node {
				return parent
			}
			return nil
		case ast.KindPrefixUnaryExpression:
			if parent.AsPrefixUnaryExpression().operator == ast.KindPlusPlusToken || parent.AsPrefixUnaryExpression().operator == ast.KindMinusMinusToken {
				return parent
			}
			return nil
		case ast.KindPostfixUnaryExpression:
			if parent.AsPostfixUnaryExpression().operator == ast.KindPlusPlusToken || parent.AsPostfixUnaryExpression().operator == ast.KindMinusMinusToken {
				return parent
			}
			return nil
		case ast.KindForInStatement, ast.KindForOfStatement:
			if parent.AsForInOrOfStatement().initializer == node {
				return parent
			}
			return nil
		case ast.KindParenthesizedExpression, ast.KindArrayLiteralExpression, ast.KindSpreadElement, ast.KindNonNullExpression:
			node = parent
		case ast.KindSpreadAssignment:
			node = parent.parent
		case ast.KindShorthandPropertyAssignment:
			if parent.AsShorthandPropertyAssignment().name != node {
				return nil
			}
			node = parent.parent
		case ast.KindPropertyAssignment:
			if parent.AsPropertyAssignment().name == node {
				return nil
			}
			node = parent.parent
		default:
			return nil
		}
	}
}

func isDeleteTarget(node *Node) bool {
	if !isAccessExpression(node) {
		return false
	}
	node = walkUpParenthesizedExpressions(node.parent)
	return node != nil && node.kind == ast.KindDeleteExpression
}

func isInCompoundLikeAssignment(node *Node) bool {
	target := getAssignmentTarget(node)
	return target != nil && isAssignmentExpression(target /*excludeCompoundAssignment*/, true) && isCompoundLikeAssignment(target)
}

func isCompoundLikeAssignment(assignment *Node) bool {
	right := skipParentheses(assignment.AsBinaryExpression().right)
	return right.kind == ast.KindBinaryExpression && isShiftOperatorOrHigher(right.AsBinaryExpression().operatorToken.kind)
}

func isPushOrUnshiftIdentifier(node *Node) bool {
	text := node.AsIdentifier().text
	return text == "push" || text == "unshift"
}

func isBooleanLiteral(node *Node) bool {
	return node.kind == ast.KindTrueKeyword || node.kind == ast.KindFalseKeyword
}

func isOptionalChain(node *Node) bool {
	kind := node.kind
	return node.flags&ast.NodeFlagsOptionalChain != 0 && (kind == ast.KindPropertyAccessExpression ||
		kind == ast.KindElementAccessExpression || kind == ast.KindCallExpression || kind == ast.KindNonNullExpression)
}

func isOptionalChainRoot(node *Node) bool {
	return isOptionalChain(node) && !isNonNullExpression(node) && getQuestionDotToken(node) != nil
}

/**
 * Determines whether a node is the expression preceding an optional chain (i.e. `a` in `a?.b`).
 */
func isExpressionOfOptionalChainRoot(node *Node) bool {
	return isOptionalChainRoot(node.parent) && node.parent.Expression() == node
}

func isEntityNameExpression(node *Node) bool {
	return node.kind == ast.KindIdentifier || isPropertyAccessEntityNameExpression(node)
}

func isPropertyAccessEntityNameExpression(node *Node) bool {
	if node.kind == ast.KindPropertyAccessExpression {
		expr := node.AsPropertyAccessExpression()
		return expr.name.kind == ast.KindIdentifier && isEntityNameExpression(expr.expression)
	}
	return false
}

func isPrologueDirective(node *Node) bool {
	return node.kind == ast.KindExpressionStatement && node.AsExpressionStatement().expression.kind == ast.KindStringLiteral
}

func nextPoolSize(size int) int {
	switch {
	case size < 16:
		return 16
	case size < 256:
		return size * 2
	}
	return size
}

func getStatementsOfBlock(block *Node) []*Statement {
	switch block.kind {
	case ast.KindBlock:
		return block.AsBlock().statements
	case ast.KindModuleBlock:
		return block.AsModuleBlock().statements
	case ast.KindSourceFile:
		return block.AsSourceFile().statements
	}
	panic("Unhandled case in getStatementsOfBlock")
}

func nodeHasName(statement *Node, id *Node) bool {
	name := statement.Name()
	if name != nil {
		return isIdentifier(name) && name.AsIdentifier().text == id.AsIdentifier().text
	}
	if isVariableStatement(statement) {
		declarations := statement.AsVariableStatement().declarationList.AsVariableDeclarationList().declarations
		return core.Some(declarations, func(d *Node) bool { return nodeHasName(d, id) })
	}
	return false
}

func isImportMeta(node *Node) bool {
	if node.kind == ast.KindMetaProperty {
		return node.AsMetaProperty().keywordToken == ast.KindImportKeyword && node.AsMetaProperty().name.AsIdentifier().text == "meta"
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
	if diagnostic.file != nil {
		fileName := diagnostic.file.fileName
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
	if diagnostic.file != nil {
		diagnostics = c.fileDiagnostics[diagnostic.file.fileName]
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
		d1.loc == d2.loc &&
		d1.code == d2.code &&
		d1.message == d2.message &&
		slices.EqualFunc(d1.messageChain, d2.messageChain, equalMessageChain) &&
		slices.EqualFunc(d1.relatedInformation, d2.relatedInformation, equalDiagnostics)
}

func equalMessageChain(c1, c2 *MessageChain) bool {
	return c1.code == c2.code &&
		c1.message == c2.message &&
		slices.EqualFunc(c1.messageChain, c2.messageChain, equalMessageChain)
}

func compareDiagnostics(d1, d2 *Diagnostic) int {
	c := strings.Compare(getDiagnosticPath(d1), getDiagnosticPath(d2))
	if c != 0 {
		return c
	}
	c = int(d1.loc.pos) - int(d2.loc.pos)
	if c != 0 {
		return c
	}
	c = int(d1.loc.end) - int(d2.loc.end)
	if c != 0 {
		return c
	}
	c = int(d1.code) - int(d2.code)
	if c != 0 {
		return c
	}
	c = strings.Compare(d1.message, d2.message)
	if c != 0 {
		return c
	}
	c = compareMessageChainSize(d1.messageChain, d2.messageChain)
	if c != 0 {
		return c
	}
	c = compareMessageChainContent(d1.messageChain, d2.messageChain)
	if c != 0 {
		return c
	}
	return compareRelatedInfo(d1.relatedInformation, d2.relatedInformation)
}

func compareMessageChainSize(c1, c2 []*MessageChain) int {
	c := len(c2) - len(c1)
	if c != 0 {
		return c
	}
	for i := range c1 {
		c = compareMessageChainSize(c1[i].messageChain, c2[i].messageChain)
		if c != 0 {
			return c
		}
	}
	return 0
}

func compareMessageChainContent(c1, c2 []*MessageChain) int {
	for i := range c1 {
		c := strings.Compare(c1[i].message, c2[i].message)
		if c != 0 {
			return c
		}
		if c1[i].messageChain != nil {
			c = compareMessageChainContent(c1[i].messageChain, c2[i].messageChain)
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
	if d.file != nil {
		return d.file.path
	}
	return ""
}

func isConstAssertion(location *Node) bool {
	switch location.kind {
	case ast.KindAsExpression:
		return isConstTypeReference(location.AsAsExpression().typeNode)
	case ast.KindTypeAssertionExpression:
		return isConstTypeReference(location.AsTypeAssertion().typeNode)
	}
	return false
}

func isConstTypeReference(node *Node) bool {
	if node.kind == ast.KindTypeReference {
		ref := node.AsTypeReference()
		return ref.typeArguments == nil && isIdentifier(ref.typeName) && ref.typeName.AsIdentifier().text == "const"
	}
	return false
}

func isModuleOrEnumDeclaration(node *Node) bool {
	return node.kind == ast.KindModuleDeclaration || node.kind == ast.KindEnumDeclaration
}

func getLocalsOfNode(node *Node) SymbolTable {
	data := node.LocalsContainerData()
	if data != nil {
		return data.locals
	}
	return nil
}

func getBodyOfNode(node *Node) *Node {
	bodyData := node.BodyData()
	if bodyData != nil {
		return bodyData.body
	}
	return nil
}

func getFlowNodeOfNode(node *Node) *FlowNode {
	flowNodeData := node.FlowNodeData()
	if flowNodeData != nil {
		return flowNodeData.flowNode
	}
	return nil
}

func isGlobalSourceFile(node *Node) bool {
	return node.kind == ast.KindSourceFile && !isExternalOrCommonJsModule(node.AsSourceFile())
}

func isParameterLikeOrReturnTag(node *Node) bool {
	switch node.kind {
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
	return isTypeNodeKind(node.kind)
}

func getLocalSymbolForExportDefault(symbol *Symbol) *Symbol {
	if !isExportDefaultSymbol(symbol) || len(symbol.declarations) == 0 {
		return nil
	}
	for _, decl := range symbol.declarations {
		localSymbol := decl.LocalSymbol()
		if localSymbol != nil {
			return localSymbol
		}
	}
	return nil
}

func isExportDefaultSymbol(symbol *Symbol) bool {
	return symbol != nil && len(symbol.declarations) > 0 && hasSyntacticModifier(symbol.declarations[0], ModifierFlagsDefault)
}

func getDeclarationOfKind(symbol *Symbol, kind ast.Kind) *Node {
	for _, declaration := range symbol.declarations {
		if declaration.kind == kind {
			return declaration
		}
	}
	return nil
}

func getIsolatedModules(options *CompilerOptions) bool {
	return options.IsolatedModules == core.TSTrue || options.VerbatimModuleSyntax == core.TSTrue
}

func findConstructorDeclaration(node *Node) *Node {
	for _, member := range node.ClassLikeData().members {
		if isConstructorDeclaration(member) && nodeIsPresent(member.AsConstructorDeclaration().body) {
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
			location = location.parent
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
					if meaning&result.flags&ast.SymbolFlagsType != 0 && lastLocation.kind != ast.KindJSDoc {
						useResult = result.flags&ast.SymbolFlagsTypeParameter != 0 && (lastLocation.flags&ast.NodeFlagsSynthesized != 0 ||
							lastLocation == location.ReturnType() ||
							isParameterLikeOrReturnTag(lastLocation))
					}
					if meaning&result.flags&ast.SymbolFlagsVariable != 0 {
						// expression inside parameter will lookup as normal variable scope when targeting es2015+
						if r.useOuterVariableScopeInParameter(result, location, lastLocation) {
							useResult = false
						} else if result.flags&ast.SymbolFlagsFunctionScopedVariable != 0 {
							// parameters are visible only inside function body, parameter list and return type
							// technically for parameter list case here we might mix parameters and variables declared in function,
							// however it is detected separately when checking initializers of parameters
							// to make sure that they reference no variables declared after them.
							useResult = lastLocation.kind == ast.KindParameter ||
								lastLocation.flags&ast.NodeFlagsSynthesized != 0 ||
								lastLocation == location.ReturnType() && findAncestor(result.valueDeclaration, isParameter) != nil
						}
					}
				} else if location.kind == ast.KindConditionalType {
					// A type parameter declared using 'infer T' in a conditional type is visible only in
					// the true branch of the conditional type.
					useResult = lastLocation == location.AsConditionalTypeNode().trueType
				}
				if useResult {
					break loop
				}
				result = nil
			}
		}
		withinDeferredContext = withinDeferredContext || getIsDeferredContext(location, lastLocation)
		switch location.kind {
		case ast.KindSourceFile:
			if !isExternalOrCommonJsModule(location.AsSourceFile()) {
				break
			}
			fallthrough
		case ast.KindModuleDeclaration:
			moduleExports := r.getSymbolOfDeclaration(location).exports
			if isSourceFile(location) || (isModuleDeclaration(location) && location.flags&ast.NodeFlagsAmbient != 0 && !isGlobalScopeAugmentation(location)) {
				// It's an external module. First see if the module has an export default and if the local
				// name of that export default matches.
				result = moduleExports[InternalSymbolNameDefault]
				if result != nil {
					localSymbol := getLocalSymbolForExportDefault(result)
					if localSymbol != nil && result.flags&meaning != 0 && localSymbol.name == name {
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
				if moduleExport != nil && moduleExport.flags == ast.SymbolFlagsAlias && (getDeclarationOfKind(moduleExport, ast.KindExportSpecifier) != nil || getDeclarationOfKind(moduleExport, ast.KindNamespaceExport) != nil) {
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
			result = r.lookup(r.getSymbolOfDeclaration(location).exports, name, meaning&ast.SymbolFlagsEnumMember)
			if result != nil {
				if nameNotFoundMessage != nil && getIsolatedModules(r.compilerOptions) && location.flags&ast.NodeFlagsAmbient == 0 && getSourceFileOfNode(location) != getSourceFileOfNode(result.valueDeclaration) {
					isolatedModulesLikeFlagName := ifElse(r.compilerOptions.VerbatimModuleSyntax == core.TSTrue, "verbatimModuleSyntax", "isolatedModules")
					r.error(originalLocation, diagnostics.Cannot_access_0_from_another_file_without_qualification_when_1_is_enabled_Use_2_instead,
						name, isolatedModulesLikeFlagName, r.getSymbolOfDeclaration(location).name+"."+name)
				}
				break loop
			}
		case ast.KindPropertyDeclaration:
			if !isStatic(location) {
				ctor := findConstructorDeclaration(location.parent)
				if ctor != nil && ctor.AsConstructorDeclaration().locals != nil {
					if r.lookup(ctor.AsConstructorDeclaration().locals, name, meaning&ast.SymbolFlagsValue) != nil {
						// Remember the property node, it will be used later to report appropriate error
						propertyWithInvalidInitializer = location
					}
				}
			}
		case ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration:
			result = r.lookup(r.getSymbolOfDeclaration(location).members, name, meaning&ast.SymbolFlagsType)
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
			if isClassExpression(location) && meaning&ast.SymbolFlagsClass != 0 {
				className := location.AsClassExpression().name
				if className != nil && name == className.AsIdentifier().text {
					result = location.AsClassExpression().symbol
					break loop
				}
			}
		case ast.KindExpressionWithTypeArguments:
			if lastLocation == location.AsExpressionWithTypeArguments().expression && isHeritageClause(location.parent) && location.parent.AsHeritageClause().token == ast.KindExtendsKeyword {
				container := location.parent.parent
				if isClassLike(container) {
					result = r.lookup(r.getSymbolOfDeclaration(container).members, name, meaning&ast.SymbolFlagsType)
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
			grandparent = location.parent.parent
			if isClassLike(grandparent) || isInterfaceDeclaration(grandparent) {
				// A reference to this grandparent's type parameters would be an error
				result = r.lookup(r.getSymbolOfDeclaration(grandparent).members, name, meaning&ast.SymbolFlagsType)
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
				functionName := location.AsFunctionExpression().name
				if functionName != nil && name == functionName.AsIdentifier().text {
					result = location.AsFunctionExpression().symbol
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
			if location.parent != nil && location.parent.kind == ast.KindParameter {
				location = location.parent
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
			if location.parent != nil && (isClassElement(location.parent) || location.parent.kind == ast.KindClassDeclaration) {
				location = location.parent
			}
		case ast.KindParameter:
			parameterDeclaration := location.AsParameterDeclaration()
			if lastLocation != nil && (lastLocation == parameterDeclaration.initializer ||
				lastLocation == parameterDeclaration.name && isBindingPattern(lastLocation)) {
				if associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case ast.KindBindingElement:
			bindingElement := location.AsBindingElement()
			if lastLocation != nil && (lastLocation == bindingElement.initializer ||
				lastLocation == bindingElement.name && isBindingPattern(lastLocation)) {
				if isPartOfParameterDeclaration(location) && associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case ast.KindInferType:
			if meaning&ast.SymbolFlagsTypeParameter != 0 {
				parameterName := location.AsInferTypeNode().typeParameter.AsTypeParameter().name
				if parameterName != nil && name == parameterName.AsIdentifier().text {
					result = location.AsInferTypeNode().typeParameter.AsTypeParameter().symbol
					break loop
				}
			}
		case ast.KindExportSpecifier:
			exportSpecifier := location.AsExportSpecifier()
			if lastLocation != nil && lastLocation == exportSpecifier.propertyName && location.parent.parent.AsExportDeclaration().moduleSpecifier != nil {
				location = location.parent.parent.parent
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
			location = location.parent
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
	if isParameter(lastLocation) {
		body := getBodyOfNode(location)
		if body != nil && result.valueDeclaration != nil && result.valueDeclaration.Pos() >= body.Pos() && result.valueDeclaration.End() <= body.End() {
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
	return r.requiresScopeChangeWorker(d.name) || d.initializer != nil && r.requiresScopeChangeWorker(d.initializer)
}

func (r *NameResolver) requiresScopeChangeWorker(node *Node) bool {
	switch node.kind {
	case ast.KindArrowFunction, ast.KindFunctionExpression, ast.KindFunctionDeclaration, ast.KindConstructor:
		return false
	case ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindPropertyAssignment:
		return r.requiresScopeChangeWorker(node.Name())
	case ast.KindPropertyDeclaration:
		if hasStaticModifier(node) {
			return !getEmitStandardClassFields(r.compilerOptions)
		}
		return r.requiresScopeChangeWorker(node.AsPropertyDeclaration().name)
	default:
		if isNullishCoalesce(node) || isOptionalChain(node) {
			return getEmitScriptTarget(r.compilerOptions) < core.ScriptTargetES2020
		}
		if isBindingElement(node) && node.AsBindingElement().dotDotDotToken != nil && isObjectBindingPattern(node.parent) {
			return getEmitScriptTarget(r.compilerOptions) < core.ScriptTargetES2017
		}
		if isTypeNode(node) {
			return false
		}
		return node.ForEachChild(r.requiresScopeChangeWorker)
	}
}

func getIsDeferredContext(location *Node, lastLocation *Node) bool {
	if location.kind != ast.KindArrowFunction && location.kind != ast.KindFunctionExpression {
		// initializers in instance property declaration of class like entities are executed in constructor and thus deferred
		// A name is evaluated within the enclosing scope - so it shouldn't count as deferred
		return isTypeQueryNode(location) ||
			(isFunctionLikeDeclaration(location) || location.kind == ast.KindPropertyDeclaration && !isStatic(location)) &&
				(lastLocation == nil || lastLocation != location.Name())
	}
	if lastLocation != nil && lastLocation == location.Name() {
		return false
	}
	// generator functions and async functions are not inlined in control flow when immediately invoked
	if location.BodyData().asteriskToken != nil || hasSyntacticModifier(location, ModifierFlagsAsync) {
		return true
	}
	return getImmediatelyInvokedFunctionExpression(location) == nil
}

func isTypeParameterSymbolDeclaredInContainer(symbol *Symbol, container *Node) bool {
	for _, decl := range symbol.declarations {
		if decl.kind == ast.KindTypeParameter {
			parent := decl.parent.parent
			if parent == container {
				return true
			}
		}
	}
	return false
}

func isSelfReferenceLocation(node *Node, lastLocation *Node) bool {
	switch node.kind {
	case ast.KindParameter:
		return lastLocation != nil && lastLocation == node.AsParameterDeclaration().name
	case ast.KindFunctionDeclaration, ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindEnumDeclaration,
		ast.KindTypeAliasDeclaration, ast.KindModuleDeclaration: // For `namespace N { N; }`
		return true
	}
	return false
}

func isTypeReferenceIdentifier(node *Node) bool {
	for node.parent.kind == ast.KindQualifiedName {
		node = node.parent
	}
	return isTypeReferenceNode(node.parent)
}

func isInTypeQuery(node *Node) bool {
	// TypeScript 1.0 spec (April 2014): 3.6.3
	// A type query consists of the keyword typeof followed by an expression.
	// The expression is restricted to a single identifier or a sequence of identifiers separated by periods
	return findAncestorOrQuit(node, func(n *Node) FindAncestorResult {
		switch n.kind {
		case ast.KindTypeQuery:
			return FindAncestorTrue
		case ast.KindIdentifier, ast.KindQualifiedName:
			return FindAncestorFalse
		}
		return FindAncestorQuit
	}) != nil
}

func nodeKindIs(node *Node, kinds ...ast.Kind) bool {
	return slices.Contains(kinds, node.kind)
}

func isTypeOnlyImportDeclaration(node *Node) bool {
	switch node.kind {
	case ast.KindImportSpecifier:
		return node.AsImportSpecifier().isTypeOnly || node.parent.parent.AsImportClause().isTypeOnly
	case ast.KindNamespaceImport:
		return node.parent.AsImportClause().isTypeOnly
	case ast.KindImportClause:
		return node.AsImportClause().isTypeOnly
	case ast.KindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().isTypeOnly
	}
	return false
}

func isTypeOnlyExportDeclaration(node *Node) bool {
	switch node.kind {
	case ast.KindExportSpecifier:
		return node.AsExportSpecifier().isTypeOnly || node.parent.parent.AsExportDeclaration().isTypeOnly
	case ast.KindExportDeclaration:
		d := node.AsExportDeclaration()
		return d.isTypeOnly && d.moduleSpecifier != nil && d.exportClause == nil
	case ast.KindNamespaceExport:
		return node.parent.AsExportDeclaration().isTypeOnly
	}
	return false
}

func isTypeOnlyImportOrExportDeclaration(node *Node) bool {
	return isTypeOnlyImportDeclaration(node) || isTypeOnlyExportDeclaration(node)
}

func getNameFromImportDeclaration(node *Node) *Node {
	switch node.kind {
	case ast.KindImportSpecifier:
		return node.AsImportSpecifier().name
	case ast.KindNamespaceImport:
		return node.AsNamespaceImport().name
	case ast.KindImportClause:
		return node.AsImportClause().name
	case ast.KindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().name
	}
	return nil
}

func isValidTypeOnlyAliasUseSite(useSite *Node) bool {
	return useSite.flags&ast.NodeFlagsAmbient != 0 ||
		isPartOfTypeQuery(useSite) ||
		isIdentifierInNonEmittingHeritageClause(useSite) ||
		isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(useSite) ||
		!(isExpressionNode(useSite) || isShorthandPropertyNameUseSite(useSite))
}

func isIdentifierInNonEmittingHeritageClause(node *Node) bool {
	if node.kind != ast.KindIdentifier {
		return false
	}
	heritageClause := findAncestorOrQuit(node.parent, func(parent *Node) FindAncestorResult {
		switch parent.kind {
		case ast.KindHeritageClause:
			return FindAncestorTrue
		case ast.KindPropertyAccessExpression, ast.KindExpressionWithTypeArguments:
			return FindAncestorFalse
		default:
			return FindAncestorQuit
		}
	})
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().token == ast.KindImmediateKeyword || heritageClause.parent.kind == ast.KindInterfaceDeclaration
	}
	return false
}

func isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(node *Node) bool {
	for nodeKindIs(node, ast.KindIdentifier, ast.KindPropertyAccessExpression) {
		node = node.parent
	}
	if node.kind != ast.KindComputedPropertyName {
		return false
	}
	if hasSyntacticModifier(node.parent, ModifierFlagsAbstract) {
		return true
	}
	return nodeKindIs(node.parent.parent, ast.KindInterfaceDeclaration, ast.KindTypeLiteral)
}

func isExpressionNode(node *Node) bool {
	switch node.kind {
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
		return !isHeritageClause(node.parent)
	case ast.KindQualifiedName:
		for node.parent.kind == ast.KindQualifiedName {
			node = node.parent
		}
		return isTypeQueryNode(node.parent) || isJSDocLinkLike(node.parent) || isJSXTagName(node)
	case ast.KindJSDocMemberName:
		return isTypeQueryNode(node.parent) || isJSDocLinkLike(node.parent) || isJSXTagName(node)
	case ast.KindPrivateIdentifier:
		return isBinaryExpression(node.parent) && node.parent.AsBinaryExpression().left == node && node.parent.AsBinaryExpression().operatorToken.kind == ast.KindInKeyword
	case ast.KindIdentifier:
		if isTypeQueryNode(node.parent) || isJSDocLinkLike(node.parent) || isJSXTagName(node) {
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
	parent := node.parent
	switch parent.kind {
	case ast.KindVariableDeclaration:
		return parent.AsVariableDeclaration().initializer == node
	case ast.KindParameter:
		return parent.AsParameterDeclaration().initializer == node
	case ast.KindPropertyDeclaration:
		return parent.AsPropertyDeclaration().initializer == node
	case ast.KindPropertySignature:
		return parent.AsPropertySignatureDeclaration().initializer == node
	case ast.KindEnumMember:
		return parent.AsEnumMember().initializer == node
	case ast.KindPropertyAssignment:
		return parent.AsPropertyAssignment().initializer == node
	case ast.KindBindingElement:
		return parent.AsBindingElement().initializer == node
	case ast.KindExpressionStatement:
		return parent.AsExpressionStatement().expression == node
	case ast.KindIfStatement:
		return parent.AsIfStatement().expression == node
	case ast.KindDoStatement:
		return parent.AsDoStatement().expression == node
	case ast.KindWhileStatement:
		return parent.AsWhileStatement().expression == node
	case ast.KindReturnStatement:
		return parent.AsReturnStatement().expression == node
	case ast.KindWithStatement:
		return parent.AsWithStatement().expression == node
	case ast.KindSwitchStatement:
		return parent.AsSwitchStatement().expression == node
	case ast.KindCaseClause, ast.KindDefaultClause:
		return parent.AsCaseOrDefaultClause().expression == node
	case ast.KindThrowStatement:
		return parent.AsThrowStatement().expression == node
	case ast.KindForStatement:
		s := parent.AsForStatement()
		return s.initializer == node && s.initializer.kind != ast.KindVariableDeclarationList || s.condition == node || s.incrementor == node
	case ast.KindForInStatement, ast.KindForOfStatement:
		s := parent.AsForInOrOfStatement()
		return s.initializer == node && s.initializer.kind != ast.KindVariableDeclarationList || s.expression == node
	case ast.KindTypeAssertionExpression:
		return parent.AsTypeAssertion().expression == node
	case ast.KindAsExpression:
		return parent.AsAsExpression().expression == node
	case ast.KindTemplateSpan:
		return parent.AsTemplateSpan().expression == node
	case ast.KindComputedPropertyName:
		return parent.AsComputedPropertyName().expression == node
	case ast.KindDecorator, ast.KindJsxExpression, ast.KindJsxSpreadAttribute, ast.KindSpreadAssignment:
		return true
	case ast.KindExpressionWithTypeArguments:
		return parent.AsExpressionWithTypeArguments().expression == node && !isPartOfTypeNode(parent)
	case ast.KindShorthandPropertyAssignment:
		return parent.AsShorthandPropertyAssignment().objectAssignmentInitializer == node
	case ast.KindSatisfiesExpression:
		return parent.AsSatisfiesExpression().expression == node
	default:
		return isExpressionNode(parent)
	}
}

func isPartOfTypeNode(node *Node) bool {
	kind := node.kind
	if kind >= ast.KindFirstTypeNode && kind <= ast.KindLastTypeNode {
		return true
	}
	switch node.kind {
	case ast.KindAnyKeyword, ast.KindUnknownKeyword, ast.KindNumberKeyword, ast.KindBigIntKeyword, ast.KindStringKeyword,
		ast.KindBooleanKeyword, ast.KindSymbolKeyword, ast.KindObjectKeyword, ast.KindUndefinedKeyword, ast.KindNullKeyword,
		ast.KindNeverKeyword:
		return true
	case ast.KindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(node)
	case ast.KindTypeParameter:
		return node.parent.kind == ast.KindMappedType || node.parent.kind == ast.KindInferType
	case ast.KindIdentifier:
		parent := node.parent
		if isQualifiedName(parent) && parent.AsQualifiedName().right == node {
			return isPartOfTypeNodeInParent(parent)
		}
		if isPropertyAccessExpression(parent) && parent.AsPropertyAccessExpression().name == node {
			return isPartOfTypeNodeInParent(parent)
		}
		return isPartOfTypeNodeInParent(node)
	case ast.KindQualifiedName, ast.KindPropertyAccessExpression, ast.KindThisKeyword:
		return isPartOfTypeNodeInParent(node)
	}
	return false
}

func isPartOfTypeNodeInParent(node *Node) bool {
	parent := node.parent
	// Do not recursively call isPartOfTypeNode on the parent. In the example:
	//
	//     let a: A.B.C;
	//
	// Calling isPartOfTypeNode would consider the qualified name A.B a type node.
	// Only C and A.B.C are type nodes.
	if parent.kind >= ast.KindFirstTypeNode && parent.kind <= ast.KindLastTypeNode {
		return true
	}
	switch parent.kind {
	case ast.KindTypeQuery:
		return false
	case ast.KindImportType:
		return !parent.AsImportTypeNode().isTypeOf
	case ast.KindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(parent)
	case ast.KindTypeParameter:
		return node == parent.AsTypeParameter().constraint
	case ast.KindPropertyDeclaration:
		return node == parent.AsPropertyDeclaration().typeNode
	case ast.KindPropertySignature:
		return node == parent.AsPropertySignatureDeclaration().typeNode
	case ast.KindParameter:
		return node == parent.AsParameterDeclaration().typeNode
	case ast.KindVariableDeclaration:
		return node == parent.AsVariableDeclaration().typeNode
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindConstructor, ast.KindMethodDeclaration,
		ast.KindMethodSignature, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindCallSignature, ast.KindConstructSignature,
		ast.KindIndexSignature:
		return node == parent.ReturnType()
	case ast.KindTypeAssertionExpression:
		return node == parent.AsTypeAssertion().typeNode
	case ast.KindCallExpression:
		return typeArgumentListContains(parent.AsCallExpression().typeArguments, node)
	case ast.KindNewExpression:
		return typeArgumentListContains(parent.AsNewExpression().typeArguments, node)
	case ast.KindTaggedTemplateExpression:
		return typeArgumentListContains(parent.AsTaggedTemplateExpression().typeArguments, node)
	}
	return false
}

func isPartOfTypeExpressionWithTypeArguments(node *Node) bool {
	parent := node.parent
	return isHeritageClause(parent) && (!isClassLike(parent.parent) || parent.AsHeritageClause().token == ast.KindImplementsKeyword)
}

func typeArgumentListContains(list *Node, node *Node) bool {
	if list != nil {
		return slices.Contains(list.AsTypeArgumentList().arguments, node)
	}
	return false
}

func isJSDocLinkLike(node *Node) bool {
	return nodeKindIs(node, ast.KindJSDocLink, ast.KindJSDocLinkCode, ast.KindJSDocLinkPlain)
}

func isJSXTagName(node *Node) bool {
	parent := node.parent
	switch parent.kind {
	case ast.KindJsxOpeningElement:
		return parent.AsJsxOpeningElement().tagName == node
	case ast.KindJsxSelfClosingElement:
		return parent.AsJsxSelfClosingElement().tagName == node
	case ast.KindJsxClosingElement:
		return parent.AsJsxClosingElement().tagName == node
	}
	return false
}

func isShorthandPropertyNameUseSite(useSite *Node) bool {
	return isIdentifier(useSite) && isShorthandPropertyAssignment(useSite.parent) && useSite.parent.AsShorthandPropertyAssignment().name == useSite
}

func isTypeDeclaration(node *Node) bool {
	switch node.kind {
	case ast.KindTypeParameter, ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindTypeAliasDeclaration, ast.KindEnumDeclaration:
		return true
	case ast.KindImportClause:
		return node.AsImportClause().isTypeOnly
	case ast.KindImportSpecifier:
		return node.parent.parent.AsImportClause().isTypeOnly
	case ast.KindExportSpecifier:
		return node.parent.parent.AsExportDeclaration().isTypeOnly
	default:
		return false
	}
}

func canHaveSymbol(node *Node) bool {
	switch node.kind {
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
	switch node.kind {
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
	return isAnyImportSyntax(node) || isExportDeclaration(node)
}

func isAnyImportSyntax(node *Node) bool {
	return nodeKindIs(node, ast.KindImportDeclaration, ast.KindImportEqualsDeclaration)
}

func getExternalModuleName(node *Node) *Node {
	switch node.kind {
	case ast.KindImportDeclaration:
		return node.AsImportDeclaration().moduleSpecifier
	case ast.KindExportDeclaration:
		return node.AsExportDeclaration().moduleSpecifier
	case ast.KindImportEqualsDeclaration:
		if node.AsImportEqualsDeclaration().moduleReference.kind == ast.KindExternalModuleReference {
			return node.AsImportEqualsDeclaration().moduleReference.AsExternalModuleReference().expression
		}
		return nil
	case ast.KindImportType:
		return getImportTypeNodeLiteral(node)
	case ast.KindCallExpression:
		return node.AsCallExpression().arguments[0]
	case ast.KindModuleDeclaration:
		if isStringLiteral(node.AsModuleDeclaration().name) {
			return node.AsModuleDeclaration().name
		}
		return nil
	}
	panic("Unhandled case in getExternalModuleName")
}

func getImportTypeNodeLiteral(node *Node) *Node {
	if isImportTypeNode(node) {
		importTypeNode := node.AsImportTypeNode()
		if isLiteralTypeNode(importTypeNode.argument) {
			literalTypeNode := importTypeNode.argument.AsLiteralTypeNode()
			if isStringLiteral(literalTypeNode.literal) {
				return literalTypeNode.literal
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
	return isShorthandAmbientModule(moduleSymbol.valueDeclaration)
}

func isShorthandAmbientModule(node *Node) bool {
	// The only kind of module that can be missing a body is a shorthand ambient module.
	return node != nil && node.kind == ast.KindModuleDeclaration && node.AsModuleDeclaration().body == nil
}

func isEntityName(node *Node) bool {
	return node.kind == ast.KindIdentifier || node.kind == ast.KindQualifiedName
}

func nodeIsSynthesized(node *Node) bool {
	return node.loc.pos < 0 || node.loc.end < 0
}

func getFirstIdentifier(node *Node) *Node {
	switch node.kind {
	case ast.KindIdentifier:
		return node
	case ast.KindQualifiedName:
		return getFirstIdentifier(node.AsQualifiedName().left)
	case ast.KindPropertyAccessExpression:
		return getFirstIdentifier(node.AsPropertyAccessExpression().expression)
	}
	panic("Unhandled case in getFirstIdentifier")
}

func getAliasDeclarationFromName(node *Node) *Node {
	switch node.kind {
	case ast.KindImportClause, ast.KindImportSpecifier, ast.KindNamespaceImport, ast.KindExportSpecifier, ast.KindExportAssignment,
		ast.KindImportEqualsDeclaration, ast.KindNamespaceExport:
		return node.parent
	case ast.KindQualifiedName:
		return getAliasDeclarationFromName(node.parent)
	}
	return nil
}

func entityNameToString(name *Node) string {
	switch name.kind {
	case ast.KindThisKeyword:
		return "this"
	case ast.KindIdentifier, ast.KindPrivateIdentifier:
		return getTextOfNode(name)
	case ast.KindQualifiedName:
		return entityNameToString(name.AsQualifiedName().left) + "." + entityNameToString(name.AsQualifiedName().right)
	case ast.KindPropertyAccessExpression:
		return entityNameToString(name.AsPropertyAccessExpression().expression) + "." + entityNameToString(name.AsPropertyAccessExpression().name)
	case ast.KindJsxNamespacedName:
		return entityNameToString(name.AsJsxNamespacedName().namespace) + ":" + entityNameToString(name.AsJsxNamespacedName().name)
	}
	panic("Unhandled case in entityNameToString")
}

func getContainingQualifiedNameNode(node *Node) *Node {
	for isQualifiedName(node.parent) {
		node = node.parent
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
	ancestor := findAncestor(node, isImportDeclaration)
	return ancestor != nil && ancestor.AsImportDeclaration().importClause == nil
}

func getExternalModuleRequireArgument(node *Node) *Node {
	if isVariableDeclarationInitializedToBareOrAccessedRequire(node) {
		return getLeftmostAccessExpression(node.AsVariableDeclaration().initializer).AsCallExpression().arguments[0]
	}
	return nil
}

func getExternalModuleImportEqualsDeclarationExpression(node *Node) *Node {
	//Debug.assert(isExternalModuleImportEqualsDeclaration(node))
	return node.AsImportEqualsDeclaration().moduleReference.AsExternalModuleReference().expression
}

func isRightSideOfQualifiedNameOrPropertyAccess(node *Node) bool {
	parent := node.parent
	switch parent.kind {
	case ast.KindQualifiedName:
		return parent.AsQualifiedName().right == node
	case ast.KindPropertyAccessExpression:
		return parent.AsPropertyAccessExpression().name == node
	case ast.KindMetaProperty:
		return parent.AsMetaProperty().name == node
	}
	return false
}

func getNamespaceDeclarationNode(node *Node) *Node {
	switch node.kind {
	case ast.KindImportDeclaration:
		importClause := node.AsImportDeclaration().importClause
		if importClause != nil && isNamespaceImport(importClause.AsImportClause().namedBindings) {
			return importClause.AsImportClause().namedBindings
		}
	case ast.KindImportEqualsDeclaration:
		return node
	case ast.KindExportDeclaration:
		exportClause := node.AsExportDeclaration().exportClause
		if exportClause != nil && isNamespaceExport(exportClause) {
			return exportClause
		}
	default:
		panic("Unhandled case in getNamespaceDeclarationNode")
	}
	return nil
}

func isImportCall(node *Node) bool {
	return isCallExpression(node) && node.AsCallExpression().expression.kind == ast.KindImportKeyword
}

func getSourceFileOfModule(module *Symbol) *SourceFile {
	declaration := module.valueDeclaration
	if declaration == nil {
		declaration = getNonAugmentationDeclaration(module)
	}
	return getSourceFileOfNode(declaration)
}

func getNonAugmentationDeclaration(symbol *Symbol) *Node {
	return core.Find(symbol.declarations, func(d *Node) bool {
		return !isExternalModuleAugmentation(d) && !(isModuleDeclaration(d) && isGlobalScopeAugmentation(d))
	})
}

func isExternalModuleAugmentation(node *Node) bool {
	return isAmbientModule(node) && isModuleAugmentationExternal(node)
}

func isJsonSourceFile(file *SourceFile) bool {
	return file.scriptKind == core.ScriptKindJSON
}

func isSyntacticDefault(node *Node) bool {
	return (isExportAssignment(node) && !node.AsExportAssignment().isExportEquals) ||
		hasSyntacticModifier(node, ModifierFlagsDefault) ||
		isExportSpecifier(node) ||
		isNamespaceExport(node)
}

func hasExportAssignmentSymbol(moduleSymbol *Symbol) bool {
	return moduleSymbol.exports[InternalSymbolNameExportEquals] != nil
}

func isImportOrExportSpecifier(node *Node) bool {
	return isImportSpecifier(node) || isExportSpecifier(node)
}

func parsePseudoBigInt(stringValue string) string {
	return stringValue // !!!
}

func isTypeAlias(node *Node) bool {
	return isTypeAliasDeclaration(node)
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
		return typeParameters.AsTypeParameterList().parameters
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
		return typeParameterList.AsTypeParameterList().parameters
	}
	return nil
}

func getTypeArgumentNodesFromNode(node *Node) []*Node {
	typeArgumentList := getTypeArgumentListFromNode(node)
	if typeArgumentList != nil {
		return typeArgumentList.AsTypeArgumentList().arguments
	}
	return nil
}

func getTypeArgumentListFromNode(node *Node) *Node {
	switch node.kind {
	case ast.KindCallExpression:
		return node.AsCallExpression().typeArguments
	case ast.KindNewExpression:
		return node.AsNewExpression().typeArguments
	case ast.KindTaggedTemplateExpression:
		return node.AsTaggedTemplateExpression().typeArguments
	case ast.KindTypeReference:
		return node.AsTypeReference().typeArguments
	case ast.KindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().typeArguments
	case ast.KindImportType:
		return node.AsImportTypeNode().typeArguments
	case ast.KindTypeQuery:
		return node.AsTypeQueryNode().typeArguments
	}
	panic("Unhandled case in getTypeArgumentListFromNode")
}

func getInitializerFromNode(node *Node) *Node {
	switch node.kind {
	case ast.KindVariableDeclaration:
		return node.AsVariableDeclaration().initializer
	case ast.KindParameter:
		return node.AsParameterDeclaration().initializer
	case ast.KindBindingElement:
		return node.AsBindingElement().initializer
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().initializer
	case ast.KindPropertyAssignment:
		return node.AsPropertyAssignment().initializer
	case ast.KindEnumMember:
		return node.AsEnumMember().initializer
	case ast.KindForStatement:
		return node.AsForStatement().initializer
	case ast.KindForInStatement, ast.KindForOfStatement:
		return node.AsForInOrOfStatement().initializer
	case ast.KindJsxAttribute:
		return node.AsJsxAttribute().initializer
	}
	return nil
}

/**
 * Gets the effective type annotation of a variable, parameter, or property. If the node was
 * parsed in a JavaScript file, gets the type annotation from JSDoc.  Also gets the type of
 * functions only the JSDoc case.
 */
func getEffectiveTypeAnnotationNode(node *Node) *Node {
	switch node.kind {
	case ast.KindVariableDeclaration:
		return node.AsVariableDeclaration().typeNode
	case ast.KindParameter:
		return node.AsParameterDeclaration().typeNode
	case ast.KindPropertySignature:
		return node.AsPropertySignatureDeclaration().typeNode
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().typeNode
	case ast.KindTypePredicate:
		return node.AsTypePredicateNode().typeNode
	case ast.KindParenthesizedType:
		return node.AsParenthesizedTypeNode().typeNode
	case ast.KindTypeOperator:
		return node.AsTypeOperatorNode().typeNode
	case ast.KindMappedType:
		return node.AsMappedTypeNode().typeNode
	case ast.KindTypeAssertionExpression:
		return node.AsTypeAssertion().typeNode
	case ast.KindAsExpression:
		return node.AsAsExpression().typeNode
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
	return node != nil && node.kind == ast.KindQuestionToken
}

func isOptionalDeclaration(declaration *Node) bool {
	switch declaration.kind {
	case ast.KindParameter:
		return declaration.AsParameterDeclaration().questionToken != nil
	case ast.KindPropertyDeclaration:
		return isQuestionToken(declaration.AsPropertyDeclaration().postfixToken)
	case ast.KindPropertySignature:
		return isQuestionToken(declaration.AsPropertySignatureDeclaration().postfixToken)
	case ast.KindMethodDeclaration:
		return isQuestionToken(declaration.AsMethodDeclaration().postfixToken)
	case ast.KindMethodSignature:
		return isQuestionToken(declaration.AsMethodSignatureDeclaration().postfixToken)
	case ast.KindPropertyAssignment:
		return isQuestionToken(declaration.AsPropertyAssignment().postfixToken)
	case ast.KindShorthandPropertyAssignment:
		return isQuestionToken(declaration.AsShorthandPropertyAssignment().postfixToken)
	}
	return false
}

func isEmptyArrayLiteral(expression *Node) bool {
	return expression.kind == ast.KindArrayLiteralExpression && len(expression.AsArrayLiteralExpression().elements) == 0
}

func declarationBelongsToPrivateAmbientMember(declaration *Node) bool {
	root := getRootDeclaration(declaration)
	memberDeclaration := root
	if root.kind == ast.KindParameter {
		memberDeclaration = root.parent
	}
	return isPrivateWithinAmbient(memberDeclaration)
}

func isPrivateWithinAmbient(node *Node) bool {
	return (hasEffectiveModifier(node, ModifierFlagsPrivate) || isPrivateIdentifierClassElementDeclaration(node)) && node.flags&ast.NodeFlagsAmbient != 0
}

func identifierToKeywordKind(node *Identifier) ast.Kind {
	return textToKeyword[node.text]
}

func isAssertionExpression(node *Node) bool {
	kind := node.kind
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
		result[symbol.name] = symbol
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
	if s1.valueDeclaration != nil && s2.valueDeclaration != nil {
		if s1.parent != nil && s2.parent != nil {
			// Symbols with the same unmerged parent are always in the same file
			if s1.parent != s2.parent {
				f1 := getSourceFileOfNode(s1.valueDeclaration)
				f2 := getSourceFileOfNode(s2.valueDeclaration)
				if f1 != f2 {
					// In different files, first compare base filename
					r := strings.Compare(filepath.Base(f1.path), filepath.Base(f2.path))
					if r == 0 {
						// Same base filename, compare the full paths (no two files should have the same full path)
						r = strings.Compare(f1.path, f2.path)
					}
					return r
				}
			}
			// In the same file, compare source positions
			return s1.valueDeclaration.Pos() - s2.valueDeclaration.Pos()
		}
	}
	// Sort by name
	r := strings.Compare(s1.name, s2.name)
	if r == 0 {
		// Same name, sort by symbol id
		r = int(getSymbolId(s1)) - int(getSymbolId(s2))
	}
	return r
}

func getClassLikeDeclarationOfSymbol(symbol *Symbol) *Node {
	return core.Find(symbol.declarations, isClassLike)
}

func isThisInTypeQuery(node *Node) bool {
	if !isThisIdentifier(node) {
		return false
	}
	for isQualifiedName(node.parent) && node.parent.AsQualifiedName().left == node {
		node = node.parent
	}
	return node.parent.kind == ast.KindTypeQuery
}

func isThisIdentifier(node *Node) bool {
	return node != nil && node.kind == ast.KindIdentifier && identifierIsThisKeyword(node)
}

func identifierIsThisKeyword(id *Node) bool {
	return id.AsIdentifier().text == "this"
}

func getDeclarationModifierFlagsFromSymbol(s *Symbol) ModifierFlags {
	return getDeclarationModifierFlagsFromSymbolEx(s, false /*isWrite*/)
}

func getDeclarationModifierFlagsFromSymbolEx(s *Symbol, isWrite bool) ModifierFlags {
	if s.valueDeclaration != nil {
		var declaration *Node
		if isWrite {
			declaration = core.Find(s.declarations, isSetAccessorDeclaration)
		}
		if declaration == nil && s.flags&ast.SymbolFlagsGetAccessor != 0 {
			declaration = core.Find(s.declarations, isGetAccessorDeclaration)
		}
		if declaration == nil {
			declaration = s.valueDeclaration
		}
		flags := getCombinedModifierFlags(declaration)
		if s.parent != nil && s.parent.flags&ast.SymbolFlagsClass != 0 {
			return flags
		}
		return flags & ^ModifierFlagsAccessibilityModifier
	}
	if s.checkFlags&ast.CheckFlagsSynthetic != 0 {
		var accessModifier ModifierFlags
		switch {
		case s.checkFlags&ast.CheckFlagsContainsPrivate != 0:
			accessModifier = ModifierFlagsPrivate
		case s.checkFlags&ast.CheckFlagsContainsPublic != 0:
			accessModifier = ModifierFlagsPublic
		default:
			accessModifier = ModifierFlagsProtected
		}
		var staticModifier ModifierFlags
		if s.checkFlags&ast.CheckFlagsContainsStatic != 0 {
			staticModifier = ModifierFlagsStatic
		}
		return accessModifier | staticModifier
	}
	if s.flags&ast.SymbolFlagsPrototype != 0 {
		return ModifierFlagsPublic | ModifierFlagsStatic
	}
	return ModifierFlagsNone
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
	return getCombinedModifierFlags(declaration)&ModifierFlagsReadonly != 0 && !isParameterPropertyDeclaration(declaration, declaration.parent)
}

func getPostfixTokenFromNode(node *Node) *Node {
	switch node.kind {
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().postfixToken
	case ast.KindPropertySignature:
		return node.AsPropertySignatureDeclaration().postfixToken
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration().postfixToken
	case ast.KindMethodSignature:
		return node.AsMethodSignatureDeclaration().postfixToken
	}
	panic("Unhandled case in getPostfixTokenFromNode")
}

func isStatic(node *Node) bool {
	// https://tc39.es/ecma262/#sec-static-semantics-isstatic
	return isClassElement(node) && hasStaticModifier(node) || isClassStaticBlockDeclaration(node)
}

func isLogicalExpression(node *Node) bool {
	for {
		if node.kind == ast.KindParenthesizedExpression {
			node = node.AsParenthesizedExpression().expression
		} else if node.kind == ast.KindPrefixUnaryExpression && node.AsPrefixUnaryExpression().operator == ast.KindExclamationToken {
			node = node.AsPrefixUnaryExpression().operand
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

type set[T comparable] struct {
	m map[T]struct{}
}

func (s *set[T]) has(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *set[T]) add(key T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[key] = struct{}{}
}

func (s *set[T]) delete(key T) {
	delete(s.m, key)
}

func (s *set[T]) len() int {
	return len(s.m)
}

func (s *set[T]) keys() map[T]struct{} {
	return s.m
}

func getContainingFunction(node *Node) *Node {
	return findAncestor(node.parent, isFunctionLike)
}

func isTypeReferenceType(node *Node) bool {
	return node.kind == ast.KindTypeReference || node.kind == ast.KindExpressionWithTypeArguments
}

func isNodeDescendantOf(node *Node, ancestor *Node) bool {
	for node != nil {
		if node == ancestor {
			return true
		}
		node = node.parent
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
	switch node.kind {
	case ast.KindIdentifier, ast.KindPrivateIdentifier, ast.KindStringLiteral, ast.KindNumericLiteral, ast.KindComputedPropertyName:
		return true
	}
	return false
}

func getPropertyNameForPropertyNameNode(name *Node) string {
	switch name.kind {
	case ast.KindIdentifier, ast.KindPrivateIdentifier, ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral,
		ast.KindNumericLiteral, ast.KindBigIntLiteral, ast.KindJsxNamespacedName:
		return name.Text()
	case ast.KindComputedPropertyName:
		nameExpression := name.AsComputedPropertyName().expression
		if isStringOrNumericLiteralLike(nameExpression) {
			return nameExpression.Text()
		}
		if isSignedNumericLiteral(nameExpression) {
			text := nameExpression.AsPrefixUnaryExpression().operand.Text()
			if nameExpression.AsPrefixUnaryExpression().operator == ast.KindMinusToken {
				text = "-" + text
			}
			return text
		}
		return InternalSymbolNameMissing
	}
	panic("Unhandled case in getPropertyNameForPropertyNameNode")
}

func isThisProperty(node *Node) bool {
	return (isPropertyAccessExpression(node) || isElementAccessExpression(node)) && node.Expression().kind == ast.KindThisKeyword
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
	if isVariableDeclaration(node) {
		return isVarConst(node) && isIdentifier(node.AsVariableDeclaration().name) && isVariableDeclarationInVariableStatement(node)
	}
	if isPropertyDeclaration(node) {
		return hasEffectiveReadonlyModifier(node) && hasStaticModifier(node)
	}
	return isPropertySignatureDeclaration(node) && hasEffectiveReadonlyModifier(node)
}

func isVarConst(node *Node) bool {
	return getCombinedNodeFlags(node)&ast.NodeFlagsBlockScoped == ast.NodeFlagsConst
}

func isVariableDeclarationInVariableStatement(node *Node) bool {
	return isVariableDeclarationList(node.parent) && isVariableStatement(node.parent.parent)
}

func isKnownSymbol(symbol *Symbol) bool {
	return isLateBoundName(symbol.name)
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
	return getSymbolTable(&symbol.members)
}

func getExports(symbol *Symbol) SymbolTable {
	return getSymbolTable(&symbol.exports)
}

func getLocals(container *Node) SymbolTable {
	return getSymbolTable(&container.LocalsContainerData().locals)
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
	heritageClause := getHeritageClause(node.AsInterfaceDeclaration().heritageClauses, ast.KindExtendsKeyword)
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().types
	}
	return nil
}

func getHeritageClause(clauses []*Node, kind ast.Kind) *Node {
	for _, clause := range clauses {
		if clause.AsHeritageClause().token == kind {
			return clause
		}
	}
	return nil
}

func getClassExtendsHeritageElement(node *Node) *Node {
	heritageClause := getHeritageClause(node.ClassLikeData().heritageClauses, ast.KindExtendsKeyword)
	if heritageClause != nil && len(heritageClause.AsHeritageClause().types) > 0 {
		return heritageClause.AsHeritageClause().types[0]
	}
	return nil
}

func concatenateDiagnosticMessageChains(headChain *MessageChain, tailChain *MessageChain) {
	lastChain := headChain
	for len(lastChain.messageChain) != 0 {
		lastChain = lastChain.messageChain[0]
	}
	lastChain.messageChain = []*MessageChain{tailChain}
}

func isObjectOrArrayLiteralType(t *Type) bool {
	return t.objectFlags&(ObjectFlagsObjectLiteral|ObjectFlagsArrayLiteral) != 0
}

func getContainingClassExcludingClassDecorators(node *Node) *ClassLikeDeclaration {
	decorator := findAncestorOrQuit(node.parent, func(n *Node) FindAncestorResult {
		if isClassLike(n) {
			return FindAncestorQuit
		}
		if isDecorator(n) {
			return FindAncestorTrue
		}
		return FindAncestorFalse
	})
	if decorator != nil && isClassLike(decorator.parent) {
		return getContainingClass(decorator.parent)
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
	switch node.kind {
	case ast.KindJsxOpeningElement, ast.KindJsxSelfClosingElement, ast.KindCallExpression, ast.KindNewExpression,
		ast.KindTaggedTemplateExpression, ast.KindDecorator:
		return true
	}
	return false
}

func isCallOrNewExpression(node *Node) bool {
	return isCallExpression(node) || isNewExpression(node)
}

func isClassInstanceProperty(node *Node) bool {
	return node.parent != nil && isClassLike(node.parent) && isPropertyDeclaration(node) && !hasAccessorModifier(node)
}

func isThisInitializedObjectBindingExpression(node *Node) bool {
	return node != nil && (isShorthandPropertyAssignment(node) || isPropertyAssignment(node)) && isBinaryExpression(node.parent.parent) &&
		node.parent.parent.AsBinaryExpression().operatorToken.kind == ast.KindEqualsToken &&
		node.parent.parent.AsBinaryExpression().right.kind == ast.KindThisKeyword
}

func isThisInitializedDeclaration(node *Node) bool {
	return node != nil && isVariableDeclaration(node) && node.AsVariableDeclaration().initializer != nil && node.AsVariableDeclaration().initializer.kind == ast.KindThisKeyword
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
	parent := node.parent
	switch parent.kind {
	case ast.KindParenthesizedExpression:
		return accessKind(parent)
	case ast.KindPrefixUnaryExpression:
		operator := parent.AsPrefixUnaryExpression().operator
		if operator == ast.KindPlusPlusToken || operator == ast.KindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case ast.KindPostfixUnaryExpression:
		operator := parent.AsPostfixUnaryExpression().operator
		if operator == ast.KindPlusPlusToken || operator == ast.KindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case ast.KindBinaryExpression:
		if parent.AsBinaryExpression().left == node {
			operator := parent.AsBinaryExpression().operatorToken
			if isAssignmentOperator(operator.kind) {
				if operator.kind == ast.KindEqualsToken {
					return AccessKindWrite
				}
				return AccessKindReadWrite
			}
		}
		return AccessKindRead
	case ast.KindPropertyAccessExpression:
		if parent.AsPropertyAccessExpression().name != node {
			return AccessKindRead
		}
		return accessKind(parent)
	case ast.KindPropertyAssignment:
		parentAccess := accessKind(parent.parent)
		// In `({ x: varname }) = { x: 1 }`, the left `x` is a read, the right `x` is a write.
		if node == parent.AsPropertyAssignment().name {
			return reverseAccessKind(parentAccess)
		}
		return parentAccess
	case ast.KindShorthandPropertyAssignment:
		// Assume it's the local variable being accessed, since we don't check public properties for --noUnusedLocals.
		if node == parent.AsShorthandPropertyAssignment().objectAssignmentInitializer {
			return AccessKindRead
		}
		return accessKind(parent.parent)
	case ast.KindArrayLiteralExpression:
		return accessKind(parent)
	case ast.KindForInStatement, ast.KindForOfStatement:
		if node == parent.AsForInOrOfStatement().initializer {
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
	return isJsxOpeningElement(node) || isJsxSelfClosingElement(node)
}

func isObjectLiteralElementLike(node *Node) bool {
	switch node.kind {
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment, ast.KindSpreadAssignment,
		ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor:
		return true
	}
	return false
}
