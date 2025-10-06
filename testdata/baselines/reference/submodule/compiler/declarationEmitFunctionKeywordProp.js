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
    const _a: boolean;
    export { _a as null };
}
declare function bar(): void;
declare namespace bar {
    const async: boolean;
}
declare namespace bar {
    const normal: boolean;
}
declare function baz(): void;
declare namespace baz {
    const _b: boolean;
    export { _b as class };
}
declare namespace baz {
    const normal: boolean;
}
