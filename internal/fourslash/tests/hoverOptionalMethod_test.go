package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestHoverOptionalMethod(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @strict: true
interface I {
    bar?(): void
}
type T = {
    bar?(): void
}
class C {
    baz?() {}
}

declare const i: I
declare const t: T
declare const c: C

i./*1*/bar
t./*2*/bar
c./*3*/baz
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Pre-fix (#3607): hover at /*3*/ returned an empty quickinfo because the optional
	// method's type was `(() => void) | undefined`, which has no call signatures via
	// getSignaturesOfType. Fall-through to `symbol.Declarations` recovers the signature.
	f.VerifyQuickInfoAt(t, "1", "(method) I.bar(): void", "")
	f.VerifyQuickInfoAt(t, "2", "(method) bar(): void", "")
	f.VerifyQuickInfoAt(t, "3", "(method) C.baz(): void", "")
}
