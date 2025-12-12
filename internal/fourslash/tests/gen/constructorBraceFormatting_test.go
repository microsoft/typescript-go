package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestConstructorBraceFormatting(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class X {
    constructor () {}/*target*/
 /**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Insert(t, "}")
	f.GoToMarker(t, "target")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `    constructor() { }`)
=======
	f.VerifyCurrentLineContentIs(t, "    constructor() { }")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
