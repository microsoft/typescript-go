//// [tests/cases/compiler/declarationEmitClassSetAccessorParamNameInJs2.ts] ////

//// [foo.js]
export class Foo {
    /**
     * Bar.
     *
     * @param {{ prop: string }} baz Baz.
     */
    set bar({}) {}
}


//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
    set bar({}) { }
}
exports.Foo = Foo;
