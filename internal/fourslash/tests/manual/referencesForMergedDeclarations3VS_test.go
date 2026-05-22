package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesForMergedDeclarations3VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[|class /*class*/[|testClass|] {
    static staticMethod() { }
    method() { }
}|]

[|module /*module*/[|testClass|] {
    export interface Bar {

    }
}|]

var c1: [|testClass|];
var c2: [|testClass|].Bar;
[|testClass|].staticMethod();
[|testClass|].prototype.method();
[|testClass|].bind(this);
new [|testClass|]();`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "module", "class")
}
