//// [tests/cases/compiler/jsdocTypesWithPrefixesAndUnionTypes1.ts] ////

//// [question.js]
/** @param {?} x */
export function f0(x) {}

/** @param {?never} x */
export function f1(x) {}

/** @param {never?} x */
export function f2(x) {}

/** @param {? | never} x */
export function f3(x) {}

/** @param {? | string} x */
export function f4(x) {}

/** @param {number | ? | string} x */
export function f5(x) {}

/** @param {number | string | ?} x */
export function f6(x) {}

/** @param {? number | string} x */
export function f7(x) {}

/** @param {number? | string} x */
export function f8(x) {}

/** @param {number | ? string} x */
export function f9(x) {}

/** @param {? { a: number } & { b: number }} x */
export function f10(x) {}

/** @param {{ a: number } & ? { b: number }} x */
export function f11(x) {}

/** @param {? { a: number } & { b: number } | string} x */
export function f12(x) {}

/** @param {{ a: number } & ? { b: number } | string} x */
export function f13(x) {}

//// [exclamation.js]
/** @param {!} x */
export function g0(x) {}

/** @param {!never} x */
export function g1(x) {}

/** @param {never!} x */
export function g2(x) {}

/** @param {! | never} x */
export function g3(x) {}

/** @param {! | string} x */
export function g4(x) {}

/** @param {number | ! | string} x */
export function g5(x) {}

/** @param {number | string | !} x */
export function g6(x) {}

/** @param {! number | string} x */
export function g7(x) {}

/** @param {number! | string} x */
export function g8(x) {}

/** @param {number | ! string} x */
export function g9(x) {}

/** @param {! { a: number } & { b: number }} x */
export function g10(x) {}

/** @param {{ a: number } & ! { b: number }} x */
export function g11(x) {}

/** @param {! { a: number } & { b: number } | string} x */
export function g12(x) {}

/** @param {{ a: number } & ! { b: number } | string} x */
export function g13(x) {}



//// [question.d.ts]
/** @param {?} x */
export declare function f0(x: any | null): void;
/** @param {?never} x */
export declare function f1(x: never | null): void;
/** @param {never?} x */
export declare function f2(x: never | null): void;
/** @param {? | never} x */
export declare function f3(x: (any | null) | never): void;
/** @param {? | string} x */
export declare function f4(x: (any | null) | string): void;
/** @param {number | ? | string} x */
export declare function f5(x: number | (any | null) | string): void;
/** @param {number | string | ?} x */
export declare function f6(x: number | string | (any | null)): void;
/** @param {? number | string} x */
export declare function f7(x: (number | null) | string): void;
/** @param {number? | string} x */
export declare function f8(x: number): void;
/** @param {number | ? string} x */
export declare function f9(x: number | (string | null)): void;
/** @param {? { a: number } & { b: number }} x */
export declare function f10(x: ({
    a: number;
} | null) & {
    b: number;
}): void;
/** @param {{ a: number } & ? { b: number }} x */
export declare function f11(x: {
    a: number;
} & ({
    b: number;
} | null)): void;
/** @param {? { a: number } & { b: number } | string} x */
export declare function f12(x: (({
    a: number;
} | null) & {
    b: number;
}) | string): void;
/** @param {{ a: number } & ? { b: number } | string} x */
export declare function f13(x: ({
    a: number;
} & ({
    b: number;
} | null)) | string): void;
//// [exclamation.d.ts]
/** @param {!} x */
export declare function g0(x: any): void;
/** @param {!never} x */
export declare function g1(x: never): void;
/** @param {never!} x */
export declare function g2(x: never): void;
/** @param {! | never} x */
export declare function g3(x: any | never): void;
/** @param {! | string} x */
export declare function g4(x: any | string): void;
/** @param {number | ! | string} x */
export declare function g5(x: number | any | string): void;
/** @param {number | string | !} x */
export declare function g6(x: number | string | any): void;
/** @param {! number | string} x */
export declare function g7(x: number | string): void;
/** @param {number! | string} x */
export declare function g8(x: number | string): void;
/** @param {number | ! string} x */
export declare function g9(x: number | string): void;
/** @param {! { a: number } & { b: number }} x */
export declare function g10(x: {
    a: number;
} & {
    b: number;
}): void;
/** @param {{ a: number } & ! { b: number }} x */
export declare function g11(x: {
    a: number;
} & {
    b: number;
}): void;
/** @param {! { a: number } & { b: number } | string} x */
export declare function g12(x: ({
    a: number;
} & {
    b: number;
}) | string): void;
/** @param {{ a: number } & ! { b: number } | string} x */
export declare function g13(x: ({
    a: number;
} & {
    b: number;
}) | string): void;
