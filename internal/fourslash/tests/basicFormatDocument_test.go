package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestBasicFormatDocument(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const   x    =     1   ;
function      foo  (  a  ,   b  )   {
return    a  +   b  ;
}
const  y =   foo(  2  ,   3  )  ;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyFormatDocument(t, &lsproto.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	})
}
