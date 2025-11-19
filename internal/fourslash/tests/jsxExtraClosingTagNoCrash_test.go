package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestJsxExtraClosingTagNoCrash(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: file.tsx
class MyComponent {
	render() {
		return <div></div></div>;
	}
}
let x/*$*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// This should not crash - verify quickinfo on a safe location
	f.VerifyQuickInfoAt(t, "$", "let x: any", "")
}
