package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test case 1: Basic applicable range test
// This test verifies that signature help is NOT provided on whitespace after closing paren
func TestSignatureHelpApplicableRangeBasic(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let obj = {
    foo(s: string): string {
        return s;
    }
};

let s = obj.foo(/*1*/"Hello, world!"/*2*/)/*3*/  
  /*4*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	
	// Use VerifyBaselineSignatureHelp to check behavior at all markers
	f.VerifyBaselineSignatureHelp(t)
}

// Test case 2: Nested calls - outer should take precedence
func TestSignatureHelpNestedCalls(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(s: string) { return s; }
function bar(s: string) { return s; }
let s = foo(/*a*//*b*/bar/*c*/(/*d*/"hello"/*e*/)/*f*/);`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	
	// Use VerifyBaselineSignatureHelp to check behavior at all markers
	f.VerifyBaselineSignatureHelp(t)
}
