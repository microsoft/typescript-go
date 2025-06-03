package format

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func rangeIsOnOneLine(node core.TextRange, file *ast.SourceFile) bool {
	startLine, _ := scanner.GetLineAndCharacterOfPosition(file, node.Pos())
	endLine, _ := scanner.GetLineAndCharacterOfPosition(file, node.End())
	return startLine == endLine
}

func getOpenTokenForList(node *ast.Node, list *ast.NodeList) ast.Kind {
	switch node.Kind {
	case ast.KindConstructor,
		ast.KindFunctionDeclaration,
		ast.KindFunctionExpression,
		ast.KindMethodDeclaration,
		ast.KindMethodSignature,
		ast.KindArrowFunction,
		ast.KindCallSignature,
		ast.KindConstructSignature,
		ast.KindFunctionType,
		ast.KindConstructorType,
		ast.KindGetAccessor,
		ast.KindSetAccessor:
		if node.TypeParameterList() == list {
			return ast.KindLessThanToken
		} else if node.ParameterList() == list {
			return ast.KindOpenParenToken
		}
	case ast.KindCallExpression, ast.KindNewExpression:
		if node.TypeArgumentList() == list {
			return ast.KindLessThanToken
		} else if node.ArgumentList() == list {
			return ast.KindOpenParenToken
		}
	case ast.KindClassDeclaration,
		ast.KindClassExpression,
		ast.KindInterfaceDeclaration,
		ast.KindTypeAliasDeclaration:
		if node.TypeParameterList() == list {
			return ast.KindLessThanToken
		}
	case ast.KindTypeReference,
		ast.KindTaggedTemplateExpression,
		ast.KindTypeQuery,
		ast.KindExpressionWithTypeArguments,
		ast.KindImportType:
		if node.TypeArgumentList() == list {
			return ast.KindLessThanToken
		}
	case ast.KindTypeLiteral:
		return ast.KindOpenBraceToken
	}

	return ast.KindUnknown
}

func getCloseTokenForOpenToken(kind ast.Kind) ast.Kind {
	// TODO: matches strada - seems like it could handle more pairs of braces, though? [] notably missing
	switch kind {
	case ast.KindOpenParenToken:
		return ast.KindCloseParenToken
	case ast.KindLessThanToken:
		return ast.KindGreaterThanToken
	case ast.KindOpenBraceToken:
		return ast.KindCloseBraceToken
	}
	return ast.KindUnknown
}

func getLineStartPositionForPosition(position int, sourceFile *ast.SourceFile) int {
	lineStarts := scanner.GetLineStarts(sourceFile)
	line, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, position)
	return int(lineStarts[line])
}

/**
 * Tests whether `child` is a grammar error on `parent`.
 * In strada, this also checked node arrays, but it is never acually called with one in practice.
 */
func isGrammarError(parent *ast.Node, child *ast.Node) bool {
	if ast.IsTypeParameterDeclaration(parent) {
		return child == parent.AsTypeParameter().Expression
	}
	if ast.IsPropertySignatureDeclaration(parent) {
		return child == parent.AsPropertySignatureDeclaration().Initializer
	}
	if ast.IsPropertyDeclaration(parent) {
		return ast.IsAutoAccessorPropertyDeclaration(parent) && child == parent.AsPropertyDeclaration().PostfixToken && child.Kind == ast.KindQuestionToken
	}
	if ast.IsPropertyAssignment(parent) {
		pa := parent.AsPropertyAssignment()
		return child == pa.PostfixToken || isGrammarErrorElement(&pa.Modifiers().NodeList, child, ast.IsModifierLike)
	}
	if ast.IsShorthandPropertyAssignment(parent) {
		sp := parent.AsShorthandPropertyAssignment()
		return child == sp.EqualsToken || child == sp.PostfixToken || isGrammarErrorElement(&parent.Modifiers().NodeList, child, ast.IsModifierLike)
	}
	if ast.IsMethodDeclaration(parent) {
		return child == parent.AsMethodDeclaration().PostfixToken && child.Kind == ast.KindExclamationToken
	}
	if ast.IsConstructorDeclaration(parent) {
		return child == parent.AsConstructorDeclaration().Type || isGrammarErrorElement(parent.AsConstructorDeclaration().TypeParameters, child, ast.IsTypeParameterDeclaration)
	}
	if ast.IsGetAccessorDeclaration(parent) {
		return isGrammarErrorElement(parent.AsGetAccessorDeclaration().TypeParameters, child, ast.IsTypeParameterDeclaration)
	}
	if ast.IsSetAccessorDeclaration(parent) {
		return child == parent.AsSetAccessorDeclaration().Type || isGrammarErrorElement(parent.AsSetAccessorDeclaration().TypeParameters, child, ast.IsTypeParameterDeclaration)
	}
	if ast.IsNamespaceExportDeclaration(parent) {
		return isGrammarErrorElement(&parent.AsNamespaceExportDeclaration().Modifiers().NodeList, child, ast.IsModifierLike)
	}
	return false
}

func isGrammarErrorElement(list *ast.NodeList, child *ast.Node, isPossibleElement func(node *ast.Node) bool) bool {
	if list == nil || len(list.Nodes) == 0 {
		return false
	}
	if !isPossibleElement(child) {
		return false
	}
	return slices.Contains(list.Nodes, child)
}
