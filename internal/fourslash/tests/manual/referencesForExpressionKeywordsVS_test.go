package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesForExpressionKeywordsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class C {
    static x = 1;
}
/*new*/new C();
/*void*/void C;
/*typeof*/typeof C;
/*delete*/delete C.x;
/*async*/async function* f() {
    /*yield*/yield C;
    /*await*/await C;
}
"x" /*in*/in C;
undefined /*instanceof*/instanceof C;
undefined /*as*/as C;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "new", "void", "typeof", "yield", "await", "in", "instanceof", "as", "delete")
}
