//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect1.ts] ////

//// [typedefModuleExportsIndirect1.js]
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
module.exports = dummy;
//// [use.js]
/** @typedef {import('./typedefModuleExportsIndirect1').C} C */
/** @type {C} */
var c


//// [typedefModuleExportsIndirect1.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
module.exports = dummy;
//// [use.js]
"use strict";
/** @typedef {import('./typedefModuleExportsIndirect1').C} C */
/** @type {C} */
var c;


//// [typedefModuleExportsIndirect1.d.ts]
export type C = {
    a: 1;
    m: 1;
};
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
export = dummy;
//// [use.d.ts]
type C = import('./typedefModuleExportsIndirect1').C;
/** @typedef {import('./typedefModuleExportsIndirect1').C} C */
/** @type {C} */
var c: C;


//// [DtsFileErrors]


dist/typedefModuleExportsIndirect1.d.ts(6,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
dist/use.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/typedefModuleExportsIndirect1.d.ts (1 errors) ====
    export type C = {
        a: 1;
        m: 1;
    };
    /** @typedef {{ a: 1, m: 1 }} C */
    const dummy = 0;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export = dummy;
    
==== dist/use.d.ts (1 errors) ====
    type C = import('./typedefModuleExportsIndirect1').C;
    /** @typedef {import('./typedefModuleExportsIndirect1').C} C */
    /** @type {C} */
    var c: C;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    