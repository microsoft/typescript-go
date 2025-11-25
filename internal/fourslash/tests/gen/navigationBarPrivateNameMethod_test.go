package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarPrivateNameMethod(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
  #foo() {
    class B {
      #bar() {
         function baz () {
         }
      }
    }
  }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "A",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "#foo",
					Kind: lsproto.SymbolKindMethod,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name: "B",
							Kind: lsproto.SymbolKindClass,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name: "#bar",
									Kind: lsproto.SymbolKindMethod,
									Children: PtrTo([]*lsproto.DocumentSymbol{
										{
											Name:     "baz",
											Kind:     lsproto.SymbolKindFunction,
											Children: nil,
										},
									}),
								},
							}),
						},
					}),
				},
			}),
		},
	})
}
