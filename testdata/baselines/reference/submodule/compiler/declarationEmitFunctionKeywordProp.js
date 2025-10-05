//// [tests/cases/compiler/declarationEmitFunctionKeywordProp.ts] ////

//// [declarationEmitFunctionKeywordProp.ts]
function foo() {}
foo.null = true;

function bar() {}
bar.async = true;
bar.normal = false;

function baz() {}
baz.class = true;
baz.normal = false;

//// [declarationEmitFunctionKeywordProp.js]
function foo() { }
foo.null = true;
function bar() { }
bar.async = true;
bar.normal = false;
function baz() { }
baz.class = true;
baz.normal = false;


//// [declarationEmitFunctionKeywordProp.d.ts]
declare function foo(): void;
declare namespace foo {
    const null_1: true;
    export { null_1 as null };
}
declare function bar(): void;
declare namespace bar {
    const async: true;
}
declare namespace bar {
    const normal: false;
}
declare function baz(): void;
declare namespace baz {
    const class_1: true;
    export { class_1 as class };
}
declare namespace baz {
    const normal: false;
}
