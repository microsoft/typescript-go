//// [tests/cases/compiler/awaitInNamespaceExportedClassComputedProperty.ts] ////

//// [awaitInNamespaceExportedClassComputedProperty.ts]
// Repro for https://github.com/microsoft/TypeScript/issues/63712
// await in computed property names inside a namespace's exported class should be an error

declare const x: string;
namespace N {
    class A { [await x]() {} }       // TS1308 (correct - always reported)
    export class B { [await x]() {} } // TS1308 (was missing - this is the bug)
}
export {};


//// [awaitInNamespaceExportedClassComputedProperty.js]
// Repro for https://github.com/microsoft/TypeScript/issues/63712
// await in computed property names inside a namespace's exported class should be an error
var N;
(function (N) {
    class A {
        [await x]() { }
    } // TS1308 (correct - always reported)
    class B {
        [await x]() { }
    } // TS1308 (was missing - this is the bug)
    N.B = B;
})(N || (N = {}));
export {};
