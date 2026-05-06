//// [tests/cases/compiler/declarationEmitInferredDefaultExportType.ts] ////

//// [declarationEmitInferredDefaultExportType.ts]
// test.ts
export default {
  foo: [],
  bar: undefined,
  baz: null
}

//// [declarationEmitInferredDefaultExportType.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// test.ts
exports.default = {
    foo: [],
    bar: undefined,
    baz: null
};


//// [declarationEmitInferredDefaultExportType.d.ts]
const _default: {
    foo: never[];
    bar: undefined;
    baz: null;
};
export default _default;


//// [DtsFileErrors]


declarationEmitInferredDefaultExportType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitInferredDefaultExportType.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo: never[];
        bar: undefined;
        baz: null;
    };
    export default _default;
    