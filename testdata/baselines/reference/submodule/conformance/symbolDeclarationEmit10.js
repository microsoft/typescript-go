//// [tests/cases/conformance/es6/Symbols/symbolDeclarationEmit10.ts] ////

//// [symbolDeclarationEmit10.ts]
var obj = {
    get [Symbol.isConcatSpreadable]() { return '' },
    set [Symbol.isConcatSpreadable](x) { }
}

//// [symbolDeclarationEmit10.js]
"use strict";
var obj = {
    get [Symbol.isConcatSpreadable]() { return ''; },
    set [Symbol.isConcatSpreadable](x) { }
};


//// [symbolDeclarationEmit10.d.ts]
var obj: {
    [Symbol.isConcatSpreadable]: string;
};


//// [DtsFileErrors]


symbolDeclarationEmit10.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolDeclarationEmit10.d.ts (1 errors) ====
    var obj: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.isConcatSpreadable]: string;
    };
    