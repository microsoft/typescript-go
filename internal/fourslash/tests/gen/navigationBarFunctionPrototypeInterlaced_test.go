package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionPrototypeInterlaced(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
var b = 1;
function A() {}; 
A.prototype.a = function() { };
A.b = function() { };
b = 2
/* Comment */
A.prototype.c = function() { }
var b = 2
A.prototype.d = function() { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "A",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "constructor",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "c",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "d",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
		{
			Name:     "b",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "b",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
