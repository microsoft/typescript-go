package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickinfoVerbosityClassWithMixinBase1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	// Expanded hover serializes Derived's generated class declaration.
	// Not every base type was serializable (e.g. Mixin did not produce a Node in the `extends` clause).
	// That would end up creating an `extends` clause with nil in place of an actual base node,
	// which would cause issues when we'd print out the constructed node.
	const content = `
class Base {}

declare const Mixin: new () => Base & { mixed: string };

class Derived/*1*/ extends Mixin {}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyBaselineHoverWithVerbosity(t, map[string][]int{"1": {0, 1}})
}
