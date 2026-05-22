package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsImportEqualsJsonFileVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @checkJs: true
// @resolveJsonModule: true
// @module: commonjs
// @Filename: /a.ts
import /*0*/j = require("/*1*/./j.json");
/*2*/j;
// @Filename: /b.js
const /*3*/j = require("/*4*/./j.json");
/*5*/j;
// @Filename: /j.json
/*6*/{ "x": 0 }`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "0", "2", "1", "4", "3", "5", "6")
}
