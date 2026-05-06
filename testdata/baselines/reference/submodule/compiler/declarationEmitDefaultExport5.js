//// [tests/cases/compiler/declarationEmitDefaultExport5.ts] ////

//// [declarationEmitDefaultExport5.ts]
export default 1 + 2;


//// [declarationEmitDefaultExport5.js]
export default 1 + 2;


//// [declarationEmitDefaultExport5.d.ts]
const _default: number;
export default _default;


//// [DtsFileErrors]


declarationEmitDefaultExport5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDefaultExport5.d.ts (1 errors) ====
    const _default: number;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    