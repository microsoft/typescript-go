package transformers

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

type Transformer struct {
	emitContext *printer.EmitContext
	factory     *printer.NodeFactory
	visitor     ast.NodeVisitor
	dispose     func()
}

func (tx *Transformer) NewTransformer(visit func(node *ast.Node) *ast.Node, emitContext *printer.EmitContext) *Transformer {
	if emitContext == nil {
		emitContext = printer.NewEmitContext()
	}
	tx.emitContext = emitContext
	tx.factory = emitContext.Factory
	if tx.visitor.Visit == nil {
		emitContext.InitNodeVisitor(&tx.visitor, visit)
	} else {
		emitContext.InitNodeVisitor(&tx.visitor, nil)
	}
	return tx
}

func (tx *Transformer) Dispose() {
	if tx.dispose != nil {
		tx.dispose()
	}
}

// SaveState returns the reusable state (dispose callback and visitor Visit function)
// that should be preserved across a pool reset.
func (tx *Transformer) SaveState() (func(), func(node *ast.Node) *ast.Node) {
	return tx.dispose, tx.visitor.Visit
}

// RestoreState restores the reusable state after a pool reset.
func (tx *Transformer) RestoreState(dispose func(), visit func(node *ast.Node) *ast.Node) {
	tx.dispose = dispose
	tx.visitor.Visit = visit
}

func (tx *Transformer) SetDispose(fn func()) {
	tx.dispose = fn
}

func (tx *Transformer) GetDispose() func() {
	return tx.dispose
}

func (tx *Transformer) EmitContext() *printer.EmitContext {
	return tx.emitContext
}

func (tx *Transformer) Visitor() *ast.NodeVisitor {
	return &tx.visitor
}

func (tx *Transformer) Factory() *printer.NodeFactory {
	return tx.factory
}

func (tx *Transformer) TransformSourceFile(file *ast.SourceFile) *ast.SourceFile {
	return tx.visitor.VisitSourceFile(file)
}
