package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationItemsExportDefaultExpression(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export default function () {}
export default function () {
    return class Foo {
    }
}

export default () => ""
export default () => {
    return class Foo {
    }
}

export default function f1() {}
export default function f2() {
    return class Foo {
    }
}

const abc = 12;
export default abc;
export default class AB {}
export default {
    a: 1,
    b: 1,
    c: {
        d: 1
    }
}

function foo(props: { x: number; y: number }) {}
export default foo({ x: 1, y: 1 });`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "default",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "default",
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
			Name:     "default",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "default",
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
			Name: "default",
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
			Name: "default",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "y",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
		{
			Name:     "AB",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name:     "abc",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "default",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "f1",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name: "f2",
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
			Name:     "foo",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
	})
}
