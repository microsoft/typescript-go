package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGetEditsForFileRename_cssImport1(t *testing.T) {
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
import styles from "./app.css";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	// We cannot request rename of the .d.css file because that would lead to a circularity of `willRenameFiles`.
	// So this case does not fully work.
	f.VerifyWillRenameFilesEdits(t, "/app.css", "/app2.css", map[string]string{
		"/a.ts": `import styles from "./app.css";`,
		"/app2.css": `.cookie-banner {
  display: none;
}`,
		"/app.d.css.ts": `declare const css: {
  cookieBanner: string;
};
export default css;`,
	}, nil /*preferences*/)
}
