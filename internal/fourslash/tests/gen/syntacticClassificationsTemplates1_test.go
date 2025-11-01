package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsTemplates1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var v = 10e0;
var x = {
    p1: ` + "`" + `hello world` + "`" + `,
    p2: ` + "`" + `goodbye ${0} cruel ${0} world` + "`" + `,
};`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration", Text: "v"},
		{Type: "variable.declaration", Text: "x"},
		{Type: "property.declaration", Text: "p1"},
		{Type: "property.declaration", Text: "p2"},
	})
}
