package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsModules2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `namespace Test.A { }

namespace Test.B {
    class Foo { }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "Test.A",
			Kind:     lsproto.SymbolKindModule,
			Children: nil,
		},
		{
			Name: "Test.B",
			Kind: lsproto.SymbolKindModule,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "Foo",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
			}),
		},
	})
}
