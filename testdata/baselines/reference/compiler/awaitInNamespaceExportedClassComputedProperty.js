//// [tests/cases/compiler/awaitInNamespaceExportedClassComputedProperty.ts] ////

//// [awaitInNamespaceExportedClassComputedProperty.ts]
declare const x: string;
namespace N {
    class A { [await x]() {} }
    export class B { [await x]() {} }
}
export class C { [await x]() {} }

{
    class D { [await x]() {} }
}

function f() {
    class E { [await x]() {} }
}

async function af() {
    class F { [await x]() {} }
}

function* gf() {
    class G { [await x]() {} }
}

async function* agf() {
    class H { [await x]() {} }
}

function switchSync() {
    switch (0) {
        case 0:
            class I { [await x]() {} }
    }
}

async function switchAsync() {
    switch (0) {
        case 0:
            class J { [await x]() {} }
    }
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
export class C {
    [await x]() { }
}
{
    class D {
        [await x]() { }
    }
}
function f() {
    class E {
        [await x]() { }
    }
}
async function af() {
    class F {
        [await x]() { }
    }
}
function* gf() {
    class G {
        [await x]() { }
    }
}
async function* agf() {
    class H {
        [await x]() { }
    }
}
function switchSync() {
    switch (0) {
        case 0:
            class I {
                [await x]() { }
            }
    }
}
async function switchAsync() {
    switch (0) {
        case 0:
            class J {
                [await x]() { }
            }
    }
}
