package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test for crash when requesting signature help for a function with binding pattern parameters.
// This verifies that signature help works for binding patterns without crashing,
// even though JSDoc parameter documentation is not available for binding patterns.
func TestSignatureHelpBindingPattern(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
function foo({}) {
}

foo(/*$*/)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	
	// We expect signature help to work without crashing
	f.GoToMarker(t, "$")
	f.VerifySignatureHelpPresent(t, &lsproto.SignatureHelpContext{
		IsRetrigger:      false,
		TriggerCharacter: PtrTo("("),
		TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
	})
}

// Test that signature help works with JSDoc on functions with binding pattern parameters.
// Note: JSDoc @param tags cannot match binding patterns, so the JSDoc comment won't
// provide parameter-specific documentation, but the signature help should still work.
func TestSignatureHelpBindingPatternWithJSDoc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
/**
 * A function with a binding pattern parameter
 */
function foo({a, b}: {a: number, b: string}) {
}

foo(/*$*/)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	
	// We expect signature help to work without crashing
	f.GoToMarker(t, "$")
	f.VerifySignatureHelpPresent(t, &lsproto.SignatureHelpContext{
		IsRetrigger:      false,
		TriggerCharacter: PtrTo("("),
		TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
	})
}

// Test array binding pattern
func TestSignatureHelpArrayBindingPattern(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
function bar([x, y]: [number, number]) {
}

bar(/*$*/)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	
	f.GoToMarker(t, "$")
	f.VerifySignatureHelpPresent(t, &lsproto.SignatureHelpContext{
		IsRetrigger:      false,
		TriggerCharacter: PtrTo("("),
		TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
	})
}
