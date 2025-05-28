package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
)

const content = `export {};
interface Point {
    x: number;
    y: number;
}
declare const p: Point;
p./*a*/`

func TestBasicInterfaceMembers(t *testing.T) {
	t.Parallel()
	cap := &lsproto.ClientCapabilities{}
	f, done := lstestutil.NewFourslash(t, cap, content, "basicInterfaceMembers.ts")
	f.VerifyCompletions(t, "a", lstestutil.VerifyCompletionsResult{
		Exact: &lsproto.CompletionList{
			IsIncomplete: false,
			ItemDefaults: &lsproto.CompletionItemDefaults{
				CommitCharacters: &lstestutil.DefaultCommitCharacters,
			},
			Items: []*lsproto.CompletionItem{
				{
					Label:            "x",
					Kind:             lstestutil.PtrTo(lsproto.CompletionItemKindField),
					SortText:         lstestutil.PtrTo(string(ls.SortTextLocationPriority)),
					FilterText:       lstestutil.PtrTo(".x"),
					InsertTextFormat: lstestutil.PtrTo(lsproto.InsertTextFormatPlainText),
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
				{
					Label:            "y",
					Kind:             lstestutil.PtrTo(lsproto.CompletionItemKindField),
					SortText:         lstestutil.PtrTo(string(ls.SortTextLocationPriority)),
					FilterText:       lstestutil.PtrTo(".y"),
					InsertTextFormat: lstestutil.PtrTo(lsproto.InsertTextFormatPlainText),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						InsertReplaceEdit: &lsproto.InsertReplaceEdit{
							NewText: "y",
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
			},
		},
	})
	done()
}
