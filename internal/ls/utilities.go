package ls

import (
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// !!! Shared (placeholder)
func isInString(file *ast.SourceFile, position int, previousToken *ast.Node) bool {
	if previousToken == nil {
		previousToken = astnav.FindPrecedingToken(file, position)
	}

	if previousToken != nil && isStringTextContainingNode(previousToken) {
	}

	return false
}

// !!! Shared (placeholder)
func isStringTextContainingNode(node *ast.Node) bool {
	return true
}

// !!! Shared (placeholder)
func getStartOfNode(node *ast.Node, file *ast.SourceFile) int {
	return node.Pos()
}

func tryGetImportFromModuleSpecifier(node *ast.StringLiteralLike) *ast.Node {
	switch node.Parent.Kind {
	case ast.KindImportDeclaration, ast.KindExportDeclaration, ast.KindJSDocImportTag:
		return node.Parent
	case ast.KindExternalModuleReference:
		return node.Parent.Parent
	case ast.KindCallExpression:
		if ast.IsImportCall(node.Parent) || ast.IsRequireCall(node.Parent, false /*requireStringLiteralLikeArgument*/) {
			return node.Parent
		}
		return nil
	case ast.KindLiteralType:
		if !ast.IsStringLiteral(node) {
			return nil
		}
		if ast.IsImportTypeNode(node.Parent.Parent) {
			return node.Parent.Parent
		}
		return nil
	}
	return nil
}

// !!! Shared (placeholder)
func isInComment(file *ast.SourceFile, position int, tokenAtPosition *ast.Node) *ast.CommentRange {
	return nil
}

// // Returns a non-nil comment range if the cursor at position in sourceFile is within a comment.
// // tokenAtPosition: must equal `getTokenAtPosition(sourceFile, position)`
// // predicate: additional predicate to test on the comment range.
// func isInComment(file *ast.SourceFile, position int, tokenAtPosition *ast.Node) *ast.CommentRange {
// 	return getRangeOfEnclosingComment(file, position, nil /*precedingToken*/, tokenAtPosition)
// }

// !!!
// Replaces last(node.getChildren(sourceFile))
func getLastChild(node *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	return nil
}

func getLastToken(node *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	if node == nil {
		return nil
	}

	assertHasRealPosition(node)

	lastChild := getLastChild(node, sourceFile)
	if lastChild == nil {
		return nil
	}

	if lastChild.Kind < ast.KindFirstNode {
		return lastChild
	} else {
		return getLastToken(lastChild, sourceFile)
	}
}

// !!!
func getFirstToken(node *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	return nil
}

func assertHasRealPosition(node *ast.Node) {
	if ast.PositionIsSynthesized(node.Pos()) || ast.PositionIsSynthesized(node.End()) {
		panic("Node must have a real position for this operation.")
	}
}

// !!!
func findChildOfKind(node *ast.Node, kind ast.Kind, sourceFile *ast.SourceFile) *ast.Node {
	return nil
}

// !!! Shared: placeholder
type PossibleTypeArgumentInfo struct {
	called         *ast.IdentifierNode
	nTypeArguments int
}

// !!! Shared: placeholder
func getPossibleTypeArgumentsInfo(tokenIn *ast.Node, sourceFile *ast.SourceFile) *PossibleTypeArgumentInfo {
	return nil
}

// !!! Shared: placeholder
func getPossibleGenericSignatures(called *ast.Expression, typeArgumentCount int, checker *checker.Checker) []*checker.Signature {
	return nil
}

func isInRightSideOfInternalImportEqualsDeclaration(node *ast.Node) bool {
	for node.Parent.Kind == ast.KindQualifiedName {
		node = node.Parent
	}

	return ast.IsInternalModuleImportEqualsDeclaration(node.Parent) && node.Parent.AsImportEqualsDeclaration().ModuleReference == node
}

func createLspRangeFromNode(node *ast.Node, file *ast.SourceFile) *lsproto.Range {
	return createLspRangeFromBounds(node.Pos(), node.End(), file)
}

func createLspRangeFromBounds(start, end int, file *ast.SourceFile) *lsproto.Range {
	// !!! needs converters access
	return nil
}

func quote(file *ast.SourceFile, preferences *UserPreferences, text string) string {
	// Editors can pass in undefined or empty string - we want to infer the preference in those cases.
	quotePreference := getQuotePreference(file, preferences)
	quoted, _ := core.StringifyJson(text, "" /*prefix*/, "" /*indent*/)
	if quotePreference == quotePreferenceSingle {
		strings.ReplaceAll(strings.ReplaceAll(core.StripQuotes(quoted), "'", `\'`), `\"`, `"`)
	}
	return quoted
}

type quotePreference int

const (
	quotePreferenceSingle quotePreference = iota
	quotePreferenceDouble
)

func getQuotePreference(file *ast.SourceFile, preferences *UserPreferences) quotePreference {
	// !!!
	return quotePreferenceDouble
}

func positionIsASICandidate(pos int, context *ast.Node, file *ast.SourceFile) bool {
	contextAncestor := ast.FindAncestorOrQuit(context, func(ancestor *ast.Node) ast.FindAncestorResult {
		if ancestor.End() != pos {
			return ast.FindAncestorQuit
		}

		return ast.ToFindAncestorResult(syntaxMayBeASICandidate(ancestor.Kind))
	})

	return contextAncestor != nil && nodeIsASICandidate(contextAncestor, file)
}

func syntaxMayBeASICandidate(kind ast.Kind) bool {
	return syntaxRequiresTrailingCommaOrSemicolonOrASI(kind) ||
		syntaxRequiresTrailingFunctionBlockOrSemicolonOrASI(kind) ||
		syntaxRequiresTrailingModuleBlockOrSemicolonOrASI(kind) ||
		syntaxRequiresTrailingSemicolonOrASI(kind)
}

func syntaxRequiresTrailingCommaOrSemicolonOrASI(kind ast.Kind) bool {
	return kind == ast.KindCallSignature ||
		kind == ast.KindConstructSignature ||
		kind == ast.KindIndexSignature ||
		kind == ast.KindPropertySignature ||
		kind == ast.KindMethodSignature
}

func syntaxRequiresTrailingFunctionBlockOrSemicolonOrASI(kind ast.Kind) bool {
	return kind == ast.KindFunctionDeclaration ||
		kind == ast.KindConstructor ||
		kind == ast.KindMethodDeclaration ||
		kind == ast.KindGetAccessor ||
		kind == ast.KindSetAccessor
}

func syntaxRequiresTrailingModuleBlockOrSemicolonOrASI(kind ast.Kind) bool {
	return kind == ast.KindModuleDeclaration
}

func syntaxRequiresTrailingSemicolonOrASI(kind ast.Kind) bool {
	return kind == ast.KindVariableStatement ||
		kind == ast.KindExpressionStatement ||
		kind == ast.KindDoStatement ||
		kind == ast.KindContinueStatement ||
		kind == ast.KindBreakStatement ||
		kind == ast.KindReturnStatement ||
		kind == ast.KindThrowStatement ||
		kind == ast.KindDebuggerStatement ||
		kind == ast.KindPropertyDeclaration ||
		kind == ast.KindTypeAliasDeclaration ||
		kind == ast.KindImportDeclaration ||
		kind == ast.KindImportEqualsDeclaration ||
		kind == ast.KindExportDeclaration ||
		kind == ast.KindNamespaceExportDeclaration ||
		kind == ast.KindExportAssignment
}

func nodeIsASICandidate(node *ast.Node, file *ast.SourceFile) bool {
	lastToken := getLastToken(node, file)
	if lastToken != nil && lastToken.Kind == ast.KindSemicolonToken {
		return false
	}

	if syntaxRequiresTrailingCommaOrSemicolonOrASI(node.Kind) {
		if lastToken != nil && lastToken.Kind == ast.KindCommaToken {
			return false
		}
	} else if syntaxRequiresTrailingModuleBlockOrSemicolonOrASI(node.Kind) {
		lastChild := getLastChild(node, file)
		if lastChild != nil && ast.IsModuleBlock(lastChild) {
			return false
		}
	} else if syntaxRequiresTrailingFunctionBlockOrSemicolonOrASI(node.Kind) {
		lastChild := getLastChild(node, file)
		if lastChild != nil && ast.IsFunctionBlock(lastChild) {
			return false
		}
	} else if !syntaxRequiresTrailingSemicolonOrASI(node.Kind) {
		return false
	}

	// See comment in parser's `parseDoStatement`
	if node.Kind == ast.KindDoStatement {
		return true
	}

	topNode := ast.FindAncestor(node, func(ancestor *ast.Node) bool { return ancestor.Parent == nil })
	nextToken := astnav.FindNextToken(node, topNode, file)
	if nextToken == nil || nextToken.Kind == ast.KindCloseBraceToken {
		return true
	}

	startLine, _ := scanner.GetLineAndCharacterOfPosition(file, node.End())
	endLine, _ := scanner.GetLineAndCharacterOfPosition(file, getStartOfNode(nextToken, file))
	return startLine != endLine
}

func isNonContextualKeyword(token ast.Kind) bool {
	return ast.IsKeywordKind(token) && !ast.IsContextualKeyword(token)
}

func probablyUsesSemicolons(file *ast.SourceFile) bool {
	withSemicolon := 0
	withoutSemicolon := 0
	nStatementsToObserve := 5

	var visit func(node *ast.Node) bool
	visit = func(node *ast.Node) bool {
		if syntaxRequiresTrailingSemicolonOrASI(node.Kind) {
			lastToken := getLastToken(node, file)
			if lastToken != nil && lastToken.Kind == ast.KindSemicolonToken {
				withSemicolon++
			} else {
				withoutSemicolon++
			}
		} else if syntaxRequiresTrailingCommaOrSemicolonOrASI(node.Kind) {
			lastToken := getLastToken(node, file)
			if lastToken != nil && lastToken.Kind == ast.KindSemicolonToken {
				withSemicolon++
			} else if lastToken != nil && lastToken.Kind != ast.KindCommaToken {
				lastTokenLine, _ := scanner.GetLineAndCharacterOfPosition(file, getStartOfNode(lastToken, file))
				nextTokenLine, _ := scanner.GetLineAndCharacterOfPosition(file, scanner.GetRangeOfTokenAtPosition(file, lastToken.End()).Pos())
				// Avoid counting missing semicolon in single-line objects:
				// `function f(p: { x: string /*no semicolon here is insignificant*/ }) {`
				if lastTokenLine != nextTokenLine {
					withoutSemicolon++
				}
			}
		}

		if withSemicolon+withoutSemicolon >= nStatementsToObserve {
			return true
		}

		return node.ForEachChild(visit)
	}

	file.ForEachChild(visit)

	// One statement missing a semicolon isn't sufficient evidence to say the user
	// doesn't want semicolons, because they may not even be done writing that statement.
	if withSemicolon == 0 && withoutSemicolon <= 1 {
		return true
	}

	// If even 2/5 places have a semicolon, the user probably wants semicolons
	return withSemicolon/withoutSemicolon > 1/nStatementsToObserve
}

var typeKeywords []ast.Kind = []ast.Kind{
	ast.KindAnyKeyword,
	ast.KindAssertsKeyword,
	ast.KindBigIntKeyword,
	ast.KindBooleanKeyword,
	ast.KindFalseKeyword,
	ast.KindInferKeyword,
	ast.KindKeyOfKeyword,
	ast.KindNeverKeyword,
	ast.KindNullKeyword,
	ast.KindNumberKeyword,
	ast.KindObjectKeyword,
	ast.KindReadonlyKeyword,
	ast.KindStringKeyword,
	ast.KindSymbolKeyword,
	ast.KindTypeOfKeyword,
	ast.KindTrueKeyword,
	ast.KindVoidKeyword,
	ast.KindUndefinedKeyword,
	ast.KindUniqueKeyword,
	ast.KindUnknownKeyword,
}

func isTypeKeyword(kind ast.Kind) bool {
	return slices.Contains(typeKeywords, kind)
}

// Returns a map of all names in the file to their positions.
// !!! cache this
func getNameTable(file *ast.SourceFile) map[string]int {
	nameTable := make(map[string]int)
	var walk func(node *ast.Node) bool

	walk = func(node *ast.Node) bool {
		if ast.IsIdentifier(node) && !isTagName(node) && node.Text() != "" ||
			ast.IsStringOrNumericLiteralLike(node) && literalIsName(node) ||
			ast.IsPrivateIdentifier(node) {
			text := node.Text()
			if _, ok := nameTable[text]; ok {
				nameTable[text] = -1
			} else {
				nameTable[text] = node.Pos()
			}
		}

		node.ForEachChild(walk)
		jsdocNodes := node.JSDoc(file)
		for _, jsdoc := range jsdocNodes {
			jsdoc.ForEachChild(walk)
		}
		return false
	}

	file.ForEachChild(walk)
	return nameTable
}

// We want to store any numbers/strings if they were a name that could be
// related to a declaration.  So, if we have 'import x = require("something")'
// then we want 'something' to be in the name table.  Similarly, if we have
// "a['propname']" then we want to store "propname" in the name table.
func literalIsName(node *ast.NumericOrStringLikeLiteral) bool {
	return ast.IsDeclarationName(node) ||
		node.Parent.Kind == ast.KindExternalModuleReference ||
		isArgumentOfElementAccessExpression(node) ||
		ast.IsLiteralComputedPropertyDeclarationName(node)
}

func isArgumentOfElementAccessExpression(node *ast.Node) bool {
	return node != nil && node.Parent != nil &&
		node.Parent.Kind == ast.KindElementAccessExpression &&
		node.Parent.AsElementAccessExpression().ArgumentExpression == node
}

func isTagName(node *ast.Node) bool {
	return node.Parent != nil && ast.IsJSDocTag(node.Parent) && node.Parent.TagName() == node
}
