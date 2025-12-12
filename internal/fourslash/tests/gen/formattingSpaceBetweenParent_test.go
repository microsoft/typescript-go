package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingSpaceBetweenParent(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/foo(() => 1);
/*2*/foo(1);
/*3*/if((true)){}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis", true)
	f.FormatDocument(t, "")
	f.GoToMarker(t, "1")
	f.VerifyCurrentLineContent(t, `foo( () => 1 );`)
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContent(t, `foo( 1 );`)
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContent(t, `if ( ( true ) ) { }`)
}
