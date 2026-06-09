package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToDefinitionMappedType2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Foo {
    /*def*/property: string
}

type JustMapIt<T> = {[P in keyof T]: 0}
type MapItWithRemap<T> = {[P in keyof T as P extends string ? ` + "`mapped_${P}`" + ` : never]: 0}

{
    let gotoDef!: JustMapIt<Foo>
    gotoDef.property
}

{
    let gotoDef!: MapItWithRemap<Foo>
    gotoDef.[|/*ref*/mapped_property|]
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, true, "ref")
}
