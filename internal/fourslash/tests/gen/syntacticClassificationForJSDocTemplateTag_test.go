package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationForJSDocTemplateTag(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/** @template T baring strait */
function ident<T>: T {
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "function.declaration", Text: "ident"},
		{Type: "typeParameter.declaration", Text: "T"},
		{Type: "typeParameter", Text: "T"},
	})
}
