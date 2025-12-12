package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatDocumentWithJSDoc(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**
 * JSDoc for things
 */
function f() {
    /** more
        jsdoc */
    var t;
    /**
     * multiline
     */
    var multiline;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `/**
 * JSDoc for things
 */
function f() {
    /** more
        jsdoc */
    var t;
    /**
     * multiline
     */
    var multiline;
}`)
}
