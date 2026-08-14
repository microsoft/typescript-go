package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGetContextNodeForOfForInSwitch(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const arr = [1, 2, 3];
for (const /*a*/x of arr) {
    /*aUse*/x;
}
const obj = { a: 1 };
for (const /*b*/y in obj) {
    /*bUse*/y;
}
function f(n: number) {
    switch (n) {
        case 0:
            const /*c*/z = 1;
            /*cUse*/z;
            break;
    }
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindAllReferences(t, "a", "aUse", "b", "bUse", "c", "cUse")
	f.VerifyBaselineGoToDefinition(t, false, "aUse", "bUse", "cUse")
}
