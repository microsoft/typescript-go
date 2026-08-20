package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoImportPatternAmbientModule(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /tsconfig.json
{ "compilerOptions": { "module": "preserve", "moduleResolution": "bundler" } }

// @Filename: /types.d.ts
declare module "*.asset" with { type: "css" } {
    export const styles: string;
}

// @Filename: /index.ts
sty/**/`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "styles",
					Data: &lsproto.CompletionItemData{
						AutoImport: &lsproto.AutoImportFix{
							ModuleSpecifier: "*.asset",
						},
					},
					AdditionalTextEdits: fourslash.AnyTextEdits,
					SortText:            new(string(ls.SortTextAutoImportSuggestions)),
				},
			},
		},
	})
	f.VerifyApplyCodeActionFromCompletion(t, new(""), &fourslash.ApplyCodeActionFromCompletionOptions{
		Name:        "styles",
		Source:      "*.asset",
		Description: "Add import from \"*.asset\"",
		NewFileContent: new(`import { styles } from "*.asset";

sty`),
		AutoImportFix: &lsproto.AutoImportFix{
			ModuleSpecifier: "*.asset",
		},
	})
}
