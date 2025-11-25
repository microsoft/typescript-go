package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsModules1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module "X.Y.Z" {}

declare module 'X2.Y2.Z2' {}

declare module "foo";

module A.B.C {
    export var x;
}

module A.B {
    export var y;
}

module A {
    export var z;
}

module A {
    module B {
        module C {
            declare var x;
        }
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "'X2.Y2.Z2'",
			Kind:     lsproto.SymbolKindModule,
			Children: nil,
		},
		{
			Name:     "\"foo\"",
			Kind:     lsproto.SymbolKindModule,
			Children: nil,
		},
		{
			Name:     "\"X.Y.Z\"",
			Kind:     lsproto.SymbolKindModule,
			Children: nil,
		},
		{
			Name: "A",
			Kind: lsproto.SymbolKindModule,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "B",
					Kind: lsproto.SymbolKindModule,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name: "C",
							Kind: lsproto.SymbolKindModule,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "x",
									Kind:     lsproto.SymbolKindVariable,
									Children: nil,
								},
							}),
						},
					}),
				},
				{
					Name:     "z",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "A.B",
			Kind: lsproto.SymbolKindModule,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "y",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "A.B.C",
			Kind: lsproto.SymbolKindModule,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
}
