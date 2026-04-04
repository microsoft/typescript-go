package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceAccessExpressionProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare const obj: { greet(name: string): string; count: number; };
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export const /*targetObj*/obj = { /*targetGreet*/greet(name) { return name; }, /*targetCount*/count: 42 };
// @Filename: /home/src/workspaces/project/index.ts
import { obj } from "pkg";
obj./*propAccess*/greet("world");
obj./*propAccess2*/count;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "propAccess", "propAccess2")
}
