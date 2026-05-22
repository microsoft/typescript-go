package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRenameImportOfReExport2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module "a" {
    [|export class /*1*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}C|] {}|]
}
declare module "b" {
    [|export { [|{| "contextRangeIndex": 2 |}C|] as /*2*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}D|] } from "a";|]
}
declare module "c" {
    [|import { /*3*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 5 |}D|] } from "b";|]
    export function f(c: [|D|]): void;
}`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3")
}
