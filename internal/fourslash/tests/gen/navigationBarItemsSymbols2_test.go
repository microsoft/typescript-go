package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsSymbols2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface I {
    [Symbol.isRegExp]: string;
    [Symbol.iterator](): string;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "I",
			Kind: lsproto.SymbolKindInterface,
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
			}),
		},
	})
}
