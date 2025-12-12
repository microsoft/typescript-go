package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingObjectLiteralOpenCurlyNewlineAssignment(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
var obj = {};
obj =
{
    prop: 3
};
 
var obj2 = obj ||
{
    prop: 0
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
var obj = {};
obj =
{
    prop: 3
};

var obj2 = obj ||
{
    prop: 0
}
`)
	f.SetFormatOption(t, "indentMultiLineObjectLiteralBeginningOnBlankLine", true)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
var obj = {};
obj =
    {
        prop: 3
    };

var obj2 = obj ||
    {
        prop: 0
    }
`)
}
