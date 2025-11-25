package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarAnonymousClassAndFunctionExpressions2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `console.log(console.log(class Y {}, class X {}), console.log(class B {}, class A {}));
console.log(class Cls { meth() {} });`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "A",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name:     "B",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name: "Cls",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "meth",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
		{
			Name:     "X",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name:     "Y",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
}
