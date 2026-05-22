package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTransitiveExportImportsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: commonjs
// @Filename: a.ts
[|class /*1*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}A|] {
}|]
[|export = [|{| "contextRangeIndex": 2 |}A|];|]
// @Filename: b.ts
[|export import /*2*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 4 |}b|] = require('./a');|]
// @Filename: c.ts
[|import /*3*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}b|] = require('./b');|]
var a = new /*4*/[|b|]./**/[|b|]();`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4")
}
