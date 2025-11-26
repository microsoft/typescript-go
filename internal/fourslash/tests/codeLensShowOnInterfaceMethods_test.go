package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func runCodeLensShowOnInterfaceMethods(t *testing.T, showOnInterfaceMethods bool) {
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
export interface I {
  methodA(): void;
}
export interface I {
  methodB(): void;
}

interface J extends I {
  methodB(): void;
  methodC(): void;
}

class C implements J {
  methodA(): void {}
  methodB(): void {}
  methodC(): void {}
}

class AbstractC implements J {
  abstract methodA(): void;
  methodB(): void {}
  abstract methodC(): void;
}
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
		ImplementationsCodeLensEnabled:                true,
		ImplementationsCodeLensShowOnInterfaceMethods: showOnInterfaceMethods,
	})
}

func TestCodeLensReferencesShowOnInterfaceMethodsTrue(t *testing.T) {
	t.Parallel()
	runCodeLensShowOnInterfaceMethods(t, true)
}

func TestCodeLensReferencesShowOnInterfaceMethodsFalse(t *testing.T) {
	t.Parallel()
	runCodeLensShowOnInterfaceMethods(t, false)
}
