package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesForMergedDeclarationsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/interface /*2*/Foo {
}

/*3*/module /*4*/Foo {
    export interface Bar { }
}

/*5*/function /*6*/Foo(): void {
}

var f1: /*7*/Foo.Bar;
var f2: /*8*/Foo;
/*9*/Foo.bind(this);`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8", "9")
}
