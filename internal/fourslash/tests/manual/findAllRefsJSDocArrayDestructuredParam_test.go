package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test for panic: "Unhandled case in Node.Text: *ast.BindingPattern"
// This tests array binding patterns specifically.
func TestFindAllRefsJSDocArrayDestructuredParam(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
/**
 * @param /*1*/arr - array destructured parameter
 */
function f([x, y]) {}
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "1")
}
