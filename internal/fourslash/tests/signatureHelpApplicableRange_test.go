package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test case 1: Basic applicable range test
// This test verifies the applicable range for signature help
// According to the issue, signature help should be provided:
// - Inside the parentheses (markers 1, 2)
// - NOT on the call target before the opening paren (marker a, b, c would test this)
// - NOT after the closing paren including whitespace (markers 3, 4)
func TestSignatureHelpApplicableRangeBasic(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let obj = {
    foo(s: string): string {
        return s;
    }
};

let s =/*a*/ /*b*/obj/*c*/./*d*/foo/*e*/(/*1*/"Hello, world!"/*2*/)/*3*/  
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
