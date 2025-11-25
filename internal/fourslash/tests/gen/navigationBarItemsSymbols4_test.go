package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsSymbols4(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @checkJs: true
// @allowJs: true
// @target: es6
// @Filename: file.js
const _sym = Symbol("_sym");
class MyClass {
    constructor() {
        // Dynamic assignment properties can't show up in navigation,
        // as they're not syntactic members
        // Additonally, late bound members are always filtered out, besides
        this[_sym] = "ok";
    }

    method() {
        this[_sym] = "yep";
        const x = this[_sym];
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "_sym",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name: "MyClass",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "constructor",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
				{
					Name: "method",
					Kind: lsproto.SymbolKindMethod,
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
	})
}
