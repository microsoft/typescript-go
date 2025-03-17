//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsReusesExistingNodesMappingJSDocTypes.ts] ////

//// [index.js]
/** @type {?} */
export const a = null;

/** @type {*} */
export const b = null;

/** @type {string?} */
export const c = null;

/** @type {string=} */
export const d = null;

/** @type {string!} */
export const e = null;

/** @type {function(string, number): object} */
export const f = null;

/** @type {function(new: object, string, number)} */
export const g = null;

/** @type {Object.<string, number>} */
export const h = null;


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.h = exports.g = exports.f = exports.e = exports.d = exports.c = exports.b = exports.a = void 0;
exports.a = null;
exports.b = null;
exports.c = null;
exports.d = null;
exports.e = null;
exports.f = null;
exports.g = null;
exports.h = null;
