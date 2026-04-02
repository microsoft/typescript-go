// @module: commonjs
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
module.exports = { real: true };

function setup(module) {
    module.exports = { configured: true };
}

// @filename: /main.ts
import dep = require("./dep.js");

export const r = dep.real;
