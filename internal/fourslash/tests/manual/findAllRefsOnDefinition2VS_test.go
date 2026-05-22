package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsOnDefinition2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: findAllRefsOnDefinition2-import.ts
export module Test{

    /*1*/export interface /*2*/start { }

    export interface stop { }
}
//@Filename: findAllRefsOnDefinition2.ts
import Second = require("./findAllRefsOnDefinition2-import");

var start: Second.Test./*3*/start;
var stop: Second.Test.stop;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3")
}
