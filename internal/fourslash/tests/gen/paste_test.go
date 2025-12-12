package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPaste(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `fn(/**/);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Paste(t, "x,y,z")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `fn(x, y, z);`)
=======
	f.VerifyCurrentLineContentIs(t, "fn(x, y, z);")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
