package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatDocumentChangesNothing(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const x = 1;
const y = 2;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.FormatDocumentChangesNothing(t)
}
