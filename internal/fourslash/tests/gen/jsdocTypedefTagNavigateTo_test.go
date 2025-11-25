package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestJsdocTypedefTagNavigateTo(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowNonTsExtensions: true
// @Filename: jsDocTypedef_form2.js

/** @typedef {(string | number)} NumberLike */
/** @typedef {(string | number | string[])} */
var NumberLike2;

/** @type {/*1*/NumberLike} */
var numberLike;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.MarkTestAsStradaServer()
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "numberLike",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "NumberLike",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
		{
			Name:     "NumberLike2",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "NumberLike2",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
}
