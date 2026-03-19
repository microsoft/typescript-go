// @module: nodenext
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
function something(module) {
    module.exports = 8;
}

export default 7;

// @filename: /main.ts
import DefaultExport from "./dep.js";

export default 3 + DefaultExport;
