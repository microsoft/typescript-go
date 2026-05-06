//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassExpressionAnonymous.ts] ////

//// [index.js]
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
}

//// [index.js]
"use strict";
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
};


//// [index.d.ts]
const _default: {
    new (p: number): import(".");
};
export = _default;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/index.d.ts(2,22): error TS1340: Module '.' does not refer to a type, but is used as a type here. Did you mean 'typeof import('.')'?


==== out/index.d.ts (2 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new (p: number): import(".");
                         ~~~~~~~~~~~
!!! error TS1340: Module '.' does not refer to a type, but is used as a type here. Did you mean 'typeof import('.')'?
    };
    export = _default;
    