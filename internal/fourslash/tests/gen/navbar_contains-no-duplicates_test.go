package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavbar_contains_no_duplicates(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module Windows {
    export module Foundation {
        export var A;
        export class Test {
            public wow();
        }
    }
}

declare module Windows {
    export module Foundation {
        export var B;
        export module Test {
            export function Boom(): number;
        }
    }
}

class ABC {
    public foo() {
        return 3;
    }
}

module ABC {
    export var x = 3;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "ABC",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "foo",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
		{
			Name: "ABC",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
		{
			Name: "Windows",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "Foundation",
					Kind: lsproto.SymbolKindNamespace,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "A",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
						{
							Name:     "B",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
						{
							Name: "Test",
							Kind: lsproto.SymbolKindClass,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "wow",
									Kind:     lsproto.SymbolKindMethod,
									Children: nil,
								},
							}),
						},
						{
							Name: "Test",
							Kind: lsproto.SymbolKindNamespace,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "Boom",
									Kind:     lsproto.SymbolKindFunction,
									Children: nil,
								},
							}),
						},
					}),
				},
			}),
		},
	})
}
