//// [tests/cases/compiler/declFileTypeofFunction.ts] ////

//// [declFileTypeofFunction.ts]
function f(n: typeof f): string;
function f(n: typeof g): string;
function f() { return undefined; }
function g(n: typeof g): number;
function g(n: typeof f): number;
function g() { return undefined; }

var b: () => typeof b;

function b1() {
    return b1;
}

function foo(): typeof foo {
    return null;
}
var foo1: typeof foo;
var foo2 = foo;

var foo3 = function () {
    return foo3;
}
var x = () => {
    return x;
}

function foo5(x: number) {
    function bar(x: number) {
        return x;
    }
    return bar;
}

//// [declFileTypeofFunction.js]
"use strict";
function f() { return undefined; }
function g() { return undefined; }
var b;
function b1() {
    return b1;
}
function foo() {
    return null;
}
var foo1;
var foo2 = foo;
var foo3 = function () {
    return foo3;
};
var x = () => {
    return x;
};
function foo5(x) {
    function bar(x) {
        return x;
    }
    return bar;
}


//// [declFileTypeofFunction.d.ts]
function f(n: typeof f): string;
function f(n: typeof g): string;
function g(n: typeof g): number;
function g(n: typeof f): number;
var b: () => typeof b;
function b1(): typeof b1;
function foo(): typeof foo;
var foo1: typeof foo;
var foo2: typeof foo;
var foo3: () => () => /*elided*/ any;
var x: () => () => /*elided*/ any;
function foo5(x: number): (x: number) => number;


//// [DtsFileErrors]


declFileTypeofFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeofFunction.d.ts (1 errors) ====
    function f(n: typeof f): string;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f(n: typeof g): string;
    function g(n: typeof g): number;
    function g(n: typeof f): number;
    var b: () => typeof b;
    function b1(): typeof b1;
    function foo(): typeof foo;
    var foo1: typeof foo;
    var foo2: typeof foo;
    var foo3: () => () => /*elided*/ any;
    var x: () => () => /*elided*/ any;
    function foo5(x: number): (x: number) => number;
    