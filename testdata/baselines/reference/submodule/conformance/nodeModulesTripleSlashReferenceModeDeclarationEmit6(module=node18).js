//// [tests/cases/conformance/node/nodeModulesTripleSlashReferenceModeDeclarationEmit6.ts] ////

//// [package.json]
{
    "name": "pkg",
    "version": "0.0.1",
    "exports": {
        "import": "./import.js",
        "require": "./require.js"
    }
}
//// [import.d.ts]
export {};
declare global {
    interface ImportInterface {}
    function getInterI(): ImportInterface;
}
//// [require.d.ts]
export {};
declare global {
    interface RequireInterface {}
    function getInterR(): RequireInterface;
}
//// [uses.ts]
/// <reference types="pkg" preserve="true" />
export default getInterR();
//// [index.ts]
import obj from "./uses.js"
export default (obj as typeof obj);

//// [uses.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/// <reference types="pkg" preserve="true" />
exports.default = getInterR();
//// [index.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const uses_js_1 = __importDefault(require("./uses.js"));
exports.default = uses_js_1.default;


//// [uses.d.ts]
/// <reference types="pkg" preserve="true" />
const _default: RequireInterface;
export default _default;
//// [index.d.ts]
import obj from "./uses.js";
const _default: typeof obj;
export default _default;


//// [DtsFileErrors]


out/index.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/uses.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    import obj from "./uses.js";
    const _default: typeof obj;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== /node_modules/pkg/package.json (0 errors) ====
    {
        "name": "pkg",
        "version": "0.0.1",
        "exports": {
            "import": "./import.js",
            "require": "./require.js"
        }
    }
==== /node_modules/pkg/import.d.ts (0 errors) ====
    export {};
    declare global {
        interface ImportInterface {}
        function getInterI(): ImportInterface;
    }
==== /node_modules/pkg/require.d.ts (0 errors) ====
    export {};
    declare global {
        interface RequireInterface {}
        function getInterR(): RequireInterface;
    }
==== out/uses.d.ts (1 errors) ====
    /// <reference types="pkg" preserve="true" />
    const _default: RequireInterface;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    