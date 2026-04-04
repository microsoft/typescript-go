package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceReExportModuleSpecifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function alpha(): string;
export declare function beta(): number;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export function /*targetAlpha*/alpha() { return "a"; }
export function /*targetBeta*/beta() { return 2; }
// @Filename: /home/src/workspaces/project/reexport.ts
export { alpha, beta } from [|"pkg"/*reExportSpecifier*/|];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "reExportSpecifier")
}
