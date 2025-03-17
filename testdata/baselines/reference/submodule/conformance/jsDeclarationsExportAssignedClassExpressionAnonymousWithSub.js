//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassExpressionAnonymousWithSub.ts] ////

//// [index.js]
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
}
module.exports.Sub = class {
    constructor() {
        this.instance = new module.exports(10);
    }
}


//// [index.js]
module.exports = class {
    constructor(p) {
        this.t = 12 + p;
    }
};
module.exports.Sub = class {
    constructor() {
        this.instance = new module.exports(10);
    }
};
