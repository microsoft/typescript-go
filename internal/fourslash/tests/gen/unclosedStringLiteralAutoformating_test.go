package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestUnclosedStringLiteralAutoformating(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var x = /*1*/"asd/*2*/
class Foo {
    /**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Insert(t, "}")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `}`)
=======
	f.VerifyCurrentLineContentIs(t, "}")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
