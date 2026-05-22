package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindReferencesJSXTagName3VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @jsx: preserve
// @Filename: /a.tsx
namespace JSX {
    export interface Element { }
    export interface IntrinsicElements {
        [|[|/*1*/div|]: any;|]
    }
}

[|const [|/*6*/Comp|] = () =>
    [|<[|/*2*/div|]>
        Some content
        [|<[|/*3*/div|]>More content</[|/*4*/div|]>|]
    </[|/*5*/div|]>|];|]

const x = [|<[|/*7*/Comp|]>
    Content
</[|/*8*/Comp|]>|];`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8")
}
