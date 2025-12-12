package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingOnSemiColon(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var  a=b+c^d-e*++f`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToEOF(t)
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, `var a = b + c ^ d - e * ++f;`)
=======
	f.VerifyCurrentFileContentIs(t, "var a = b + c ^ d - e * ++f;")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
