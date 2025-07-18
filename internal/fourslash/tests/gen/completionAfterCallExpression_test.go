package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionAfterCallExpression(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	
	// Test case 1: Simple call expression followed by cursor
	const content1 = `let x = someCall() /**/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content1)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"satisfies",
				"as",
			},
		},
	})
	
	// Test case 2: Call expression with arguments followed by cursor
	const content2 = `let y = anotherCall(1, 2, 3) /**/`
	f2 := fourslash.NewFourslash(t, nil /*capabilities*/, content2)
	f2.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"satisfies",
				"as",
			},
		},
	})
}
