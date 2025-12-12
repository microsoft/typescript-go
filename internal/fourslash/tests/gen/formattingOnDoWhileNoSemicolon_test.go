package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingOnDoWhileNoSemicolon(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*2*/do {
/*3*/    for (var i = 0; i < 10; i++)
/*4*/        i -= 2
/*5*/        }/*1*/while (1 !== 1)`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "\n")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `while (1 !== 1)`)
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContent(t, `do {`)
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContent(t, `    for (var i = 0; i < 10; i++)`)
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContent(t, `        i -= 2`)
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContent(t, `}`)
=======
	f.VerifyCurrentLineContentIs(t, "while (1 !== 1)")
	f.GoToMarker(t, "2")
	f.VerifyCurrentLineContentIs(t, "do {")
	f.GoToMarker(t, "3")
	f.VerifyCurrentLineContentIs(t, "    for (var i = 0; i < 10; i++)")
	f.GoToMarker(t, "4")
	f.VerifyCurrentLineContentIs(t, "        i -= 2")
	f.GoToMarker(t, "5")
	f.VerifyCurrentLineContentIs(t, "}")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
