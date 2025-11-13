//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect1.ts] ////

//// [typedefModuleExportsIndirect1.js]
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
module.exports = dummy;
//// [use.js]
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c


//// [typedefModuleExportsIndirect1.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
module.exports = dummy;
//// [use.js]
"use strict";
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c;


//// [typedefModuleExportsIndirect1.d.ts]
export type C = {
    a: 1;
    m: 1;
};
export = dummy;
//// [use.d.ts]
type C = import('./controlFlowJSClassProperty').C;
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
declare var c: C;
