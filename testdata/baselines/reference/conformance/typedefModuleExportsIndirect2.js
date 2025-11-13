//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect2.ts] ////

//// [controlFlowJSClassProperty.js]
/** @typedef {{ a: 1, m: 1 }} C */
const f = function() {};
module.exports = f;
//// [use.js]
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c


//// [controlFlowJSClassProperty.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const f = function () { };
module.exports = f;
//// [use.js]
"use strict";
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c;


//// [controlFlowJSClassProperty.d.ts]
export type C = {
    a: 1;
    m: 1;
};
export = f;
//// [use.d.ts]
type C = import('./controlFlowJSClassProperty').C;
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
declare var c: C;


//// [DtsFileErrors]


dist/controlFlowJSClassProperty.d.ts(5,1): error TS2309: An export assignment cannot be used in a module with other exported elements.
dist/controlFlowJSClassProperty.d.ts(5,10): error TS2304: Cannot find name 'f'.
dist/use.d.ts(1,49): error TS2694: Namespace 'unknown' has no exported member 'C'.


==== dist/controlFlowJSClassProperty.d.ts (2 errors) ====
    export type C = {
        a: 1;
        m: 1;
    };
    export = f;
    ~~~~~~~~~~~
!!! error TS2309: An export assignment cannot be used in a module with other exported elements.
             ~
!!! error TS2304: Cannot find name 'f'.
    
==== dist/use.d.ts (1 errors) ====
    type C = import('./controlFlowJSClassProperty').C;
                                                    ~
!!! error TS2694: Namespace 'unknown' has no exported member 'C'.
    /** @typedef {import('./controlFlowJSClassProperty').C} C */
    /** @type {C} */
    declare var c: C;
    