package ast_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"gotest.tools/v3/assert"
)

func TestIdAllocatorPagesAreExclusive(t *testing.T) {
	t.Parallel()

	var allocators [4]ast.IdAllocator
	pages := make(map[ast.SymbolId]int)
	for i := range 100000 {
		a := i % len(allocators)
		symbol := &ast.Symbol{}
		id := allocators[a].GetSymbolId(symbol)
		assert.Equal(t, ast.GetSymbolId(symbol), id)
		page := id / core.LinkStorePageSize
		owner, seen := pages[page]
		if seen {
			assert.Equal(t, owner, a, "page %d holds ids from two allocators", page)
		} else {
			pages[page] = a
		}
	}
}
