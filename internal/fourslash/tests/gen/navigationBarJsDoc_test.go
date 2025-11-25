package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarJsDoc(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.js
/** @typedef {(number|string)} NumberLike */
/** @typedef {(string|number)} */
const x = 0;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "NumberLike",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name:     "x",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "x",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
}
