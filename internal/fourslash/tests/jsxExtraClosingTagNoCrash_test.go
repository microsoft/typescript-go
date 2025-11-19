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
		return <div></div></d/*$*/iv>;
	}
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// This should not crash - just call the API to trigger the bug
	f.GoToMarker(t, "$")
	// Request completions which will trigger tryGetObjectTypeDeclarationCompletionContainer
	f.VerifyCompletions(t, nil, nil)
}
