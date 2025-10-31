package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationCallableVariables(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A { onEvent: () => void; }
const x = new A().onEvent;
const match = (s: any) => x();
const other = match;
match({ other });
interface B = { (): string; }; var b: B
var s: String;
var t: { (): string; foo: string};`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "class.declaration", Text: "A"}, {Type: "method.declaration", Text: "onEvent"}, {Type: "function.declaration.readonly", Text: "x"}, {Type: "class", Text: "A"}, {Type: "method", Text: "onEvent"}, {Type: "function.declaration.readonly", Text: "match"}, {Type: "parameter.declaration", Text: "s"}, {Type: "function.readonly", Text: "x"}, {Type: "function.declaration.readonly", Text: "other"}, {Type: "function.readonly", Text: "match"}, {Type: "function.readonly", Text: "match"}, {Type: "method.declaration", Text: "other"}, {Type: "interface.declaration", Text: "B"}, {Type: "variable.declaration", Text: "b"}, {Type: "interface", Text: "B"}, {Type: "variable.declaration", Text: "s"}, {Type: "interface.defaultLibrary", Text: "String"}, {Type: "variable.declaration", Text: "t"}, {Type: "property.declaration", Text: "foo"}})
}
