// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: box.js
/** @template T */
class Box {
    /** @param {T} value */
    constructor(value) {
        this.value = value;
    }
}
module.exports = Box;
/** @typedef {string} Extra */

// @filename: index.js
const Box = require("./box");
const value = new Box("value");
/** @type {import("./box")<string>} */
const typed = value;
