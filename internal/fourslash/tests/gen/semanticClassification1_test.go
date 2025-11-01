package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassification1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module /*0*/M {
    export interface /*1*/I {
    }
}
interface /*2*/X extends /*3*/M./*4*/I { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "namespace.declaration", Text: "M"},
		{Type: "interface.declaration", Text: "I"},
		{Type: "interface.declaration", Text: "X"},
		{Type: "namespace", Text: "M"},
		{Type: "interface", Text: "I"},
	})
}
