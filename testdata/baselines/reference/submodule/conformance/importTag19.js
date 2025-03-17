//// [tests/cases/conformance/jsdoc/importTag19.ts] ////

//// [a.ts]
export interface Foo {}

//// [b.js]
/**
 * @import { Foo }
 * from "./a"
 */

/**
 * @param {Foo} a
 */
export function foo(a) {}


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo = foo;
function foo(a) { }
