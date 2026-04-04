package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceNonDeclarationFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 191: declaration is in a .ts file (not .d.ts),
	// so mapDeclarationToSourceDefinitions returns it as-is.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/utils.ts
export function /*target*/helper() { return 1; }
// @Filename: /home/src/workspaces/project/index.ts
import { helper } from "./utils";
helper/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourceNoImplementationFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers lines 103, 196: no implementation file can be resolved
	// (types-only package with no .js).
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function typesOnly(): void;
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/typesOnly } from "pkg";
typesOnly/*callSite*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName", "callSite")
}

func TestGoToSourceEmptyNamesEntryFallback(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 212: getCandidateSourceDeclarationNames returns empty names,
	// so mapDeclarationToSourceDefinitions falls through to entry declarations.
	// This happens when the declaration has no name and originalNode is not
	// an identifier (e.g., a namespace import on the module specifier side).
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
declare const _default: { run(): void };
export default _default;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export default { run() {} };
// @Filename: /home/src/workspaces/project/index.ts
import /*defaultImport*/pkg from "pkg";
pkg.run();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "defaultImport")
}

func TestGoToSourceSubpathNotIndex(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 292: subpath resolution succeeds when the d.ts is NOT index.d.ts
	// (non-tryPackageRootFirst path).
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "types": "./lib/utils.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/lib/utils.d.ts
export declare function util(): void;
// @Filename: /home/src/workspaces/project/node_modules/pkg/lib/utils.js
export function /*target*/util() {}
// @Filename: /home/src/workspaces/project/index.ts
import { util } from "pkg";
util/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourceNamedExportsSpecifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers lines 134-135 (KindNamedExports in isImportOrExportName),
	// line 119 (names found, not import/export name), and
	// line 156 (shouldPreferModuleSpecifierResult — not import/export name).
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function foo(): string;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export function /*target*/foo() { return "ok"; }
// @Filename: /home/src/workspaces/project/index.ts
import { foo } from "pkg";
const result = foo/*valueUsage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "valueUsage")
}

func TestGoToSourceExportAssignmentExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 452 (ExportAssignment match in findDeclarationNodesByName)
	// and line 506 (default function/class match).
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export default function createThing(): { value: number };
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export default function createThing() { return { value: 42 }; }
// @Filename: /home/src/workspaces/project/index.ts
import /*defaultName*/createThing from "pkg";
createThing/*callDefault*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "defaultName", "callDefault")
}

func TestGoToSourcePackageRootFallsBackToSubpath(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers lines 298-301: tryPackageRootFirst is true (index.d.ts),
	// root resolution fails, falls back to subpath.
	// Package has no "main" and package name resolution fails,
	// but "pkg/index" subpath can find the .js.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function work(): void;
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export function /*target*/work() {}
// @Filename: /home/src/workspaces/project/index.ts
import { work } from "pkg";
work/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}
