package format

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type IndentStyle int

const (
	IndentStyleNone IndentStyle = iota
	IndentStyleBlock
	IndentStyleSmart
)

type SemicolonPreference string

const (
	SemicolonPreferenceIgnore SemicolonPreference = "ignore"
	SemicolonPreferenceInsert SemicolonPreference = "insert"
	SemicolonPreferenceRemove SemicolonPreference = "remove"
)

type EditorSettings struct {
	BaseIndentSize         int
	IndentSize             int
	TabSize                int
	NewLineCharacter       string
	ConvertTabsToSpaces    bool
	IndentStyle            IndentStyle
	TrimTrailingWhitespace bool
}

type FormatCodeSettings struct {
	EditorSettings
	InsertSpaceAfterCommaDelimiter                              bool
	InsertSpaceAfterSemicolonInForStatements                    bool
	InsertSpaceBeforeAndAfterBinaryOperators                    bool
	InsertSpaceAfterConstructor                                 bool
	InsertSpaceAfterKeywordsInControlFlowStatements             bool
	InsertSpaceAfterFunctionKeywordForAnonymousFunctions        bool
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis  bool
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets     bool
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces       bool
	InsertSpaceAfterOpeningAndBeforeClosingEmptyBraces          bool
	InsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces bool
	InsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces  bool
	InsertSpaceAfterTypeAssertion                               bool
	InsertSpaceBeforeFunctionParenthesis                        bool
	PlaceOpenBraceOnNewLineForFunctions                         bool
	PlaceOpenBraceOnNewLineForControlBlocks                     bool
	InsertSpaceBeforeTypeAnnotation                             bool
	IndentMultiLineObjectLiteralBeginningOnBlankLine            bool
	Semicolons                                                  SemicolonPreference
	IndentSwitchCase                                            bool
}

type formattingContext struct {
	currentTokenSpan   *TextRangeWithKind
	nextTokenSpan      *TextRangeWithKind
	contextNode        *ast.Node
	currentTokenParent *ast.Node
	nextTokenParent    *ast.Node

	contextNodeAllOnSameLine    core.Tristate
	nextNodeAllOnSameLine       core.Tristate
	tokensAreOnSameLine         core.Tristate
	contextNodeBlockIsOnOneLine core.Tristate
	nextNodeBlockIsOnOneLine    core.Tristate

	SourceFile            *ast.SourceFile
	FormattingRequestKind FormatRequestKind
	Options               *FormatCodeSettings

	scanner *scanner.Scanner
}

func NewFormattingContext(file *ast.SourceFile, kind FormatRequestKind, options *FormatCodeSettings) *formattingContext {
	res := &formattingContext{
		SourceFile:            file,
		FormattingRequestKind: kind,
		Options:               options,
		scanner:               scanner.NewScanner(),
	}
	res.scanner.SetText(file.Text())
	res.scanner.SetSkipTrivia(true)
	return res
}

func (this *formattingContext) UpdateContext(cur *TextRangeWithKind, curParent *ast.Node, next *TextRangeWithKind, nextParent *ast.Node, commonParent *ast.Node) {
	if cur == nil {
		panic("nil current range in update context")
	}
	if curParent == nil {
		panic("nil current range node parent in update context")
	}
	if next == nil {
		panic("nil next range in update context")
	}
	if nextParent == nil {
		panic("nil next range node parent in update context")
	}
	if commonParent == nil {
		panic("nil common parent node in update context")
	}
	this.currentTokenSpan = cur
	this.currentTokenParent = curParent
	this.nextTokenSpan = next
	this.nextTokenParent = nextParent
	this.contextNode = commonParent

	// drop cached results
	this.contextNodeAllOnSameLine = core.TSUnknown
	this.nextNodeAllOnSameLine = core.TSUnknown
	this.tokensAreOnSameLine = core.TSUnknown
	this.contextNodeBlockIsOnOneLine = core.TSUnknown
	this.nextNodeBlockIsOnOneLine = core.TSUnknown
}

func (this *formattingContext) rangeIsOnOneLine(node core.TextRange) core.Tristate {
	startLine, _ := scanner.GetLineAndCharacterOfPosition(this.SourceFile, node.Pos())
	endLine, _ := scanner.GetLineAndCharacterOfPosition(this.SourceFile, node.End())
	if startLine == endLine {
		return core.TSTrue
	}
	return core.TSFalse
}

func (this *formattingContext) blockIsOnOneLine(node *ast.Node) core.Tristate {
	// !!! in strada, this relies on token child manifesting - we just use the scanner here,
	// so this will have a differing performance profile. Is this OK? Needs profiling to know.
	this.scanner.ResetPos(node.Pos())
	end := node.End()
	firstOpenBrace := -1
	lastCloseBrace := -1
	for this.scanner.TokenEnd() < end {
		if firstOpenBrace == -1 && this.scanner.Token() == ast.KindOpenBraceToken {
			firstOpenBrace = this.scanner.TokenFullStart()
		} else if this.scanner.Token() == ast.KindCloseBraceToken {
			lastCloseBrace = this.scanner.TokenFullStart()
		}
		this.scanner.Scan()
	}
	if firstOpenBrace != -1 && lastCloseBrace != -1 {
		return this.rangeIsOnOneLine(core.NewTextRange(firstOpenBrace, lastCloseBrace))
	}
	return core.TSFalse
}

func (this *formattingContext) ContextNodeAllOnSameLine() bool {
	if this.contextNodeAllOnSameLine == core.TSUnknown {
		this.contextNodeAllOnSameLine = this.rangeIsOnOneLine(this.contextNode.Loc)
	}
	return this.contextNodeAllOnSameLine == core.TSTrue
}

func (this *formattingContext) NextNodeAllOnSameLine() bool {
	if this.nextNodeAllOnSameLine == core.TSUnknown {
		this.nextNodeAllOnSameLine = this.rangeIsOnOneLine(this.nextTokenParent.Loc)
	}
	return this.nextNodeAllOnSameLine == core.TSTrue
}

func (this *formattingContext) TokensAreOnSameLine() bool {
	if this.tokensAreOnSameLine == core.TSUnknown {
		this.tokensAreOnSameLine = this.rangeIsOnOneLine(core.NewTextRange(this.currentTokenSpan.Loc.Pos(), this.nextTokenSpan.Loc.End()))
	}
	return this.tokensAreOnSameLine == core.TSTrue
}

func (this *formattingContext) ContextNodeBlockIsOnOneLine() bool {
	if this.contextNodeBlockIsOnOneLine == core.TSUnknown {
		this.contextNodeBlockIsOnOneLine = this.blockIsOnOneLine(this.contextNode)
	}
	return this.contextNodeBlockIsOnOneLine == core.TSTrue
}

func (this *formattingContext) NextNodeBlockIsOnOneLine() bool {
	if this.nextNodeBlockIsOnOneLine == core.TSUnknown {
		this.nextNodeBlockIsOnOneLine = this.blockIsOnOneLine(this.nextTokenParent)
	}
	return this.nextNodeBlockIsOnOneLine == core.TSTrue
}
