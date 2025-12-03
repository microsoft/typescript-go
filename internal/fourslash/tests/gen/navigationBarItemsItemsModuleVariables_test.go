package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsItemsModuleVariables(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: navigationItemsModuleVariables_0.ts
 /*file1*/
module Module1 {
    export var x = 0;
}
// @Filename: navigationItemsModuleVariables_1.ts
 /*file2*/
module Module1.SubModule {
    export var y = 0;
}
// @Filename: navigationItemsModuleVariables_2.ts
 /*file3*/
module Module1 {
    export var z = 0;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToMarker(t, "file1")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "Module1",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
	f.GoToMarker(t, "file2")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "Module1.SubModule",
			Kind: lsproto.SymbolKindNamespace,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "y",
					Kind:     lsproto.SymbolKindVariable,
					Children: nil,
				},
			}),
		},
	})
}
