package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCodeFixEnumComputedPropertyName(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = "// @Filename: test.ts\n" +
		"enum CHAR {\n" +
		"    ['\\t']/*1*/ = 0x09,\n" +
		"}"
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.VerifyImportFixAtPosition(t, []string{
		"enum CHAR {\n" +
			"    \"\\t\" = 0x09,\n" +
			"}",
	}, nil /*preferences*/)
}
