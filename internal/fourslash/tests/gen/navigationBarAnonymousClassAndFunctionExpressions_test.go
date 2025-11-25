package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarAnonymousClassAndFunctionExpressions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `global.cls = class { };
(function() {
    const x = () => {
        // Presence of inner function causes x to be a top-level function.
        function xx() {}
    };
    const y = {
        // This is not a top-level function (contains nothing, but shows up in childItems of its parent.)
        foo: function() {}
    };
    (function nest() {
        function moreNest() {}
    })();
})();
(function() { // Different anonymous functions are not merged
    // These will only show up as childItems.
    function z() {}
    console.log(function() {})
    describe("this", 'function', ` + "`" + `is a function` + "`" + `, ` + "`" + `with template literal ${"a"}` + "`" + `, () => {});
    [].map(() => {});
})
(function classes() {
    // Classes show up in top-level regardless of whether they have names or inner declarations.
    const cls2 = class { };
    console.log(class cls3 {});
    (class { });
})`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "<function>",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "nest",
					Kind: lsproto.SymbolKindFunction,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "moreNest",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
					}),
				},
				{
					Name: "x",
					Kind: lsproto.SymbolKindVariable,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "xx",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
					}),
				},
				{
					Name: "y",
					Kind: lsproto.SymbolKindVariable,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "foo",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
			}),
		},
		{
			Name: "<function>",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "console.log() callback",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "describe(\"this\", 'function', `is a function`, `with template literal ${\"a\"}`) callback",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "map() callback",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "z",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
		{
			Name: "classes",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "<class>",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
				{
					Name:     "cls2",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
				{
					Name:     "cls3",
					Kind:     lsproto.SymbolKindClass,
					Children: nil,
				},
			}),
		},
		{
			Name:     "cls",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
}
