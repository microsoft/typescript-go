package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestHoverQualifiedGenericNames(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
function makeBox<T>(x: T) {
    class Box {
        value = x
    }
    return new Box()
}

let box/**/ = makeBox("hello")
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "")
	f.VerifyQuickInfoIs(t, "let box: makeBox<string>.Box", "")
}
