package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionLikePropertyAssignments(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var functions = {
    a: 0,
    b: function () { },
    c: function x() { },
    d: () => { },
    e: y(),
    f() { }
};`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "functions",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "c",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "d",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "e",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "f",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
	})
}
