package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarInitializerSpans(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// get the name for the navbar from the variable name rather than the function name
const [|[|x|] = () => { var [|a|]; }|];
const [|[|f|] = function f() { var [|b|]; }|];
const [|[|y|] = { [|[|z|]: function z() { var [|c|]; }|] }|];`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "f",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "x",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "y",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "z",
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
	})
}
