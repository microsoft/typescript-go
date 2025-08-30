package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoOnVariableDeclaration(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noEmit: true
// @allowJS: true
// @checkJs: true
// @filename: /a.js

/** @type {number} */
const /*1*/x, /*2*/y;

try {} catch(/*3*/error) {}`

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyQuickInfoAt(t, "1", "const x: number", "")
	f.VerifyQuickInfoAt(t, "2", "const y: any", "")
	f.VerifyQuickInfoAt(t, "3", "var error: any", "")
}
