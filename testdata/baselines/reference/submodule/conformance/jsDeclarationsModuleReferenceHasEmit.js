//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsModuleReferenceHasEmit.ts] ////

//// [index.js]
/**
 * @module A
 */
class A {}


/**
 * Target element
 * @type {module:A}
 */
export let el = null;

export default A;

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.el = void 0;
class A {
}
exports.el = null;
exports.default = A;
