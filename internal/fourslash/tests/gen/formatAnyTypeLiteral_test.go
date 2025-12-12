package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatAnyTypeLiteral(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(x: { } /*objLit*/){
/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Insert(t, "}")
	f.GoToMarker(t, "objLit")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `function foo(x: {}) {`)
=======
	f.VerifyCurrentLineContentIs(t, "function foo(x: {}) {")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
