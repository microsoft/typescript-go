//// [tests/cases/compiler/declarationEmitCastReusesTypeNode4.ts] ////

//// [input.js]
/**
 * @typedef {{ } & { name?: string }} P
 */

const something = /** @type {*} */(null);

export let vLet = /** @type {P} */(something);
export const vConst = /** @type {P} */(something);

export function fn(p = /** @type {P} */(something)) {}

/** @param {number} req */
export function fnWithRequiredDefaultParam(p = /** @type {P} */(something), req) {}

export class C {
    field = /** @type {P} */(something);
    /** @optional */ optField = /** @type {P} */(something); // not a thing
    /** @readonly */ roFiled = /** @type {P} */(something);
    method(p = /** @type {P} */(something)) {}
    /** @param {number} req */
    methodWithRequiredDefault(p = /** @type {P} */(something), req) {}

    constructor(ctorField = /** @type {P} */(something)) {}

    get x() { return /** @type {P} */(something) }
    set x(v) { }
}

export default /** @type {P} */(something);

// allows `undefined` on the input side, thanks to the initializer
/**
 * 
 * @param {P} x
 * @param {number} b
 */
export function fnWithPartialAnnotationOnDefaultparam(x = /** @type {P} */(something), b) {}

//// [input.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = exports.vConst = exports.vLet = void 0;
exports.fn = fn;
exports.fnWithRequiredDefaultParam = fnWithRequiredDefaultParam;
exports.fnWithPartialAnnotationOnDefaultparam = fnWithPartialAnnotationOnDefaultparam;
const something = (null);
exports.vLet = (something);
exports.vConst = (something);
function fn(p = (something)) { }
function fnWithRequiredDefaultParam(p = (something), req) { }
class C {
    field = (something);
    optField = (something);
    roFiled = (something);
    method(p = (something)) { }
    methodWithRequiredDefault(p = (something), req) { }
    constructor(ctorField = (something)) { }
    get x() { return (something); }
    set x(v) { }
}
exports.C = C;
exports.default = (something);
function fnWithPartialAnnotationOnDefaultparam(x = (something), b) { }
