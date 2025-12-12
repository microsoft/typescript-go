package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingOnEnter(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class foo { }
class bar {/**/ }
// new line here`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.InsertLine(t, "")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, `class foo { }
class bar {
}
// new line here`)
=======
	f.VerifyCurrentFileContentIs(t, "class foo { }\nclass bar {\n}\n// new line here")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
