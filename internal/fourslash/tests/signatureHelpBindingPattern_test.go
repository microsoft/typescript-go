package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Tests for signature help with binding pattern parameters.
// This covers the crash fix for binding patterns and various combinations
// as requested in the issue.
func TestSignatureHelpBindingPattern(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// Empty object binding pattern
function emptyObj({}) {}
emptyObj(/*emptyObj*/)

// Empty array binding pattern
function emptyArr([]) {}
emptyArr(/*emptyArr*/)

// Non-empty object binding pattern
function nonEmptyObj({a, b}: {a: number, b: string}) {}
nonEmptyObj(/*nonEmptyObj*/)

// Non-empty array binding pattern
function nonEmptyArr([x, y]: [number, string]) {}
nonEmptyArr(/*nonEmptyArr*/)

// Identifiers leading, binding pattern trailing
function idLeading(first: number, {a, b}: {a: number, b: string}) {}
idLeading(/*idLeading*/)

// Binding pattern leading, identifiers trailing
function bindingLeading({a, b}: {a: number, b: string}, last: number) {}
bindingLeading(/*bindingLeading*/)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineSignatureHelp(t)
}
