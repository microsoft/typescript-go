package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsObjectLiteral(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var v = 10e0;
var x = {
    p1: 1,
    p2: 2,
    any: 3,
    function: 4,
    var: 5,
    void: void 0,
    v: v += v,
};`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration", Text: "v"},
		{Type: "variable.declaration", Text: "x"},
		{Type: "property.declaration", Text: "p1"},
		{Type: "property.declaration", Text: "p2"},
		{Type: "property.declaration", Text: "any"},
		{Type: "property.declaration", Text: "function"},
		{Type: "property.declaration", Text: "var"},
		{Type: "property.declaration", Text: "void"},
		{Type: "property.declaration", Text: "v"},
		{Type: "variable", Text: "v"},
		{Type: "variable", Text: "v"},
	})
}
