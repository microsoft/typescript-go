package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestParserCorruptionAfterMapInClass(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @target: esnext
// @lib: es2015
// @strict: true
class C {
    map = new Set<string, number>/*$*/

    foo() {

    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToMarker(t, "$")
	f.Insert(t, "()")
	f.VerifyDiagnostics(t, nil)
}
