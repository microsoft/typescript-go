//// [tests/cases/compiler/declarationEmitAccessorNilParent.ts] ////

//// [declarationEmitAccessorNilParent.ts]
// Regression test: when an accessor property symbol has a nil parent
// (e.g. from a union/intersection of classes with differing value declarations),
// the node builder should not crash with a nil pointer dereference.

class A {
    get x(): number { return 1; }
    set x(v: number) { }
}

class B {
    get x(): number { return 2; }
    set x(v: number) { }
}

// Force inlined structural type emit by using non-exported local classes
function make() {
    class C {
        get foo(): number { return 1; }
        set foo(v: number) { }
    }
    return new C();
}
export const val = make();


//// [declarationEmitAccessorNilParent.js]
// Regression test: when an accessor property symbol has a nil parent
// (e.g. from a union/intersection of classes with differing value declarations),
// the node builder should not crash with a nil pointer dereference.
class A {
    get x() { return 1; }
    set x(v) { }
}
class B {
    get x() { return 2; }
    set x(v) { }
}
// Force inlined structural type emit by using non-exported local classes
function make() {
    class C {
        get foo() { return 1; }
        set foo(v) { }
    }
    return new C();
}
export const val = make();


//// [declarationEmitAccessorNilParent.d.ts]
export declare const val: {
    get foo(): number;
    set foo(v: number);
};
