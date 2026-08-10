// @allowJs: true
// @checkJs: true
// @target: esnext
// @module: preserve
// @outDir: ./out
// @declaration: true
// @emitDeclarationOnly: true

// @filename: namedExports.js
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

// @filename: inlineExports.js
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

// @filename: defaultExport.js
/**
 * jsdoc description on `ExpandoArrowDefaultExport`
 */
const ExpandoArrowDefaultExport = () => { return null; };
ExpandoArrowDefaultExport.args = { num: 3 };

export default ExpandoArrowDefaultExport;
