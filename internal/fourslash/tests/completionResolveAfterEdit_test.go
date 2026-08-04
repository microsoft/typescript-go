package fourslash_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionResolveAfterEdit(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: /index.ts
interface Point {
	x: number;
	y: number;
}
declare const p: Point;
/*a*/

// @filename: /foo.ts
/*b*/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Step 1: Get completions at the marker.
	f.GoToMarker(t, "a")
	completions := f.GetCompletions(t, nil /*userPreferences*/)
	if completions == nil || len(completions.Items) == 0 {
		t.Fatal("Expected completions but got none")
	}
	firstItem := completions.Items[0]

	// Step 2: Make a file change (insert a comment after marker).
	f.GoToMarker(t, "b")
	f.Insert(t, "1")

	// Step 3: Resolve the first completion item from the original list.
	resolved := f.ResolveCompletionItem(t, firstItem)
	if resolved == nil {
		t.Fatal("Expected resolved completion item but got nil")
	}
}

func TestClassMemberCompletionResolveAfterEdit(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: /a.ts
export type T = string;
export interface I {
	method(value: T): T;
}

// @filename: /b.ts
import { I } from "./a";
class C implements I {
	/*a*/
}

// @filename: /c.ts
/*b*/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsWithClassMemberSnippets = core.TSTrue
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	itemToResolve := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "method" && item.Data != nil && item.Data.Source == ls.SourceClassMemberSnippet
	})
	if itemToResolve == nil {
		t.Fatal("Expected class member completion for method")
	}

	f.GoToMarker(t, "b")
	f.Insert(t, "1")

	resolved := f.ResolveCompletionItem(t, itemToResolve)
	if resolved == nil || resolved.AdditionalTextEdits == nil || len(*resolved.AdditionalTextEdits) == 0 {
		t.Fatal("Expected resolved class member completion with import edits")
	}
	if resolved.Detail == nil || !strings.Contains(*resolved.Detail, "Includes imports of types referenced by") {
		t.Fatal("Expected resolved class member completion to describe its import edits")
	}
}
