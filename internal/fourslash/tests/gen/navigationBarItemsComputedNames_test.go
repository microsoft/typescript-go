package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsComputedNames(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const enum E {
	A = 'A',
}
const a = '';

class C {
    [a]() {
        return 1;
    }

    [E.A]() {
        return 1;
    }

    [1]() {
        return 1;
    },

    ["foo"]() {
        return 1;
    },
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "a",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name: "C",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "[a]",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "[E.A]",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "[1]",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "[\"foo\"]",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
		{
			Name: "E",
			Kind: lsproto.SymbolKindEnum,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "A",
					Kind:     lsproto.SymbolKindEnumMember,
					Children: nil,
				},
			}),
		},
	})
}
