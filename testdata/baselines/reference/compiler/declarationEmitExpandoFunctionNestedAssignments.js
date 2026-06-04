//// [tests/cases/compiler/declarationEmitExpandoFunctionNestedAssignments.ts] ////

//// [declarationEmitExpandoFunctionNestedAssignments.ts]
function Foo(): void { }
Foo.top = 1;
let d: number = (Foo.inInitializer = 2);
if (true) {
    Foo.inBlock = 3;
}


//// [declarationEmitExpandoFunctionNestedAssignments.js]
"use strict";
function Foo() { }
Foo.top = 1;
let d = (Foo.inInitializer = 2);
if (true) {
    Foo.inBlock = 3;
}


//// [declarationEmitExpandoFunctionNestedAssignments.d.ts]
declare function Foo(): void;
declare namespace Foo {
    var top: number;
    var inInitializer: number;
    var inBlock: number;
}
declare let d: number;
