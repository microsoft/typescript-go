package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestDocumentSymbolPrivateName(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: first.ts
class A {
  #foo: () => {
    class B {
      #bar: () => {   
         function baz () {
         }
      }
    }
  }
}
// @Filename: second.ts
class Foo {
	#privateMethod() {}
}
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineDocumentSymbols(t)
	f.GoToFile(t, "second.ts")
	f.VerifyBaselineDocumentSymbols(t)
}
