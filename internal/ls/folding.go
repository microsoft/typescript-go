package ls

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func (l *LanguageService) ProvideFoldingRange(ctx context.Context, documentURI lsproto.DocumentUri) []*lsproto.FoldingRange {
	_, sourceFile := l.getProgramAndFile(documentURI)
	res := l.addNodeOutliningSpans(sourceFile)
	res = append(res, l.addRegionOutliningSpans(sourceFile)...)
	sort.Slice(res, func(i, j int) bool {
		if res[i] == nil && res[j] == nil {
			return false
		}
		if res[i] == nil {
			return false
		}
		if res[j] == nil {
			return true
		}
		return res[i].StartLine < res[j].StartLine
	})
	return res
}

func (l *LanguageService) addNodeOutliningSpans(sourceFile *ast.SourceFile) []*lsproto.FoldingRange {
	depthRemaining := 40
	current := 0
	// Includes the EOF Token so that comments which aren't attached to statements are included
	statements := sourceFile.Statements //!!! sourceFile.endOfFileToken
	n := len(statements.Nodes)
	foldingRange := make([]*lsproto.FoldingRange, n)
	for current < n {
		for current < n && !ast.IsAnyImportSyntax(statements.Nodes[current]) {
			foldingRange = append(foldingRange, visitNode(statements.Nodes[current], depthRemaining, sourceFile, l)...)
			current++
		}
		if current == n {
			break
		}
		firstImport := current
		for current < n && ast.IsAnyImportSyntax(statements.Nodes[current]) {
			foldingRange = append(foldingRange, visitNode(statements.Nodes[current], depthRemaining, sourceFile, l)...)
			current++
		}
		lastImport := current - 1
		if lastImport != firstImport {
			foldingRangeKind := lsproto.FoldingRangeKindImports
			foldingRange = append(foldingRange, createFoldingRangeFromBounds(
				astnav.GetStartOfNode(findChildOfKind(statements.Nodes[firstImport],
					ast.KindImportKeyword, sourceFile), sourceFile, false /*includeJSDoc*/),
				statements.Nodes[lastImport].End(),
				&foldingRangeKind,
				sourceFile,
				l))
		}
	}
	return foldingRange
}

func (l *LanguageService) addRegionOutliningSpans(sourceFile *ast.SourceFile) []*lsproto.FoldingRange {
	regions := []*lsproto.FoldingRange{}
	lineStarts := scanner.GetLineStarts(sourceFile)
	for _, currentLineStart := range lineStarts {
		lineEnd := scanner.GetLineEndOfPosition(sourceFile, int(currentLineStart))
		lineText := sourceFile.Text()[currentLineStart:lineEnd]
		result := parseRegionDelimiter(lineText)
		if result == nil || isInComment(sourceFile, int(currentLineStart), nil) != nil {
			continue
		}
		if result.isStart {
			span := l.createLspRangeFromBounds(strings.Index(sourceFile.Text()[currentLineStart:lineEnd], "//"), lineEnd, sourceFile)
			foldingRangeKindRegion := lsproto.FoldingRangeKindRegion
			collapsedTest := "#region"
			if result.name != "" {
				collapsedTest = result.name
			}
			regions = append(regions, createFoldingRange(span, &foldingRangeKindRegion, nil, collapsedTest))
		} else {
			// if len(regions) > 0 {
			// 	region := regions[len(regions)-1]
			// 	regions = regions[:len(regions)-1]
			// 	if region != nil {
			// 		region.StartLine = uint32(lineEnd - int(region.StartLine)) // !!! test
			// 		region.EndLine = lineEnd - region.HintSpan.Start
			// 		out = append(out, region)
			// 	}
			// }
		}
	}
	return regions
}

func visitNode(n *ast.Node, depthRemaining int, sourceFile *ast.SourceFile, l *LanguageService) []*lsproto.FoldingRange {
	if depthRemaining == 0 {
		return nil
	}
	// cancellationToken.throwIfCancellationRequested();
	var foldingRange []*lsproto.FoldingRange
	if ast.IsDeclaration(n) || ast.IsVariableStatement(n) || ast.IsReturnStatement(n) || ast.IsCallOrNewExpression(n) || n.Kind == ast.KindEndOfFile {
		foldingRange = append(foldingRange, addOutliningForLeadingCommentsForNode(n, sourceFile)...)
	}
	if ast.IsFunctionLike(n) && n.Parent != nil && ast.IsBinaryExpression(n.Parent) && n.Parent.AsBinaryExpression().Left != nil && ast.IsPropertyAccessExpression(n.Parent.AsBinaryExpression().Left) {
		foldingRange = append(foldingRange, addOutliningForLeadingCommentsForNode(n.Parent.AsBinaryExpression().Left, sourceFile)...)
	}
	if ast.IsBlock(n) || ast.IsModuleBlock(n) {
		statements := n.Statements()
		foldingRange = append(foldingRange, addOutliningForLeadingCommentsForPos(statements[len(statements)-1].End(), sourceFile)...)
	}
	if ast.IsClassLike(n) || ast.IsInterfaceDeclaration(n) {
		members := n.Members()
		foldingRange = append(foldingRange, addOutliningForLeadingCommentsForPos(members[len(members)-1].End(), sourceFile)...)
	}

	span := getOutliningSpanForNode(n, sourceFile, l)
	if span != nil {
		foldingRange = append(foldingRange, span)
	}

	depthRemaining--
	if ast.IsCallExpression(n) {
		depthRemaining++
		expressionNodes := visitNode(n.Expression(), depthRemaining, sourceFile, l)
		if expressionNodes != nil {
			foldingRange = append(foldingRange, expressionNodes...)
		}
		depthRemaining--
		for _, arg := range n.Arguments() {
			if arg != nil {
				foldingRange = append(foldingRange, visitNode(arg, depthRemaining, sourceFile, l)...)
			}
		}
		typeArguments := n.TypeArguments()
		for _, typeArg := range typeArguments {
			if typeArg != nil {
				foldingRange = append(foldingRange, visitNode(typeArg, depthRemaining, sourceFile, l)...)
			}
		}
	} else if ast.IsIfStatement(n) && n.AsIfStatement().ElseStatement != nil && ast.IsIfStatement(n.AsIfStatement().ElseStatement) {
		// Consider an 'else if' to be on the same depth as the 'if'.
		ifStatement := n.AsIfStatement()
		expressionNodes := visitNode(n.Expression(), depthRemaining, sourceFile, l)
		if expressionNodes != nil {
			foldingRange = append(foldingRange, expressionNodes...)
		}
		thenNode := visitNode(ifStatement.ThenStatement, depthRemaining, sourceFile, l)
		if thenNode != nil {
			foldingRange = append(foldingRange, thenNode...)
		}
		depthRemaining++
		elseNode := visitNode(ifStatement.ElseStatement, depthRemaining, sourceFile, l)
		if elseNode != nil {
			foldingRange = append(foldingRange, elseNode...)
		}
		depthRemaining--
	} else {
		var visit func(node *ast.Node) bool
		visit = func(node *ast.Node) bool {
			childNode := visitNode(node, depthRemaining, sourceFile, l)
			if childNode != nil {
				foldingRange = append(foldingRange, childNode...)
			}
			return false
		}
		n.ForEachChild(visit)
	}
	depthRemaining++
	return foldingRange
}

func addOutliningForLeadingCommentsForNode(n *ast.Node, sourceFile *ast.SourceFile) []*lsproto.FoldingRange {
	if ast.IsJsxText(n) {
		return nil
	}
	return addOutliningForLeadingCommentsForPos(n.Pos(), sourceFile)
}

func addOutliningForLeadingCommentsForPos(pos int, sourceFile *ast.SourceFile) []*lsproto.FoldingRange {
	c := &printer.EmitContext{}
	comments := scanner.GetLeadingCommentRanges(&printer.NewNodeFactory(c).NodeFactory, sourceFile.Text(), pos)
	if comments == nil {
		return nil
	}

	foldingRange := []*lsproto.FoldingRange{}
	firstSingleLineCommentStart := -1
	lastSingleLineCommentEnd := -1
	singleLineCommentCount := 0
	foldingRangeKindComment := lsproto.FoldingRangeKindComment

	combineAndAddMultipleSingleLineComments := func() *lsproto.FoldingRange {
		// Only outline spans of two or more consecutive single line comments
		if singleLineCommentCount > 1 {

			return createFoldingRangeFromBounds(
				firstSingleLineCommentStart, lastSingleLineCommentEnd, &foldingRangeKindComment, sourceFile, nil)
		}
		return nil
	}

	sourceText := sourceFile.Text()
	for comment := range comments {
		pos := comment.Pos()
		end := comment.End()
		// cancellationToken.throwIfCancellationRequested();
		switch comment.Kind {
		case ast.KindSingleLineCommentTrivia:
			// never fold region delimiters into single-line comment regions
			commentText := sourceText[pos:end]
			if parseRegionDelimiter(commentText) != nil {
				comments := combineAndAddMultipleSingleLineComments()
				if comments != nil {
					foldingRange = append(foldingRange, comments)
				}
				singleLineCommentCount = 0
				break
			}

			// For single line comments, combine consecutive ones (2 or more) into
			// a single span from the start of the first till the end of the last
			if singleLineCommentCount == 0 {
				firstSingleLineCommentStart = pos
			}
			lastSingleLineCommentEnd = end
			singleLineCommentCount++
			break
		case ast.KindMultiLineCommentTrivia:
			comments := combineAndAddMultipleSingleLineComments()
			if comments != nil {
				foldingRange = append(foldingRange, comments)
			}
			foldingRange = append(foldingRange, createFoldingRangeFromBounds(pos, end, &foldingRangeKindComment, sourceFile, nil))
			singleLineCommentCount = 0
			break
		default:
			// Debug.assertNever(kind);
		}
	}
	return foldingRange
}

var regionDelimiterRegExp = regexp.MustCompile(`^#(end)?region(.*)\r?$`)

type regionDelimiterResult struct {
	isStart bool
	name    string
}

func parseRegionDelimiter(lineText string) *regionDelimiterResult {
	// We trim the leading whitespace and // without the regex since the
	// multiple potential whitespace matches can make for some gnarly backtracking behavior
	lineText = strings.TrimLeft(lineText, " \t")
	if !strings.HasPrefix(lineText, "//") {
		return nil
	}
	lineText = strings.TrimLeft(lineText[2:], " \t")
	result := regionDelimiterRegExp.FindStringSubmatch(lineText)
	if result != nil {
		return &regionDelimiterResult{
			isStart: result[1] == "",
			name:    strings.TrimSpace(result[2]),
		}
	}
	return nil
}

func getOutliningSpanForNode(n *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	switch n.Kind {
	case ast.KindBlock:
		if ast.IsFunctionLike(n.Parent) {
			return functionSpan(n.Parent, n, sourceFile, l)
		}
		// Check if the block is standalone, or 'attached' to some parent statement.
		// If the latter, we want to collapse the block, but consider its hint span
		// to be the entire span of the parent.
		switch n.Parent.Kind {
		case ast.KindDoStatement, ast.KindForInStatement, ast.KindForOfStatement, ast.KindForStatement, ast.KindIfStatement, ast.KindWhileStatement, ast.KindWithStatement, ast.KindCatchClause:
			return spanForNode(n, n.Parent, ast.KindOpenBraceToken, true /*useFullStart */, sourceFile, l)
		case ast.KindTryStatement:
			// Could be the try-block, or the finally-block.
			tryStatement := n.Parent.AsTryStatement()
			if tryStatement.TryBlock == n {
				return spanForNode(n, n.Parent, ast.KindOpenBraceToken, true /*useFullStart */, sourceFile, l)
			} else if tryStatement.FinallyBlock == n {
				node := findChildOfKind(n.Parent, ast.KindFinallyKeyword, sourceFile)
				if node != nil {
					return spanForNode(n, node, ast.KindOpenBraceToken, true /*useFullStart */, sourceFile, l)
				}
			}
		default:
			// Block was a standalone block.  In this case we want to only collapse
			// the span of the block, independent of any parent span.
			return createFoldingRange(l.createLspRangeFromNode(n, sourceFile), nil, nil, "")
		}
	case ast.KindModuleBlock:
		return spanForNode(n, n.Parent, ast.KindOpenBraceToken, true /*useFullStart */, sourceFile, l)
	case ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration, ast.KindEnumDeclaration, ast.KindCaseBlock, ast.KindTypeLiteral, ast.KindObjectBindingPattern:
		return spanForNode(n, n, ast.KindOpenBraceToken, true /*useFullStart */, sourceFile, l)
	case ast.KindTupleType:
		return spanForNode(n, n, ast.KindOpenBracketToken, !ast.IsTupleTypeNode(n.Parent) /*useFullStart */, sourceFile, l)
	case ast.KindCaseClause, ast.KindDefaultClause:
		return spanForNodeArray(n.AsCaseOrDefaultClause().Statements, sourceFile, l)
	case ast.KindObjectLiteralExpression:
		return spanForNode(n, n, ast.KindOpenBraceToken, !ast.IsArrayLiteralExpression(n.Parent) && !ast.IsCallExpression(n.Parent) /*useFullStart */, sourceFile, l)
	case ast.KindArrayLiteralExpression:
		return spanForNode(n, n, ast.KindOpenBracketToken, !ast.IsArrayLiteralExpression(n.Parent) && !ast.IsCallExpression(n.Parent) /*useFullStart */, sourceFile, l)
	case ast.KindJsxElement, ast.KindJsxFragment:
		return spanForJSXElement(n, sourceFile, l)
	case ast.KindJsxSelfClosingElement, ast.KindJsxOpeningElement:
		return spanForJSXAttributes(n, sourceFile, l)
	case ast.KindTemplateExpression, ast.KindNoSubstitutionTemplateLiteral:
		return spanForTemplateLiteral(n, sourceFile, l)
	case ast.KindArrayBindingPattern:
		return spanForNode(n, n, ast.KindOpenBracketToken, !ast.IsBindingElement(n.Parent) /*useFullStart */, sourceFile, l)
	case ast.KindArrowFunction:
		return spanForArrowFunction(n, sourceFile, l)
	case ast.KindCallExpression:
		return spanForCallExpression(n, sourceFile, l)
	case ast.KindParenthesizedExpression:
		return spanForParenthesizedExpression(n, sourceFile, l)
	case ast.KindNamedImports, ast.KindNamedExports, ast.KindImportAttributes:
		return spanForImportExportElements(n, sourceFile, l)
	}
	return nil
}

func spanForImportExportElements(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	if len(node.Elements()) == 0 {
		return nil
	}
	openToken := findChildOfKind(node, ast.KindOpenBraceToken, sourceFile)
	closeToken := findChildOfKind(node, ast.KindCloseBraceToken, sourceFile)
	if openToken == nil || closeToken == nil || printer.PositionsAreOnSameLine(openToken.Pos(), closeToken.End(), sourceFile) {
		return nil
	}
	return spanBetweenTokens(openToken, closeToken, node, sourceFile, false /*useFullStart*/, l)
}

func spanForParenthesizedExpression(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	start := astnav.GetStartOfNode(node, sourceFile, false /*includeJSDoc*/)
	if printer.PositionsAreOnSameLine(start, node.End(), sourceFile) {
		return nil
	}
	textRange := l.createLspRangeFromBounds(start, node.End(), sourceFile)
	return createFoldingRange(textRange, nil, l.createLspRangeFromNode(node, sourceFile), "")
}

func spanForCallExpression(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	if node.AsCallExpression().Arguments == nil {
		return nil
	}
	openToken := findChildOfKind(node, ast.KindOpenParenToken, sourceFile)
	closeToken := findChildOfKind(node, ast.KindCloseParenToken, sourceFile)
	if openToken == nil || closeToken == nil || printer.PositionsAreOnSameLine(openToken.Pos(), closeToken.End(), sourceFile) {
		return nil
	}

	return spanBetweenTokens(openToken, closeToken, node, sourceFile, true /*useFullStart*/, l)
}

func spanForArrowFunction(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	arrowFunctionNode := node.AsArrowFunction()
	if ast.IsBlock(arrowFunctionNode.Body) || ast.IsParenthesizedExpression(arrowFunctionNode.Body) || printer.PositionsAreOnSameLine(arrowFunctionNode.Body.Pos(), arrowFunctionNode.Body.End(), sourceFile) {
		return nil
	}
	textRange := l.createLspRangeFromBounds(arrowFunctionNode.Body.Pos(), arrowFunctionNode.Body.End(), sourceFile)
	return createFoldingRange(textRange, nil, l.createLspRangeFromNode(node, sourceFile), "")
}

func spanForTemplateLiteral(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	if node.Kind == ast.KindNoSubstitutionTemplateLiteral && len(node.Text()) == 0 {
		return nil
	}
	return createFoldingRangeFromBounds(astnav.GetStartOfNode(node, sourceFile, false /*includeJSDoc*/), node.End(), nil, sourceFile, l)
}

func spanForJSXElement(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	var openingElement *ast.Node
	if node.Kind == ast.KindJsxElement {
		openingElement = node.AsJsxElement().OpeningElement
	} else {
		openingElement = node.AsJsxFragment().OpeningFragment
	}
	textRange := l.createLspRangeFromBounds(openingElement.Pos(), openingElement.End(), sourceFile)
	tagName := openingElement.TagName().Text()
	var bannerText strings.Builder
	if node.Kind == ast.KindJsxElement {
		bannerText.WriteString("<")
		bannerText.WriteString(tagName)
		bannerText.WriteString(">...</")
		bannerText.WriteString(tagName)
		bannerText.WriteString(">")
	} else {
		bannerText.WriteString("<>...</>")
	}

	return createFoldingRange(textRange, nil, nil, bannerText.String())
}

func spanForJSXAttributes(node *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	var attributes *ast.JsxAttributesNode
	if node.Kind == ast.KindJsxSelfClosingElement {
		attributes = node.AsJsxSelfClosingElement().Attributes
	} else {
		attributes = node.AsJsxOpeningElement().Attributes
	}
	if len(attributes.Properties()) == 0 {
		return nil
	}
	return createFoldingRangeFromBounds(astnav.GetStartOfNode(node, sourceFile, false /*includeJSDoc*/), node.End(), nil, sourceFile, l)
}

func spanForNodeArray(statements *ast.NodeList, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	if statements != nil && len(statements.Nodes) != 0 {
		return createFoldingRange(l.createLspRangeFromBounds(statements.Pos(), statements.End(), sourceFile), nil, nil, "")
	}
	return nil
}

func spanForNode(node *ast.Node, hintSpanNode *ast.Node, open ast.Kind, useFullStart bool, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	close := ast.KindCloseBraceToken
	if open != ast.KindOpenBraceToken {
		close = ast.KindCloseBracketToken
	}
	openToken := findChildOfKind(node, open, sourceFile)
	closeToken := findChildOfKind(node, close, sourceFile)
	if openToken != nil && closeToken != nil {
		return spanBetweenTokens(openToken, closeToken, hintSpanNode, sourceFile, useFullStart, l)
	}
	return nil
}

func spanBetweenTokens(openToken *ast.Node, closeToken *ast.Node, hintSpanNode *ast.Node, sourceFile *ast.SourceFile, useFullStart bool, l *LanguageService) *lsproto.FoldingRange {
	var textRange *lsproto.Range
	if useFullStart {
		textRange = l.createLspRangeFromBounds(openToken.Pos(), closeToken.End(), sourceFile)
	} else {
		textRange = l.createLspRangeFromBounds(astnav.GetStartOfNode(openToken, sourceFile, false /*includeJSDoc*/), closeToken.End(), sourceFile)
	}
	return createFoldingRange(textRange, nil, l.createLspRangeFromNode(hintSpanNode, sourceFile), "")
}

func createFoldingRange(textRange *lsproto.Range, foldingRangeKind *lsproto.FoldingRangeKind, hintRange *lsproto.Range, collapsedText string) *lsproto.FoldingRange {
	if hintRange == nil {
		hintRange = textRange
	}
	if collapsedText == "" {
		defaultText := "..."
		collapsedText = defaultText
	}
	return &lsproto.FoldingRange{
		StartLine:     textRange.Start.Line,
		EndLine:       textRange.End.Line, // !!! needs to be adjusted for in vscode repo
		Kind:          foldingRangeKind,
		CollapsedText: &collapsedText,
	}
}

// func adjustFoldingRange(textRange lsproto.Range, sourceFile *ast.SourceFile) {
// 	if textRange.End.Character > 0 {
// 		foldEndCharacter := sourceFile.Text()[textRange.End.Line : textRange.End]
// 	}
// }

func createFoldingRangeFromBounds(pos int, end int, foldingRangeKind *lsproto.FoldingRangeKind, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	return createFoldingRange(l.createLspRangeFromBounds(pos, end, sourceFile), foldingRangeKind, nil, "")
}

func functionSpan(node *ast.Node, body *ast.Node, sourceFile *ast.SourceFile, l *LanguageService) *lsproto.FoldingRange {
	openToken := tryGetFunctionOpenToken(node, body, sourceFile)
	closeToken := findChildOfKind(body, ast.KindCloseBraceToken, sourceFile)
	if openToken != nil && closeToken != nil {
		return spanBetweenTokens(openToken, closeToken, node, sourceFile, true /*useFullStart*/, l)
	}
	return nil
}

func tryGetFunctionOpenToken(node *ast.SignatureDeclaration, body *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	if isNodeArrayMultiLine(node.Parameters(), sourceFile) {
		openParenToken := findChildOfKind(node, ast.KindOpenParenToken, sourceFile)
		if openParenToken != nil {
			return openParenToken
		}
	}
	return findChildOfKind(body, ast.KindOpenBraceToken, sourceFile)
}

func isNodeArrayMultiLine(list []*ast.Node, sourceFile *ast.SourceFile) bool {
	if len(list) == 0 {
		return false
	}
	return !printer.PositionsAreOnSameLine(list[0].Pos(), list[len(list)-1].End(), sourceFile)
}
