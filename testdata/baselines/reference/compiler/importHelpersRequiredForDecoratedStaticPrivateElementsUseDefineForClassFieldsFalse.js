//// [tests/cases/compiler/importHelpersRequiredForDecoratedStaticPrivateElementsUseDefineForClassFieldsFalse.ts] ////

//// [main.ts]
export declare var dec: any;

// At useDefineForClassFields: false, native decorators are still lowered even at esnext (this is the
// one case where the esnext decorator transform isn't skipped by target alone - see
// newESDecoratorTransformer's skip condition), so a decorated class's static private elements still
// need the classPrivateField* helpers here too.
@dec
export class Foo {
    #instanceField = 1;
    static #staticField = 1;

    getInstanceField() {
        return this.#instanceField;
    }
    static getStaticField() {
        return Foo.#staticField;
    }
}

//// [index.d.ts]
// Provides only the decorator helpers, deliberately omitting __classPrivateFieldGet, to confirm that
// accessing the decorated class's static private field still requires it.
export declare function __esDecorate(ctor: any, descriptorIn: any, decorators: any[], contextIn: any, initializers: any, extraInitializers: any): void;
export declare function __runInitializers(thisArg: any, initializers: any[], value?: any): any;
export declare function __setFunctionName(f: any, name: any, prefix?: string): any;


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
const tslib_1 = require("tslib");
// At useDefineForClassFields: false, native decorators are still lowered even at esnext (this is the
// one case where the esnext decorator transform isn't skipped by target alone - see
// newESDecoratorTransformer's skip condition), so a decorated class's static private elements still
// need the classPrivateField* helpers here too.
let Foo = (() => {
    var _Foo_staticField;
    let _classDecorators = [exports.dec];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    var Foo = class {
        static { _classThis = this; }
        static { tslib_1.__setFunctionName(this, "Foo"); }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            tslib_1.__esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            Foo = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        #instanceField = 1;
        static {
            _Foo_staticField = { value: 1 };
        }
        getInstanceField() {
            return this.#instanceField;
        }
        static getStaticField() {
            return tslib_1.__classPrivateFieldGet(Foo, _classThis, "f", _Foo_staticField);
        }
        static {
            tslib_1.__runInitializers(_classThis, _classExtraInitializers);
        }
    };
    return Foo = _classThis;
})();
exports.Foo = Foo;
