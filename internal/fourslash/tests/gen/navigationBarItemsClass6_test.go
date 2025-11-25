package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsClass6(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function Z() { }

Z.foo = 42

class Z { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "Z",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "constructor",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
				{
					Name:     "foo",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
		{
			Name:     "Z",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
}
