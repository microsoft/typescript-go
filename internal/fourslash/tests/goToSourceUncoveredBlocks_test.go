package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceAtTypesPackage(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the GetResolvedModuleFromModuleSpecifier fallback path:
	// NoDts resolver can't resolve "foo" to any .js (only @types/foo has .d.ts),
	// so resolveImplementationFromModuleName returns "". Then the program's cached
	// resolution to @types/foo/index.d.ts is used, and findImplementationFileFromDtsFileName
	// maps @types/foo → foo and finds the .js via the real foo package.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/@types/foo/package.json
{ "name": "@types/foo", "version": "1.0.0" }
// @Filename: /home/src/workspaces/project/node_modules/@types/foo/index.d.ts
export declare function bar(): string;
// @Filename: /home/src/workspaces/project/node_modules/foo/package.json
{ "name": "foo", "version": "1.0.0", "main": "./index.js" }
// @Filename: /home/src/workspaces/project/node_modules/foo/index.js
export function /*target*/bar() { return "hello"; }
// @Filename: /home/src/workspaces/project/index.ts
import { bar } from "foo";
bar/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourceTripleSlashReference(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the getReferenceAtPosition fallback: cursor on a
	// /// <reference path="..."/> directive pointing to a .js file.
	const content = `// @allowJs: true
// @Filename: /home/src/workspaces/project/helper.js
/*target*/function helper() { return 1; }
// @Filename: /home/src/workspaces/project/index.ts
/// <reference path="./[|helper.js/*refPath*/|]" />
declare function helper(): number;
helper();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "refPath")
}

func TestGoToSourceDeclarationMapSourceMap(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the tryGetSourcePosition path: .d.ts has a sourcemap
	// pointing back to the original .ts source.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./dist/index.js", "types": "./dist/index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/src/index.ts
export function /*target*/greet() { return "hi"; }
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.d.ts
export declare function greet(): string;
//# sourceMappingURL=index.d.ts.map
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.d.ts.map
{"version":3,"file":"index.d.ts","sourceRoot":"","sources":["../src/index.ts"],"names":[],"mappings":"AAAA,wBAAgB,KAAK,WAAY"}
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.greet = greet;
function greet() { return "hi"; }
// @Filename: /home/src/workspaces/project/index.ts
import { greet } from "pkg";
greet/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}

func TestGoToSourcePackageRootThenSubpath(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers the tryPackageRootFirst fallback: the .d.ts is at index.d.ts
	// in the package root. Root package resolution ("pkg") fails because there's
	// no main entry, but subpath resolution ("pkg/index") succeeds.
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
