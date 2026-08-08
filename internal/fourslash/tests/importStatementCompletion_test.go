package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestImportStatementCompletionUsesNamedImport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
export interface I {}
// @Filename: 1.ts
import * as u from "./a";
import I/*a*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsForModuleExports = core.TSFalse
	preferences.IncludeCompletionsForImportStatements = core.TSTrue
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	item := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "I"
	})
	if item == nil {
		t.Fatal("Expected import statement completion for I")
	}
	if item.TextEdit == nil || item.TextEdit.TextEdit == nil || item.TextEdit.TextEdit.NewText != `import { I } from "./a";` {
		t.Fatalf("Expected named import text edit, got %#v", item.TextEdit)
	}
}
