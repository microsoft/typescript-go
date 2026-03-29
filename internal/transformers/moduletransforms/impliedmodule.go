package moduletransforms

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type ImpliedModuleTransformer struct {
	transformers.Transformer
	opts                      *transformers.TransformOptions
	resolver                  binder.ReferenceResolver
	getEmitModuleFormatOfFile func(file ast.HasFileName) core.ModuleKind
	cjsTransformer            *transformers.Transformer
	esmTransformer            *transformers.Transformer
}

var impliedModulePool = sync.Pool{New: func() any { return &ImpliedModuleTransformer{} }}

func getImpliedModuleTransformer() *ImpliedModuleTransformer {
	return impliedModulePool.Get().(*ImpliedModuleTransformer)
}

func putImpliedModuleTransformer(tx *ImpliedModuleTransformer) {
	dispose, visit := tx.SaveState()
	*tx = ImpliedModuleTransformer{}
	tx.RestoreState(dispose, visit)
	impliedModulePool.Put(tx)
}

func NewImpliedModuleTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := getImpliedModuleTransformer()
	tx.opts = opts
	tx.resolver = opts.Resolver
	tx.getEmitModuleFormatOfFile = opts.GetEmitModuleFormatOfFile
	if tx.GetDispose() == nil {
		tx.SetDispose(func() { putImpliedModuleTransformer(tx) })
	}
	return tx.NewTransformer(tx.visit, opts.Context)
}

func (tx *ImpliedModuleTransformer) visit(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindSourceFile:
		node = tx.visitSourceFile(node.AsSourceFile())
	}
	return node
}

func (tx *ImpliedModuleTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.IsDeclarationFile {
		return node.AsNode()
	}

	format := tx.getEmitModuleFormatOfFile(node)

	var transformer *transformers.Transformer
	if format >= core.ModuleKindES2015 {
		if tx.esmTransformer == nil {
			tx.esmTransformer = NewESModuleTransformer(tx.opts)
		}
		transformer = tx.esmTransformer
	} else {
		if tx.cjsTransformer == nil {
			tx.cjsTransformer = NewCommonJSModuleTransformer(tx.opts)
		}
		transformer = tx.cjsTransformer
	}

	return transformer.TransformSourceFile(node).AsNode()
}
