package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestWhiteSpaceTrimming3(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let t = "foo \
bar     \   
"/*1*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, "let t = \"foo \\\nbar     \\   \n\";")
=======
	f.VerifyCurrentFileContentIs(t, "let t = \"foo \\\nbar     \\   \n\";")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
