package estransforms

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/transformers"
)

var useStrictPool = sync.Pool{New: func() any { return &useStrictTransformer{} }}

func getUseStrictTransformer() *useStrictTransformer {
	return useStrictPool.Get().(*useStrictTransformer)
}

func putUseStrictTransformer(tx *useStrictTransformer) {
	dispose, visit := tx.SaveState()
	*tx = useStrictTransformer{}
	tx.RestoreState(dispose, visit)
	useStrictPool.Put(tx)
}

func NewUseStrictTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := getUseStrictTransformer()
	tx.compilerOptions = opts.CompilerOptions
	tx.getEmitModuleFormatOfFile = opts.GetEmitModuleFormatOfFile
	if tx.GetDispose() == nil {
		tx.SetDispose(func() { putUseStrictTransformer(tx) })
	}
	return tx.NewTransformer(tx.visit, opts.Context)
}

type useStrictTransformer struct {
	transformers.Transformer
	compilerOptions           *core.CompilerOptions
	getEmitModuleFormatOfFile func(file ast.HasFileName) core.ModuleKind
}

func (tx *useStrictTransformer) visit(node *ast.Node) *ast.Node {
	if node.Kind != ast.KindSourceFile {
		return node
	}
	return tx.visitSourceFile(node.AsSourceFile())
}

func (tx *useStrictTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.ScriptKind == core.ScriptKindJSON {
		return node.AsNode()
	}

	isExternalModule := ast.IsExternalModule(node)
	moduleKind := tx.compilerOptions.GetEmitModuleKind()
	format := tx.getEmitModuleFormatOfFile(node)

	// ESM is always strict. If the file is ESM, and CJS emit
	// has not been requested, then skip adding "use strict".
	if isExternalModule && moduleKind >= core.ModuleKindES2015 &&
		(moduleKind == core.ModuleKindPreserve || format >= core.ModuleKindES2015) {
		return node.AsNode()
	}

	statements := tx.Factory().EnsureUseStrict(node.Statements.Nodes)
	statementList := tx.Factory().NewNodeList(statements)
	statementList.Loc = node.Statements.Loc
	return tx.Factory().UpdateSourceFile(node, statementList, node.EndOfFileToken).AsSourceFile().AsNode()
}
