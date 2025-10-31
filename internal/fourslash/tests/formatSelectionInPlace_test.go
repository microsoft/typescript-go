package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatSelectionInPlace(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const x = 1;
/*1*/function foo(a,b){return a+b;}/*2*/
const y = 2;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.FormatSelection(t, "1", "2")
	f.CurrentFileContentIs(t, `const x = 1;
function foo(a, b) { return a + b; }
const y = 2;`)
}
