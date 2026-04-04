package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceMappedTypeProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers lines 72-74: getDeclarationsFromLocation returns empty for a property
	// that exists only via a mapped type (no explicit declaration), so
	// GetPropertyOfType fallback is used.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
type Keys = "a" | "b";
export declare const obj: { [K in Keys]: number };
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export const obj = { a: 1, /*target*/b: 2 };
// @Filename: /home/src/workspaces/project/index.ts
import { obj } from "pkg";
obj./*propAccess*/b;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "propAccess")
}

func TestGoToSourceForwardedNonConcreteMerge(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 415: forwarded declarations are non-concrete,
	// so they merge with the initial non-concrete declarations.
	// The barrel index.js re-exports from types.js which only has type re-exports.
	const content = `// @moduleResolution: bundler
// @allowJs: true
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export { Config } from "./types";
// @Filename: /home/src/workspaces/project/node_modules/pkg/types.d.ts
export interface Config { enabled: boolean; }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
export { Config } from "./types.js";
// @Filename: /home/src/workspaces/project/node_modules/pkg/types.js
// Config is a type, no runtime value
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/Config } from "pkg";
let c: Config;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName")
}

func TestGoToSourceFilterPreferredFallbackAll(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 527: filterPreferredSourceDeclarations returns all declarations
	// when none are property-like and none are concrete. This happens with
	// re-export specifiers matching the name in the .js file.
	const content = `// @moduleResolution: bundler
// @allowJs: true
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./barrel.js", "types": "./barrel.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/barrel.d.ts
export { value } from "./impl";
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.d.ts
export declare const value: number;
// @Filename: /home/src/workspaces/project/node_modules/pkg/barrel.js
export { value } from "./impl.js";
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.js
export const /*target*/value = 42;
// @Filename: /home/src/workspaces/project/index.ts
import { /*importName*/value } from "pkg";
console.log(value);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "importName")
}

func TestGoToSourceDeclarationMapFallback(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Covers line 608: findClosestDeclarationNode walks up parents and finds
	// no declaration, returns entry node. This happens when source map
	// points to a position that's not inside any declaration.
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./dist/index.js", "types": "./dist/index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/src/index.ts
/*target*/console.log("side effect");
export function greet() { return "hi"; }
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.d.ts
export declare function greet(): string;
//# sourceMappingURL=index.d.ts.map
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.d.ts.map
{"version":3,"file":"index.d.ts","sourceRoot":"","sources":["../src/index.ts"],"names":[],"mappings":"AAC6B"}
// @Filename: /home/src/workspaces/project/node_modules/pkg/dist/index.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.greet = greet;
console.log("side effect");
function greet() { return "hi"; }
// @Filename: /home/src/workspaces/project/index.ts
import { greet } from "pkg";
greet/*usage*/();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "usage")
}
