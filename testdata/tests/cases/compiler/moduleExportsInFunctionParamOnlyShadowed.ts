// @module: nodenext
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
function setup(module) {
    module.exports = { configured: true };
}

export default 7;

// @filename: /main.ts
import DefaultExport from "./dep.js";

export default 3 + DefaultExport;
