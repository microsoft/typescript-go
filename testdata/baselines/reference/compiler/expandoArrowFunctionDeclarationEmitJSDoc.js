//// [tests/cases/compiler/expandoArrowFunctionDeclarationEmitJSDoc.ts] ////

//// [namedExports.js]
// An expando host declared as a function declaration keeps its JSDoc in the
// declaration file; one declared as an arrow function or function expression
// should keep it too, since both are emitted as a `declare function`.

/**
 * jsdoc description on `ExpandoFunctionNamedExport`
 */
function ExpandoFunctionNamedExport() { return null; }
ExpandoFunctionNamedExport.args = { num: 3 };

/**
 * jsdoc description on `ExpandoArrowNamedExport`
 */
const ExpandoArrowNamedExport = () => { return null; };
ExpandoArrowNamedExport.args = { num: 3 };

/**
 * jsdoc description on `ExpandoFunctionExpressionNamedExport`
 */
const ExpandoFunctionExpressionNamedExport = function () { return null; };
ExpandoFunctionExpressionNamedExport.args = { num: 3 };

export { ExpandoFunctionNamedExport, ExpandoArrowNamedExport, ExpandoFunctionExpressionNamedExport };

//// [inlineExports.js]
/**
 * jsdoc description on `ExpandoFunctionInlineNamedExport`
 */
export function ExpandoFunctionInlineNamedExport() { return null; }
ExpandoFunctionInlineNamedExport.args = { num: 3 };

/**
 * jsdoc description on `ExpandoArrowInlineNamedExport`
 */
export const ExpandoArrowInlineNamedExport = () => { return null; };
ExpandoArrowInlineNamedExport.args = { num: 3 };

//// [defaultExport.js]
/**
 * jsdoc description on `ExpandoArrowDefaultExport`
 */
const ExpandoArrowDefaultExport = () => { return null; };
ExpandoArrowDefaultExport.args = { num: 3 };

export default ExpandoArrowDefaultExport;




//// [namedExports.d.ts]
/**
 * jsdoc description on `ExpandoFunctionNamedExport`
 */
declare function ExpandoFunctionNamedExport(): null;
declare namespace ExpandoFunctionNamedExport {
    var args: {
        num: number;
    };
}
/**
 * jsdoc description on `ExpandoArrowNamedExport`
 */
declare function ExpandoArrowNamedExport(): null;
declare namespace ExpandoArrowNamedExport {
    var args: {
        num: number;
    };
}
/**
 * jsdoc description on `ExpandoFunctionExpressionNamedExport`
 */
declare function ExpandoFunctionExpressionNamedExport(): null;
declare namespace ExpandoFunctionExpressionNamedExport {
    var args: {
        num: number;
    };
}
export { ExpandoFunctionNamedExport, ExpandoArrowNamedExport, ExpandoFunctionExpressionNamedExport };
//// [inlineExports.d.ts]
/**
 * jsdoc description on `ExpandoFunctionInlineNamedExport`
 */
export declare function ExpandoFunctionInlineNamedExport(): null;
export declare namespace ExpandoFunctionInlineNamedExport {
    var args: {
        num: number;
    };
}
/**
 * jsdoc description on `ExpandoArrowInlineNamedExport`
 */
export declare function ExpandoArrowInlineNamedExport(): null;
export declare namespace ExpandoArrowInlineNamedExport {
    var args: {
        num: number;
    };
}
//// [defaultExport.d.ts]
/**
 * jsdoc description on `ExpandoArrowDefaultExport`
 */
declare function ExpandoArrowDefaultExport(): null;
declare namespace ExpandoArrowDefaultExport {
    var args: {
        num: number;
    };
}
export default ExpandoArrowDefaultExport;
