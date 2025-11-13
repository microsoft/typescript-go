// @filename: controlFlowJSClassProperty.js
// @checkJs: true
// @strict: true
// @outdir: dist
// @declaration: true
/** @typedef {{ a: 1, m: 1 }} C */
const dummy = 0;
module.exports = dummy;
// @filename: use.js
/** @typedef {import('./controlFlowJSClassProperty').C} C */
/** @type {C} */
var c
