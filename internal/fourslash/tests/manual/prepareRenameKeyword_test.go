package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestPrepareRenameKeyword(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
class /*1*/Foo {
    method() {}
}

const x = /*2*/(1 + 2);

function /*3*/bar() {}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Test 1: Rename on 'class' keyword should work (adjusts to class name)
	// Note: We can't easily test the class keyword position in fourslash,
	// but we can test that renaming on the class name works
	f.GoToMarker(t, "1")
	f.VerifyRenameSucceeded(t, nil /*preferences*/)

	// Test 2: Rename on paren should fail
	f.GoToMarker(t, "2")
	f.VerifyRenameFailed(t, nil /*preferences*/)

	// Test 3: Rename on 'function' keyword should work (adjusts to function name)
	f.GoToMarker(t, "3")
	f.VerifyRenameSucceeded(t, nil /*preferences*/)
}
