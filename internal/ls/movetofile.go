package ls

// movetofile.go implements the core types and logic for the "Move to file"
// and "Move to a new file" refactorings.
//
// This is a partial implementation. Complex parts requiring deep type checker
// integration (symbol analysis, module resolution) are marked with !!! TODO.

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/core"
)

// StatementRange represents a contiguous range of statements in a source file.
type StatementRange struct {
	First     *ast.Node // The first statement in the range
	AfterLast *ast.Node // The statement after the last in the range (may be nil)
}

// ToMove holds the set of statements to move and their ranges in the original file.
type ToMove struct {
	// All statements that will be moved.
	All []*ast.Node
	// Ranges of those statements in the original file.
	Ranges []StatementRange
}

// rangeToMove is an internal helper type for getRangeToMove.
type rangeToMove struct {
	toMove    []*ast.Node
	afterLast *ast.Node
}

// getRangeToMove finds the contiguous block of statements in the source file
// that overlap with the given span.
func getRangeToMove(file *ast.SourceFile, span core.TextRange) *rangeToMove {
	statements := file.Statements.Nodes

	// Find the first statement whose end is after span.Pos().
	startNodeIndex := -1
	for i, s := range statements {
		if s.End() > span.Pos() {
			startNodeIndex = i
			break
		}
	}
	if startNodeIndex == -1 {
		return nil
	}

	// !!! TODO: Handle overload ranges (needs type checker for symbol.declarations)

	// Find the last statement whose end is >= span.End(), starting from startNodeIndex.
	endNodeIndex := -1
	for i := startNodeIndex; i < len(statements); i++ {
		if statements[i].End() >= span.End() {
			endNodeIndex = i
			break
		}
	}

	// If the range ends before the start of the end statement, go back one.
	if endNodeIndex != -1 && span.End() <= astnav.GetStartOfNode(statements[endNodeIndex], file, false) {
		endNodeIndex--
	}

	// !!! TODO: Handle ending overload ranges (needs type checker)

	var toMove []*ast.Node
	var afterLast *ast.Node
	if endNodeIndex == -1 {
		toMove = statements[startNodeIndex:]
	} else {
		toMove = statements[startNodeIndex : endNodeIndex+1]
		if endNodeIndex+1 < len(statements) {
			afterLast = statements[endNodeIndex+1]
		}
	}

	return &rangeToMove{toMove: toMove, afterLast: afterLast}
}

// getRangesWhere calls cb for each contiguous range of elements in arr for which pred returns true.
// cb receives (startIndex, afterEndIndex).
func getRangesWhere(arr []*ast.Node, pred func(*ast.Node) bool, cb func(start, afterEnd int)) {
	start := -1
	for i, elem := range arr {
		if pred(elem) {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				cb(start, i)
				start = -1
			}
		}
	}
	if start != -1 {
		cb(start, len(arr))
	}
}

// getStatementsToMove determines which statements in the given span should be
// moved, filtering out imports and prologue directives.
func getStatementsToMove(file *ast.SourceFile, span core.TextRange) *ToMove {
	r := getRangeToMove(file, span)
	if r == nil {
		return nil
	}

	var all []*ast.Node
	var ranges []StatementRange

	getRangesWhere(r.toMove, isAllowedStatementToMove, func(start, afterEnd int) {
		for i := start; i < afterEnd; i++ {
			all = append(all, r.toMove[i])
		}
		// afterLast is the statement immediately following this range:
		// if the range ends before the end of toMove, it is toMove[afterEnd];
		// otherwise it is the statement after the entire toMove block.
		var afterLast *ast.Node
		if afterEnd < len(r.toMove) {
			afterLast = r.toMove[afterEnd]
		} else {
			afterLast = r.afterLast
		}
		ranges = append(ranges, StatementRange{
			First:     r.toMove[start],
			AfterLast: afterLast,
		})
	})

	if len(all) == 0 {
		return nil
	}
	return &ToMove{All: all, Ranges: ranges}
}

// isAllowedStatementToMove returns true if a statement can be moved to a new
// file. Pure imports and prologue directives are excluded.
func isAllowedStatementToMove(statement *ast.Node) bool {
	return !isPureImport(statement) && !ast.IsPrologueDirective(statement)
}

// isPureImport returns true if a node is a pure import statement (not exported).
func isPureImport(node *ast.Node) bool {
	if node == nil {
		return false
	}
	switch node.Kind {
	case ast.KindImportDeclaration:
		return true
	case ast.KindImportEqualsDeclaration:
		return !ast.HasSyntacticModifier(node, ast.ModifierFlagsExport)
	case ast.KindVariableStatement:
		// !!! TODO: check for require() calls using isVariableDeclarationInitializedToRequire
		return false
	default:
		return false
	}
}
