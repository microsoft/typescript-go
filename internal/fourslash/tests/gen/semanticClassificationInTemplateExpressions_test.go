package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationInTemplateExpressions(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module /*0*/M {
    export class /*1*/C {
        static x;
    }
    export enum /*2*/E {
        E1 = 0
    }
}
` + "`" + `abcd${ /*3*/M./*4*/C.x + /*5*/M./*6*/E.E1}efg` + "`" + ``
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "namespace.declaration", Text: "M"},
		{Type: "class.declaration", Text: "C"},
		{Type: "property.declaration.static", Text: "x"},
		{Type: "enum.declaration", Text: "E"},
		{Type: "enumMember.declaration.readonly", Text: "E1"},
		{Type: "namespace", Text: "M"},
		{Type: "class", Text: "C"},
		{Type: "property.static", Text: "x"},
		{Type: "namespace", Text: "M"},
		{Type: "enum", Text: "E"},
		{Type: "enumMember.readonly", Text: "E1"},
	})
}
