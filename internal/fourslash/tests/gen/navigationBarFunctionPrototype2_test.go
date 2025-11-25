package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarFunctionPrototype2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
A.prototype.a = function() { };
A.prototype.b = function() { };
function A() {}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "A",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "b",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "constructor",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
			}),
		},
	})
}
