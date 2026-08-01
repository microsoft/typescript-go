// @target: es2015
// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true

// @filename: requires.d.ts
declare var module: { exports: any };
declare function require(name: string): any;

// @filename: mod.js
module.exports = function main() {};
/** @param {number} value */
module.exports.helper = function helper(value) {};

// @filename: plugin.js
class Plugin {}
module.exports = Plugin;
module.exports.helper = function helper() {};

// @filename: defined.js
module.exports = class Defined {};
Object.defineProperty(module.exports, "value", { value: 1 });

// @filename: arrow.js
const arrow = () => {};
module.exports = arrow;
module.exports.helper = function helper() {};

// @filename: index.js
const mod = require("./mod");
const Plugin = require("./plugin");
const Defined = require("./defined");
const arrow = require("./arrow");

mod();
mod.helper(1);
Plugin.helper();
Defined.value.toFixed();
arrow.helper();
