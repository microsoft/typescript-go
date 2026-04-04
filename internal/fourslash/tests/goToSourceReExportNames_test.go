package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceReExportNames(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function foo(): string;
export declare function bar(): number;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export function /*targetFoo*/foo() { return "ok"; }
export function /*targetBar*/bar() { return 42; }
// @Filename: /home/src/workspaces/project/reexport.ts
export { /*reExportFoo*/foo, /*reExportBar*/bar } from "pkg";
// @Filename: /home/src/workspaces/project/index.ts
import { foo, bar } from [|"pkg"/*moduleSpecifier*/|];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "reExportFoo", "reExportBar", "moduleSpecifier")
}
