//// [tests/cases/compiler/jsDeclarationEmitBoundMethodInConstructor.ts] ////

//// [main.js]
export class C {
    constructor() {
        this.foo = this.foo.bind(this);
    }
    foo() {}
}

export class D {
    constructor() {
        this.bar = 1;
    }
    static bar() {}
}

export class E {
    static init() {
        this.baz = 1;
    }
    baz() {}
}




//// [main.d.ts]
export declare class C {
    constructor();
    foo(): void;
}
export declare class D {
    bar: number;
    constructor();
    static bar(): void;
}
export declare class E {
    static baz: number | undefined;
    static init(): void;
    baz(): void;
}
