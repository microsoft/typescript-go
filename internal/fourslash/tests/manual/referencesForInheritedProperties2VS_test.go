package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesForInheritedProperties2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface interface1 {
    /*1*/doStuff(): void;
}

interface interface2 {
    doStuff(): void;
}

interface interface2 extends interface1 {
}

class class1 implements interface2 {
    doStuff() {

    }
}

class class2 extends class1 {

}

var v: class2;
v.doStuff();`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1")
}
