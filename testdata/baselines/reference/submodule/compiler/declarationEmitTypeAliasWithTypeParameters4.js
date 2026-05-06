//// [tests/cases/compiler/declarationEmitTypeAliasWithTypeParameters4.ts] ////

//// [declarationEmitTypeAliasWithTypeParameters4.ts]
type Foo<T, Y> = {
    foo<U, J>(): Foo<U, J>
};
type SubFoo<R> = Foo<string, R>;

function foo() {
    return {} as SubFoo<number>;
}


//// [declarationEmitTypeAliasWithTypeParameters4.js]
"use strict";
function foo() {
    return {};
}


//// [declarationEmitTypeAliasWithTypeParameters4.d.ts]
type Foo<T, Y> = {
    foo<U, J>(): Foo<U, J>;
};
type SubFoo<R> = Foo<string, R>;
function foo(): SubFoo<number>;


//// [DtsFileErrors]


declarationEmitTypeAliasWithTypeParameters4.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTypeAliasWithTypeParameters4.d.ts (1 errors) ====
    type Foo<T, Y> = {
        foo<U, J>(): Foo<U, J>;
    };
    type SubFoo<R> = Foo<string, R>;
    function foo(): SubFoo<number>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    