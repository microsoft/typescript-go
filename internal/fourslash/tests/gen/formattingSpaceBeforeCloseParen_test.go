package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingSpaceBeforeCloseParen(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/({});
/*2*/(  {});
/*3*/({foo:42});
/*4*/(  {foo:42}  );
/*5*/var bar = (function (a) { });`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis", true)
	f.FormatDocument(t, "")
	f.GoToMarker(t, "1")
	f.VerifyCurrentLineContent(t, `( {} );`)
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContent(t, `( {} );`)
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContent(t, `( { foo: 42 } );`)
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContent(t, `( { foo: 42 } );`)
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContent(t, `var bar = ( function( a ) { } );`)
	f.SetFormatOption(t, "InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis", false)
	f.FormatDocument(t, "")
	f.GoToMarker(t, "1")
	f.VerifyCurrentLineContent(t, `({});`)
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContent(t, `({});`)
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContent(t, `({ foo: 42 });`)
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContent(t, `({ foo: 42 });`)
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContent(t, `var bar = (function(a) { });`)
}
