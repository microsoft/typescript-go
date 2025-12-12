package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingInComment(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()

=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
foo(              ); // /*1*/
}
function foo() {       var x;       } // /*2*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `foo(              ); // ;`)
	f.GoToMarker(t, "2")
	f.Insert(t, "}")
	f.VerifyCurrentLineContent(t, `function foo() {       var x;       } // }`)
=======
	f.VerifyCurrentLineContentIs(t, "foo(              ); // ;")
	f.GoToMarker(t, "2")
	f.Insert(t, "}")
	f.VerifyCurrentLineContentIs(t, "function foo() {       var x;       } // }")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
