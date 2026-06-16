//// [tests/cases/compiler/jsDeclarationEmitBoundMethodInConstructor.ts] ////

//// [main.js]
export class C {
    constructor() {
        this.foo = this.foo.bind(this);
    }
    foo() {}
}




//// [main.d.ts]
export declare class C {
    constructor();
    foo(): void;
}
