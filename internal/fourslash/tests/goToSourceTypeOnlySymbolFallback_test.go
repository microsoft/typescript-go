package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceTypeOnlySymbolFallback(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// When a type-only symbol (type alias) is imported with a regular import and used
	// in a value position, source definition should fall back to regular definition
	// (the .d.ts declaration) since there's no concrete JS implementation.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/types.d.ts
export interface Config { enabled: boolean; }
// @Filename: /home/src/workspaces/project/node_modules/pkg/types.js
// no runtime content for Config interface
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export { Config } from "./types";
export declare function makeConfig(): Config;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export { Config } from "./types.js";
export function makeConfig() { return { enabled: true }; }
// @Filename: /home/src/workspaces/project/index.ts
import { Config, makeConfig } from "pkg";
let c: /*typeRef*/Config;
makeConfig/*callRef*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "typeRef", "callRef")
}
