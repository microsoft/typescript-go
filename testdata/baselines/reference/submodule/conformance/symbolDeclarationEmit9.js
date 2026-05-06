//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit9.ts] ////

//// [symbolDeclarationEmit9.ts]
var obj = {
    [Symbol.isConcatSpreadable]() { }
}

//// [symbolDeclarationEmit9.js]
"use strict";
var obj = {
    [Symbol.isConcatSpreadable]() { }
};


//// [symbolDeclarationEmit9.d.ts]
var obj: {
    [Symbol.isConcatSpreadable](): void;
};


//// [DtsFileErrors]


symbolDeclarationEmit9.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit9.d.ts (1 errors) ====
    var obj: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.isConcatSpreadable](): void;
    };
    