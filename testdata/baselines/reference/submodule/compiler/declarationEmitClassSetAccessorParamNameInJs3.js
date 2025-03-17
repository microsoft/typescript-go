//// [tests/cases/compiler/declarationEmitClassSetAccessorParamNameInJs3.ts] ////

//// [foo.js]
export class Foo {
    /**
     * Bar.
     *
     * @param {{ prop: string | undefined }} baz Baz.
     */
    set bar({ prop = 'foo' }) {}
}


//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
    set bar({ prop = 'foo' }) { }
}
exports.Foo = Foo;
