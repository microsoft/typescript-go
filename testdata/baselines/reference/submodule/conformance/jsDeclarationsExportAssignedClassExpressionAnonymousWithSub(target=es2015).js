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
"use strict";
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
};
module.exports.Sub = class {
    constructor() {
        this.instance = new module.exports(10);
    }
};


//// [index.d.ts]
export = _default;
declare class _default {
    t: number;
    /**
     * @param {number} p
     */
    constructor(p: number);
}
declare namespace _ns {
    export class _class {
        instance: import(".");
        constructor();
    }
}
import _Sub = _ns._class;
declare namespace _default {
    export { _Sub as Sub };
}
