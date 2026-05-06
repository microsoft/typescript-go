//// [tests/cases/conformance/declarationEmit/exportDefaultExpressionComments.ts] ////

//// [exportDefaultExpressionComments.ts]
/**
 * JSDoc Comments
 */
export default null


//// [exportDefaultExpressionComments.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * JSDoc Comments
 */
exports.default = null;


//// [exportDefaultExpressionComments.d.ts]
/**
 * JSDoc Comments
 */
const _default: null;
export default _default;


//// [DtsFileErrors]


exportDefaultExpressionComments.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== exportDefaultExpressionComments.d.ts (1 errors) ====
    /**
     * JSDoc Comments
     */
    const _default: null;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    