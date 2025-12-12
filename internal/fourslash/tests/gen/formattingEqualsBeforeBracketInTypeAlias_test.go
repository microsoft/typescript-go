package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingEqualsBeforeBracketInTypeAlias(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type X    =     [number]/*1*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `type X = [number];`)
=======
	f.VerifyCurrentLineContentIs(t, "type X = [number];")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
