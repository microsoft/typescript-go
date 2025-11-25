package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarMerging_grandchildren(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// Should not merge grandchildren with property assignments
const o = {
    a: {
        m() {},
    },
    b: {
        m() {},
    },
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "o",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "a",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "m",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
				{
					Name: "b",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "m",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
