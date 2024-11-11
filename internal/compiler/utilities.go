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

	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/compiler/textpos"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/utils"
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
	if node.Id == 0 {
		node.Id = NodeId(nextNodeId.Add(1))
	}
	return node.Id
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

func getOperatorPrecedence(nodeKind SyntaxKind, operatorKind SyntaxKind, hasArguments bool) OperatorPrecedence {
	switch nodeKind {
	case SyntaxKindCommaListExpression:
		return OperatorPrecedenceComma
	case SyntaxKindSpreadElement:
		return OperatorPrecedenceSpread
	case SyntaxKindYieldExpression:
		return OperatorPrecedenceYield
	case SyntaxKindConditionalExpression:
		return OperatorPrecedenceConditional
	case SyntaxKindBinaryExpression:
		switch operatorKind {
		case SyntaxKindCommaToken:
			return OperatorPrecedenceComma
		case SyntaxKindEqualsToken, SyntaxKindPlusEqualsToken, SyntaxKindMinusEqualsToken, SyntaxKindAsteriskAsteriskEqualsToken,
			SyntaxKindAsteriskEqualsToken, SyntaxKindSlashEqualsToken, SyntaxKindPercentEqualsToken, SyntaxKindLessThanLessThanEqualsToken,
			SyntaxKindGreaterThanGreaterThanEqualsToken, SyntaxKindGreaterThanGreaterThanGreaterThanEqualsToken, SyntaxKindAmpersandEqualsToken,
			SyntaxKindCaretEqualsToken, SyntaxKindBarEqualsToken, SyntaxKindBarBarEqualsToken, SyntaxKindAmpersandAmpersandEqualsToken,
			SyntaxKindQuestionQuestionEqualsToken:
			return OperatorPrecedenceAssignment
		}
		return getBinaryOperatorPrecedence(operatorKind)
	// TODO: Should prefix `++` and `--` be moved to the `Update` precedence?
	case SyntaxKindTypeAssertionExpression, SyntaxKindNonNullExpression, SyntaxKindPrefixUnaryExpression, SyntaxKindTypeOfExpression,
		SyntaxKindVoidExpression, SyntaxKindDeleteExpression, SyntaxKindAwaitExpression:
		return OperatorPrecedenceUnary
	case SyntaxKindPostfixUnaryExpression:
		return OperatorPrecedenceUpdate
	case SyntaxKindCallExpression:
		return OperatorPrecedenceLeftHandSide
	case SyntaxKindNewExpression:
		if hasArguments {
			return OperatorPrecedenceMember
		}
		return OperatorPrecedenceLeftHandSide
	case SyntaxKindTaggedTemplateExpression, SyntaxKindPropertyAccessExpression, SyntaxKindElementAccessExpression, SyntaxKindMetaProperty:
		return OperatorPrecedenceMember
	case SyntaxKindAsExpression, SyntaxKindSatisfiesExpression:
		return OperatorPrecedenceRelational
	case SyntaxKindThisKeyword, SyntaxKindSuperKeyword, SyntaxKindIdentifier, SyntaxKindPrivateIdentifier, SyntaxKindNullKeyword,
		SyntaxKindTrueKeyword, SyntaxKindFalseKeyword, SyntaxKindNumericLiteral, SyntaxKindBigIntLiteral, SyntaxKindStringLiteral,
		SyntaxKindArrayLiteralExpression, SyntaxKindObjectLiteralExpression, SyntaxKindFunctionExpression, SyntaxKindArrowFunction,
		SyntaxKindClassExpression, SyntaxKindRegularExpressionLiteral, SyntaxKindNoSubstitutionTemplateLiteral, SyntaxKindTemplateExpression,
		SyntaxKindParenthesizedExpression, SyntaxKindOmittedExpression, SyntaxKindJsxElement, SyntaxKindJsxSelfClosingElement, SyntaxKindJsxFragment:
		return OperatorPrecedencePrimary
	}
	return OperatorPrecedenceInvalid
}

func getBinaryOperatorPrecedence(kind SyntaxKind) OperatorPrecedence {
	switch kind {
	case SyntaxKindQuestionQuestionToken:
		return OperatorPrecedenceCoalesce
	case SyntaxKindBarBarToken:
		return OperatorPrecedenceLogicalOR
	case SyntaxKindAmpersandAmpersandToken:
		return OperatorPrecedenceLogicalAND
	case SyntaxKindBarToken:
		return OperatorPrecedenceBitwiseOR
	case SyntaxKindCaretToken:
		return OperatorPrecedenceBitwiseXOR
	case SyntaxKindAmpersandToken:
		return OperatorPrecedenceBitwiseAND
	case SyntaxKindEqualsEqualsToken, SyntaxKindExclamationEqualsToken, SyntaxKindEqualsEqualsEqualsToken, SyntaxKindExclamationEqualsEqualsToken:
		return OperatorPrecedenceEquality
	case SyntaxKindLessThanToken, SyntaxKindGreaterThanToken, SyntaxKindLessThanEqualsToken, SyntaxKindGreaterThanEqualsToken,
		SyntaxKindInstanceOfKeyword, SyntaxKindInKeyword, SyntaxKindAsKeyword, SyntaxKindSatisfiesKeyword:
		return OperatorPrecedenceRelational
	case SyntaxKindLessThanLessThanToken, SyntaxKindGreaterThanGreaterThanToken, SyntaxKindGreaterThanGreaterThanGreaterThanToken:
		return OperatorPrecedenceShift
	case SyntaxKindPlusToken, SyntaxKindMinusToken:
		return OperatorPrecedenceAdditive
	case SyntaxKindAsteriskToken, SyntaxKindSlashToken, SyntaxKindPercentToken:
		return OperatorPrecedenceMultiplicative
	case SyntaxKindAsteriskAsteriskToken:
		return OperatorPrecedenceExponentiation
	}
	// -1 is lower than all other precedences.  Returning it will cause binary expression
	// parsing to stop.
	return OperatorPrecedenceInvalid
}

func formatStringFromArgs(text string, args []any) string {
	return utils.MakeRegexp(`{(\d+)}`).ReplaceAllStringFunc(text, func(match string) string {
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

func boolToTristate(b bool) Tristate {
	if b {
		return TSTrue
	}
	return TSFalse
}

func modifierToFlag(token SyntaxKind) ModifierFlags {
	switch token {
	case SyntaxKindStaticKeyword:
		return ModifierFlagsStatic
	case SyntaxKindPublicKeyword:
		return ModifierFlagsPublic
	case SyntaxKindProtectedKeyword:
		return ModifierFlagsProtected
	case SyntaxKindPrivateKeyword:
		return ModifierFlagsPrivate
	case SyntaxKindAbstractKeyword:
		return ModifierFlagsAbstract
	case SyntaxKindAccessorKeyword:
		return ModifierFlagsAccessor
	case SyntaxKindExportKeyword:
		return ModifierFlagsExport
	case SyntaxKindDeclareKeyword:
		return ModifierFlagsAmbient
	case SyntaxKindConstKeyword:
		return ModifierFlagsConst
	case SyntaxKindDefaultKeyword:
		return ModifierFlagsDefault
	case SyntaxKindAsyncKeyword:
		return ModifierFlagsAsync
	case SyntaxKindReadonlyKeyword:
		return ModifierFlagsReadonly
	case SyntaxKindOverrideKeyword:
		return ModifierFlagsOverride
	case SyntaxKindInKeyword:
		return ModifierFlagsIn
	case SyntaxKindOutKeyword:
		return ModifierFlagsOut
	case SyntaxKindImmediateKeyword:
		return ModifierFlagsImmediate
	case SyntaxKindDecorator:
		return ModifierFlagsDecorator
	}
	return ModifierFlagsNone
}

func modifiersToFlags(modifierList *Node) ModifierFlags {
	flags := ModifierFlagsNone
	if modifierList != nil {
		for _, modifier := range modifierList.AsModifierList().Modifiers {
			flags |= modifierToFlag(modifier.Kind)
		}
	}
	return flags
}

func nodeIsMissing(node *Node) bool {
	return node == nil || node.Loc.pos == node.Loc.end && node.Loc.pos >= 0 && node.Kind != SyntaxKindEndOfFile
}

func nodeIsPresent(node *Node) bool {
	return !nodeIsMissing(node)
}

func isLeftHandSideExpression(node *Node) bool {
	return isLeftHandSideExpressionKind(node.Kind)
}

func isLeftHandSideExpressionKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindPropertyAccessExpression, SyntaxKindElementAccessExpression, SyntaxKindNewExpression, SyntaxKindCallExpression,
		SyntaxKindJsxElement, SyntaxKindJsxSelfClosingElement, SyntaxKindJsxFragment, SyntaxKindTaggedTemplateExpression, SyntaxKindArrayLiteralExpression,
		SyntaxKindParenthesizedExpression, SyntaxKindObjectLiteralExpression, SyntaxKindClassExpression, SyntaxKindFunctionExpression, SyntaxKindIdentifier,
		SyntaxKindPrivateIdentifier, SyntaxKindRegularExpressionLiteral, SyntaxKindNumericLiteral, SyntaxKindBigIntLiteral, SyntaxKindStringLiteral,
		SyntaxKindNoSubstitutionTemplateLiteral, SyntaxKindTemplateExpression, SyntaxKindFalseKeyword, SyntaxKindNullKeyword, SyntaxKindThisKeyword,
		SyntaxKindTrueKeyword, SyntaxKindSuperKeyword, SyntaxKindNonNullExpression, SyntaxKindExpressionWithTypeArguments, SyntaxKindMetaProperty,
		SyntaxKindImportKeyword, SyntaxKindMissingDeclaration:
		return true
	}
	return false
}

func isUnaryExpression(node *Node) bool {
	return isUnaryExpressionKind(node.Kind)
}

func isUnaryExpressionKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindPrefixUnaryExpression, SyntaxKindPostfixUnaryExpression, SyntaxKindDeleteExpression, SyntaxKindTypeOfExpression,
		SyntaxKindVoidExpression, SyntaxKindAwaitExpression, SyntaxKindTypeAssertionExpression:
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

func isExpressionKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindConditionalExpression, SyntaxKindYieldExpression, SyntaxKindArrowFunction, SyntaxKindBinaryExpression,
		SyntaxKindSpreadElement, SyntaxKindAsExpression, SyntaxKindOmittedExpression, SyntaxKindCommaListExpression,
		SyntaxKindPartiallyEmittedExpression, SyntaxKindSatisfiesExpression:
		return true
	}
	return isUnaryExpressionKind(kind)
}

func isAssignmentOperator(token SyntaxKind) bool {
	return token >= SyntaxKindFirstAssignment && token <= SyntaxKindLastAssignment
}

func isExpressionWithTypeArguments(node *Node) bool {
	return node.Kind == SyntaxKindExpressionWithTypeArguments
}

func isNonNullExpression(node *Node) bool {
	return node.Kind == SyntaxKindNonNullExpression
}

func isStringLiteralLike(node *Node) bool {
	return node.Kind == SyntaxKindStringLiteral || node.Kind == SyntaxKindNoSubstitutionTemplateLiteral
}

func isNumericLiteral(node *Node) bool {
	return node.Kind == SyntaxKindNumericLiteral
}

func isStringOrNumericLiteralLike(node *Node) bool {
	return isStringLiteralLike(node) || isNumericLiteral(node)
}

func isSignedNumericLiteral(node *Node) bool {
	if node.Kind == SyntaxKindPrefixUnaryExpression {
		node := node.AsPrefixUnaryExpression()
		return (node.Operator == SyntaxKindPlusToken || node.Operator == SyntaxKindMinusToken) && isNumericLiteral(node.Operand)
	}
	return false
}

func ifElse[T any](b bool, whenTrue T, whenFalse T) T {
	if b {
		return whenTrue
	}
	return whenFalse
}

func tokenIsIdentifierOrKeyword(token SyntaxKind) bool {
	return token >= SyntaxKindIdentifier
}

func tokenIsIdentifierOrKeywordOrGreaterThan(token SyntaxKind) bool {
	return token == SyntaxKindGreaterThanToken || tokenIsIdentifierOrKeyword(token)
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
	return node.Kind == SyntaxKindBinaryExpression
}

func isAccessExpression(node *Node) bool {
	return node.Kind == SyntaxKindPropertyAccessExpression || node.Kind == SyntaxKindElementAccessExpression
}

func isInJSFile(node *Node) bool {
	return node != nil && node.Flags&NodeFlagsJavaScriptFile != 0
}

func isEffectiveModuleDeclaration(node *Node) bool {
	return IsModuleDeclaration(node) || IsIdentifier(node)
}

func isObjectLiteralOrClassExpressionMethodOrAccessor(node *Node) bool {
	kind := node.Kind
	return (kind == SyntaxKindMethodDeclaration || kind == SyntaxKindGetAccessor || kind == SyntaxKindSetAccessor) &&
		(node.Parent.Kind == SyntaxKindObjectLiteralExpression || node.Parent.Kind == SyntaxKindClassExpression)
}

func isFunctionLike(node *Node) bool {
	return node != nil && isFunctionLikeKind(node.Kind)
}

func isFunctionLikeKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindMethodSignature, SyntaxKindCallSignature, SyntaxKindJSDocSignature, SyntaxKindConstructSignature, SyntaxKindIndexSignature,
		SyntaxKindFunctionType, SyntaxKindJSDocFunctionType, SyntaxKindConstructorType:
		return true
	}
	return isFunctionLikeDeclarationKind(kind)
}

func isFunctionLikeDeclaration(node *Node) bool {
	return node != nil && isFunctionLikeDeclarationKind(node.Kind)
}

func isFunctionLikeDeclarationKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindFunctionDeclaration, SyntaxKindMethodDeclaration, SyntaxKindConstructor, SyntaxKindGetAccessor, SyntaxKindSetAccessor,
		SyntaxKindFunctionExpression, SyntaxKindArrowFunction:
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
	case SyntaxKindParenthesizedExpression:
		return kinds&OEKParentheses != 0 && !(kinds&OEKExcludeJSDocTypeAssertion != 0 && isJSDocTypeAssertion(node))
	case SyntaxKindTypeAssertionExpression, SyntaxKindAsExpression, SyntaxKindSatisfiesExpression:
		return kinds&OEKTypeAssertions != 0
	case SyntaxKindExpressionWithTypeArguments:
		return kinds&OEKExpressionsWithTypeArguments != 0
	case SyntaxKindNonNullExpression:
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
	for node != nil && node.Kind == SyntaxKindParenthesizedType {
		node = node.Parent
	}
	return node
}

func walkUpParenthesizedExpressions(node *Node) *Node {
	for node != nil && node.Kind == SyntaxKindParenthesizedExpression {
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
	case SyntaxKindPropertyDeclaration, SyntaxKindPropertySignature, SyntaxKindMethodDeclaration, SyntaxKindMethodSignature, SyntaxKindGetAccessor,
		SyntaxKindSetAccessor, SyntaxKindEnumMember, SyntaxKindPropertyAssignment, SyntaxKindPropertyAccessExpression:
		return parent.GetName() == node
	case SyntaxKindQualifiedName:
		return parent.AsQualifiedName().Right == node
	case SyntaxKindBindingElement:
		return parent.AsBindingElement().PropertyName == node
	case SyntaxKindImportSpecifier:
		return parent.AsImportSpecifier().PropertyName == node
	case SyntaxKindExportSpecifier, SyntaxKindJsxAttribute, SyntaxKindJsxSelfClosingElement, SyntaxKindJsxOpeningElement, SyntaxKindJsxClosingElement:
		return true
	}
	return false
}

func getSourceFileOfNode(node *Node) *SourceFile {
	for {
		if node == nil {
			return nil
		}
		if node.Kind == SyntaxKindSourceFile {
			return node.Data.(*SourceFile)
		}
		node = node.Parent
	}
}

/** @internal */
func getErrorRangeForNode(sourceFile *SourceFile, node *Node) TextRange {
	errorNode := node
	switch node.Kind {
	case SyntaxKindSourceFile:
		pos := skipTrivia(sourceFile.Text, 0)
		if pos == len(sourceFile.Text) {
			return NewTextRange(0, 0)
		}
		return getRangeOfTokenAtPosition(sourceFile, pos)
	// This list is a work in progress. Add missing node kinds to improve their error spans
	case SyntaxKindVariableDeclaration, SyntaxKindBindingElement, SyntaxKindClassDeclaration, SyntaxKindClassExpression, SyntaxKindInterfaceDeclaration,
		SyntaxKindModuleDeclaration, SyntaxKindEnumDeclaration, SyntaxKindEnumMember, SyntaxKindFunctionDeclaration, SyntaxKindFunctionExpression,
		SyntaxKindMethodDeclaration, SyntaxKindGetAccessor, SyntaxKindSetAccessor, SyntaxKindTypeAliasDeclaration, SyntaxKindPropertyDeclaration,
		SyntaxKindPropertySignature, SyntaxKindNamespaceImport:
		errorNode = getNameOfDeclaration(node)
	case SyntaxKindArrowFunction:
		return getErrorRangeForArrowFunction(sourceFile, node)
	case SyntaxKindCaseClause:
	case SyntaxKindDefaultClause:
		start := skipTrivia(sourceFile.Text, node.Pos())
		end := node.End()
		statements := node.Data.(*CaseOrDefaultClause).Statements
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return NewTextRange(start, end)
	case SyntaxKindReturnStatement, SyntaxKindYieldExpression:
		pos := skipTrivia(sourceFile.Text, node.Pos())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case SyntaxKindSatisfiesExpression:
		pos := skipTrivia(sourceFile.Text, node.AsSatisfiesExpression().Expression.End())
		return getRangeOfTokenAtPosition(sourceFile, pos)
	case SyntaxKindConstructor:
		scanner := getScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.tokenStart
		for scanner.token != SyntaxKindConstructorKeyword && scanner.token != SyntaxKindStringLiteral && scanner.token != SyntaxKindEndOfFile {
			scanner.Scan()
		}
		return NewTextRange(start, scanner.pos)
		// !!!
		// case SyntaxKindJSDocSatisfiesTag:
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
	return NewTextRange(pos, errorNode.End())
}

func getErrorRangeForArrowFunction(sourceFile *SourceFile, node *Node) TextRange {
	pos := skipTrivia(sourceFile.Text, node.Pos())
	body := node.AsArrowFunction().Body
	if body != nil && body.Kind == SyntaxKindBlock {
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
	return node != nil && (node.Kind == SyntaxKindClassDeclaration || node.Kind == SyntaxKindClassExpression)
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
		if (IsClassDeclaration(parent) || IsFunctionDeclaration(parent)) && parent.GetName() == node {
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
		case SyntaxKindComputedPropertyName:
			if includeClassComputedPropertyName && isClassLike(node.Parent.Parent) {
				return node
			}
			node = node.Parent.Parent
		case SyntaxKindDecorator:
			if node.Parent.Kind == SyntaxKindParameter && isClassElement(node.Parent.Parent) {
				// If the decorator's parent is a Parameter, we resolve the this container from
				// the grandparent class declaration.
				node = node.Parent.Parent
			} else if isClassElement(node.Parent) {
				// If the decorator's parent is a class element, we resolve the 'this' container
				// from the parent class declaration.
				node = node.Parent
			}
		case SyntaxKindArrowFunction:
			if includeArrowFunctions {
				return node
			}
		case SyntaxKindFunctionDeclaration, SyntaxKindFunctionExpression, SyntaxKindModuleDeclaration, SyntaxKindClassStaticBlockDeclaration,
			SyntaxKindPropertyDeclaration, SyntaxKindPropertySignature, SyntaxKindMethodDeclaration, SyntaxKindMethodSignature, SyntaxKindConstructor,
			SyntaxKindGetAccessor, SyntaxKindSetAccessor, SyntaxKindCallSignature, SyntaxKindConstructSignature, SyntaxKindIndexSignature,
			SyntaxKindEnumDeclaration, SyntaxKindSourceFile:
			return node
		}
	}
}

func isClassElement(node *Node) bool {
	switch node.Kind {
	case SyntaxKindConstructor, SyntaxKindPropertyDeclaration, SyntaxKindMethodDeclaration, SyntaxKindGetAccessor, SyntaxKindSetAccessor,
		SyntaxKindIndexSignature, SyntaxKindClassStaticBlockDeclaration, SyntaxKindSemicolonClassElement:
		return true
	}
	return false
}

func isPartOfTypeQuery(node *Node) bool {
	for node.Kind == SyntaxKindQualifiedName || node.Kind == SyntaxKindIdentifier {
		node = node.Parent
	}
	return node.Kind == SyntaxKindTypeQuery
}

func getModifierFlags(node *Node) ModifierFlags {
	modifiers := node.GetModifiers()
	if modifiers != nil {
		return modifiers.AsModifierList().ModifierFlags
	}
	return ModifierFlagsNone
}

func getNodeFlags(node *Node) NodeFlags {
	return node.Flags
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
	if fn.Kind == SyntaxKindFunctionExpression || fn.Kind == SyntaxKindArrowFunction {
		prev := fn
		parent := fn.Parent
		for parent.Kind == SyntaxKindParenthesizedExpression {
			prev = parent
			parent = parent.Parent
		}
		if parent.Kind == SyntaxKindCallExpression && parent.AsCallExpression().Expression == prev {
			return parent
		}
	}
	return nil
}

// Does not handle signed numeric names like `a[+0]` - handling those would require handling prefix unary expressions
// throughout late binding handling as well, which is awkward (but ultimately probably doable if there is demand)
func getElementOrPropertyAccessArgumentExpressionOrName(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().Name
	case SyntaxKindElementAccessExpression:
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
	case SyntaxKindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().QuestionDotToken
	case SyntaxKindElementAccessExpression:
		return node.AsElementAccessExpression().QuestionDotToken
	case SyntaxKindCallExpression:
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
	case SyntaxKindComputedPropertyName:
		expr = name.AsComputedPropertyName().Expression
	case SyntaxKindElementAccessExpression:
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
	case SyntaxKindBinaryExpression:
		if isFunctionPropertyAssignment(declaration) {
			return getElementOrPropertyAccessArgumentExpressionOrName(declaration.AsBinaryExpression().Left)
		}
		return nil
	case SyntaxKindExportAssignment:
		expr := declaration.AsExportAssignment().Expression
		if IsIdentifier(expr) {
			return expr
		}
		return nil
	}
	return declaration.GetName()
}

func getAssignedName(node *Node) *Node {
	parent := node.Parent
	if parent != nil {
		switch parent.Kind {
		case SyntaxKindPropertyAssignment:
			return parent.AsPropertyAssignment().Name
		case SyntaxKindBindingElement:
			return parent.AsBindingElement().Name
		case SyntaxKindBinaryExpression:
			if node == parent.AsBinaryExpression().Right {
				left := parent.AsBinaryExpression().Left
				switch left.Kind {
				case SyntaxKindIdentifier:
					return left
				case SyntaxKindPropertyAccessExpression:
					return left.AsPropertyAccessExpression().Name
				case SyntaxKindElementAccessExpression:
					arg := skipParentheses(left.AsElementAccessExpression().ArgumentExpression)
					if isStringOrNumericLiteralLike(arg) {
						return arg
					}
				}
			}
		case SyntaxKindVariableDeclaration:
			name := parent.AsVariableDeclaration().Name
			if IsIdentifier(name) {
				return name
			}
		}
	}
	return nil
}

func isFunctionPropertyAssignment(node *Node) bool {
	if node.Kind == SyntaxKindBinaryExpression {
		expr := node.AsBinaryExpression()
		if expr.OperatorToken.Kind == SyntaxKindEqualsToken {
			switch expr.Left.Kind {
			case SyntaxKindPropertyAccessExpression:
				// F.id = expr
				return IsIdentifier(expr.Left.AsPropertyAccessExpression().Expression) && IsIdentifier(expr.Left.AsPropertyAccessExpression().Name)
			case SyntaxKindElementAccessExpression:
				// F[xxx] = expr
				return IsIdentifier(expr.Left.AsElementAccessExpression().Expression)
			}
		}
	}
	return false
}

func isAssignmentExpression(node *Node, excludeCompoundAssignment bool) bool {
	if node.Kind == SyntaxKindBinaryExpression {
		expr := node.AsBinaryExpression()
		return (expr.OperatorToken.Kind == SyntaxKindEqualsToken || !excludeCompoundAssignment && isAssignmentOperator(expr.OperatorToken.Kind)) &&
			isLeftHandSideExpression(expr.Left)
	}
	return false
}

func isBlockOrCatchScoped(declaration *Node) bool {
	return getCombinedNodeFlags(declaration)&NodeFlagsBlockScoped != 0 || isCatchClauseVariableDeclarationOrBindingElement(declaration)
}

func isCatchClauseVariableDeclarationOrBindingElement(declaration *Node) bool {
	node := getRootDeclaration(declaration)
	return node.Kind == SyntaxKindVariableDeclaration && node.Parent.Kind == SyntaxKindCatchClause
}

func isAmbientModule(node *Node) bool {
	return IsModuleDeclaration(node) && (node.AsModuleDeclaration().Name.Kind == SyntaxKindStringLiteral || isGlobalScopeAugmentation(node))
}

func isGlobalScopeAugmentation(node *Node) bool {
	return node.Flags&NodeFlagsGlobalAugmentation != 0
}

func isPropertyNameLiteral(node *Node) bool {
	switch node.Kind {
	case SyntaxKindIdentifier, SyntaxKindStringLiteral, SyntaxKindNoSubstitutionTemplateLiteral, SyntaxKindNumericLiteral:
		return true
	}
	return false
}

func isMemberName(node *Node) bool {
	return node.Kind == SyntaxKindIdentifier || node.Kind == SyntaxKindPrivateIdentifier
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
	if node.Kind == SyntaxKindVariableDeclaration {
		node = node.Parent
	}
	if node != nil && node.Kind == SyntaxKindVariableDeclarationList {
		flags |= getFlags(node)
		node = node.Parent
	}
	if node != nil && node.Kind == SyntaxKindVariableStatement {
		flags |= getFlags(node)
	}
	return flags
}

func getCombinedModifierFlags(node *Node) ModifierFlags {
	return getCombinedFlags(node, getModifierFlags)
}

func getCombinedNodeFlags(node *Node) NodeFlags {
	return getCombinedFlags(node, getNodeFlags)
}

func isBindingPattern(node *Node) bool {
	return node != nil && (node.Kind == SyntaxKindArrayBindingPattern || node.Kind == SyntaxKindObjectBindingPattern)
}

func isParameterPropertyDeclaration(node *Node, parent *Node) bool {
	return IsParameter(node) && hasSyntacticModifier(node, ModifierFlagsParameterPropertyModifier) && parent.Kind == SyntaxKindConstructor
}

/**
 * Like {@link isVariableDeclarationInitializedToRequire} but allows things like `require("...").foo.bar` or `require("...")["baz"]`.
 */
func isVariableDeclarationInitializedToBareOrAccessedRequire(node *Node) bool {
	return isVariableDeclarationInitializedWithRequireHelper(node, true /*allowAccessedRequire*/)
}

func isVariableDeclarationInitializedWithRequireHelper(node *Node, allowAccessedRequire bool) bool {
	if node.Kind == SyntaxKindVariableDeclaration && node.AsVariableDeclaration().Initializer != nil {
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
	return getRootDeclaration(node).Kind == SyntaxKindParameter
}

func getRootDeclaration(node *Node) *Node {
	for node.Kind == SyntaxKindBindingElement {
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
	case SyntaxKindFunctionDeclaration, SyntaxKindFunctionExpression, SyntaxKindArrowFunction, SyntaxKindMethodDeclaration:
		data := node.BodyData()
		return data.Body != nil && data.AsteriskToken == nil && hasSyntacticModifier(node, ModifierFlagsAsync)
	}
	return false
}

func isObjectLiteralMethod(node *Node) bool {
	return node != nil && node.Kind == SyntaxKindMethodDeclaration && node.Parent.Kind == SyntaxKindObjectLiteralExpression
}

func symbolName(symbol *Symbol) string {
	if symbol.valueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(symbol.valueDeclaration) {
		return symbol.valueDeclaration.GetName().AsPrivateIdentifier().Text
	}
	return symbol.name
}

func isStaticPrivateIdentifierProperty(s *Symbol) bool {
	return s.valueDeclaration != nil && isPrivateIdentifierClassElementDeclaration(s.valueDeclaration) && isStatic(s.valueDeclaration)
}

func isPrivateIdentifierClassElementDeclaration(node *Node) bool {
	return (IsPropertyDeclaration(node) || isMethodOrAccessor(node)) && IsPrivateIdentifier(node.GetName())
}

func isMethodOrAccessor(node *Node) bool {
	switch node.Kind {
	case SyntaxKindMethodDeclaration, SyntaxKindGetAccessor, SyntaxKindSetAccessor:
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
	case SyntaxKindSourceFile:
		return isExternalModule(node.Parent.AsSourceFile())
	case SyntaxKindModuleBlock:
		grandParent := node.Parent.Parent
		return isAmbientModule(grandParent) && IsSourceFile(grandParent.Parent) && !isExternalModule(grandParent.Parent.AsSourceFile())
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
func isDeclarationStatementKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindFunctionDeclaration, SyntaxKindMissingDeclaration, SyntaxKindClassDeclaration, SyntaxKindInterfaceDeclaration,
		SyntaxKindTypeAliasDeclaration, SyntaxKindEnumDeclaration, SyntaxKindModuleDeclaration, SyntaxKindImportDeclaration,
		SyntaxKindImportEqualsDeclaration, SyntaxKindExportDeclaration, SyntaxKindExportAssignment, SyntaxKindNamespaceExportDeclaration:
		return true
	}
	return false
}

func isDeclarationStatement(node *Node) bool {
	return isDeclarationStatementKind(node.Kind)
}

func isStatementKindButNotDeclarationKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindBreakStatement, SyntaxKindContinueStatement, SyntaxKindDebuggerStatement, SyntaxKindDoStatement, SyntaxKindExpressionStatement,
		SyntaxKindEmptyStatement, SyntaxKindForInStatement, SyntaxKindForOfStatement, SyntaxKindForStatement, SyntaxKindIfStatement,
		SyntaxKindLabeledStatement, SyntaxKindReturnStatement, SyntaxKindSwitchStatement, SyntaxKindThrowStatement, SyntaxKindTryStatement,
		SyntaxKindVariableStatement, SyntaxKindWhileStatement, SyntaxKindWithStatement, SyntaxKindNotEmittedStatement:
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
	if node.Kind != SyntaxKindBlock {
		return false
	}
	if node.Parent != nil && (node.Parent.Kind == SyntaxKindTryStatement || node.Parent.Kind == SyntaxKindCatchClause) {
		return false
	}
	return !isFunctionBlock(node)
}

func isFunctionBlock(node *Node) bool {
	return node != nil && node.Kind == SyntaxKindBlock && isFunctionLike(node.Parent)
}

func shouldPreserveConstEnums(options *CompilerOptions) bool {
	return options.PreserveConstEnums == TSTrue || options.IsolatedModules == TSTrue
}

func exportAssignmentIsAlias(node *Node) bool {
	return isAliasableExpression(getExportAssignmentExpression(node))
}

func getExportAssignmentExpression(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindExportAssignment:
		return node.AsExportAssignment().Expression
	case SyntaxKindBinaryExpression:
		return node.AsBinaryExpression().Right
	}
	panic("Unhandled case in getExportAssignmentExpression")
}

func isAliasableExpression(e *Node) bool {
	return isEntityNameExpression(e) || IsClassExpression(e)
}

func isEmptyObjectLiteral(expression *Node) bool {
	return expression.Kind == SyntaxKindObjectLiteralExpression && len(expression.AsObjectLiteralExpression().Properties) == 0
}

func isFunctionSymbol(symbol *Symbol) bool {
	d := symbol.valueDeclaration
	return d != nil && (IsFunctionDeclaration(d) || IsVariableDeclaration(d) && isFunctionLike(d.AsVariableDeclaration().Initializer))
}

func isLogicalOrCoalescingAssignmentOperator(token SyntaxKind) bool {
	return token == SyntaxKindBarBarEqualsToken || token == SyntaxKindAmpersandAmpersandEqualsToken || token == SyntaxKindQuestionQuestionEqualsToken
}

func isLogicalOrCoalescingAssignmentExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingAssignmentOperator(expr.AsBinaryExpression().OperatorToken.Kind)
}

func isLogicalOrCoalescingBinaryOperator(token SyntaxKind) bool {
	return isBinaryLogicalOperator(token) || token == SyntaxKindQuestionQuestionToken
}

func isLogicalOrCoalescingBinaryExpression(expr *Node) bool {
	return isBinaryExpression(expr) && isLogicalOrCoalescingBinaryOperator(expr.AsBinaryExpression().OperatorToken.Kind)
}

func isBinaryLogicalOperator(token SyntaxKind) bool {
	return token == SyntaxKindBarBarToken || token == SyntaxKindAmpersandAmpersandToken
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
	return node.Kind == SyntaxKindBinaryExpression && node.AsBinaryExpression().OperatorToken.Kind == SyntaxKindQuestionQuestionToken
}

func isDottedName(node *Node) bool {
	switch node.Kind {
	case SyntaxKindIdentifier, SyntaxKindThisKeyword, SyntaxKindSuperKeyword, SyntaxKindMetaProperty:
		return true
	case SyntaxKindPropertyAccessExpression, SyntaxKindParenthesizedExpression:
		return isDottedName(node.Expression())
	}
	return false
}

func unusedLabelIsError(options *CompilerOptions) bool {
	return options.AllowUnusedLabels == TSFalse
}

func unreachableCodeIsError(options *CompilerOptions) bool {
	return options.AllowUnreachableCode == TSFalse
}

func isDestructuringAssignment(node *Node) bool {
	if isAssignmentExpression(node, true /*excludeCompoundAssignment*/) {
		kind := node.AsBinaryExpression().Left.Kind
		return kind == SyntaxKindObjectLiteralExpression || kind == SyntaxKindArrayLiteralExpression
	}
	return false
}

func isTopLevelLogicalExpression(node *Node) bool {
	for IsParenthesizedExpression(node.Parent) || IsPrefixUnaryExpression(node.Parent) && node.Parent.AsPrefixUnaryExpression().Operator == SyntaxKindExclamationToken {
		node = node.Parent
	}
	return !isStatementCondition(node) && !isLogicalExpression(node.Parent) && !(isOptionalChain(node.Parent) && node.Parent.Expression() == node)
}

func isStatementCondition(node *Node) bool {
	switch node.Parent.Kind {
	case SyntaxKindIfStatement:
		return node.Parent.AsIfStatement().Expression == node
	case SyntaxKindWhileStatement:
		return node.Parent.AsWhileStatement().Expression == node
	case SyntaxKindDoStatement:
		return node.Parent.AsDoStatement().Expression == node
	case SyntaxKindForStatement:
		return node.Parent.AsForStatement().Condition == node
	case SyntaxKindConditionalExpression:
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
	case SyntaxKindBinaryExpression:
		binaryOperator := target.AsBinaryExpression().OperatorToken.Kind
		if binaryOperator == SyntaxKindEqualsToken || isLogicalOrCoalescingAssignmentOperator(binaryOperator) {
			return AssignmentKindDefinite
		}
		return AssignmentKindCompound
	case SyntaxKindPrefixUnaryExpression, SyntaxKindPostfixUnaryExpression:
		return AssignmentKindCompound
	case SyntaxKindForInStatement, SyntaxKindForOfStatement:
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
		case SyntaxKindBinaryExpression:
			if isAssignmentOperator(parent.AsBinaryExpression().OperatorToken.Kind) && parent.AsBinaryExpression().Left == node {
				return parent
			}
			return nil
		case SyntaxKindPrefixUnaryExpression:
			if parent.AsPrefixUnaryExpression().Operator == SyntaxKindPlusPlusToken || parent.AsPrefixUnaryExpression().Operator == SyntaxKindMinusMinusToken {
				return parent
			}
			return nil
		case SyntaxKindPostfixUnaryExpression:
			if parent.AsPostfixUnaryExpression().Operator == SyntaxKindPlusPlusToken || parent.AsPostfixUnaryExpression().Operator == SyntaxKindMinusMinusToken {
				return parent
			}
			return nil
		case SyntaxKindForInStatement, SyntaxKindForOfStatement:
			if parent.AsForInOrOfStatement().Initializer == node {
				return parent
			}
			return nil
		case SyntaxKindParenthesizedExpression, SyntaxKindArrayLiteralExpression, SyntaxKindSpreadElement, SyntaxKindNonNullExpression:
			node = parent
		case SyntaxKindSpreadAssignment:
			node = parent.Parent
		case SyntaxKindShorthandPropertyAssignment:
			if parent.AsShorthandPropertyAssignment().Name != node {
				return nil
			}
			node = parent.Parent
		case SyntaxKindPropertyAssignment:
			if parent.AsPropertyAssignment().Name == node {
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
	return node != nil && node.Kind == SyntaxKindDeleteExpression
}

func isInCompoundLikeAssignment(node *Node) bool {
	target := getAssignmentTarget(node)
	return target != nil && isAssignmentExpression(target /*excludeCompoundAssignment*/, true) && isCompoundLikeAssignment(target)
}

func isCompoundLikeAssignment(assignment *Node) bool {
	right := skipParentheses(assignment.AsBinaryExpression().Right)
	return right.Kind == SyntaxKindBinaryExpression && isShiftOperatorOrHigher(right.AsBinaryExpression().OperatorToken.Kind)
}

func isPushOrUnshiftIdentifier(node *Node) bool {
	text := node.AsIdentifier().Text
	return text == "push" || text == "unshift"
}

func isBooleanLiteral(node *Node) bool {
	return node.Kind == SyntaxKindTrueKeyword || node.Kind == SyntaxKindFalseKeyword
}

func isOptionalChain(node *Node) bool {
	kind := node.Kind
	return node.Flags&NodeFlagsOptionalChain != 0 && (kind == SyntaxKindPropertyAccessExpression ||
		kind == SyntaxKindElementAccessExpression || kind == SyntaxKindCallExpression || kind == SyntaxKindNonNullExpression)
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
	return node.Kind == SyntaxKindIdentifier || isPropertyAccessEntityNameExpression(node)
}

func isPropertyAccessEntityNameExpression(node *Node) bool {
	if node.Kind == SyntaxKindPropertyAccessExpression {
		expr := node.AsPropertyAccessExpression()
		return expr.Name.Kind == SyntaxKindIdentifier && isEntityNameExpression(expr.Expression)
	}
	return false
}

func isPrologueDirective(node *Node) bool {
	return node.Kind == SyntaxKindExpressionStatement && node.AsExpressionStatement().Expression.Kind == SyntaxKindStringLiteral
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
	switch block.Kind {
	case SyntaxKindBlock:
		return block.AsBlock().Statements
	case SyntaxKindModuleBlock:
		return block.AsModuleBlock().Statements
	case SyntaxKindSourceFile:
		return block.AsSourceFile().Statements
	}
	panic("Unhandled case in getStatementsOfBlock")
}

func nodeHasName(statement *Node, id *Node) bool {
	name := statement.GetName()
	if name != nil {
		return IsIdentifier(name) && name.AsIdentifier().Text == id.AsIdentifier().Text
	}
	if IsVariableStatement(statement) {
		declarations := statement.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations
		return utils.Some(declarations, func(d *Node) bool { return nodeHasName(d, id) })
	}
	return false
}

func isImportMeta(node *Node) bool {
	if node.Kind == SyntaxKindMetaProperty {
		return node.AsMetaProperty().KeywordToken == SyntaxKindImportKeyword && node.AsMetaProperty().Name.AsIdentifier().Text == "meta"
	}
	return false
}

func ensureScriptKind(fileName string, scriptKind ScriptKind) ScriptKind {
	// Using scriptKind as a condition handles both:
	// - 'scriptKind' is unspecified and thus it is `undefined`
	// - 'scriptKind' is set and it is `Unknown` (0)
	// If the 'scriptKind' is 'undefined' or 'Unknown' then we attempt
	// to get the ScriptKind from the file name. If it cannot be resolved
	// from the file name then the default 'TS' script kind is returned.
	if scriptKind == ScriptKindUnknown {
		scriptKind = getScriptKindFromFileName(fileName)
	}
	if scriptKind == ScriptKindUnknown {
		scriptKind = ScriptKindTS
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

func getScriptKindFromFileName(fileName string) ScriptKind {
	dotPos := strings.LastIndex(fileName, ".")
	if dotPos >= 0 {
		switch strings.ToLower(fileName[dotPos:]) {
		case ExtensionJs, ExtensionCjs, ExtensionMjs:
			return ScriptKindJS
		case ExtensionJsx:
			return ScriptKindJSX
		case ExtensionTs, ExtensionCts, ExtensionMts:
			return ScriptKindTS
		case ExtensionTsx:
			return ScriptKindTSX
		case ExtensionJson:
			return ScriptKindJSON
		}
	}
	return ScriptKindUnknown
}

func getLanguageVariant(scriptKind ScriptKind) LanguageVariant {
	switch scriptKind {
	case ScriptKindTSX, ScriptKindJSX, ScriptKindJS, ScriptKindJSON:
		// .tsx and .jsx files are treated as jsx language variant.
		return LanguageVariantJSX
	}
	return LanguageVariantStandard
}

func getEmitScriptTarget(options *CompilerOptions) ScriptTarget {
	if options.Target != ScriptTargetNone {
		return options.Target
	}
	return ScriptTargetES5
}

func getEmitModuleKind(options *CompilerOptions) ModuleKind {
	if options.ModuleKind != ModuleKindNone {
		return options.ModuleKind
	}
	if options.Target >= ScriptTargetES2015 {
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
	if options.ESModuleInterop != TSUnknown {
		return options.ESModuleInterop == TSTrue
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
	if options.AllowSyntheticDefaultImports != TSUnknown {
		return options.AllowSyntheticDefaultImports == TSTrue
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
		c.fileDiagnostics[fileName] = utils.InsertSorted(c.fileDiagnostics[fileName], diagnostic, compareDiagnostics)
	} else {
		c.nonFileDiagnostics = utils.InsertSorted(c.nonFileDiagnostics, diagnostic, compareDiagnostics)
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
	switch location.Kind {
	case SyntaxKindAsExpression:
		return isConstTypeReference(location.AsAsExpression().TypeNode)
	case SyntaxKindTypeAssertionExpression:
		return isConstTypeReference(location.AsTypeAssertion().TypeNode)
	}
	return false
}

func isConstTypeReference(node *Node) bool {
	if node.Kind == SyntaxKindTypeReference {
		ref := node.AsTypeReference()
		return ref.TypeArguments == nil && IsIdentifier(ref.TypeName) && ref.TypeName.AsIdentifier().Text == "const"
	}
	return false
}

func isModuleOrEnumDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindModuleDeclaration || node.Kind == SyntaxKindEnumDeclaration
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
	return node.Kind == SyntaxKindSourceFile && !isExternalOrCommonJsModule(node.AsSourceFile())
}

func isParameterLikeOrReturnTag(node *Node) bool {
	switch node.Kind {
	case SyntaxKindParameter, SyntaxKindTypeParameter, SyntaxKindJSDocParameterTag, SyntaxKindJSDocReturnTag:
		return true
	}
	return false
}

func getEmitStandardClassFields(options *CompilerOptions) bool {
	return options.UseDefineForClassFields != TSFalse && getEmitScriptTarget(options) >= ScriptTargetES2022
}

func isTypeNodeKind(kind SyntaxKind) bool {
	switch kind {
	case SyntaxKindAnyKeyword, SyntaxKindUnknownKeyword, SyntaxKindNumberKeyword, SyntaxKindBigIntKeyword, SyntaxKindObjectKeyword,
		SyntaxKindBooleanKeyword, SyntaxKindStringKeyword, SyntaxKindSymbolKeyword, SyntaxKindVoidKeyword, SyntaxKindUndefinedKeyword,
		SyntaxKindNeverKeyword, SyntaxKindIntrinsicKeyword, SyntaxKindExpressionWithTypeArguments, SyntaxKindJSDocAllType, SyntaxKindJSDocUnknownType,
		SyntaxKindJSDocNullableType, SyntaxKindJSDocNonNullableType, SyntaxKindJSDocOptionalType, SyntaxKindJSDocFunctionType, SyntaxKindJSDocVariadicType:
		return true
	}
	return kind >= SyntaxKindFirstTypeNode && kind <= SyntaxKindLastTypeNode
}

func isTypeNode(node *Node) bool {
	return isTypeNodeKind(node.Kind)
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

func getDeclarationOfKind(symbol *Symbol, kind SyntaxKind) *Node {
	for _, declaration := range symbol.declarations {
		if declaration.Kind == kind {
			return declaration
		}
	}
	return nil
}

func getIsolatedModules(options *CompilerOptions) bool {
	return options.IsolatedModules == TSTrue || options.VerbatimModuleSyntax == TSTrue
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
	lookup                           func(symbols SymbolTable, name string, meaning SymbolFlags) *Symbol
	setRequiresScopeChangeCache      func(node *Node, value Tristate)
	getRequiresScopeChangeCache      func(node *Node) Tristate
	onPropertyWithInvalidInitializer func(location *Node, name string, declaration *Node, result *Symbol) bool
	onFailedToResolveSymbol          func(location *Node, name string, meaning SymbolFlags, nameNotFoundMessage *diagnostics.Message)
	onSuccessfullyResolvedSymbol     func(location *Node, result *Symbol, meaning SymbolFlags, lastLocation *Node, associatedDeclarationForContainingInitializerOrBindingName *Node, withinDeferredContext bool)
}

func (r *NameResolver) resolve(location *Node, name string, meaning SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *Symbol {
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
		if isModuleOrEnumDeclaration(location) && lastLocation != nil && location.GetName() == lastLocation {
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
					if meaning&result.flags&SymbolFlagsType != 0 && lastLocation.Kind != SyntaxKindJSDoc {
						useResult = result.flags&SymbolFlagsTypeParameter != 0 && (lastLocation.Flags&NodeFlagsSynthesized != 0 ||
							lastLocation == location.ReturnType() ||
							isParameterLikeOrReturnTag(lastLocation))
					}
					if meaning&result.flags&SymbolFlagsVariable != 0 {
						// expression inside parameter will lookup as normal variable scope when targeting es2015+
						if r.useOuterVariableScopeInParameter(result, location, lastLocation) {
							useResult = false
						} else if result.flags&SymbolFlagsFunctionScopedVariable != 0 {
							// parameters are visible only inside function body, parameter list and return type
							// technically for parameter list case here we might mix parameters and variables declared in function,
							// however it is detected separately when checking initializers of parameters
							// to make sure that they reference no variables declared after them.
							useResult = lastLocation.Kind == SyntaxKindParameter ||
								lastLocation.Flags&NodeFlagsSynthesized != 0 ||
								lastLocation == location.ReturnType() && findAncestor(result.valueDeclaration, IsParameter) != nil
						}
					}
				} else if location.Kind == SyntaxKindConditionalType {
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
		case SyntaxKindSourceFile:
			if !isExternalOrCommonJsModule(location.AsSourceFile()) {
				break
			}
			fallthrough
		case SyntaxKindModuleDeclaration:
			moduleExports := r.getSymbolOfDeclaration(location).exports
			if IsSourceFile(location) || (IsModuleDeclaration(location) && location.Flags&NodeFlagsAmbient != 0 && !isGlobalScopeAugmentation(location)) {
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
				//     2. We check === SymbolFlags.Alias in order to check that the symbol is *purely*
				//        an alias. If we used &, we'd be throwing out symbols that have non alias aspects,
				//        which is not the desired behavior.
				moduleExport := moduleExports[name]
				if moduleExport != nil && moduleExport.flags == SymbolFlagsAlias && (getDeclarationOfKind(moduleExport, SyntaxKindExportSpecifier) != nil || getDeclarationOfKind(moduleExport, SyntaxKindNamespaceExport) != nil) {
					break
				}
			}
			if name != InternalSymbolNameDefault {
				result = r.lookup(moduleExports, name, meaning&SymbolFlagsModuleMember)
				if result != nil {
					break loop
				}
			}
		case SyntaxKindEnumDeclaration:
			result = r.lookup(r.getSymbolOfDeclaration(location).exports, name, meaning&SymbolFlagsEnumMember)
			if result != nil {
				if nameNotFoundMessage != nil && getIsolatedModules(r.compilerOptions) && location.Flags&NodeFlagsAmbient == 0 && getSourceFileOfNode(location) != getSourceFileOfNode(result.valueDeclaration) {
					isolatedModulesLikeFlagName := ifElse(r.compilerOptions.VerbatimModuleSyntax == TSTrue, "verbatimModuleSyntax", "isolatedModules")
					r.error(originalLocation, diagnostics.Cannot_access_0_from_another_file_without_qualification_when_1_is_enabled_Use_2_instead,
						name, isolatedModulesLikeFlagName, r.getSymbolOfDeclaration(location).name+"."+name)
				}
				break loop
			}
		case SyntaxKindPropertyDeclaration:
			if !isStatic(location) {
				ctor := findConstructorDeclaration(location.Parent)
				if ctor != nil && ctor.AsConstructorDeclaration().locals != nil {
					if r.lookup(ctor.AsConstructorDeclaration().locals, name, meaning&SymbolFlagsValue) != nil {
						// Remember the property node, it will be used later to report appropriate error
						propertyWithInvalidInitializer = location
					}
				}
			}
		case SyntaxKindClassDeclaration, SyntaxKindClassExpression, SyntaxKindInterfaceDeclaration:
			result = r.lookup(r.getSymbolOfDeclaration(location).members, name, meaning&SymbolFlagsType)
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
			if IsClassExpression(location) && meaning&SymbolFlagsClass != 0 {
				className := location.AsClassExpression().Name
				if className != nil && name == className.AsIdentifier().Text {
					result = location.AsClassExpression().Symbol
					break loop
				}
			}
		case SyntaxKindExpressionWithTypeArguments:
			if lastLocation == location.AsExpressionWithTypeArguments().Expression && IsHeritageClause(location.Parent) && location.Parent.AsHeritageClause().Token == SyntaxKindExtendsKeyword {
				container := location.Parent.Parent
				if isClassLike(container) {
					result = r.lookup(r.getSymbolOfDeclaration(container).members, name, meaning&SymbolFlagsType)
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
		case SyntaxKindComputedPropertyName:
			grandparent = location.Parent.Parent
			if isClassLike(grandparent) || IsInterfaceDeclaration(grandparent) {
				// A reference to this grandparent's type parameters would be an error
				result = r.lookup(r.getSymbolOfDeclaration(grandparent).members, name, meaning&SymbolFlagsType)
				if result != nil {
					if nameNotFoundMessage != nil {
						r.error(originalLocation, diagnostics.A_computed_property_name_cannot_reference_a_type_parameter_from_its_containing_type)
					}
					return nil
				}
			}
		case SyntaxKindArrowFunction:
			// when targeting ES6 or higher there is no 'arguments' in an arrow function
			// for lower compile targets the resolved symbol is used to emit an error
			if getEmitScriptTarget(r.compilerOptions) >= ScriptTargetES2015 {
				break
			}
			fallthrough
		case SyntaxKindMethodDeclaration, SyntaxKindConstructor, SyntaxKindGetAccessor, SyntaxKindSetAccessor, SyntaxKindFunctionDeclaration:
			if meaning&SymbolFlagsVariable != 0 && name == "arguments" {
				result = r.argumentsSymbol
				break loop
			}
		case SyntaxKindFunctionExpression:
			if meaning&SymbolFlagsVariable != 0 && name == "arguments" {
				result = r.argumentsSymbol
				break loop
			}
			if meaning&SymbolFlagsFunction != 0 {
				functionName := location.AsFunctionExpression().Name
				if functionName != nil && name == functionName.AsIdentifier().Text {
					result = location.AsFunctionExpression().Symbol
					break loop
				}
			}
		case SyntaxKindDecorator:
			// Decorators are resolved at the class declaration. Resolving at the parameter
			// or member would result in looking up locals in the method.
			//
			//   function y() {}
			//   class C {
			//       method(@y x, y) {} // <-- decorator y should be resolved at the class declaration, not the parameter.
			//   }
			//
			if location.Parent != nil && location.Parent.Kind == SyntaxKindParameter {
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
			if location.Parent != nil && (isClassElement(location.Parent) || location.Parent.Kind == SyntaxKindClassDeclaration) {
				location = location.Parent
			}
		case SyntaxKindParameter:
			parameterDeclaration := location.AsParameterDeclaration()
			if lastLocation != nil && (lastLocation == parameterDeclaration.Initializer ||
				lastLocation == parameterDeclaration.Name && isBindingPattern(lastLocation)) {
				if associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case SyntaxKindBindingElement:
			bindingElement := location.AsBindingElement()
			if lastLocation != nil && (lastLocation == bindingElement.Initializer ||
				lastLocation == bindingElement.Name && isBindingPattern(lastLocation)) {
				if isPartOfParameterDeclaration(location) && associatedDeclarationForContainingInitializerOrBindingName == nil {
					associatedDeclarationForContainingInitializerOrBindingName = location
				}
			}
		case SyntaxKindInferType:
			if meaning&SymbolFlagsTypeParameter != 0 {
				parameterName := location.AsInferTypeNode().TypeParameter.AsTypeParameter().Name
				if parameterName != nil && name == parameterName.AsIdentifier().Text {
					result = location.AsInferTypeNode().TypeParameter.AsTypeParameter().Symbol
					break loop
				}
			}
		case SyntaxKindExportSpecifier:
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
		if body != nil && result.valueDeclaration != nil && result.valueDeclaration.Pos() >= body.Pos() && result.valueDeclaration.End() <= body.End() {
			// check for several cases where we introduce temporaries that require moving the name/initializer of the parameter to the body
			// - static field in a class expression
			// - optional chaining pre-es2020
			// - nullish coalesce pre-es2020
			// - spread assignment in binding pattern pre-es2017
			target := getEmitScriptTarget(r.compilerOptions)
			if target >= ScriptTargetES2015 {
				functionLocation := location
				declarationRequiresScopeChange := r.getRequiresScopeChangeCache(functionLocation)
				if declarationRequiresScopeChange == TSUnknown {
					declarationRequiresScopeChange = boolToTristate(utils.Some(functionLocation.Parameters(), r.requiresScopeChange))
					r.setRequiresScopeChangeCache(functionLocation, declarationRequiresScopeChange)
				}
				return declarationRequiresScopeChange == TSTrue
			}
		}
	}
	return false
}

func (r *NameResolver) requiresScopeChange(node *Node) bool {
	d := node.AsParameterDeclaration()
	return r.requiresScopeChangeWorker(d.Name) || d.Initializer != nil && r.requiresScopeChangeWorker(d.Initializer)
}

func (r *NameResolver) requiresScopeChangeWorker(node *Node) bool {
	switch node.Kind {
	case SyntaxKindArrowFunction, SyntaxKindFunctionExpression, SyntaxKindFunctionDeclaration, SyntaxKindConstructor:
		return false
	case SyntaxKindMethodDeclaration, SyntaxKindGetAccessor, SyntaxKindSetAccessor, SyntaxKindPropertyAssignment:
		return r.requiresScopeChangeWorker(node.GetName())
	case SyntaxKindPropertyDeclaration:
		if hasStaticModifier(node) {
			return !getEmitStandardClassFields(r.compilerOptions)
		}
		return r.requiresScopeChangeWorker(node.AsPropertyDeclaration().Name)
	default:
		if isNullishCoalesce(node) || isOptionalChain(node) {
			return getEmitScriptTarget(r.compilerOptions) < ScriptTargetES2020
		}
		if IsBindingElement(node) && node.AsBindingElement().DotDotDotToken != nil && IsObjectBindingPattern(node.Parent) {
			return getEmitScriptTarget(r.compilerOptions) < ScriptTargetES2017
		}
		if isTypeNode(node) {
			return false
		}
		return node.ForEachChild(r.requiresScopeChangeWorker)
	}
}

func getIsDeferredContext(location *Node, lastLocation *Node) bool {
	if location.Kind != SyntaxKindArrowFunction && location.Kind != SyntaxKindFunctionExpression {
		// initializers in instance property declaration of class like entities are executed in constructor and thus deferred
		// A name is evaluated within the enclosing scope - so it shouldn't count as deferred
		return IsTypeQueryNode(location) ||
			(isFunctionLikeDeclaration(location) || location.Kind == SyntaxKindPropertyDeclaration && !isStatic(location)) &&
				(lastLocation == nil || lastLocation != location.GetName())
	}
	if lastLocation != nil && lastLocation == location.GetName() {
		return false
	}
	// generator functions and async functions are not inlined in control flow when immediately invoked
	if location.BodyData().AsteriskToken != nil || hasSyntacticModifier(location, ModifierFlagsAsync) {
		return true
	}
	return getImmediatelyInvokedFunctionExpression(location) == nil
}

func isTypeParameterSymbolDeclaredInContainer(symbol *Symbol, container *Node) bool {
	for _, decl := range symbol.declarations {
		if decl.Kind == SyntaxKindTypeParameter {
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
	case SyntaxKindParameter:
		return lastLocation != nil && lastLocation == node.AsParameterDeclaration().Name
	case SyntaxKindFunctionDeclaration, SyntaxKindClassDeclaration, SyntaxKindInterfaceDeclaration, SyntaxKindEnumDeclaration,
		SyntaxKindTypeAliasDeclaration, SyntaxKindModuleDeclaration: // For `namespace N { N; }`
		return true
	}
	return false
}

func isTypeReferenceIdentifier(node *Node) bool {
	for node.Parent.Kind == SyntaxKindQualifiedName {
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
		case SyntaxKindTypeQuery:
			return FindAncestorTrue
		case SyntaxKindIdentifier, SyntaxKindQualifiedName:
			return FindAncestorFalse
		}
		return FindAncestorQuit
	}) != nil
}

func nodeKindIs(node *Node, kinds ...SyntaxKind) bool {
	return slices.Contains(kinds, node.Kind)
}

func isTypeOnlyImportDeclaration(node *Node) bool {
	switch node.Kind {
	case SyntaxKindImportSpecifier:
		return node.AsImportSpecifier().IsTypeOnly || node.Parent.Parent.AsImportClause().IsTypeOnly
	case SyntaxKindNamespaceImport:
		return node.Parent.AsImportClause().IsTypeOnly
	case SyntaxKindImportClause:
		return node.AsImportClause().IsTypeOnly
	case SyntaxKindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().IsTypeOnly
	}
	return false
}

func isTypeOnlyExportDeclaration(node *Node) bool {
	switch node.Kind {
	case SyntaxKindExportSpecifier:
		return node.AsExportSpecifier().IsTypeOnly || node.Parent.Parent.AsExportDeclaration().IsTypeOnly
	case SyntaxKindExportDeclaration:
		d := node.AsExportDeclaration()
		return d.IsTypeOnly && d.ModuleSpecifier != nil && d.ExportClause == nil
	case SyntaxKindNamespaceExport:
		return node.Parent.AsExportDeclaration().IsTypeOnly
	}
	return false
}

func isTypeOnlyImportOrExportDeclaration(node *Node) bool {
	return isTypeOnlyImportDeclaration(node) || isTypeOnlyExportDeclaration(node)
}

func getNameFromImportDeclaration(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindImportSpecifier:
		return node.AsImportSpecifier().Name
	case SyntaxKindNamespaceImport:
		return node.AsNamespaceImport().Name
	case SyntaxKindImportClause:
		return node.AsImportClause().Name
	case SyntaxKindImportEqualsDeclaration:
		return node.AsImportEqualsDeclaration().Name
	}
	return nil
}

func isValidTypeOnlyAliasUseSite(useSite *Node) bool {
	return useSite.Flags&NodeFlagsAmbient != 0 ||
		isPartOfTypeQuery(useSite) ||
		isIdentifierInNonEmittingHeritageClause(useSite) ||
		isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(useSite) ||
		!(isExpressionNode(useSite) || isShorthandPropertyNameUseSite(useSite))
}

func isIdentifierInNonEmittingHeritageClause(node *Node) bool {
	if node.Kind != SyntaxKindIdentifier {
		return false
	}
	heritageClause := findAncestorOrQuit(node.Parent, func(parent *Node) FindAncestorResult {
		switch parent.Kind {
		case SyntaxKindHeritageClause:
			return FindAncestorTrue
		case SyntaxKindPropertyAccessExpression, SyntaxKindExpressionWithTypeArguments:
			return FindAncestorFalse
		default:
			return FindAncestorQuit
		}
	})
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().Token == SyntaxKindImmediateKeyword || heritageClause.Parent.Kind == SyntaxKindInterfaceDeclaration
	}
	return false
}

func isPartOfPossiblyValidTypeOrAbstractComputedPropertyName(node *Node) bool {
	for nodeKindIs(node, SyntaxKindIdentifier, SyntaxKindPropertyAccessExpression) {
		node = node.Parent
	}
	if node.Kind != SyntaxKindComputedPropertyName {
		return false
	}
	if hasSyntacticModifier(node.Parent, ModifierFlagsAbstract) {
		return true
	}
	return nodeKindIs(node.Parent.Parent, SyntaxKindInterfaceDeclaration, SyntaxKindTypeLiteral)
}

func isExpressionNode(node *Node) bool {
	switch node.Kind {
	case SyntaxKindSuperKeyword, SyntaxKindNullKeyword, SyntaxKindTrueKeyword, SyntaxKindFalseKeyword, SyntaxKindRegularExpressionLiteral,
		SyntaxKindArrayLiteralExpression, SyntaxKindObjectLiteralExpression, SyntaxKindPropertyAccessExpression, SyntaxKindElementAccessExpression,
		SyntaxKindCallExpression, SyntaxKindNewExpression, SyntaxKindTaggedTemplateExpression, SyntaxKindAsExpression, SyntaxKindTypeAssertionExpression,
		SyntaxKindSatisfiesExpression, SyntaxKindNonNullExpression, SyntaxKindParenthesizedExpression, SyntaxKindFunctionExpression,
		SyntaxKindClassExpression, SyntaxKindArrowFunction, SyntaxKindVoidExpression, SyntaxKindDeleteExpression, SyntaxKindTypeOfExpression,
		SyntaxKindPrefixUnaryExpression, SyntaxKindPostfixUnaryExpression, SyntaxKindBinaryExpression, SyntaxKindConditionalExpression,
		SyntaxKindSpreadElement, SyntaxKindTemplateExpression, SyntaxKindOmittedExpression, SyntaxKindJsxElement, SyntaxKindJsxSelfClosingElement,
		SyntaxKindJsxFragment, SyntaxKindYieldExpression, SyntaxKindAwaitExpression, SyntaxKindMetaProperty:
		return true
	case SyntaxKindExpressionWithTypeArguments:
		return !IsHeritageClause(node.Parent)
	case SyntaxKindQualifiedName:
		for node.Parent.Kind == SyntaxKindQualifiedName {
			node = node.Parent
		}
		return IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node)
	case SyntaxKindJSDocMemberName:
		return IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node)
	case SyntaxKindPrivateIdentifier:
		return isBinaryExpression(node.Parent) && node.Parent.AsBinaryExpression().Left == node && node.Parent.AsBinaryExpression().OperatorToken.Kind == SyntaxKindInKeyword
	case SyntaxKindIdentifier:
		if IsTypeQueryNode(node.Parent) || isJSDocLinkLike(node.Parent) || isJSXTagName(node) {
			return true
		}
		fallthrough
	case SyntaxKindNumericLiteral, SyntaxKindBigIntLiteral, SyntaxKindStringLiteral, SyntaxKindNoSubstitutionTemplateLiteral, SyntaxKindThisKeyword:
		return isInExpressionContext(node)
	default:
		return false
	}
}

func isInExpressionContext(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case SyntaxKindVariableDeclaration:
		return parent.AsVariableDeclaration().Initializer == node
	case SyntaxKindParameter:
		return parent.AsParameterDeclaration().Initializer == node
	case SyntaxKindPropertyDeclaration:
		return parent.AsPropertyDeclaration().Initializer == node
	case SyntaxKindPropertySignature:
		return parent.AsPropertySignatureDeclaration().Initializer == node
	case SyntaxKindEnumMember:
		return parent.AsEnumMember().Initializer == node
	case SyntaxKindPropertyAssignment:
		return parent.AsPropertyAssignment().Initializer == node
	case SyntaxKindBindingElement:
		return parent.AsBindingElement().Initializer == node
	case SyntaxKindExpressionStatement:
		return parent.AsExpressionStatement().Expression == node
	case SyntaxKindIfStatement:
		return parent.AsIfStatement().Expression == node
	case SyntaxKindDoStatement:
		return parent.AsDoStatement().Expression == node
	case SyntaxKindWhileStatement:
		return parent.AsWhileStatement().Expression == node
	case SyntaxKindReturnStatement:
		return parent.AsReturnStatement().Expression == node
	case SyntaxKindWithStatement:
		return parent.AsWithStatement().Expression == node
	case SyntaxKindSwitchStatement:
		return parent.AsSwitchStatement().Expression == node
	case SyntaxKindCaseClause, SyntaxKindDefaultClause:
		return parent.AsCaseOrDefaultClause().Expression == node
	case SyntaxKindThrowStatement:
		return parent.AsThrowStatement().Expression == node
	case SyntaxKindForStatement:
		s := parent.AsForStatement()
		return s.Initializer == node && s.Initializer.Kind != SyntaxKindVariableDeclarationList || s.Condition == node || s.Incrementor == node
	case SyntaxKindForInStatement, SyntaxKindForOfStatement:
		s := parent.AsForInOrOfStatement()
		return s.Initializer == node && s.Initializer.Kind != SyntaxKindVariableDeclarationList || s.Expression == node
	case SyntaxKindTypeAssertionExpression:
		return parent.AsTypeAssertion().Expression == node
	case SyntaxKindAsExpression:
		return parent.AsAsExpression().Expression == node
	case SyntaxKindTemplateSpan:
		return parent.AsTemplateSpan().Expression == node
	case SyntaxKindComputedPropertyName:
		return parent.AsComputedPropertyName().Expression == node
	case SyntaxKindDecorator, SyntaxKindJsxExpression, SyntaxKindJsxSpreadAttribute, SyntaxKindSpreadAssignment:
		return true
	case SyntaxKindExpressionWithTypeArguments:
		return parent.AsExpressionWithTypeArguments().Expression == node && !isPartOfTypeNode(parent)
	case SyntaxKindShorthandPropertyAssignment:
		return parent.AsShorthandPropertyAssignment().ObjectAssignmentInitializer == node
	case SyntaxKindSatisfiesExpression:
		return parent.AsSatisfiesExpression().Expression == node
	default:
		return isExpressionNode(parent)
	}
}

func isPartOfTypeNode(node *Node) bool {
	kind := node.Kind
	if kind >= SyntaxKindFirstTypeNode && kind <= SyntaxKindLastTypeNode {
		return true
	}
	switch node.Kind {
	case SyntaxKindAnyKeyword, SyntaxKindUnknownKeyword, SyntaxKindNumberKeyword, SyntaxKindBigIntKeyword, SyntaxKindStringKeyword,
		SyntaxKindBooleanKeyword, SyntaxKindSymbolKeyword, SyntaxKindObjectKeyword, SyntaxKindUndefinedKeyword, SyntaxKindNullKeyword,
		SyntaxKindNeverKeyword:
		return true
	case SyntaxKindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(node)
	case SyntaxKindTypeParameter:
		return node.Parent.Kind == SyntaxKindMappedType || node.Parent.Kind == SyntaxKindInferType
	case SyntaxKindIdentifier:
		parent := node.Parent
		if IsQualifiedName(parent) && parent.AsQualifiedName().Right == node {
			return isPartOfTypeNodeInParent(parent)
		}
		if IsPropertyAccessExpression(parent) && parent.AsPropertyAccessExpression().Name == node {
			return isPartOfTypeNodeInParent(parent)
		}
		return isPartOfTypeNodeInParent(node)
	case SyntaxKindQualifiedName, SyntaxKindPropertyAccessExpression, SyntaxKindThisKeyword:
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
	if parent.Kind >= SyntaxKindFirstTypeNode && parent.Kind <= SyntaxKindLastTypeNode {
		return true
	}
	switch parent.Kind {
	case SyntaxKindTypeQuery:
		return false
	case SyntaxKindImportType:
		return !parent.AsImportTypeNode().IsTypeOf
	case SyntaxKindExpressionWithTypeArguments:
		return isPartOfTypeExpressionWithTypeArguments(parent)
	case SyntaxKindTypeParameter:
		return node == parent.AsTypeParameter().Constraint
	case SyntaxKindPropertyDeclaration:
		return node == parent.AsPropertyDeclaration().TypeNode
	case SyntaxKindPropertySignature:
		return node == parent.AsPropertySignatureDeclaration().TypeNode
	case SyntaxKindParameter:
		return node == parent.AsParameterDeclaration().TypeNode
	case SyntaxKindVariableDeclaration:
		return node == parent.AsVariableDeclaration().TypeNode
	case SyntaxKindFunctionDeclaration, SyntaxKindFunctionExpression, SyntaxKindArrowFunction, SyntaxKindConstructor, SyntaxKindMethodDeclaration,
		SyntaxKindMethodSignature, SyntaxKindGetAccessor, SyntaxKindSetAccessor, SyntaxKindCallSignature, SyntaxKindConstructSignature,
		SyntaxKindIndexSignature:
		return node == parent.ReturnType()
	case SyntaxKindTypeAssertionExpression:
		return node == parent.AsTypeAssertion().TypeNode
	case SyntaxKindCallExpression:
		return typeArgumentListContains(parent.AsCallExpression().TypeArguments, node)
	case SyntaxKindNewExpression:
		return typeArgumentListContains(parent.AsNewExpression().TypeArguments, node)
	case SyntaxKindTaggedTemplateExpression:
		return typeArgumentListContains(parent.AsTaggedTemplateExpression().TypeArguments, node)
	}
	return false
}

func isPartOfTypeExpressionWithTypeArguments(node *Node) bool {
	parent := node.Parent
	return IsHeritageClause(parent) && (!isClassLike(parent.Parent) || parent.AsHeritageClause().Token == SyntaxKindImplementsKeyword)
}

func typeArgumentListContains(list *Node, node *Node) bool {
	if list != nil {
		return slices.Contains(list.AsTypeArgumentList().Arguments, node)
	}
	return false
}

func isJSDocLinkLike(node *Node) bool {
	return nodeKindIs(node, SyntaxKindJSDocLink, SyntaxKindJSDocLinkCode, SyntaxKindJSDocLinkPlain)
}

func isJSXTagName(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case SyntaxKindJsxOpeningElement:
		return parent.AsJsxOpeningElement().TagName == node
	case SyntaxKindJsxSelfClosingElement:
		return parent.AsJsxSelfClosingElement().TagName == node
	case SyntaxKindJsxClosingElement:
		return parent.AsJsxClosingElement().TagName == node
	}
	return false
}

func isShorthandPropertyNameUseSite(useSite *Node) bool {
	return IsIdentifier(useSite) && IsShorthandPropertyAssignment(useSite.Parent) && useSite.Parent.AsShorthandPropertyAssignment().Name == useSite
}

func isTypeDeclaration(node *Node) bool {
	switch node.Kind {
	case SyntaxKindTypeParameter, SyntaxKindClassDeclaration, SyntaxKindInterfaceDeclaration, SyntaxKindTypeAliasDeclaration, SyntaxKindEnumDeclaration:
		return true
	case SyntaxKindImportClause:
		return node.AsImportClause().IsTypeOnly
	case SyntaxKindImportSpecifier:
		return node.Parent.Parent.AsImportClause().IsTypeOnly
	case SyntaxKindExportSpecifier:
		return node.Parent.Parent.AsExportDeclaration().IsTypeOnly
	default:
		return false
	}
}

func canHaveSymbol(node *Node) bool {
	switch node.Kind {
	case SyntaxKindArrowFunction, SyntaxKindBinaryExpression, SyntaxKindBindingElement, SyntaxKindCallExpression, SyntaxKindCallSignature,
		SyntaxKindClassDeclaration, SyntaxKindClassExpression, SyntaxKindClassStaticBlockDeclaration, SyntaxKindConstructor, SyntaxKindConstructorType,
		SyntaxKindConstructSignature, SyntaxKindElementAccessExpression, SyntaxKindEnumDeclaration, SyntaxKindEnumMember, SyntaxKindExportAssignment,
		SyntaxKindExportDeclaration, SyntaxKindExportSpecifier, SyntaxKindFunctionDeclaration, SyntaxKindFunctionExpression, SyntaxKindFunctionType,
		SyntaxKindGetAccessor, SyntaxKindIdentifier, SyntaxKindImportClause, SyntaxKindImportEqualsDeclaration, SyntaxKindImportSpecifier,
		SyntaxKindIndexSignature, SyntaxKindInterfaceDeclaration, SyntaxKindJSDocCallbackTag, SyntaxKindJSDocEnumTag, SyntaxKindJSDocFunctionType,
		SyntaxKindJSDocParameterTag, SyntaxKindJSDocPropertyTag, SyntaxKindJSDocSignature, SyntaxKindJSDocTypedefTag, SyntaxKindJSDocTypeLiteral,
		SyntaxKindJsxAttribute, SyntaxKindJsxAttributes, SyntaxKindJsxSpreadAttribute, SyntaxKindMappedType, SyntaxKindMethodDeclaration,
		SyntaxKindMethodSignature, SyntaxKindModuleDeclaration, SyntaxKindNamedTupleMember, SyntaxKindNamespaceExport, SyntaxKindNamespaceExportDeclaration,
		SyntaxKindNamespaceImport, SyntaxKindNewExpression, SyntaxKindNoSubstitutionTemplateLiteral, SyntaxKindNumericLiteral, SyntaxKindObjectLiteralExpression,
		SyntaxKindParameter, SyntaxKindPropertyAccessExpression, SyntaxKindPropertyAssignment, SyntaxKindPropertyDeclaration, SyntaxKindPropertySignature,
		SyntaxKindSetAccessor, SyntaxKindShorthandPropertyAssignment, SyntaxKindSourceFile, SyntaxKindSpreadAssignment, SyntaxKindStringLiteral,
		SyntaxKindTypeAliasDeclaration, SyntaxKindTypeLiteral, SyntaxKindTypeParameter, SyntaxKindVariableDeclaration:
		return true
	}
	return false
}

func canHaveLocals(node *Node) bool {
	switch node.Kind {
	case SyntaxKindArrowFunction, SyntaxKindBlock, SyntaxKindCallSignature, SyntaxKindCaseBlock, SyntaxKindCatchClause,
		SyntaxKindClassStaticBlockDeclaration, SyntaxKindConditionalType, SyntaxKindConstructor, SyntaxKindConstructorType,
		SyntaxKindConstructSignature, SyntaxKindForStatement, SyntaxKindForInStatement, SyntaxKindForOfStatement, SyntaxKindFunctionDeclaration,
		SyntaxKindFunctionExpression, SyntaxKindFunctionType, SyntaxKindGetAccessor, SyntaxKindIndexSignature, SyntaxKindJSDocCallbackTag,
		SyntaxKindJSDocEnumTag, SyntaxKindJSDocFunctionType, SyntaxKindJSDocSignature, SyntaxKindJSDocTypedefTag, SyntaxKindMappedType,
		SyntaxKindMethodDeclaration, SyntaxKindMethodSignature, SyntaxKindModuleDeclaration, SyntaxKindSetAccessor, SyntaxKindSourceFile,
		SyntaxKindTypeAliasDeclaration:
		return true
	}
	return false
}

func isAnyImportOrReExport(node *Node) bool {
	return isAnyImportSyntax(node) || IsExportDeclaration(node)
}

func isAnyImportSyntax(node *Node) bool {
	return nodeKindIs(node, SyntaxKindImportDeclaration, SyntaxKindImportEqualsDeclaration)
}

func getExternalModuleName(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindImportDeclaration:
		return node.AsImportDeclaration().ModuleSpecifier
	case SyntaxKindExportDeclaration:
		return node.AsExportDeclaration().ModuleSpecifier
	case SyntaxKindImportEqualsDeclaration:
		if node.AsImportEqualsDeclaration().ModuleReference.Kind == SyntaxKindExternalModuleReference {
			return node.AsImportEqualsDeclaration().ModuleReference.AsExternalModuleReference().Expression
		}
		return nil
	case SyntaxKindImportType:
		return getImportTypeNodeLiteral(node)
	case SyntaxKindCallExpression:
		return node.AsCallExpression().Arguments[0]
	case SyntaxKindModuleDeclaration:
		if IsStringLiteral(node.AsModuleDeclaration().Name) {
			return node.AsModuleDeclaration().Name
		}
		return nil
	}
	panic("Unhandled case in getExternalModuleName")
}

func getImportTypeNodeLiteral(node *Node) *Node {
	if isImportTypeNode(node) {
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
	return utils.MakeRegexp(`^\.\.?(?:$|[\\/])`).MatchString(path)
}

func extensionIsTs(ext string) bool {
	return ext == ExtensionTs || ext == ExtensionTsx || ext == ExtensionDts || ext == ExtensionMts || ext == ExtensionDmts || ext == ExtensionCts || ext == ExtensionDcts || len(ext) >= 7 && ext[:3] == ".d." && ext[len(ext)-3:] == ".ts"
}

func isShorthandAmbientModuleSymbol(moduleSymbol *Symbol) bool {
	return isShorthandAmbientModule(moduleSymbol.valueDeclaration)
}

func isShorthandAmbientModule(node *Node) bool {
	// The only kind of module that can be missing a body is a shorthand ambient module.
	return node != nil && node.Kind == SyntaxKindModuleDeclaration && node.AsModuleDeclaration().Body == nil
}

func isEntityName(node *Node) bool {
	return node.Kind == SyntaxKindIdentifier || node.Kind == SyntaxKindQualifiedName
}

func nodeIsSynthesized(node *Node) bool {
	return node.Loc.pos < 0 || node.Loc.end < 0
}

func getFirstIdentifier(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindIdentifier:
		return node
	case SyntaxKindQualifiedName:
		return getFirstIdentifier(node.AsQualifiedName().Left)
	case SyntaxKindPropertyAccessExpression:
		return getFirstIdentifier(node.AsPropertyAccessExpression().Expression)
	}
	panic("Unhandled case in getFirstIdentifier")
}

func getAliasDeclarationFromName(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindImportClause, SyntaxKindImportSpecifier, SyntaxKindNamespaceImport, SyntaxKindExportSpecifier, SyntaxKindExportAssignment,
		SyntaxKindImportEqualsDeclaration, SyntaxKindNamespaceExport:
		return node.Parent
	case SyntaxKindQualifiedName:
		return getAliasDeclarationFromName(node.Parent)
	}
	return nil
}

func entityNameToString(name *Node) string {
	switch name.Kind {
	case SyntaxKindThisKeyword:
		return "this"
	case SyntaxKindIdentifier, SyntaxKindPrivateIdentifier:
		return getTextOfNode(name)
	case SyntaxKindQualifiedName:
		return entityNameToString(name.AsQualifiedName().Left) + "." + entityNameToString(name.AsQualifiedName().Right)
	case SyntaxKindPropertyAccessExpression:
		return entityNameToString(name.AsPropertyAccessExpression().Expression) + "." + entityNameToString(name.AsPropertyAccessExpression().Name)
	case SyntaxKindJsxNamespacedName:
		return entityNameToString(name.AsJsxNamespacedName().Namespace) + ":" + entityNameToString(name.AsJsxNamespacedName().Name)
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
	return node.AsImportEqualsDeclaration().ModuleReference.AsExternalModuleReference().Expression
}

func isRightSideOfQualifiedNameOrPropertyAccess(node *Node) bool {
	parent := node.Parent
	switch parent.Kind {
	case SyntaxKindQualifiedName:
		return parent.AsQualifiedName().Right == node
	case SyntaxKindPropertyAccessExpression:
		return parent.AsPropertyAccessExpression().Name == node
	case SyntaxKindMetaProperty:
		return parent.AsMetaProperty().Name == node
	}
	return false
}

func getNamespaceDeclarationNode(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindImportDeclaration:
		importClause := node.AsImportDeclaration().ImportClause
		if importClause != nil && IsNamespaceImport(importClause.AsImportClause().NamedBindings) {
			return importClause.AsImportClause().NamedBindings
		}
	case SyntaxKindImportEqualsDeclaration:
		return node
	case SyntaxKindExportDeclaration:
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
	return IsCallExpression(node) && node.AsCallExpression().Expression.Kind == SyntaxKindImportKeyword
}

func getSourceFileOfModule(module *Symbol) *SourceFile {
	declaration := module.valueDeclaration
	if declaration == nil {
		declaration = getNonAugmentationDeclaration(module)
	}
	return getSourceFileOfNode(declaration)
}

func getNonAugmentationDeclaration(symbol *Symbol) *Node {
	return utils.Find(symbol.declarations, func(d *Node) bool {
		return !isExternalModuleAugmentation(d) && !(IsModuleDeclaration(d) && isGlobalScopeAugmentation(d))
	})
}

func isExternalModuleAugmentation(node *Node) bool {
	return isAmbientModule(node) && isModuleAugmentationExternal(node)
}

func isJsonSourceFile(file *SourceFile) bool {
	return file.ScriptKind == ScriptKindJSON
}

func isSyntacticDefault(node *Node) bool {
	return (IsExportAssignment(node) && !node.AsExportAssignment().IsExportEquals) ||
		hasSyntacticModifier(node, ModifierFlagsDefault) ||
		IsExportSpecifier(node) ||
		IsNamespaceExport(node)
}

func hasExportAssignmentSymbol(moduleSymbol *Symbol) bool {
	return moduleSymbol.exports[InternalSymbolNameExportEquals] != nil
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
	// 	Debug.assert(node.parent.kind == SyntaxKindJSDoc)
	// 	return flatMap(node.parent.tags, func(tag JSDocTag) *NodeArray[TypeParameterDeclaration] {
	// 		if isJSDocTemplateTag(tag) {
	// 			return tag.typeParameters
	// 		} else {
	// 			return nil
	// 		}
	// 	})
	// }
	typeParameters := node.GetTypeParameters()
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
	typeParameterList := node.GetTypeParameters()
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
	case SyntaxKindCallExpression:
		return node.AsCallExpression().TypeArguments
	case SyntaxKindNewExpression:
		return node.AsNewExpression().TypeArguments
	case SyntaxKindTaggedTemplateExpression:
		return node.AsTaggedTemplateExpression().TypeArguments
	case SyntaxKindTypeReference:
		return node.AsTypeReference().TypeArguments
	case SyntaxKindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().TypeArguments
	case SyntaxKindImportType:
		return node.AsImportTypeNode().TypeArguments
	case SyntaxKindTypeQuery:
		return node.AsTypeQueryNode().TypeArguments
	}
	panic("Unhandled case in getTypeArgumentListFromNode")
}

func getInitializerFromNode(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindVariableDeclaration:
		return node.AsVariableDeclaration().Initializer
	case SyntaxKindParameter:
		return node.AsParameterDeclaration().Initializer
	case SyntaxKindBindingElement:
		return node.AsBindingElement().Initializer
	case SyntaxKindPropertyDeclaration:
		return node.AsPropertyDeclaration().Initializer
	case SyntaxKindPropertyAssignment:
		return node.AsPropertyAssignment().Initializer
	case SyntaxKindEnumMember:
		return node.AsEnumMember().Initializer
	case SyntaxKindForStatement:
		return node.AsForStatement().Initializer
	case SyntaxKindForInStatement, SyntaxKindForOfStatement:
		return node.AsForInOrOfStatement().Initializer
	case SyntaxKindJsxAttribute:
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
	case SyntaxKindVariableDeclaration:
		return node.AsVariableDeclaration().TypeNode
	case SyntaxKindParameter:
		return node.AsParameterDeclaration().TypeNode
	case SyntaxKindPropertySignature:
		return node.AsPropertySignatureDeclaration().TypeNode
	case SyntaxKindPropertyDeclaration:
		return node.AsPropertyDeclaration().TypeNode
	case SyntaxKindTypePredicate:
		return node.AsTypePredicateNode().TypeNode
	case SyntaxKindParenthesizedType:
		return node.AsParenthesizedTypeNode().TypeNode
	case SyntaxKindTypeOperator:
		return node.AsTypeOperatorNode().TypeNode
	case SyntaxKindMappedType:
		return node.AsMappedTypeNode().TypeNode
	case SyntaxKindTypeAssertionExpression:
		return node.AsTypeAssertion().TypeNode
	case SyntaxKindAsExpression:
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
	return node != nil && node.Kind == SyntaxKindQuestionToken
}

func isOptionalDeclaration(declaration *Node) bool {
	switch declaration.Kind {
	case SyntaxKindParameter:
		return declaration.AsParameterDeclaration().QuestionToken != nil
	case SyntaxKindPropertyDeclaration:
		return isQuestionToken(declaration.AsPropertyDeclaration().PostfixToken)
	case SyntaxKindPropertySignature:
		return isQuestionToken(declaration.AsPropertySignatureDeclaration().PostfixToken)
	case SyntaxKindMethodDeclaration:
		return isQuestionToken(declaration.AsMethodDeclaration().PostfixToken)
	case SyntaxKindMethodSignature:
		return isQuestionToken(declaration.AsMethodSignatureDeclaration().PostfixToken)
	case SyntaxKindPropertyAssignment:
		return isQuestionToken(declaration.AsPropertyAssignment().PostfixToken)
	case SyntaxKindShorthandPropertyAssignment:
		return isQuestionToken(declaration.AsShorthandPropertyAssignment().PostfixToken)
	}
	return false
}

func isEmptyArrayLiteral(expression *Node) bool {
	return expression.Kind == SyntaxKindArrayLiteralExpression && len(expression.AsArrayLiteralExpression().Elements) == 0
}

func declarationBelongsToPrivateAmbientMember(declaration *Node) bool {
	root := getRootDeclaration(declaration)
	memberDeclaration := root
	if root.Kind == SyntaxKindParameter {
		memberDeclaration = root.Parent
	}
	return isPrivateWithinAmbient(memberDeclaration)
}

func isPrivateWithinAmbient(node *Node) bool {
	return (hasEffectiveModifier(node, ModifierFlagsPrivate) || isPrivateIdentifierClassElementDeclaration(node)) && node.Flags&NodeFlagsAmbient != 0
}

func identifierToKeywordKind(node *Identifier) SyntaxKind {
	return textToKeyword[node.Text]
}

func isAssertionExpression(node *Node) bool {
	kind := node.Kind
	return kind == SyntaxKindTypeAssertionExpression || kind == SyntaxKindAsExpression
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
	return utils.Find(symbol.declarations, isClassLike)
}

func isThisInTypeQuery(node *Node) bool {
	if !isThisIdentifier(node) {
		return false
	}
	for IsQualifiedName(node.Parent) && node.Parent.AsQualifiedName().Left == node {
		node = node.Parent
	}
	return node.Parent.Kind == SyntaxKindTypeQuery
}

func isThisIdentifier(node *Node) bool {
	return node != nil && node.Kind == SyntaxKindIdentifier && identifierIsThisKeyword(node)
}

func identifierIsThisKeyword(id *Node) bool {
	return id.AsIdentifier().Text == "this"
}

func getDeclarationModifierFlagsFromSymbol(s *Symbol) ModifierFlags {
	return getDeclarationModifierFlagsFromSymbolEx(s, false /*isWrite*/)
}

func getDeclarationModifierFlagsFromSymbolEx(s *Symbol, isWrite bool) ModifierFlags {
	if s.valueDeclaration != nil {
		var declaration *Node
		if isWrite {
			declaration = utils.Find(s.declarations, IsSetAccessorDeclaration)
		}
		if declaration == nil && s.flags&SymbolFlagsGetAccessor != 0 {
			declaration = utils.Find(s.declarations, IsGetAccessorDeclaration)
		}
		if declaration == nil {
			declaration = s.valueDeclaration
		}
		flags := getCombinedModifierFlags(declaration)
		if s.parent != nil && s.parent.flags&SymbolFlagsClass != 0 {
			return flags
		}
		return flags & ^ModifierFlagsAccessibilityModifier
	}
	if s.checkFlags&CheckFlagsSynthetic != 0 {
		var accessModifier ModifierFlags
		switch {
		case s.checkFlags&CheckFlagsContainsPrivate != 0:
			accessModifier = ModifierFlagsPrivate
		case s.checkFlags&CheckFlagsContainsPublic != 0:
			accessModifier = ModifierFlagsPublic
		default:
			accessModifier = ModifierFlagsProtected
		}
		var staticModifier ModifierFlags
		if s.checkFlags&CheckFlagsContainsStatic != 0 {
			staticModifier = ModifierFlagsStatic
		}
		return accessModifier | staticModifier
	}
	if s.flags&SymbolFlagsPrototype != 0 {
		return ModifierFlagsPublic | ModifierFlagsStatic
	}
	return ModifierFlagsNone
}

func isExponentiationOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindAsteriskAsteriskToken
}

func isMultiplicativeOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindAsteriskToken || kind == SyntaxKindSlashToken || kind == SyntaxKindPercentToken
}

func isMultiplicativeOperatorOrHigher(kind SyntaxKind) bool {
	return isExponentiationOperator(kind) || isMultiplicativeOperator(kind)
}

func isAdditiveOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindPlusToken || kind == SyntaxKindMinusToken
}

func isAdditiveOperatorOrHigher(kind SyntaxKind) bool {
	return isAdditiveOperator(kind) || isMultiplicativeOperatorOrHigher(kind)
}

func isShiftOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindLessThanLessThanToken || kind == SyntaxKindGreaterThanGreaterThanToken ||
		kind == SyntaxKindGreaterThanGreaterThanGreaterThanToken
}

func isShiftOperatorOrHigher(kind SyntaxKind) bool {
	return isShiftOperator(kind) || isAdditiveOperatorOrHigher(kind)
}

func isRelationalOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindLessThanToken || kind == SyntaxKindLessThanEqualsToken || kind == SyntaxKindGreaterThanToken ||
		kind == SyntaxKindGreaterThanEqualsToken || kind == SyntaxKindInstanceOfKeyword || kind == SyntaxKindInKeyword
}

func isRelationalOperatorOrHigher(kind SyntaxKind) bool {
	return isRelationalOperator(kind) || isShiftOperatorOrHigher(kind)
}

func isEqualityOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindEqualsEqualsToken || kind == SyntaxKindEqualsEqualsEqualsToken ||
		kind == SyntaxKindExclamationEqualsToken || kind == SyntaxKindExclamationEqualsEqualsToken
}

func isEqualityOperatorOrHigher(kind SyntaxKind) bool {
	return isEqualityOperator(kind) || isRelationalOperatorOrHigher(kind)
}

func isBitwiseOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindAmpersandToken || kind == SyntaxKindBarToken || kind == SyntaxKindCaretToken
}

func isBitwiseOperatorOrHigher(kind SyntaxKind) bool {
	return isBitwiseOperator(kind) || isEqualityOperatorOrHigher(kind)
}

// NOTE: The version in utilities includes ExclamationToken, which is not a binary operator.
func isLogicalOperator(kind SyntaxKind) bool {
	return kind == SyntaxKindAmpersandAmpersandToken || kind == SyntaxKindBarBarToken
}

func isLogicalOperatorOrHigher(kind SyntaxKind) bool {
	return isLogicalOperator(kind) || isBitwiseOperatorOrHigher(kind)
}

func isAssignmentOperatorOrHigher(kind SyntaxKind) bool {
	return kind == SyntaxKindQuestionQuestionToken || isLogicalOperatorOrHigher(kind) || isAssignmentOperator(kind)
}

func isBinaryOperator(kind SyntaxKind) bool {
	return isAssignmentOperatorOrHigher(kind) || kind == SyntaxKindCommaToken
}

func isObjectLiteralType(t *Type) bool {
	return t.objectFlags&ObjectFlagsObjectLiteral != 0
}

func isDeclarationReadonly(declaration *Node) bool {
	return getCombinedModifierFlags(declaration)&ModifierFlagsReadonly != 0 && !isParameterPropertyDeclaration(declaration, declaration.Parent)
}

func getPostfixTokenFromNode(node *Node) *Node {
	switch node.Kind {
	case SyntaxKindPropertyDeclaration:
		return node.AsPropertyDeclaration().PostfixToken
	case SyntaxKindPropertySignature:
		return node.AsPropertySignatureDeclaration().PostfixToken
	case SyntaxKindMethodDeclaration:
		return node.AsMethodDeclaration().PostfixToken
	case SyntaxKindMethodSignature:
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
		if node.Kind == SyntaxKindParenthesizedExpression {
			node = node.AsParenthesizedExpression().Expression
		} else if node.Kind == SyntaxKindPrefixUnaryExpression && node.AsPrefixUnaryExpression().Operator == SyntaxKindExclamationToken {
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
	return findAncestor(node.Parent, isFunctionLike)
}

func isTypeReferenceType(node *Node) bool {
	return node.Kind == SyntaxKindTypeReference || node.Kind == SyntaxKindExpressionWithTypeArguments
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
	case SyntaxKindIdentifier, SyntaxKindPrivateIdentifier, SyntaxKindStringLiteral, SyntaxKindNumericLiteral, SyntaxKindComputedPropertyName:
		return true
	}
	return false
}

func getPropertyNameForPropertyNameNode(name *Node) string {
	switch name.Kind {
	case SyntaxKindIdentifier, SyntaxKindPrivateIdentifier, SyntaxKindStringLiteral, SyntaxKindNoSubstitutionTemplateLiteral,
		SyntaxKindNumericLiteral, SyntaxKindBigIntLiteral, SyntaxKindJsxNamespacedName:
		return name.Text()
	case SyntaxKindComputedPropertyName:
		nameExpression := name.AsComputedPropertyName().Expression
		if isStringOrNumericLiteralLike(nameExpression) {
			return nameExpression.Text()
		}
		if isSignedNumericLiteral(nameExpression) {
			text := nameExpression.AsPrefixUnaryExpression().Operand.Text()
			if nameExpression.AsPrefixUnaryExpression().Operator == SyntaxKindMinusToken {
				text = "-" + text
			}
			return text
		}
		return InternalSymbolNameMissing
	}
	panic("Unhandled case in getPropertyNameForPropertyNameNode")
}

func isThisProperty(node *Node) bool {
	return (IsPropertyAccessExpression(node) || IsElementAccessExpression(node)) && node.Expression().Kind == SyntaxKindThisKeyword
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
		return isVarConst(node) && IsIdentifier(node.AsVariableDeclaration().Name) && isVariableDeclarationInVariableStatement(node)
	}
	if IsPropertyDeclaration(node) {
		return hasEffectiveReadonlyModifier(node) && hasStaticModifier(node)
	}
	return IsPropertySignatureDeclaration(node) && hasEffectiveReadonlyModifier(node)
}

func isVarConst(node *Node) bool {
	return getCombinedNodeFlags(node)&NodeFlagsBlockScoped == NodeFlagsConst
}

func isVariableDeclarationInVariableStatement(node *Node) bool {
	return IsVariableDeclarationList(node.Parent) && IsVariableStatement(node.Parent.Parent)
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
	return isThisIdentifier(parameter.GetName())
}

func getInterfaceBaseTypeNodes(node *Node) []*Node {
	heritageClause := getHeritageClause(node.AsInterfaceDeclaration().HeritageClauses, SyntaxKindExtendsKeyword)
	if heritageClause != nil {
		return heritageClause.AsHeritageClause().Types
	}
	return nil
}

func getHeritageClause(clauses []*Node, kind SyntaxKind) *Node {
	for _, clause := range clauses {
		if clause.AsHeritageClause().Token == kind {
			return clause
		}
	}
	return nil
}

func getClassExtendsHeritageElement(node *Node) *Node {
	heritageClause := getHeritageClause(node.ClassLikeData().HeritageClauses, SyntaxKindExtendsKeyword)
	if heritageClause != nil && len(heritageClause.AsHeritageClause().Types) > 0 {
		return heritageClause.AsHeritageClause().Types[0]
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
	case SyntaxKindJsxOpeningElement, SyntaxKindJsxSelfClosingElement, SyntaxKindCallExpression, SyntaxKindNewExpression,
		SyntaxKindTaggedTemplateExpression, SyntaxKindDecorator:
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
		node.Parent.Parent.AsBinaryExpression().OperatorToken.Kind == SyntaxKindEqualsToken &&
		node.Parent.Parent.AsBinaryExpression().Right.Kind == SyntaxKindThisKeyword
}

func isThisInitializedDeclaration(node *Node) bool {
	return node != nil && IsVariableDeclaration(node) && node.AsVariableDeclaration().Initializer != nil && node.AsVariableDeclaration().Initializer.Kind == SyntaxKindThisKeyword
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
	case SyntaxKindParenthesizedExpression:
		return accessKind(parent)
	case SyntaxKindPrefixUnaryExpression:
		operator := parent.AsPrefixUnaryExpression().Operator
		if operator == SyntaxKindPlusPlusToken || operator == SyntaxKindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case SyntaxKindPostfixUnaryExpression:
		operator := parent.AsPostfixUnaryExpression().Operator
		if operator == SyntaxKindPlusPlusToken || operator == SyntaxKindMinusMinusToken {
			return AccessKindReadWrite
		}
		return AccessKindRead
	case SyntaxKindBinaryExpression:
		if parent.AsBinaryExpression().Left == node {
			operator := parent.AsBinaryExpression().OperatorToken
			if isAssignmentOperator(operator.Kind) {
				if operator.Kind == SyntaxKindEqualsToken {
					return AccessKindWrite
				}
				return AccessKindReadWrite
			}
		}
		return AccessKindRead
	case SyntaxKindPropertyAccessExpression:
		if parent.AsPropertyAccessExpression().Name != node {
			return AccessKindRead
		}
		return accessKind(parent)
	case SyntaxKindPropertyAssignment:
		parentAccess := accessKind(parent.Parent)
		// In `({ x: varname }) = { x: 1 }`, the left `x` is a read, the right `x` is a write.
		if node == parent.AsPropertyAssignment().Name {
			return reverseAccessKind(parentAccess)
		}
		return parentAccess
	case SyntaxKindShorthandPropertyAssignment:
		// Assume it's the local variable being accessed, since we don't check public properties for --noUnusedLocals.
		if node == parent.AsShorthandPropertyAssignment().ObjectAssignmentInitializer {
			return AccessKindRead
		}
		return accessKind(parent.Parent)
	case SyntaxKindArrayLiteralExpression:
		return accessKind(parent)
	case SyntaxKindForInStatement, SyntaxKindForOfStatement:
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
	case SyntaxKindPropertyAssignment, SyntaxKindShorthandPropertyAssignment, SyntaxKindSpreadAssignment,
		SyntaxKindMethodDeclaration, SyntaxKindGetAccessor, SyntaxKindSetAccessor:
		return true
	}
	return false
}
