package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsImports(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import d1 from "a";

import { a } from "a";

import { b as B } from "a" 

import d2, { c, d as D } from "a" 

import e = require("a");

import * as ns from "a";`
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
			Name:     "c",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "D",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "d1",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "d2",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "e",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name:     "ns",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
	})
}
