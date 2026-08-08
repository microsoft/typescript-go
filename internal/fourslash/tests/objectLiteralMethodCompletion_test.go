package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestObjectLiteralMethodCompletionPreservesOptionalLabel(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface I {
	method?(): void
}
const u: I = {
	/*a*/
}`
	f, done := fourslash.NewFourslash(t, fourslash.GetDefaultCapabilitiesWithOptions(&fourslash.ClientCapabilitiesOptions{
		CompletionItem: &lsproto.ClientCompletionItemOptions{
			LabelDetailsSupport: new(true),
		},
	}), content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsWithObjectLiteralMethodSnippets = core.TSTrue
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	count := 0
	for _, item := range completions.Items {
		if item.Label == "method?" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("Expected two optional method completions, got %d", count)
	}
}
