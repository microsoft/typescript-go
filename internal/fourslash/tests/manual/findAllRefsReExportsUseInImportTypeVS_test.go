package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsReExportsUseInImportTypeVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /foo/types/types.ts
[|export type /*full0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}Full|] = { prop: string; };|]
// @Filename: /foo/types/index.ts
[|import * as /*foo0*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}foo|] from './types';|]
[|export { /*foo1*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 4 |}foo|] };|]
// @Filename: /app.ts
[|import { /*foo2*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}foo|] } from './foo/types';|]
export type fullType = /*foo3*/[|foo|]./*full1*/[|Full|];
type namespaceImport = typeof import('./foo/types');
type fullType2 = import('./foo/types')./*foo4*/[|foo|]./*full2*/[|Full|];`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "full0", "full1", "full2", "foo0", "foo1", "foo2", "foo3", "foo4")
}
