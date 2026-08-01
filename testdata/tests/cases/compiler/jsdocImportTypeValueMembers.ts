// @allowJs: true
// @checkJs: true
// @noEmit: true
// @lib: es2022

// @filename: values.js
const VALUE = "value";
module.exports.VALUE = VALUE;

// @filename: mixed.js
class Mixed {}
const idsSymbol = Symbol("ids");
module.exports = Mixed;
module.exports.idsSymbol = idsSymbol;

// @filename: index.js
/** @typedef {import("./values").VALUE} VALUE */
/** @typedef {import("./mixed").idsSymbol} IDsSymbol */

/** @type {VALUE} */
const value = "value";
/** @type {IDsSymbol} */
const id = require("./mixed").idsSymbol;
/** @type {typeof import("./values").VALUE} */
const valueType = "value";
