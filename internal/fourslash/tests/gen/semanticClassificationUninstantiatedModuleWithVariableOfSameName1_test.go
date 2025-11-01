package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationUninstantiatedModuleWithVariableOfSameName1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module /*0*/M {
    interface /*1*/I {

    }
}

var M = { I: 10 };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable", Text: "M"},
		{Type: "interface.declaration", Text: "I"},
		{Type: "variable.declaration", Text: "M"},
		{Type: "property.declaration", Text: "I"},
	})
}
