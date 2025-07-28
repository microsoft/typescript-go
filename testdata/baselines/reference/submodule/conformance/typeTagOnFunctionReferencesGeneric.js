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
/**
 * @typedef {<T>(m : T) => T} IFn
 */
/**@type {IFn}*/
function inJs(l) {
    return l;
}
inJs(1); // lints error. Why?
/**@type {IFn}*/
const inJsArrow = (j) => {
    return j;
};
inJsArrow(2); // no error gets linted as expected


//// [typeTagOnFunctionReferencesGeneric.d.ts]
export type IFn = <T>(m: T) => T;
/**
 * @typedef {<T>(m : T) => T} IFn
 */
/**@type {IFn}*/
export declare function inJs(l: T): T;


//// [DtsFileErrors]


out/typeTagOnFunctionReferencesGeneric.d.ts(6,33): error TS2304: Cannot find name 'T'.
out/typeTagOnFunctionReferencesGeneric.d.ts(6,37): error TS2304: Cannot find name 'T'.


==== out/typeTagOnFunctionReferencesGeneric.d.ts (2 errors) ====
    export type IFn = <T>(m: T) => T;
    /**
     * @typedef {<T>(m : T) => T} IFn
     */
    /**@type {IFn}*/
    export declare function inJs(l: T): T;
                                    ~
!!! error TS2304: Cannot find name 'T'.
                                        ~
!!! error TS2304: Cannot find name 'T'.
    