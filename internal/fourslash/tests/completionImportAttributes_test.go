package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionImportAttributes(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @target: esnext
// @module: esnext
// @filename: main.ts
import yadda1 from "yadda" with {/*attr*/}
import yadda2 from "yadda" with {attr/*attrEnd1*/: true}
import yadda3 from "yadda" with {attr: /*attrValue*/}

// @filename: yadda
export default {};
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Test completion at empty attributes
	// This should not panic
	f.VerifyCompletions(t, "attr", nil)

	// Test completion after attribute name
	// This should not panic
	f.VerifyCompletions(t, "attrEnd1", nil)

	// Test completion at attribute value position
	// This should not panic
	f.VerifyCompletions(t, "attrValue", nil)
}
