//// [tests/cases/conformance/dynamicImport/importCallExpressionDeclarationEmit2.ts] ////

//// [0.ts]
export function foo() { return "foo"; }

//// [1.ts]
var p1 = import("./0");


//// [0.js]
export function foo() { return "foo"; }
//// [1.js]
"use strict";
var p1 = import("./0");


//// [0.d.ts]
export function foo(): string;
//// [1.d.ts]
var p1: Promise<typeof import("./0")>;


//// [DtsFileErrors]


1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== 0.d.ts (0 errors) ====
    export function foo(): string;
    
==== 1.d.ts (1 errors) ====
    var p1: Promise<typeof import("./0")>;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    