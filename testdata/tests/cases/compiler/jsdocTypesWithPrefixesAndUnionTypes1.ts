// @checkJs: true
// @declaration: true
// @isolatedDeclarations: true

// @ts-check

// **Test motivation**
// Previously `?` was a valid "standalone" type that just meant `any`.
// This meant that it was possible to write types like `? | string` which would be equivalent to `any | string`.
// `?` is also allowed as a prefix operator though, so `?number` is equivalent to `number | null`.
//
// In the 7.0 port, `?` ceased to be a valid standalone type,
// but somehow this allowed constructs like `? | string` to be valid,
// which is not what we want. This test is meant to validate how `?` is
// handled when written in different contexts - as a prefix/postfix operator, as a standalone type, and in union types.

// @filename: question.js
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

// @filename: exclamation.js
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