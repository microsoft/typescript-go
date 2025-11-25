package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionPrototype(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
function f() {}
f.prototype.x = 0;
f.y = 0;
f.prototype.method = function () {};
Object.defineProperty(f, 'staticProp', { 
    set: function() {}, 
    get: function(){
    } 
});
Object.defineProperty(f.prototype, 'name', { 
    set: function() {}, 
    get: function(){
    } 
}); `
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "f",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "constructor",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
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
				{
					Name:     "method",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name: "staticProp",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "get",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "set",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
				{
					Name: "name",
					Kind: lsproto.SymbolKindProperty,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "get",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "set",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
