package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestBasicFormatOnType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo()  {/*a*/
const x=1;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// Verify formatting after typing opening curly brace
	f.VerifyFormatOnType(t, "a", "{", &lsproto.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	})
}
