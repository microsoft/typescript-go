// @module: commonjs
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
exports.greeting = "hello";

function setup(module) {
    module.exports = { configured: true };
}

const fs = require("fs");

// @filename: /main.ts
import dep = require("./dep.js");

export const g = dep.greeting;
