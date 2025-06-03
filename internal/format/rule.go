package format

type Rule interface {
	String() string
	Context() []ContextPredicate
	Action() RuleAction
	Flags() RuleFlags
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
