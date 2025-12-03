package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsMultilineStringIdentifiers3(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module 'MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters' { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "'MoreThanOneHundredAndFiftyCharactersMoreThanOneHundredAndFiftyCharactersMoreThanOneHundredAndFiftyCharactersMoreThanOneHundredAndFiftyCharacter...",
			Kind:     lsproto.SymbolKindNamespace,
			Children: nil,
		},
	})
}
