package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionPrototypeBroken(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
function A() {}
A. // Started typing something here
A.prototype.a = function() { };
G. // Started typing something here
A.prototype.a = function() { };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "G",
			Kind: lsproto.SymbolKindMethod,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "A",
					Kind: lsproto.SymbolKindMethod,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "a",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
					}),
				},
			}),
		},
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
					Name: "A",
					Kind: lsproto.SymbolKindMethod,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "a",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
