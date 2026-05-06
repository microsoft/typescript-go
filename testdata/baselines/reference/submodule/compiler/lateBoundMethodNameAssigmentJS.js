//// [tests/cases/compiler/lateBoundMethodNameAssigmentJS.ts] ////

//// [lateBoundMethodNameAssigmentJS.js]
const _symbol = Symbol("_sym");
export class MyClass {
    constructor() {
        this[_symbol] = this[_symbol].bind(this);
    }

    async [_symbol]() { }
}



//// [lateBoundMethodNameAssigmentJS.d.ts]
const _symbol: unique symbol;
export class MyClass {
    constructor();
    [_symbol](): Promise<void>;
}
export {};


//// [DtsFileErrors]


lateBoundMethodNameAssigmentJS.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== lateBoundMethodNameAssigmentJS.d.ts (1 errors) ====
    const _symbol: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export class MyClass {
        constructor();
        [_symbol](): Promise<void>;
    }
    export {};
    