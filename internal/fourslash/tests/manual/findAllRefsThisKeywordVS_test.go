package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsThisKeywordVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noLib: true
/*1*/this;
function f(/*2*/this) {
    return /*3*/this;
    function g(/*4*/this) { return /*5*/this; }
}
class C {
    static x() {
        /*6*/this;
    }
    static y() {
        () => /*7*/this;
    }
    constructor() {
        /*8*/this;
    }
    method() {
        () => /*9*/this;
    }
}
// These are *not* real uses of the 'this' keyword, they are identifiers.
const x = { /*10*/this: 0 }
x./*11*/this;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11")
}
