package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceReExportedImplementation(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts", "type": "module" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export { foo } from "./foo";
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export { foo } from "./foo.js";
// @Filename: /home/src/workspaces/project/node_modules/pkg/foo.d.ts
export declare function foo(): string;
// @Filename: /home/src/workspaces/project/node_modules/pkg/foo.js
export function /*target*/foo() { return "ok"; }
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/foo } from "pkg";
foo/*start*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName", "start")
}
