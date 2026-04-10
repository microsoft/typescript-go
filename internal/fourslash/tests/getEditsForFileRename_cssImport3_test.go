package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGetEditsForFileRename_cssImport3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @Filename: /tsconfig.json
{ "compilerOptions": { "allowArbitraryExtensions": true } }
// @Filename: /app.css
.cookie-banner {
  display: none;
}
// @Filename: /app.d.css.ts
declare const css: {
  cookieBanner: string;
};
export default css;
// @Filename: /a.ts
import styles from ".//*rename*/app.css";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "rename")
	f.RenameAtCaret(t, "/app2.css")
	f.GoToFile(t, "/a.ts")
	f.VerifyCurrentFileContent(t, `import styles from "./app2.css";`)
	f.GoToFile(t, "/app2.d.css.ts")
	f.VerifyCurrentFileContent(t, `declare const css: {
  cookieBanner: string;
};
export default css;`)
	f.GoToFile(t, "/app2.css")
	f.VerifyCurrentFileContent(t, `.cookie-banner {
  display: none;
}`)
}
