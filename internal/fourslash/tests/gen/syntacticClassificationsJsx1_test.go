package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsJsx1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: file1.tsx
let x  = <div a = "some-value" b = {1}>
    some jsx text
</div>;

let y = <element attr="123"/>`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration", Text: "x"},
		{Type: "variable.declaration", Text: "y"},
	})
}
