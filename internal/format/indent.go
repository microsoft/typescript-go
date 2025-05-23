package format

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func GetIndentationForNode(n *ast.Node, ignoreActualIndentationRange core.TextRange, sourceFile *ast.SourceFile, options *EditorSettings) int {
	return 0 // !!!
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

	return false // !!!
}

func NodeWillIndentChild(settings *FormatCodeSettings, parent *ast.Node, child *ast.Node, sourceFile *ast.SourceFile, indentByDefault bool) bool {
	return false // !!!
}
