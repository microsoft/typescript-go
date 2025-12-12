package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatDebuggerStatement(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if(false){debugger;}
  if    (   false   )   {    debugger  ;   }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.GoToBOF(t)
	f.VerifyCurrentLineContent(t, `if (false) { debugger; }`)
	f.GoToEOF(t)
	f.VerifyCurrentLineContent(t, `if (false) { debugger; }`)
}
