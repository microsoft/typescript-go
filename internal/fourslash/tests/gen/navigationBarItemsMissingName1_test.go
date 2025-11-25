package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsMissingName1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export function
class C {
    foo() {}
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "<function>",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "C",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "foo",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
	})
}
