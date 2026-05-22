package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTransitiveExportImports3VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
[|export function /*f*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}f|]() {}|]
// @Filename: b.ts
[|export { [|{| "contextRangeIndex": 2 |}f|] as /*g0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}g|] } from "./a";|]
[|import { /*f2*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 5 |}f|] } from "./a";|]
[|import { /*g1*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 7 |}g|] } from "./b";|]`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "f", "g0", "g1", "f2")
}
