package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsForDefaultExport03VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/function /*2*/f() {
    return 100;
}

/*3*/export default /*4*/f;

var x: typeof /*5*/f;

var y = /*6*/f();

/*7*/namespace /*8*/f {
    var local = 100;
}`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8")
}
