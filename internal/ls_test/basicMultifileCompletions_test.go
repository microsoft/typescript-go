package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
)

func TestBasicMultifileCompletions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export const foo = { bar: 'baz' };

// @Filename: /b.ts
import { foo } from './a';
const test = foo./*1*/`
	f, done := lstestutil.NewFourslash(t, nil /*capabilities*/, content, "basicMultifileCompletions.ts")
	f.VerifyCompletions(t, "1", &lstestutil.VerifyCompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &lsproto.CompletionItemDefaults{
			CommitCharacters: &lstestutil.DefaultCommitCharacters,
		},
		Items: &lstestutil.VerifyCompletionsExpectedItems{
			Includes: []lstestutil.ExpectedCompletionItem{
				&lsproto.CompletionItem{
					Label:      "bar",
					Kind:       lstestutil.PtrTo(lsproto.CompletionItemKindField),
					SortText:   lstestutil.PtrTo(string(ls.SortTextLocationPriority)),
					FilterText: lstestutil.PtrTo(".bar"),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						InsertReplaceEdit: &lsproto.InsertReplaceEdit{
							NewText: "bar",
							Insert: lsproto.Range{
								Start: lsproto.Position{Line: 1, Character: 17},
								End:   lsproto.Position{Line: 1, Character: 17},
							},
							Replace: lsproto.Range{
								Start: lsproto.Position{Line: 1, Character: 17},
								End:   lsproto.Position{Line: 1, Character: 17},
							},
						},
					},
				},
			},
		},
	})
	done()
}
