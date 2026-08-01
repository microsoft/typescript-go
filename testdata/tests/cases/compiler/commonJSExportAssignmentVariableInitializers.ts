// @target: es2015
// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true

// @filename: requires.d.ts
declare var module: { exports: any };
declare function require(name: string): any;

// @filename: uninitialized.js
var uninitialized;
module.exports = uninitialized;

// @filename: object.js
const object = { value: 1 };
module.exports = object;
module.exports.helper = function helper() {};

// @filename: index.js
require("./uninitialized");
const object = require("./object");
object.helper();
