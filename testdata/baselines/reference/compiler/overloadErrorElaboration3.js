//// [tests/cases/compiler/overloadErrorElaboration3.ts] ////

//// [overloadErrorElaboration3.ts]
// Three overloads — should show per-overload errors
function bar(x: string): string;
function bar(x: number): number;
function bar(x: boolean): boolean;
function bar(x: any): any {
    return x;
}

var y = bar({ a: 1 });


//// [overloadErrorElaboration3.js]
"use strict";
function bar(x) {
    return x;
}
var y = bar({ a: 1 });
