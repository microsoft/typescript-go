package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarAnonymousClassAndFunctionExpressions3(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `describe('foo', () => {
    test(` + "`" + `a ${1} b ${2}` + "`" + `, () => {})
})

const a = 1;
const b = 2;
describe('foo', () => {
    test(` + "`" + `a ${a} b {b}` + "`" + `, () => {})
})`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
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
			Name: "describe('foo') callback",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "test(`a ${1} b ${2}`) callback",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
		{
			Name: "describe('foo') callback",
			Kind: lsproto.SymbolKindFunction,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "test(`a ${a} b {b}`) callback",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
			}),
		},
	})
}
