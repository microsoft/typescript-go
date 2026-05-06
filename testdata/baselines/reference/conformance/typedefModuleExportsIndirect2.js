//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect2.ts] ////

//// [typedefModuleExportsIndirect2.js]
/** @typedef {{ a: 1, m: 1 }} C */
const f = function() {};
module.exports = f;
//// [use.js]
/** @typedef {import('./typedefModuleExportsIndirect2').C} C */
/** @type {C} */
var c


//// [typedefModuleExportsIndirect2.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const f = function () { };
module.exports = f;
//// [use.js]
"use strict";
/** @typedef {import('./typedefModuleExportsIndirect2').C} C */
/** @type {C} */
var c;


//// [typedefModuleExportsIndirect2.d.ts]
export type C = {
    a: 1;
    m: 1;
};
/** @typedef {{ a: 1, m: 1 }} C */
const f: () => void;
export = f;
//// [use.d.ts]
type C = import('./typedefModuleExportsIndirect2').C;
/** @typedef {import('./typedefModuleExportsIndirect2').C} C */
/** @type {C} */
var c: C;


//// [DtsFileErrors]


dist/typedefModuleExportsIndirect2.d.ts(6,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
dist/use.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/typedefModuleExportsIndirect2.d.ts (1 errors) ====
    export type C = {
        a: 1;
        m: 1;
    };
    /** @typedef {{ a: 1, m: 1 }} C */
    const f: () => void;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export = f;
    
==== dist/use.d.ts (1 errors) ====
    type C = import('./typedefModuleExportsIndirect2').C;
    /** @typedef {import('./typedefModuleExportsIndirect2').C} C */
    /** @type {C} */
    var c: C;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    