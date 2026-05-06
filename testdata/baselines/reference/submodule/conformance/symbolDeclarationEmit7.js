//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit7.ts] ////

//// [symbolDeclarationEmit7.ts]
var obj: {
    [Symbol.isConcatSpreadable]: string;
}

//// [symbolDeclarationEmit7.js]
"use strict";
var obj;


//// [symbolDeclarationEmit7.d.ts]
var obj: {
    [Symbol.isConcatSpreadable]: string;
};


//// [DtsFileErrors]


symbolDeclarationEmit7.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit7.d.ts (1 errors) ====
    var obj: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.isConcatSpreadable]: string;
    };
    