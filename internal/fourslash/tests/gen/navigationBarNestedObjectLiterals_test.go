package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarNestedObjectLiterals(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var a = {
    b: 0,
    c: {},
    d: {
        e: 1,
    },
    f: {
        g: 2,
        h: {
            i: 3,
        },
    },
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "a",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
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
				{
					Name: "d",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "e",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
					}),
				},
				{
					Name: "f",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "g",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name: "h",
							Kind: lsproto.SymbolKindProperty,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "i",
									Kind:     lsproto.SymbolKindProperty,
									Children: nil,
								},
							}),
						},
					}),
				},
			}),
		},
	})
}
