package ls

// movetofile.go implements the core types and logic for the "Move to file"
// and "Move to a new file" refactorings.
//
// !!! TODO: This is a stub implementation. The full implementation requires
// complex interactions with the type checker, module resolution, and import
// management that have not yet been fully ported.

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

// StatementRange represents a contiguous range of statements in a source file.
type StatementRange struct {
	First    *ast.Node // The first statement in the range
	AfterLast *ast.Node // The statement after the last in the range (may be nil)
}

// ToMove holds the set of statements to move and their ranges in the original file.
type ToMove struct {
	// All statements that will be moved.
	All []*ast.Node
	// Ranges of those statements in the original file.
	Ranges []StatementRange
}

// getStatementsToMove determines which statements in the given span should be
// moved, filtering out imports and prologue directives.
//
// !!! TODO: Implement full logic from TypeScript's getStatementsToMove.
func getStatementsToMove(file *ast.SourceFile, span core.TextRange) *ToMove {
	// !!! TODO: Implement
	// This should:
	// 1. Find the statements that overlap with span
	// 2. Filter out pure imports and prologue directives
	// 3. Handle overload ranges
	// 4. Return a ToMove with All and Ranges populated
	return nil
}

// isAllowedStatementToMove returns true if a statement can be moved to a new
// file. Pure imports and prologue directives are excluded.
//
// !!! TODO: Implement
func isAllowedStatementToMove(statement *ast.Node) bool {
	return !isPureImport(statement) && !ast.IsPrologueDirective(statement)
}

// isPureImport returns true if a node is a pure import statement (not exported).
//
// !!! TODO: Implement fully
func isPureImport(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindImportDeclaration:
		return true
	case ast.KindImportEqualsDeclaration:
		return !ast.HasSyntacticModifier(node, ast.ModifierFlagsExport)
	case ast.KindVariableStatement:
		// !!! TODO: check for require() calls
		return false
	default:
		return false
	}
}
