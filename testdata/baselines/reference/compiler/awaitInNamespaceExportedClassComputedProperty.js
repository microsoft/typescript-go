//// [tests/cases/compiler/awaitInNamespaceExportedClassComputedProperty.ts] ////

//// [awaitInNamespaceExportedClassComputedProperty.ts]
declare const x: string;
namespace N {
    class A { [await x]() {} }
    export class B { [await x]() {} }
}
export {};


//// [awaitInNamespaceExportedClassComputedProperty.js]
var N;
(function (N) {
    class A {
        [await x]() { }
    }
    class B {
        [await x]() { }
    }
    N.B = B;
})(N || (N = {}));
export {};
