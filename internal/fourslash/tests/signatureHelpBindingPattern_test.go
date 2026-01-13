package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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
emptyObj(/*1*/)

// Empty array binding pattern
function emptyArr([]) {}
emptyArr(/*2*/)

// Non-empty object binding pattern
function nonEmptyObj({a, b}: {a: number, b: string}) {}
nonEmptyObj(/*3*/)

// Non-empty array binding pattern
function nonEmptyArr([x, y]: [number, string]) {}
nonEmptyArr(/*4*/)

// Identifiers leading, binding pattern trailing
function idLeading(first: number, {a, b}: {a: number, b: string}) {}
idLeading(/*5*/)

// Binding pattern leading, identifiers trailing
function bindingLeading({a, b}: {a: number, b: string}, last: number) {}
bindingLeading(/*6*/)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	ctx := &lsproto.SignatureHelpContext{
		IsRetrigger:      false,
		TriggerCharacter: PtrTo("("),
		TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
	}

	// Test all markers - each should work without crashing
	f.GoToMarker(t, "1")
	f.VerifySignatureHelpPresent(t, ctx)

	f.GoToMarker(t, "2")
	f.VerifySignatureHelpPresent(t, ctx)

	f.GoToMarker(t, "3")
	f.VerifySignatureHelpPresent(t, ctx)

	f.GoToMarker(t, "4")
	f.VerifySignatureHelpPresent(t, ctx)

	f.GoToMarker(t, "5")
	f.VerifySignatureHelpPresent(t, ctx)

	f.GoToMarker(t, "6")
	f.VerifySignatureHelpPresent(t, ctx)
}
