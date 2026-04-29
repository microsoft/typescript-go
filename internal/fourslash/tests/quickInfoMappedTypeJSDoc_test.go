package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoMappedTypeJSDoc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	// https://github.com/microsoft/typescript-go/issues/3659
	const content = `
type Test = {
    /** a's comment */
    a: string;
};

type Mapped = {
    [K in keyof Test]: number;
};

const x: Mapped = {
    a: 123
};

x./*1*/a
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "(property) a: number", "a's comment")
}
