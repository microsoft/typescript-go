package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAllowRenameOfImportPath2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /src/example.ts
import brushPackageJson from './visx-brush/[|package|].json';
// @Filename: /src/visx-brush/package.json
{ "name": "brush" }
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	prefs := &lsutil.UserPreferences{
		AllowRenameOfImportPath: core.TSTrue,
	}
	f.VerifyBaselineRename(t, prefs, f.Ranges()[0])
}
