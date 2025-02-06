package transformers

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

type Transformer struct {
	ast.NodeVisitor
	EmitContext *printer.EmitContext
}

func (tx *Transformer) newTransformer(visit func(node *ast.Node) *ast.Node, emitContext *printer.EmitContext) *Transformer {
	if tx.Visit != nil || tx.EmitContext != nil || tx.Factory != nil {
		panic("Transformer already initialized")
	}
	if emitContext == nil {
		emitContext = printer.NewEmitContext()
	}
	tx.Visit = visit
	tx.EmitContext = emitContext
	tx.Factory = emitContext.Factory
	tx.Hooks.SetOriginal = emitContext.SetOriginal
	return tx
}
