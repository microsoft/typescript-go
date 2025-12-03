package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsMultilineStringIdentifiers1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module "Multiline\r\nMadness" {
}

declare module "Multiline\
Madness" {
}
declare module "MultilineMadness" {}

declare module "Multiline\
Madness2" {
}

interface Foo {
    "a1\\\r\nb";
    "a2\
    \
    b"(): Foo;
}

class Bar implements Foo {
    'a1\\\r\nb': Foo;

    'a2\
    \
    b'(): Foo {
        return this;
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "\"MultilineMadness\"",
			Kind:     lsproto.SymbolKindNamespace,
			Children: nil,
		},
		{
			Name:     "\"MultilineMadness2\"",
			Kind:     lsproto.SymbolKindNamespace,
			Children: nil,
		},
		{
			Name:     "\"Multiline\\r\\nMadness\"",
			Kind:     lsproto.SymbolKindNamespace,
			Children: nil,
		},
		{
			Name: "Bar",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "'a1\\\\\\r\\nb'",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "'a2        b'",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
		{
			Name: "Foo",
			Kind: lsproto.SymbolKindInterface,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "\"a1\\\\\\r\\nb\"",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
				{
					Name:     "\"a2        b\"",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
			}),
		},
	})
}
