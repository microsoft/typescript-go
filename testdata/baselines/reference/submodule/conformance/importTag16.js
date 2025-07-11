//// [tests/cases/conformance/jsdoc/importTag16.ts] ////

//// [a.ts]
export default interface Foo {}
export interface I {}

//// [b.js]
/** @import Foo, { I } from "./a" */

/**
 * @param {Foo} a
 * @param {I} b
 */
export function foo(a, b) {}




//// [a.d.ts]
export default interface Foo {
}
export interface I {
}
//// [b.d.ts]
import type Foo, { I } from "./a";
/** @import Foo, { I } from "./a" */
/**
 * @param {Foo} a
 * @param {I} b
 */
export declare function foo(a: Foo, b: I): void;
