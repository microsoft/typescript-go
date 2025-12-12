package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingOnNestedDoWhileByEnter(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*2*/do{
/*3*/do/*1*/{
/*4*/do{
/*5*/}while(a!==b)
/*6*/}while(a!==b)
/*7*/}while(a!==b)`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "\n")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `    {`)
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContent(t, `do{`)
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContent(t, `    do`)
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContent(t, `do{`)
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContent(t, `}while(a!==b)`)
	f.GoToMarker(t, "6")
	f.VerifyCurrentLineContent(t, `}while(a!==b)`)
	f.GoToMarker(t, "7")
	f.VerifyCurrentLineContent(t, `}while(a!==b)`)
=======
	f.VerifyCurrentLineContentIs(t, "    {")
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContentIs(t, "do{")
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContentIs(t, "    do")
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContentIs(t, "do{")
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContentIs(t, "}while(a!==b)")
	f.GoToMarker(t, "6")
	f.VerifyCurrentLineContentIs(t, "}while(a!==b)")
	f.GoToMarker(t, "7")
	f.VerifyCurrentLineContentIs(t, "}while(a!==b)")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
