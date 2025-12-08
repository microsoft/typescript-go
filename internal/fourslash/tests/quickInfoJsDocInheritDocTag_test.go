package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoJsDocInheritDocTag(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noEmit: true
// @allowJs: true
// @Filename: quickInfoJsDocInheritDocTag.js
abstract class A {
  /**
   * A.f description
   * @returns {string} A.f return value.
   */
  public static f(props?: any): string {
    throw new Error("Must be implemented by subclass");
  }
}

class B extends A {
  /**
   * B.f description
   * @inheritDoc
   * @param {{ a: string; b: string; }} [props] description of props
   */
  public static /**/f(props?: { a: string; b: string }): string {}
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineHover(t)
}
