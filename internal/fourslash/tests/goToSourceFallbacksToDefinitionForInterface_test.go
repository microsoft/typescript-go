package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceFallbacksToDefinitionForInterface(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export interface /*target*/Config {
    enabled: boolean;
}
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
exports.makeConfig = () => ({ enabled: true });
// @Filename: /home/src/workspaces/project/index.ts
import type { /*importName*/Config } from "pkg";
let value: /*typeRef*/Config;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName", "typeRef")
}
