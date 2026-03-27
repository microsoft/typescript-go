package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoMappedType4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type ToGet<T> = T extends string ? ` + "`get${Capitalize<T>}`" + ` : never;
type Getters<T> = { 
    /** @inheritDoc desc on Getters */
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
	f.VerifyQuickInfoAt(t, "3", "(property) getD: () => string", "")
}
