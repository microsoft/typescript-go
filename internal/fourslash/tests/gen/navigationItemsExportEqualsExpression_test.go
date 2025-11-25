package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationItemsExportEqualsExpression(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export = function () {}
export = function () {
    return class Foo {
    }
}

export = () => ""
export = () => {
    return class Foo {
    }
}

export = function f1() {}
export = function f2() {
    return class Foo {
    }
}

const abc = 12;
export = abc;
export = class AB {}
export = {
    a: 1,
    b: 1,
    c: {
        d: 1
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "export=",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "export=",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "Foo",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
			}),
		},
		{
			Name:     "export=",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "export=",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "Foo",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
			}),
		},
		{
			Name:     "export=",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "export=",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "Foo",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
			}),
		},
		{
			Name:     "export=",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name: "export=",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name: "c",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "d",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
					}),
				},
			}),
		},
		{
			Name:     "abc",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "export=",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
