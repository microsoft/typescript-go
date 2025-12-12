package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatControlFlowConstructs(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if (true)/**/
{     
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.GoToMarker(t, "")
	f.VerifyCurrentLineContent(t, `if (true) {`)
}
