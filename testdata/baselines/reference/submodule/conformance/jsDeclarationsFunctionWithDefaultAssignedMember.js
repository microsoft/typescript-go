//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionWithDefaultAssignedMember.ts] ////

//// [index.js]
function foo() {}

foo.foo = foo;
foo.default = foo;
module.exports = foo;

//// [index.js]
"use strict";
function foo() { }
foo.foo = foo;
foo.default = foo;
module.exports = foo;


//// [index.d.ts]
function foo(): void;
declare namespace foo {
    var foo: typeof import(".");
}
declare namespace foo {
    var _a: typeof import(".");
    export { _a as default };
}
export = foo;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    function foo(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    declare namespace foo {
        var foo: typeof import(".");
    }
    declare namespace foo {
        var _a: typeof import(".");
        export { _a as default };
    }
    export = foo;
    