package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceFallbackToModuleSpecifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// When the specific name can't be found in the .js implementation file
	// (because the JS uses a different export pattern), shouldFallbackToModuleSpecifier
	// triggers and re-resolves via the module specifier with nil names, returning
	// the entry declaration of the .js file.
	const content = `// @moduleResolution: bundler
// @allowJs: true
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function internalHelper(): void;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
/*entryPoint*/Object.defineProperty(exports, "internalHelper", { value: function() {} });
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/internalHelper } from "pkg";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName")
}
