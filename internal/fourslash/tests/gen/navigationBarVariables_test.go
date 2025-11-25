package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarVariables(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var x = 0;
let y = 1;
const z = 2;
// @Filename: file2.ts
var {a} = 0;
let {a: b} = 0;
const [c] = 0;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "x",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "y",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "z",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
	f.GoToFile(t, "file2.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "a",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "b",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "c",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
