//// [tests/cases/compiler/declarationEmitForMixinsWithStaticAccessors2.ts] ////

//// [a.ts]
function mix<T extends new (...args: any[]) => any, U extends new (...args: any[]) => any>(
    base1: T, base2: U
): T & U {
    return null as any;
}

class A {
    static get shared(): number { return 1; }
    static set shared(v: number) { }
    x: string = "";
}

class B {
    static get shared(): number { return 2; }
    static set shared(v: number) { }
    y: number = 0;
}

function make() {
    class C extends mix(A, B) {
        z: boolean = true;
    }
    return C;
}

export const MixedClass = make();


//// [a.js]
function mix(base1, base2) {
    return null;
}
class A {
    static get shared() { return 1; }
    static set shared(v) { }
    x = "";
}
class B {
    static get shared() { return 2; }
    static set shared(v) { }
    y = 0;
}
function make() {
    class C extends mix(A, B) {
        z = true;
    }
    return C;
}
export const MixedClass = make();


//// [a.d.ts]
export declare const MixedClass: {
    new (): {
        z: boolean;
        x: string;
    };
    new (): {
        z: boolean;
        x: string;
    };
    shared: number;
};
