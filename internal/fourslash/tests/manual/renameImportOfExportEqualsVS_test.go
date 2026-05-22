package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRenameImportOfExportEqualsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[|declare namespace /*N*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 0 |}N|] {
    [|export var /*x*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 2 |}x|]: number;|]
}|]
declare module "mod" {
    [|export = [|{| "contextRangeIndex": 4 |}N|];|]
}
declare module "a" {
    [|import * as /*a*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 6 |}N|] from "mod";|]
    [|export { [|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 8 |}N|] };|] // Renaming N here would rename
}
declare module "b" {
    [|import { /*b*/[|{| "isWriteAccess": true, "isDefinition": true, "contextRangeIndex": 10 |}N|] } from "a";|]
    export const y: typeof [|N|].[|x|];
}`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "N", "a", "b", "x")
}
