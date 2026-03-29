package transformers

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
)

func ExtractModifiers(emitContext *printer.EmitContext, modifiers *ast.ModifierList, allowed ast.ModifierFlags) *ast.ModifierList {
	if modifiers == nil {
		return nil
	}

	filtered := core.Filter(modifiers.Nodes, func(node *ast.Node) bool {
		flags := ast.ModifierToFlag(node.Kind)
		return flags == ast.ModifierFlagsNone || flags&allowed != 0
	})

	if core.Same(filtered, modifiers.Nodes) {
		return modifiers
	}

	list := emitContext.Factory.NewModifierList(filtered)
	list.Loc = modifiers.Loc
	return list
}
