//// [tests/cases/compiler/declarationEmitDefaultExport8.ts] ////

//// [declarationEmitDefaultExport8.ts]
var _default = 1;
export {_default as d}
export default 1 + 2;


//// [declarationEmitDefaultExport8.js]
var _default = 1;
export { _default as d };
export default 1 + 2;


//// [declarationEmitDefaultExport8.d.ts]
var _default: number;
export { _default as d };
const _default_1: number;
export default _default_1;


//// [DtsFileErrors]


declarationEmitDefaultExport8.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDefaultExport8.d.ts (1 errors) ====
    var _default: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export { _default as d };
    const _default_1: number;
    export default _default_1;
    