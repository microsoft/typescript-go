// @declaration: true
// @incremental: true
// @target: esnext
// @tsBuildInfoFile: /a.tsbuildinfo

// @Filename: /a.ts
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
