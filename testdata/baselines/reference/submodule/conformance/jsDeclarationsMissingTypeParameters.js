//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsMissingTypeParameters.ts] ////

//// [file.js]
/**
  * @param {Array=} y desc
  */
function x(y) { }

// @ts-ignore
/** @param {function (Array)} func Invoked
 */
function y(func) { return; }

/**
 * @return {(Array.<> | null)} list of devices
 */
function z() { return null ;}

/**
 * 
 * @return {?Promise} A promise
 */
function w() { return null; }

//// [file.js]
function x(y) { }
function y(func) { return; }
function z() { return null; }
function w() { return null; }
