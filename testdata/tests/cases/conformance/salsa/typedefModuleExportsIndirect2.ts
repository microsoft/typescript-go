// @filename: typedefModuleExportsIndirect2.js
// @checkJs: true
// @strict: true
// @outdir: dist
// @declaration: true
/** @typedef {{ a: 1, m: 1 }} C */
const f = function() {};
module.exports = f;
// @filename: use.js
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c
