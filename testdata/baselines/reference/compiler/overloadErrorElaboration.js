//// [tests/cases/compiler/overloadErrorElaboration.ts] ////

//// [overloadErrorElaboration.ts]
// Two overloads
function foo(bar: string): string;
function foo(bar: number): number;
function foo(bar: any): any {
    return bar;
}

var x = foo(true);


//// [overloadErrorElaboration.js]
"use strict";
function foo(bar) {
    return bar;
}
var x = foo(true);
