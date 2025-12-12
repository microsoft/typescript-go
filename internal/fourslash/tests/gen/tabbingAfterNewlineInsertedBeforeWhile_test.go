package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTabbingAfterNewlineInsertedBeforeWhile(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo() {
    /**/while (true) { }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.InsertLine(t, "")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `    while (true) { }`)
=======
	f.VerifyCurrentLineContentIs(t, "    while (true) { }")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
