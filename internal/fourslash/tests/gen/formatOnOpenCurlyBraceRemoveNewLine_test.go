package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatOnOpenCurlyBraceRemoveNewLine(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if(true)
/**/ }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "PlaceOpenBraceOnNewLineForControlBlocks", false)
	f.GoToMarker(t, "")
	f.Insert(t, "{")
	f.VerifyCurrentFileContent(t, `if (true) { }`)
}
