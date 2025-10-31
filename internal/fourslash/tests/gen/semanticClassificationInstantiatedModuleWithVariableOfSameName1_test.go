package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationInstantiatedModuleWithVariableOfSameName1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module /*0*/M {
    export interface /*1*/I {
    }
    var x = 10;
}

var /*2*/M = {
    foo: 10,
    bar: 20
}

var v: /*3*/M./*4*/I;

var x = /*5*/M;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "namespace.declaration", Text: "M"},
		{Type: "interface.declaration", Text: "I"},
		{Type: "variable.declaration.local", Text: "x"},
		{Type: "variable.declaration", Text: "M"},
		{Type: "property.declaration", Text: "foo"},
		{Type: "property.declaration", Text: "bar"},
		{Type: "variable.declaration", Text: "v"},
		{Type: "namespace", Text: "M"},
		{Type: "interface", Text: "I"},
		{Type: "variable.declaration", Text: "x"},
		{Type: "namespace", Text: "M"},
	})
}
