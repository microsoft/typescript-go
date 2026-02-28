package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionAfterTrailingAtInJSDoc1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: /atTagPosition.js
/**
 * @/*1*/
 */
function foo(x) {}

// @Filename: /atAfterExistingParam.js
/**
 * @param {string} x ok
 * @/*2*/
 */
function bar(x, y) {}

// @Filename: /atMidLine.js
/**
 * some text @/*3*/
 */
function baz(y) {}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	expectTagCompletions := &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "param",
					Kind:  new(lsproto.CompletionItemKindKeyword),
				},
			},
		},
	}
	f.VerifyCompletions(t, "1", expectTagCompletions)
	f.VerifyCompletions(t, "2", expectTagCompletions)
	f.VerifyCompletions(t, "3", nil)
}
