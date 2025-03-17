//// [tests/cases/compiler/jsFileMethodOverloads5.ts] ////

//// [a.js]
/**
 * @overload
 * @param {string} a
 * @return {void}
 */

/**
 * @overload
 * @param {number} a
 * @param {number} [b]
 * @return {void}
 */

/**
 * @param {string | number} a
 * @param {number} [b]
 */
export const foo = function (a, b) { }


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo = void 0;
const foo = function (a, b) { };
exports.foo = foo;
