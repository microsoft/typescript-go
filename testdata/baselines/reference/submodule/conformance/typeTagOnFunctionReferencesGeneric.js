//// [tests/cases/conformance/salsa/typeTagOnFunctionReferencesGeneric.ts] ////

//// [typeTagOnFunctionReferencesGeneric.js]
/**
 * @typedef {<T>(m : T) => T} IFn
 */

/**@type {IFn}*/
export function inJs(l) {
    return l;
}
inJs(1); // lints error. Why?

/**@type {IFn}*/
const inJsArrow = (j) => {
    return j;
}
inJsArrow(2); // no error gets linted as expected


//// [typeTagOnFunctionReferencesGeneric.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.inJs = inJs;
function inJs(l) {
    return l;
}
inJs(1);
const inJsArrow = (j) => {
    return j;
};
inJsArrow(2);
