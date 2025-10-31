package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationObjectProperties(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let x = 1, y = 1;
const a1 = { e: 1 };
var a2 = { x };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "variable.declaration", Text: "x"}, {Type: "variable.declaration", Text: "y"}, {Type: "variable.declaration.readonly", Text: "a1"}, {Type: "property.declaration", Text: "e"}, {Type: "variable.declaration", Text: "a2"}, {Type: "property.declaration", Text: "x"}})
}
