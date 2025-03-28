//// [tests/cases/compiler/declarationFileNoCrashOnExtraExportModifier.ts] ////

//// [input.ts]
export = exports;
declare class exports {
    constructor(p: number);
    t: number;
}
export class Sub {
    instance!: {
        t: number;
    };
}
declare namespace exports {
    export { Sub };
}

//// [input.js]
"use strict";
exports.Sub = void 0;
class Sub {
    instance;
}
module.exports = exports;
