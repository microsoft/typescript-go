package binder

import (
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func setParent(child *ast.Node, parent *ast.Node) {
	if child != nil {
		child.Parent = parent
	}
}

func SetParentInChildren(node *ast.Node) {
	node.ForEachChild(func(child *ast.Node) bool {
		child.Parent = node
		SetParentInChildren(child)
		return false
	})
}

func isSignedNumericLiteral(node *ast.Node) bool {
	if node.Kind == ast.KindPrefixUnaryExpression {
		node := node.AsPrefixUnaryExpression()
		return (node.Operator == ast.KindPlusToken || node.Operator == ast.KindMinusToken) && ast.IsNumericLiteral(node.Operand)
	}
	return false
}

// Atomic ids

var nextNodeId atomic.Uint32
var nextSymbolId atomic.Uint32

func GetNodeId(node *ast.Node) ast.NodeId {
	if node.Id == 0 {
		node.Id = ast.NodeId(nextNodeId.Add(1))
	}
	return node.Id
}

func GetSymbolId(symbol *ast.Symbol) ast.SymbolId {
	if symbol.Id == 0 {
		symbol.Id = ast.SymbolId(nextSymbolId.Add(1))
	}
	return symbol.Id
}

func GetSymbolTable(data *ast.SymbolTable) ast.SymbolTable {
	if *data == nil {
		*data = make(ast.SymbolTable)
	}
	return *data
}

func GetMembers(symbol *ast.Symbol) ast.SymbolTable {
	return GetSymbolTable(&symbol.Members)
}

func GetExports(symbol *ast.Symbol) ast.SymbolTable {
	return GetSymbolTable(&symbol.Exports)
}

func GetLocals(container *ast.Node) ast.SymbolTable {
	return GetSymbolTable(&container.LocalsContainerData().Locals)
}

func isFunctionPropertyAssignment(node *ast.Node) bool {
	if node.Kind == ast.KindBinaryExpression {
		expr := node.AsBinaryExpression()
		if expr.OperatorToken.Kind == ast.KindEqualsToken {
			switch expr.Left.Kind {
			case ast.KindPropertyAccessExpression:
				// F.id = expr
				return ast.IsIdentifier(expr.Left.Expression()) && ast.IsIdentifier(expr.Left.Name())
			case ast.KindElementAccessExpression:
				// F[xxx] = expr
				return ast.IsIdentifier(expr.Left.Expression())
			}
		}
	}
	return false
}

func getPostfixTokenFromNode(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindPropertyDeclaration:
		return node.AsPropertyDeclaration().PostfixToken
	case ast.KindPropertySignature:
		return node.AsPropertySignatureDeclaration().PostfixToken
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration().PostfixToken
	case ast.KindMethodSignature:
		return node.AsMethodSignatureDeclaration().PostfixToken
	}
	panic("Unhandled case in getPostfixTokenFromNode")
}

func ExportAssignmentIsAlias(node *ast.Node) bool {
	return isAliasableExpression(getExportAssignmentExpression(node))
}

func getExportAssignmentExpression(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindExportAssignment:
		return node.AsExportAssignment().Expression
	case ast.KindBinaryExpression:
		return node.AsBinaryExpression().Right
	}
	panic("Unhandled case in getExportAssignmentExpression")
}

func isAliasableExpression(e *ast.Node) bool {
	return ast.IsEntityNameExpression(e) || ast.IsClassExpression(e)
}

func nodeHasName(statement *ast.Node, id *ast.Node) bool {
	name := statement.Name()
	if name != nil {
		return ast.IsIdentifier(name) && name.AsIdentifier().Text == id.AsIdentifier().Text
	}
	if ast.IsVariableStatement(statement) {
		declarations := statement.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes
		return core.Some(declarations, func(d *ast.Node) bool { return nodeHasName(d, id) })
	}
	return false
}

func isAsyncFunction(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction, ast.KindMethodDeclaration:
		data := node.BodyData()
		return data.Body != nil && data.AsteriskToken == nil && ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync)
	}
	return false
}

func SymbolName(symbol *ast.Symbol) string {
	if symbol.ValueDeclaration != nil && ast.IsPrivateIdentifierClassElementDeclaration(symbol.ValueDeclaration) {
		return symbol.ValueDeclaration.Name().AsPrivateIdentifier().Text
	}
	return symbol.Name
}

func isFunctionSymbol(symbol *ast.Symbol) bool {
	d := symbol.ValueDeclaration
	if d != nil {
		if ast.IsFunctionDeclaration(d) {
			return true
		}
		if ast.IsVariableDeclaration(d) {
			varDecl := d.AsVariableDeclaration()
			if varDecl.Initializer != nil {
				return ast.IsFunctionLike(varDecl.Initializer)
			}
		}
	}
	return false
}

func unreachableCodeIsError(options *core.CompilerOptions) bool {
	return options.AllowUnreachableCode == core.TSFalse
}

func unusedLabelIsError(options *core.CompilerOptions) bool {
	return options.AllowUnusedLabels == core.TSFalse
}

func isStatementCondition(node *ast.Node) bool {
	switch node.Parent.Kind {
	case ast.KindIfStatement:
		return node.Parent.AsIfStatement().Expression == node
	case ast.KindWhileStatement:
		return node.Parent.AsWhileStatement().Expression == node
	case ast.KindDoStatement:
		return node.Parent.AsDoStatement().Expression == node
	case ast.KindForStatement:
		return node.Parent.AsForStatement().Condition == node
	case ast.KindConditionalExpression:
		return node.Parent.AsConditionalExpression().Condition == node
	}
	return false
}

func isTopLevelLogicalExpression(node *ast.Node) bool {
	for ast.IsParenthesizedExpression(node.Parent) || ast.IsPrefixUnaryExpression(node.Parent) && node.Parent.AsPrefixUnaryExpression().Operator == ast.KindExclamationToken {
		node = node.Parent
	}
	return !isStatementCondition(node) && !ast.IsLogicalExpression(node.Parent) && !(ast.IsOptionalChain(node.Parent) && node.Parent.Expression() == node)
}

func isAssignmentDeclaration(decl *ast.Node) bool {
	return ast.IsBinaryExpression(decl) || ast.IsAccessExpression(decl) || ast.IsIdentifier(decl) || ast.IsCallExpression(decl)
}

func isEffectiveModuleDeclaration(node *ast.Node) bool {
	return ast.IsModuleDeclaration(node) || ast.IsIdentifier(node)
}

func getErrorRangeForArrowFunction(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	pos := scanner.SkipTrivia(sourceFile.Text, node.Pos())
	body := node.AsArrowFunction().Body
	if body != nil && body.Kind == ast.KindBlock {
		startLine, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, body.Pos())
		endLine, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, body.End())
		if startLine < endLine {
			// The arrow function spans multiple lines,
			// make the error span be the first line, inclusive.
			return core.NewTextRange(pos, scanner.GetEndLinePosition(sourceFile, startLine))
		}
	}
	return core.NewTextRange(pos, node.End())
}

func GetErrorRangeForNode(sourceFile *ast.SourceFile, node *ast.Node) core.TextRange {
	errorNode := node
	switch node.Kind {
	case ast.KindSourceFile:
		pos := scanner.SkipTrivia(sourceFile.Text, 0)
		if pos == len(sourceFile.Text) {
			return core.NewTextRange(0, 0)
		}
		return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	// This list is a work in progress. Add missing node kinds to improve their error spans
	case ast.KindVariableDeclaration, ast.KindBindingElement, ast.KindClassDeclaration, ast.KindClassExpression, ast.KindInterfaceDeclaration,
		ast.KindModuleDeclaration, ast.KindEnumDeclaration, ast.KindEnumMember, ast.KindFunctionDeclaration, ast.KindFunctionExpression,
		ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor, ast.KindTypeAliasDeclaration, ast.KindPropertyDeclaration,
		ast.KindPropertySignature, ast.KindNamespaceImport:
		errorNode = ast.GetNameOfDeclaration(node)
	case ast.KindArrowFunction:
		return getErrorRangeForArrowFunction(sourceFile, node)
	case ast.KindCaseClause, ast.KindDefaultClause:
		start := scanner.SkipTrivia(sourceFile.Text, node.Pos())
		end := node.End()
		statements := node.AsCaseOrDefaultClause().Statements.Nodes
		if len(statements) != 0 {
			end = statements[0].Pos()
		}
		return core.NewTextRange(start, end)
	case ast.KindReturnStatement, ast.KindYieldExpression:
		pos := scanner.SkipTrivia(sourceFile.Text, node.Pos())
		return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindSatisfiesExpression:
		pos := scanner.SkipTrivia(sourceFile.Text, node.AsSatisfiesExpression().Expression.End())
		return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	case ast.KindConstructor:
		scanner := scanner.GetScannerForSourceFile(sourceFile, node.Pos())
		start := scanner.TokenStart()
		for scanner.Token() != ast.KindConstructorKeyword && scanner.Token() != ast.KindStringLiteral && scanner.Token() != ast.KindEndOfFile {
			scanner.Scan()
		}
		return core.NewTextRange(start, scanner.TokenEnd())
		// !!!
		// case KindJSDocSatisfiesTag:
		// 	pos := scanner.SkipTrivia(sourceFile.text, node.tagName.pos)
		// 	return scanner.GetRangeOfTokenAtPosition(sourceFile, pos)
	}
	if errorNode == nil {
		// If we don't have a better node, then just set the error on the first token of
		// construct.
		return scanner.GetRangeOfTokenAtPosition(sourceFile, node.Pos())
	}
	pos := errorNode.Pos()
	if !ast.NodeIsMissing(errorNode) {
		pos = scanner.SkipTrivia(sourceFile.Text, pos)
	}
	return core.NewTextRange(pos, errorNode.End())
}
