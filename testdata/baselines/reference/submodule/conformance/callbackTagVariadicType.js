//// [tests/cases/conformance/jsdoc/callbackTagVariadicType.ts] ////

//// [callbackTagVariadicType.js]
/**
 * @callback Foo
 * @param {...string} args
 * @returns {number}
 */

/** @type {Foo} */
export const x = () => 1
var res = x('a', 'b')


//// [callbackTagVariadicType.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
const x = () => 1;
exports.x = x;
var res = (0, exports.x)('a', 'b');
