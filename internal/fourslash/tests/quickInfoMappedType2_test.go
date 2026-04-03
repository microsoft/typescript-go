package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoMappedType2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Tests that @inheritDoc on a MappedTypeNode causes hover to show documentation.
	// TODO: once @inheritDoc resolution is fully implemented, the expected documentation
	// should be "desc on Getters\nhello" (combined from the mapped type and the source property).
	const content = `type ToGet<T> = T extends string ? ` + "`get${Capitalize<T>}`" + ` : never;
type Getters<T> = /** @inheritDoc desc on Getters */ { 
    [P in keyof T as ToGet<P>]: () => T[P]
};

type Y = {
    /** hello  */
    d: string;
}

type T50 = Getters<Y>;  // { getFoo: () => string, getBar: () => number }

declare let y: T50;
y.get/*3*/D;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	// Current Go behavior: shows the raw @inheritDoc tag text from the mapped type declaration.
	f.VerifyQuickInfoAt(t, "3", "(property) getD: () => string", "\n\n*@inheritDoc* \u2014 desc on Getters ")
}
