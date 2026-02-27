//// [tests/cases/conformance/esDecorators/metadata/esDecoratorsMetadata2.ts] ////

//// [foo.ts]
function meta(key: string, value: string) {
    return (_, context) => {
        context.metadata[key] = value;
    };
}

@meta('a', 'x')
class C {
    @meta('b', 'y')
    m() {}
}

C[Symbol.metadata].a; // 'x'
C[Symbol.metadata].b; // 'y'

class D extends C {
    @meta('b', 'z')
    m() {}
}

D[Symbol.metadata].a; // 'x'
D[Symbol.metadata].b; // 'z'


//// [foo.js]
"use strict";
function meta(key, value) {
    return (_, context) => {
        context.metadata[key] = value;
    };
}
let C = (() => {
    let _classDecorators = [meta('a', 'x')];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    let _instanceExtraInitializers = [];
    let _m_decorators;
    var C = class {
        static { _classThis = this; }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _m_decorators = [meta('b', 'y')];
            __esDecorate(this, null, _m_decorators, { kind: "method", name: "m", static: false, private: false, access: { has: obj => "m" in obj, get: obj => obj.m }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            C = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
            __runInitializers(_classThis, _classExtraInitializers);
        }
        m() { }
        constructor() {
            __runInitializers(this, _instanceExtraInitializers);
        }
    };
    return C = _classThis;
})();
C[Symbol.metadata].a; // 'x'
C[Symbol.metadata].b; // 'y'
let D = (() => {
    let _classSuper = C;
    let _instanceExtraInitializers = [];
    let _m_decorators;
    return class D extends _classSuper {
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(_classSuper[Symbol.metadata] ?? null) : void 0;
            _m_decorators = [meta('b', 'z')];
            __esDecorate(this, null, _m_decorators, { kind: "method", name: "m", static: false, private: false, access: { has: obj => "m" in obj, get: obj => obj.m }, metadata: _metadata }, null, _instanceExtraInitializers);
            if (_metadata) Object.defineProperty(this, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        m() { }
        constructor() {
            super(...arguments);
            __runInitializers(this, _instanceExtraInitializers);
        }
    };
})();
D[Symbol.metadata].a; // 'x'
D[Symbol.metadata].b; // 'z'
