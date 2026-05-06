//// [tests/cases/compiler/genericClasses0.ts] ////

//// [genericClasses0.ts]
class C<T> {
	public x: T;
}

var v1 : C<string>;

var y = v1.x; // should be 'string'

//// [genericClasses0.js]
"use strict";
class C {
}
var v1;
var y = v1.x; // should be 'string'


//// [genericClasses0.d.ts]
class C<T> {
    x: T;
}
var v1: C<string>;
var y: string;
