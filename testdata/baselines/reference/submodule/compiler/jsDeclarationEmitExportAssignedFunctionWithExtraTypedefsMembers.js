//// [tests/cases/compiler/jsDeclarationEmitExportAssignedFunctionWithExtraTypedefsMembers.ts] ////

//// [index.js]
/**
 * @typedef Options
 * @property {string} opt
 */

/**
 * @param {Options} options
 */
module.exports = function loader(options) {}


//// [index.js]
module.exports = function loader(options) { };
