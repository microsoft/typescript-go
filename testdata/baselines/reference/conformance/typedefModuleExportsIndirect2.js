//// [tests/cases/conformance/salsa/typedefModuleExportsIndirect2.ts] ////

//// [typedefModuleExportsIndirect2.js]
/** @typedef {{ a: 1, m: 1 }} C */
const f = function() {};
module.exports = f;
//// [use.js]
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c


//// [typedefModuleExportsIndirect2.js]
"use strict";
/** @typedef {{ a: 1, m: 1 }} C */
const f = function () { };
module.exports = f;
//// [use.js]
"use strict";
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c;


//// [typedefModuleExportsIndirect2.d.ts]
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
