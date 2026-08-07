//// [tests/cases/compiler/importHelpersNotRequiredForDecoratedStaticPrivateElementsAtESNext.ts] ////

//// [main.ts]
export declare var dec: any;

// At target: esnext with the default useDefineForClassFields: true, native decorators need no
// transform at all (see newESDecoratorTransformer's skip condition), so even a decorated class's
// static private elements stay fully native and never need the classPrivateField* helpers.
@dec
export class Foo {
    static #staticField = 1;
    static getStaticField() {
        return Foo.#staticField;
    }
}


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
// At target: esnext with the default useDefineForClassFields: true, native decorators need no
// transform at all (see newESDecoratorTransformer's skip condition), so even a decorated class's
// static private elements stay fully native and never need the classPrivateField* helpers.
@exports.dec
class Foo {
    static #staticField = 1;
    static getStaticField() {
        return Foo.#staticField;
    }
}
exports.Foo = Foo;
