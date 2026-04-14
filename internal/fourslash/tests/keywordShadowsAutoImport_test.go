package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test that global keywords like `function`, `class`, and `const` shadow
// auto-import completions with the same name, rather than being
// shadowed by them. See: https://github.com/microsoft/typescript-go/issues/1379
func TestKeywordShadowsAutoImport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @Filename: /mod.ts
const value = 1;
export { value as function }

// @Filename: /index.ts
function/**/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	// The keyword `function` should appear, but NOT the auto-import `function` from ./mod.
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:    "function",
					Kind:     new(lsproto.CompletionItemKindKeyword),
					SortText: new(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
			Excludes: []string{},
		},
	})
}

// Test that contextual keywords like `type` do NOT shadow auto-imports with the same name.
// Both the keyword and auto-import should appear.
func TestContextualKeywordDoesNotShadowAutoImport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @Filename: /mod.ts
const value = 1;
export { value as type }

// @Filename: /index.ts
type/**/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	// Both the keyword `type` and the auto-import `type` should appear.
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				// The auto-import "type" from ./mod should appear (contextual keyword, coexists)
				&lsproto.CompletionItem{
					Label: "type",
					Data: &lsproto.CompletionItemData{
						AutoImport: &lsproto.AutoImportFix{
							ModuleSpecifier: "./mod",
						},
					},
					LabelDetails: &lsproto.CompletionItemLabelDetails{
						Description: new("./mod"),
					},
					SortText:            new(string(ls.SortTextAutoImportSuggestions)),
					AdditionalTextEdits: fourslash.AnyTextEdits,
				},
			},
		},
	})
}
