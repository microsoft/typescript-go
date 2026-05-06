//// [tests/cases/compiler/declarationEmitIndexTypeArray.ts] ////

//// [declarationEmitIndexTypeArray.ts]
function doSomethingWithKeys<T>(...keys: (keyof T)[]) { }

const utilityFunctions = {
  doSomethingWithKeys
};


//// [declarationEmitIndexTypeArray.js]
"use strict";
function doSomethingWithKeys(...keys) { }
const utilityFunctions = {
    doSomethingWithKeys
};


//// [declarationEmitIndexTypeArray.d.ts]
function doSomethingWithKeys<T>(...keys: (keyof T)[]): void;
const utilityFunctions: {
    doSomethingWithKeys: typeof doSomethingWithKeys;
};


//// [DtsFileErrors]


declarationEmitIndexTypeArray.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitIndexTypeArray.d.ts (1 errors) ====
    function doSomethingWithKeys<T>(...keys: (keyof T)[]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const utilityFunctions: {
        doSomethingWithKeys: typeof doSomethingWithKeys;
    };
    