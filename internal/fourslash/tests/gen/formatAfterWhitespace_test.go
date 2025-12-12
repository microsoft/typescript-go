package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatAfterWhitespace(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo()
{
    var bar;
    /*1*/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.InsertLine(t, "")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, `function foo()
{
    var bar;


}`)
=======
	f.VerifyCurrentFileContentIs(t, "function foo()\n{\n    var bar;\n\n\n}")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
