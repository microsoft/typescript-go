package format

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func GetIndentationForNode(n *ast.Node, ignoreActualIndentationRange *core.TextRange, sourceFile *ast.SourceFile, options *FormatCodeSettings) int {
	startline, startpos := scanner.GetLineAndCharacterOfPosition(sourceFile, scanner.GetTokenPosOfNode(n, sourceFile, false))
	return getIndentationForNodeWorker(n, startline, startpos, ignoreActualIndentationRange /*indentationDelta*/, 0, sourceFile /*isNextChild*/, false, options)
}

func getIndentationForNodeWorker(
	current *ast.Node,
	currentStartLine int,
	currentStartCharacter int,
	ignoreActualIndentationRange *core.TextRange,
	indentationDelta int,
	sourceFile *ast.SourceFile,
	isNextChild bool,
	options *FormatCodeSettings,
) int {
	parent := current.Parent

	// Walk up the tree and collect indentation for parent-child node pairs. Indentation is not added if
	// * parent and child nodes start on the same line, or
	// * parent is an IfStatement and child starts on the same line as an 'else clause'.
	for parent != nil {
		useActualIndentation := true
		if ignoreActualIndentationRange != nil {
			start := scanner.GetTokenPosOfNode(current, sourceFile, false)
			useActualIndentation = start < ignoreActualIndentationRange.Pos() || start > ignoreActualIndentationRange.End()
		}

		containingListOrParentStartLine, containingListOrParentStartCharacter := getContainingListOrParentStart(parent, current, sourceFile)
		parentAndChildShareLine := containingListOrParentStartLine == currentStartLine ||
			childStartsOnTheSameLineWithElseInIfStatement(parent, current, currentStartLine, sourceFile)

		if useActualIndentation {
			// check if current node is a list item - if yes, take indentation from it
			var firstListChild *ast.Node
			containerList := getContainingList(current, sourceFile)
			if containerList != nil {
				firstListChild = core.FirstOrNil(containerList.Nodes)
			}
			// A list indents its children if the children begin on a later line than the list itself:
			//
			// f1(               L0 - List start
			//   {               L1 - First child start: indented, along with all other children
			//     prop: 0
			//   },
			//   {
			//     prop: 1
			//   }
			// )
			//
			// f2({             L0 - List start and first child start: children are not indented.
			//   prop: 0             Object properties are indented only one level, because the list
			// }, {                  itself contributes nothing.
			//   prop: 1        L3 - The indentation of the second object literal is best understood by
			// })                    looking at the relationship between the list and *first* list item.
			listIndentsChild := firstListChild != nil && getStartLineAndCharacterForNode(firstListChild, sourceFile).line > containingListOrParentStartLine
			actualIndentation := getActualIndentationForListItem(current, sourceFile, options, listIndentsChild)
			if actualIndentation != -1 {
				return actualIndentation + indentationDelta
			}

			// try to fetch actual indentation for current node from source text
			actualIndentation = getActualIndentationForNode(current, parent, currentStartLine, currentStartCharacter, parentAndChildShareLine, sourceFile, options)
			if actualIndentation != -1 {
				return actualIndentation + indentationDelta
			}
		}

		// increase indentation if parent node wants its content to be indented and parent and child nodes don't start on the same line
		if ShouldIndentChildNode(options, parent, current, sourceFile, isNextChild) && !parentAndChildShareLine {
			indentationDelta += options.IndentSize
		}

		// In our AST, a call argument's `parent` is the call-expression, not the argument list.
		// We would like to increase indentation based on the relationship between an argument and its argument-list,
		// so we spoof the starting position of the (parent) call-expression to match the (non-parent) argument-list.
		// But, the spoofed start-value could then cause a problem when comparing the start position of the call-expression
		// to *its* parent (in the case of an iife, an expression statement), adding an extra level of indentation.
		//
		// Instead, when at an argument, we unspoof the starting position of the enclosing call expression
		// *after* applying indentation for the argument.

		useTrueStart := isArgumentAndStartLineOverlapsExpressionBeingCalled(parent, current, currentStartLine, sourceFile)

		current = parent
		parent = current.Parent

		if useTrueStart {
			currentStartLine, currentStartCharacter = scanner.GetLineAndCharacterOfPosition(sourceFile, scanner.GetTokenPosOfNode(current, sourceFile, false))
		} else {
			currentStartLine = containingListOrParentStartLine
			currentStartCharacter = containingListOrParentStartCharacter
		}
	}

	return indentationDelta + options.BaseIndentSize
}

func getContainingList(node *ast.Node, sourceFile *ast.SourceFile) *ast.NodeList {
	if node.Parent == nil {
		return nil
	}
	return getListByRange(scanner.GetTokenPosOfNode(node, sourceFile, false), node.End(), node.Parent, sourceFile)
}

func getListByPosition(pos int, node *ast.Node, sourceFile *ast.SourceFile) *ast.NodeList {
	if node == nil {
		return nil
	}
	return getListByRange(pos, pos, node, sourceFile)
}

func getListByRange(start int, end int, node *ast.Node, sourceFile *ast.SourceFile) *ast.NodeList {
	// !!!
	return nil
}

func getVisualListRange(node *ast.Node, list core.TextRange, sourceFile *ast.SourceFile) core.TextRange {
	// !!! In strada, this relied on the services .getChildren method, which manifested synthetic token nodes
	// _however_, the logic boils down to "find the child with the matching span and adjust its start to the
	// previous (possibly token) child's end and its end to the token start of the following element" - basically
	// expanding the range to encompass all the neighboring non-token trivia
	// Now, we perform that logic with the scanner instead
	// !!!
	prior := astnav.FindPrecedingToken(sourceFile, list.Pos())
	return core.NewTextRange(-1, -1) // !!!
}

func getContainingListOrParentStart(parent *ast.Node, child *ast.Node, sourceFile *ast.SourceFile) (line int, character int) {
	containingList := getContainingList(child, sourceFile)
	var startPos int
	if containingList != nil {
		startPos = containingList.Loc.Pos()
	} else {
		startPos = scanner.GetTokenPosOfNode(parent, sourceFile, false)
	}
	return scanner.GetLineAndCharacterOfPosition(sourceFile, startPos)
}

func isControlFlowEndingStatement(kind ast.Kind, parentKind ast.Kind) bool {
	switch kind {
	case ast.KindReturnStatement, ast.KindThrowStatement, ast.KindContinueStatement, ast.KindBreakStatement:
		return parentKind != ast.KindBlock
	default:
		return false
	}
}

/**
* True when the parent node should indent the given child by an explicit rule.
* @param isNextChild If true, we are judging indent of a hypothetical child *after* this one, not the current child.
 */
func ShouldIndentChildNode(settings *FormatCodeSettings, parent *ast.Node, child *ast.Node, sourceFile *ast.SourceFile, isNextChildArg ...bool) bool {
	isNextChild := false
	if len(isNextChildArg) > 0 {
		isNextChild = isNextChildArg[0]
	}

	return NodeWillIndentChild(settings, parent, child, sourceFile, false) && !(isNextChild && child != nil && isControlFlowEndingStatement(child.Kind, parent.Kind))
}

func NodeWillIndentChild(settings *FormatCodeSettings, parent *ast.Node, child *ast.Node, sourceFile *ast.SourceFile, indentByDefault bool) bool {
	childKind := ast.KindUnknown
	if child != nil {
		childKind = child.Kind
	}

	switch parent.Kind {
	case ast.KindExpressionStatement,
		ast.KindClassDeclaration,
		ast.KindClassExpression,
		ast.KindInterfaceDeclaration,
		ast.KindEnumDeclaration,
		ast.KindTypeAliasDeclaration,
		ast.KindArrayLiteralExpression,
		ast.KindBlock,
		ast.KindModuleBlock,
		ast.KindObjectLiteralExpression,
		ast.KindTypeLiteral,
		ast.KindMappedType,
		ast.KindTupleType,
		ast.KindParenthesizedExpression,
		ast.KindPropertyAccessExpression,
		ast.KindCallExpression,
		ast.KindNewExpression,
		ast.KindVariableStatement,
		ast.KindExportAssignment,
		ast.KindReturnStatement,
		ast.KindConditionalExpression,
		ast.KindArrayBindingPattern,
		ast.KindObjectBindingPattern,
		ast.KindJsxOpeningElement,
		ast.KindJsxOpeningFragment,
		ast.KindJsxSelfClosingElement,
		ast.KindJsxExpression,
		ast.KindMethodSignature,
		ast.KindCallSignature,
		ast.KindConstructSignature,
		ast.KindParameter,
		ast.KindFunctionType,
		ast.KindConstructorType,
		ast.KindParenthesizedType,
		ast.KindTaggedTemplateExpression,
		ast.KindAwaitExpression,
		ast.KindNamedExports,
		ast.KindNamedImports,
		ast.KindExportSpecifier,
		ast.KindImportSpecifier,
		ast.KindPropertyDeclaration,
		ast.KindCaseClause,
		ast.KindDefaultClause:
		return true
	case ast.KindCaseBlock:
		return settings.IndentSwitchCase.IsTrueOrUnknown()
	case ast.KindVariableDeclaration, ast.KindPropertyAssignment, ast.KindBinaryExpression:
		if settings.IndentMultiLineObjectLiteralBeginningOnBlankLine.IsFalseOrUnknown() && sourceFile != nil && childKind == ast.KindObjectLiteralExpression {
			return rangeIsOnOneLine(child.Loc, sourceFile)
		}
		if parent.Kind == ast.KindBinaryExpression && sourceFile != nil && childKind == ast.KindJsxElement {
			parentStartLine, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, scanner.SkipTrivia(sourceFile.Text(), parent.Pos()))
			childStartLine, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, scanner.SkipTrivia(sourceFile.Text(), child.Pos()))
			return parentStartLine != childStartLine
		}
		if parent.Kind != ast.KindBinaryExpression {
			return true
		}
		return indentByDefault
	case ast.KindDoStatement,
		ast.KindWhileStatement,
		ast.KindForInStatement,
		ast.KindForOfStatement,
		ast.KindForStatement,
		ast.KindIfStatement,
		ast.KindFunctionDeclaration,
		ast.KindFunctionExpression,
		ast.KindMethodDeclaration,
		ast.KindConstructor,
		ast.KindGetAccessor,
		ast.KindSetAccessor:
		return childKind != ast.KindBlock
	case ast.KindArrowFunction:
		if sourceFile != nil && childKind == ast.KindParenthesizedExpression {
			return rangeIsOnOneLine(child.Loc, sourceFile)
		}
		return childKind != ast.KindBlock
	case ast.KindExportDeclaration:
		return childKind != ast.KindNamedExports
	case ast.KindImportDeclaration:
		return childKind != ast.KindImportClause || (child.AsImportClause().NamedBindings != nil && child.AsImportClause().NamedBindings.Kind != ast.KindNamedImports)
	case ast.KindJsxElement:
		return childKind != ast.KindJsxClosingElement
	case ast.KindJsxFragment:
		return childKind != ast.KindJsxClosingFragment
	case ast.KindIntersectionType, ast.KindUnionType, ast.KindSatisfiesExpression:
		if childKind == ast.KindTypeLiteral || childKind == ast.KindTupleType || childKind == ast.KindMappedType {
			return false
		}
		return indentByDefault
	case ast.KindTryStatement:
		if childKind == ast.KindBlock {
			return false
		}
		return indentByDefault
	}

	// No explicit rule for given nodes so the result will follow the default value argument
	return indentByDefault
}
