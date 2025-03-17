//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassExpression.ts] ////

//// [index.js]
module.exports = class Thing {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
}

//// [index.js]
module.exports = class Thing {
    constructor(p) {
        this.t = 12 + p;
    }
};
