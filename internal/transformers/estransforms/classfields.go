package estransforms

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type classFieldsTransformer struct {
	transformers.Transformer
}

func (ch *classFieldsTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newClassFieldsTransformer(ctx context.Context) *transformers.Transformer {
	tx := &classFieldsTransformer{}
	return tx.NewTransformer(tx.visit, transformers.GetEmitContextFromContext(ctx))
}
