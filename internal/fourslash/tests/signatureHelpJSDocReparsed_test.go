package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSignatureHelpJSDocReparsed(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @checkJs: true
// @Filename: test.js
/**
 * @param {string} x
 */
function foo(x) {
    foo(/*1*/);
}
/**
 * @this {string}
 */
function bar() {
    bar(/*2*/);
}
/**
 * @type {function(string): void}
 */
function qux(x) {
    qux(/*3*/);
}
/**
 * @template T
 * @param {T} x
 * @returns {T}
 */
function identity(x) {
    identity(/*4*/);
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineSignatureHelp(t)
}
