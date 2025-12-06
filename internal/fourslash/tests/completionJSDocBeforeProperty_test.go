package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionJSDocBeforeProperty(t *testing.T) {
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
	// Should not crash when requesting completions at the beginning of property name preceded by JSDoc
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
	})
	f.VerifyCompletions(t, "b", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
	})
}
