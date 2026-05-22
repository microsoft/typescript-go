package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsNonModuleVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @checkJs: true
// @Filename: /script.ts
console.log("I'm a script!");
// @Filename: /import.ts
import "./script/*1*/";
// @Filename: /require.js
require("./script/*2*/");
console.log("./script/*3*/");
// @Filename: /tripleSlash.ts
/// <reference path="script.ts" />
// @Filename: /stringLiteral.ts
console.log("./script");`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3")
}
