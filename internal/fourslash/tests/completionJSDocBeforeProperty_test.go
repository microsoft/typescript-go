package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test for: Crash completing beginning of property name when preceded by JSDoc
// This test verifies that requesting completions at the beginning of a property name
// preceded by JSDoc does not cause a crash. The crash was caused by a nil pointer
// dereference when contextToken was nil.
func TestCompletionJSDocBeforePropertyNoCrash(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export class SomeInterface {
    /** ruh-roh! */
    /*a*/property: string;
}

export class SomeClass {
    /** ruh-roh! */
    /*b*/property = "value";
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	// The primary goal of this test is to ensure no panic occurs when requesting completions.
	// The testutil.RecoverAndFail defer will catch any panic and fail the test.
	// We verify completions can be requested successfully without checking specific completion items.
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{".", ",", ";"},
			EditRange:        fourslash.Ignored{},
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{},
		},
	})
	f.VerifyCompletions(t, "b", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{".", ",", ";"},
			EditRange:        fourslash.Ignored{},
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{},
		},
	})
}
