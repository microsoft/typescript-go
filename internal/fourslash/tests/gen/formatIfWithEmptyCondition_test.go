package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatIfWithEmptyCondition(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if () {
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "PlaceOpenBraceOnNewLineForControlBlocks", true)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `if ()
{
}`)
}
