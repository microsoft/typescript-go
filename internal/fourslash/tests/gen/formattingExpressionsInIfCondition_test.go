package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingExpressionsInIfCondition(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if (a === 1 ||
    /*0*/b === 2 ||/*1*/
    c === 3) {
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "\n")
	f.GoToMarker(t, "0")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `    b === 2 ||`)
=======
	f.VerifyCurrentLineContentIs(t, "    b === 2 ||")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
