package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferenceToClassVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: referenceToClass_1.ts
class /*1*/foo {
    public n: /*2*/foo;
    public foo: number;
}

class bar {
    public n: /*3*/foo;
    public k = new /*4*/foo();
}

namespace mod {
    var k: /*5*/foo = null;
}
// @Filename: referenceToClass_2.ts
var k: /*6*/foo;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6")
}
