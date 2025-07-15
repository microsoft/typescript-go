package estransforms

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type forawaitTransformer struct {
	transformers.Transformer
}

func (ch *forawaitTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newforawaitTransformer(ctx context.Context) *transformers.Transformer {
	tx := &forawaitTransformer{}
	return tx.NewTransformer(tx.visit, transformers.GetEmitContextFromContext(ctx))
}
