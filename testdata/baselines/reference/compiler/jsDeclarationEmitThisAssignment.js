//// [tests/cases/compiler/jsDeclarationEmitThisAssignment.ts] ////

//// [main.js]
export class Foo {
    static {
        this.bar = 10;
    }

    constructor() {
        this.baz = "hello";
    }
}

export class Bar {
    constructor() {
        this.x = 42;
        this.y = true;
    }
}




//// [main.d.ts]
export class Foo {
    static bar: number | undefined;
    baz: string;
    constructor();
}
export class Bar {
    x: number;
    y: boolean;
    constructor();
}
