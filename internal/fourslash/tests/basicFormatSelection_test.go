package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestBasicFormatSelection(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const x = 1;
[|function foo(a,b){return a+b;}|]
const y = 2;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// Format only the function declaration
	rangeToFormat := f.Ranges()[0]
	f.VerifyFormatSelection(t, rangeToFormat, &lsproto.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	})
}
