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
class A {
}
var a = (x) => x;
function a2(x) { return x; }
var a3 = (x) => new A();
function a4(x) { return new A(); }
var b = (x) => x;
function b2(x) { return x; }


//// [declarationEmitTypeParameterNameInOuterScope.d.ts]
declare class A {
}
declare var a: <A>(x: A) => A;
declare function a2<A>(x: A): A;
declare var a3: <A>(x: A) => globalThis.A;
declare function a4<A>(x: A): globalThis.A;
interface B {
}
declare var b: <B>(x: B) => B;
declare function b2<B>(x: B): B;
