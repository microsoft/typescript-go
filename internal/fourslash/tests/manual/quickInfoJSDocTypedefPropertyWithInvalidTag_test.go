package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoJSDocTypedefPropertyWithInvalidTag(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: /a.js
/**
 * @typedef {Object} MyType1
 * @property {string} name
 * @-rule
 * @property {number} age
 */

/**
 * @typedef {Object} MyType2
 * @property {string} name
 * some comment
 * @property {number} age
 */

/** @type {/*t1*/MyType1} */
const obj1 = { /*1n*/name: "", /*1a*/age: 10 };

/** @type {/*t2*/MyType2} */
const obj2 = { /*2n*/name: "", /*2a*/age: 10 };
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyQuickInfoAt(t, "t1", "type MyType1 = { name: string; age: number; }", "")
	f.VerifyQuickInfoAt(t, "t2", "type MyType2 = { name: string; age: number; }", "")

	f.VerifyQuickInfoAt(t, "1n", "(property) name: string", "@-rule\n")
	f.VerifyQuickInfoAt(t, "2n", "(property) name: string", "some comment\n")

	f.VerifyQuickInfoAt(t, "1a", "(property) age: number", "")
	f.VerifyQuickInfoAt(t, "2a", "(property) age: number", "")
}
