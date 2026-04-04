package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceCallThroughImport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// When calling an imported function, the checker returns both the import specifier
	// (in the current file) and the call signature target (from .d.ts → mapped to .js).
	// filterImportLikeDeclarations should remove the import specifier, keeping the .js target.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare class Widget {
    constructor(name: string);
    render(): void;
}
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export class /*targetWidget*/Widget {
    constructor(name) { this.name = name; }
    /*targetRender*/render() {}
}
// @Filename: /home/src/workspaces/project/index.ts
import { Widget } from "pkg";
const w = new /*constructorCall*/Widget("test");
w./*methodCall*/render();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "constructorCall", "methodCall")
}
