//// [tests/cases/compiler/declarationEmitInferredTypeAlias8.ts] ////

//// [declarationEmitInferredTypeAlias8.ts]
type Foo<T> = T | { x: Foo<T> };
var x: Foo<number[]>;

function returnSomeGlobalValue() {
    return x;
}

//// [declarationEmitInferredTypeAlias8.js]
"use strict";
var x;
function returnSomeGlobalValue() {
    return x;
}


//// [declarationEmitInferredTypeAlias8.d.ts]
type Foo<T> = T | {
    x: Foo<T>;
};
var x: Foo<number[]>;
function returnSomeGlobalValue(): Foo<number[]>;


//// [DtsFileErrors]


declarationEmitInferredTypeAlias8.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitInferredTypeAlias8.d.ts (1 errors) ====
    type Foo<T> = T | {
        x: Foo<T>;
    };
    var x: Foo<number[]>;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function returnSomeGlobalValue(): Foo<number[]>;
    