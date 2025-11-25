package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarPropertyDeclarations(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
    public A1 = class {
        public x = 1;
        private y() {}
        protected z() {}
    }

    public A2 = {
        x: 1,
        y() {},
        z() {}
    }

    public A3 = function () {}
    public A4 = () => {}
    public A5 = 1;
    public A6 = "A6";

    public ["A7"] = class {
        public x = 1;
        private y() {}
        protected z() {}
    }

    public [1] = {
        x: 1,
        y() {},
        z() {}
    }

    public [1 + 1] = 1;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "A",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "[1]",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "x",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "z",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
				{
					Name: "A1",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "x",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "z",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
				{
					Name: "A2",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "x",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "z",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
				{
					Name:     "A3",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "A4",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "A5",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "A6",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name: "[\"A7\"]",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "x",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "z",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
