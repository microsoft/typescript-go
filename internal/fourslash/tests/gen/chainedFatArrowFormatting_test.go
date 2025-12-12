package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestChainedFatArrowFormatting(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var fn = () => () => null/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `var fn = () => () => null;`)
=======
	f.VerifyCurrentLineContentIs(t, "var fn = () => () => null;")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
