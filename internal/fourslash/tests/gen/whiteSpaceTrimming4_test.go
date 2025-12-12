package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestWhiteSpaceTrimming4(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var re = /\w+   /*1*//;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "\n")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, "var re = /\\w+\n    /;")
=======
	f.VerifyCurrentFileContentIs(t, "var re = /\\w+\n    /;")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
