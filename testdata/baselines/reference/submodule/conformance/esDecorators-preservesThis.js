//// [tests/cases/conformance/esDecorators/esDecorators-preservesThis.ts] ////

//// [esDecorators-preservesThis.ts]
// https://github.com/microsoft/TypeScript/issues/53752

declare class DecoratorProvider {
    decorate<T>(this: DecoratorProvider, v: T, ctx: DecoratorContext): T;
}

declare const instance: DecoratorProvider;

// preserve `this` for access
class C {
    @instance.decorate
    method1() { }

    @(instance["decorate"])
    method2() { }

    // even in parens
    @((instance.decorate))
    method3() { }
}

// preserve `this` for `super` access
class D extends DecoratorProvider {
    m() {
        class C {
            @(super.decorate)
            method1() { }

            @(super["decorate"])
            method2() { }

            @((super.decorate))
            method3() { }
        }
    }
}


//// [esDecorators-preservesThis.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/53752
// preserve `this` for access
let C = (() => {
    let _instanceExtraInitializers = [];
    let _method1_decorators;
    let _method2_decorators;
    let _method3_decorators;
    return class C {
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _method1_decorators = [(_a = instance).decorate.bind(_a)];
            _method2_decorators = [((_b = instance)["decorate"].bind(_b))];
            _method3_decorators = [(((_c = instance).decorate.bind(_c)))];
            __esDecorate(this, null, _method1_decorators, { kind: "method", name: "method1", static: false, private: false, access: { has: obj => "method1" in obj, get: obj => obj.method1 }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(this, null, _method2_decorators, { kind: "method", name: "method2", static: false, private: false, access: { has: obj => "method2" in obj, get: obj => obj.method2 }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(this, null, _method3_decorators, { kind: "method", name: "method3", static: false, private: false, access: { has: obj => "method3" in obj, get: obj => obj.method3 }, metadata: _metadata }, null, _instanceExtraInitializers);
            if (_metadata) Object.defineProperty(this, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        method1() { }
        method2() { }
        // even in parens
        method3() { }
        constructor() {
            __runInitializers(this, _instanceExtraInitializers);
        }
    };
})();
// preserve `this` for `super` access
class D extends DecoratorProvider {
    m() {
        let C = (() => {
            let _outerThis = this;
            let _instanceExtraInitializers = [];
            let _method1_decorators;
            let _method2_decorators;
            let _method3_decorators;
            return class C {
                static {
                    const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
                    _method1_decorators = [(super.decorate.bind(_outerThis))];
                    _method2_decorators = [(super["decorate"].bind(_outerThis))];
                    _method3_decorators = [((super.decorate.bind(_outerThis)))];
                    __esDecorate(this, null, _method1_decorators, { kind: "method", name: "method1", static: false, private: false, access: { has: obj => "method1" in obj, get: obj => obj.method1 }, metadata: _metadata }, null, _instanceExtraInitializers);
                    __esDecorate(this, null, _method2_decorators, { kind: "method", name: "method2", static: false, private: false, access: { has: obj => "method2" in obj, get: obj => obj.method2 }, metadata: _metadata }, null, _instanceExtraInitializers);
                    __esDecorate(this, null, _method3_decorators, { kind: "method", name: "method3", static: false, private: false, access: { has: obj => "method3" in obj, get: obj => obj.method3 }, metadata: _metadata }, null, _instanceExtraInitializers);
                    if (_metadata) Object.defineProperty(this, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
                }
                method1() { }
                method2() { }
                method3() { }
                constructor() {
                    __runInitializers(this, _instanceExtraInitializers);
                }
            };
        })();
    }
}
