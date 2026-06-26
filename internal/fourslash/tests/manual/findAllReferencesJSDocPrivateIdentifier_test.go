package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllReferencesJSDocPrivateIdentifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
class C {
    /*1*/#field;
    /**
     * @see C./*2*/#field
     */
    method() {}
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyBaselineFindAllReferences(t, "1", "2")
}
