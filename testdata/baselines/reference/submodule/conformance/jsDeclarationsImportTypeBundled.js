//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsImportTypeBundled.ts] ////

//// [mod1.js]
/**
 * @typedef {{x: number}} Item
 */
/**
 * @type {Item};
 */
const x = {x: 12};
module.exports = x;
//// [index.js]
/** @type {(typeof import("./folder/mod1"))[]} */
const items = [{x: 12}];
module.exports = items;

//// [index.js]
const items = [{ x: 12 }];
module.exports = items;
//// [mod1.js]
const x = { x: 12 };
module.exports = x;
