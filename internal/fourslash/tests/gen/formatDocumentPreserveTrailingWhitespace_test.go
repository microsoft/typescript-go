package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatDocumentPreserveTrailingWhitespace(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
var a;     
var b     
     
//     
function b(){     
    while(true){     
    }     
}     
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "trimTrailingWhitespace", false)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
var a;     
var b     
     
//     
function b() {     
    while (true) {     
    }     
}     
`)
}
