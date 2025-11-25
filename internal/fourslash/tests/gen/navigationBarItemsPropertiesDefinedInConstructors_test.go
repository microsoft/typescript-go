package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsPropertiesDefinedInConstructors(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class List<T> {
    constructor(public a: boolean, private b: T, readonly c: string, d: number) {
        var local = 0;
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "List",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "constructor",
					Kind: lsproto.SymbolKindConstructor,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "local",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
					}),
				},
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "c",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
	})
}
