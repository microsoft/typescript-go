package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationModules(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module /*0*/M {
    export var v;
    export interface /*1*/I {
    }
}

var x: /*2*/M./*3*/I = /*4*/M.v;
var y = /*5*/M;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "namespace.declaration", Text: "M"}, {Type: "variable.declaration.local", Text: "v"}, {Type: "interface.declaration", Text: "I"}, {Type: "variable.declaration", Text: "x"}, {Type: "namespace", Text: "M"}, {Type: "interface", Text: "I"}, {Type: "namespace", Text: "M"}, {Type: "variable.local", Text: "v"}, {Type: "variable.declaration", Text: "y"}, {Type: "namespace", Text: "M"}})
}
