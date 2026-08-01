// @allowJs: true
// @checkJs: true
// @noEmit: true
// @lib: es2022

// @filename: keys.js
const idsSymbol = /** @type {symbol} */ (Symbol("ids"));
module.exports.idsSymbol = idsSymbol;

// @filename: serializables.js
module.exports = {
    "a/B": 1,
    "c/D": 2,
};

// @filename: index.js
/** @typedef {import("./keys").idsSymbol} IDsSymbol */
/** @typedef {keyof import("./serializables")} SerializableKey */

/** @type {Record<IDsSymbol, string[]>} */
const values = {};
values[require("./keys").idsSymbol] = [];

const serializables = require("./serializables");
/** @param {SerializableKey} key */
const load = key => serializables[key];
load("a/B");
