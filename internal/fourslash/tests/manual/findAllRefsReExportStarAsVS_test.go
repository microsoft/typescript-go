package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsReExportStarAsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /leafModule.ts
export const /*helloDef*/hello = () => 'Hello';
// @Filename: /exporting.ts
export * as /*leafDef*/Leaf from './leafModule';
// @Filename: /importing.ts
 import { /*leafImportDef*/Leaf } from './exporting';
 /*leafUse*/[|Leaf|]./*helloUse*/[|hello|]()`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "helloDef", "helloUse", "leafDef", "leafImportDef", "leafUse")
}
