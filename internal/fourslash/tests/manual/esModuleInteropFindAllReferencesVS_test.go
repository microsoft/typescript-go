package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestEsModuleInteropFindAllReferencesVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @esModuleInterop: true
// @Filename: /abc.d.ts
declare module "a" {
    /*1*/export const /*2*/x: number;
}
// @Filename: /b.ts
import a from "a";
a./*3*/x;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3")
}
