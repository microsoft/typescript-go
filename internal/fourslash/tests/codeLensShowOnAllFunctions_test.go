package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func runCodeLensShowOnAllFunctions(t *testing.T, showOnAllFunctions bool) {
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
export function f1(): void {}

function f2(): void {}

export const f3 = () => {};

const f4 = () => {};

const f5 = function() {};
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
		ReferencesCodeLensEnabled:            true,
		ReferencesCodeLensShowOnAllFunctions: showOnAllFunctions,
	})
}

func TestCodeLensReferencesShowOnAllFunctionsTrue(t *testing.T) {
	t.Parallel()
	runCodeLensShowOnAllFunctions(t, true)
}

func TestCodeLensReferencesShowOnAllFunctionsFalse(t *testing.T) {
	t.Parallel()
	runCodeLensShowOnAllFunctions(t, false)
}
