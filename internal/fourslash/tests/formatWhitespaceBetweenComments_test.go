package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatWhitespaceBetweenComments(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
  const x = "wont format"
//
 
//
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
const x = "wont format"
//

//
`)
}

func TestFormatWhitespaceBetweenCommentsPreserveTrailingWhitespace(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
  const x = "wont format"
//
 
//
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	opts := f.GetOptions()
	opts.FormatCodeSettings.TrimTrailingWhitespace = core.TSFalse
	f.Configure(t, opts)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
const x = "wont format"
//
 
//
`)
}
