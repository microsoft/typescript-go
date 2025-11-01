package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassifications1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// comment
module M {
    var v = 0 + 1;
    var s = "string";

    class C<T> {
    }

    enum E {
    }

    interface I {
    }

    module M1.M2 {
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "namespace.declaration", Text: "M"},
		{Type: "variable.declaration.local", Text: "v"},
		{Type: "variable.declaration.local", Text: "s"},
		{Type: "class.declaration", Text: "C"},
		{Type: "typeParameter.declaration", Text: "T"},
		{Type: "enum.declaration", Text: "E"},
		{Type: "interface.declaration", Text: "I"},
		{Type: "namespace.declaration", Text: "M1"},
		{Type: "namespace.declaration", Text: "M2"},
	})
}
