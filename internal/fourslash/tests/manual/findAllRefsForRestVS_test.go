package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsForRestVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Gen {
    x: number
    /*1*/parent: Gen;
    millennial: string;
}
let t: Gen;
var { x, ...rest } = t;
rest./*2*/parent;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2")
}
