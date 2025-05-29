package format

import (
	"context"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type FormatRequestKind int

const (
	FormatDocument FormatRequestKind = iota
	FormatSelection
	FormatOnEnter
	FormatOnSemicolon
	FormatOnOpeningCurlyBrace
	FormatOnClosingCurlyBrace
)

/** find node that fully contains given text range */
func findEnclosingNode(r core.TextRange, sourceFile *ast.SourceFile) *ast.Node {
	var find func(*ast.Node) *ast.Node
	find = func(n *ast.Node) *ast.Node {
		var candidate *ast.Node
		n.ForEachChild(func(c *ast.Node) bool {
			if r.ContainedBy(withTokenStart(c, sourceFile)) {
				candidate = c
				return true
			}
			return false
		})
		if candidate != nil {
			result := find(candidate)
			if result != nil {
				return result
			}
		}

		return n
	}
	return find(sourceFile.AsNode())
}

/**
 * Start of the original range might fall inside the comment - scanner will not yield appropriate results
 * This function will look for token that is located before the start of target range
 * and return its end as start position for the scanner.
 */
func getScanStartPosition(enclosingNode *ast.Node, originalRange core.TextRange, sourceFile *ast.SourceFile) int {
	adjusted := withTokenStart(enclosingNode, sourceFile)
	start := adjusted.Pos()
	if start == originalRange.Pos() && enclosingNode.End() == originalRange.End() {
		return start
	}

	precedingToken := astnav.FindPrecedingToken(sourceFile, originalRange.Pos())
	if precedingToken == nil {
		// no preceding token found - start from the beginning of enclosing node
		return enclosingNode.Pos()
	}

	// preceding token ends after the start of original range (i.e when originalRange.pos falls in the middle of literal)
	// start from the beginning of enclosingNode to handle the entire 'originalRange'
	if precedingToken.End() >= originalRange.Pos() {
		return enclosingNode.Pos()
	}

	return precedingToken.End()
}

/*
 * For cases like
 * if (a ||
 *     b ||$
 *     c) {...}
 * If we hit Enter at $ we want line '    b ||' to be indented.
 * Formatting will be applied to the last two lines.
 * Node that fully encloses these lines is binary expression 'a ||...'.
 * Initial indentation for this node will be 0.
 * Binary expressions don't introduce new indentation scopes, however it is possible
 * that some parent node on the same line does - like if statement in this case.
 * Note that we are considering parents only from the same line with initial node -
 * if parent is on the different line - its delta was already contributed
 * to the initial indentation.
 */
func getOwnOrInheritedDelta(n *ast.Node, options *FormatCodeSettings, sourceFile *ast.SourceFile) int {
	previousLine := -1
	var child *ast.Node
	for n != nil {
		line, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, withTokenStart(n, sourceFile).Pos())
		if previousLine != -1 && line != previousLine {
			break
		}

		if ShouldIndentChildNode(options, n, child, sourceFile) {
			return options.IndentSize // !!! nil check???
		}

		previousLine = line
		child = n
		n = n.Parent
	}
	return 0
}

func FormatSpan(ctx context.Context, span core.TextRange, file *ast.SourceFile, kind FormatRequestKind) []core.TextChange {
	// find the smallest node that fully wraps the range and compute the initial indentation for the node
	enclosingNode := findEnclosingNode(span, file)
	opts := ctx.Value(formatOptionsKey).(*FormatCodeSettings)

	return newFormattingScanner(
		file.Text(),
		file.LanguageVariant,
		getScanStartPosition(enclosingNode, span, file),
		span.End(),
		newFormatSpanWorker(
			ctx,
			span,
			enclosingNode,
			GetIndentationForNode(enclosingNode, &span, file, opts),
			getOwnOrInheritedDelta(enclosingNode, opts, file),
			kind,
			prepareRangeContainsErrorFunction(file.Diagnostics(), span),
			file,
		),
	)
}

func rangeHasNoErrors(_ core.TextRange) bool {
	return false
}

func prepareRangeContainsErrorFunction(errors []*ast.Diagnostic, originalRange core.TextRange) func(r core.TextRange) bool {
	if len(errors) == 0 {
		return rangeHasNoErrors
	}

	// pick only errors that fall in range
	sorted := core.Filter(errors, func(d *ast.Diagnostic) bool {
		return originalRange.Overlaps(d.Loc())
	})
	if len(sorted) == 0 {
		return rangeHasNoErrors
	}
	slices.SortStableFunc(sorted, func(a *ast.Diagnostic, b *ast.Diagnostic) int { return a.Pos() - b.Pos() })

	index := 0
	return func(r core.TextRange) bool {
		// in current implementation sequence of arguments [r1, r2...] is monotonically increasing.
		// 'index' tracks the index of the most recent error that was checked.
		for true {
			if index >= len(sorted) {
				// all errors in the range were already checked -> no error in specified range
				return false
			}

			err := sorted[index]

			if r.End() <= err.Pos() {
				// specified range ends before the error referred by 'index' - no error in range
				return false
			}

			if r.Overlaps(err.Loc()) {
				// specified range overlaps with error range
				return true
			}

			index++
		}
		return false // unreachable
	}
}

type formatSpanWorker struct {
	originalRange      core.TextRange
	enclosingNode      *ast.Node
	initialIndentation int
	delta              int
	requestKind        FormatRequestKind
	rangeContainsError func(r core.TextRange) bool
	sourceFile         *ast.SourceFile

	ctx context.Context

	formattingScanner *formattingScanner
	formattingContext *formattingContext

	edits                  []core.TextChange
	previousRange          *TextRangeWithKind
	previousRangeTriviaEnd int
	previousParent         *ast.Node
	previousRangeStartLine int
}

func newFormatSpanWorker(
	ctx context.Context,
	originalRange core.TextRange,
	enclosingNode *ast.Node,
	initialIndentation int,
	delta int,
	requestKind FormatRequestKind,
	rangeContainsError func(r core.TextRange) bool,
	sourceFile *ast.SourceFile,
) *formatSpanWorker {
	return &formatSpanWorker{
		ctx:                ctx,
		originalRange:      originalRange,
		enclosingNode:      enclosingNode,
		initialIndentation: initialIndentation,
		delta:              delta,
		requestKind:        requestKind,
		rangeContainsError: rangeContainsError,
		sourceFile:         sourceFile,
	}
}

type formatContextKey int

const (
	formatOptionsKey formatContextKey = iota
	formatNewlineKey
)

func NewContext(ctx context.Context, options *FormatCodeSettings, newLine string) context.Context {
	ctx = context.WithValue(ctx, formatOptionsKey, options)
	ctx = context.WithValue(ctx, formatNewlineKey, newLine)
	// In strada, the rules map was both globally cached *and* cached into the context, for some reason. We skip that here and just use the global one.
	return ctx
}

func getNewLineOrDefaultFromContext(ctx context.Context) string { // TODO: Move into broader LS - more than just the formatter uses the newline editor setting/host new line
	opt := ctx.Value(formatOptionsKey).(*FormatCodeSettings)
	if opt != nil && len(opt.NewLineCharacter) > 0 {
		return opt.NewLineCharacter
	}
	host := ctx.Value(formatNewlineKey).(string)
	if len(host) > 0 {
		return host
	}
	return "\n"
}

func getNonDecoratorTokenPosOfNode(node *ast.Node, file *ast.SourceFile) int {
	var lastDecorator *ast.Node
	if ast.HasDecorators(node) {
		lastDecorator = core.FindLast(node.Modifiers().Nodes, ast.IsDecorator)
	}
	if file == nil {
		file = ast.GetSourceFileOfNode(node)
	}
	if lastDecorator == nil {
		return withTokenStart(node, file).Pos()
	}
	return scanner.SkipTrivia(file.Text(), lastDecorator.End())
}

func (w *formatSpanWorker) execute(s *formattingScanner) []core.TextChange {
	w.formattingScanner = s
	opt := w.ctx.Value(formatOptionsKey).(*FormatCodeSettings)
	w.formattingContext = NewFormattingContext(w.sourceFile, w.requestKind, opt)
	// formatting context is used by rules provider

	w.formattingScanner.advance()

	if w.formattingScanner.isOnToken() {
		startLine, _ := scanner.GetLineAndCharacterOfPosition(w.sourceFile, withTokenStart(w.enclosingNode, w.sourceFile).Pos())
		undecoratedStartLine := startLine
		if ast.HasDecorators(w.enclosingNode) {
			undecoratedStartLine, _ = scanner.GetLineAndCharacterOfPosition(w.sourceFile, getNonDecoratorTokenPosOfNode(w.enclosingNode, w.sourceFile))
		}

		w.processNode(w.enclosingNode, w.enclosingNode, startLine, undecoratedStartLine, w.initialIndentation, w.delta)
	}

	// Leading trivia items get attached to and processed with the token that proceeds them. If the
	// range ends in the middle of some leading trivia, the token that proceeds them won't be in the
	// range and thus won't get processed. So we process those remaining trivia items here.
	remainingTrivia := w.formattingScanner.getCurrentLeadingTrivia()
	if len(remainingTrivia) > 0 {
		indentation := w.initialIndentation
		if NodeWillIndentChild(w.formattingContext.Options, w.enclosingNode, nil, w.sourceFile, false) {
			indentation += opt.IndentSize // !!! TODO: nil check???
		}

		w.indentTriviaItems(remainingTrivia, indentation, true, func(item *TextRangeWithKind) {
			startLine, startChar := scanner.GetLineAndCharacterOfPosition(w.sourceFile, item.Loc.Pos())
			w.processRange(item, startLine, startChar, w.enclosingNode, w.enclosingNode, nil)
			w.insertIndentation(item.Loc.Pos(), indentation, false)
		})

		if opt.TrimTrailingWhitespace != false {
			w.trimTrailingWhitespacesForRemainingRange(remainingTrivia)
		}
	}

	if w.previousRange != nil && w.formattingScanner.getTokenFullStart() >= w.originalRange.End() {
		// Formatting edits happen by looking at pairs of contiguous tokens (see `processPair`),
		// typically inserting or deleting whitespace between them. The recursive `processNode`
		// logic above bails out as soon as it encounters a token that is beyond the end of the
		// range we're supposed to format (or if we reach the end of the file). But this potentially
		// leaves out an edit that would occur *inside* the requested range but cannot be discovered
		// without looking at one token *beyond* the end of the range: consider the line `x = { }`
		// with a selection from the beginning of the line to the space inside the curly braces,
		// inclusive. We would expect a format-selection would delete the space (if rules apply),
		// but in order to do that, we need to process the pair ["{", "}"], but we stopped processing
		// just before getting there. This block handles this trailing edit.
		var tokenInfo *TextRangeWithKind
		if w.formattingScanner.isOnEOF() {
			tokenInfo = w.formattingScanner.readEOFTokenRange()
		} else if w.formattingScanner.isOnToken() {
			tokenInfo = w.formattingScanner.readTokenInfo(w.enclosingNode).token
		}

		if tokenInfo != nil && tokenInfo.Loc.Pos() == w.previousRangeTriviaEnd {
			// We need to check that tokenInfo and previousRange are contiguous: the `originalRange`
			// may have ended in the middle of a token, which means we will have stopped formatting
			// on that token, leaving `previousRange` pointing to the token before it, but already
			// having moved the formatting scanner (where we just got `tokenInfo`) to the next token.
			// If this happens, our supposed pair [previousRange, tokenInfo] actually straddles the
			// token that intersects the end of the range we're supposed to format, so the pair will
			// produce bogus edits if we try to `processPair`. Recall that the point of this logic is
			// to perform a trailing edit at the end of the selection range: but there can be no valid
			// edit in the middle of a token where the range ended, so if we have a non-contiguous
			// pair here, we're already done and we can ignore it.
			parent := astnav.FindPrecedingToken(w.sourceFile, tokenInfo.Loc.End())
			if parent == nil {
				parent = w.previousParent
			}
			line, _ := scanner.GetLineAndCharacterOfPosition(w.sourceFile, tokenInfo.Loc.Pos())
			w.processPair(
				tokenInfo,
				line,
				parent,
				w.previousRange,
				w.previousRangeStartLine,
				w.previousParent,
				parent,
				nil,
			)
		}
	}

	return w.edits
}

func (w *formatSpanWorker) getProcessNodeVisitor(node *ast.Node, indenter *dynamicIndenter, nodeStartLine int, undecoratedNodeStartLine int, childContextNode *ast.Node) *ast.NodeVisitor {
	processChildNode := func(
		child *ast.Node,
		inheritedIndentation int,
		parent *ast.Node,
		parentDynamicIndentation *dynamicIndenter,
		parentStartLine int,
		undecoratedParentStartLine int,
		isListItem bool,
		isFirstListItem bool,
	) {
		// !!!
	}

	processChildNodes := func(
		nodes *ast.NodeList,
		parent *ast.Node,
		parentStartLine int,
		parentDynamicIndentation *dynamicIndenter,
	) {
		// !!!
	}

	return ast.NewNodeVisitor(func(child *ast.Node) *ast.Node {
		processChildNode(child, -1, node, indenter, nodeStartLine, undecoratedNodeStartLine, false, false)
		return node
	}, &ast.NodeFactory{}, ast.NodeVisitorHooks{
		VisitNodes: func(nodes *ast.NodeList, v *ast.NodeVisitor) *ast.NodeList {
			processChildNodes(nodes, node, nodeStartLine, indenter)
			return nodes
		},
	})
}

func (w *formatSpanWorker) processNode(node *ast.Node, contextNode *ast.Node, nodeStartLine int, undecoratedNodeStartLine int, indentation int, delta int) {
	if !w.originalRange.Overlaps(withTokenStart(node, w.sourceFile)) {
		return
	}

	nodeDynamicIndentation := w.getDynamicIndentation(node, nodeStartLine, indentation, delta)

	// a useful observations when tracking context node
	//        /
	//      [a]
	//   /   |   \
	//  [b] [c] [d]
	// node 'a' is a context node for nodes 'b', 'c', 'd'
	// except for the leftmost leaf token in [b] - in this case context node ('e') is located somewhere above 'a'
	// this rule can be applied recursively to child nodes of 'a'.
	//
	// context node is set to parent node value after processing every child node
	// context node is set to parent of the token after processing every token

	childContextNode := contextNode

	// if there are any tokens that logically belong to node and interleave child nodes
	// such tokens will be consumed in processChildNode for the child that follows them
	v := w.getProcessNodeVisitor(node, nodeDynamicIndentation, nodeStartLine, undecoratedNodeStartLine, childContextNode)
	node.VisitEachChild(v)

	// proceed any tokens in the node that are located after child nodes
	for w.formattingScanner.isOnToken() && w.formattingScanner.getTokenFullStart() < w.originalRange.End() {
		tokenInfo := w.formattingScanner.readTokenInfo(node)
		if tokenInfo.token.Loc.End() > min(node.End(), w.originalRange.End()) {
			break
		}
		w.consumeTokenAndAdvanceScanner(tokenInfo, node, nodeDynamicIndentation, node, false)
	}
}

func (w *formatSpanWorker) processPair(currentItem *TextRangeWithKind, currentStartLine int, currentParent *ast.Node, previousItem *TextRangeWithKind, previousStartLine int, previousParent *ast.Node, contextNode *ast.Node, dynamicIndentation *dynamicIndenter) {
	// !!!
}

type LineAction int

const (
	LineActionNone LineAction = iota
	LineActionLineAdded
	LineActionLineRemoved
)

func (w *formatSpanWorker) processRange(r *TextRangeWithKind, rangeStartLine int, rangeStartCharacter int, parent *ast.Node, contextNode *ast.Node, dynamicIndentation *dynamicIndenter) LineAction {
	// !!!
	return LineActionNone
}

/**
* Trimming will be done for lines after the previous range.
* Exclude comments as they had been previously processed.
 */
func (w *formatSpanWorker) trimTrailingWhitespacesForRemainingRange(trivias []*TextRangeWithKind) {
	// !!!
}

func (w *formatSpanWorker) insertIndentation(pos int, indentation int, lineAdded bool) {
	// !!!
}

func (w *formatSpanWorker) indentTriviaItems(trivia []*TextRangeWithKind, commentIndentation int, indentNextTokenOrTrivia bool, indentSingleLine func(item *TextRangeWithKind)) {
	// !!!
}

func (w *formatSpanWorker) consumeTokenAndAdvanceScanner(currentTokenInfo *tokenInfo, parent *ast.Node, dynamicIndenation *dynamicIndenter, container *ast.Node, isListEndToken bool) {
	// assert(currentTokenInfo.token.Loc.ContainedBy(parent.Loc)) // !!!
	// !!!
}

type dynamicIndenter struct {
	node          *ast.Node
	nodeStartLine int
	indentation   int
	delta         int

	options    *FormatCodeSettings
	sourceFile *ast.SourceFile
}

func (i *dynamicIndenter) getIndentationForComment(kind ast.Kind, tokenIndentation int, container *ast.Node) int {
	switch kind {
	// preceding comment to the token that closes the indentation scope inherits the indentation from the scope
	// ..  {
	//     // comment
	// }
	case ast.KindCloseBraceToken, ast.KindCloseBracketToken, ast.KindCloseParenToken:
		return i.indentation + i.getDelta(container)
	}
	return i.indentation
}

// if list end token is LessThanToken '>' then its delta should be explicitly suppressed
// so that LessThanToken as a binary operator can still be indented.
// foo.then
//
//	<
//	    number,
//	    string,
//	>();
//
// vs
// var a = xValue
//
//	> yValue;
func (i *dynamicIndenter) getIndentationForToken(line int, kind ast.Kind, container *ast.Node, suppressDelta bool) int {
	if !suppressDelta && i.shouldAddDelta(line, kind, container) {
		return i.indentation + i.getDelta(container)
	}
	return i.indentation
}

func (i *dynamicIndenter) getIndentation() int {
	return i.indentation
}

func (i *dynamicIndenter) getDelta(child *ast.Node) int {
	// Delta value should be zero when the node explicitly prevents indentation of the child node
	if NodeWillIndentChild(i.options, i.node, child, i.sourceFile, true) {
		return i.delta
	}
	return 0
}

func (i *dynamicIndenter) recomputeIndentation(lineAdded bool, parent *ast.Node) {
	if ShouldIndentChildNode(i.options, parent, i.node, i.sourceFile) {
		if lineAdded {
			i.indentation += i.options.IndentSize // !!! no nil check???
		} else {
			i.indentation -= i.options.IndentSize // !!! no nil check???
		}
		if ShouldIndentChildNode(i.options, i.node, nil, nil) {
			i.delta = i.options.IndentSize
		} else {
			i.delta = 0
		}
	}
}

func (i *dynamicIndenter) shouldAddDelta(line int, kind ast.Kind, container *ast.Node) bool {
	switch kind {
	// open and close brace, 'else' and 'while' (in do statement) tokens has indentation of the parent
	case ast.KindOpenBraceToken, ast.KindCloseBraceToken, ast.KindCloseParenToken, ast.KindElseKeyword, ast.KindWhileKeyword, ast.KindAtToken:
		return false
	case ast.KindSlashToken, ast.KindGreaterThanToken:
		switch container.Kind {
		case ast.KindJsxOpeningElement, ast.KindJsxClosingElement, ast.KindJsxSelfClosingElement:
			return false
		}
		break
	case ast.KindOpenBracketToken, ast.KindCloseBracketToken:
		if container.Kind != ast.KindMappedType {
			return false
		}
		break
	}
	// if token line equals to the line of containing node (this is a first token in the node) - use node indentation
	return i.nodeStartLine != line &&
		// if this token is the first token following the list of decorators, we do not need to indent
		!(ast.HasDecorators(i.node) && kind == getFirstNonDecoratorTokenOfNode(i.node))
}

func getFirstNonDecoratorTokenOfNode(node *ast.Node) ast.Kind {
	if ast.CanHaveModifiers(node) {
		modifier := core.Find(node.Modifiers().Nodes[core.FindIndex(node.Modifiers().Nodes, ast.IsDecorator):], ast.IsModifier)
		if modifier != nil {
			return modifier.Kind
		}
	}

	switch node.Kind {
	case ast.KindClassDeclaration:
		return ast.KindClassKeyword
	case ast.KindInterfaceDeclaration:
		return ast.KindInterfaceKeyword
	case ast.KindFunctionDeclaration:
		return ast.KindFunctionKeyword
	case ast.KindEnumDeclaration:
		return ast.KindEnumDeclaration
	case ast.KindGetAccessor:
		return ast.KindGetKeyword
	case ast.KindSetAccessor:
		return ast.KindSetKeyword
	case ast.KindMethodDeclaration:
		if node.AsMethodDeclaration().AsteriskToken != nil {
			return ast.KindAsteriskToken
		}
		fallthrough

	case ast.KindPropertyDeclaration, ast.KindParameter:
		name := ast.GetNameOfDeclaration(node)
		if name != nil {
			return name.Kind
		}
	}

	return ast.KindUnknown
}

func (w *formatSpanWorker) getDynamicIndentation(node *ast.Node, nodeStartLine int, indentation int, delta int) *dynamicIndenter {
	return &dynamicIndenter{
		node:          node,
		nodeStartLine: nodeStartLine,
		indentation:   indentation,
		delta:         delta,
		options:       w.formattingContext.Options,
		sourceFile:    w.sourceFile,
	}
}
