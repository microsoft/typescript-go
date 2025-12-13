package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCodeLensReferencesIncludeImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
// @Filename: /abc.ts
export function abc() { }

// @Filename: /other.ts
import { abc } from "./abc";
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
		CodeLens: lsutil.CodeLensUserPreferences{
			ReferencesCodeLensEnabled:            true,
			ReferencesCodeLensShowOnAllFunctions: true,
		},
	})
}
