// @target: esnext
// @module: esnext

declare const x: string;
namespace N {
    class A { [await x]() {} }
    export class B { [await x]() {} }
}
export class C { [await x]() {} }

{
    export class D { [await x]() {} }
}

function f() {
    export class E { [await x]() {} }
}

async function af() {
    export class F { [await x]() {} }
}

function* gf() {
    export class G { [await x]() {} }
}

async function* agf() {
    export class H { [await x]() {} }
}

export {};
