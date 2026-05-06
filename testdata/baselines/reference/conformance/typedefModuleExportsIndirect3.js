//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect3.ts] ////

//// [typedefModuleExportsIndirect3.js]
/** @typedef {{ a: 1, m: 1 }} C */
const o = {};
module.exports = o;
//// [use.js]
/** @typedef {import('./typedefModuleExportsIndirect3').C} C */
/** @type {C} */
var c


//// [typedefModuleExportsIndirect3.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const o = {};
module.exports = o;
//// [use.js]
"use strict";
/** @typedef {import('./typedefModuleExportsIndirect3').C} C */
/** @type {C} */
var c;


//// [typedefModuleExportsIndirect3.d.ts]
export type C = {
    a: 1;
    m: 1;
};
/** @typedef {{ a: 1, m: 1 }} C */
const o: {};
export = o;
//// [use.d.ts]
type C = import('./typedefModuleExportsIndirect3').C;
/** @typedef {import('./typedefModuleExportsIndirect3').C} C */
/** @type {C} */
var c: C;


//// [DtsFileErrors]


dist/typedefModuleExportsIndirect3.d.ts(6,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
dist/use.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/typedefModuleExportsIndirect3.d.ts (1 errors) ====
    export type C = {
        a: 1;
        m: 1;
    };
    /** @typedef {{ a: 1, m: 1 }} C */
    const o: {};
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export = o;
    
==== dist/use.d.ts (1 errors) ====
    type C = import('./typedefModuleExportsIndirect3').C;
    /** @typedef {import('./typedefModuleExportsIndirect3').C} C */
    /** @type {C} */
    var c: C;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    