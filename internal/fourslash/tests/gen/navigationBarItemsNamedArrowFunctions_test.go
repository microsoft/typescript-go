package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsNamedArrowFunctions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export const value = 2;
export const func = () => 2;
export const func2 = function() { };
export function exportedFunction() { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "exportedFunction",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name:     "func",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "func2",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "value",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
