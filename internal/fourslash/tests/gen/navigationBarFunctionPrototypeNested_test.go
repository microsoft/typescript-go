package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionPrototypeNested(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
function A() {}
A.B = function () {  } 
A.B.prototype.d = function () {  }  
Object.defineProperty(A.B.prototype, "x", {
    get() {}
})
A.prototype.D = function () {  } 
A.prototype.D.prototype.d = function () {  } `
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
					Name: "B",
					Kind: lsproto.SymbolKindClass,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "constructor",
							Kind:     lsproto.SymbolKindConstructor,
							Children: nil,
						},
						{
							Name:     "d",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
						{
							Name: "x",
							Kind: lsproto.SymbolKindProperty,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "get",
									Kind:     lsproto.SymbolKindMethod,
									Children: nil,
								},
							}),
						},
					}),
				},
				{
					Name: "D",
					Kind: lsproto.SymbolKindClass,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "constructor",
							Kind:     lsproto.SymbolKindConstructor,
							Children: nil,
						},
						{
							Name:     "d",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
