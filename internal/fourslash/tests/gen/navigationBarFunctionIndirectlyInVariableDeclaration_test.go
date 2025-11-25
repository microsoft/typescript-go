package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionIndirectlyInVariableDeclaration(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var a = {
    propA: function() {
        var c;
    }
};
var b;
b = {
    propB: function() {
    // function must not have an empty body to appear top level
        var d;
    }
};`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "a",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "propA",
					Kind: lsproto.SymbolKindMethod,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "c",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
					}),
				},
			}),
		},
		{
			Name:     "b",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name: "propB",
			Kind: lsproto.SymbolKindMethod,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "d",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
}
