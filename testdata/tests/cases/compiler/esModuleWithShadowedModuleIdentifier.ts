// @allowJs: true
// @noEmit: true

// https://github.com/microsoft/typescript-go/issues/2656

// @filename: dep.js
function something(module) {
    module.exports = 8;
}

export default 7;

// @filename: main.ts
import DefaultExport from "./dep.js";

export default 3 + DefaultExport;
