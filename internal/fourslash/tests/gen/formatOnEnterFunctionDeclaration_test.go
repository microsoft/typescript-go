package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatOnEnterFunctionDeclaration(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*0*/function listAPIFiles(path: string): string[] {/*1*/ }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.InsertLine(t, "")
	f.GoToMarker(t, "0")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `function listAPIFiles(path: string): string[] {`)
=======
	f.VerifyCurrentLineContentIs(t, "function listAPIFiles(path: string): string[] {")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
