package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsFunctions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo() {
    var x = 10;
    function bar() {
        var y = 10;
        function biz() {
            var z = 10;
        }
        function qux() {
            // A function with an empty body should not be top level
        }
    }
}

function baz() {
    var v = 10;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "baz",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "v",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "foo",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "bar",
					Kind: lsproto.SymbolKindFunction,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name: "biz",
							Kind: lsproto.SymbolKindFunction,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "z",
									Kind:     lsproto.SymbolKindVariable,
									Children: nil,
								},
							}),
						},
						{
							Name:     "qux",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
					}),
				},
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
}
