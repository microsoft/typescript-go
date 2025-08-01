package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionListInvalidMemberNames(t *testing.T) {
	t.Parallel()
	// t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var x = {
    "foo ": "space in the name",
    "bar": "valid identifier name",
    "break": "valid identifier name (matches a keyword)",
    "any": "valid identifier name (matches a typescript keyword)",
    "#": "invalid identifier name",
    "$": "valid identifier name",
    "\u0062": "valid unicode identifier name (b)",
    "\u0031\u0062": "invalid unicode identifier name (1b)"
};

x[|./*a*/|];
x["[|/*b*/|]"];`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "b", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Unsorted: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "foo ",
				},
				&lsproto.CompletionItem{
					Label: "bar",
				},
				&lsproto.CompletionItem{
					Label: "break",
				},
				&lsproto.CompletionItem{
					Label: "any",
				},
				&lsproto.CompletionItem{
					Label: "#",
				},
				&lsproto.CompletionItem{
					Label: "$",
				},
				&lsproto.CompletionItem{
					Label: "b",
				},
				&lsproto.CompletionItem{
					Label: "1b",
				},
			},
		},
	})
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Unsorted: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:      "foo ",
					InsertText: PtrTo("[\"foo \"]"),
				},
				"bar",
				"break",
				"any",
				&lsproto.CompletionItem{
					Label:      "#",
					InsertText: PtrTo("[\"#\"]"),
				},
				"$",
				"b",
				&lsproto.CompletionItem{
					Label:      "1b",
					InsertText: PtrTo("[\"1b\"]"),
				},
			},
		},
	})
}
