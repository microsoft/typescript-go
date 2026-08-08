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

func TestClassMemberCompletionIncludesImportEdits(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: a.ts
export type T = string;
export interface I {
	method(value: T): T;
}

// @filename: 1.ts
import { I } from "./a";
class C implements I {
	/*a*/
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsWithClassMemberSnippets = core.TSTrue
	preferences.IncludeCompletionsForModuleExports = core.TSFalse
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	item := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "method" && item.Data != nil && item.Data.Source == ls.SourceClassMemberSnippet
	})
	if item == nil {
		t.Fatal("Expected class member completion for method")
	}
	if item.AdditionalTextEdits == nil || len(*item.AdditionalTextEdits) == 0 {
		t.Fatal("Expected class member completion with import edits")
	}
}

func TestClassMemberCompletionKeepsNameFallback(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class B {
	constructor(public value: string) {}
}
class C extends B {
	/*a*/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsWithClassMemberSnippets = core.TSTrue
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	item := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "value"
	})
	if item == nil {
		t.Fatal("Expected class member completion for value")
	}
	if item.InsertText == nil || *item.InsertText != "value" {
		t.Fatalf("Expected name-only insertion, got %#v", item.InsertText)
	}
}

func TestClassMemberCompletionAddsRequiredOverride(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noImplicitOverride: true
class B {
	method() {}
}
class C extends B {
	/*a*/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "a")
	preferences := lsutil.NewDefaultUserPreferences()
	preferences.IncludeCompletionsWithClassMemberSnippets = core.TSTrue
	completions := f.GetCompletions(t, &preferences)
	if completions == nil {
		t.Fatal("Expected completions but got none")
	}
	item := core.Find(completions.Items, func(item *lsproto.CompletionItem) bool {
		return item.Label == "method"
	})
	if item == nil || item.InsertText == nil || !strings.Contains(*item.InsertText, "override method()") {
		t.Fatalf("Expected completion with override modifier, got %#v", item)
	}
}

func TestImplementClassFixDoesNotAddInvalidOverride(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noImplicitOverride: true
class B {
    method() {}
}
class C implements B {[| |]}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Implement interface 'B'",
		NewFileContent: `class B {
    method() {}
}
class C implements B {
    method(): void {
        throw new Error("Method not implemented.");
    }
}`,
		Index: 0,
	})
}
