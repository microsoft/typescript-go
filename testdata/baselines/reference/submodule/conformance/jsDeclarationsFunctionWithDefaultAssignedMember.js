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
declare namespace foo {
    const foo: typeof import(".");
}
declare namespace foo {
    const default_1: typeof import(".");
    export { default_1 as default };
}
export = foo;
