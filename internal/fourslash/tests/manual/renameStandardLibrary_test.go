package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRenameStandardLibrary(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// Test that standard library symbols cannot be renamed
/*1*/setTimeout(() => {}, 100);
/*2*/console.log("test");
const arr = [1, 2, 3];
arr./*3*/push(4);
const str = "test";
str./*4*/substring(0, 1);
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Try to rename standard library functions - should fail
	markers := []string{"1", "2", "3", "4"}
	for _, marker := range markers {
		f.GoToMarker(t, marker)
		f.VerifyRenameFailed(t, nil /*preferences*/)
	}
}
