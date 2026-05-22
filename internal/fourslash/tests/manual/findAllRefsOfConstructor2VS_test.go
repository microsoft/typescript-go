package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsOfConstructor2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
    /*a*/constructor(s: string) {}
}
class B extends A {
    /*b*/constructor() { super(""); }
}
class C extends B {
    /*c*/constructor() {
        super();
    }
}
class D extends B { }
const a = new A("a");
const b = new B();
const c = new C();
const d = new D();`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "a", "b", "c")
}
