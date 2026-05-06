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



//// [input.d.ts]
/**
 * @typedef {{ } & { name?: string }} P
 */
export type P = {} & {
    name?: string;
};
export let vLet: P;
export const vConst: P;
export function fn(p?: P): void;
/** @param {number} req */
export function fnWithRequiredDefaultParam(p: P, req: number): void;
export class C {
    field: P;
    /** @optional */ optField: P;
    /** @readonly */ readonly roFiled: P;
    method(p?: P): void;
    /** @param {number} req */
    methodWithRequiredDefault(p: P, req: number): void;
    constructor(ctorField?: P);
    get x(): P;
    set x(v: P);
}
const _default: P;
export default _default;
/**
 *
 * @param {P} x
 * @param {number} b
 */
export function fnWithPartialAnnotationOnDefaultparam(x: P, b: number): void;


//// [DtsFileErrors]


input.d.ts(23,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== input.d.ts (1 errors) ====
    /**
     * @typedef {{ } & { name?: string }} P
     */
    export type P = {} & {
        name?: string;
    };
    export let vLet: P;
    export const vConst: P;
    export function fn(p?: P): void;
    /** @param {number} req */
    export function fnWithRequiredDefaultParam(p: P, req: number): void;
    export class C {
        field: P;
        /** @optional */ optField: P;
        /** @readonly */ readonly roFiled: P;
        method(p?: P): void;
        /** @param {number} req */
        methodWithRequiredDefault(p: P, req: number): void;
        constructor(ctorField?: P);
        get x(): P;
        set x(v: P);
    }
    const _default: P;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    /**
     *
     * @param {P} x
     * @param {number} b
     */
    export function fnWithPartialAnnotationOnDefaultparam(x: P, b: number): void;
    