//// [tests/cases/compiler/functionDeclarationWithArgumentOfTypeFunctionTypeArray.ts] ////

//// [functionDeclarationWithArgumentOfTypeFunctionTypeArray.ts]
function foo(args: { (x): number }[]) {
    return args.length;
}


//// [functionDeclarationWithArgumentOfTypeFunctionTypeArray.js]
"use strict";
function foo(args) {
    return args.length;
}


//// [functionDeclarationWithArgumentOfTypeFunctionTypeArray.d.ts]
function foo(args: {
    (x: any): number;
}[]): number;


//// [DtsFileErrors]


functionDeclarationWithArgumentOfTypeFunctionTypeArray.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== functionDeclarationWithArgumentOfTypeFunctionTypeArray.d.ts (1 errors) ====
    function foo(args: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        (x: any): number;
    }[]): number;
    