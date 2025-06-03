package format

import (
	"context"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

func formatNodeLines(ctx context.Context, sourceFile *ast.SourceFile, node *ast.Node, requestKind FormatRequestKind) []core.TextChange {
	if node == nil {
		return nil
	}
	tokenStart := scanner.GetTokenPosOfNode(node, sourceFile, false)
	lineStart := getLineStartPositionForPosition(tokenStart, sourceFile)
	span := core.NewTextRange(lineStart, node.End())
	return FormatSpan(ctx, span, sourceFile, requestKind)
}

func FormatDocument(ctx context.Context, sourceFile *ast.SourceFile) []core.TextChange {
	return FormatSpan(ctx, core.NewTextRange(0, sourceFile.End()), sourceFile, FormatRequestKindFormatDocument)
}

func FormatSelection(ctx context.Context, sourceFile *ast.SourceFile, start int, end int) []core.TextChange {
	return FormatSpan(ctx, core.NewTextRange(getLineStartPositionForPosition(start, sourceFile), end), sourceFile, FormatRequestKindFormatSelection)
}

func FormatOnOpeningCurly(ctx context.Context, sourceFile *ast.SourceFile, position int) []core.TextChange {
	openingCurly := findImmediatelyPrecedingTokenOfKind(position, ast.KindOpenBraceToken, sourceFile)
	if openingCurly == nil {
		return nil
	}
	curlyBraceRange := openingCurly.Parent
	outermostNode := findOutermostNodeWithinListLevel(curlyBraceRange)
	/**
	 * We limit the span to end at the opening curly to handle the case where
	 * the brace matched to that just typed will be incorrect after further edits.
	 * For example, we could type the opening curly for the following method
	 * body without brace-matching activated:
	 * ```
	 * class C {
	 *     foo()
	 * }
	 * ```
	 * and we wouldn't want to move the closing brace.
	 */
	textRange := core.NewTextRange(getLineStartPositionForPosition(scanner.GetTokenPosOfNode(outermostNode, sourceFile, false), sourceFile), position)
	return FormatSpan(ctx, textRange, sourceFile, FormatRequestKindFormatOnOpeningCurlyBrace)
}

func FormatOnSemicolon(ctx context.Context, sourceFile *ast.SourceFile, position int) []core.TextChange {
	semicolon := findImmediatelyPrecedingTokenOfKind(position, ast.KindSemicolonToken, sourceFile)
	return formatNodeLines(ctx, sourceFile, findOutermostNodeWithinListLevel(semicolon), FormatRequestKindFormatOnSemicolon)
}

func FormatOnEnter(ctx context.Context, sourceFile *ast.SourceFile, position int) []core.TextChange {
	line, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, position)
	if line == 0 {
		return nil
	}
	// After the enter key, the cursor is now at a new line. The new line may or may not contain non-whitespace characters.
	// If the new line has only whitespaces, we won't want to format this line, because that would remove the indentation as
	// trailing whitespaces. So the end of the formatting span should be the later one between:
	//  1. the end of the previous line
	//  2. the last non-whitespace character in the current line
	endOfFormatSpan := scanner.GetEndLinePosition(sourceFile, line)
	for endOfFormatSpan > 0 {
		ch, s := utf8.DecodeRuneInString(sourceFile.Text()[endOfFormatSpan:])
		if s == 0 || stringutil.IsWhiteSpaceSingleLine(ch) { // on multibyte character keep backing up
			endOfFormatSpan--
			continue
		}
		break
	}

	span := core.NewTextRange(
		// get start position for the previous line
		int(scanner.GetLineStarts(sourceFile)[line-1]),
		// end value is exclusive so add 1 to the result
		endOfFormatSpan+1,
	)

	return FormatSpan(ctx, span, sourceFile, FormatRequestKindFormatOnEnter)
}
