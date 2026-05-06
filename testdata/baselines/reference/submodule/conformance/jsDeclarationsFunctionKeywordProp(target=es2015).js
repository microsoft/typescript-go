//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionKeywordProp.ts] ////

//// [source.js]
function foo() {}
foo.null = true;

function bar() {}
bar.async = true;
bar.normal = false;

function baz() {}
baz.class = true;
baz.normal = false;

//// [source.js]
"use strict";
function foo() { }
foo.null = true;
function bar() { }
bar.async = true;
bar.normal = false;
function baz() { }
baz.class = true;
baz.normal = false;


//// [source.d.ts]
function foo(): void;
declare namespace foo {
    var _a: boolean;
    export { _a as null };
}
function bar(): void;
declare namespace bar {
    var async: boolean;
}
declare namespace bar {
    var normal: boolean;
}
function baz(): void;
declare namespace baz {
    var _b: boolean;
    export { _b as class };
}
declare namespace baz {
    var normal: boolean;
}


//// [DtsFileErrors]


out/source.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/source.d.ts (1 errors) ====
    function foo(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    declare namespace foo {
        var _a: boolean;
        export { _a as null };
    }
    function bar(): void;
    declare namespace bar {
        var async: boolean;
    }
    declare namespace bar {
        var normal: boolean;
    }
    function baz(): void;
    declare namespace baz {
        var _b: boolean;
        export { _b as class };
    }
    declare namespace baz {
        var normal: boolean;
    }
    