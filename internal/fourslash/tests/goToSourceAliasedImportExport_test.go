package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceAliasedImportExport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare const foo: number;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
exports./*target*/foo = 1;
// @Filename: /home/src/workspaces/project/index.ts
import { foo as /*importAlias*/bar } from "pkg";
bar;
// @Filename: /home/src/workspaces/project/reexport.ts
export { foo as /*reExportAlias*/bar } from "pkg";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importAlias", "reExportAlias")
}
