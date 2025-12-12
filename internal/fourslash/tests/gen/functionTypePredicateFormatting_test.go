package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFunctionTypePredicateFormatting(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**/function bar(a: A):     a        is       B    {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.FormatDocument(t, "")
	f.VerifyCurrentLineContent(t, `function bar(a: A): a is B { }`)
}
