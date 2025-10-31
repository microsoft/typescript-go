package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationInfinityAndNaN(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = ` Infinity;
 NaN;

// Regular properties

const obj1 = {
    Infinity: 100,
    NaN: 200,
    "-Infinity": 300
};

obj1.Infinity;
obj1.NaN;
obj1["-Infinity"];

// Shorthand properties

const obj2 = {
    Infinity,
    NaN,
}

obj2.Infinity;
obj2.NaN;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "variable.declaration.readonly", Text: "obj1"},
		{Type: "variable.readonly", Text: "obj1"},
		{Type: "variable.readonly", Text: "obj1"},
		{Type: "variable.readonly", Text: "obj1"},
		{Type: "variable.declaration.readonly", Text: "obj2"},
		{Type: "variable.readonly", Text: "obj2"},
		{Type: "variable.readonly", Text: "obj2"},
	})
}
