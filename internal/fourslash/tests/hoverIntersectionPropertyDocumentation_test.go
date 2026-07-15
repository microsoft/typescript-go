package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestHoverIntersectionPropertyDocumentation(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
interface A {
    /** interface A */
    prop: () => number;
}

interface B {
    /** interface B */
    prop: () => number;
}

type C = A & B;
declare const c: C;

c./*1*/prop;
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "(property) prop: () => number", "interface A\ninterface B")
}
