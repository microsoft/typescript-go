package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsExports(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export { a } from "a";

export { b as B } from "a" 

export import e = require("a");

export * from "a"; // no bindings here`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "a",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "B",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "e",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
