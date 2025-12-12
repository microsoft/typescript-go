package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingIfInElseBlock(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if (true) {
}
else {
    if (true) {
        /*1*/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "}")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `    }`)
=======
	f.VerifyCurrentLineContentIs(t, "    }")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
