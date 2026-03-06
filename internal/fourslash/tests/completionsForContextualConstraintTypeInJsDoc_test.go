package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionsForContextualConstraintTypeInJsDoc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
// @allowJs: true
// @filename: a.ts
export interface Blah<T extends { a: "hello" | "world" }> {
}

// @filename: b.js
/** @import * as a from "./a" */

/** @type {a.Blah<{ a: /*$*/ }>} */
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// The assertion here is simply "does not crash/panic".
	f.VerifyCompletions(t, "$", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{CommitCharacters: &[]string{".", ",", ";"}},
		Items:        &fourslash.CompletionsExpectedItems{},
	})
}
