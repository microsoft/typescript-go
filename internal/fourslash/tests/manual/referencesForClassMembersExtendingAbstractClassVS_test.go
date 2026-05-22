package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesForClassMembersExtendingAbstractClassVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `abstract class Base {
    abstract /*a1*/a: number;
    abstract /*method1*/method(): void;
}
class MyClass extends Base {
    /*a2*/a;
    /*method2*/method() { }
}

var c: MyClass;
c./*a3*/a;
c./*method3*/method();`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "a1", "a2", "a3", "method1", "method2", "method3")
}
