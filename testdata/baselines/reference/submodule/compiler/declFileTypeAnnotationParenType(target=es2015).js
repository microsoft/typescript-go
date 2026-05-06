//// [tests/cases/compiler/declFileTypeAnnotationParenType.ts] ////

//// [declFileTypeAnnotationParenType.ts]
class c {
    private p: string;
}

var x: (() => c)[] = [() => new c()];
var y = [() => new c()];

var k: (() => c) | string = (() => new c()) || "";
var l = (() => new c()) || "";

//// [declFileTypeAnnotationParenType.js]
"use strict";
class c {
}
var x = [() => new c()];
var y = [() => new c()];
var k = (() => new c()) || "";
var l = (() => new c()) || "";


//// [declFileTypeAnnotationParenType.d.ts]
class c {
    private p;
}
var x: (() => c)[];
var y: (() => c)[];
var k: (() => c) | string;
var l: () => c;
