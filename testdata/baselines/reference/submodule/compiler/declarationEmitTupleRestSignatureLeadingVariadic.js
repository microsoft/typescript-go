//// [tests/cases/compiler/declarationEmitTupleRestSignatureLeadingVariadic.ts] ////

//// [declarationEmitTupleRestSignatureLeadingVariadic.ts]
const f = <TFirstArgs extends any[], TLastArg>(...args: [...TFirstArgs, TLastArg]): void => {};

//// [declarationEmitTupleRestSignatureLeadingVariadic.js]
"use strict";
const f = (...args) => { };


//// [declarationEmitTupleRestSignatureLeadingVariadic.d.ts]
const f: <TFirstArgs extends any[], TLastArg>(...args: [...TFirstArgs, TLastArg]) => void;


//// [DtsFileErrors]


declarationEmitTupleRestSignatureLeadingVariadic.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTupleRestSignatureLeadingVariadic.d.ts (1 errors) ====
    const f: <TFirstArgs extends any[], TLastArg>(...args: [...TFirstArgs, TLastArg]) => void;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    