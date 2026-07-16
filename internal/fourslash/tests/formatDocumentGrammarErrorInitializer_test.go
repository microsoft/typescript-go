package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// A malformed intersection type with an unclosed object type (`x as A & { a: { b: C }`) makes
// the parser swallow the following statements into the still-open type literal as recovered
// members. One becomes a PropertySignature whose (grammar-error) initializer is an array literal
// containing a template literal. isGrammarError skips that initializer, so the formatting scanner
// must be advanced past it; otherwise it stays parked on the template backtick and mis-scans it,
// failing the "token end is child end" span invariant on a later sibling.
func TestFormatDocumentGrammarErrorInitializer(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := `// @Filename: /a.ts
const f = () => {
  const v = x as A & {
    a: { b: C
  }
  const m: T[] = [
    { g: () => { nav(` + "`${z}`" + `) } },
  ]
  const n: T[] =
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "/a.ts")
}
