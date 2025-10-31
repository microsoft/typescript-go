package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsFunctionWithComments(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**
 * This is my function.
 * There are many like it, but this one is mine.
 */
function myFunction(/* x */ x: any) {
    var y = x ? x++ : ++x;
}
// end of file`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "function.declaration", Text: "myFunction"},
		{Type: "parameter.declaration", Text: "x"},
		{Type: "variable.declaration.local", Text: "y"},
		{Type: "parameter", Text: "x"},
		{Type: "parameter", Text: "x"},
		{Type: "parameter", Text: "x"},
	})
}
