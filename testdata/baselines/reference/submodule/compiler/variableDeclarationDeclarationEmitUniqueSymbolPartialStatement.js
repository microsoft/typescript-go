//// [tests/cases/compiler/variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.ts] ////

//// [variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.ts]
const key = Symbol(), value = 12;

export class Foo {
    [key] = value;
}

//// [variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.js]
"use strict";
var _a;
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
const key = Symbol(), value = 12;
class Foo {
    constructor() {
        this[_a] = value;
    }
}
exports.Foo = Foo;
_a = key;


//// [variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.d.ts]
const key: unique symbol;
export class Foo {
    [key]: number;
}
export {};


//// [DtsFileErrors]


variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== variableDeclarationDeclarationEmitUniqueSymbolPartialStatement.d.ts (1 errors) ====
    const key: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export class Foo {
        [key]: number;
    }
    export {};
    