package ls

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

type changeNodeOptions struct {
	// Text to be inserted before the new node
	prefix string

	// Text to be inserted after the new node
	suffix string

	// Text of inserted node will be formatted with this indentation, otherwise indentation will be inferred from the old node
	indentation *int

	// Text of inserted node will be formatted with this delta, otherwise delta will be inferred from the new node kind
	delta *int

	leadingTriviaOption
	trailingTriviaOption
	joiner string
}

type leadingTriviaOption int

const (
	leadingTriviaOptionNone       leadingTriviaOption = 0
	leadingTriviaOptionExclude    leadingTriviaOption = 1
	leadingTriviaOptionIncludeAll leadingTriviaOption = 2
	leadingTriviaOptionJSDoc      leadingTriviaOption = 3
	leadingTriviaOptionStartLine  leadingTriviaOption = 4
)

type trailingTriviaOption int

const (
	trailingTriviaOptionNone              trailingTriviaOption = 0
	trailingTriviaOptionExclude           trailingTriviaOption = 1
	trailingTriviaOptionExcludeWhitespace trailingTriviaOption = 2
	trailingTriviaOptionInclude           trailingTriviaOption = 3
)

type trackerEditKind int

const (
	trackerEditKindText                     trackerEditKind = 1
	trackerEditKindRemove                   trackerEditKind = 2
	trackerEditKindReplaceWithSingleNode    trackerEditKind = 3
	trackerEditKindReplaceWithMultipleNodes trackerEditKind = 4
)

type trackerEdit struct {
	kind trackerEditKind
	lsproto.Range

	NewText string // kind == text

	*ast.Node             // single
	nodes     []*ast.Node // multiple
	options   changeNodeOptions
}

type changeTracker struct {
	// initialized with
	formatSettings *format.FormatCodeSettings
	newLine        string
	ls             *LanguageService
	ctx            context.Context
	*printer.EmitContext

	*ast.NodeFactory
	changes *collections.MultiMap[*ast.SourceFile, *trackerEdit]

	// created during call to getChanges
	writer *printer.ChangeTrackerWriter
	// printer
}

func (ls *LanguageService) newChangeTracker(ctx context.Context) *changeTracker {
	emitContext := printer.NewEmitContext()
	newLine := ls.GetProgram().Options().NewLine.GetNewLineCharacter()
	formatCodeSettings := format.GetDefaultFormatCodeSettings(newLine) // !!! format.GetFormatCodeSettingsFromContext(ctx),
	ctx = format.WithFormatCodeSettings(ctx, formatCodeSettings, newLine)
	return &changeTracker{
		ls:             ls,
		EmitContext:    emitContext,
		NodeFactory:    &emitContext.Factory.NodeFactory,
		changes:        &collections.MultiMap[*ast.SourceFile, *trackerEdit]{},
		ctx:            ctx,
		formatSettings: formatCodeSettings,
		newLine:        newLine,
	}
}

// !!! address strada note
//   - Note: after calling this, the TextChanges object must be discarded!
func (ct *changeTracker) getChanges() map[string][]*lsproto.TextEdit {
	// !!! finishDeleteDeclarations
	// !!! finishClassesWithNodesInsertedAtStart
	changes := ct.getTextChangesFromChanges()
	// !!! changes for new files
	return changes
}

func (ct *changeTracker) getTextChangesFromChanges() map[string][]*lsproto.TextEdit {
	changes := map[string][]*lsproto.TextEdit{}
	for sourceFile, changesInFile := range ct.changes.M {
		// order changes by start position
		// If the start position is the same, put the shorter range first, since an empty range (x, x) may precede (x, y) but not vice-versa.
		slices.SortStableFunc(changesInFile, func(a, b *trackerEdit) int { return CompareRanges(ptrTo(a.Range), ptrTo(b.Range)) })
		// verify that change intervals do not overlap, except possibly at end points.
		for i := range len(changesInFile) - 1 {
			if ComparePositions(changesInFile[i].Range.End, changesInFile[i+1].Range.Start) > 0 {
				// assert change[i].End <= change[i + 1].Start
				panic(fmt.Sprintf("changes overlap: %v and %v", changesInFile[i].Range, changesInFile[i+1].Range))
			}
		}

		textChanges := core.MapNonNil(changesInFile, func(change *trackerEdit) *lsproto.TextEdit {
			// !!! targetSourceFile

			newText := ct.computeNewText(change, sourceFile, sourceFile)
			// span := createTextSpanFromRange(c.Range)
			// !!!
			// Filter out redundant changes.
			// if (span.length == newText.length && stringContainsAt(targetSourceFile.text, newText, span.start)) { return nil }

			return &lsproto.TextEdit{
				NewText: newText,
				Range:   change.Range,
			}
		})

		if len(textChanges) > 0 {
			changes[sourceFile.FileName()] = textChanges
		}
	}
	return changes
}

func (ct *changeTracker) computeNewText(change *trackerEdit, targetSourceFile *ast.SourceFile, sourceFile *ast.SourceFile) string {
	switch change.kind {
	case trackerEditKindRemove:
		return ""
	case trackerEditKindText:
		return change.NewText
	}

	pos := int(ct.ls.converters.LineAndCharacterToPosition(sourceFile, change.Range.Start))
	targetFileLineMap := targetSourceFile.LineMap()
	format := func(n *ast.Node) string {
		return ct.getFormattedTextOfNode(n, targetSourceFile, sourceFile, pos, targetFileLineMap, change.options)
	}

	var text string
	switch change.kind {

	case trackerEditKindReplaceWithMultipleNodes:
		if change.options.joiner == "" {
			change.options.joiner = ct.newLine
		}
		text = strings.Join(core.Map(change.nodes, func(n *ast.Node) string { return strings.TrimSuffix(format(n), ct.newLine) }), change.options.joiner)
	case trackerEditKindReplaceWithSingleNode:
		text = format(change.Node)
	default:
		panic(fmt.Sprintf("change kind %d should have been handled earlier", change.kind))
	}
	// strip initial indentation (spaces or tabs) if text will be inserted in the middle of the line
	noIndent := text
	if !(change.options.indentation != nil && *change.options.indentation != 0 || scanner.GetLineStartPositionForPosition(pos, targetFileLineMap) == pos) {
		noIndent = strings.TrimLeftFunc(text, unicode.IsSpace)
	}
	return change.options.prefix + noIndent // !!!  +((!options.suffix || endsWith(noIndent, options.suffix)) ? "" : options.suffix);
}

/** Note: this may mutate `nodeIn`. */
func (ct *changeTracker) getFormattedTextOfNode(nodeIn *ast.Node, targetSourceFile *ast.SourceFile, sourceFile *ast.SourceFile, pos int, targetFileLineMap []core.TextPos, options changeNodeOptions) string {
	text, node := ct.getNonformattedText(nodeIn, targetSourceFile)
	// !!! if (validate) validate(node, text);
	formatOptions := getFormatCodeSettingsForWriting(ct.formatSettings, targetSourceFile)

	var initialIndentation, delta int
	if options.indentation == nil {
		// !!! indentation for position
		// initialIndentation = format.GetIndentationForPos(pos, sourceFile, formatOptions, options.prefix == ct.newLine || scanner.GetLineStartPositionForPosition(pos, targetFileLineMap) == pos);
	} else {
		initialIndentation = *options.indentation
	}

	if options.delta != nil {
		delta = *options.delta
	} else if formatOptions.IndentSize != 0 && format.ShouldIndentChildNode(formatOptions, nodeIn, nil, nil) {
		delta = formatOptions.IndentSize
	}

	changes := format.FormatNodeGivenIndentation(ct.ctx, node, targetSourceFile, targetSourceFile.LanguageVariant, initialIndentation, delta)
	return applyTextChanges(text, changes)
}

func getFormatCodeSettingsForWriting(options *format.FormatCodeSettings, sourceFile *ast.SourceFile) *format.FormatCodeSettings {
	shouldAutoDetectSemicolonPreference := options.Semicolons == format.SemicolonPreferenceIgnore
	shouldRemoveSemicolons := options.Semicolons == format.SemicolonPreferenceRemove || shouldAutoDetectSemicolonPreference && !probablyUsesSemicolons(sourceFile)
	if shouldRemoveSemicolons {
		options.Semicolons = format.SemicolonPreferenceRemove
	}

	return options
}

/** Note: output node may be mutated input node. */
func (ct *changeTracker) getNonformattedText(node *ast.Node, sourceFile *ast.SourceFile) (string, *ast.Node) {
	writer := printer.NewChangeTrackerWriter(ct.newLine)

	printer.NewPrinter(
		printer.PrinterOptions{
			NewLine:                       core.GetNewLineKind(ct.newLine),
			NeverAsciiEscape:              true,
			PreserveSourceNewlines:        true,
			TerminateUnterminatedLiterals: true,
		},
		writer.GetPrintHandlers(),
		ct.EmitContext,
	).Write(node, sourceFile, writer, nil)

	text := writer.String()

	return text, writer.AssignPositionsToNode(node, ct.NodeFactory)
}

func (ct *changeTracker) Write(s string) { ct.writer.Write(s) }

func (ct *changeTracker) replaceNode(sourceFile *ast.SourceFile, oldNode *ast.Node, newNode *ast.Node, options *changeNodeOptions) {
	if options == nil {
		// defaults to `useNonAdjustedPositions`
		options = &changeNodeOptions{
			leadingTriviaOption:  leadingTriviaOptionExclude,
			trailingTriviaOption: trailingTriviaOptionExclude,
		}
	}
	ct.replaceRange(sourceFile, ct.getAdjustedRange(sourceFile, oldNode, oldNode, options.leadingTriviaOption, options.trailingTriviaOption), newNode, *options)
}

func (ct *changeTracker) replaceRange(sourceFile *ast.SourceFile, lsprotoRange lsproto.Range, newNode *ast.Node, options changeNodeOptions) {
	ct.changes.Add(sourceFile, &trackerEdit{kind: trackerEditKindReplaceWithSingleNode, Range: lsprotoRange, options: options, Node: newNode})
}

func (ct *changeTracker) replaceRangeWithText(sourceFile *ast.SourceFile, lsprotoRange lsproto.Range, text string) {
	ct.changes.Add(sourceFile, &trackerEdit{kind: trackerEditKindText, Range: lsprotoRange, NewText: text})
}

func (ct *changeTracker) replaceRangeWithNodes(sourceFile *ast.SourceFile, lsprotoRange lsproto.Range, newNodes []*ast.Node, options changeNodeOptions) {
	if len(newNodes) == 1 {
		ct.replaceRange(sourceFile, lsprotoRange, newNodes[0], options)
		return
	}
	ct.changes.Add(sourceFile, &trackerEdit{kind: trackerEditKindReplaceWithMultipleNodes, Range: lsprotoRange, nodes: newNodes, options: options})
}

func (ct *changeTracker) insertText(sourceFile *ast.SourceFile, pos lsproto.Position, text string) {
	ct.replaceRangeWithText(sourceFile, lsproto.Range{Start: pos, End: pos}, text)
}

func (ct *changeTracker) insertNodeAt(sourceFile *ast.SourceFile, pos core.TextPos, newNode *ast.Node, options changeNodeOptions) {
	lsPos := ct.ls.converters.PositionToLineAndCharacter(sourceFile, pos)
	ct.replaceRange(sourceFile, lsproto.Range{Start: lsPos, End: lsPos}, newNode, options)
}

func (ct *changeTracker) insertNodesAt(sourceFile *ast.SourceFile, pos core.TextPos, newNodes []*ast.Node, options changeNodeOptions) {
	lsPos := ct.ls.converters.PositionToLineAndCharacter(sourceFile, pos)
	ct.replaceRangeWithNodes(sourceFile, lsproto.Range{Start: lsPos, End: lsPos}, newNodes, options)
}

func (ct *changeTracker) insertNodeAfter(sourceFile *ast.SourceFile, after *ast.Node, newNode *ast.Node) {
	endPosition := ct.endPosForInsertNodeAfter(sourceFile, after, newNode)
	ct.insertNodeAt(sourceFile, endPosition, newNode, ct.getInsertNodeAfterOptions(sourceFile, after))
}

func (ct *changeTracker) insertNodesAfter(sourceFile *ast.SourceFile, after *ast.Node, newNodes []*ast.Node) {
	endPosition := ct.endPosForInsertNodeAfter(sourceFile, after, newNodes[0])
	ct.insertNodesAt(sourceFile, endPosition, newNodes, ct.getInsertNodeAfterOptions(sourceFile, after))
}

func (ct *changeTracker) endPosForInsertNodeAfter(sourceFile *ast.SourceFile, after *ast.Node, newNode *ast.Node) core.TextPos {
	if (needSemicolonBetween(after, newNode)) && (rune(sourceFile.Text()[after.End()-1]) != ';') {
		// check if previous statement ends with semicolon
		// if not - insert semicolon to preserve the code from changing the meaning due to ASI
		endPos := ct.ls.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(after.End()))
		ct.replaceRange(sourceFile,
			lsproto.Range{Start: endPos, End: endPos},
			sourceFile.GetOrCreateToken(ast.KindSemicolonToken, after.End(), after.End(), after.Parent),
			changeNodeOptions{},
		)
	}
	return core.TextPos(ct.getAdjustedEndPosition(sourceFile, after, trailingTriviaOptionNone))
}

func (ct *changeTracker) getInsertNodeAfterOptions(sourceFile *ast.SourceFile, node *ast.Node) changeNodeOptions {
	newLineChar := ct.newLine
	var options changeNodeOptions
	switch node.Kind {
	case ast.KindParameter:
		// default opts
		options = changeNodeOptions{}
	case ast.KindClassDeclaration, ast.KindModuleDeclaration:
		options = changeNodeOptions{prefix: newLineChar, suffix: newLineChar}

	case ast.KindVariableDeclaration, ast.KindStringLiteral, ast.KindIdentifier:
		options = changeNodeOptions{prefix: ", "}

	case ast.KindPropertyAssignment:
		options = changeNodeOptions{suffix: "," + newLineChar}

	case ast.KindExportKeyword:
		options = changeNodeOptions{prefix: " "}

	default:
		if !(ast.IsStatement(node) || ast.IsClassOrTypeElement(node)) {
			// Else we haven't handled this kind of node yet -- add it
			panic("unimplemented node type " + node.Kind.String() + " in changeTracker.getInsertNodeAfterOptions")
		}
		options = changeNodeOptions{suffix: newLineChar}
	}
	if node.End() == sourceFile.End() && ast.IsStatement(node) {
		options.prefix = "\n" + options.prefix
	}

	return options
}

/**
* This function should be used to insert nodes in lists when nodes don't carry separators as the part of the node range,
* i.e. arguments in arguments lists, parameters in parameter lists etc.
* Note that separators are part of the node in statements and class elements.
 */
func (ct *changeTracker) insertNodeInListAfter(sourceFile *ast.SourceFile, after *ast.Node, newNode *ast.Node, containingList []*ast.Node) {
	if len(containingList) == 0 {
		containingList = format.GetContainingList(after, sourceFile).Nodes
	}
	index := slices.Index(containingList, after)
	if index < 0 {
		return
	}
	if index != len(containingList)-1 {
		// any element except the last one
		// use next sibling as an anchor
		if nextToken := astnav.GetTokenAtPosition(sourceFile, after.End()); nextToken != nil && isSeparator(after, nextToken) {
			// for list
			// a, b, c
			// create change for adding 'e' after 'a' as
			// - find start of next element after a (it is b)
			// - use next element start as start and end position in final change
			// - build text of change by formatting the text of node + whitespace trivia of b

			// in multiline case it will work as
			//   a,
			//   b,
			//   c,
			// result - '*' denotes leading trivia that will be inserted after new text (displayed as '#')
			//   a,
			//   insertedtext<separator>#
			// ###b,
			//   c,
			nextNode := containingList[index+1]
			startPos := scanner.SkipTriviaEx(sourceFile.Text(), nextNode.Pos(), &scanner.SkipTriviaOptions{StopAfterLineBreak: true, StopAtComments: false})

			// write separator and leading trivia of the next element as suffix
			suffix := scanner.TokenToString(nextToken.Kind) + sourceFile.Text()[nextNode.End():startPos]
			ct.insertNodeAt(sourceFile, core.TextPos(startPos), newNode, changeNodeOptions{suffix: suffix})
		}
		return
	}

	afterStart := astnav.GetStartOfNode(after, sourceFile, false)
	lineMap := sourceFile.LineMap()
	afterStartLinePosition := scanner.GetLineStartPositionForPosition(afterStart, lineMap)

	// insert element after the last element in the list that has more than one item
	// pick the element preceding the after element to:
	// - pick the separator
	// - determine if list is a multiline
	multilineList := false

	// if list has only one element then we'll format is as multiline if node has comment in trailing trivia, or as singleline otherwise
	// i.e. var x = 1 // this is x
	//     | new element will be inserted at this position
	separator := ast.KindCommaToken // SyntaxKind.CommaToken | SyntaxKind.SemicolonToken
	if len(containingList) != 1 {
		// otherwise, if list has more than one element, pick separator from the list
		tokenBeforeInsertPosition := astnav.FindPrecedingToken(sourceFile, after.Pos())
		separator = core.IfElse(isSeparator(after, tokenBeforeInsertPosition), tokenBeforeInsertPosition.Kind, ast.KindCommaToken)
		// determine if list is multiline by checking lines of after element and element that precedes it.
		afterMinusOneStartLinePosition := scanner.GetLineStartPositionForPosition(astnav.GetStartOfNode(containingList[index-1], sourceFile, false), lineMap)
		multilineList = afterMinusOneStartLinePosition != afterStartLinePosition
	}
	if hasCommentsBeforeLineBreak(sourceFile.Text(), after.End()) || printer.GetLinesBetweenPositions(sourceFile, containingList[0].Pos(), containingList[len(containingList)-1].End()) != 0 {
		// in this case we'll always treat containing list as multiline
		multilineList = true
	}

	separatorString := scanner.TokenToString(separator)
	end := ct.ls.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(after.End()))
	if !multilineList {
		ct.replaceRange(sourceFile, lsproto.Range{Start: end, End: end}, newNode, changeNodeOptions{prefix: separatorString})
		return
	}

	// insert separator immediately following the 'after' node to preserve comments in trailing trivia
	// !!! formatcontext
	ct.replaceRange(sourceFile, lsproto.Range{Start: end, End: end}, sourceFile.GetOrCreateToken(separator, after.End(), after.End()+len(separatorString), after.Parent), changeNodeOptions{})
	// use the same indentation as 'after' item
	indentation := format.FindFirstNonWhitespaceColumn(afterStartLinePosition, afterStart, sourceFile, ct.formatSettings)
	// insert element before the line break on the line that contains 'after' element
	insertPos := scanner.SkipTriviaEx(sourceFile.Text(), after.End(), &scanner.SkipTriviaOptions{StopAfterLineBreak: true, StopAtComments: false})
	// find position before "\n" or "\r\n"
	for insertPos != after.End() && stringutil.IsLineBreak(rune(sourceFile.Text()[insertPos-1])) {
		insertPos--
	}
	insertLSPos := ct.ls.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(insertPos))
	ct.replaceRange(
		sourceFile,
		lsproto.Range{Start: insertLSPos, End: insertLSPos},
		newNode,
		changeNodeOptions{
			indentation: ptrTo(indentation),
			prefix:      ct.newLine,
		},
	)
}

func (ct *changeTracker) insertAtTopOfFile(sourceFile *ast.SourceFile, insert []*ast.Statement, blankLineBetween bool) {
	pos := ct.getInsertionPositionAtSourceFileTop(sourceFile)
	options := changeNodeOptions{}
	if pos != 0 {
		options.prefix = ct.newLine
	}
	if !stringutil.IsLineBreak(rune(sourceFile.Text()[pos])) {
		options.suffix = ct.newLine
	}
	if blankLineBetween {
		options.suffix += ct.newLine
	}

	if len(insert) == 0 {
		ct.insertNodeAt(sourceFile, core.TextPos(pos), insert[0], options)
	} else {
		ct.insertNodesAt(sourceFile, core.TextPos(pos), insert, options)
	}
}

// method on the changeTracker because use of converters
func (ct *changeTracker) getAdjustedRange(sourceFile *ast.SourceFile, startNode *ast.Node, endNode *ast.Node, leadingOption leadingTriviaOption, trailingOption trailingTriviaOption) lsproto.Range {
	return *ct.ls.createLspRangeFromBounds(
		ct.getAdjustedStartPosition(sourceFile, startNode, leadingOption, false),
		ct.getAdjustedEndPosition(sourceFile, endNode, trailingOption),
		sourceFile,
	)
}

// method on the changeTracker because use of converters
func (ct *changeTracker) getAdjustedStartPosition(sourceFile *ast.SourceFile, node *ast.Node, leadingOption leadingTriviaOption, hasTrailingComment bool) int {
	if leadingOption == leadingTriviaOptionJSDoc {
		if JSDocComments := parser.GetJSDocCommentRanges(ct.NodeFactory, nil, node, sourceFile.Text()); len(JSDocComments) > 0 {
			return scanner.GetLineStartPositionForPosition(JSDocComments[0].Pos(), sourceFile.LineMap())
		}
	}

	start := astnav.GetStartOfNode(node, sourceFile, false)
	lineStarts := sourceFile.LineMap()
	startOfLinePos := scanner.GetLineStartPositionForPosition(start, lineStarts)

	switch leadingOption {
	case leadingTriviaOptionExclude:
		return start
	case leadingTriviaOptionStartLine:
		if ast.NodeRangeContainsPosition(node, startOfLinePos) {
			return startOfLinePos
		}
		return start
	}

	fullStart := node.Pos()
	if fullStart == start {
		return start
	}
	fullStartLineIndex := scanner.ComputeLineOfPosition(lineStarts, fullStart)
	fullStartLinePos := scanner.GetLineStartPositionForPosition(fullStart, lineStarts)
	if startOfLinePos == fullStartLinePos {
		// full start and start of the node are on the same line
		//   a,     b;
		//    ^     ^
		//    |   start
		// fullstart
		// when b is replaced - we usually want to keep the leading trvia
		// when b is deleted - we delete it
		if leadingOption == leadingTriviaOptionIncludeAll {
			return fullStart
		}
		return start
	}

	// if node has a trailing comments, use comment end position as the text has already been included.
	if hasTrailingComment {
		// Check first for leading comments as if the node is the first import, we want to exclude the trivia;
		// otherwise we get the trailing comments.
		comments := slices.Collect(scanner.GetLeadingCommentRanges(ct.NodeFactory, sourceFile.Text(), fullStart))
		if len(comments) == 0 {
			comments = slices.Collect(scanner.GetTrailingCommentRanges(ct.NodeFactory, sourceFile.Text(), fullStart))
		}
		if len(comments) > 0 {
			return scanner.SkipTriviaEx(sourceFile.Text(), comments[0].End(), &scanner.SkipTriviaOptions{StopAfterLineBreak: true, StopAtComments: true})
		}
	}

	// get start position of the line following the line that contains fullstart position
	// (but only if the fullstart isn't the very beginning of the file)
	nextLineStart := core.IfElse(fullStart > 0, 1, 0)
	adjustedStartPosition := int(lineStarts[fullStartLineIndex+nextLineStart])
	// skip whitespaces/newlines
	adjustedStartPosition = scanner.SkipTriviaEx(sourceFile.Text(), adjustedStartPosition, &scanner.SkipTriviaOptions{StopAtComments: true})
	return int(lineStarts[scanner.ComputeLineOfPosition(lineStarts, adjustedStartPosition)])
}

// method on the changeTracker because of converters
// Return the end position of a multiline comment of it is on another line; otherwise returns `undefined`;
func (ct *changeTracker) getEndPositionOfMultilineTrailingComment(sourceFile *ast.SourceFile, node *ast.Node, trailingOpt trailingTriviaOption) int {
	if trailingOpt == trailingTriviaOptionInclude {
		// If the trailing comment is a multiline comment that extends to the next lines,
		// return the end of the comment and track it for the next nodes to adjust.
		lineStarts := sourceFile.LineMap()
		nodeEndLine := scanner.ComputeLineOfPosition(lineStarts, node.End())
		for comment := range scanner.GetTrailingCommentRanges(ct.NodeFactory, sourceFile.Text(), node.End()) {
			// Single line can break the loop as trivia will only be this line.
			// Comments on subsequest lines are also ignored.
			if comment.Kind == ast.KindSingleLineCommentTrivia || scanner.ComputeLineOfPosition(lineStarts, comment.Pos()) > nodeEndLine {
				break
			}

			// Get the end line of the comment and compare against the end line of the node.
			// If the comment end line position and the multiline comment extends to multiple lines,
			// then is safe to return the end position.
			if commentEndLine := scanner.ComputeLineOfPosition(lineStarts, comment.End()); commentEndLine > nodeEndLine {
				return scanner.SkipTriviaEx(sourceFile.Text(), comment.End(), &scanner.SkipTriviaOptions{StopAfterLineBreak: true, StopAtComments: true})
			}
		}
	}

	return 0
}

// method on the changeTracker because of converters
func (ct *changeTracker) getAdjustedEndPosition(sourceFile *ast.SourceFile, node *ast.Node, trailingTriviaOption trailingTriviaOption) int {
	if trailingTriviaOption == trailingTriviaOptionExclude {
		return node.End()
	}
	if trailingTriviaOption == trailingTriviaOptionExcludeWhitespace {
		if comments := slices.AppendSeq(
			slices.Collect(scanner.GetTrailingCommentRanges(ct.NodeFactory, sourceFile.Text(), node.End())),
			scanner.GetLeadingCommentRanges(ct.NodeFactory, sourceFile.Text(), node.End()),
		); len(comments) > 0 {
			if realEnd := comments[len(comments)-1].End(); realEnd != 0 {
				return realEnd
			}
		}
		return node.End()
	}

	if multilineEndPosition := ct.getEndPositionOfMultilineTrailingComment(sourceFile, node, trailingTriviaOption); multilineEndPosition != 0 {
		return multilineEndPosition
	}

	newEnd := scanner.SkipTriviaEx(sourceFile.Text(), node.End(), &scanner.SkipTriviaOptions{StopAfterLineBreak: true})

	if newEnd != node.End() && (trailingTriviaOption == trailingTriviaOptionInclude || stringutil.IsLineBreak(rune(sourceFile.Text()[newEnd-1]))) {
		return newEnd
	}
	return node.End()
}

// ============= utilities =============

func hasCommentsBeforeLineBreak(text string, start int) bool {
	for _, ch := range []rune(text[start:]) {
		if !stringutil.IsWhiteSpaceSingleLine(ch) {
			return ch == '/'
		}
	}
	return false
}

func needSemicolonBetween(a, b *ast.Node) bool {
	return (ast.IsPropertySignatureDeclaration(a) || ast.IsPropertyDeclaration(a)) &&
		ast.IsClassOrTypeElement(b) &&
		b.Name().Kind == ast.KindComputedPropertyName ||
		ast.IsStatementButNotDeclaration(a) &&
			ast.IsStatementButNotDeclaration(b) // TODO: only if b would start with a `(` or `[`
}

func (ct *changeTracker) getInsertionPositionAtSourceFileTop(sourceFile *ast.SourceFile) int {
	var lastPrologue *ast.Node
	for _, node := range sourceFile.Statements.Nodes {
		if ast.IsPrologueDirective(node) {
			lastPrologue = node
		} else {
			break
		}
	}

	position := 0
	text := sourceFile.Text()
	advancePastLineBreak := func() {
		if position >= len(text) {
			return
		}
		if char := rune(text[position]); stringutil.IsLineBreak(char) {
			position++
			if position < len(text) && char == '\r' && rune(text[position]) == '\n' {
				position++
			}
		}
	}
	if lastPrologue != nil {
		position = lastPrologue.End()
		advancePastLineBreak()
		return position
	}

	shebang := scanner.GetShebang(text)
	if shebang != "" {
		position = len(shebang)
		advancePastLineBreak()
	}

	ranges := slices.Collect(scanner.GetLeadingCommentRanges(ct.NodeFactory, text, position))
	if len(ranges) == 0 {
		return position
	}
	// Find the first attached comment to the first node and add before it
	var lastComment *ast.CommentRange
	pinnedOrTripleSlash := false
	firstNodeLine := -1

	lenStatements := len(sourceFile.Statements.Nodes)
	lineMap := sourceFile.LineMap()
	for _, r := range ranges {
		if r.Kind == ast.KindMultiLineCommentTrivia {
			if printer.IsPinnedComment(text, r) {
				lastComment = &r
				pinnedOrTripleSlash = true
				continue
			}
		} else if printer.IsRecognizedTripleSlashComment(text, r) {
			lastComment = &r
			pinnedOrTripleSlash = true
			continue
		}

		if lastComment != nil {
			// Always insert after pinned or triple slash comments
			if pinnedOrTripleSlash {
				break
			}

			// There was a blank line between the last comment and this comment.
			// This comment is not part of the copyright comments
			commentLine := scanner.ComputeLineOfPosition(lineMap, r.Pos())
			lastCommentEndLine := scanner.ComputeLineOfPosition(lineMap, lastComment.End())
			if commentLine >= lastCommentEndLine+2 {
				break
			}
		}

		if lenStatements > 0 {
			if firstNodeLine == -1 {
				firstNodeLine = scanner.ComputeLineOfPosition(lineMap, astnav.GetStartOfNode(sourceFile.Statements.Nodes[0], sourceFile, false))
			}
			commentEndLine := scanner.ComputeLineOfPosition(lineMap, r.End())
			if firstNodeLine < commentEndLine+2 {
				break
			}
		}
		lastComment = &r
		pinnedOrTripleSlash = false
	}

	if lastComment != nil {
		position = lastComment.End()
		advancePastLineBreak()
	}
	return position
}
