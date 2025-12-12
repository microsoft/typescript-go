package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatTsxClosingAfterJsxText(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: foo.tsx

const a = (
    <div>
        text
               </div>
)
const b = (
    <div>
        text
      twice
               </div>
)
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
const a = (
    <div>
        text
    </div>
)
const b = (
    <div>
        text
        twice
    </div>
)
`)
}
