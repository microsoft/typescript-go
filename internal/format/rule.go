package format

import "github.com/microsoft/typescript-go/internal/ast"

type Rule struct {
	debugName string
	context   []ContextPredicate
	action    RuleAction
	flags     RuleFlags
}

func (r Rule) Action() RuleAction {
	return r.action
}

func (r Rule) Context() []ContextPredicate {
	return r.context
}

func (r Rule) Flags() RuleFlags {
	return r.flags
}

func (r Rule) String() string {
	return r.debugName
}

type tokenRange struct {
	tokens     []ast.Kind
	isSpecific bool
}

type ruleSpec struct {
	leftTokenRange  tokenRange
	rightTokenRange tokenRange
	rule            *Rule
}

/**
 * A rule takes a two tokens (left/right) and a particular context
 * for which you're meant to look at them. You then declare what should the
 * whitespace annotation be between these tokens via the action param.
 *
 * @param debugName Name to print
 * @param left The left side of the comparison
 * @param right The right side of the comparison
 * @param context A set of filters to narrow down the space in which this formatter rule applies
 * @param action a declaration of the expected whitespace
 * @param flags whether the rule deletes a line or not, defaults to no-op
 */
func rule(debugName string, left any, right any, context []ContextPredicate, action RuleAction, flags ...RuleFlags) ruleSpec {
	flag := RuleFlagsNone
	if len(flags) > 0 {
		flag = flags[0]
	}
	leftRange := toTokenRange(left)
	rightRange := toTokenRange(right)
	rule := &Rule{
		debugName: debugName,
		context:   context,
		action:    action,
		flags:     flag,
	}
	return ruleSpec{
		leftTokenRange:  leftRange,
		rightTokenRange: rightRange,
		rule:            rule,
	}
}

func toTokenRange(e any) tokenRange {
	switch t := e.(type) {
	case ast.Kind:
		return tokenRange{isSpecific: true, tokens: []ast.Kind{t}}
	case []ast.Kind:
		return tokenRange{isSpecific: true, tokens: t}
	case tokenRange:
		return t
	}
	panic("Unknown argument type passed to toTokenRange - only ast.Kind, []ast.Kind, and tokenRange supported")
}

type ContextPredicate = func(ctx *formattingContext) bool

var anyContext = []ContextPredicate{}

type RuleAction int

const (
	RuleActionNone                       RuleAction = 0
	RuleActionStopProcessingSpaceActions RuleAction = 1 << 0
	RuleActionStopProcessingTokenActions RuleAction = 1 << 1
	RuleActionInsertSpace                RuleAction = 1 << 2
	RuleActionInsertNewLine              RuleAction = 1 << 3
	RuleActionDeleteSpace                RuleAction = 1 << 4
	RuleActionDeleteToken                RuleAction = 1 << 5
	RuleActionInsertTrailingSemicolon    RuleAction = 1 << 6

	RuleActionStopAction        RuleAction = RuleActionStopProcessingSpaceActions | RuleActionStopProcessingTokenActions
	RuleActionModifySpaceAction RuleAction = RuleActionInsertSpace | RuleActionInsertNewLine | RuleActionDeleteSpace
	RuleActionModifyTokenAction RuleAction = RuleActionDeleteToken | RuleActionInsertTrailingSemicolon
)

type RuleFlags int

const (
	RuleFlagsNone RuleFlags = iota
	RuleFlagsCanDeleteNewLines
)
