package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatDocumentInPlace(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/for (;;) { }
/*2*/for (var x;x<0;x++) { }
/*3*/for (var x ;x<0 ;x++) { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.FormatDocument(t)
	// After formatting, verify the formatted content
	f.CurrentFileContentIs(t, `for (; ;) { }
for (var x; x < 0; x++) { }
for (var x; x < 0; x++) { } `)
}
