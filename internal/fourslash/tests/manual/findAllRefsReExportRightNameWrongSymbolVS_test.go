package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsReExportRightNameWrongSymbolVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
[|export const /*a*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}x|] = 0;|]
// @Filename: /b.ts
[|export const /*b*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}x|] = 0;|]
//@Filename: /c.ts
[|export { /*cFromB*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 4 |}x|] } from "./b";|]
[|import { /*cFromA*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}x|] } from "./a";|]
/*cUse*/[|x|];
// @Filename: /d.ts
[|import { /*d*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 9 |}x|] } from "./c";|]`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "a", "b", "cFromB", "cFromA", "cUse", "d")
}
