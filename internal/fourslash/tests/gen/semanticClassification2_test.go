package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassification2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface /*0*/Thing {
    toExponential(): number;
}

var Thing = 0;
Thing.toExponential();`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "interface.declaration", Text: "Thing"},
		{Type: "method.declaration", Text: "toExponential"},
		{Type: "variable.declaration", Text: "Thing"},
		{Type: "variable", Text: "Thing"},
		{Type: "method.defaultLibrary", Text: "toExponential"},
	})
}
