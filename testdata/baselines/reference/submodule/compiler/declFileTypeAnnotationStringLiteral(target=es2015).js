//// [tests/cases/compiler/declFileTypeAnnotationStringLiteral.ts] ////

//// [declFileTypeAnnotationStringLiteral.ts]
function foo(a: "hello"): number;
function foo(a: "name"): string;
function foo(a: string): string | number;
function foo(a: string): string | number {
    if (a === "hello") {
        return a.length;
    }

    return a;
}

//// [declFileTypeAnnotationStringLiteral.js]
"use strict";
function foo(a) {
    if (a === "hello") {
        return a.length;
    }
    return a;
}


//// [declFileTypeAnnotationStringLiteral.d.ts]
function foo(a: "hello"): number;
function foo(a: "name"): string;
function foo(a: string): string | number;


//// [DtsFileErrors]


declFileTypeAnnotationStringLiteral.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeAnnotationStringLiteral.d.ts (1 errors) ====
    function foo(a: "hello"): number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function foo(a: "name"): string;
    function foo(a: string): string | number;
    