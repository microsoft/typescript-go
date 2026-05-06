//// [tests/cases/compiler/es6ExportDefaultExpression.ts] ////

//// [es6ExportDefaultExpression.ts]
export default (1 + 2);


//// [es6ExportDefaultExpression.js]
export default (1 + 2);


//// [es6ExportDefaultExpression.d.ts]
const _default: number;
export default _default;


//// [DtsFileErrors]


es6ExportDefaultExpression.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es6ExportDefaultExpression.d.ts (1 errors) ====
    const _default: number;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    