package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarMerging(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: file1.ts
module a {
    function foo() {}
}
module b {
    function foo() {}
}
module a {
    function bar() {}
}
// @Filename: file2.ts
module a {}
function a() {}
// @Filename: file3.ts
module a {
    interface A {
        foo: number;
    }
}
module a {
    interface A {
        bar: number;
    }
}
// @Filename: file4.ts
module A { export var x; }
module A.B { export var y; }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "a",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "bar",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "foo",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
		{
			Name: "b",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "foo",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
	})
	f.GoToFile(t, "file2.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "a",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name:     "a",
			Kind:     lsproto.SymbolKindNamespace,
			Children: nil,
		},
	})
	f.GoToFile(t, "file3.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "a",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "A",
					Kind: lsproto.SymbolKindInterface,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "bar",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "foo",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
	f.GoToFile(t, "file4.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "A",
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
			Name: "A.B",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "y",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
}
