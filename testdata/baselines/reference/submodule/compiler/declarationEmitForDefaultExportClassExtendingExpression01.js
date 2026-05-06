//// [tests/cases/compiler/declarationEmitForDefaultExportClassExtendingExpression01.ts] ////

//// [declarationEmitForDefaultExportClassExtendingExpression01.ts]
interface Greeter {
    getGreeting(): string;
}

interface GreeterConstructor {
    new (): Greeter;
}

class A {
    getGreeting() {
        return 'hello';
    }
}

const getGreeterBase = (): GreeterConstructor => A;

export default class extends getGreeterBase() {
}



//// [declarationEmitForDefaultExportClassExtendingExpression01.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class A {
    getGreeting() {
        return 'hello';
    }
}
const getGreeterBase = () => A;
class default_1 extends getGreeterBase() {
}
exports.default = default_1;


//// [declarationEmitForDefaultExportClassExtendingExpression01.d.ts]
interface Greeter {
    getGreeting(): string;
}
interface GreeterConstructor {
    new (): Greeter;
}
const default_base: GreeterConstructor;
export default class extends default_base {
}
export {};


//// [DtsFileErrors]


declarationEmitForDefaultExportClassExtendingExpression01.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitForDefaultExportClassExtendingExpression01.d.ts (1 errors) ====
    interface Greeter {
        getGreeting(): string;
    }
    interface GreeterConstructor {
        new (): Greeter;
    }
    const default_base: GreeterConstructor;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default class extends default_base {
    }
    export {};
    