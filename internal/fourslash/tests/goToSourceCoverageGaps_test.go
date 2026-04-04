package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceIndexSignatureProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// When accessing a property defined via index signature, getDeclarationsFromLocation
	// returns empty, so the GetPropertyOfType fallback is used.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare const config: { readonly [key: string]: string; name: string };
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export const config = { /*targetName*/name: "test" };
// @Filename: /home/src/workspaces/project/index.ts
import { config } from "pkg";
config./*propAccess*/name;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "propAccess")
}

func TestGoToSourceAliasedImportSpecifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the propertyName branch in getImportNamesForModuleSpecifier
	// when using `import { original as alias }`.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function original(): string;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export function /*target*/original() { return "ok"; }
// @Filename: /home/src/workspaces/project/index.ts
import { original as /*aliasedImport*/renamed } from "pkg";
renamed();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "aliasedImport")
}

func TestGoToSourceDtsReExport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// The .d.ts declaration itself re-exports from another module,
	// so findContainingModuleSpecifier(declaration) finds that specifier.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.d.ts
export declare function helper(): void;
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.js
export function /*target*/helper() {}
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export { helper } from "./impl";
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export { helper } from "./impl.js";
// @Filename: /home/src/workspaces/project/index.ts
import { helper } from "pkg";
helper/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourcePackageIndexDts(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// When the .d.ts is index.d.ts, tryPackageRootFirst is true,
	// so package root resolution is tried before subpath.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./lib/index.js", "types": "./lib/index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/lib/index.d.ts
export declare function greet(): string;
// @Filename: /home/src/workspaces/project/node_modules/pkg/lib/index.js
export function /*target*/greet() { return "hi"; }
// @Filename: /home/src/workspaces/project/index.ts
import { greet } from "pkg";
greet/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourceExportAssignmentDefault(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the ExportAssignment/default path in findDeclarationNodesByName
	// and getCandidateSourceDeclarationNames.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
declare const _default: { run(): void };
export default _default;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
/*target*/export default { run() {} };
// @Filename: /home/src/workspaces/project/index.ts
import pkg from "pkg";
pkg/*usage*/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourceBarrelReExportChain(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Tests forwarded implementation files: index.js re-exports from impl.js,
	// causing getForwardedImplementationFiles to follow the chain.
	// Also tests the non-concrete forwarded merge path.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.js
export function /*target*/doWork() { return 42; }
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.d.ts
export declare function doWork(): number;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export { doWork } from "./impl";
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export { doWork } from "./impl.js";
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/doWork } from "pkg";
doWork/*callSite*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName", "callSite")
}
