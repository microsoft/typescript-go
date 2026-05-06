//// [tests/cases/compiler/declarationEmitTypeAliasWithTypeParameters3.ts] ////

//// [declarationEmitTypeAliasWithTypeParameters3.ts]
type Foo<T> = {
    foo<U>(): Foo<U>
};
function bar() {
    return {} as Foo<number>;
}


//// [declarationEmitTypeAliasWithTypeParameters3.js]
"use strict";
function bar() {
    return {};
}


//// [declarationEmitTypeAliasWithTypeParameters3.d.ts]
type Foo<T> = {
    foo<U>(): Foo<U>;
};
function bar(): Foo<number>;


//// [DtsFileErrors]


declarationEmitTypeAliasWithTypeParameters3.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTypeAliasWithTypeParameters3.d.ts (1 errors) ====
    type Foo<T> = {
        foo<U>(): Foo<U>;
    };
    function bar(): Foo<number>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    