package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsTemplates2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var tiredOfCanonicalExamples =
` + "`" + `goodbye "${ ` + "`" + `hello world` + "`" + ` }" 
and ${ ` + "`" + `good${ " " }riddance` + "`" + ` }` + "`" + `;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration", Text: "tiredOfCanonicalExamples"},
	})
}
