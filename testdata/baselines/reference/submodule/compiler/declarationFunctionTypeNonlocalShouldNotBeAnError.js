//// [tests/cases/compiler/declarationFunctionTypeNonlocalShouldNotBeAnError.ts] ////

//// [declarationFunctionTypeNonlocalShouldNotBeAnError.ts]
namespace foo {
    function bar(): void {}

    export const obj = {
        bar
    }
}


//// [declarationFunctionTypeNonlocalShouldNotBeAnError.js]
"use strict";
var foo;
(function (foo) {
    function bar() { }
    foo.obj = {
        bar
    };
})(foo || (foo = {}));


//// [declarationFunctionTypeNonlocalShouldNotBeAnError.d.ts]
namespace foo {
    function bar(): void;
    export const obj: {
        bar: typeof bar;
    };
    export {};
}


//// [DtsFileErrors]


declarationFunctionTypeNonlocalShouldNotBeAnError.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationFunctionTypeNonlocalShouldNotBeAnError.d.ts (1 errors) ====
    namespace foo {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        function bar(): void;
        export const obj: {
            bar: typeof bar;
        };
        export {};
    }
    