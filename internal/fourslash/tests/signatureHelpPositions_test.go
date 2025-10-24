package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Simple test to understand positions
func TestSignatureHelpPositions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(s: string) { return s; }
let s = foo(/*1*/"hello"/*2*/)/*3*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	
	f.VerifyBaselineSignatureHelp(t)
}
