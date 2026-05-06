//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedConstructorFunctionWithSub.ts] ////

//// [jsDeclarationsExportAssignedConstructorFunctionWithSub.js]
/**
 * @param {number} p
 */
module.exports = function (p) {
    this.t = 12 + p;
}
module.exports.Sub = function() {
    this.instance = new module.exports(10);
}
module.exports.Sub.prototype = { }


//// [jsDeclarationsExportAssignedConstructorFunctionWithSub.js]
"use strict";
/**
 * @param {number} p
 */
module.exports = function (p) {
    this.t = 12 + p;
};
module.exports.Sub = function () {
    this.instance = new module.exports(10);
};
module.exports.Sub.prototype = {};


//// [jsDeclarationsExportAssignedConstructorFunctionWithSub.d.ts]
/**
 * @param {number} p
 */
const _default: (p: number) => void;
export = _default;
export var Sub: () => void;
