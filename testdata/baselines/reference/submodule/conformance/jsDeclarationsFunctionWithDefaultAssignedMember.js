//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionWithDefaultAssignedMember.ts] ////

//// [index.js]
function foo() {}

foo.foo = foo;
foo.default = foo;
module.exports = foo;

//// [index.js]
function foo() { }
foo.foo = foo;
foo.default = foo;
export = foo;
module.exports = foo;


//// [index.d.ts]
declare function foo(): void;
export = foo;
declare namespace foo {
    var foo: typeof foo;
    var default_1: typeof foo;
    export { default_1 as default };
}


//// [DtsFileErrors]


out/index.d.ts(4,9): error TS2502: 'foo' is referenced directly or indirectly in its own type annotation.


==== out/index.d.ts (1 errors) ====
    declare function foo(): void;
    export = foo;
    declare namespace foo {
        var foo: typeof foo;
            ~~~
!!! error TS2502: 'foo' is referenced directly or indirectly in its own type annotation.
        var default_1: typeof foo;
        export { default_1 as default };
    }
    