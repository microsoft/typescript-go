//// [tests/cases/compiler/jsDeclarationExportDefaultAssignmentCrash.ts] ////

//// [index.js]
exports.default = () => {
    return 1234;
}


//// [index.js]
"use strict";
exports.default = () => {
    return 1234;
};


//// [index.d.ts]
const _default: () => number;
export default _default;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    const _default: () => number;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    