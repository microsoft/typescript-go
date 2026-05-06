//// [tests/cases/compiler/nodeNextCjsNamespaceImportDefault2.ts] ////

//// [a.cts]
export const a: number = 1;
export default 'string';
//// [foo.mts]
import d, {a} from './a.cjs';
import * as ns from './a.cjs';
export {d, a, ns};

d.a;
ns.default.a;

//// [a.cjs]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.a = void 0;
exports.a = 1;
exports.default = 'string';
//// [foo.mjs]
import d, { a } from './a.cjs';
import * as ns from './a.cjs';
export { d, a, ns };
d.a;
ns.default.a;


//// [a.d.cts]
export const a: number;
const _default = "string";
export default _default;
//// [foo.d.mts]
import d, { a } from './a.cjs';
import * as ns from './a.cjs';
export { d, a, ns };


//// [DtsFileErrors]


out/a.d.cts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/a.d.cts (1 errors) ====
    export const a: number;
    const _default = "string";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== out/foo.d.mts (0 errors) ====
    import d, { a } from './a.cjs';
    import * as ns from './a.cjs';
    export { d, a, ns };
    