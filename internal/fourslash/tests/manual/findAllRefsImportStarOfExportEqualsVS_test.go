package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsImportStarOfExportEqualsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowSyntheticDefaultimports: true
// @Filename: /node_modules/a/index.d.ts
[|declare function /*a0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}a|](): void;|]
[|declare namespace /*a1*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}a|] {
    export const x: number;
}|]
[|export = /*a2*/[|{| "contextRangeIndex": 4 |}a|];|]
// @Filename: /b.ts
[|import /*b0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}b|] from "a";|]
/*b1*/[|b|]();
[|b|].x;
// @Filename: /c.ts
[|import /*c0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 10 |}a|] from "a";|]
/*c1*/[|a|]();
/*c2*/[|a|].x;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "a0", "a1", "a2", "b0", "b1", "c0", "c1", "c2")
}
