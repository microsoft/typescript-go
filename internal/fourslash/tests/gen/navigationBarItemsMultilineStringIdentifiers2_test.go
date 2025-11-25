package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsMultilineStringIdentifiers2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function f(p1: () => any, p2: string) { }
f(() => { }, ` + "`" + `line1\
line2\
line3` + "`" + `);

class c1 {
    const a = ' ''line1\
        line2';
}

f(() => { }, ` + "`" + `unterminated backtick 1
unterminated backtick 2
unterminated backtick 3`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "c1",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "a",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "'line1        line2'",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
		{
			Name:     "f",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name:     "f(`line1line2line3`) callback",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
		{
			Name:     "f(`unterminated backtick 1unterminated backtick 2unterminated backtick 3) callback",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
	})
}
