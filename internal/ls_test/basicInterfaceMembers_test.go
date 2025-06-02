package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
)

func TestBasicInterfaceMembers(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export {};
interface Point {
	x: number;
	y: number;
}
declare const p: Point;
p./*a*/`
	f, done := lstestutil.NewFourslash(t, nil /*capabilities*/, content, "basicInterfaceMembers.ts")
	f.VerifyCompletions(t, "a", &lstestutil.VerifyCompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &lsproto.CompletionItemDefaults{
			CommitCharacters: &lstestutil.DefaultCommitCharacters,
		},
		Items: &lstestutil.VerifyCompletionsExpectedItems{
			Exact: []lstestutil.ExpectedCompletionItem{
				&lsproto.CompletionItem{
					Label:      "x",
					Kind:       lstestutil.PtrTo(lsproto.CompletionItemKindField),
					SortText:   lstestutil.PtrTo(string(ls.SortTextLocationPriority)),
					FilterText: lstestutil.PtrTo(".x"),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						InsertReplaceEdit: &lsproto.InsertReplaceEdit{
							NewText: "x",
							Insert: lsproto.Range{
								Start: lsproto.Position{Line: 6, Character: 2},
								End:   lsproto.Position{Line: 6, Character: 2},
							},
							Replace: lsproto.Range{
								Start: lsproto.Position{Line: 6, Character: 2},
								End:   lsproto.Position{Line: 6, Character: 2},
							},
						},
					},
				},
				"y",
			},
		},
	})
	done()
}
