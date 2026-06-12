package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

// symbolLinkStore is a links store keyed by symbol id rather than by symbol pointer.
// Indexing by id avoids map overhead in the hottest link stores.
type symbolLinkStore[V any] struct {
	store core.IdLinkStore[V]
}

func (s *symbolLinkStore[V]) Get(symbol *ast.Symbol) *V {
	return s.store.Get(uint64(ast.GetSymbolId(symbol)))
}

func (s *symbolLinkStore[V]) Has(symbol *ast.Symbol) bool {
	return s.store.Has(uint64(ast.GetSymbolId(symbol)))
}

func (s *symbolLinkStore[V]) TryGet(symbol *ast.Symbol) *V {
	return s.store.TryGet(uint64(ast.GetSymbolId(symbol)))
}

// nodeLinkStore is a links store keyed by node id rather than by node pointer.
type nodeLinkStore[V any] struct {
	store core.IdLinkStore[V]
}

func (s *nodeLinkStore[V]) Get(node *ast.Node) *V {
	return s.store.Get(uint64(ast.GetNodeId(node)))
}

func (s *nodeLinkStore[V]) Has(node *ast.Node) bool {
	return s.store.Has(uint64(ast.GetNodeId(node)))
}

func (s *nodeLinkStore[V]) TryGet(node *ast.Node) *V {
	return s.store.TryGet(uint64(ast.GetNodeId(node)))
}
