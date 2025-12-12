package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatIfTryCatchBlocks(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `try {
}
catch {
}

try {
}
catch (e) {
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "PlaceOpenBraceOnNewLineForControlBlocks", true)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `try
{
}
catch
{
}

try
{
}
catch (e)
{
}`)
}
