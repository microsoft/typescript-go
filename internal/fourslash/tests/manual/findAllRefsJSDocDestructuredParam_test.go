package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test for panic: "Unhandled case in Node.Text: *ast.BindingPattern"
// This occurred when requesting findAllReferences with a JSDoc parameter tag
// that referenced a destructured parameter.
func TestFindAllRefsJSDocDestructuredParam(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
/**
 * @param /*1*/obj - object destructured parameter
 */
function f({ x, y }) {}
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "1")
}
