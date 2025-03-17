//// [tests/cases/compiler/declarationEmitClassSetAccessorParamNameInJs.ts] ////

//// [foo.js]
// https://github.com/microsoft/TypeScript/issues/55391

export class Foo {
    /**
     * Bar.
     *
     * @param {string} baz Baz.
     */
    set bar(baz) {}
}


//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
    set bar(baz) { }
}
exports.Foo = Foo;
