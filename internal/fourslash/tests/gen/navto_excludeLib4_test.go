package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavto_excludeLib4(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @filename: /node_modules/bar/index.d.ts
import { someOtherName } from "./baz";
export const [|someName: number|];
// @filename: /node_modules/bar/baz.d.ts
export const someOtherName: string;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyWorkspaceSymbol(t, []*fourslash.VerifyWorkspaceSymbolCase{
		{
			Preferences: &ls.UserPreferences{ExcludeLibrarySymbols: PtrTo(true)},
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "someName",
					Kind:     lsproto.SymbolKindVariable,
					Location: f.Ranges()[0].LSLocation(),
				},
			},
		},
	})
}
