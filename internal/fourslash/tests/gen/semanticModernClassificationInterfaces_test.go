package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationInterfaces(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Pos { x: number, y: number };
const p = { x: 1, y: 2 } as Pos;
const foo = (o: Pos) => o.x + o.y;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "interface.declaration", Text: "Pos"}, {Type: "property.declaration", Text: "x"}, {Type: "property.declaration", Text: "y"}, {Type: "variable.declaration.readonly", Text: "p"}, {Type: "property.declaration", Text: "x"}, {Type: "property.declaration", Text: "y"}, {Type: "interface", Text: "Pos"}, {Type: "function.declaration.readonly", Text: "foo"}, {Type: "parameter.declaration", Text: "o"}, {Type: "interface", Text: "Pos"}, {Type: "parameter", Text: "o"}, {Type: "property", Text: "x"}, {Type: "parameter", Text: "o"}, {Type: "property", Text: "y"}})
}
