package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavto_serverExcludeLib(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /home/src/workspaces/project/index.ts
import { weirdName as otherName } from "bar";
const [|weirdName: number = 1|];
// @Filename: /home/src/workspaces/project/tsconfig.json
{}
// @Filename: /home/src/workspaces/project/node_modules/bar/index.d.ts
export const [|weirdName: number|];
// @Filename: /home/src/workspaces/project/node_modules/bar/package.json
{}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyWorkspaceSymbol(t, []*fourslash.VerifyWorkspaceSymbolCase{
		{
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "weirdName",
					Kind:     lsproto.SymbolKindVariable,
					Location: f.Ranges()[1].LSLocation(),
				},
				{
					Name:     "weirdName",
					Kind:     lsproto.SymbolKindVariable,
					Location: f.Ranges()[0].LSLocation(),
				},
			},
		},
	})
	f.VerifyWorkspaceSymbol(t, []*fourslash.VerifyWorkspaceSymbolCase{
		{
			Preferences: &ls.UserPreferences{ExcludeLibrarySymbols: PtrTo(true)},
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "weirdName",
					Kind:     lsproto.SymbolKindVariable,
					Location: f.Ranges()[0].LSLocation(),
				},
			},
		},
	})
}
