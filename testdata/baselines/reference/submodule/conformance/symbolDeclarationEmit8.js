//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit8.ts] ////

//// [symbolDeclarationEmit8.ts]
var obj = {
    [Symbol.isConcatSpreadable]: 0
}

//// [symbolDeclarationEmit8.js]
"use strict";
var obj = {
    [Symbol.isConcatSpreadable]: 0
};


//// [symbolDeclarationEmit8.d.ts]
var obj: {
    [Symbol.isConcatSpreadable]: number;
};


//// [DtsFileErrors]


symbolDeclarationEmit8.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit8.d.ts (1 errors) ====
    var obj: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.isConcatSpreadable]: number;
    };
    