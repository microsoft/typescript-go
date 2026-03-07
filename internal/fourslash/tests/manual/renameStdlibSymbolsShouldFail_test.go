package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRenameStdlibSymbolsShouldFail(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const x: /*1*/Array<string> = [];
const y: /*2*/Promise<void> = Promise.resolve();
/*3*/console.log("hello");
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Standard library symbols should not be renameable
	markers := []string{"1", "2", "3"}
	for _, marker := range markers {
		f.GoToMarker(t, marker)
		f.VerifyRenameFailed(t, nil /*preferences*/)
	}
}
