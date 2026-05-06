//// [tests/cases/compiler/declarationEmitTypeAliasWithTypeParameters6.ts] ////

//// [declarationEmitTypeAliasWithTypeParameters6.ts]
type Foo<T, Y> = {
    foo<U, J>(): Foo<U, J>
};
type SubFoo<R, S> = Foo<S, R>;

function foo() {
    return {} as SubFoo<number, string>;
}


//// [declarationEmitTypeAliasWithTypeParameters6.js]
"use strict";
function foo() {
    return {};
}


//// [declarationEmitTypeAliasWithTypeParameters6.d.ts]
type Foo<T, Y> = {
    foo<U, J>(): Foo<U, J>;
};
type SubFoo<R, S> = Foo<S, R>;
function foo(): SubFoo<number, string>;


//// [DtsFileErrors]


declarationEmitTypeAliasWithTypeParameters6.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTypeAliasWithTypeParameters6.d.ts (1 errors) ====
    type Foo<T, Y> = {
        foo<U, J>(): Foo<U, J>;
    };
    type SubFoo<R, S> = Foo<S, R>;
    function foo(): SubFoo<number, string>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    