package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToDefinitionMappedType3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Source {
  /*def*/alpha: number;
  beta: string;
}

// Transforming interface field names with a suffix
type Transformed<T> = {
  [K in keyof T as ` + "`${K & string}Suffix`" + `]: () => T[K];
};

type Result = Transformed<Source>;
/*
  Expected:
  {
    alphaSuffix: () => number;
    betaSuffix: () => string;
  }
  */

const obj: Result = {
  alphaSuffix: () => 42,
  betaSuffix: () => "hello",
};

obj.[|/*ref*/alphaSuffix|]();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, true, "ref")
}
