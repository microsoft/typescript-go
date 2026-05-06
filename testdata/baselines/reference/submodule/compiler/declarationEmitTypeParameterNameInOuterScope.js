//// [tests/cases/compiler/declarationEmitTypeParameterNameInOuterScope.ts] ////

//// [declarationEmitTypeParameterNameInOuterScope.ts]
class A { }

var a = <A,>(x: A) => x;
function a2<A,>(x: A) { return x }

var a3 = <A,>(x: A) => new A();
function a4<A,>(x: A) { return new A() }


interface B { }

var b = <B,>(x: B) => x;
function b2<B,>(x: B) { return x }


//// [declarationEmitTypeParameterNameInOuterScope.js]
"use strict";
class A {
}
var a = (x) => x;
function a2(x) { return x; }
var a3 = (x) => new A();
function a4(x) { return new A(); }
var b = (x) => x;
function b2(x) { return x; }


//// [declarationEmitTypeParameterNameInOuterScope.d.ts]
class A {
}
var a: <A>(x: A) => A;
function a2<A>(x: A): A;
var a3: <A>(x: A) => globalThis.A;
function a4<A>(x: A): globalThis.A;
interface B {
}
var b: <B>(x: B) => B;
function b2<B>(x: B): B;


//// [DtsFileErrors]


declarationEmitTypeParameterNameInOuterScope.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTypeParameterNameInOuterScope.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    var a: <A>(x: A) => A;
    function a2<A>(x: A): A;
    var a3: <A>(x: A) => globalThis.A;
    function a4<A>(x: A): globalThis.A;
    interface B {
    }
    var b: <B>(x: B) => B;
    function b2<B>(x: B): B;
    