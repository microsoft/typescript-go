package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionResolveAfterEdit(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: a.ts
interface I {
	x: number;
	y: number;
}
declare const u: I;
/*a*/

// @filename: 1.ts
/*b*/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	completions := f.GetCompletions(t, nil /*userPreferences*/)
	if completions == nil || len(completions.Items) == 0 {
		t.Fatal("Expected completions but got none")
	}
	firstItem := completions.Items[0]

	f.GoToMarker(t, "b")
	f.Insert(t, "1")

	resolved := f.ResolveCompletionItem(t, firstItem)
	if resolved == nil {
		t.Fatal("Expected resolved completion item but got nil")
	}
}

func TestResolveImportStatementCompletion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: a.ts
export const u = 1;

// @filename: 1.ts
[|import u/*a*/|]
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	completions := f.GetCompletions(t, nil /*userPreferences*/)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	item := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "u"
	})
	if item == nil {
		t.Fatal("Expected import statement completion for u")
	}
	if item.Data == nil || item.Data.AutoImport == nil || !item.Data.IsImportStatementCompletion {
		t.Fatalf("Expected import statement completion data, got %#v", item.Data)
	}
	resolved := f.ResolveCompletionItem(t, item)
	if resolved == nil {
		t.Fatal("Expected resolved import statement completion")
	}
	if resolved.AdditionalTextEdits != nil && len(*resolved.AdditionalTextEdits) != 0 {
		t.Fatal("Expected import statement completion to have no additional import edits")
	}
}
