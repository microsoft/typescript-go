//// [tests/cases/compiler/importHelpersRequiredForDecoratedStaticPrivateElements.ts] ////

//// [main.ts]
export declare var dec: any;

// Class decorators hoist a class's static private/auto-accessor elements out of the class body, so
// accessing them still requires the classPrivateField* helpers even at targets (like es2022) that
// otherwise support private fields natively. Instance private fields are unaffected by this and stay
// fully native even in a decorated class.
@dec
export class Foo {
    #instanceField = 1;
    static #staticField = 1;
    static #staticMethod() { return 1; }
    static get #staticAccessor() { return 1; }
    static set #staticAccessor(v: number) {}
    static accessor #staticAutoAccessor = 1;

    getInstanceField() {
        return this.#instanceField;
    }
    static getStaticField() {
        return Foo.#staticField;
    }
    static callStaticMethod() {
        return Foo.#staticMethod();
    }
    static useStaticAccessor() {
        Foo.#staticAccessor = Foo.#staticAccessor;
    }
    static useStaticAutoAccessor() {
        Foo.#staticAutoAccessor = Foo.#staticAutoAccessor;
    }
    static hasStaticField(x: object) {
        return #staticField in x;
    }
}

//// [index.d.ts]
// Provides only the decorator helpers, deliberately omitting __classPrivateField{Get,Set,In}, to
// confirm that accessing a decorated class's static private elements still requires them.
export declare function __esDecorate(ctor: any, descriptorIn: any, decorators: any[], contextIn: any, initializers: any, extraInitializers: any): void;
export declare function __runInitializers(thisArg: any, initializers: any[], value?: any): any;
export declare function __setFunctionName(f: any, name: any, prefix?: string): any;


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
const tslib_1 = require("tslib");
// Class decorators hoist a class's static private/auto-accessor elements out of the class body, so
// accessing them still requires the classPrivateField* helpers even at targets (like es2022) that
// otherwise support private fields natively. Instance private fields are unaffected by this and stay
// fully native even in a decorated class.
let Foo = (() => {
    var _Foo_staticField, _Foo_staticMethod, _Foo_staticAccessor_get, _Foo_staticAccessor_set, _Foo_staticAutoAccessor_get, _Foo_staticAutoAccessor_set, _Foo_staticAutoAccessor_accessor_storage;
    let _classDecorators = [exports.dec];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    var Foo = class {
        static { _classThis = this; }
        static { tslib_1.__setFunctionName(this, "Foo"); }
        static { _Foo_staticMethod = function _Foo_staticMethod() { return 1; }, _Foo_staticAccessor_get = function _Foo_staticAccessor_get() { return 1; }, _Foo_staticAccessor_set = function _Foo_staticAccessor_set(v) { }, _Foo_staticAutoAccessor_get = function _Foo_staticAutoAccessor_get() { return tslib_1.__classPrivateFieldGet(_classThis, _classThis, "f", _Foo_staticAutoAccessor_accessor_storage); }, _Foo_staticAutoAccessor_set = function _Foo_staticAutoAccessor_set(value) { tslib_1.__classPrivateFieldSet(_classThis, _classThis, value, "f", _Foo_staticAutoAccessor_accessor_storage); }; }
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
        static {
            _Foo_staticAutoAccessor_accessor_storage = { value: 1 };
        }
        getInstanceField() {
            return this.#instanceField;
        }
        static getStaticField() {
            return tslib_1.__classPrivateFieldGet(Foo, _classThis, "f", _Foo_staticField);
        }
        static callStaticMethod() {
            return tslib_1.__classPrivateFieldGet(Foo, _classThis, "m", _Foo_staticMethod).call(Foo);
        }
        static useStaticAccessor() {
            tslib_1.__classPrivateFieldSet(Foo, _classThis, tslib_1.__classPrivateFieldGet(Foo, _classThis, "a", _Foo_staticAccessor_get), "a", _Foo_staticAccessor_set);
        }
        static useStaticAutoAccessor() {
            tslib_1.__classPrivateFieldSet(Foo, _classThis, tslib_1.__classPrivateFieldGet(Foo, _classThis, "a", _Foo_staticAutoAccessor_get), "a", _Foo_staticAutoAccessor_set);
        }
        static hasStaticField(x) {
            return tslib_1.__classPrivateFieldIn(_classThis, x);
        }
        static {
            tslib_1.__runInitializers(_classThis, _classExtraInitializers);
        }
    };
    return Foo = _classThis;
})();
exports.Foo = Foo;
