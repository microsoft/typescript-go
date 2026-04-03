package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoMappedType3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type Getters<Type> =  /** @inheritDoc desc on Getters */  {
  [Property in keyof Type as ` + "`" + `get${Capitalize<
    string & Property
  >}` + "`" + `]: () => Type[Property];
};

interface Person {
  // ✅ When hovering here, the documentation is displayed, as it should.
  /**
   * Person's name.
   * @example "John Doe"
   */
  name: string;

  // ✅ When hovering here, the documentation is displayed, as it should.
  /**
   * Person's Age.
   * @example 30
   */
  age: number;

  // ✅ When hovering here, the documentation is displayed, as it should.
  /**
   * Person's Location.
   * @example "Brazil"
   */
  location: string;
}

type LazyPerson = Getters<Person>;

const me: LazyPerson = {
  // ❌ When hovering here, the documentation is NOT displayed.
  /*1*/getName: () => "Jake Carter",
  // ❌ When hovering here, the documentation is NOT displayed.
  /*2*/getAge: () => 35,
  // ❌ When hovering here, the documentation is NOT displayed.
  /*3*/getLocation: () => "United States",
};

// ❌ When hovering here, the documentation is NOT displayed.
me./*4*/getName();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineHover(t)
}
