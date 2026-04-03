package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToSourceCJSReExportViaDefineProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @moduleResolution: bundler
// @Filename: /home/src/workspaces/project/node_modules/pkg/package.json
{ "name": "pkg", "main": "./index.js", "types": "./index.d.ts" }
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.d.ts
export declare function greet(name: string): string;
export declare enum TargetPopulation {
    Team = "team",
    Public = "public",
}
// @Filename: /home/src/workspaces/project/node_modules/pkg/index.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TargetPopulation = exports.greet = void 0;
var impl_1 = require("./impl");
Object.defineProperty(exports, "greet", { enumerable: true, get: function () { return impl_1.greet; } });
var types_1 = require("./types");
Object.defineProperty(exports, "TargetPopulation", { enumerable: true, get: function () { return types_1.TargetPopulation; } });
// @Filename: /home/src/workspaces/project/node_modules/pkg/impl.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.greet = void 0;
function /*greetImpl*/greet(name) { return "Hello, " + name; }
exports.greet = greet;
// @Filename: /home/src/workspaces/project/node_modules/pkg/types.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TargetPopulation = void 0;
var /*targetPopulationImpl*/TargetPopulation;
(function (TargetPopulation) {
    TargetPopulation["Team"] = "team";
    TargetPopulation["Public"] = "public";
})(TargetPopulation || (exports.TargetPopulation = TargetPopulation = {}));
// @Filename: /home/src/workspaces/project/index.ts
import { /*namedImport*/greet, /*enumImport*/TargetPopulation } from "pkg";
greet/*call*/("world");
TargetPopulation/*enumAccess*/.Team;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToSourceDefinition(t, "namedImport", "enumImport", "call", "enumAccess")
}
