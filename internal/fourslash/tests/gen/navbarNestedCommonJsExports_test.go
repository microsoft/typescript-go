package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavbarNestedCommonJsExports(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: /a.js
exports.a = exports.b = exports.c = 0;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "a",
			Kind: lsproto.SymbolKindVariable,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "b",
					Kind: lsproto.SymbolKindVariable,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "c",
							Kind:     lsproto.SymbolKindVariable,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
