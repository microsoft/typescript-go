package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsForOfKeyword(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `for (var of of of) { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration", Text: "of"},
		{Type: "variable", Text: "of"},
	})
}
