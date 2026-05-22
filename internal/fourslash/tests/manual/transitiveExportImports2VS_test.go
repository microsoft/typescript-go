package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTransitiveExportImports2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
[|namespace /*A*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}A|] {
    export const x = 0;
}|]
// @Filename: b.ts
[|export import /*B*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}B|] = [|A|];|]
[|B|].x;
// @Filename: c.ts
[|import { /*C*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}B|] } from "./b";|]`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "A", "B", "C")
}
