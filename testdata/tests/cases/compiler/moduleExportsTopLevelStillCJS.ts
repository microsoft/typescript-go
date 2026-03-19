// @module: commonjs
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
module.exports = 8;

// @filename: /main.ts
import dep = require("./dep.js");

export const result = dep + 1;
