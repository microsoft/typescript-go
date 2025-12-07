package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test that string literal completions are suggested in tuple contexts
// even without typing a quote character.
func TestCompletionsInArrayLiteralWithContextualType(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	// Test 1: Completions after `[` in a tuple should suggest string literals
	const content1 = `let y: ["foo" | "bar", string] = [/*a*/];`
	f1, done1 := fourslash.NewFourslash(t, nil /*capabilities*/, content1)
	defer done1()
	f1.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"\"foo\"",
				"\"bar\"",
			},
		},
	})

	// Test 2: Completions after `,` in a tuple should provide contextual type for second element
	const content2 = `let z: ["a", "b" | "c"] = ["a", /*b*/];`
	f2, done2 := fourslash.NewFourslash(t, nil /*capabilities*/, content2)
	defer done2()
	f2.VerifyCompletions(t, "b", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"\"b\"",
				"\"c\"",
			},
		},
	})
}
