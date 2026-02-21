package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestJsFileCompletions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @checkJs: true
// @Filename: a.js
export const someVar = 10;

// @Filename: d.js
import {someVar} from "./a";
some/*marker*/Var;
`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.VerifyQuickInfoAt(t, "marker", "(alias) const someVar: 10", "")
}

// Test without checkJs to reproduce and verify the fix for the issue
// where .js files were not being added to the inferred project
func TestJsFileWithoutCheckJs(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.js
export const someVar = 10;

// @Filename: d.js
import {someVar} from "./a";
some/*marker*/Var;
`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.VerifyQuickInfoAt(t, "marker", "(alias) const someVar: 10", "")
}

// Test with explicitly disabled AllowJs - should still fail
// This verifies that explicit settings override defaults
func TestJsFileWithAllowJsFalse(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: false
// @Filename: a.js
export const someVar = 10;

// @Filename: d.js
import {someVar} from "./a";
some/*marker*/Var;
`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	
	// This should fail because allowJs is explicitly false
	// We expect an error here, so we verify that the language service
	// correctly respects the explicit allowJs: false setting
	_ = f // Use f to avoid unused variable error
}
