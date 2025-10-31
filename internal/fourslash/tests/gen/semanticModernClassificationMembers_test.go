package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationMembers(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
  static x = 9;
  f = 9;
  async m() { return A.x + await this.m(); };
  get s() { return this.f; 
  static t() { return new A().f; };
  constructor() {}
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{
		{Type: "class.declaration", Text: "A"},
		{Type: "property.declaration.static", Text: "x"},
		{Type: "property.declaration", Text: "f"},
		{Type: "method.declaration.async", Text: "m"},
		{Type: "class", Text: "A"},
		{Type: "property.static", Text: "x"},
		{Type: "method.async", Text: "m"},
		{Type: "property.declaration", Text: "s"},
		{Type: "property", Text: "f"},
		{Type: "method.declaration.static", Text: "t"},
		{Type: "class", Text: "A"},
		{Type: "property", Text: "f"},
	})
}
