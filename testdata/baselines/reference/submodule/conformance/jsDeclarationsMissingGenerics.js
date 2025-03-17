//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsMissingGenerics.ts] ////

//// [file.js]
/**
 * @param {Array} x
 */
function x(x) {}
/**
 * @param {Promise} x
 */
function y(x) {}

//// [file.js]
function x(x) { }
function y(x) { }
