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


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo = foo;
function foo(a, b) { }
