package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsSymbols1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class C {
    [Symbol.isRegExp] = 0;
    [Symbol.iterator]() { }
    get [Symbol.isConcatSpreadable]() { }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "C",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "[Symbol.isRegExp]",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "[Symbol.iterator]",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "[Symbol.isConcatSpreadable]",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
	})
}
