package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationConstructorTypes(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `Object.create(null);
const x = Promise.resolve(Number.MAX_VALUE);
if (x instanceof Promise) {}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "class.defaultLibrary", Text: "Object"},
		{Type: "method.defaultLibrary", Text: "create"},
		{Type: "variable.declaration.readonly", Text: "x"},
		{Type: "class.defaultLibrary", Text: "Number"},
		{Type: "property.readonly.defaultLibrary", Text: "MAX_VALUE"},
		{Type: "variable.readonly", Text: "x"},
	})
}
