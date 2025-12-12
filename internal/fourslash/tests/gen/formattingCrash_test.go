package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingCrash(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**/module Default{ 
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "PlaceOpenBraceOnNewLineForFunctions", true)
	f.SetFormatOption(t, "PlaceOpenBraceOnNewLineForControlBlocks", true)
	f.FormatDocument(t, "")
	f.GoToMarker(t, "")
	f.VerifyCurrentLineContent(t, `module Default`)
}
