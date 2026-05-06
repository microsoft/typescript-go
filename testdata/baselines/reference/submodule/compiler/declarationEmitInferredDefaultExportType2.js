//// [tests/cases/compiler/declarationEmitInferredDefaultExportType2.ts] ////

//// [declarationEmitInferredDefaultExportType2.ts]
// test.ts
export = {
  foo: [],
  bar: undefined,
  baz: null
}

//// [declarationEmitInferredDefaultExportType2.js]
"use strict";
module.exports = {
    foo: [],
    bar: undefined,
    baz: null
};


//// [declarationEmitInferredDefaultExportType2.d.ts]
const _default: {
    foo: never[];
    bar: undefined;
    baz: null;
};
export = _default;


//// [DtsFileErrors]


declarationEmitInferredDefaultExportType2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitInferredDefaultExportType2.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo: never[];
        bar: undefined;
        baz: null;
    };
    export = _default;
    