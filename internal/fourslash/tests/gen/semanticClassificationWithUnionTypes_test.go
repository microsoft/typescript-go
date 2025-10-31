package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationWithUnionTypes(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module /*0*/M {
    export interface /*1*/I {
    }
}

interface /*2*/I {
}
class /*3*/C {
}

var M: /*4*/M./*5*/I | /*6*/I | /*7*/C;
var I: typeof M | typeof /*8*/C;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "variable", Text: "M"}, {Type: "interface.declaration", Text: "I"}, {Type: "interface.declaration", Text: "I"}, {Type: "class.declaration", Text: "C"}, {Type: "variable.declaration", Text: "M"}, {Type: "variable", Text: "M"}, {Type: "interface", Text: "I"}, {Type: "interface", Text: "I"}, {Type: "class", Text: "C"}, {Type: "class.declaration", Text: "I"}, {Type: "variable", Text: "M"}, {Type: "class", Text: "C"}})
}
