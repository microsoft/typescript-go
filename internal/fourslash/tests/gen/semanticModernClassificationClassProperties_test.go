package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationClassProperties(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A { 
  private y: number;
  constructor(public x : number, _y : number) { this.y = _y; }
  get z() : number { return this.x + this.y; }
  set a(v: number) { }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "class.declaration", Text: "A"}, {Type: "property.declaration", Text: "y"}, {Type: "parameter.declaration", Text: "x"}, {Type: "parameter.declaration", Text: "_y"}, {Type: "property", Text: "y"}, {Type: "parameter", Text: "_y"}, {Type: "property.declaration", Text: "z"}, {Type: "property", Text: "x"}, {Type: "property", Text: "y"}, {Type: "property.declaration", Text: "a"}, {Type: "parameter.declaration", Text: "v"}})
}
