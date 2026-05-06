//// [tests/cases/compiler/es5ExportDefaultExpression.ts] ////

//// [es5ExportDefaultExpression.ts]
export default (1 + 2);


//// [es5ExportDefaultExpression.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = (1 + 2);


//// [es5ExportDefaultExpression.d.ts]
const _default: number;
export default _default;


//// [DtsFileErrors]


es5ExportDefaultExpression.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es5ExportDefaultExpression.d.ts (1 errors) ====
    const _default: number;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    