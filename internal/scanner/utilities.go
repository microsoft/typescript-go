package scanner

import (
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func tokenIsIdentifierOrKeyword(token ast.Kind) bool {
	return token >= ast.KindIdentifier
}

func IdentifierToKeywordKind(node *ast.Identifier) ast.Kind {
	return textToKeyword[node.Text]
}

func GetSourceTextOfNodeFromSourceFile(sourceFile *ast.SourceFile, node *ast.Node, includeTrivia bool) string {
	return GetTextOfNodeFromSourceText(sourceFile.Text(), node, includeTrivia)
}

func GetTextOfNodeFromSourceText(sourceText string, node *ast.Node, includeTrivia bool) string {
	if ast.NodeIsMissing(node) {
		return ""
	}
	pos := node.Pos()
	if !includeTrivia {
		pos = SkipTrivia(sourceText, pos)
	}
	text := sourceText[pos:node.End()]
	// if (isJSDocTypeExpressionOrChild(node)) {
	//     // strip space + asterisk at line start
	//     text = text.split(/\r\n|\n|\r/).map(line => line.replace(/^\s*\*/, "").trimStart()).join("\n");
	// }
	return text
}

func GetTextOfNode(node *ast.Node) string {
	return GetSourceTextOfNodeFromSourceFile(ast.GetSourceFileOfNode(node), node, false /*includeTrivia*/)
}

func DeclarationNameToString(name *ast.Node) string {
	if name == nil || name.Pos() == name.End() {
		return "(Missing)"
	}
	return GetTextOfNode(name)
}

func IsIdentifierText(name string, languageVariant core.LanguageVariant) bool {
	ch, size := utf8.DecodeRuneInString(name)
	if !IsIdentifierStart(ch) {
		return false
	}
	for i := size; i < len(name); {
		ch, size = utf8.DecodeRuneInString(name[i:])
		if !IsIdentifierPartEx(ch, languageVariant) {
			return false
		}
		i += size
	}
	return true
}

func IsIntrinsicJsxName(name string) bool {
	return len(name) != 0 && (name[0] >= 'a' && name[0] <= 'z' || strings.ContainsRune(name, '-'))
}

func getErrorRangeForArrowFunction(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	pos := SkipTrivia(sourceFile.Text(), node.Pos())
	body := node.AsArrowFunction().Body
	if body != nil && body.Kind == ast.KindBlock {
		startLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.Pos())
		endLine, _ := GetLineAndCharacterOfPosition(sourceFile, body.End())
		if startLine < endLine {
			// The arrow function spans multiple lines,
			// make the error span be the first line, inclusive.
			return core.NewTextRange(pos, GetEndLinePosition(sourceFile, startLine))
		}
	}
	return core.NewTextRange(pos, node.End())
}

func GetErrorRangeForNode(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	errorNode := node
	switch node.Kind {
	case ast.KindSourceFile:
		pos := SkipTrivia(sourceFile.Text(), 0)
		if pos == len(sourceFile.Text()) {
			return core.NewTextRange(0, 0)
		}
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	// This list is a work in progress. Add missing node kinds to improve their error spans
	case ast.KindFunctionDeclaration, ast.KindMethodDeclaration:
		if node.Flags&ast.NodeFlagsReparsed != 0 {
			errorNode = node
			break
		}
		fallthrough
	case ast.KindVariableDeclaration, ast.KindBindingElement, ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration,
		ast.KindModuleDeclaration, ast.KindEnumDeclaration, ast.KindEnumMember, ast.KindFunctionExpression,
		ast.KindGetAccessor, ast.KindSetAccessor, ast.KindTypeAliasDeclaration, ast.KindJSTypeAliasDeclaration, ast.KindPropertyDeclaration,
		ast.KindPropertySignature, ast.KindNamespaceImport:
		errorNode = ast.GetNameOfDeclaration(node)
	case ast.KindArrowFunction:
		return getErrorRangeForArrowFunction(sourceFile, node)
	case ast.KindCaseClause, ast.KindDefaultClause:
		start := SkipTrivia(sourceFile.Text(), node.Pos())
		end := node.End()
		statements := node.AsCaseOrDefaultClause().Statements.Nodes
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return core.NewTextRange(start, end)
	case ast.KindReturnStatement, ast.KindYieldExpression:
		pos := SkipTrivia(sourceFile.Text(), node.Pos())
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindSatisfiesExpression:
		pos := SkipTrivia(sourceFile.Text(), node.AsSatisfiesExpression().Expression.End())
		return GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindConstructor:
		if node.Flags&ast.NodeFlagsReparsed != 0 {
			errorNode = node
			break
		}
		scanner := GetScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.TokenStart()
		for scanner.Token() != ast.KindConstructorKeyword && scanner.Token() != ast.KindStringLiteral && scanner.Token() != ast.KindEndOfFile {
			scanner.Scan()
		}
		return core.NewTextRange(start, scanner.TokenEnd())
		// !!!
		// case KindJSDocSatisfiesTag:
		// 	pos := scanner.SkipTrivia(sourceFile.Text(), node.tagName.pos)
		// 	return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	}
	if errorNode == nil {
		// If we don't have a better node, then just set the error on the first token of
		// construct.
		return GetRangeOfTokenAtPosition(sourceFile, node.Pos())
	}
	pos := errorNode.Pos()
	if !ast.NodeIsMissing(errorNode) && !ast.IsJsxText(errorNode) {
		pos = SkipTrivia(sourceFile.Text(), pos)
	}
	return core.NewTextRange(pos, errorNode.End())
}
