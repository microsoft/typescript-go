// @module: nodenext
// @allowJs: true
// @checkJs: true
// @noEmit: true
// @strict: true

// @filename: /dep.js
function something(module) {
    module.exports = 8;
}

export const foo = 7;

// @filename: /main.ts
import { foo } from "./dep.js";

export const bar = foo + 1;
