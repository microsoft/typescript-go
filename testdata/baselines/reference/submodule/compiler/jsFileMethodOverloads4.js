//// [tests/cases/compiler/jsFileMethodOverloads4.ts] ////

//// [a.js]
export function Foo() { }

/**
 * @overload
 * @param {string} a
 * @return {void}
 */

/**
 * @overload
 * @param {number} a
 * @param {string} b
 * @return {void}
 */

/**
 * @param {string | number} a
 * @param {string} [b]
 * @return {void}
 */
Foo.prototype.bar = function (a, b) { }


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = Foo;
function Foo() { }
Foo.prototype.bar = function (a, b) { };
