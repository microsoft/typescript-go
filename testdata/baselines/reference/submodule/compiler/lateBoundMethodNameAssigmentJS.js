//// [tests/cases/compiler/lateBoundMethodNameAssigmentJS.ts] ////

//// [lateBoundMethodNameAssigmentJS.js]
const _symbol = Symbol("_sym");
export class MyClass {
    constructor() {
        this[_symbol] = this[_symbol].bind(this);
    }

    async [_symbol]() { }
}

//// [lateBoundMethodNameAssigmentJS.js]
const _symbol = Symbol("_sym");
export class MyClass {
    constructor() {
        this[_symbol] = this[_symbol].bind(this);
    }
    async [_symbol]() { }
}
