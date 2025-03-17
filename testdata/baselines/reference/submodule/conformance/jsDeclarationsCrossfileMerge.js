//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsCrossfileMerge.ts] ////

//// [index.js]
const m = require("./exporter");

module.exports = m.default;
module.exports.memberName = "thing";

//// [exporter.js]
function validate() {}

export default validate;


//// [index.js]
const m = require("./exporter");
module.exports = m.default;
module.exports.memberName = "thing";
//// [exporter.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function validate() { }
exports.default = validate;
