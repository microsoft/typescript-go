package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsBindingPatterns(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `'use strict'
var foo, {}
var bar, []
let foo1, {a, b}
const bar1, [c, d]
var {e, x: [f, g]} = {a:1, x:[]};
var { h: i = function j() {} } = obj;`
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
			Name:     "bar",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "bar1",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "c",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "d",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "e",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "f",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "foo",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "foo1",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "g",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "i",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
